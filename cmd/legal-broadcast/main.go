package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/infrai-examples/legaltech-broadcast/broadcast"
)

type server struct {
	client *broadcast.Client
}

func main() {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		log.Fatal("INFRAI_API_KEY is required")
	}

	s := &server{client: broadcast.NewClient(key)}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /announcements", s.announce)

	addr := envOr("ADDR", ":8080")
	httpServer := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Printf("legal broadcast listening on %s", addr)
	log.Fatal(httpServer.ListenAndServe())
}

func (s *server) announce(w http.ResponseWriter, r *http.Request) {
	var input broadcast.Announcement
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	dispatch, err := broadcast.Plan(input)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	requestID := fmt.Sprintf("matter-%s-%s", input.MatterID, input.Stage)
	if err := s.client.Publish(r.Context(), dispatch, requestID); err != nil {
		var apiErr *broadcast.APIError
		if errors.As(err, &apiErr) && apiErr.HTTPStatus >= 400 && apiErr.HTTPStatus < 500 {
			writeJSON(w, apiErr.HTTPStatus, map[string]string{"error": apiErr.Message, "code": apiErr.Code})
			return
		}
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "announcement was not accepted"})
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"matter_id": input.MatterID,
		"channel":   dispatch.Channel,
		"event":     dispatch.Event,
		"status":    "published",
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
