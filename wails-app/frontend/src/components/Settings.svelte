<script lang="ts">
  import { timer } from "../lib/timerStore.svelte";
</script>

<!-- Settings view -->
<div class="texture"></div>

<button class="settings-back-btn" onclick={() => timer.closeSettings()} aria-label="Back to timer">
  <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
    <path d="M10 2L4 8L10 14" />
  </svg>
</button>

<button class="window-btn minimize-btn" onclick={() => timer.handleMinimize()} aria-label="Minimize">
  <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
    <line x1="2" y1="6" x2="10" y2="6" />
  </svg>
</button>

<button class="window-btn close-btn" onclick={() => timer.handleClose()} aria-label="Close">
  <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
    <line x1="2" y1="2" x2="10" y2="10" />
    <line x1="10" y1="2" x2="2" y2="10" />
  </svg>
</button>

<div class="settings-panel">
  <h2 class="settings-title">Settings</h2>

  <div class="settings-field">
    <label class="settings-label" for="duration">Default start time</label>
    <div class="settings-input-row">
      <input
        id="duration"
        type="number"
        min="1"
        max="180"
        bind:value={timer.settingsDurationMin}
        class="settings-input"
      />
      <span class="settings-unit">min</span>
    </div>
  </div>

  <div class="settings-field">
    <label class="settings-label" for="increment">Time increment</label>
    <div class="settings-input-row">
      <input
        id="increment"
        type="number"
        min="1"
        max="60"
        bind:value={timer.settingsIncrementMin}
        class="settings-input"
      />
      <span class="settings-unit">min</span>
    </div>
  </div>

  <div class="settings-field">
    <label class="settings-label" for="secondary-increment">Secondary increment</label>
    <div class="settings-input-row">
      <input
        id="secondary-increment"
        type="number"
        min="1"
        max="60"
        bind:value={timer.settingsSecondaryIncrementMin}
        class="settings-input"
      />
      <span class="settings-unit">min</span>
    </div>
  </div>

  <div class="settings-field">
    <label class="settings-label" for="tertiary-increment">Tertiary increment</label>
    <div class="settings-input-row">
      <input
        id="tertiary-increment"
        type="number"
        min="1"
        max="1440"
        bind:value={timer.settingsTertiaryIncrementMin}
        class="settings-input"
      />
      <span class="settings-unit">min</span>
    </div>
  </div>

  <div class="settings-field">
    <label class="settings-label" for="mini-window-width">Mini window width</label>
    <div class="settings-input-row">
      <input
        id="mini-window-width"
        type="number"
        min="50"
        max="400"
        bind:value={timer.settingsMiniWindowWidth}
        class="settings-input"
      />
      <span class="settings-unit">px</span>
    </div>
  </div>

  <div class="settings-field">
    <label class="settings-label" for="mini-window-height">Mini window height</label>
    <div class="settings-input-row">
      <input
        id="mini-window-height"
        type="number"
        min="50"
        max="400"
        bind:value={timer.settingsMiniWindowHeight}
        class="settings-input"
      />
      <span class="settings-unit">px</span>
    </div>
  </div>

  <button class="settings-save-btn" onclick={() => timer.saveSettings()}>
    Save
  </button>
</div>
