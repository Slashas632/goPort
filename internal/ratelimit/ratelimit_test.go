package ratelimit

import "testing"

func TestInit_SetsLimiter(t *testing.T) {
	Init(500)
	if Limiter == nil {
		t.Fatal("expected Limiter to be non-nil after Init")
	}
	if got := Limiter.Burst(); got != 250 {
		t.Errorf("Burst() = %d, want 250 for 500 workers", got)
	}
}

// Regression test: workers/2 truncates to 0 for workers in {0, 1}, and a
// rate.Limiter with burst 0 rejects every Wait() call immediately instead
// of throttling — silently disabling rate limiting. Init must floor the
// burst at 1.
func TestInit_BurstNeverZero(t *testing.T) {
	for _, workers := range []int{0, 1, -3} {
		Init(workers)
		if Limiter == nil {
			t.Fatalf("Init(%d): Limiter is nil", workers)
		}
		if got := Limiter.Burst(); got < 1 {
			t.Errorf("Init(%d): Burst() = %d, want >= 1", workers, got)
		}
	}
}

func TestInit_OddWorkerCountRoundsDown(t *testing.T) {
	Init(7)
	if got := Limiter.Burst(); got != 3 {
		t.Errorf("Burst() = %d, want 3 for 7 workers (7/2 truncated)", got)
	}
}
