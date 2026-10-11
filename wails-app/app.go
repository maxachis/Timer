package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"wails-app/internal/history"
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
	// history is derived by diffing timer state at each sync point (every
	// mutation and every status poll), so no mutator has to name its event.
	historyLog     *history.Log
	historyTracker *history.Tracker[*timer.CountdownTimer]
}

func NewApp() *App {
	store, err := settings.OpenStore()
	if err != nil {
		store = nil
	}
	return newAppFromStore(store)
}

// newAppFromStore builds the app from a store (nil when the config dir is
// unavailable, in which case nothing persists).
func newAppFromStore(store *settings.Store) *App {
	s := settings.Default()
	var collection *timer.TimerCollection
	var log *history.Log
	if store != nil {
		s = store.LoadSettings()
		var snap timer.CollectionSnapshot
		if ok, err := store.LoadTimerStates(&snap); err == nil && ok {
			collection = timer.RestoreCollection(snap, s.DefaultDurationSecs)
		} else {
			collection = timer.FromNames(store.LoadTimerNames(), s.DefaultDurationSecs)
		}
		var entries []history.Entry
		if ok, err := store.LoadHistory(&entries); err == nil && ok {
			log = history.FromEntries(entries)
		}
	} else {
		collection = timer.FromNames(nil, s.DefaultDurationSecs)
	}
	if log == nil {
		log = history.NewLog()
	}
	a := &App{
		timers:     collection,
		settings:   s,
		store:      store,
		historyLog: log,
	}
	a.historyTracker = history.NewTracker[*timer.CountdownTimer](log)
	a.historyTracker.Seed(a.observedTimers())
	return a
}

// observedTimers is the tracker's view of the collection: identity, display
// name, and run state per timer.
func (a *App) observedTimers() []history.Observed[*timer.CountdownTimer] {
	names := a.timers.TimerNames()
	out := make([]history.Observed[*timer.CountdownTimer], len(names))
	for i, name := range names {
		t := a.timers.TimerAt(i)
		out[i] = history.Observed[*timer.CountdownTimer]{ID: t, Name: displayName(name, i), State: t.StateName()}
	}
	return out
}

// displayName matches the frontend's label for an unnamed timer.
func displayName(name string, index int) string {
	if name == "" {
		return fmt.Sprintf("Timer %d", index+1)
	}
	return name
}

// syncHistory logs whatever changed since the last sync and persists the log
// if anything did. Callers hold a.mu.
func (a *App) syncHistory() {
	if !a.historyTracker.Sync(time.Now(), a.observedTimers()) {
		return
	}
	if a.store != nil {
		_ = a.store.SaveHistory(a.historyLog.Entries())
	}
}

// commitChange runs after every timer mutation: it logs any history the
// change produced and saves the timer collection, so the change survives a
// crash rather than waiting for the next autosave. Callers hold a.mu.
func (a *App) commitChange() {
	a.syncHistory()
	if a.store != nil {
		_ = a.store.SaveTimerStates(a.timers.Snapshot())
	}
}

func (a *App) GetHistory() []history.Entry {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.historyLog.Entries()
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
	a.syncHistory()
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

func (a *App) StartTimer() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.Active().Start()
}
func (a *App) PauseTimer() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.Active().Pause()
}
func (a *App) ResumeTimer() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.Active().Resume()
}
func (a *App) ResetTimer() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.Active().Reset()
}

func (a *App) AddTime() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.Active().AddTime(time.Duration(a.settings.DefaultIncrementSecs) * time.Second)
}

func (a *App) RemoveTime() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.Active().RemoveTime(time.Duration(a.settings.DefaultIncrementSecs) * time.Second)
}

func (a *App) AddTimeSecondary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.Active().AddTime(time.Duration(a.settings.SecondaryIncrementSecs) * time.Second)
}

func (a *App) RemoveTimeSecondary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.Active().RemoveTime(time.Duration(a.settings.SecondaryIncrementSecs) * time.Second)
}

func (a *App) AddTimeTertiary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.Active().AddTime(time.Duration(a.settings.TertiaryIncrementSecs) * time.Second)
}

func (a *App) RemoveTimeTertiary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.Active().RemoveTime(time.Duration(a.settings.TertiaryIncrementSecs) * time.Second)
}

func (a *App) AddTimeCustom(seconds uint64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.Active().AddTime(time.Duration(seconds) * time.Second)
}

func (a *App) RemoveTimeCustom(seconds uint64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.Active().RemoveTime(time.Duration(seconds) * time.Second)
}

func (a *App) CreateTimer(durationSecs uint64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	old := a.timers.Active()
	a.timers.ReplaceActive(durationSecs)
	a.historyTracker.Rekey(old, a.timers.Active())
}

func (a *App) AddNewTimer(name string) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	idx, err := a.timers.AddTimer(name, a.settings.DefaultDurationSecs)
	if err != nil {
		return 0, err
	}
	return idx, nil
}

func (a *App) RemoveExistingTimer(index int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	if err := a.timers.RemoveTimer(index); err != nil {
		return err
	}
	return nil
}

func (a *App) SwitchTimer(index int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	return a.timers.SwitchTo(index)
}

func (a *App) SwitchTimerNext() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.SwitchNext()
}
func (a *App) SwitchTimerPrev() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.SwitchPrev()
}

func (a *App) RenameTimer(index int, name string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	if err := a.timers.RenameTimer(index, name); err != nil {
		return err
	}
	return nil
}

func (a *App) SetTimerTint(index int, tint string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	if err := a.timers.SetTint(index, tint); err != nil {
		return err
	}
	return nil
}

func (a *App) ToggleBreakTimer() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.ToggleBreakDesignation()
}

// SetMaxTime commits the active timer's current time as its maximum, so that
// resetting returns to it and the stacked bar scales against it.
func (a *App) SetMaxTime() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.SetMaxToCurrent()
}

func (a *App) ToggleBreak() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	return a.timers.ToggleBreak()
}

func (a *App) AddBreakTime() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.AddBreakTime(time.Duration(a.settings.DefaultIncrementSecs) * time.Second)
}

func (a *App) RemoveBreakTime() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.RemoveBreakTime(time.Duration(a.settings.DefaultIncrementSecs) * time.Second)
}

func (a *App) AddBreakTimeSecondary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.AddBreakTime(time.Duration(a.settings.SecondaryIncrementSecs) * time.Second)
}

func (a *App) RemoveBreakTimeSecondary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.RemoveBreakTime(time.Duration(a.settings.SecondaryIncrementSecs) * time.Second)
}

func (a *App) AddBreakTimeTertiary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.AddBreakTime(time.Duration(a.settings.TertiaryIncrementSecs) * time.Second)
}

func (a *App) RemoveBreakTimeTertiary() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
	a.timers.RemoveBreakTime(time.Duration(a.settings.TertiaryIncrementSecs) * time.Second)
}

func (a *App) ResetBreakTime() {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer a.commitChange()
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
