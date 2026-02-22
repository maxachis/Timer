## Context

The timer app currently fires a one-shot audio beep and system notification when the countdown reaches zero. If the user is away or misses both, there is no persistent indicator that the timer has finished. The app already imports `getCurrentWindow` from `@tauri-apps/api/window` for the always-on-top feature.

Tauri 2 exposes `WebviewWindow.requestUserAttention(type)` on the frontend. Passing `"critical"` makes the taskbar icon flash persistently (Windows) or bounce the dock icon (macOS). Passing `null` cancels the attention request.

## Goals / Non-Goals

**Goals:**
- Flash the taskbar icon when the timer finishes, persisting until the user resets
- Cancel the flash immediately when the timer is reset
- Minimal code change — leverage the existing `getCurrentWindow()` import

**Non-Goals:**
- Custom tray icon or overlay badge
- Flashing for any state other than "finished"
- Making flash behavior configurable in settings (can be added later)

## Decisions

### Use `requestUserAttention("critical")` from the frontend

The Tauri window API is already available in the frontend (`getCurrentWindow` is imported). Calling `requestUserAttention` with the `"critical"` attention type provides persistent flashing on Windows and a dock bounce on macOS. No Rust-side changes are needed.

**Alternative considered**: Implementing attention via a Rust command — rejected because the frontend already has the window handle and a direct call is simpler.

### Trigger on the existing finish-detection branch

The `fetchStatus()` function already detects the `idle/running/paused → finished` transition (line 122 in App.svelte). The attention request will be added alongside the existing `playAlertSound()` and `sendCompletionNotification()` calls.

### Cancel attention on reset

`handleReset()` is the single code path for returning from "finished" to "idle". Calling `requestUserAttention(null)` there ensures the flash stops regardless of how reset is triggered (button click, spacebar, Escape key).

## Risks / Trade-offs

- **Platform variance** → On Linux/Wayland, `requestUserAttention` may be a no-op depending on the compositor. This is acceptable — it's a best-effort enhancement.
- **Permission requirement** → Tauri 2 requires `core:window:allow-request-user-attention` in the capability config. If missing, the call silently fails. Mitigation: add the permission to `default.json`.
