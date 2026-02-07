## 1. Module Setup

- [x] 1.1 Create src/timer.rs and declare the module in src/main.rs

## 2. Core Types

- [x] 2.1 Define TimerState enum (Idle, Running, Paused, Finished)
- [x] 2.2 Define CountdownTimer struct with duration, state, and elapsed fields
- [x] 2.3 Define DEFAULT_INCREMENT constant (300 seconds)

## 3. Timer Operations

- [x] 3.1 Implement CountdownTimer::new(duration_secs) constructor
- [x] 3.2 Implement start() — transition from Idle to Running
- [x] 3.3 Implement pause() — transition from Running to Paused
- [x] 3.4 Implement resume() — transition from Paused to Running
- [x] 3.5 Implement remaining() — calculate remaining time from Instant
- [x] 3.6 Implement is_finished() — check if remaining is zero
- [x] 3.7 Implement add_time(Option<Duration>) — add to remaining
- [x] 3.8 Implement remove_time(Option<Duration>) — subtract, clamping to zero
- [x] 3.9 Implement reset() — return to Idle with original duration

## 4. Tests

- [x] 4.1 Test creation and initial state
- [x] 4.2 Test start/pause/resume lifecycle
- [x] 4.3 Test add_time and remove_time (including clamp to zero)
- [x] 4.4 Test reset from running and paused states
- [x] 4.5 Test is_finished detection
