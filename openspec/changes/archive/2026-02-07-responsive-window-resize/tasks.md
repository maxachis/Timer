## 1. Tauri Window Config

- [x] 1.1 Add `minWidth: 100` and `minHeight: 100` to the window definition in `src-tauri/tauri.conf.json`

## 2. Fluid Sizing

- [x] 2.1 Replace the fixed `width: 300px; height: 300px` on `.timer-face` with `clamp(80px, 70vmin, 300px)` for both dimensions
- [x] 2.2 Replace fixed font sizes on `.digits` and `.separator` with `clamp()`-based values (`clamp(1.2rem, 12vmin, 4.5rem)` for digits, proportional for separator)
- [x] 2.3 Replace fixed `width: 64px; height: 64px` on `.primary-btn` with `clamp(32px, 14vmin, 64px)`
- [x] 2.4 Scale `.adjust-btn` min-width and padding with `clamp()` to shrink at smaller sizes

## 3. Progressive Visibility

- [x] 3.1 Add a `@media` block for the tiny tier (`max-width: 179px` or `max-height: 179px`) that hides `.progress-ring`, `.adjust-row .adjust-btn`, `.reset-btn`, `.pin-btn`, `.settings-btn`, and tick marks
- [x] 3.2 Add a `@media` block for the compact tier (between 180px and 349px width / 180px–399px height) that hides tick marks, `.adjust-btn`, `.reset-btn`, `.pin-btn`, and `.settings-btn` but keeps the progress ring visible
- [x] 3.3 Ensure the full tier (≥350px wide and ≥400px tall) shows all elements — verify this is the default behavior with no media query needed

## 4. Layout Adjustments

- [x] 4.1 Update `main` padding to use `clamp()` so it shrinks at small viewports (e.g. `clamp(0.5rem, 3vmin, 2rem)` vertical, `clamp(0.25rem, 2vmin, 1.5rem)` horizontal)
- [x] 4.2 Update `.timer-face` margin-bottom to scale down at small sizes (e.g. `clamp(0.5rem, 4vmin, 2.5rem)`)
- [x] 4.3 Update `.controls` gap to scale with viewport

## 5. Verification

- [x] 5.1 Test at 100x100 — confirm only time digits and primary button visible, text readable (manual)
- [x] 5.2 Test at 250x300 — confirm ring visible without ticks, no adjust/reset/toolbar buttons (manual)
- [x] 5.3 Test at 400x600 — confirm pixel-identical to current UI (manual)
- [x] 5.4 Test at a large size (800x1000) — confirm UI doesn't grow beyond max bounds (manual)
- [x] 5.5 Test settings view remains usable at ≥250px width (manual)
