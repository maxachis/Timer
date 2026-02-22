## MODIFIED Requirements

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
