package broadcast

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPublishRetriesRateLimitWithSameIdempotencyKey(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodPost {
			t.Errorf("method = %s", r.Method)
		}
		if got := r.Header.Get("Idempotency-Key"); got != "matter-MAT-42-intake" {
			t.Errorf("idempotency key = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		if requests == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"ok":false,"data":null,"error":{"code":"RATE_LIMITED","message":"retry later"},"metadata":null}`))
			return
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"published":true},"error":null,"metadata":null}`))
	}))
	defer server.Close()

	client := NewClient("test-key")
	client.baseURL = server.URL
	client.sleep = func(context.Context, time.Duration) error { return nil }
	dispatch, err := Plan(Announcement{MatterID: "MAT-42", Stage: StageIntake, Summary: "New conflicts review"})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Publish(context.Background(), dispatch, "matter-MAT-42-intake"); err != nil {
		t.Fatal(err)
	}
	if requests != 2 {
		t.Fatalf("requests = %d, want 2", requests)
	}
}
