## Context

The timer window currently behaves like any other window — it goes behind other apps when the user switches focus. Tauri 2 provides a built-in `setAlwaysOnTop()` API on the window object, callable from the frontend with no Rust-side changes.

## Goals / Non-Goals

**Goals:**
- Toggle button in the UI to pin/unpin the window
- `T` keyboard shortcut to toggle
- Visual indicator showing current pin state
- State defaults to unpinned on app launch

**Non-Goals:**
- Persisting the pin state across app restarts (can be added later with a settings store)
- Mini/compact mode when pinned (separate feature)
- System tray integration

## Decisions

### Decision 1: Frontend-only implementation

**Approach:** Use `getCurrentWindow().setAlwaysOnTop(true/false)` from `@tauri-apps/api/window`. Track the pinned state with a Svelte `$state` variable. No Tauri commands or Rust code needed.

**Why not a Tauri command?** The window API is already exposed to the frontend. Wrapping it in a custom command would add unnecessary indirection.

### Decision 2: Small pin button in the top-right corner

**Approach:** Place a small pin icon button in the top-right corner of the app, outside the main timer layout. When active, the icon visually changes (e.g., filled vs outline) to indicate the pinned state.

**Why top-right?** It's the conventional location for window utility controls. It stays out of the way of the timer display and action buttons.

### Decision 3: T key for toggle shortcut

**Approach:** Add `t`/`T` to the existing `handleKeydown` switch statement. It toggles always-on-top regardless of timer state.

**Why T?** Mnemonic for "top". Not conflicting with existing shortcuts (Space, R, arrows, +/-, Escape).

## Risks / Trade-offs

- [No persistence] → The pin state resets on app restart. Acceptable for now — users can press T again. Persistence would require a settings store (e.g., `tauri-plugin-store`).
