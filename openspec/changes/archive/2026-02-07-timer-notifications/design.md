## Context

The Timer app has a working countdown UI that polls `get_timer_status` every 80ms while
running. When the timer finishes, the display turns amber and shows "complete" — but only
if the user is looking at the window. There is no notification or sound to alert a user
who has switched focus or minimized the app.

Tauri 2 provides an official `tauri-plugin-notification` for OS-level toast notifications.
For audio, the Web Audio API in the webview is the simplest approach — no Rust-side audio
crate needed.

## Goals / Non-Goals

**Goals:**
- OS-level toast notification when the timer finishes (works even when app is minimized)
- Audible alert sound on completion
- Both triggered from the frontend when the finished state is first detected
- Request notification permission on first use

**Non-Goals:**
- Custom notification sounds (OS default is fine)
- Notification actions (snooze, restart from notification)
- Sound settings or mute toggle (can be added later)
- Repeating/persistent alarms
- Notification when app is fully closed (requires background service)

## Decisions

### Decision 1: Trigger from frontend polling, not backend

**Approach:** The frontend already polls `get_timer_status` every 80ms. When it detects
the transition to `finished` state (previous state was not finished, current state is),
it fires both the notification and the sound.

**Why not backend-initiated?** The backend would need a separate timer thread to detect
completion and emit a Tauri event. The frontend already has the polling loop with all
the state — adding a transition check there is one `if` statement. No new architecture.

### Decision 2: Web Audio API for sound, not a Rust audio crate

**Approach:** Generate a short alert tone programmatically using the Web Audio API
(oscillator + gain envelope). No audio file to bundle.

**Why not an audio file?** Bundling an MP3/WAV adds asset management complexity and
file size. A synthesized tone is ~15 lines of code, has zero dependencies, and sounds
clean. It also avoids cross-platform codec issues.

**Why not a Rust audio crate?** Adding `rodio` or similar would bloat the binary and
add native audio complexity. The webview already has a working audio stack.

### Decision 3: Use tauri-plugin-notification for OS toasts

**Approach:** Add `tauri-plugin-notification` and `@tauri-apps/plugin-notification`.
Call `sendNotification()` from the frontend JS when the timer finishes. Request
permission once on first completion attempt.

**Why not tauri-plugin-dialog?** Dialogs are modal and block the app — wrong UX for
a background alert. Toast notifications are non-blocking and appear system-wide.

### Decision 4: Permission request on first completion

**Approach:** Check `isPermissionGranted()` when the timer first finishes. If not
granted, call `requestPermission()`. Cache the result so subsequent completions
don't re-prompt.

**Why not request on app launch?** Requesting permission before the user has context
("why does this timer app want notifications?") leads to denials. Requesting at the
moment of completion — when the value is obvious — gets better grant rates.

## Risks / Trade-offs

- [Dev mode notification appearance] → On Windows, notifications in dev mode show
  as PowerShell rather than the app name. This is a Tauri limitation, fixed in
  production builds.
- [Autoplay restrictions] → Web Audio API requires prior user interaction before
  playing. Since the user clicked "Start" to begin the timer, this is satisfied.
- [No sound if app is minimized] → Some OS/webview combos may not play audio from
  a minimized window. The system notification serves as the backup alert.
