package broadcast

import "testing"

func TestPlanSelectsLegalWorkflowEvent(t *testing.T) {
	tests := []struct {
		name      string
		input     Announcement
		wantEvent string
		wantError bool
	}{
		{name: "intake", input: Announcement{MatterID: "MAT-42", Stage: StageIntake, Summary: "New conflicts review"}, wantEvent: "matter.intake.received"},
		{name: "signed delivery", input: Announcement{MatterID: "MAT-42", Stage: StageSignedDelivery, Summary: "Engagement letter signed", Document: "DOC-9"}, wantEvent: "document.signed.delivered"},
		{name: "deadline", input: Announcement{MatterID: "MAT-42", Stage: StageDeadline, Summary: "Response due", Deadline: "2026-10-05T09:00:00Z"}, wantEvent: "matter.deadline.follow_up"},
		{name: "signed delivery needs document", input: Announcement{MatterID: "MAT-42", Stage: StageSignedDelivery, Summary: "Engagement letter signed"}, wantError: true},
		{name: "deadline needs timestamp", input: Announcement{MatterID: "MAT-42", Stage: StageDeadline, Summary: "Response due", Deadline: "next Friday"}, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Plan(tt.input)
			if (err != nil) != tt.wantError {
				t.Fatalf("Plan() error = %v, wantError %v", err, tt.wantError)
			}
			if err == nil && (got.Channel != LegalTechChannel || got.Event != tt.wantEvent) {
				t.Fatalf("Plan() = channel %q event %q", got.Channel, got.Event)
			}
		})
	}
}
