package timer

import (
	"testing"
	"time"
)

func TestCreationAndInitialState(t *testing.T) {
	tm := NewCountdownTimer(300)
	if tm.Remaining() != 300*time.Second {
		t.Fatalf("expected 300s remaining, got %v", tm.Remaining())
	}
	if tm.IsFinished() {
		t.Fatal("should not be finished")
	}
}

func TestStartPauseResumeLifecycle(t *testing.T) {
	tm := NewCountdownTimer(300)
	tm.Start()
	time.Sleep(50 * time.Millisecond)
	r1 := tm.Remaining()
	if r1 >= 300*time.Second {
		t.Fatal("remaining should have decreased")
	}
	tm.Pause()
	paused := tm.Remaining()
	time.Sleep(50 * time.Millisecond)
	if tm.Remaining() != paused {
		t.Fatal("paused remaining should not change")
	}
	tm.Resume()
	time.Sleep(50 * time.Millisecond)
	if tm.Remaining() >= paused {
		t.Fatal("remaining should decrease after resume")
	}
}

func TestAddTimeDefault(t *testing.T) {
	tm := NewCountdownTimer(120)
	tm.Start()
	tm.AddTime(0)
	if tm.Remaining() <= 400*time.Second {
		t.Fatalf("expected >400s, got %v", tm.Remaining())
	}
}

func TestAddTimeCustomWhilePaused(t *testing.T) {
	tm := NewCountdownTimer(120)
	tm.Start()
	tm.Pause()
	before := tm.Remaining()
	tm.AddTime(60 * time.Second)
	after := tm.Remaining()
	if after-before != 60*time.Second {
		t.Fatalf("expected +60s, got %v", after-before)
	}
}

func TestRemoveTimeDefault(t *testing.T) {
	tm := NewCountdownTimer(600)
	tm.Start()
	tm.RemoveTime(0)
	if tm.Remaining() >= 305*time.Second {
		t.Fatalf("expected <305s, got %v", tm.Remaining())
	}
}

func TestRemoveTimeClampsToZero(t *testing.T) {
	tm := NewCountdownTimer(60)
	tm.Start()
	tm.RemoveTime(120 * time.Second)
	if !tm.IsFinished() {
		t.Fatal("should be finished")
	}
	if tm.Remaining() != 0 {
		t.Fatalf("expected 0, got %v", tm.Remaining())
	}
}

func TestResetFromRunning(t *testing.T) {
	tm := NewCountdownTimer(300)
	tm.Start()
	time.Sleep(50 * time.Millisecond)
	tm.Reset()
	if tm.Remaining() != 300*time.Second {
		t.Fatal("reset should restore duration")
	}
	if tm.IsFinished() {
		t.Fatal("should not be finished after reset")
	}
	tm.Start()
	time.Sleep(10 * time.Millisecond)
	if tm.Remaining() >= 300*time.Second {
		t.Fatal("should tick after restart")
	}
}

func TestResetFromPaused(t *testing.T) {
	tm := NewCountdownTimer(300)
	tm.Start()
	time.Sleep(50 * time.Millisecond)
	tm.Pause()
	tm.Reset()
	if tm.Remaining() != 300*time.Second {
		t.Fatal("reset should restore duration")
	}
	if tm.IsFinished() {
		t.Fatal("should not be finished")
	}
}

func TestIsFinishedDetection(t *testing.T) {
	tm := NewCountdownTimer(0)
	tm.Start()
	if !tm.IsFinished() {
		t.Fatal("zero-duration timer should be finished")
	}
	tm2 := NewCountdownTimer(300)
	if tm2.IsFinished() {
		t.Fatal("fresh timer should not be finished")
	}
}

func TestCollectionNewHasOneTimer(t *testing.T) {
	c := NewCollection(300)
	if c.Count() != 1 || c.ActiveIndex() != 0 || c.ActiveName() != "" {
		t.Fatal("collection init bad")
	}
	if c.Active().Remaining() != 300*time.Second {
		t.Fatal("bad initial remaining")
	}
}

func TestCollectionFromNames(t *testing.T) {
	c := FromNames([]string{"Work", "Break"}, 300)
	if c.Count() != 2 || c.ActiveName() != "Work" {
		t.Fatal("from_names failed")
	}
}

func TestCollectionFromEmptyNames(t *testing.T) {
	c := FromNames(nil, 300)
	if c.Count() != 1 {
		t.Fatal("empty names should yield 1 timer")
	}
}

func TestAddTimer(t *testing.T) {
	c := NewCollection(300)
	idx, err := c.AddTimer("Second", 600)
	if err != nil || idx != 1 || c.Count() != 2 {
		t.Fatal("add failed")
	}
}

func TestAddTimerMaxLimit(t *testing.T) {
	c := NewCollection(300)
	c.AddTimer("Two", 300)
	c.AddTimer("Three", 300)
	if _, err := c.AddTimer("Four", 300); err == nil {
		t.Fatal("expected max limit error")
	}
	if c.Count() != 3 {
		t.Fatal("count should be 3")
	}
}

func TestRemoveTimer(t *testing.T) {
	c := NewCollection(300)
	c.AddTimer("Second", 600)
	if err := c.RemoveTimer(1); err != nil {
		t.Fatal(err)
	}
	if c.Count() != 1 {
		t.Fatal("should have 1 timer")
	}
}

func TestRemoveLastTimerFails(t *testing.T) {
	c := NewCollection(300)
	if err := c.RemoveTimer(0); err == nil {
		t.Fatal("expected error")
	}
}

func TestRemoveActiveAdjustsIndex(t *testing.T) {
	c := NewCollection(300)
	c.AddTimer("Second", 300)
	c.AddTimer("Third", 300)
	c.SwitchTo(2)
	if c.ActiveIndex() != 2 {
		t.Fatal("switch failed")
	}
	c.RemoveTimer(2)
	if c.ActiveIndex() != 1 {
		t.Fatalf("expected 1, got %d", c.ActiveIndex())
	}
}

func TestRemoveBeforeActiveAdjustsIndex(t *testing.T) {
	c := NewCollection(300)
	c.AddTimer("Second", 300)
	c.AddTimer("Third", 300)
	c.SwitchTo(2)
	c.RemoveTimer(0)
	if c.ActiveIndex() != 1 {
		t.Fatalf("expected 1, got %d", c.ActiveIndex())
	}
}

func TestSwitchPausesCurrent(t *testing.T) {
	c := NewCollection(300)
	c.AddTimer("Second", 300)
	c.Active().Start()
	if c.Active().StateName() != "running" {
		t.Fatal("should be running")
	}
	c.SwitchTo(1)
	if c.TimerAt(0).StateName() != "paused" {
		t.Fatal("previous should be paused")
	}
	if c.ActiveIndex() != 1 {
		t.Fatal("active index wrong")
	}
}

func TestSwitchToSameIsNoop(t *testing.T) {
	c := NewCollection(300)
	c.Active().Start()
	c.SwitchTo(0)
	if c.Active().StateName() != "running" {
		t.Fatal("should still be running")
	}
}

func TestSwitchNextWraps(t *testing.T) {
	c := NewCollection(300)
	c.AddTimer("Two", 300)
	c.AddTimer("Three", 300)
	c.SwitchNext()
	if c.ActiveIndex() != 1 {
		t.Fatal()
	}
	c.SwitchNext()
	if c.ActiveIndex() != 2 {
		t.Fatal()
	}
	c.SwitchNext()
	if c.ActiveIndex() != 0 {
		t.Fatal("should wrap")
	}
}

func TestSwitchPrevWraps(t *testing.T) {
	c := NewCollection(300)
	c.AddTimer("Two", 300)
	c.SwitchPrev()
	if c.ActiveIndex() != 1 {
		t.Fatal("should wrap to end")
	}
	c.SwitchPrev()
	if c.ActiveIndex() != 0 {
		t.Fatal()
	}
}

func TestSwitchSingleTimerNoop(t *testing.T) {
	c := NewCollection(300)
	c.SwitchNext()
	if c.ActiveIndex() != 0 {
		t.Fatal()
	}
	c.SwitchPrev()
	if c.ActiveIndex() != 0 {
		t.Fatal()
	}
}

func TestRenameTimer(t *testing.T) {
	c := NewCollection(300)
	c.RenameTimer(0, "Focus")
	if c.ActiveName() != "Focus" {
		t.Fatal()
	}
}

func TestTimerList(t *testing.T) {
	c := NewCollection(300)
	c.AddTimer("Break", 120)
	list := c.TimerList()
	if len(list) != 2 || !list[0].IsActive || list[1].IsActive || list[1].Name != "Break" {
		t.Fatalf("bad list: %+v", list)
	}
}

func TestTimerNames(t *testing.T) {
	c := NewCollection(300)
	c.RenameTimer(0, "Work")
	c.AddTimer("Break", 120)
	names := c.TimerNames()
	if len(names) != 2 || names[0] != "Work" || names[1] != "Break" {
		t.Fatalf("bad names: %v", names)
	}
}

func TestSnapshotIdleTimer(t *testing.T) {
	c := NewCollection(300)
	c.RenameTimer(0, "Work")
	snap := c.Snapshot()
	if len(snap.Timers) != 1 || snap.ActiveIndex != 0 {
		t.Fatalf("bad snapshot: %+v", snap)
	}
	ts := snap.Timers[0]
	if ts.Name != "Work" || ts.State != "idle" || ts.ElapsedMs != 0 {
		t.Fatalf("bad timer snapshot: %+v", ts)
	}
	if ts.OriginalDurationMs != 300_000 || ts.DurationMs != 300_000 {
		t.Fatalf("bad durations: %+v", ts)
	}
}

func TestSnapshotRunningCollapsesToPaused(t *testing.T) {
	c := NewCollection(300)
	c.Active().Start()
	time.Sleep(50 * time.Millisecond)
	snap := c.Snapshot()
	ts := snap.Timers[0]
	if ts.State != "paused" {
		t.Fatalf("running should snapshot as paused, got %q", ts.State)
	}
	if ts.ElapsedMs == 0 {
		t.Fatal("elapsed should be > 0 for a timer that has run")
	}
}

func TestSnapshotPreservesAddTimeAdjustments(t *testing.T) {
	c := NewCollection(120)
	c.Active().Start()
	c.Active().AddTime(60 * time.Second)
	c.Active().Pause()
	snap := c.Snapshot()
	ts := snap.Timers[0]
	if ts.OriginalDurationMs != 120_000 {
		t.Fatalf("original drifted: %d", ts.OriginalDurationMs)
	}
	if ts.DurationMs != 180_000 {
		t.Fatalf("duration after +60s should be 180000ms, got %d", ts.DurationMs)
	}
}

func TestSnapshotFinishedTimer(t *testing.T) {
	c := NewCollection(60)
	c.Active().Start()
	c.Active().RemoveTime(120 * time.Second)
	snap := c.Snapshot()
	if snap.Timers[0].State != "finished" {
		t.Fatalf("expected finished, got %q", snap.Timers[0].State)
	}
}

func TestRestoreRoundtripIdle(t *testing.T) {
	c := NewCollection(300)
	c.RenameTimer(0, "Work")
	c.AddTimer("Break", 600)
	snap := c.Snapshot()

	r := RestoreCollection(snap, 300)
	if r.Count() != 2 || r.ActiveIndex() != 0 {
		t.Fatalf("bad restored shape")
	}
	if r.TimerAt(0).Remaining() != 300*time.Second || r.TimerAt(1).Remaining() != 600*time.Second {
		t.Fatal("bad restored remaining")
	}
	if r.TimerAt(0).StateName() != "idle" || r.TimerAt(1).StateName() != "idle" {
		t.Fatal("idle timers should restore as idle")
	}
}

func TestRestoreFreezesRunningAsPaused(t *testing.T) {
	c := NewCollection(300)
	c.Active().Start()
	time.Sleep(50 * time.Millisecond)
	expected := c.Active().Remaining()

	snap := c.Snapshot()
	r := RestoreCollection(snap, 300)

	if r.Active().StateName() != "paused" {
		t.Fatalf("expected paused, got %q", r.Active().StateName())
	}
	got := r.Active().Remaining()
	// Tolerate small drift from snapshot timing.
	if diff := expected - got; diff > 5*time.Millisecond || diff < -5*time.Millisecond {
		t.Fatalf("remaining drifted: expected %v, got %v", expected, got)
	}
	// Confirm resume picks up from frozen elapsed.
	r.Active().Resume()
	time.Sleep(20 * time.Millisecond)
	if r.Active().Remaining() >= got {
		t.Fatal("resume should tick down from frozen elapsed")
	}
}

func TestRestoreActiveIndexPreserved(t *testing.T) {
	c := NewCollection(300)
	c.AddTimer("B", 300)
	c.AddTimer("C", 300)
	c.SwitchTo(2)
	snap := c.Snapshot()
	r := RestoreCollection(snap, 300)
	if r.ActiveIndex() != 2 {
		t.Fatalf("active index lost: %d", r.ActiveIndex())
	}
}

func TestRestoreEmptySnapshotFallsBackToDefault(t *testing.T) {
	r := RestoreCollection(CollectionSnapshot{}, 300)
	if r.Count() != 1 || r.Active().Remaining() != 300*time.Second {
		t.Fatal("empty snapshot should fall back to default single timer")
	}
}

func TestRestoreClampsActiveIndex(t *testing.T) {
	snap := CollectionSnapshot{
		Timers: []TimerSnapshot{
			{Name: "A", OriginalDurationMs: 300_000, DurationMs: 300_000, State: "idle"},
		},
		ActiveIndex: 5,
	}
	r := RestoreCollection(snap, 300)
	if r.ActiveIndex() != 0 {
		t.Fatalf("out-of-range active index should clamp to 0, got %d", r.ActiveIndex())
	}
}

func TestRestoreCapsToMaxTimers(t *testing.T) {
	timers := make([]TimerSnapshot, MaxTimers+2)
	for i := range timers {
		timers[i] = TimerSnapshot{OriginalDurationMs: 300_000, DurationMs: 300_000, State: "idle"}
	}
	r := RestoreCollection(CollectionSnapshot{Timers: timers}, 300)
	if r.Count() != MaxTimers {
		t.Fatalf("expected count capped at %d, got %d", MaxTimers, r.Count())
	}
}
