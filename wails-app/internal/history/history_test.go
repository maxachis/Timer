package history

import (
	"testing"
	"time"
)

var t0 = time.Date(2026, 9, 18, 14, 0, 0, 0, time.UTC)

func kinds(entries []Entry) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Kind
	}
	return out
}

func assertKinds(t *testing.T, got []Entry, want ...string) {
	t.Helper()
	g := kinds(got)
	if len(g) != len(want) {
		t.Fatalf("kinds = %v, want %v", g, want)
	}
	for i := range want {
		if g[i] != want[i] {
			t.Fatalf("kinds = %v, want %v", g, want)
		}
	}
}

func TestLogRecordsOldestFirst(t *testing.T) {
	l := NewLog()
	l.Record(KindCreated, "A", t0)
	l.Record(KindStarted, "A", t0.Add(time.Second))

	got := l.Entries()
	assertKinds(t, got, KindCreated, KindStarted)
	if got[0].TimerName != "A" || !got[0].At.Equal(t0) {
		t.Fatalf("first entry = %+v", got[0])
	}
}

func TestLogDropsOldestPastCap(t *testing.T) {
	l := NewLog()
	for i := 0; i < MaxEntries+3; i++ {
		l.Record(KindStarted, "A", t0.Add(time.Duration(i)*time.Second))
	}
	got := l.Entries()
	if len(got) != MaxEntries {
		t.Fatalf("len = %d, want %d", len(got), MaxEntries)
	}
	if !got[0].At.Equal(t0.Add(3 * time.Second)) {
		t.Fatalf("oldest surviving entry at %v, want %v", got[0].At, t0.Add(3*time.Second))
	}
}

func TestLogEntriesIsACopy(t *testing.T) {
	l := NewLog()
	l.Record(KindCreated, "A", t0)
	got := l.Entries()
	got[0].TimerName = "mutated"
	if l.Entries()[0].TimerName != "A" {
		t.Fatal("Entries() exposed internal slice")
	}
}

func TestFromEntriesRestoresAndKeepsCap(t *testing.T) {
	in := make([]Entry, MaxEntries+2)
	for i := range in {
		in[i] = Entry{At: t0.Add(time.Duration(i) * time.Second), Kind: KindStarted, TimerName: "A"}
	}
	l := FromEntries(in)
	got := l.Entries()
	if len(got) != MaxEntries {
		t.Fatalf("len = %d, want %d", len(got), MaxEntries)
	}
	if !got[len(got)-1].At.Equal(in[len(in)-1].At) {
		t.Fatal("newest entry was dropped instead of oldest")
	}
}

// Tracker tests key timers by a plain int stand-in for the timer pointer the
// app uses; the tracker only needs identity.

func obs(id int, name, state string) Observed[int] {
	return Observed[int]{ID: id, Name: name, State: state}
}

func newTracker() (*Tracker[int], *Log) {
	l := NewLog()
	return NewTracker[int](l), l
}

func TestSeedDoesNotLogExistingTimers(t *testing.T) {
	tr, l := newTracker()
	tr.Seed([]Observed[int]{obs(1, "A", StateIdle), obs(2, "B", StatePaused)})
	if n := len(l.Entries()); n != 0 {
		t.Fatalf("seed logged %d entries", n)
	}
}

func TestSyncLogsCreatedAndRemoved(t *testing.T) {
	tr, l := newTracker()
	tr.Seed([]Observed[int]{obs(1, "A", StateIdle)})

	tr.Sync(t0, []Observed[int]{obs(1, "A", StateIdle), obs(2, "B", StateIdle)})
	assertKinds(t, l.Entries(), KindCreated)
	if l.Entries()[0].TimerName != "B" {
		t.Fatalf("created entry names %q", l.Entries()[0].TimerName)
	}

	tr.Sync(t0.Add(time.Second), []Observed[int]{obs(2, "B", StateIdle)})
	assertKinds(t, l.Entries(), KindCreated, KindRemoved)
	if got := l.Entries()[1]; got.TimerName != "A" || !got.At.Equal(t0.Add(time.Second)) {
		t.Fatalf("removed entry = %+v", got)
	}
}

func TestSyncLogsRunStateTransitions(t *testing.T) {
	tr, l := newTracker()
	tr.Seed([]Observed[int]{obs(1, "A", StateIdle)})

	steps := []struct {
		state string
		want  string
	}{
		{StateRunning, KindStarted},
		{StatePaused, KindPaused},
		{StateRunning, KindResumed},
		{StateFinished, KindFinished},
		{StateIdle, KindReset},
	}
	var want []string
	for i, s := range steps {
		tr.Sync(t0.Add(time.Duration(i)*time.Second), []Observed[int]{obs(1, "A", s.state)})
		want = append(want, s.want)
		assertKinds(t, l.Entries(), want...)
	}
}

func TestSyncIsIdempotentWhenNothingChanged(t *testing.T) {
	tr, l := newTracker()
	tr.Seed([]Observed[int]{obs(1, "A", StateRunning)})
	tr.Sync(t0, []Observed[int]{obs(1, "A", StateRunning)})
	tr.Sync(t0, []Observed[int]{obs(1, "A", StateRunning)})
	if n := len(l.Entries()); n != 0 {
		t.Fatalf("unchanged sync logged %d entries", n)
	}
}

func TestSyncReportsWhetherAnythingWasLogged(t *testing.T) {
	tr, _ := newTracker()
	tr.Seed([]Observed[int]{obs(1, "A", StateIdle)})
	if tr.Sync(t0, []Observed[int]{obs(1, "A", StateIdle)}) {
		t.Fatal("no-op sync reported a change")
	}
	if !tr.Sync(t0, []Observed[int]{obs(1, "A", StateRunning)}) {
		t.Fatal("transition sync reported no change")
	}
}

func TestRenameUpdatesNameWithoutLogging(t *testing.T) {
	tr, l := newTracker()
	tr.Seed([]Observed[int]{obs(1, "A", StateIdle)})
	tr.Sync(t0, []Observed[int]{obs(1, "Deep Work", StateIdle)})
	if n := len(l.Entries()); n != 0 {
		t.Fatalf("rename logged %d entries", n)
	}
	tr.Sync(t0, []Observed[int]{})
	if got := l.Entries()[0].TimerName; got != "Deep Work" {
		t.Fatalf("removed entry used stale name %q", got)
	}
}

func TestRekeyTreatsReplacementAsSameTimer(t *testing.T) {
	tr, l := newTracker()
	tr.Seed([]Observed[int]{obs(1, "A", StateIdle)})
	tr.Rekey(1, 9)
	tr.Sync(t0, []Observed[int]{obs(9, "A", StateIdle)})
	if n := len(l.Entries()); n != 0 {
		t.Fatalf("rekeyed replacement logged %d entries", n)
	}
}
