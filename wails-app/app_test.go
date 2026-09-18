package main

import (
	"path/filepath"
	"testing"
	"time"

	"wails-app/internal/history"
	"wails-app/internal/settings"
	"wails-app/internal/timer"
)

func newTestApp(t *testing.T) (*App, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	store, err := settings.OpenStoreAt(path)
	if err != nil {
		t.Fatal(err)
	}
	return newAppFromStore(store), path
}

func loadSnap(t *testing.T, path string) (timer.CollectionSnapshot, bool) {
	t.Helper()
	store, err := settings.OpenStoreAt(path)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	var snap timer.CollectionSnapshot
	ok, err := store.LoadTimerStates(&snap)
	if err != nil {
		t.Fatalf("LoadTimerStates: %v", err)
	}
	return snap, ok
}

// TestAddNewTimerPersistsFullState verifies that adding a timer writes a full
// timer_states snapshot so the new timer survives a crash (no graceful shutdown).
func TestAddNewTimerPersistsFullState(t *testing.T) {
	app, path := newTestApp(t)
	if _, err := app.AddNewTimer("Focus"); err != nil {
		t.Fatal(err)
	}

	snap, ok := loadSnap(t, path)
	if !ok {
		t.Fatal("timer_states not written after AddNewTimer")
	}
	if len(snap.Timers) != 2 {
		t.Fatalf("expected 2 timers, got %d", len(snap.Timers))
	}
	if snap.Timers[1].Name != "Focus" {
		t.Fatalf("expected name %q, got %q", "Focus", snap.Timers[1].Name)
	}
}

// TestRemoveExistingTimerPersistsFullState verifies that removing a timer
// writes a full snapshot, preventing the removed timer from reappearing after a crash.
func TestRemoveExistingTimerPersistsFullState(t *testing.T) {
	app, path := newTestApp(t)
	if _, err := app.AddNewTimer("Extra"); err != nil {
		t.Fatal(err)
	}
	if err := app.RemoveExistingTimer(1); err != nil {
		t.Fatal(err)
	}

	snap, ok := loadSnap(t, path)
	if !ok {
		t.Fatal("timer_states not written after RemoveExistingTimer")
	}
	if len(snap.Timers) != 1 {
		t.Fatalf("expected 1 timer after remove, got %d", len(snap.Timers))
	}
}

// TestRenameTimerPersistsFullState verifies that renaming a timer writes a full
// snapshot so the rename survives a crash.
func TestRenameTimerPersistsFullState(t *testing.T) {
	app, path := newTestApp(t)
	if err := app.RenameTimer(0, "Deep Work"); err != nil {
		t.Fatal(err)
	}

	snap, ok := loadSnap(t, path)
	if !ok {
		t.Fatal("timer_states not written after RenameTimer")
	}
	if snap.Timers[0].Name != "Deep Work" {
		t.Fatalf("expected name %q, got %q", "Deep Work", snap.Timers[0].Name)
	}
}

// TestAutosaveWritesElapsedProgress verifies that the periodic autosave captures
// elapsed time from a running timer, so progress survives a crash between saves.
func TestAutosaveWritesElapsedProgress(t *testing.T) {
	app, path := newTestApp(t)
	app.StartTimer()

	cancel := app.startAutosave(10 * time.Millisecond)
	defer cancel()

	// Poll rather than sleep a fixed amount: under load the ticker goroutine
	// may lag arbitrarily. The store writes non-atomically, so a read racing
	// the autosave can see a partial file — treat that as "not written yet".
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		store, err := settings.OpenStoreAt(path)
		if err == nil {
			var snap timer.CollectionSnapshot
			if ok, err := store.LoadTimerStates(&snap); err == nil && ok &&
				len(snap.Timers) > 0 && snap.Timers[0].ElapsedMs > 0 {
				return
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("autosave never wrote a snapshot with elapsed progress")
}

func historyKinds(app *App) []string {
	entries := app.GetHistory()
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Kind
	}
	return out
}

func assertHistoryKinds(t *testing.T, app *App, want ...string) {
	t.Helper()
	got := historyKinds(app)
	if len(got) != len(want) {
		t.Fatalf("history kinds = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("history kinds = %v, want %v", got, want)
		}
	}
}

func TestHistoryLogsStartPauseResumeReset(t *testing.T) {
	app, _ := newTestApp(t)
	app.StartTimer()
	app.PauseTimer()
	app.ResumeTimer()
	app.ResetTimer()
	assertHistoryKinds(t, app,
		history.KindStarted, history.KindPaused, history.KindResumed, history.KindReset)
}

func TestHistoryLogsCreatedAndRemovedWithNames(t *testing.T) {
	app, _ := newTestApp(t)
	if _, err := app.AddNewTimer("Focus"); err != nil {
		t.Fatal(err)
	}
	if err := app.RemoveExistingTimer(1); err != nil {
		t.Fatal(err)
	}
	entries := app.GetHistory()
	assertHistoryKinds(t, app, history.KindCreated, history.KindRemoved)
	for _, e := range entries {
		if e.TimerName != "Focus" {
			t.Fatalf("entry %+v should name the Focus timer", e)
		}
	}
}

func TestHistoryNamesUnnamedTimersByPosition(t *testing.T) {
	app, _ := newTestApp(t)
	app.StartTimer()
	if got := app.GetHistory()[0].TimerName; got != "Timer 1" {
		t.Fatalf("unnamed timer logged as %q, want %q", got, "Timer 1")
	}
}

func TestHistoryLogsFinishFromRemovingTime(t *testing.T) {
	app, _ := newTestApp(t)
	app.StartTimer()
	app.RemoveTimeCustom(settings.Default().DefaultDurationSecs + 1)
	assertHistoryKinds(t, app, history.KindStarted, history.KindFinished)
}

// A timer that runs out on its own has no App call to hook; the status poll
// is where the finish is noticed, so polling must sync the log.
func TestHistoryStatusPollNoticesTransitions(t *testing.T) {
	app, _ := newTestApp(t)
	app.timers.Active().Start()
	assertHistoryKinds(t, app)
	app.GetTimerStatus()
	assertHistoryKinds(t, app, history.KindStarted)
}

func TestHistoryLogsPauseWhenSwitchingAwayFromRunningTimer(t *testing.T) {
	app, _ := newTestApp(t)
	if _, err := app.AddNewTimer("Other"); err != nil {
		t.Fatal(err)
	}
	app.StartTimer()
	if err := app.SwitchTimer(1); err != nil {
		t.Fatal(err)
	}
	entries := app.GetHistory()
	assertHistoryKinds(t, app, history.KindCreated, history.KindStarted, history.KindPaused)
	if entries[2].TimerName != "Timer 1" {
		t.Fatalf("pause logged against %q, want the timer that was running", entries[2].TimerName)
	}
}

// Saving settings while idle swaps the active timer for one with the new
// default duration; that's not the user creating or removing anything.
func TestHistoryIgnoresIdleTimerReplacement(t *testing.T) {
	app, _ := newTestApp(t)
	app.CreateTimer(600)
	assertHistoryKinds(t, app)
}

func TestHistorySurvivesRestartWithoutReplayingRestoredTimers(t *testing.T) {
	app, path := newTestApp(t)
	if _, err := app.AddNewTimer("Focus"); err != nil {
		t.Fatal(err)
	}
	app.StartTimer()

	store, err := settings.OpenStoreAt(path)
	if err != nil {
		t.Fatal(err)
	}
	reopened := newAppFromStore(store)
	assertHistoryKinds(t, reopened, history.KindCreated, history.KindStarted)
	// The restored timers are the baseline, not fresh creations.
	reopened.GetTimerStatus()
	assertHistoryKinds(t, reopened, history.KindCreated, history.KindStarted)
}
