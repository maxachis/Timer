## ADDED Requirements

### Requirement: Window body acts as a drag region
The application's main container SHALL be marked as a Tauri drag region so that clicking and dragging on non-interactive areas moves the window.

#### Scenario: Dragging on empty space moves the window
- **WHEN** the user clicks and drags on the background area of the timer UI (not on a button or input)
- **THEN** the window SHALL move with the cursor, behaving as a native window drag

#### Scenario: Buttons remain clickable
- **WHEN** the user clicks on any button (start, pause, reset, add time, remove time, settings, pin)
- **THEN** the button action SHALL fire normally
- **AND** no window drag SHALL be initiated

#### Scenario: Settings inputs remain interactive
- **WHEN** the settings panel is open and the user clicks on an input field
- **THEN** the input SHALL receive focus and accept input
- **AND** no window drag SHALL be initiated

### Requirement: Drag region permission
The Tauri capability configuration SHALL include the `core:window:allow-start-dragging` permission for the main window.

#### Scenario: Permission declared
- **WHEN** `src-tauri/capabilities/default.json` is inspected
- **THEN** the `permissions` array SHALL include `"core:window:allow-start-dragging"`
