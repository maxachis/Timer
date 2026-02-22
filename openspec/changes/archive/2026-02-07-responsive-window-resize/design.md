## Context

The timer app runs in a dedicated Tauri window at a fixed 400x600. The UI is a single Svelte component (`App.svelte`) with a 300x300 SVG progress ring, fixed-size text, and absolutely-positioned toolbar buttons. Because every dimension is hardcoded in pixels, the window cannot be resized meaningfully. The user wants the window to scale from 100x100 up to any size, adapting which elements are visible and how large they are.

The SVG progress ring already uses a `viewBox` so it scales for free once its container is fluid. The main work is making the container, text, and controls respond to available space.

## Goals / Non-Goals

**Goals:**
- Window resizable from 100x100 up, with the current 400x600 as the default
- Core time display always visible at every size
- Graceful progressive disclosure: more space → more UI elements
- Pure CSS solution — no JS resize observers or listeners
- Current default appearance is pixel-identical at 400x600

**Non-Goals:**
- Mobile/touch layout — this is a desktop Tauri app
- Persistent user-chosen size — Tauri already remembers window size on restart if configured; not our concern here
- Redesigning the settings view for tiny sizes — settings can simply require a minimum usable size (the back button and fields already work at ~250px wide)

## Decisions

### 1. Scaling unit: `vmin` with `clamp()`

**Decision:** Size the timer face and text relative to `vmin` (the smaller of viewport width/height), bounded by `clamp()`.

**Why over alternatives:**
- *Media queries alone* → step-wise jumps; doesn't feel fluid
- *Container queries* → the app is full-viewport, so container = viewport; adds complexity for no gain
- *JS resize observer* → unnecessary overhead; CSS handles this natively
- `vmin` ensures the timer face fits whether the window is wide-and-short or narrow-and-tall

**Specifics:**
- Timer face: `clamp(80px, 70vmin, 300px)` for width/height
- Time digits: `clamp(1.2rem, 12vmin, 4.5rem)`
- Separator: scale proportionally with digits
- Primary button: `clamp(32px, 14vmin, 64px)`

### 2. Progressive visibility via media queries

**Decision:** Use `@media` queries on viewport dimensions to show/hide secondary elements at three tiers.

| Tier | Condition | Visible elements |
|------|-----------|-----------------|
| Tiny | `max-width: 179px` or `max-height: 179px` | Time digits + primary button only (no ring, no adjust buttons, no toolbar) |
| Compact | `min-width: 180px` and `min-height: 180px` up to ~349px | Ring + digits + primary button; hide tick marks, adjust buttons, reset text, toolbar buttons |
| Full | `min-width: 350px` and `min-height: 400px` | Everything (current UI) |

**Why media queries for visibility:** They're declarative, easy to maintain, and naturally work with the standalone Tauri window viewport. Container queries would work identically here since the app fills 100vh/100vw but add syntactic overhead.

### 3. Tauri config: `minWidth` / `minHeight`

**Decision:** Add `"minWidth": 100` and `"minHeight": 100` to the window config in `tauri.conf.json`. Keep `"width": 400` and `"height": 600` as defaults.

No Rust backend changes needed — this is a config-only change.

### 4. Layout structure: keep existing, add fluid overrides

**Decision:** Keep the existing CSS as-is and add responsive overrides below it, rather than rewriting the layout.

**Why:** This preserves pixel-identical behavior at default size and isolates the responsive logic in one clearly marked section at the bottom of the `<style>` block.

### 5. Toolbar buttons at small sizes

**Decision:** Hide the pin and settings buttons below the compact threshold. At small sizes the user interacts via keyboard shortcuts (already implemented: T for pin toggle). The settings gear reappears once the window is large enough for it.

## Risks / Trade-offs

- **Fonts at tiny sizes:** Below ~120px, the time digits may become hard to read even when scaled. → Mitigation: `clamp()` ensures a minimum readable font size of 1.2rem (~19px). At 100x100 the display is intentionally minimal.
- **Progress ring clipping:** If the ring container shrinks below the stroke width, it could clip. → Mitigation: At tiny tier the ring is hidden entirely; at compact tier the minimum size (80px) is well above the stroke width.
- **Settings view at small sizes:** The settings panel has fixed-width inputs (200px) that won't fit at 100px. → Mitigation: non-goal; settings require ≥250px width which is reasonable. The user can resize up, adjust settings, resize back down.
- **Tick marks at compact size:** 60 tick marks on a small ring look cluttered. → Mitigation: hide tick marks below the full tier threshold.
