package timer

import (
	"errors"
	"time"
)

const DefaultIncrement = 5 * time.Minute
const MaxTimers = 3

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

type NamedTimer struct {
	Name  string
	Timer *CountdownTimer
}

type TimerInfo struct {
	Index         int     `json:"index"`
	Name          string  `json:"name"`
	State         string  `json:"state"`
	RemainingSecs float64 `json:"remaining_secs"`
	IsActive      bool    `json:"is_active"`
}

type TimerCollection struct {
	timers      []*NamedTimer
	activeIndex int
}

func NewCollection(defaultDurationSecs uint64) *TimerCollection {
	return &TimerCollection{
		timers:      []*NamedTimer{{Name: "", Timer: NewCountdownTimer(defaultDurationSecs)}},
		activeIndex: 0,
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
	return &TimerCollection{timers: ts, activeIndex: 0}
}

func (c *TimerCollection) Active() *CountdownTimer { return c.timers[c.activeIndex].Timer }
func (c *TimerCollection) ActiveName() string      { return c.timers[c.activeIndex].Name }
func (c *TimerCollection) ActiveIndex() int        { return c.activeIndex }
func (c *TimerCollection) Count() int              { return len(c.timers) }

func (c *TimerCollection) AddTimer(name string, durationSecs uint64) (int, error) {
	if len(c.timers) >= MaxTimers {
		return 0, errors.New("Maximum of 3 timers reached. Remove a timer before adding another.")
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
	c.timers = append(c.timers[:index], c.timers[index+1:]...)
	if c.activeIndex >= len(c.timers) {
		c.activeIndex = len(c.timers) - 1
	} else if index < c.activeIndex {
		c.activeIndex--
	}
	return nil
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
		out[i] = TimerInfo{
			Index:         i,
			Name:          t.Name,
			State:         t.Timer.StateName(),
			RemainingSecs: t.Timer.Remaining().Seconds(),
			IsActive:      i == c.activeIndex,
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
