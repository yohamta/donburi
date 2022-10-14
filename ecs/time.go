package ecs

import "time"

const maxDaltaTime = time.Second / 30

// Time manages the time of the world.
type Time struct {
	// IsPaused is true if the world is paused.
	IsPaused bool
	// TimeScale is the scale of the time.
	TimeScale float64
	// DeltaTime is the time between the last update and the current update
	DeltaTime time.Duration
	// Sleep is the time to sleep
	Sleep time.Duration

	prevTime time.Time
}

// NewTIme creates a new Time.
func NewTime() *Time {
	return &Time{
		prevTime: time.Now(),
	}
}

// Update updates the time.
func (t *Time) Update() {
	now := time.Now()
	if t.IsPaused {
		t.prevTime = now
	}

	t.DeltaTime = (now.Sub(t.prevTime))

	if t.TimeScale != 1 {
		ms := float64(t.DeltaTime.Milliseconds())
		ms *= t.TimeScale
		t.DeltaTime = time.Duration(ms) * time.Millisecond
	}

	if t.DeltaTime > maxDaltaTime {
		t.DeltaTime = maxDaltaTime
	}

	if t.Sleep > 0 {
		d := t.DeltaTime
		t.DeltaTime -= t.Sleep
		if t.DeltaTime < 0 {
			t.DeltaTime = 0
		}
		t.Sleep -= d
		if t.Sleep < 0 {
			t.Sleep = 0
		}
	}

	t.prevTime = now
}
