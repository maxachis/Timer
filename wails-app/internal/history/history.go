// Package history keeps a bounded, chronological log of what happened to the
// user's timers — created, started, paused, resumed, reset, finished, removed —
// and derives those events by diffing timer state between sync points, so the
// callers that mutate timers never have to say which event they caused.
package history

import "time"

// Event kinds, as serialised to the frontend and to disk.
const (
	KindCreated  = "created"
	KindStarted  = "started"
	KindPaused   = "paused"
	KindResumed  = "resumed"
	KindReset    = "reset"
	KindFinished = "finished"
	KindRemoved  = "removed"
)

// Timer run states as reported by timer.CountdownTimer.StateName. Mirrored
// here so this package stays independent of the timer package.
const (
	StateIdle     = "idle"
	StateRunning  = "running"
	StatePaused   = "paused"
	StateFinished = "finished"
)

// MaxEntries bounds the log: old activity falls off the front so the
// persisted file and the modal's scroll list stay a manageable size.
const MaxEntries = 500

type Entry struct {
	At        time.Time `json:"at"`
	Kind      string    `json:"kind"`
	TimerName string    `json:"timer_name"`
}

// Log is an append-only, capped list of entries, oldest first.
type Log struct {
	entries []Entry
}

func NewLog() *Log { return &Log{} }

// FromEntries rehydrates a log from persisted entries, keeping the newest
// MaxEntries if the stored list is longer.
func FromEntries(entries []Entry) *Log {
	if len(entries) > MaxEntries {
		entries = entries[len(entries)-MaxEntries:]
	}
	l := &Log{entries: make([]Entry, len(entries))}
	copy(l.entries, entries)
	return l
}

func (l *Log) Record(kind, timerName string, at time.Time) {
	l.entries = append(l.entries, Entry{At: at, Kind: kind, TimerName: timerName})
	if len(l.entries) > MaxEntries {
		l.entries = l.entries[len(l.entries)-MaxEntries:]
	}
}

// Entries returns a copy, oldest first.
func (l *Log) Entries() []Entry {
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Observed is one timer as seen at a sync point. K identifies the timer
// across syncs (the app uses the timer's pointer); indices won't do because
// removing a timer shifts everything after it.
type Observed[K comparable] struct {
	ID    K
	Name  string
	State string
}

type seen struct {
	name  string
	state string
}

// Tracker turns successive snapshots of the timer set into log entries.
type Tracker[K comparable] struct {
	log  *Log
	last map[K]seen
}

func NewTracker[K comparable](log *Log) *Tracker[K] {
	return &Tracker[K]{log: log, last: map[K]seen{}}
}

// Seed records the current timer set as the baseline without logging, for
// timers restored from disk at startup.
func (t *Tracker[K]) Seed(timers []Observed[K]) {
	t.last = map[K]seen{}
	for _, o := range timers {
		t.last[o.ID] = seen{name: o.Name, state: o.State}
	}
}

// Rekey tells the tracker that the timer known as old now lives under id,
// so replacing a timer's underlying object isn't logged as removed+created.
func (t *Tracker[K]) Rekey(old, id K) {
	if s, ok := t.last[old]; ok {
		delete(t.last, old)
		t.last[id] = s
	}
}

// Sync diffs the given timer set against the last one and logs every
// difference at the given time. Returns true if anything was logged.
func (t *Tracker[K]) Sync(now time.Time, timers []Observed[K]) bool {
	logged := false
	present := make(map[K]bool, len(timers))
	for _, o := range timers {
		present[o.ID] = true
		prev, known := t.last[o.ID]
		if !known {
			t.log.Record(KindCreated, o.Name, now)
			logged = true
		} else if kind := transitionKind(prev.state, o.State); kind != "" {
			t.log.Record(kind, o.Name, now)
			logged = true
		}
		t.last[o.ID] = seen{name: o.Name, state: o.State}
	}
	for id, prev := range t.last {
		if !present[id] {
			t.log.Record(KindRemoved, prev.name, now)
			delete(t.last, id)
			logged = true
		}
	}
	return logged
}

// transitionKind names the event a run-state change represents, or "" when
// the state is unchanged.
func transitionKind(from, to string) string {
	if from == to {
		return ""
	}
	switch to {
	case StateRunning:
		if from == StateIdle {
			return KindStarted
		}
		return KindResumed
	case StatePaused:
		return KindPaused
	case StateIdle:
		return KindReset
	case StateFinished:
		return KindFinished
	}
	return ""
}
