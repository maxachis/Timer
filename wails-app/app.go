package main

import (
	"context"
	"sync"
	"time"

	"wails-app/internal/settings"
	"wails-app/internal/timer"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const autosaveInterval = 30 * time.Second

type TimerStatus struct {
	RemainingSecs float64 `json:"remaining_secs"`
	State         string  `json:"state"`
	IsFinished    bool    `json:"is_finished"`
	ActiveIndex   int     `json:"active_index"`
	ActiveName    string  `json:"active_name"`
	TimerCount    int     `json:"timer_count"`
	OnBreak       bool    `json:"on_break"`
	BreakIndex    int     `json:"break_index"`
}

type App struct {
	ctx            context.Context
	mu             sync.Mutex
	timers         *timer.TimerCollection
	settings       settings.AppSettings
	store          *settings.Store
	cancelAutosave context.CancelFunc
}

func NewApp() *App {
	store, err := settings.OpenStore()
	s := settings.Default()
	var collection *timer.TimerCollection
	if err == nil {
		s = store.LoadSettings()
		var snap timer.CollectionSnapshot
		if ok, err := store.LoadTimerStates(&snap); err == nil && ok {
			collection = timer.RestoreCollection(snap, s.DefaultDurationSecs)
		} else {
			collection = timer.FromNames(store.LoadTimerNames(), s.DefaultDurationSecs)
		}
	} else {
		collection = timer.FromNames(nil, s.DefaultDurationSecs)
	}
	return &App{
		timers:   collection,
		settings: s,
		store:    store,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.cancelAutosave = a.startAutosave(autosaveInterval)
}

func (a *App) startAutosave(interval time.Duration) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				a.mu.Lock()
				if a.store != nil {
					_ = a.store.SaveTimerStates(a.timers.Snapshot())
				}
				a.mu.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()
	return cancel
}

func (a *App) domReady(ctx context.Context) {
	installMinSizeOverride(ctx, "Timer", 50, 50)
}

// PersistTimerStates writes the current timer collection snapshot to disk.
// Called from the Wails OnShutdown hook so that elapsed time and per-timer
// state survive app restarts.
func (a *App) PersistTimerStates() {
	if a.cancelAutosave != nil {
		a.cancelAutosave()
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.store == nil {
		return
	}
	_ = a.store.SaveTimerStates(a.timers.Snapshot())
}

func (a *App) GetTimerStatus() TimerStatus {
	a.mu.Lock()
	defer a.mu.Unlock()
	t := a.timers.Active()
	remaining := t.Remaining()
	// While on break the countdown is the break allocation left within the
	// block, not the block's own remaining time.
	if a.timers.IsOnBreak() {
		remaining = a.timers.BreakAccrued()
	}
	return TimerStatus{
		RemainingSecs: remaining.Seconds(),
		State:         t.StateName(),
		IsFinished:    t.IsFinished(),
		ActiveIndex:   a.timers.ActiveIndex(),
		ActiveName:    a.timers.ActiveName(),
		TimerCount:    a.timers.Count(),
		OnBreak:       a.timers.IsOnBreak(),
		BreakIndex:    a.timers.BreakIndex(),
	}
}

func (a *App) GetTimerList() []timer.TimerInfo {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.timers.TimerList()
}

func (a *App) StartTimer()  { a.mu.Lock(); defer a.mu.Unlock(); a.timers.Active().Start() }
func (a *App) PauseTimer()  { a.mu.Lock(); defer a.mu.Unlock(); a.timers.Active().Pause() }
func (a *App) ResumeTimer() { a.mu.Lock(); defer a.mu.Unlock(); a.timers.Active().Resume() }
func (a *App) ResetTimer()  { a.mu.Lock(); defer a.mu.Unlock(); a.timers.Active().Reset() }

func (a *App) AddTime() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.Active().AddTime(time.Duration(a.settings.DefaultIncrementSecs) * time.Second)
}

func (a *App) RemoveTime() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.Active().RemoveTime(time.Duration(a.settings.DefaultIncrementSecs) * time.Second)
}

func (a *App) AddTimeSecondary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.Active().AddTime(time.Duration(a.settings.SecondaryIncrementSecs) * time.Second)
}

func (a *App) RemoveTimeSecondary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.Active().RemoveTime(time.Duration(a.settings.SecondaryIncrementSecs) * time.Second)
}

func (a *App) AddTimeTertiary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.Active().AddTime(time.Duration(a.settings.TertiaryIncrementSecs) * time.Second)
}

func (a *App) RemoveTimeTertiary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.Active().RemoveTime(time.Duration(a.settings.TertiaryIncrementSecs) * time.Second)
}

func (a *App) AddTimeCustom(seconds uint64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.Active().AddTime(time.Duration(seconds) * time.Second)
}

func (a *App) RemoveTimeCustom(seconds uint64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.Active().RemoveTime(time.Duration(seconds) * time.Second)
}

func (a *App) CreateTimer(durationSecs uint64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.ReplaceActive(durationSecs)
}

func (a *App) AddNewTimer(name string) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	idx, err := a.timers.AddTimer(name, a.settings.DefaultDurationSecs)
	if err != nil {
		return 0, err
	}
	if a.store != nil {
		_ = a.store.SaveTimerStates(a.timers.Snapshot())
	}
	return idx, nil
}

func (a *App) RemoveExistingTimer(index int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.timers.RemoveTimer(index); err != nil {
		return err
	}
	if a.store != nil {
		_ = a.store.SaveTimerStates(a.timers.Snapshot())
	}
	return nil
}

func (a *App) SwitchTimer(index int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.timers.SwitchTo(index)
}

func (a *App) SwitchTimerNext() { a.mu.Lock(); defer a.mu.Unlock(); a.timers.SwitchNext() }
func (a *App) SwitchTimerPrev() { a.mu.Lock(); defer a.mu.Unlock(); a.timers.SwitchPrev() }

func (a *App) RenameTimer(index int, name string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.timers.RenameTimer(index, name); err != nil {
		return err
	}
	if a.store != nil {
		_ = a.store.SaveTimerStates(a.timers.Snapshot())
	}
	return nil
}

func (a *App) SetTimerTint(index int, tint string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.timers.SetTint(index, tint); err != nil {
		return err
	}
	if a.store != nil {
		_ = a.store.SaveTimerStates(a.timers.Snapshot())
	}
	return nil
}

func (a *App) ToggleBreakTimer() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.ToggleBreakDesignation()
	if a.store != nil {
		_ = a.store.SaveTimerStates(a.timers.Snapshot())
	}
}

// SetMaxTime commits the active timer's current time as its maximum, so that
// resetting returns to it and the stacked bar scales against it.
func (a *App) SetMaxTime() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.SetMaxToCurrent()
	if a.store != nil {
		_ = a.store.SaveTimerStates(a.timers.Snapshot())
	}
}

func (a *App) ToggleBreak() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.timers.ToggleBreak()
}

func (a *App) AddBreakTime() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.AddBreakTime(time.Duration(a.settings.DefaultIncrementSecs) * time.Second)
}

func (a *App) RemoveBreakTime() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.RemoveBreakTime(time.Duration(a.settings.DefaultIncrementSecs) * time.Second)
}

func (a *App) AddBreakTimeSecondary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.AddBreakTime(time.Duration(a.settings.SecondaryIncrementSecs) * time.Second)
}

func (a *App) RemoveBreakTimeSecondary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.RemoveBreakTime(time.Duration(a.settings.SecondaryIncrementSecs) * time.Second)
}

func (a *App) AddBreakTimeTertiary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.AddBreakTime(time.Duration(a.settings.TertiaryIncrementSecs) * time.Second)
}

func (a *App) RemoveBreakTimeTertiary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.RemoveBreakTime(time.Duration(a.settings.TertiaryIncrementSecs) * time.Second)
}

func (a *App) ResetBreakTime() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timers.ResetBreakTime()
}

func (a *App) GetSettings() settings.AppSettings {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.settings
}

func (a *App) SetAlwaysOnTop(on bool) {
	wruntime.WindowSetAlwaysOnTop(a.ctx, on)
}

func (a *App) SetWindowSize(width, height int) {
	wruntime.WindowSetSize(a.ctx, width, height)
}

func (a *App) SaveSettings(s settings.AppSettings) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := settings.Validate(s); err != nil {
		return err
	}
	if a.store != nil {
		if err := a.store.SaveSettings(s); err != nil {
			return err
		}
	}
	a.settings = s
	return nil
}
