# timer-commands Specification

## Purpose
Tauri IPC commands that bridge the frontend UI to the backend CountdownTimer state.

## Requirements

### Requirement: Timer state management
The application SHALL manage a single `CountdownTimer` instance as Tauri app state, initialized with the user's configured default duration at startup. If no settings are saved, the default duration SHALL be 300 seconds (5 minutes).

#### Scenario: Default timer on startup with saved settings
- **WHEN** the application starts and saved settings specify a default duration of 3600 seconds
- **THEN** a `CountdownTimer` with 3600 seconds duration SHALL be available via managed state

#### Scenario: Default timer on startup without saved settings
- **WHEN** the application starts and no saved settings exist
- **THEN** a `CountdownTimer` with 300 seconds duration SHALL be available via managed state

### Requirement: Create timer command
The application SHALL expose a `create_timer` Tauri command that reinitializes the timer with a specified duration.

#### Scenario: Create a new timer
- **WHEN** the frontend invokes `create_timer` with `duration_secs: 600`
- **THEN** the managed timer SHALL be replaced with a new timer of 600 seconds
- **AND** the timer SHALL be in idle state

### Requirement: Start timer command
The application SHALL expose a `start_timer` Tauri command.

#### Scenario: Start the timer
- **WHEN** the frontend invokes `start_timer`
- **THEN** the managed timer SHALL transition to running state

### Requirement: Pause timer command
The application SHALL expose a `pause_timer` Tauri command.

#### Scenario: Pause the timer
- **WHEN** the frontend invokes `pause_timer`
- **THEN** the managed timer SHALL transition to paused state

### Requirement: Resume timer command
The application SHALL expose a `resume_timer` Tauri command.

#### Scenario: Resume the timer
- **WHEN** the frontend invokes `resume_timer`
- **THEN** the managed timer SHALL transition from paused back to running state

### Requirement: Reset timer command
The application SHALL expose a `reset_timer` Tauri command.

#### Scenario: Reset the timer
- **WHEN** the frontend invokes `reset_timer`
- **THEN** the managed timer SHALL return to idle state with its original duration

### Requirement: Add time command
The application SHALL expose an `add_time` Tauri command that adds the configured default increment to the timer.

#### Scenario: Add time uses configured increment
- **WHEN** the frontend invokes `add_time` and the configured increment is 600 seconds
- **THEN** 600 seconds SHALL be added to the timer's remaining time

### Requirement: Remove time command
The application SHALL expose a `remove_time` Tauri command that removes the configured default decrement from the timer.

#### Scenario: Remove time uses configured increment
- **WHEN** the frontend invokes `remove_time` and the configured increment is 600 seconds
- **THEN** 600 seconds SHALL be removed from the timer's remaining time
- **AND** if remaining time would go below zero, the timer SHALL finish

### Requirement: Get timer status command
The application SHALL expose a `get_timer_status` Tauri command that returns the timer's current state as a serializable struct.

#### Scenario: Get status of running timer
- **WHEN** the frontend invokes `get_timer_status`
- **THEN** the response SHALL include `remaining_secs` as a floating-point number
- **AND** the response SHALL include `state` as a string ("idle", "running", "paused", or "finished")
- **AND** the response SHALL include `is_finished` as a boolean
