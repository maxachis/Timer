## Why

The timer window is small and mostly empty space around the timer face. Users should be able to grab anywhere on the window body (not just the title bar) to drag it around the screen. This is especially useful when the window is resized very small and the title bar becomes a tiny drag target.

## What Changes

- Add Tauri's `data-tauri-drag-region` attribute to the main container element so that clicking and dragging on non-interactive areas moves the window.
- Interactive elements (buttons, inputs) remain clickable and do not trigger window dragging.

## Capabilities

### New Capabilities
- `window-drag-region`: Defines which areas of the window act as drag handles, allowing the user to move the window by clicking and dragging on non-button areas of the UI.

### Modified Capabilities

_(none — no existing spec requirements change)_

## Impact

- `src/App.svelte`: Add `data-tauri-drag-region` to the outer container. CSS may need `user-select: none` on drag regions.
- `src-tauri/capabilities/default.json`: May need `core:window:allow-start-dragging` permission.
