# always-on-top Specification

## Purpose
Always-on-top window pinning, allowing the timer to stay above all other windows.

## Requirements

### Requirement: Always-on-top toggle
The application SHALL allow the user to pin the timer window above all other windows.

#### Scenario: Pin the window
- **WHEN** the user activates the always-on-top toggle
- **THEN** the timer window SHALL remain above all other windows

#### Scenario: Unpin the window
- **WHEN** the user deactivates the always-on-top toggle
- **THEN** the timer window SHALL return to normal window stacking behavior

#### Scenario: Default state on launch
- **WHEN** the application starts
- **THEN** always-on-top SHALL be disabled by default

### Requirement: Pin button in the UI
The application SHALL display a toggle button indicating the current always-on-top state.

#### Scenario: Button shows unpinned state
- **WHEN** always-on-top is disabled
- **THEN** the button SHALL display an outline pin icon

#### Scenario: Button shows pinned state
- **WHEN** always-on-top is enabled
- **THEN** the button SHALL display a filled pin icon to indicate the active state

### Requirement: Always-on-top permissions
The application SHALL declare the required Tauri window permissions for always-on-top control.

#### Scenario: Permissions configured
- **WHEN** the application starts
- **THEN** `core:window:allow-set-always-on-top` and `core:window:allow-is-always-on-top` SHALL be declared in capabilities
