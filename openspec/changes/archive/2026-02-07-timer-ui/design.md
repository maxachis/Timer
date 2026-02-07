## Context

The Timer app has a working Rust `CountdownTimer` module and a Tauri 2 shell with a
placeholder Svelte frontend. The timer logic (start, pause, resume, reset, add/remove
time, completion detection) is fully implemented and tested in `src-tauri/src/timer.rs`.
The frontend currently shows only an `<h1>Timer</h1>` placeholder.

We need to bridge the Rust logic to the Svelte frontend via Tauri commands, then build
an interactive UI.

## Goals / Non-Goals

**Goals:**
- Tauri commands exposing all timer operations to the frontend
- Thread-safe timer state managed by Tauri's state system
- Clean countdown display (MM:SS) that updates in real-time while running
- Controls: start/pause toggle, reset, +5 min, -5 min
- Visual indication when timer finishes
- Default timer duration of 5 minutes on app launch

**Non-Goals:**
- Custom duration input (user can't type an arbitrary time — they adjust with +/- buttons)
- Multiple concurrent timers
- Sound or system notifications on completion (separate change)
- Persistent timer state across app restarts
- Theming or dark mode

## Decisions

### Decision 1: Tauri managed state with Mutex

**Approach:** Wrap `CountdownTimer` in `Mutex<CountdownTimer>` and register it via
`tauri::Builder::manage()`. Each Tauri command locks the mutex, performs the operation,
and returns.

**Why not channels or async?** The timer operations are instant (no I/O, no waiting).
A simple mutex lock-unlock per command is the simplest correct approach. No risk of
contention since the frontend is the only caller.

### Decision 2: Polling for time updates (frontend setInterval)

**Approach:** The frontend uses `setInterval` (100ms interval) to call a
`get_timer_status` command that returns remaining seconds and the current state
(idle/running/paused/finished). The interval runs only while the timer is running.

**Why not events/push?** Tauri supports events, but polling is simpler for a single
timer with one consumer. A 100ms poll gives smooth visual updates without
complexity. The overhead is negligible — each poll is a mutex lock + duration
calculation.

### Decision 3: Single status command returning a struct

**Approach:** Instead of separate `get_remaining` and `get_state` commands, use one
`get_timer_status` command returning `{ remaining_secs: f64, state: String, is_finished: bool }`.
This avoids race conditions from two separate calls and reduces IPC round trips.

### Decision 4: State-driven UI with Svelte 5 reactivity

**Approach:** Use Svelte 5 `$state` runes to track timer status. Button visibility
and labels change based on state:
- Idle → show "Start" button, +/- buttons, reset disabled
- Running → show "Pause" button, +/- buttons, reset enabled
- Paused → show "Resume" button, +/- buttons, reset enabled
- Finished → show "Reset" button, finished message, +/- disabled

### Decision 5: Timer created at app startup

**Approach:** Initialize a `CountdownTimer::new(300)` (5 minutes) when the Tauri app
starts, stored in managed state. A `create_timer` command allows reinitializing with
a different duration if needed later.

## Risks / Trade-offs

- [100ms polling] → Acceptable CPU cost for a desktop app. If battery concern arises,
  increase interval to 250ms or switch to events.
- [Mutex contention] → Impossible in practice with a single frontend caller and
  instant operations. No mitigation needed.
- [No duration input] → Users can only set time in 5-minute increments. Acceptable
  for v1 — custom input can be added later.
