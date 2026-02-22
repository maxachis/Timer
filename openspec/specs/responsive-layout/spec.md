# responsive-layout Specification

## Purpose
Fluid scaling rules and breakpoint-based visibility for the timer UI, defining which elements appear at which window sizes and how dimensions adapt.

## Requirements

### Requirement: Fluid timer face scaling
The timer face container SHALL scale fluidly based on the viewport's smaller dimension using `vmin` units, bounded by `clamp()` to maintain readability at all sizes.

#### Scenario: Timer face at default window size
- **WHEN** the window is 400x600 (default)
- **THEN** the timer face SHALL render at 300x300 pixels, identical to the current layout

#### Scenario: Timer face at minimum window size
- **WHEN** the window is 100x100
- **THEN** the timer face container SHALL shrink to approximately 80px, showing only the time digits without the progress ring

#### Scenario: Timer face at intermediate size
- **WHEN** the window is 250x300
- **THEN** the timer face SHALL scale proportionally between 80px and 300px based on `vmin`

### Requirement: Fluid text scaling
The time display digits and separator SHALL scale fluidly with the viewport using `clamp()` to remain readable at all window sizes.

#### Scenario: Digits at default size
- **WHEN** the window is 400x600
- **THEN** the digit font size SHALL be 4.5rem

#### Scenario: Digits at minimum size
- **WHEN** the window is 100x100
- **THEN** the digit font size SHALL be no smaller than 1.2rem

### Requirement: Fluid control scaling
The primary action button SHALL scale fluidly between a minimum and maximum size based on viewport dimensions.

#### Scenario: Primary button at default size
- **WHEN** the window is 400x600
- **THEN** the primary button SHALL be 64x64 pixels

#### Scenario: Primary button at small size
- **WHEN** the window is 150x150
- **THEN** the primary button SHALL scale down to approximately 32px while remaining tappable

### Requirement: Progressive visibility tiers
The UI SHALL progressively show or hide elements based on three viewport size tiers: tiny, compact, and full.

#### Scenario: Tiny tier
- **WHEN** the viewport width is below 180px or the viewport height is below 180px
- **THEN** the UI SHALL display only the time digits and the primary action button
- **AND** the progress ring, tick marks, adjust buttons, reset button, pin button, and settings button SHALL be hidden

#### Scenario: Compact tier
- **WHEN** the viewport width is between 180px and 349px and the viewport height is between 180px and 399px
- **THEN** the UI SHALL display the progress ring (without tick marks), the time digits, and the primary action button
- **AND** the adjust buttons, reset button, pin button, and settings button SHALL be hidden

#### Scenario: Full tier
- **WHEN** the viewport width is 350px or greater and the viewport height is 400px or greater
- **THEN** all UI elements SHALL be visible, matching the current default layout

### Requirement: CSS-only implementation
All responsive scaling and visibility changes SHALL be implemented using CSS features only (media queries, `clamp()`, `vmin` units). No JavaScript resize listeners or observers SHALL be used.

#### Scenario: No JS resize handling
- **WHEN** the window is resized
- **THEN** the layout SHALL adapt without any JavaScript event handlers firing for resize events
