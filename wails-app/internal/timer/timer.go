package timer

import (
	"errors"
	"fmt"
	"time"
)

const DefaultIncrement = 5 * time.Minute
const MaxTimers = 7

type stateKind int

const (
	stateIdle stateKind = iota
	stateRunning
	statePaused
	stateFinished
)

type CountdownTimer struct {
	originalDuration  time.Duration
	duration          time.Duration
	kind              stateKind
	startedAt         time.Time
	previouslyElapsed time.Duration
	pausedElapsed     time.Duration
}

func NewCountdownTimer(durationSecs uint64) *CountdownTimer {
	d := time.Duration(durationSecs) * time.Second
	return &CountdownTimer{originalDuration: d, duration: d, kind: stateIdle}
}

func (t *CountdownTimer) Start() {
	if t.kind == stateIdle {
		t.kind = stateRunning
		t.startedAt = time.Now()
		t.previouslyElapsed = 0
	}
}

func (t *CountdownTimer) Pause() {
	if t.kind == stateRunning {
		elapsed := t.previouslyElapsed + time.Since(t.startedAt)
		t.kind = statePaused
		t.pausedElapsed = elapsed
	}
}

func (t *CountdownTimer) Resume() {
	if t.kind == statePaused {
		t.kind = stateRunning
		t.startedAt = time.Now()
		t.previouslyElapsed = t.pausedElapsed
	}
}

func (t *CountdownTimer) Remaining() time.Duration {
	var elapsed time.Duration
	switch t.kind {
	case stateIdle:
		return t.duration
	case stateRunning:
		elapsed = t.previouslyElapsed + time.Since(t.startedAt)
	case statePaused:
		elapsed = t.pausedElapsed
	case stateFinished:
		return 0
	}
	if elapsed >= t.duration {
		return 0
	}
	return t.duration - elapsed
}

func (t *CountdownTimer) IsFinished() bool {
	return t.kind == stateFinished || t.Remaining() == 0
}

func (t *CountdownTimer) AddTime(amount time.Duration) {
	if amount == 0 {
		amount = DefaultIncrement
	}
	if t.kind == stateRunning || t.kind == statePaused {
		t.duration += amount
	}
}

func (t *CountdownTimer) RemoveTime(amount time.Duration) {
	if amount == 0 {
		amount = DefaultIncrement
	}
	if t.kind == stateRunning || t.kind == statePaused {
		remaining := t.Remaining()
		if remaining <= amount {
			t.kind = stateFinished
		} else {
			t.duration -= amount
		}
	}
}

func (t *CountdownTimer) StateName() string {
	if t.IsFinished() {
		return "finished"
	}
	switch t.kind {
	case stateIdle:
		return "idle"
	case stateRunning:
		return "running"
	case statePaused:
		return "paused"
	default:
		return "finished"
	}
}

func (t *CountdownTimer) Reset() {
	t.duration = t.originalDuration
	t.kind = stateIdle
	t.previouslyElapsed = 0
	t.pausedElapsed = 0
}

// SetMaxToRemaining commits the timer's current remaining time as its
// original duration: both the value Reset() restores and the maximum its
// stacked-bar lane is scaled against. The timer keeps its run state and its
// displayed time; only the baseline moves. A no-op on a finished/zero timer,
// since a zero maximum would collapse the lane.
func (t *CountdownTimer) SetMaxToRemaining() {
	remaining := t.Remaining()
	if remaining <= 0 {
		return
	}
	t.originalDuration = remaining
	t.duration = remaining
	t.previouslyElapsed = 0
	t.pausedElapsed = 0
	if t.kind == stateRunning {
		t.startedAt = time.Now()
	} else if t.kind == stateFinished {
		t.kind = stateIdle
	}
}

type NamedTimer struct {
	Name  string
	Timer *CountdownTimer
	Tint  string
}

type TimerInfo struct {
	Index                int     `json:"index"`
	Name                 string  `json:"name"`
	State                string  `json:"state"`
	RemainingSecs        float64 `json:"remaining_secs"`
	IsActive             bool    `json:"is_active"`
	Tint                 string  `json:"tint"`
	OriginalDurationSecs float64 `json:"original_duration_secs"`
	IsBreakTimer         bool    `json:"is_break_timer"`
	BreakAccruedSecs     float64 `json:"break_accrued_secs"`
}

var validTints = map[string]bool{
	"": true, "none": true, "rose": true, "sage": true,
	"sky": true, "lavender": true, "amber": true, "teal": true, "clay": true,
}

// SetTint sets the background tint id for the timer at index. "none" is
// normalised to an empty string. Unknown ids or out-of-range indexes error.
func (c *TimerCollection) SetTint(index int, tint string) error {
	if index < 0 || index >= len(c.timers) {
		return errors.New("Invalid timer index")
	}
	if !validTints[tint] {
		return errors.New("Invalid tint")
	}
	if tint == "none" {
		tint = ""
	}
	c.timers[index].Tint = tint
	return nil
}

type TimerCollection struct {
	timers          []*NamedTimer
	activeIndex     int
	breakIndex      int
	onBreak         bool
	preBreakActive  int
	preBreakRunning bool
	breakAccrued    time.Duration
	// breakStartRemaining is the break timer's remaining time when the current
	// break began; the difference against its live remaining is how much of the
	// accrued allocation this break has consumed so far.
	breakStartRemaining time.Duration
}

func NewCollection(defaultDurationSecs uint64) *TimerCollection {
	return &TimerCollection{
		timers:      []*NamedTimer{{Name: "", Timer: NewCountdownTimer(defaultDurationSecs)}},
		activeIndex: 0,
		breakIndex:  -1,
	}
}

func FromNames(names []string, defaultDurationSecs uint64) *TimerCollection {
	if len(names) == 0 {
		return NewCollection(defaultDurationSecs)
	}
	n := len(names)
	if n > MaxTimers {
		n = MaxTimers
	}
	ts := make([]*NamedTimer, 0, n)
	for i := 0; i < n; i++ {
		ts = append(ts, &NamedTimer{Name: names[i], Timer: NewCountdownTimer(defaultDurationSecs)})
	}
	return &TimerCollection{timers: ts, activeIndex: 0, breakIndex: -1}
}

func (c *TimerCollection) Active() *CountdownTimer { return c.timers[c.activeIndex].Timer }
func (c *TimerCollection) ActiveName() string      { return c.timers[c.activeIndex].Name }
func (c *TimerCollection) ActiveIndex() int        { return c.activeIndex }
func (c *TimerCollection) Count() int              { return len(c.timers) }

func (c *TimerCollection) AddTimer(name string, durationSecs uint64) (int, error) {
	if len(c.timers) >= MaxTimers {
		return 0, fmt.Errorf("Maximum of %d timers reached. Remove a timer before adding another.", MaxTimers)
	}
	c.timers = append(c.timers, &NamedTimer{Name: name, Timer: NewCountdownTimer(durationSecs)})
	return len(c.timers) - 1, nil
}

func (c *TimerCollection) RemoveTimer(index int) error {
	if len(c.timers) <= 1 {
		return errors.New("Cannot remove the last timer")
	}
	if index < 0 || index >= len(c.timers) {
		return errors.New("Invalid timer index")
	}
	if index == c.breakIndex {
		if c.onBreak {
			c.EndBreak()
		}
		c.breakIndex = -1
		c.breakAccrued = 0
	} else if index < c.breakIndex {
		c.breakIndex--
	}
	c.timers = append(c.timers[:index], c.timers[index+1:]...)
	if c.activeIndex >= len(c.timers) {
		c.activeIndex = len(c.timers) - 1
	} else if index < c.activeIndex {
		c.activeIndex--
	}
	return nil
}

// BreakIndex returns the index of the designated break timer, or -1 if none.
func (c *TimerCollection) BreakIndex() int { return c.breakIndex }

// IsOnBreak reports whether the collection is currently in break mode.
func (c *TimerCollection) IsOnBreak() bool { return c.onBreak }

// SetMaxToCurrent commits the active timer's current time as its maximum.
func (c *TimerCollection) SetMaxToCurrent() {
	c.timers[c.activeIndex].Timer.SetMaxToRemaining()
}

// ToggleBreakDesignation designates the currently-active timer as the break
// timer, or clears the designation if it is already the break timer.
func (c *TimerCollection) ToggleBreakDesignation() {
	c.breakAccrued = 0
	if c.breakIndex == c.activeIndex {
		c.breakIndex = -1
		if c.onBreak {
			c.EndBreak()
		}
	} else {
		c.breakIndex = c.activeIndex
	}
}

// StartBreak pauses the active timer, remembers it, and switches to and
// starts/resumes the designated break timer. The break draws down the accrued
// break allocation, so it errors if no break timer is designated or the
// allocation is empty. A no-op if already on break.
func (c *TimerCollection) StartBreak() error {
	if c.breakIndex < 0 {
		return errors.New("No break timer designated")
	}
	if c.onBreak {
		return nil
	}
	if c.BreakAccrued() <= 0 {
		return errors.New("No break time accrued")
	}
	c.preBreakActive = c.activeIndex
	c.preBreakRunning = c.timers[c.activeIndex].Timer.kind == stateRunning
	c.timers[c.activeIndex].Timer.Pause()
	c.onBreak = true
	c.activeIndex = c.breakIndex
	bt := c.timers[c.breakIndex].Timer
	switch bt.kind {
	case stateIdle:
		bt.Start()
	case statePaused:
		bt.Resume()
	}
	c.breakStartRemaining = bt.Remaining()
	return nil
}

// breakConsumed reports how much break time the in-progress break has used:
// the break timer runs while on break, so its drop in remaining time since the
// break began is exactly the allocation consumed. Zero when not on break.
func (c *TimerCollection) breakConsumed() time.Duration {
	if !c.onBreak {
		return 0
	}
	consumed := c.breakStartRemaining - c.timers[c.breakIndex].Timer.Remaining()
	if consumed < 0 {
		return 0
	}
	return consumed
}

// EndBreak commits the break time consumed, pauses the break timer, and
// restores the previously-active timer, resuming it if it had been running.
// A no-op if not on break.
func (c *TimerCollection) EndBreak() {
	if !c.onBreak {
		return
	}
	c.breakAccrued -= c.breakConsumed()
	if c.breakAccrued < 0 {
		c.breakAccrued = 0
	}
	c.timers[c.breakIndex].Timer.Pause()
	c.onBreak = false
	if c.preBreakActive >= 0 && c.preBreakActive < len(c.timers) {
		c.activeIndex = c.preBreakActive
		if c.preBreakRunning {
			c.timers[c.activeIndex].Timer.Resume()
		}
	}
}

// ToggleBreak starts break mode if not on break, or ends it if already on break.
func (c *TimerCollection) ToggleBreak() error {
	if c.onBreak {
		c.EndBreak()
		return nil
	}
	return c.StartBreak()
}

// AddBreakTime grows the accrued break allocation for the designated break
// timer, clamped to at most the break timer's remaining time (parity).
func (c *TimerCollection) AddBreakTime(amount time.Duration) {
	if c.breakIndex < 0 {
		return
	}
	c.breakAccrued += amount
	if c.breakAccrued < 0 {
		c.breakAccrued = 0
	}
	// The un-consumed allocation can never exceed the break block's remaining
	// time; mid-break, the stored value still carries the consumed portion.
	maxAllocation := c.timers[c.breakIndex].Timer.Remaining() + c.breakConsumed()
	if c.breakAccrued > maxAllocation {
		c.breakAccrued = maxAllocation
	}
}

// RemoveBreakTime shrinks the accrued break allocation, floored at zero.
func (c *TimerCollection) RemoveBreakTime(amount time.Duration) {
	if c.breakIndex < 0 {
		return
	}
	c.breakAccrued -= amount
	if c.breakAccrued < 0 {
		c.breakAccrued = 0
	}
}

// ResetBreakTime clears the accrued break allocation back to zero.
func (c *TimerCollection) ResetBreakTime() {
	if c.breakIndex < 0 {
		return
	}
	c.breakAccrued = 0
}

// BreakAccrued returns the break allocation still available on the designated
// break timer: the stored allocation minus whatever an in-progress break has
// consumed, clamped to the break block's remaining time. Zero if no break
// timer is designated.
func (c *TimerCollection) BreakAccrued() time.Duration {
	if c.breakIndex < 0 {
		return 0
	}
	acc := c.breakAccrued - c.breakConsumed()
	if rem := c.timers[c.breakIndex].Timer.Remaining(); acc > rem {
		acc = rem
	}
	if acc < 0 {
		return 0
	}
	return acc
}

func (c *TimerCollection) SwitchTo(index int) error {
	if index < 0 || index >= len(c.timers) {
		return errors.New("Invalid timer index")
	}
	if index == c.activeIndex {
		return nil
	}
	c.timers[c.activeIndex].Timer.Pause()
	c.activeIndex = index
	c.timers[c.activeIndex].Timer.Pause()
	return nil
}

func (c *TimerCollection) SwitchNext() {
	if len(c.timers) <= 1 {
		return
	}
	_ = c.SwitchTo((c.activeIndex + 1) % len(c.timers))
}

func (c *TimerCollection) SwitchPrev() {
	if len(c.timers) <= 1 {
		return
	}
	prev := c.activeIndex - 1
	if prev < 0 {
		prev = len(c.timers) - 1
	}
	_ = c.SwitchTo(prev)
}

func (c *TimerCollection) RenameTimer(index int, name string) error {
	if index < 0 || index >= len(c.timers) {
		return errors.New("Invalid timer index")
	}
	c.timers[index].Name = name
	return nil
}

func (c *TimerCollection) TimerList() []TimerInfo {
	out := make([]TimerInfo, len(c.timers))
	for i, t := range c.timers {
		breakAccruedSecs := 0.0
		if i == c.breakIndex {
			breakAccruedSecs = c.BreakAccrued().Seconds()
		}
		out[i] = TimerInfo{
			Index:                i,
			Name:                 t.Name,
			State:                t.Timer.StateName(),
			RemainingSecs:        t.Timer.Remaining().Seconds(),
			IsActive:             i == c.activeIndex,
			Tint:                 t.Tint,
			OriginalDurationSecs: t.Timer.originalDuration.Seconds(),
			IsBreakTimer:         i == c.breakIndex,
			BreakAccruedSecs:     breakAccruedSecs,
		}
	}
	return out
}

func (c *TimerCollection) TimerNames() []string {
	out := make([]string, len(c.timers))
	for i, t := range c.timers {
		out[i] = t.Name
	}
	return out
}

// TimerAt exposes a timer at index (used by tests).
func (c *TimerCollection) TimerAt(i int) *CountdownTimer { return c.timers[i].Timer }

// ReplaceActive replaces the active timer with a new timer of the given duration.
func (c *TimerCollection) ReplaceActive(durationSecs uint64) {
	c.timers[c.activeIndex].Timer = NewCountdownTimer(durationSecs)
}

// TimerSnapshot is a serialisable, point-in-time view of a single timer.
// A running timer is captured as paused so reopening the app never lets a
// timer expire while it was closed. Durations are stored in milliseconds
// to preserve sub-second precision across save/restore.
type TimerSnapshot struct {
	Name               string `json:"name"`
	OriginalDurationMs uint64 `json:"original_duration_ms"`
	DurationMs         uint64 `json:"duration_ms"`
	ElapsedMs          uint64 `json:"elapsed_ms"`
	State              string `json:"state"` // "idle" | "paused" | "finished"
	Tint               string `json:"tint"`
}

// CollectionSnapshot is the serialisable view of a TimerCollection.
type CollectionSnapshot struct {
	Timers      []TimerSnapshot `json:"timers"`
	ActiveIndex int             `json:"active_index"`
}

func msOf(d time.Duration) uint64 { return uint64(d / time.Millisecond) }

func (t *CountdownTimer) snapshot() (durationMs, elapsedMs uint64, state string) {
	dur := msOf(t.duration)
	switch {
	case t.IsFinished():
		return dur, dur, "finished"
	case t.kind == stateIdle:
		return dur, 0, "idle"
	case t.kind == stateRunning:
		return dur, msOf(t.previouslyElapsed + time.Since(t.startedAt)), "paused"
	case t.kind == statePaused:
		return dur, msOf(t.pausedElapsed), "paused"
	}
	return dur, 0, "idle"
}

func (c *TimerCollection) Snapshot() CollectionSnapshot {
	out := CollectionSnapshot{Timers: make([]TimerSnapshot, len(c.timers)), ActiveIndex: c.activeIndex}
	for i, nt := range c.timers {
		dur, elapsed, state := nt.Timer.snapshot()
		out.Timers[i] = TimerSnapshot{
			Name:               nt.Name,
			OriginalDurationMs: msOf(nt.Timer.originalDuration),
			DurationMs:         dur,
			ElapsedMs:          elapsed,
			State:              state,
			Tint:               nt.Tint,
		}
	}
	return out
}

func restoreTimer(s TimerSnapshot, defaultDurationSecs uint64) *CountdownTimer {
	original := time.Duration(s.OriginalDurationMs) * time.Millisecond
	if original == 0 {
		original = time.Duration(defaultDurationSecs) * time.Second
	}
	duration := time.Duration(s.DurationMs) * time.Millisecond
	if duration == 0 {
		duration = original
	}
	t := &CountdownTimer{originalDuration: original, duration: duration}
	switch s.State {
	case "paused":
		t.kind = statePaused
		t.pausedElapsed = time.Duration(s.ElapsedMs) * time.Millisecond
		if t.pausedElapsed >= t.duration {
			t.kind = stateFinished
		}
	case "finished":
		t.kind = stateFinished
	default:
		t.kind = stateIdle
	}
	return t
}

// RestoreCollection rebuilds a TimerCollection from a snapshot. An empty or
// invalid snapshot falls back to a fresh single-timer collection. Running
// timers in the snapshot are revived as paused.
func RestoreCollection(snap CollectionSnapshot, defaultDurationSecs uint64) *TimerCollection {
	if len(snap.Timers) == 0 {
		return NewCollection(defaultDurationSecs)
	}
	n := len(snap.Timers)
	if n > MaxTimers {
		n = MaxTimers
	}
	ts := make([]*NamedTimer, 0, n)
	for i := 0; i < n; i++ {
		ts = append(ts, &NamedTimer{
			Name:  snap.Timers[i].Name,
			Timer: restoreTimer(snap.Timers[i], defaultDurationSecs),
			Tint:  snap.Timers[i].Tint,
		})
	}
	active := snap.ActiveIndex
	if active < 0 || active >= len(ts) {
		active = 0
	}
	return &TimerCollection{timers: ts, activeIndex: active, breakIndex: -1}
}
