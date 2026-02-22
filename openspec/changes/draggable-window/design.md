## Context

The Timer app is a Tauri 2 desktop application with a Svelte 5 frontend. The window defaults to 400x600 and can be resized down to 100x100. At small sizes the OS title bar becomes a difficult drag target. The window body is mostly non-interactive space around the timer face and progress ring.

## Goals / Non-Goals

**Goals:**
- Allow the user to drag the window by clicking and holding on any non-interactive area of the UI.
- Buttons, inputs, and other interactive elements remain fully clickable without triggering a drag.

**Non-Goals:**
- Removing the native title bar / going to a custom title bar.
- Adding a custom drag handle widget.

## Decisions

### Use Tauri's `data-tauri-drag-region` HTML attribute

Tauri 2 natively supports marking HTML elements as window drag regions via the `data-tauri-drag-region` attribute. When the user clicks and drags on a marked element, Tauri initiates a native window move. Child elements that are interactive (buttons, inputs) naturally consume click events and do not trigger the drag.

**Why this over alternatives:**
- **vs. calling `startDragging()` from JS**: The attribute approach requires zero JS, no event listeners, and no edge-case handling. It's declarative and built into Tauri.
- **vs. custom title bar with `decorations: false`**: Far more invasive, loses native minimize/maximize/close controls, and is out of scope.

### Apply the attribute to the outermost container

Place `data-tauri-drag-region` on the root `.container` element in `App.svelte`. This makes the entire window body a drag region by default. Interactive children (buttons, inputs, links) automatically exclude themselves because they handle their own pointer events.

### Add `core:window:allow-start-dragging` permission

Tauri 2 requires explicit capability permissions. The `data-tauri-drag-region` attribute triggers `startDragging` internally, so the permission must be granted in `capabilities/default.json`.

## Risks / Trade-offs

- **Text selection disabled in drag regions**: `data-tauri-drag-region` prevents text selection on drag-region elements. This is acceptable since the timer UI has no user-selectable text.
- **Future interactive elements**: Any new interactive element added to the UI will need to properly handle pointer events to avoid being swallowed by the drag region. Standard HTML buttons and inputs do this automatically.
