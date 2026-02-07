## ADDED Requirements

### Requirement: Timer creation with duration
The system SHALL allow creating a CountdownTimer with a specified duration in seconds. The timer SHALL store the original duration for reset purposes.

#### Scenario: Create a timer
- **WHEN** a CountdownTimer is created with a duration of 300 seconds
- **THEN** the timer's remaining time SHALL be 300 seconds
- **AND** the timer SHALL NOT be running

### Requirement: Timer start
The system SHALL allow starting an idle timer, which begins counting down.

#### Scenario: Start an idle timer
- **WHEN** the timer is started
- **THEN** the timer SHALL be in the running state
- **AND** the remaining time SHALL decrease over time

### Requirement: Timer pause
The system SHALL allow pausing a running timer.

#### Scenario: Pause a running timer
- **WHEN** a running timer is paused
- **THEN** the timer SHALL stop counting down
- **AND** the remaining time SHALL be preserved at the moment of pause

### Requirement: Timer resume
The system SHALL allow resuming a paused timer.

#### Scenario: Resume a paused timer
- **WHEN** a paused timer is resumed
- **THEN** the timer SHALL continue counting down from where it was paused

### Requirement: Add time to timer
The system SHALL allow adding time to a running or paused timer. The default increment SHALL be 300 seconds (5 minutes). A custom increment MAY be specified.

#### Scenario: Add default increment while running
- **WHEN** time is added to a running timer with 120 seconds remaining using the default increment
- **THEN** the remaining time SHALL become approximately 420 seconds

#### Scenario: Add custom increment while paused
- **WHEN** 60 seconds is added to a paused timer with 120 seconds remaining
- **THEN** the remaining time SHALL become 180 seconds

### Requirement: Remove time from timer
The system SHALL allow removing time from a running or paused timer. The default decrement SHALL be 300 seconds (5 minutes). A custom decrement MAY be specified. Remaining time SHALL NOT go below zero.

#### Scenario: Remove default increment with sufficient remaining
- **WHEN** time is removed from a timer with 600 seconds remaining using the default decrement
- **THEN** the remaining time SHALL become approximately 300 seconds

#### Scenario: Remove time that would go negative
- **WHEN** 120 seconds is removed from a timer with 60 seconds remaining
- **THEN** the remaining time SHALL become 0 seconds
- **AND** the timer SHALL be considered finished

### Requirement: Timer reset
The system SHALL allow resetting the timer back to its original duration. A reset timer SHALL be in the idle (not running) state.

#### Scenario: Reset a running timer
- **WHEN** a running timer with original duration 300 seconds is reset
- **THEN** the remaining time SHALL be 300 seconds
- **AND** the timer SHALL NOT be running

#### Scenario: Reset a paused timer
- **WHEN** a paused timer with original duration 300 seconds is reset
- **THEN** the remaining time SHALL be 300 seconds
- **AND** the timer SHALL NOT be running

### Requirement: Timer completion detection
The system SHALL report when the countdown has reached zero.

#### Scenario: Timer finishes
- **WHEN** the remaining time reaches zero
- **THEN** is_finished SHALL return true

#### Scenario: Timer still running
- **WHEN** the remaining time is greater than zero
- **THEN** is_finished SHALL return false
