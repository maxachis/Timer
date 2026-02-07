## Context

The timer app has button-based controls for all actions (start, pause, resume, reset, add/remove time). Users must click specific buttons to interact. Adding keyboard shortcuts makes the app faster to use, especially when switching back from another window.

All timer actions already exist as `handle*` functions in `App.svelte`. The shortcuts just need to dispatch to them based on current state.

## Goals / Non-Goals

**Goals:**
- Keyboard shortcuts for all timer actions (start, pause, resume, reset, add/remove time)
- Context-aware Space bar (maps to the correct action based on timer state)
- Arrow keys as intuitive alternatives (left/right for time adjust, down for reset)

**Non-Goals:**
- Configurable key bindings (hardcoded is fine for now)
- Visual shortcut hints on buttons (can be added later)
- Global/system-wide shortcuts when app is not focused (would require Tauri global shortcut plugin)

## Decisions

### Decision 1: `svelte:window` keydown binding

**Approach:** Use Svelte's `<svelte:window on:keydown={handleKeydown}>` to listen for keyboard events at the window level. A single handler function maps `event.key` to the appropriate `handle*` function based on current `timerState`.

**Why not per-button key listeners?** A single centralized handler is simpler and avoids focus management issues. The app has no text inputs, so there's no conflict risk.

### Decision 2: Context-aware Space bar

**Approach:** Space maps to the "primary action" — the same action as the center button:
- `idle` → `handleStart()`
- `running` → `handlePause()`
- `paused` → `handleResume()`
- `finished` → `handleReset()`

**Why not separate keys for each?** The primary button already changes meaning based on state. Space mirrors that behavior, so the user only needs to remember one key.

### Decision 3: Prevent default on handled keys

**Approach:** Call `event.preventDefault()` for all handled keys to suppress browser defaults (e.g., Space scrolling the page, arrow keys moving focus).

## Risks / Trade-offs

- [No text inputs currently] → If text inputs are added later (e.g., custom duration), the global keydown handler will need to check `event.target` to avoid capturing keystrokes meant for inputs. Low risk since no inputs exist today.
