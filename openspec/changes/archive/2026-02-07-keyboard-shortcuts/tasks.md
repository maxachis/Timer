## 1. Keyboard handler implementation

- [x] 1.1 Add `handleKeydown(event)` function in `App.svelte` that maps keys to existing handler functions based on timer state
- [x] 1.2 Add `<svelte:window>` keydown binding in the template to wire up the handler

## 2. Verify

- [x] 2.1 Run `npx vite build` to verify frontend compiles
- [x] 2.2 Run `cargo tauri dev` and verify all keyboard shortcuts work
