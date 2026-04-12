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

/// A named timer for use in a collection.
#[derive(Debug)]
pub struct NamedTimer {
    pub name: String,
    pub timer: CountdownTimer,
}

/// Info about a timer in the collection, suitable for serialization.
#[derive(Debug, Clone)]
pub struct TimerInfo {
    pub index: usize,
    pub name: String,
    pub state: String,
    pub remaining_secs: f64,
    pub is_active: bool,
}

/// Maximum number of timers allowed.
const MAX_TIMERS: usize = 3;

/// A collection of up to 3 named countdown timers with one active at a time.
#[derive(Debug)]
pub struct TimerCollection {
    timers: Vec<NamedTimer>,
    active_index: usize,
}

impl TimerCollection {
    /// Create a new collection with a single unnamed timer.
    pub fn new(default_duration_secs: u64) -> Self {
        Self {
            timers: vec![NamedTimer {
                name: String::new(),
                timer: CountdownTimer::new(default_duration_secs),
            }],
            active_index: 0,
        }
    }

    /// Create a collection from persisted timer names.
    pub fn from_names(names: &[String], default_duration_secs: u64) -> Self {
        if names.is_empty() {
            return Self::new(default_duration_secs);
        }
        let timers = names
            .iter()
            .take(MAX_TIMERS)
            .map(|name| NamedTimer {
                name: name.clone(),
                timer: CountdownTimer::new(default_duration_secs),
            })
            .collect();
        Self {
            timers,
            active_index: 0,
        }
    }

    pub fn active(&self) -> &CountdownTimer {
        &self.timers[self.active_index].timer
    }

    pub fn active_mut(&mut self) -> &mut CountdownTimer {
        &mut self.timers[self.active_index].timer
    }

    pub fn active_name(&self) -> &str {
        &self.timers[self.active_index].name
    }

    pub fn active_index(&self) -> usize {
        self.active_index
    }

    pub fn count(&self) -> usize {
        self.timers.len()
    }

    /// Add a new timer. Returns the new timer's index, or an error if at max capacity.
    pub fn add_timer(&mut self, name: String, duration_secs: u64) -> Result<usize, &'static str> {
        if self.timers.len() >= MAX_TIMERS {
            return Err("Maximum of 3 timers reached. Remove a timer before adding another.");
        }
        self.timers.push(NamedTimer {
            name,
            timer: CountdownTimer::new(duration_secs),
        });
        Ok(self.timers.len() - 1)
    }

    /// Remove a timer by index. Cannot remove the last timer.
    pub fn remove_timer(&mut self, index: usize) -> Result<(), &'static str> {
        if self.timers.len() <= 1 {
            return Err("Cannot remove the last timer");
        }
        if index >= self.timers.len() {
            return Err("Invalid timer index");
        }
        self.timers.remove(index);
        // Adjust active_index after removal
        if self.active_index >= self.timers.len() {
            self.active_index = self.timers.len() - 1;
        } else if index < self.active_index {
            self.active_index -= 1;
        }
        Ok(())
    }

    /// Switch to a different timer. Pauses the current timer first.
    pub fn switch_to(&mut self, index: usize) -> Result<(), &'static str> {
        if index >= self.timers.len() {
            return Err("Invalid timer index");
        }
        if index == self.active_index {
            return Ok(());
        }
        // Pause the current timer
        self.timers[self.active_index].timer.pause();
        self.active_index = index;
        // Pause the target too (safety: shouldn't be running, but enforce invariant)
        self.timers[self.active_index].timer.pause();
        Ok(())
    }

    /// Switch to the next timer (wrapping).
    pub fn switch_next(&mut self) {
        if self.timers.len() <= 1 {
            return;
        }
        let next = (self.active_index + 1) % self.timers.len();
        let _ = self.switch_to(next);
    }

    /// Switch to the previous timer (wrapping).
    pub fn switch_prev(&mut self) {
        if self.timers.len() <= 1 {
            return;
        }
        let prev = if self.active_index == 0 {
            self.timers.len() - 1
        } else {
            self.active_index - 1
        };
        let _ = self.switch_to(prev);
    }

    /// Rename a timer.
    pub fn rename_timer(&mut self, index: usize, name: String) -> Result<(), &'static str> {
        if index >= self.timers.len() {
            return Err("Invalid timer index");
        }
        self.timers[index].name = name;
        Ok(())
    }

    /// Get info about all timers.
    pub fn timer_list(&self) -> Vec<TimerInfo> {
        self.timers
            .iter()
            .enumerate()
            .map(|(i, t)| TimerInfo {
                index: i,
                name: t.name.clone(),
                state: t.timer.state_name().to_string(),
                remaining_secs: t.timer.remaining().as_secs_f64(),
                is_active: i == self.active_index,
            })
            .collect()
    }

    /// Get just the timer names (for persistence).
    pub fn timer_names(&self) -> Vec<String> {
        self.timers.iter().map(|t| t.name.clone()).collect()
    }
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

    // --- TimerCollection tests ---

    #[test]
    fn test_collection_new_has_one_timer() {
        let col = TimerCollection::new(300);
        assert_eq!(col.count(), 1);
        assert_eq!(col.active_index(), 0);
        assert_eq!(col.active_name(), "");
        assert_eq!(col.active().remaining(), Duration::from_secs(300));
    }

    #[test]
    fn test_collection_from_names() {
        let names = vec!["Work".into(), "Break".into()];
        let col = TimerCollection::from_names(&names, 300);
        assert_eq!(col.count(), 2);
        assert_eq!(col.active_name(), "Work");
    }

    #[test]
    fn test_collection_from_empty_names() {
        let col = TimerCollection::from_names(&[], 300);
        assert_eq!(col.count(), 1);
    }

    #[test]
    fn test_add_timer() {
        let mut col = TimerCollection::new(300);
        let idx = col.add_timer("Second".into(), 600).unwrap();
        assert_eq!(idx, 1);
        assert_eq!(col.count(), 2);
    }

    #[test]
    fn test_add_timer_max_limit() {
        let mut col = TimerCollection::new(300);
        col.add_timer("Two".into(), 300).unwrap();
        col.add_timer("Three".into(), 300).unwrap();
        let result = col.add_timer("Four".into(), 300);
        assert!(result.is_err());
        assert_eq!(col.count(), 3);
    }

    #[test]
    fn test_remove_timer() {
        let mut col = TimerCollection::new(300);
        col.add_timer("Second".into(), 600).unwrap();
        col.remove_timer(1).unwrap();
        assert_eq!(col.count(), 1);
    }

    #[test]
    fn test_remove_last_timer_fails() {
        let mut col = TimerCollection::new(300);
        let result = col.remove_timer(0);
        assert!(result.is_err());
    }

    #[test]
    fn test_remove_active_adjusts_index() {
        let mut col = TimerCollection::new(300);
        col.add_timer("Second".into(), 300).unwrap();
        col.add_timer("Third".into(), 300).unwrap();
        // Switch to last timer
        col.switch_to(2).unwrap();
        assert_eq!(col.active_index(), 2);
        // Remove it — active should fall back
        col.remove_timer(2).unwrap();
        assert_eq!(col.active_index(), 1);
    }

    #[test]
    fn test_remove_before_active_adjusts_index() {
        let mut col = TimerCollection::new(300);
        col.add_timer("Second".into(), 300).unwrap();
        col.add_timer("Third".into(), 300).unwrap();
        col.switch_to(2).unwrap();
        // Remove timer 0 (before active)
        col.remove_timer(0).unwrap();
        assert_eq!(col.active_index(), 1); // was 2, shifted down by 1
    }

    #[test]
    fn test_switch_pauses_current() {
        let mut col = TimerCollection::new(300);
        col.add_timer("Second".into(), 300).unwrap();
        col.active_mut().start();
        assert_eq!(col.active().state_name(), "running");
        col.switch_to(1).unwrap();
        // Previous timer should now be paused
        assert_eq!(col.timers[0].timer.state_name(), "paused");
        assert_eq!(col.active_index(), 1);
    }

    #[test]
    fn test_switch_to_same_is_noop() {
        let mut col = TimerCollection::new(300);
        col.active_mut().start();
        col.switch_to(0).unwrap();
        // Should still be running (not paused)
        assert_eq!(col.active().state_name(), "running");
    }

    #[test]
    fn test_switch_next_wraps() {
        let mut col = TimerCollection::new(300);
        col.add_timer("Two".into(), 300).unwrap();
        col.add_timer("Three".into(), 300).unwrap();
        assert_eq!(col.active_index(), 0);
        col.switch_next();
        assert_eq!(col.active_index(), 1);
        col.switch_next();
        assert_eq!(col.active_index(), 2);
        col.switch_next();
        assert_eq!(col.active_index(), 0); // wrapped
    }

    #[test]
    fn test_switch_prev_wraps() {
        let mut col = TimerCollection::new(300);
        col.add_timer("Two".into(), 300).unwrap();
        assert_eq!(col.active_index(), 0);
        col.switch_prev();
        assert_eq!(col.active_index(), 1); // wrapped to end
        col.switch_prev();
        assert_eq!(col.active_index(), 0);
    }

    #[test]
    fn test_switch_single_timer_noop() {
        let mut col = TimerCollection::new(300);
        col.switch_next();
        assert_eq!(col.active_index(), 0);
        col.switch_prev();
        assert_eq!(col.active_index(), 0);
    }

    #[test]
    fn test_rename_timer() {
        let mut col = TimerCollection::new(300);
        col.rename_timer(0, "Focus".into()).unwrap();
        assert_eq!(col.active_name(), "Focus");
    }

    #[test]
    fn test_timer_list() {
        let mut col = TimerCollection::new(300);
        col.add_timer("Break".into(), 120).unwrap();
        let list = col.timer_list();
        assert_eq!(list.len(), 2);
        assert!(list[0].is_active);
        assert!(!list[1].is_active);
        assert_eq!(list[1].name, "Break");
    }

    #[test]
    fn test_timer_names() {
        let mut col = TimerCollection::new(300);
        col.rename_timer(0, "Work".into()).unwrap();
        col.add_timer("Break".into(), 120).unwrap();
        assert_eq!(col.timer_names(), vec!["Work", "Break"]);
    }
}
