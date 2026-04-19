package main

import (
	"context"
	"sync"
	"time"

	"wails-app/internal/settings"
	"wails-app/internal/timer"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type TimerStatus struct {
	RemainingSecs float64 `json:"remaining_secs"`
	State         string  `json:"state"`
	IsFinished    bool    `json:"is_finished"`
	ActiveIndex   int     `json:"active_index"`
	ActiveName    string  `json:"active_name"`
	TimerCount    int     `json:"timer_count"`
}

type App struct {
	ctx      context.Context
	mu       sync.Mutex
	timers   *timer.TimerCollection
	settings settings.AppSettings
	store    *settings.Store
}

func NewApp() *App {
	store, err := settings.OpenStore()
	s := settings.Default()
	var names []string
	if err == nil {
		s = store.LoadSettings()
		names = store.LoadTimerNames()
	}
	return &App{
		timers:   timer.FromNames(names, s.DefaultDurationSecs),
		settings: s,
		store:    store,
	}
}

func (a *App) startup(ctx context.Context) { a.ctx = ctx }

func (a *App) domReady(ctx context.Context) {
	installMinSizeOverride(ctx, "Timer", 50, 50)
}

func (a *App) GetTimerStatus() TimerStatus {
	a.mu.Lock()
	defer a.mu.Unlock()
	t := a.timers.Active()
	return TimerStatus{
		RemainingSecs: t.Remaining().Seconds(),
		State:         t.StateName(),
		IsFinished:    t.IsFinished(),
		ActiveIndex:   a.timers.ActiveIndex(),
		ActiveName:    a.timers.ActiveName(),
		TimerCount:    a.timers.Count(),
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
		_ = a.store.SaveTimerNames(a.timers.TimerNames())
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
		_ = a.store.SaveTimerNames(a.timers.TimerNames())
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
		_ = a.store.SaveTimerNames(a.timers.TimerNames())
	}
	return nil
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
