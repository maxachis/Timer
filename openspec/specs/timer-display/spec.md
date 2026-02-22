# timer-display Specification

## Purpose
Frontend countdown UI displaying time, controls, and visual feedback.

## Requirements

### Requirement: Countdown time display
The UI SHALL display the remaining time in MM:SS format, prominently centered in the window. The display SHALL scale fluidly with the window size while remaining readable at all supported sizes.

#### Scenario: Display while running
- **WHEN** the timer is running with 125 seconds remaining
- **THEN** the display SHALL show "02:05"

#### Scenario: Display at zero
- **WHEN** the timer has finished
- **THEN** the display SHALL show "00:00"

#### Scenario: Display at small window size
- **WHEN** the window is resized below 180px in either dimension
- **THEN** the time display SHALL remain visible and readable as the primary UI element

### Requirement: Real-time display updates
The UI SHALL update the time display at approximately 100ms intervals while the timer is running.

#### Scenario: Smooth countdown
- **WHEN** the timer is running
- **THEN** the displayed time SHALL update smoothly without visible lag or jumps

#### Scenario: Polling stops when not running
- **WHEN** the timer is idle, paused, or finished
- **THEN** the UI SHALL NOT poll the backend continuously

### Requirement: Start and pause toggle
The UI SHALL display a single button that toggles between start and pause based on timer state.

#### Scenario: Idle state shows start
- **WHEN** the timer is in idle state
- **THEN** the button SHALL display "Start"

#### Scenario: Running state shows pause
- **WHEN** the timer is running
- **THEN** the button SHALL display "Pause"

#### Scenario: Paused state shows resume
- **WHEN** the timer is paused
- **THEN** the button SHALL display "Resume"

### Requirement: Reset button
The UI SHALL display a reset button that returns the timer to its initial duration.

#### Scenario: Reset available when not idle
- **WHEN** the timer is running, paused, or finished
- **THEN** the reset button SHALL be enabled

#### Scenario: Reset disabled when idle
- **WHEN** the timer is in idle state
- **THEN** the reset button SHALL be disabled

### Requirement: Add time button
The UI SHALL display a "+5:00" button that adds 5 minutes to the timer.

#### Scenario: Add time while running or paused
- **WHEN** the timer is running or paused and the user clicks "+5:00"
- **THEN** 300 seconds SHALL be added to the remaining time

#### Scenario: Add time disabled when finished or idle
- **WHEN** the timer is finished or idle
- **THEN** the "+5:00" button SHALL be disabled

### Requirement: Remove time button
The UI SHALL display a "-5:00" button that removes 5 minutes from the timer.

#### Scenario: Remove time while running or paused
- **WHEN** the timer is running or paused and the user clicks "-5:00"
- **THEN** 300 seconds SHALL be removed from the remaining time

#### Scenario: Remove time disabled when finished or idle
- **WHEN** the timer is finished or idle
- **THEN** the "-5:00" button SHALL be disabled

### Requirement: Finished state visual feedback
The UI SHALL provide clear visual feedback when the timer finishes. At small window sizes, the visual feedback SHALL be conveyed through the time display itself (e.g., color change) since the progress ring may be hidden.

#### Scenario: Timer completes at full size
- **WHEN** the timer reaches zero and the window is at default size (400x600)
- **THEN** the display SHALL visually indicate completion (e.g., color change or message)
- **AND** a "Reset" action SHALL be prominently available

#### Scenario: Timer completes at tiny size
- **WHEN** the timer reaches zero and the window is below 180px in either dimension
- **THEN** the time digits SHALL visually indicate completion via color change
- **AND** the primary button SHALL show the reset action
