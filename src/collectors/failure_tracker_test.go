package collectors

import (
	"testing"
	"time"
)

// newTestTracker returns an isolated tracker so tests don't share the package global.
func newTestTracker() *failureTracker {
	return &failureTracker{states: make(map[string]*deviceFailureState)}
}

// TestRecordFailureDisabled verifies a zero warning threshold never alerts.
func TestRecordFailureDisabled(t *testing.T) {
	tr := newTestTracker()
	now := time.Unix(0, 0)
	for i := 0; i < 10; i++ {
		if tr.recordFailure("dev", 0, 5, 24*time.Hour, now) {
			t.Fatalf("expected no alert with warning=0, got one on failure %d", i+1)
		}
	}
}

// TestRecordFailureInitialThreshold verifies the first alert lands exactly at `warning`
// failures and then repeats every `warning` failures before backoff.
func TestRecordFailureInitialThreshold(t *testing.T) {
	tr := newTestTracker()
	now := time.Unix(0, 0)
	warning := 5
	var alerts []int
	for i := 1; i <= 15; i++ {
		if tr.recordFailure("dev", warning, 5, 24*time.Hour, now) {
			alerts = append(alerts, i)
		}
	}
	want := []int{5, 10, 15}
	if len(alerts) != len(want) {
		t.Fatalf("expected alerts at %v, got %v", want, alerts)
	}
	for i := range want {
		if alerts[i] != want[i] {
			t.Fatalf("expected alerts at %v, got %v", want, alerts)
		}
	}
}

// TestRecordFailureBackoff verifies that once backoffCount alerts have fired, further
// alerts are rate limited to one per backoffInterval regardless of failure cadence.
func TestRecordFailureBackoff(t *testing.T) {
	tr := newTestTracker()
	warning := 1        // alert on every failure until backoff
	backoffCount := 3   // engage backoff after 3 alerts
	interval := time.Hour
	base := time.Unix(0, 0)

	// First three failures (times don't matter yet) should each alert.
	for i := 0; i < 3; i++ {
		if !tr.recordFailure("dev", warning, backoffCount, interval, base) {
			t.Fatalf("expected alert on pre-backoff failure %d", i+1)
		}
	}

	// Backoff is now engaged. A failure at the same instant must be suppressed.
	if tr.recordFailure("dev", warning, backoffCount, interval, base) {
		t.Fatal("expected backoff to suppress alert within the interval")
	}
	// Still inside the interval.
	if tr.recordFailure("dev", warning, backoffCount, interval, base.Add(30*time.Minute)) {
		t.Fatal("expected backoff to suppress alert 30m in")
	}
	// Past the interval - one alert allowed.
	if !tr.recordFailure("dev", warning, backoffCount, interval, base.Add(time.Hour)) {
		t.Fatal("expected alert once the backoff interval elapsed")
	}
	// Immediately after, suppressed again.
	if tr.recordFailure("dev", warning, backoffCount, interval, base.Add(time.Hour+time.Minute)) {
		t.Fatal("expected suppression right after a backoff alert")
	}
}

// TestRecordSuccessResets verifies a success clears the episode so counting restarts.
func TestRecordSuccessResets(t *testing.T) {
	tr := newTestTracker()
	now := time.Unix(0, 0)
	warning := 3

	// Two failures, below threshold, no alert yet.
	tr.recordFailure("dev", warning, 5, 24*time.Hour, now)
	tr.recordFailure("dev", warning, 5, 24*time.Hour, now)

	// Success resets the counter.
	tr.recordSuccess("dev")

	// It should now take a full `warning` failures again to alert.
	if tr.recordFailure("dev", warning, 5, 24*time.Hour, now) {
		t.Fatal("expected no alert on first failure after reset")
	}
	if tr.recordFailure("dev", warning, 5, 24*time.Hour, now) {
		t.Fatal("expected no alert on second failure after reset")
	}
	if !tr.recordFailure("dev", warning, 5, 24*time.Hour, now) {
		t.Fatal("expected alert on the third failure after reset")
	}
}
