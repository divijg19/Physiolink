package clock

import (
	"testing"
	"time"
)

func TestRealClock_Now(t *testing.T) {
	c := NewReal()
	now := c.Now()
	if now.IsZero() {
		t.Fatal("expected non-zero time")
	}
}

func TestFakeClock_Now(t *testing.T) {
	tm := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	c := NewFake(tm)
	if !c.Now().Equal(tm) {
		t.Fatalf("expected %v, got %v", tm, c.Now())
	}
}

func TestFakeClock_Set(t *testing.T) {
	c := NewFake(time.Time{})
	tm := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	c.Set(tm)
	if !c.Now().Equal(tm) {
		t.Fatalf("expected %v, got %v", tm, c.Now())
	}
}
