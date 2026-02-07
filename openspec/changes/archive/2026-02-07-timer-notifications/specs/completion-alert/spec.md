## ADDED Requirements

### Requirement: System notification on timer completion
The application SHALL send an OS-level toast notification when the countdown timer reaches zero.

#### Scenario: Timer finishes while app is focused
- **WHEN** the timer reaches zero
- **THEN** a system notification SHALL be displayed with the title "Timer Complete"
- **AND** the notification body SHALL indicate the timer has finished

#### Scenario: Timer finishes while app is minimized
- **WHEN** the timer reaches zero and the application window is not focused
- **THEN** a system notification SHALL appear in the OS notification area

#### Scenario: Notification fires only once per completion
- **WHEN** the timer reaches zero
- **THEN** the notification SHALL be sent exactly once
- **AND** subsequent polling cycles SHALL NOT trigger additional notifications

### Requirement: Notification permission handling
The application SHALL request notification permission from the OS before sending the first notification.

#### Scenario: Permission not yet granted
- **WHEN** the timer finishes for the first time and notification permission has not been granted
- **THEN** the application SHALL request permission from the OS
- **AND** if granted, the notification SHALL be sent immediately

#### Scenario: Permission already granted
- **WHEN** the timer finishes and notification permission is already granted
- **THEN** the notification SHALL be sent without prompting

#### Scenario: Permission denied
- **WHEN** notification permission is denied by the OS
- **THEN** the application SHALL NOT attempt to send the notification
- **AND** the audio alert SHALL still play

### Requirement: Audio alert on timer completion
The application SHALL play an audible alert sound when the countdown timer reaches zero.

#### Scenario: Alert sound plays on completion
- **WHEN** the timer reaches zero
- **THEN** a short alert tone SHALL play

#### Scenario: Sound plays only once per completion
- **WHEN** the timer reaches zero
- **THEN** the alert sound SHALL play exactly once
- **AND** subsequent polling cycles SHALL NOT replay the sound

### Requirement: Notification plugin registration
The application SHALL register `tauri-plugin-notification` and include the required permissions in the capability configuration.

#### Scenario: Plugin and permissions configured
- **WHEN** the application starts
- **THEN** the notification plugin SHALL be initialized
- **AND** `notification:default` permission SHALL be declared in `src-tauri/capabilities/default.json`
