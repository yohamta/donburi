package ecs

import "time"

const maxDaltaTime = time.Second / 30

// Time manages the time of the world.
type Time struct {
	// timeScale is the scale of the time.
	timeScale float64
	// deltaTime is the time between the last update and the current update
	deltaTime time.Duration
	// sleep is the time to sleep
	sleep time.Duration

	prevTime time.Time
	isPaused bool
}

// NewTIme creates a new Time.
func NewTime() *Time {
	return &Time{
		prevTime:  time.Now(),
		timeScale: 1,
	}
}

// Update updates the time.
func (t *Time) Update() {
	now := time.Now()
	if t.isPaused {
		t.prevTime = now
	}

	t.deltaTime = (now.Sub(t.prevTime))

	if t.timeScale != 1 {
		ms := float64(t.deltaTime.Milliseconds())
		ms *= t.timeScale
		t.deltaTime = time.Duration(ms) * time.Millisecond
	}

	if t.deltaTime > maxDaltaTime {
		t.deltaTime = maxDaltaTime
	}

	if t.sleep > 0 {
		d := t.deltaTime
		t.deltaTime -= t.sleep
		if t.deltaTime < 0 {
			t.deltaTime = 0
		}
		t.sleep -= d
		if t.sleep < 0 {
			t.sleep = 0
		}
	}

	t.prevTime = now
}

// SetSleep sets the time to sleep.
func (t *Time) SetSleep(d time.Duration) {
	t.sleep = d
}

// SetTimeScale sets the time scale.
func (t *Time) SetTimeScale(scale float64) {
	t.timeScale = scale
}

// TimeScale returns the time scale.
func (t *Time) TimeScale() float64 {
	return t.timeScale
}

// DeltaTime returns the time between the last update and the current update.
func (t *Time) DeltaTime() time.Duration {
	return t.deltaTime
}

// Pause pauses the time.
func (t *Time) Pause() {
	t.isPaused = true
}

// Resume resumes the time.
func (t *Time) Resume() {
	t.isPaused = false
}
