## Why

The timer window is fixed at 400x600 and the UI uses hardcoded pixel sizes, so it cannot be resized to fit different use cases — a compact corner widget, a small always-on-top overlay, or a larger display. Users should be able to resize the window freely (down to 100x100) and still see the core countdown, with the UI gracefully adapting what it shows at each size.

## What Changes

- Set minimum window size to 100x100 in Tauri config
- Make the timer face (progress ring + digits) scale fluidly with the window
- At very small sizes (below ~200px wide), hide secondary controls (adjust buttons, reset) and show only the time display and primary action button
- At compact sizes (~200–350px), show a simplified layout without the progress ring tick marks
- At default size (400x600) and above, display the full current UI unchanged
- Settings view adapts similarly — simplified at small sizes
- All sizing is CSS-driven (viewport units, clamp, container queries or media queries), no JavaScript resize listeners needed

## Capabilities

### New Capabilities
- `responsive-layout`: Fluid scaling rules and breakpoint-based visibility for the timer UI, defining which elements appear at which window sizes and how dimensions adapt

### Modified Capabilities
- `tauri-app-shell`: Add minimum window size constraint (100x100) to the window configuration requirement
- `timer-display`: Add responsive sizing behavior — the time display, progress ring, and controls must scale and show/hide based on available space

## Impact

- **Frontend**: `src/App.svelte` — replace hardcoded pixel sizes with fluid/responsive values; add breakpoint-based visibility classes
- **Tauri config**: `src-tauri/tauri.conf.json` — add `minWidth`/`minHeight` to the window definition
- **No backend changes** — all changes are purely presentational
- **No dependency additions** — uses standard CSS features
