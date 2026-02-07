<script lang="ts">
  import { invoke } from "@tauri-apps/api/core";
  import { getCurrentWindow } from "@tauri-apps/api/window";
  import {
    isPermissionGranted,
    requestPermission,
    sendNotification,
  } from "@tauri-apps/plugin-notification";

  interface TimerStatus {
    remaining_secs: number;
    state: string;
    is_finished: boolean;
  }

  let remainingSecs = $state(300);
  let timerState = $state("idle");
  let isFinished = $state(false);
  let pollInterval: ReturnType<typeof setInterval> | null = $state(null);
  let totalDuration = $state(300);
  let prevState = $state("idle");
  let alwaysOnTop = $state(false);
  let showSettings = $state(false);
  let settingsDurationMin = $state(5);
  let settingsIncrementMin = $state(5);
  let incrementSecs = $state(300);

  let incrementLabel = $derived(Math.round(incrementSecs / 60) + " min");

  async function loadSettings() {
    const s = await invoke<{ default_duration_secs: number; default_increment_secs: number }>("get_settings");
    settingsDurationMin = Math.round(s.default_duration_secs / 60);
    settingsIncrementMin = Math.round(s.default_increment_secs / 60);
    incrementSecs = s.default_increment_secs;
  }

  async function openSettings() {
    await loadSettings();
    showSettings = true;
  }

  function closeSettings() {
    showSettings = false;
  }

  async function saveSettings() {
    const durationMin = Math.max(1, Math.min(180, settingsDurationMin));
    const incrementMin = Math.max(1, Math.min(60, settingsIncrementMin));
    const newSettings = {
      default_duration_secs: durationMin * 60,
      default_increment_secs: incrementMin * 60,
    };
    await invoke("save_settings", { newSettings });
    incrementSecs = newSettings.default_increment_secs;
    showSettings = false;
    if (timerState === "idle") {
      await invoke("create_timer", { durationSecs: newSettings.default_duration_secs });
      await fetchStatus();
    }
  }

  async function toggleAlwaysOnTop() {
    alwaysOnTop = !alwaysOnTop;
    await getCurrentWindow().setAlwaysOnTop(alwaysOnTop);
  }

  // Track total duration for the progress ring
  let progress = $derived(
    totalDuration > 0 ? remainingSecs / totalDuration : 0
  );

  // SVG circle math
  const RADIUS = 130;
  const CIRCUMFERENCE = 2 * Math.PI * RADIUS;
  let strokeOffset = $derived(CIRCUMFERENCE * (1 - progress));

  function formatMinutes(secs: number): string {
    return String(Math.floor(secs / 60)).padStart(2, "0");
  }

  function formatSeconds(secs: number): string {
    return String(Math.floor(secs % 60)).padStart(2, "0");
  }

  function playAlertSound() {
    const ctx = new AudioContext();
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();
    osc.connect(gain);
    gain.connect(ctx.destination);
    osc.type = "sine";
    osc.frequency.setValueAtTime(880, ctx.currentTime);
    osc.frequency.setValueAtTime(660, ctx.currentTime + 0.15);
    osc.frequency.setValueAtTime(880, ctx.currentTime + 0.3);
    gain.gain.setValueAtTime(0.3, ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.5);
    osc.start(ctx.currentTime);
    osc.stop(ctx.currentTime + 0.5);
  }

  async function sendCompletionNotification() {
    let permitted = await isPermissionGranted();
    if (!permitted) {
      const permission = await requestPermission();
      permitted = permission === "granted";
    }
    if (permitted) {
      sendNotification({
        title: "Timer Complete",
        body: "Your timer has finished.",
      });
    }
  }

  async function fetchStatus() {
    const status = await invoke<TimerStatus>("get_timer_status");
    remainingSecs = status.remaining_secs;
    timerState = status.state;
    isFinished = status.is_finished;

    // Detect transition to finished state
    if (prevState !== "finished" && status.state === "finished") {
      playAlertSound();
      sendCompletionNotification();
    }
    prevState = status.state;

    // Update total duration when idle (after reset or create)
    if (status.state === "idle") {
      totalDuration = status.remaining_secs;
    }

    if (status.state !== "running" && pollInterval) {
      clearInterval(pollInterval);
      pollInterval = null;
    }
  }

  function startPolling() {
    if (!pollInterval) {
      pollInterval = setInterval(fetchStatus, 80);
    }
  }

  async function handleStart() {
    await invoke("start_timer");
    await fetchStatus();
    startPolling();
  }

  async function handlePause() {
    await invoke("pause_timer");
    await fetchStatus();
  }

  async function handleResume() {
    await invoke("resume_timer");
    await fetchStatus();
    startPolling();
  }

  async function handleReset() {
    await invoke("reset_timer");
    await fetchStatus();
  }

  async function handleAddTime() {
    await invoke("add_time");
    totalDuration += incrementSecs;
    await fetchStatus();
  }

  async function handleRemoveTime() {
    await invoke("remove_time");
    await fetchStatus();
  }

  function handleKeydown(event: KeyboardEvent) {
    if (showSettings) return;
    switch (event.key) {
      case " ":
        event.preventDefault();
        if (timerState === "idle") handleStart();
        else if (timerState === "running") handlePause();
        else if (timerState === "paused") handleResume();
        else if (timerState === "finished") handleReset();
        break;
      case "r":
      case "R":
      case "ArrowDown":
        if (timerState !== "idle") {
          event.preventDefault();
          handleReset();
        }
        break;
      case "+":
      case "=":
      case "ArrowRight":
        if (timerState === "running" || timerState === "paused") {
          event.preventDefault();
          handleAddTime();
        }
        break;
      case "-":
      case "ArrowLeft":
        if (timerState === "running" || timerState === "paused") {
          event.preventDefault();
          handleRemoveTime();
        }
        break;
      case "Escape":
        if (timerState === "finished") {
          event.preventDefault();
          handleReset();
        }
        break;
      case "t":
      case "T":
        event.preventDefault();
        toggleAlwaysOnTop();
        break;
    }
  }

  fetchStatus();
  loadSettings();
</script>

<svelte:window onkeydown={handleKeydown} />

<svelte:head>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link href="https://fonts.googleapis.com/css2?family=DM+Serif+Display&family=Anybody:wght@300;400;500&display=swap" rel="stylesheet" />
</svelte:head>

<main class:finished={isFinished} class:running={timerState === "running"} class:paused={timerState === "paused"}>
  {#if showSettings}
    <!-- Settings view -->
    <div class="texture"></div>

    <button class="settings-back-btn" onclick={closeSettings} aria-label="Back to timer">
      <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <path d="M10 2L4 8L10 14" />
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
            bind:value={settingsDurationMin}
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
            bind:value={settingsIncrementMin}
            class="settings-input"
          />
          <span class="settings-unit">min</span>
        </div>
      </div>

      <button class="settings-save-btn" onclick={saveSettings}>
        Save
      </button>
    </div>
  {:else}
    <!-- Pin toggle -->
    <button
      class="pin-btn"
      class:pinned={alwaysOnTop}
      onclick={toggleAlwaysOnTop}
      aria-label={alwaysOnTop ? "Unpin window" : "Pin window on top"}
    >
      <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        {#if alwaysOnTop}
          <path d="M5.5 2.5L10.5 2.5L11 7L13 8.5V10H8.5V14.5L8 15.5L7.5 14.5V10H3V8.5L5 7Z" fill="currentColor" />
        {:else}
          <path d="M5.5 2.5L10.5 2.5L11 7L13 8.5V10H8.5V14.5L8 15.5L7.5 14.5V10H3V8.5L5 7Z" />
        {/if}
      </svg>
    </button>

    <!-- Settings toggle -->
    <button class="settings-btn" onclick={openSettings} aria-label="Settings">
      <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="8" cy="8" r="2" />
        <path d="M13.5 8a5.5 5.5 0 0 0-.1-.8l1.3-1-.7-1.2-1.5.5a5.5 5.5 0 0 0-1.2-.7L11 3.3h-1.4l-.3 1.5a5.5 5.5 0 0 0-1.2.7l-1.5-.5-.7 1.2 1.3 1a5.5 5.5 0 0 0 0 1.6l-1.3 1 .7 1.2 1.5-.5a5.5 5.5 0 0 0 1.2.7l.3 1.5H11l.3-1.5a5.5 5.5 0 0 0 1.2-.7l1.5.5.7-1.2-1.3-1a5.5 5.5 0 0 0 .1-.8z" />
      </svg>
    </button>

    <!-- Subtle texture overlay -->
    <div class="texture"></div>

    <!-- Timer face -->
    <div class="timer-face">
      <svg class="progress-ring" viewBox="0 0 300 300">
        <!-- Track -->
        <circle
          cx="150" cy="150" r={RADIUS}
          fill="none"
          stroke="var(--ring-track)"
          stroke-width="4"
        />
        <!-- Progress arc -->
        <circle
          class="progress-arc"
          cx="150" cy="150" r={RADIUS}
          fill="none"
          stroke="var(--ring-color)"
          stroke-width="6"
          stroke-linecap="round"
          stroke-dasharray={CIRCUMFERENCE}
          stroke-dashoffset={strokeOffset}
          transform="rotate(-90 150 150)"
        />
        <!-- Tick marks -->
        {#each Array(60) as _, i}
          {@const angle = (i / 60) * 360 - 90}
          {@const isMajor = i % 5 === 0}
          {@const r1 = isMajor ? 116 : 120}
          {@const r2 = 124}
          {@const x1 = 150 + r1 * Math.cos((angle * Math.PI) / 180)}
          {@const y1 = 150 + r1 * Math.sin((angle * Math.PI) / 180)}
          {@const x2 = 150 + r2 * Math.cos((angle * Math.PI) / 180)}
          {@const y2 = 150 + r2 * Math.sin((angle * Math.PI) / 180)}
          <line
            {x1} {y1} {x2} {y2}
            stroke="var(--tick-color)"
            stroke-width={isMajor ? 1.5 : 0.5}
            opacity={isMajor ? 0.5 : 0.2}
          />
        {/each}
      </svg>

      <!-- Time display -->
      <div class="time-display">
        <span class="digits">{formatMinutes(remainingSecs)}</span>
        <span class="separator">:</span>
        <span class="digits">{formatSeconds(remainingSecs)}</span>
      </div>

      {#if isFinished}
        <div class="finished-label">complete</div>
      {:else if timerState === "paused"}
        <div class="state-label">paused</div>
      {/if}
    </div>

    <!-- Controls -->
    <div class="controls">
      <!-- Time adjust -->
      <div class="adjust-row">
        <button
          class="adjust-btn"
          onclick={handleRemoveTime}
          disabled={timerState === "idle" || isFinished}
          aria-label="Remove time"
        >
          <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
            <line x1="2" y1="7" x2="12" y2="7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
          <span>{incrementLabel}</span>
        </button>

        <!-- Primary action -->
        {#if timerState === "idle"}
          <button class="primary-btn" onclick={handleStart} aria-label="Start timer">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
              <polygon points="5,3 17,10 5,17"/>
            </svg>
          </button>
        {:else if timerState === "running"}
          <button class="primary-btn active" onclick={handlePause} aria-label="Pause timer">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
              <rect x="4" y="3" width="4" height="14" rx="1"/>
              <rect x="12" y="3" width="4" height="14" rx="1"/>
            </svg>
          </button>
        {:else if timerState === "paused"}
          <button class="primary-btn" onclick={handleResume} aria-label="Resume timer">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
              <polygon points="5,3 17,10 5,17"/>
            </svg>
          </button>
        {:else}
          <button class="primary-btn" onclick={handleReset} aria-label="Reset timer">
            <svg width="18" height="18" viewBox="0 0 18 18" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
              <path d="M3 9a6 6 0 1 1 1.5 4"/>
              <polyline points="1,7 3,9.5 5.5,7.5"/>
            </svg>
          </button>
        {/if}

        <button
          class="adjust-btn"
          onclick={handleAddTime}
          disabled={timerState === "idle" || isFinished}
          aria-label="Add time"
        >
          <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
            <line x1="7" y1="2" x2="7" y2="12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            <line x1="2" y1="7" x2="12" y2="7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
          <span>{incrementLabel}</span>
        </button>
      </div>

      <!-- Reset -->
      {#if timerState !== "idle" && !isFinished}
        <button class="reset-btn" onclick={handleReset}>
          Reset
        </button>
      {/if}
    </div>
  {/if}
</main>

<style>
  :root {
    --cream: #f5f0e8;
    --warm-white: #faf7f2;
    --charcoal: #2c2c2c;
    --stone: #8a8578;
    --brass: #c4a265;
    --brass-dark: #a68942;
    --brass-glow: #dbb978;
    --ring-track: #e8e2d8;
    --ring-color: #c4a265;
    --tick-color: #b8b0a4;
    --shadow-soft: 0 2px 20px rgba(44, 44, 44, 0.06);
    --shadow-btn: 0 1px 3px rgba(44, 44, 44, 0.1), 0 4px 12px rgba(44, 44, 44, 0.05);
    --shadow-btn-hover: 0 2px 6px rgba(44, 44, 44, 0.12), 0 8px 24px rgba(44, 44, 44, 0.08);
    --shadow-primary: 0 2px 8px rgba(196, 162, 101, 0.3), 0 6px 20px rgba(196, 162, 101, 0.15);
    --amber-alert: #d4763a;
    --amber-glow: rgba(212, 118, 58, 0.08);
  }

  :global(*) {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
  }

  :global(body) {
    margin: 0;
    background: var(--warm-white);
    color: var(--charcoal);
    font-family: "Anybody", sans-serif;
    font-weight: 400;
    -webkit-font-smoothing: antialiased;
    overflow: hidden;
    user-select: none;
  }

  main {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    padding: 2rem 1.5rem;
    position: relative;
    transition: background-color 0.6s ease;
  }

  main.finished {
    background-color: var(--amber-glow);
    --ring-color: var(--amber-alert);
  }

  .texture {
    position: fixed;
    inset: 0;
    background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)' opacity='0.03'/%3E%3C/svg%3E");
    pointer-events: none;
    z-index: 0;
  }

  /* Timer face */
  .timer-face {
    position: relative;
    width: 300px;
    height: 300px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    z-index: 1;
    margin-bottom: 2.5rem;
  }

  .progress-ring {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
  }

  .progress-arc {
    transition: stroke-dashoffset 0.15s ease-out;
    filter: drop-shadow(0 0 6px rgba(196, 162, 101, 0.25));
  }

  main.finished .progress-arc {
    filter: drop-shadow(0 0 10px rgba(212, 118, 58, 0.4));
    animation: pulse-ring 2s ease-in-out infinite;
  }

  @keyframes pulse-ring {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.6; }
  }

  /* Time digits */
  .time-display {
    display: flex;
    align-items: baseline;
    gap: 0;
    line-height: 1;
    margin-top: -4px;
  }

  .digits {
    font-family: "DM Serif Display", serif;
    font-size: 4.5rem;
    color: var(--charcoal);
    letter-spacing: -0.02em;
    transition: color 0.4s ease;
  }

  .separator {
    font-family: "DM Serif Display", serif;
    font-size: 4rem;
    color: var(--stone);
    margin: 0 0.1rem;
    opacity: 0.6;
  }

  main.running .separator {
    animation: blink 1.2s step-end infinite;
  }

  @keyframes blink {
    0%, 100% { opacity: 0.6; }
    50% { opacity: 0.15; }
  }

  main.finished .digits {
    color: var(--amber-alert);
  }

  main.finished .separator {
    color: var(--amber-alert);
    opacity: 0.5;
  }

  /* State labels */
  .finished-label,
  .state-label {
    font-family: "Anybody", sans-serif;
    font-weight: 300;
    font-size: 0.85rem;
    letter-spacing: 0.2em;
    text-transform: uppercase;
    margin-top: 0.75rem;
  }

  .finished-label {
    color: var(--amber-alert);
    animation: fade-in 0.4s ease;
  }

  .state-label {
    color: var(--stone);
    animation: fade-in 0.3s ease;
  }

  @keyframes fade-in {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }

  /* Controls */
  .controls {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.25rem;
    z-index: 1;
  }

  .adjust-row {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }

  /* Adjust buttons */
  .adjust-btn {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.3rem;
    padding: 0.6rem 0.8rem;
    border: 1px solid #e0dbd3;
    border-radius: 12px;
    background: white;
    color: var(--stone);
    cursor: pointer;
    font-family: "Anybody", sans-serif;
    font-size: 0.7rem;
    font-weight: 500;
    letter-spacing: 0.05em;
    text-transform: uppercase;
    box-shadow: var(--shadow-btn);
    transition: all 0.2s ease;
    min-width: 56px;
  }

  .adjust-btn span {
    line-height: 1;
  }

  .adjust-btn:hover:not(:disabled) {
    box-shadow: var(--shadow-btn-hover);
    border-color: #d4cfc6;
    color: var(--charcoal);
    transform: translateY(-1px);
  }

  .adjust-btn:active:not(:disabled) {
    transform: translateY(0);
    box-shadow: 0 1px 2px rgba(44, 44, 44, 0.08);
  }

  .adjust-btn:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }

  /* Primary button */
  .primary-btn {
    width: 64px;
    height: 64px;
    border-radius: 50%;
    border: none;
    background: linear-gradient(145deg, var(--brass-glow), var(--brass-dark));
    color: white;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: var(--shadow-primary);
    transition: all 0.2s ease;
  }

  .primary-btn:hover {
    transform: scale(1.06);
    box-shadow: 0 3px 12px rgba(196, 162, 101, 0.4), 0 8px 28px rgba(196, 162, 101, 0.2);
  }

  .primary-btn:active {
    transform: scale(0.97);
    box-shadow: 0 1px 4px rgba(196, 162, 101, 0.3);
  }

  .primary-btn.active {
    background: linear-gradient(145deg, #8a8578, #6a6560);
    box-shadow: 0 2px 8px rgba(44, 44, 44, 0.2);
  }

  .primary-btn svg {
    flex-shrink: 0;
  }

  /* Reset text button */
  .reset-btn {
    border: none;
    background: none;
    color: var(--stone);
    font-family: "Anybody", sans-serif;
    font-size: 0.8rem;
    font-weight: 400;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    cursor: pointer;
    padding: 0.4rem 0.8rem;
    border-radius: 6px;
    transition: all 0.2s ease;
    animation: fade-in 0.3s ease;
  }

  .reset-btn:hover {
    color: var(--charcoal);
    background: rgba(44, 44, 44, 0.04);
  }

  /* Pin button */
  .pin-btn {
    position: absolute;
    top: 1rem;
    right: 1rem;
    z-index: 2;
    width: 32px;
    height: 32px;
    border: 1px solid #e0dbd3;
    border-radius: 8px;
    background: white;
    color: var(--stone);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: var(--shadow-btn);
    transition: all 0.2s ease;
    opacity: 0.6;
  }

  .pin-btn:hover {
    opacity: 1;
    box-shadow: var(--shadow-btn-hover);
    transform: translateY(-1px);
  }

  .pin-btn.pinned {
    opacity: 1;
    color: var(--brass);
    border-color: var(--brass);
    background: rgba(196, 162, 101, 0.08);
  }

  /* Settings button (cog) */
  .settings-btn {
    position: absolute;
    top: 1rem;
    left: 1rem;
    z-index: 2;
    width: 32px;
    height: 32px;
    border: 1px solid #e0dbd3;
    border-radius: 8px;
    background: white;
    color: var(--stone);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: var(--shadow-btn);
    transition: all 0.2s ease;
    opacity: 0.6;
  }

  .settings-btn:hover {
    opacity: 1;
    box-shadow: var(--shadow-btn-hover);
    transform: translateY(-1px);
  }

  /* Settings back button */
  .settings-back-btn {
    position: absolute;
    top: 1rem;
    left: 1rem;
    z-index: 2;
    width: 32px;
    height: 32px;
    border: 1px solid #e0dbd3;
    border-radius: 8px;
    background: white;
    color: var(--stone);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: var(--shadow-btn);
    transition: all 0.2s ease;
    opacity: 0.6;
  }

  .settings-back-btn:hover {
    opacity: 1;
    box-shadow: var(--shadow-btn-hover);
    transform: translateY(-1px);
  }

  /* Settings panel */
  .settings-panel {
    z-index: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.75rem;
    animation: fade-in 0.3s ease;
  }

  .settings-title {
    font-family: "DM Serif Display", serif;
    font-size: 1.5rem;
    font-weight: 400;
    color: var(--charcoal);
    letter-spacing: 0.02em;
  }

  .settings-field {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    width: 200px;
  }

  .settings-label {
    font-family: "Anybody", sans-serif;
    font-size: 0.75rem;
    font-weight: 500;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--stone);
  }

  .settings-input-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .settings-input {
    width: 100%;
    padding: 0.6rem 0.8rem;
    border: 1px solid #e0dbd3;
    border-radius: 10px;
    background: white;
    color: var(--charcoal);
    font-family: "Anybody", sans-serif;
    font-size: 1rem;
    font-weight: 400;
    box-shadow: var(--shadow-btn);
    outline: none;
    transition: border-color 0.2s ease;
    -moz-appearance: textfield;
  }

  .settings-input:focus {
    border-color: var(--brass);
  }

  .settings-input::-webkit-outer-spin-button,
  .settings-input::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }

  .settings-unit {
    font-family: "Anybody", sans-serif;
    font-size: 0.8rem;
    font-weight: 400;
    color: var(--stone);
    flex-shrink: 0;
  }

  .settings-save-btn {
    margin-top: 0.5rem;
    padding: 0.7rem 2.5rem;
    border: none;
    border-radius: 12px;
    background: linear-gradient(145deg, var(--brass-glow), var(--brass-dark));
    color: white;
    font-family: "Anybody", sans-serif;
    font-size: 0.85rem;
    font-weight: 500;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    cursor: pointer;
    box-shadow: var(--shadow-primary);
    transition: all 0.2s ease;
  }

  .settings-save-btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 3px 12px rgba(196, 162, 101, 0.4), 0 8px 28px rgba(196, 162, 101, 0.2);
  }

  .settings-save-btn:active {
    transform: translateY(0);
    box-shadow: 0 1px 4px rgba(196, 162, 101, 0.3);
  }
</style>
