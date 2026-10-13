package collectors

import (
	"sync"
	"time"
)

// deviceFailureState tracks the failure and notification history for a single device
// across daemon run intervals.
type deviceFailureState struct {
	// consecutiveFailures is the number of failed runs since the last success.
	consecutiveFailures int
	// notifyCount is the number of alerts sent during the current failure episode.
	notifyCount int
	// lastNotified is when the most recent alert was sent.
	lastNotified time.Time
}

// failureTracker holds failure state for all devices, keyed by device name.
type failureTracker struct {
	mu     sync.Mutex
	states map[string]*deviceFailureState
}

// globalFailureTracker persists device failure state for the lifetime of the process.
// It lives at package scope on purpose: CollectorWorker is recreated on every run, and
// the config that carries per-device settings is rebuilt on every daemon interval by
// ReloadConfig. State kept on the (reloaded) DeviceConfig would reset each interval, so
// the failure counter has to live somewhere that outlives both.
var globalFailureTracker = &failureTracker{
	states: make(map[string]*deviceFailureState),
}

// recordSuccess clears any failure state for the named device, ending its failure episode.
func (t *failureTracker) recordSuccess(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.states, name)
}

// recordFailure registers a failed collection for the device and reports whether an alert
// should be sent for this failure.
//
// warning is the number of consecutive failed runs before the first alert (0 disables
// alerts entirely). Until backoffCount alerts have been sent in the current episode, an
// alert fires on every warning failed runs. After that, backoff engages and further
// alerts are rate limited to at most one per backoffInterval. now is supplied by the
// caller so the clock can be controlled in tests.
func (t *failureTracker) recordFailure(name string, warning, backoffCount int, backoffInterval time.Duration, now time.Time) bool {
	// A zero (or negative) warning threshold disables alerts for this device.
	if warning <= 0 {
		return false
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	state, ok := t.states[name]
	if !ok {
		state = &deviceFailureState{}
		t.states[name] = state
	}
	state.consecutiveFailures++

	// Below the initial threshold, stay quiet.
	if state.consecutiveFailures < warning {
		return false
	}

	notify := false
	if state.notifyCount < backoffCount {
		// Before backoff: alert once every `warning` consecutive failed runs. Each run
		// increments consecutiveFailures by exactly one, so a modulo on the overage gives
		// a steady count-based cadence (first alert lands exactly at `warning`).
		if (state.consecutiveFailures-warning)%warning == 0 {
			notify = true
		}
	} else {
		// Backoff engaged: alert at most once per interval.
		if now.Sub(state.lastNotified) >= backoffInterval {
			notify = true
		}
	}

	if notify {
		state.notifyCount++
		state.lastNotified = now
	}
	return notify
}
