use std::time::{Duration, Instant};

/// Default time increment/decrement: 5 minutes (300 seconds).
pub const DEFAULT_INCREMENT: Duration = Duration::from_secs(300);

/// The lifecycle state of a countdown timer.
#[derive(Debug)]
enum TimerState {
    /// Timer has been created but not started (or has been reset).
    Idle,
    /// Timer is actively counting down.
    Running {
        started_at: Instant,
        previously_elapsed: Duration,
    },
    /// Timer is paused; elapsed time is frozen.
    Paused { elapsed: Duration },
    /// Timer has completed (remaining reached zero).
    Finished,
}

/// A countdown timer that tracks remaining time using monotonic `Instant`.
#[derive(Debug)]
pub struct CountdownTimer {
    /// The original duration this timer was created with (used for reset).
    original_duration: Duration,
    /// The current target duration (may differ from original after add/remove).
    duration: Duration,
    /// Current lifecycle state.
    state: TimerState,
}

impl CountdownTimer {
    /// Create a new countdown timer with the given duration in seconds.
    pub fn new(duration_secs: u64) -> Self {
        Self {
            original_duration: Duration::from_secs(duration_secs),
            duration: Duration::from_secs(duration_secs),
            state: TimerState::Idle,
        }
    }

    /// Start the timer. Only transitions from Idle.
    pub fn start(&mut self) {
        if let TimerState::Idle = self.state {
            self.state = TimerState::Running {
                started_at: Instant::now(),
                previously_elapsed: Duration::ZERO,
            };
        }
    }

    /// Pause a running timer, freezing the elapsed time.
    pub fn pause(&mut self) {
        if let TimerState::Running {
            started_at,
            previously_elapsed,
        } = self.state
        {
            let elapsed = previously_elapsed + started_at.elapsed();
            self.state = TimerState::Paused { elapsed };
        }
    }

    /// Resume a paused timer.
    pub fn resume(&mut self) {
        if let TimerState::Paused { elapsed } = self.state {
            self.state = TimerState::Running {
                started_at: Instant::now(),
                previously_elapsed: elapsed,
            };
        }
    }

    /// Get the remaining time on the timer.
    pub fn remaining(&self) -> Duration {
        let elapsed = match &self.state {
            TimerState::Idle => return self.duration,
            TimerState::Running {
                started_at,
                previously_elapsed,
            } => *previously_elapsed + started_at.elapsed(),
            TimerState::Paused { elapsed } => *elapsed,
            TimerState::Finished => return Duration::ZERO,
        };

        self.duration.saturating_sub(elapsed)
    }

    /// Check if the timer has finished (remaining time is zero).
    pub fn is_finished(&self) -> bool {
        matches!(self.state, TimerState::Finished) || self.remaining().is_zero()
    }

    /// Add time to the timer. Uses `DEFAULT_INCREMENT` (5 minutes) if `None`.
    pub fn add_time(&mut self, amount: Option<Duration>) {
        let increment = amount.unwrap_or(DEFAULT_INCREMENT);
        match &self.state {
            TimerState::Running { .. } | TimerState::Paused { .. } => {
                self.duration += increment;
            }
            _ => {}
        }
    }

    /// Remove time from the timer. Uses `DEFAULT_INCREMENT` (5 minutes) if `None`.
    /// Remaining time will not go below zero — if it does, the timer finishes.
    pub fn remove_time(&mut self, amount: Option<Duration>) {
        let decrement = amount.unwrap_or(DEFAULT_INCREMENT);
        match &self.state {
            TimerState::Running { .. } | TimerState::Paused { .. } => {
                // Check if removing this amount would finish the timer
                let remaining = self.remaining();
                if remaining <= decrement {
                    self.state = TimerState::Finished;
                } else {
                    self.duration = self.duration.saturating_sub(decrement);
                }
            }
            _ => {}
        }
    }

    /// Get the current state name as a string.
    pub fn state_name(&self) -> &'static str {
        if self.is_finished() {
            return "finished";
        }
        match &self.state {
            TimerState::Idle => "idle",
            TimerState::Running { .. } => "running",
            TimerState::Paused { .. } => "paused",
            TimerState::Finished => "finished",
        }
    }

    /// Reset the timer to its original duration in an idle state.
    pub fn reset(&mut self) {
        self.duration = self.original_duration;
        self.state = TimerState::Idle;
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::thread;

    #[test]
    fn test_creation_and_initial_state() {
        let timer = CountdownTimer::new(300);
        assert_eq!(timer.remaining(), Duration::from_secs(300));
        assert!(!timer.is_finished());
    }

    #[test]
    fn test_start_pause_resume_lifecycle() {
        let mut timer = CountdownTimer::new(300);

        // Start
        timer.start();
        thread::sleep(Duration::from_millis(50));
        let r1 = timer.remaining();
        assert!(r1 < Duration::from_secs(300));

        // Pause — remaining should freeze
        timer.pause();
        let paused_remaining = timer.remaining();
        thread::sleep(Duration::from_millis(50));
        assert_eq!(timer.remaining(), paused_remaining);

        // Resume — remaining should continue decreasing
        timer.resume();
        thread::sleep(Duration::from_millis(50));
        assert!(timer.remaining() < paused_remaining);
    }

    #[test]
    fn test_add_time_default() {
        let mut timer = CountdownTimer::new(120);
        timer.start();
        timer.add_time(None); // adds 300s (5 minutes)
        assert!(timer.remaining() > Duration::from_secs(400));
    }

    #[test]
    fn test_add_time_custom_while_paused() {
        let mut timer = CountdownTimer::new(120);
        timer.start();
        timer.pause();
        let before = timer.remaining();
        timer.add_time(Some(Duration::from_secs(60)));
        let after = timer.remaining();
        assert_eq!(after - before, Duration::from_secs(60));
    }

    #[test]
    fn test_remove_time_default() {
        let mut timer = CountdownTimer::new(600);
        timer.start();
        timer.remove_time(None); // removes 300s (5 minutes)
        assert!(timer.remaining() < Duration::from_secs(305));
    }

    #[test]
    fn test_remove_time_clamps_to_zero() {
        let mut timer = CountdownTimer::new(60);
        timer.start();
        timer.remove_time(Some(Duration::from_secs(120)));
        assert!(timer.is_finished());
        assert_eq!(timer.remaining(), Duration::ZERO);
    }

    #[test]
    fn test_reset_from_running() {
        let mut timer = CountdownTimer::new(300);
        timer.start();
        thread::sleep(Duration::from_millis(50));
        timer.reset();
        assert_eq!(timer.remaining(), Duration::from_secs(300));
        assert!(!timer.is_finished());
        // Verify it's idle — starting should work
        timer.start();
        thread::sleep(Duration::from_millis(10));
        assert!(timer.remaining() < Duration::from_secs(300));
    }

    #[test]
    fn test_reset_from_paused() {
        let mut timer = CountdownTimer::new(300);
        timer.start();
        thread::sleep(Duration::from_millis(50));
        timer.pause();
        timer.reset();
        assert_eq!(timer.remaining(), Duration::from_secs(300));
        assert!(!timer.is_finished());
    }

    #[test]
    fn test_is_finished_detection() {
        let mut timer = CountdownTimer::new(0);
        timer.start();
        assert!(timer.is_finished());

        let timer2 = CountdownTimer::new(300);
        assert!(!timer2.is_finished());
    }
}
