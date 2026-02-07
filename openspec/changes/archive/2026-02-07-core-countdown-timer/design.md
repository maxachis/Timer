## Context

The Timer project is a greenfield Rust application (edition 2024). Currently
`src/main.rs` contains only the default hello-world. We need a timer module
that can be used by any future frontend (Tauri, CLI, etc.) without coupling
to a specific UI framework.

## Goals / Non-Goals

**Goals:**
- Accurate countdown tracking using monotonic clock (`std::time::Instant`)
- Clean state machine: Idle → Running → Paused → Finished
- Zero external dependencies — pure `std` library
- Testable without sleeps via a time-source abstraction

**Non-Goals:**
- UI or display rendering
- Notification/alert system
- Persistence or serialization
- Thread safety (can be added later with Arc<Mutex<>>)

## Decisions

### Decision 1: Use `Instant`-based elapsed tracking

**Approach:** Track `started_at: Instant` when running, accumulate
`elapsed_while_paused: Duration` on pause. Remaining = duration - total_elapsed.

**Why not tick-based?** Tick/sleep approaches drift over time and couple the
timer to a specific update frequency. Instant-based calculation is always
accurate regardless of when you query it.

### Decision 2: State enum for timer lifecycle

**Approach:** Use an enum `TimerState { Idle, Running { started_at: Instant },
Paused { elapsed: Duration }, Finished }` to make invalid states unrepresentable.

**Why not booleans?** With `is_running` + `is_paused` booleans, you can have
impossible combinations (running AND paused). An enum prevents this at
compile time.

### Decision 3: Default increment as a constant

**Approach:** Define `DEFAULT_INCREMENT: Duration = Duration::from_secs(300)`
as a module-level constant. Methods accept `Option<Duration>` to allow
custom increments while defaulting to 5 minutes.

## Risks / Trade-offs

- [No thread safety] → Acceptable for now; the timer will be owned by a
  single thread. Can wrap in Arc<Mutex<>> when Tauri integration happens.
- [Instant is not mockable in tests] → We'll use short durations and
  `std::thread::sleep` for brief periods in tests. For a pure logic module
  this is acceptable.
