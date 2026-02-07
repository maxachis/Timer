## MODIFIED Requirements

### Requirement: Default browser behavior prevention
The application SHALL prevent default browser behavior for all handled keyboard shortcuts.

#### Scenario: Space does not scroll the page
- **WHEN** the user presses Space to control the timer
- **THEN** the page SHALL NOT scroll

#### Scenario: Arrow keys do not move focus
- **WHEN** the user presses arrow keys to adjust time or reset
- **THEN** the default browser focus navigation SHALL NOT occur

#### Scenario: T key does not type in inputs
- **WHEN** the user presses T to toggle always-on-top
- **THEN** default browser behavior SHALL be prevented

## ADDED Requirements

### Requirement: Always-on-top shortcut
The application SHALL allow toggling always-on-top via the T key.

#### Scenario: T key toggles always-on-top on
- **WHEN** always-on-top is disabled and the user presses T
- **THEN** always-on-top SHALL be enabled

#### Scenario: T key toggles always-on-top off
- **WHEN** always-on-top is enabled and the user presses T
- **THEN** always-on-top SHALL be disabled

#### Scenario: T key works in any timer state
- **WHEN** the user presses T regardless of whether the timer is idle, running, paused, or finished
- **THEN** the always-on-top state SHALL toggle
