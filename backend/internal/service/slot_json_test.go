package service

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestSlotJSONShapeMatchesOpenAPI locks the wire format of Slot to the
// openapi.yaml Appointment schema. The Flutter client parses this response
// strictly (json['_id'] as String, DateTime.parse(json['startTime'])), so
// renaming a Go field or dropping a tag silently breaks the mobile app.
func TestSlotJSONShapeMatchesOpenAPI(t *testing.T) {
	slot := Slot{
		ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		TherapistID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		StartTs:     time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC).Format(time.RFC3339),
		EndTs:       time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC).Format(time.RFC3339),
		Status:      "open",
	}

	raw, err := json.Marshal(slot)
	if err != nil {
		t.Fatalf("marshal Slot: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal Slot: %v", err)
	}

	// Keys the spec (and the Dart model) require.
	for _, key := range []string{"_id", "startTime", "endTime", "status"} {
		if _, ok := got[key]; !ok {
			t.Errorf("Slot JSON missing required key %q; got %s", key, raw)
		}
	}

	// Go field names must never leak: the Dart client would see null and throw.
	for _, key := range []string{"ID", "TherapistID", "StartTs", "EndTs"} {
		if _, ok := got[key]; ok {
			t.Errorf("Go-style key %q leaked into the wire format; got %s", key, raw)
		}
	}

	// The Dart model does DateTime.parse on these, so they must be RFC3339.
	for _, key := range []string{"startTime", "endTime"} {
		value, ok := got[key].(string)
		if !ok {
			t.Errorf("%q is not a string, got %T", key, got[key])
			continue
		}
		if _, err := time.Parse(time.RFC3339, value); err != nil {
			t.Errorf("%q is not RFC3339 (%v), DateTime.parse would throw", key, err)
		}
	}
}
