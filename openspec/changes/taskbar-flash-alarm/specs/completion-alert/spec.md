## ADDED Requirements

### Requirement: Taskbar icon flash on timer completion
The application SHALL request OS-level window attention (taskbar flash) when the countdown timer reaches zero, persisting until the alarm is resolved.

#### Scenario: Taskbar flashes when timer finishes
- **WHEN** the timer transitions to the "finished" state
- **THEN** the application SHALL call `requestUserAttention` with the "critical" attention type
- **AND** the taskbar icon SHALL flash persistently on the user's OS

#### Scenario: Taskbar flash persists until reset
- **WHEN** the timer is in the "finished" state and the user has not yet reset
- **THEN** the taskbar icon SHALL continue flashing

#### Scenario: Taskbar flash stops on reset
- **WHEN** the user resets the timer from the "finished" state
- **THEN** the application SHALL cancel the attention request by calling `requestUserAttention` with `null`
- **AND** the taskbar icon SHALL stop flashing

#### Scenario: Flash fires only once per completion
- **WHEN** the timer reaches zero
- **THEN** `requestUserAttention` SHALL be called exactly once during the finish transition
- **AND** subsequent polling cycles SHALL NOT trigger additional attention requests

### Requirement: Window attention permission
The application SHALL declare the `core:window:allow-request-user-attention` permission in the Tauri capability configuration.

#### Scenario: Permission declared in capabilities
- **WHEN** `src-tauri/capabilities/default.json` is inspected
- **THEN** the permissions array SHALL include `core:window:allow-request-user-attention`
