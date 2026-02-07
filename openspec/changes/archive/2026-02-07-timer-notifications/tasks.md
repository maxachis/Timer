## 1. Backend dependencies and configuration

- [x] 1.1 Add `tauri-plugin-notification` to `src-tauri/Cargo.toml` dependencies
- [x] 1.2 Register `tauri_plugin_notification::init()` in `src-tauri/src/lib.rs`
- [x] 1.3 Add `notification:default` permission to `src-tauri/capabilities/default.json`

## 2. Frontend dependencies

- [x] 2.1 Install `@tauri-apps/plugin-notification` npm package

## 3. Notification logic in App.svelte

- [x] 3.1 Add a `prevState` variable to track state transitions (detect when timer first becomes finished)
- [x] 3.2 Implement `playAlertSound()` using Web Audio API — short synthesized tone with gain envelope
- [x] 3.3 Implement `sendCompletionNotification()` — check/request permission, then send notification with title "Timer Complete"
- [x] 3.4 Call both `playAlertSound()` and `sendCompletionNotification()` on the finished transition (prevState !== "finished" && current === "finished")

## 4. Verify

- [x] 4.1 Run `cargo build` in src-tauri to verify notification plugin compiles
- [x] 4.2 Run `cargo test` in src-tauri to verify existing timer tests still pass
- [x] 4.3 Run `cargo tauri dev` and verify notification + sound fire on timer completion
