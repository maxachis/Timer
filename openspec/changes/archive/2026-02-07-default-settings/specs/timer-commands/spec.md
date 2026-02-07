## MODIFIED Requirements

### Requirement: Timer state management
The application SHALL manage a single `CountdownTimer` instance as Tauri app state, initialized with the user's configured default duration at startup. If no settings are saved, the default duration SHALL be 300 seconds (5 minutes).

#### Scenario: Default timer on startup with saved settings
- **WHEN** the application starts and saved settings specify a default duration of 3600 seconds
- **THEN** a `CountdownTimer` with 3600 seconds duration SHALL be available via managed state

#### Scenario: Default timer on startup without saved settings
- **WHEN** the application starts and no saved settings exist
- **THEN** a `CountdownTimer` with 300 seconds duration SHALL be available via managed state

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
