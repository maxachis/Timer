## ADDED Requirements

### Requirement: Settings navigation via cog icon
The application SHALL display a cog icon button for navigating to the settings view.

#### Scenario: Cog button position and visibility
- **WHEN** the timer view is displayed
- **THEN** a cog icon button SHALL be visible in the top-left corner of the window

#### Scenario: Cog button opens settings
- **WHEN** the user clicks the cog icon button
- **THEN** the timer view SHALL be replaced by the settings view

### Requirement: Settings view layout
The application SHALL display a settings form with fields for configuring default start time and default increment.

#### Scenario: Default start time field
- **WHEN** the settings view is displayed
- **THEN** a numeric input labeled "Default start time" SHALL be shown with the value in minutes

#### Scenario: Default increment field
- **WHEN** the settings view is displayed
- **THEN** a numeric input labeled "Time increment" SHALL be shown with the value in minutes

#### Scenario: Fields show current values
- **WHEN** the settings view is opened
- **THEN** the fields SHALL be populated with the currently saved settings

### Requirement: Settings input validation
The application SHALL enforce valid ranges for settings inputs.

#### Scenario: Default start time range
- **WHEN** the user enters a default start time
- **THEN** the value SHALL be constrained to 1–180 minutes

#### Scenario: Default increment range
- **WHEN** the user enters a default increment
- **THEN** the value SHALL be constrained to 1–60 minutes

### Requirement: Save and return from settings
The application SHALL allow the user to save settings and return to the timer view.

#### Scenario: Save button persists settings
- **WHEN** the user clicks the save button
- **THEN** the settings SHALL be persisted and the view SHALL return to the timer

#### Scenario: Back button without saving
- **WHEN** the user clicks the back button without saving
- **THEN** the view SHALL return to the timer without changing settings

#### Scenario: Timer resets after settings saved
- **WHEN** settings are saved and the timer is idle
- **THEN** the timer SHALL reinitialize with the new default start time

### Requirement: Keyboard shortcuts disabled in settings
The application SHALL not process timer keyboard shortcuts while the settings view is active.

#### Scenario: Space key in settings view
- **WHEN** the settings view is active and the user presses Space
- **THEN** no timer action SHALL be taken

#### Scenario: Other shortcut keys in settings view
- **WHEN** the settings view is active and the user presses R, T, arrow keys, or Escape
- **THEN** no timer action SHALL be taken

### Requirement: Dynamic adjust button labels
The application SHALL display the configured increment on the time adjust buttons.

#### Scenario: Adjust buttons reflect saved increment
- **WHEN** the timer view is displayed and the default increment is set to 10 minutes
- **THEN** the adjust buttons SHALL display "10 min" instead of "5 min"
