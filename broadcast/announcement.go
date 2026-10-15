package broadcast

import (
	"errors"
	"fmt"
	"time"
)

const LegalTechChannel = "legal-tech-members"

type MatterStage string

const (
	StageIntake         MatterStage = "matter_intake"
	StageSignedDelivery MatterStage = "signed_document_delivery"
	StageDeadline       MatterStage = "deadline_follow_up"
)

type Announcement struct {
	MatterID  string      `json:"matter_id"`
	Stage     MatterStage `json:"stage"`
	Summary   string      `json:"summary"`
	Document  string      `json:"document_id,omitempty"`
	Deadline  string      `json:"deadline,omitempty"`
	AccountID string      `json:"account_id,omitempty"`
}

type Dispatch struct {
	Channel string       `json:"channel"`
	Event   string       `json:"event"`
	Data    Announcement `json:"data"`
}

func Plan(a Announcement) (Dispatch, error) {
	if a.MatterID == "" || a.Summary == "" {
		return Dispatch{}, errors.New("matter_id and summary are required")
	}

	var event string
	switch a.Stage {
	case StageIntake:
		event = "matter.intake.received"
	case StageSignedDelivery:
		if a.Document == "" {
			return Dispatch{}, errors.New("document_id is required for signed document delivery")
		}
		event = "document.signed.delivered"
	case StageDeadline:
		if a.Deadline == "" {
			return Dispatch{}, errors.New("deadline is required for deadline follow-up")
		}
		if _, err := time.Parse(time.RFC3339, a.Deadline); err != nil {
			return Dispatch{}, fmt.Errorf("deadline must be RFC3339: %w", err)
		}
		event = "matter.deadline.follow_up"
	default:
		return Dispatch{}, fmt.Errorf("unsupported stage %q", a.Stage)
	}

	return Dispatch{Channel: LegalTechChannel, Event: event, Data: a}, nil
}
