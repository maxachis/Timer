## 1. Tauri Permissions

- [x] 1.1 Add `core:window:allow-request-user-attention` to the permissions array in `src-tauri/capabilities/default.json`

## 2. Frontend Implementation

- [x] 2.1 Call `getCurrentWindow().requestUserAttention("critical")` in the finish-transition branch of `fetchStatus()` (alongside `playAlertSound` and `sendCompletionNotification`)
- [x] 2.2 Call `getCurrentWindow().requestUserAttention(null)` in `handleReset()` to cancel the flash when the timer is reset

## 3. Verification

- [x] 3.1 Build the app and confirm the taskbar icon flashes when the timer finishes
- [x] 3.2 Confirm the flash stops immediately after resetting the timer
