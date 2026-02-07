## MODIFIED Requirements

### Requirement: Add time to timer
The system SHALL allow adding time to a running or paused timer. The increment SHALL be configurable. A custom increment MAY be specified per call.

#### Scenario: Add default increment while running
- **WHEN** time is added to a running timer with 120 seconds remaining using the default increment
- **THEN** the remaining time SHALL increase by the configured default increment

#### Scenario: Add custom increment while paused
- **WHEN** 60 seconds is added to a paused timer with 120 seconds remaining
- **THEN** the remaining time SHALL become 180 seconds

### Requirement: Remove time from timer
The system SHALL allow removing time from a running or paused timer. The decrement SHALL be configurable. A custom decrement MAY be specified per call. Remaining time SHALL NOT go below zero.

#### Scenario: Remove default increment with sufficient remaining
- **WHEN** time is removed from a timer with 600 seconds remaining using the configured default decrement
- **THEN** the remaining time SHALL decrease by the configured default decrement

#### Scenario: Remove time that would go negative
- **WHEN** 120 seconds is removed from a timer with 60 seconds remaining
- **THEN** the remaining time SHALL become 0 seconds
- **AND** the timer SHALL be considered finished
