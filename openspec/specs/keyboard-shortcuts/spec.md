# keyboard-shortcuts Specification

## Purpose
Keyboard shortcut bindings for timer control actions, enabling keyboard-driven operation.

## Requirements

### Requirement: Space bar toggles primary timer action
The application SHALL use the Space bar to trigger the context-appropriate primary action based on the current timer state.

#### Scenario: Space starts an idle timer
- **WHEN** the timer is idle and the user presses Space
- **THEN** the timer SHALL start

#### Scenario: Space pauses a running timer
- **WHEN** the timer is running and the user presses Space
- **THEN** the timer SHALL pause

#### Scenario: Space resumes a paused timer
- **WHEN** the timer is paused and the user presses Space
- **THEN** the timer SHALL resume

#### Scenario: Space resets a finished timer
- **WHEN** the timer is finished and the user presses Space
- **THEN** the timer SHALL reset to idle

### Requirement: Reset shortcut
The application SHALL allow resetting the timer via keyboard when the timer is not idle.

#### Scenario: R key resets the timer
- **WHEN** the timer is running, paused, or finished and the user presses R
- **THEN** the timer SHALL reset to idle

#### Scenario: Down arrow resets the timer
- **WHEN** the timer is running, paused, or finished and the user presses the down arrow key
- **THEN** the timer SHALL reset to idle

#### Scenario: Reset ignored when idle
- **WHEN** the timer is idle and the user presses R or the down arrow key
- **THEN** no action SHALL be taken

### Requirement: Add time shortcut
The application SHALL allow adding 5 minutes via keyboard when the timer is running or paused.

#### Scenario: Plus key adds time
- **WHEN** the timer is running or paused and the user presses +
- **THEN** 5 minutes SHALL be added to the timer

#### Scenario: Equals key adds time
- **WHEN** the timer is running or paused and the user presses =
- **THEN** 5 minutes SHALL be added to the timer

#### Scenario: Right arrow adds time
- **WHEN** the timer is running or paused and the user presses the right arrow key
- **THEN** 5 minutes SHALL be added to the timer

#### Scenario: Add time ignored when idle or finished
- **WHEN** the timer is idle or finished and the user presses +, =, or right arrow
- **THEN** no action SHALL be taken

### Requirement: Remove time shortcut
The application SHALL allow removing 5 minutes via keyboard when the timer is running or paused.

#### Scenario: Minus key removes time
- **WHEN** the timer is running or paused and the user presses -
- **THEN** 5 minutes SHALL be removed from the timer

#### Scenario: Left arrow removes time
- **WHEN** the timer is running or paused and the user presses the left arrow key
- **THEN** 5 minutes SHALL be removed from the timer

#### Scenario: Remove time ignored when idle or finished
- **WHEN** the timer is idle or finished and the user presses - or left arrow
- **THEN** no action SHALL be taken

### Requirement: Escape resets a finished timer
The application SHALL allow resetting a finished timer using the Escape key.

#### Scenario: Escape resets from finished state
- **WHEN** the timer is finished and the user presses Escape
- **THEN** the timer SHALL reset to idle

#### Scenario: Escape ignored when not finished
- **WHEN** the timer is not finished and the user presses Escape
- **THEN** no action SHALL be taken

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
