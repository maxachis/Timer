package main

import (
	"path/filepath"
	"testing"
	"time"

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
	s := settings.Default()
	return &App{
		timers:   timer.NewCollection(s.DefaultDurationSecs),
		settings: s,
		store:    store,
	}, path
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
