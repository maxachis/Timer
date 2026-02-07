## 1. Permissions

- [x] 1.1 Add `core:window:allow-set-always-on-top` and `core:window:allow-is-always-on-top` to `src-tauri/capabilities/default.json`

## 2. Frontend implementation

- [x] 2.1 Add `alwaysOnTop` state variable and `toggleAlwaysOnTop()` function using `getCurrentWindow().setAlwaysOnTop()`
- [x] 2.2 Add `T` key case to `handleKeydown()` that calls `toggleAlwaysOnTop()`
- [x] 2.3 Add pin toggle button to the UI (top-right corner, outline/filled icon based on state)

## 3. Verify

- [x] 3.1 Run `npx vite build` to verify frontend compiles
- [x] 3.2 Run `cargo tauri dev` and verify pin button and T shortcut toggle always-on-top
