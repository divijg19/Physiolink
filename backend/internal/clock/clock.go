package clock

import "time"

// Clock is a simple interface for getting the current time.
type Clock interface {
	Now() time.Time
}

// RealClock implements Clock using the real time package.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now().UTC() }

// NewReal returns a RealClock instance.
func NewReal() Clock { return RealClock{} }

// FakeClock is useful for tests.
type FakeClock struct {
	t time.Time
}

func NewFake(t time.Time) *FakeClock { return &FakeClock{t: t} }

func (f *FakeClock) Now() time.Time { return f.t }

func (f *FakeClock) Set(t time.Time) { f.t = t }
