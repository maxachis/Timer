## Why

The Timer app needs core countdown logic before any UI or notification features can be built. A well-designed timer module provides the foundation for all future features — display, controls, alerts.

## What Changes

- Add a `CountdownTimer` struct with accurate time tracking using `std::time::Instant`
- Support start, pause, resume, and completion detection
- Support adding/removing time in configurable increments while the timer is running or paused
- Add unit tests to verify all timer states and transitions

## Capabilities

### New Capabilities
- `countdown-timer`: A Rust module providing countdown timer logic with start/pause/resume, time adjustment, and completion detection

### Modified Capabilities
<!-- None — this is a greenfield addition -->

## Impact

- `src/timer.rs`: New module containing the `CountdownTimer` struct and implementation
- `src/main.rs`: Updated to declare the timer module
