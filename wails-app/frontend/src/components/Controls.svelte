<script lang="ts">
  import { timer } from "../lib/timerStore.svelte";
</script>

<!-- Controls -->
<div class="controls">
  <!-- Tertiary time adjust -->
  <div class="adjust-row tertiary">
    <button
      class="adjust-btn tertiary"
      onclick={() => timer.handleRemoveTimeTertiary()}
      disabled={timer.timerState === "idle" || timer.isFinished}
      aria-label="Remove time (tertiary)"
    >
      <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
        <line x1="2" y1="7" x2="12" y2="7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
      </svg>
      <span>{timer.tertiaryIncrementLabel}</span>
    </button>

    <button
      class="adjust-btn tertiary"
      onclick={() => timer.handleAddTimeTertiary()}
      disabled={timer.timerState === "idle" || timer.isFinished}
      aria-label="Add time (tertiary)"
    >
      <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
        <line x1="7" y1="2" x2="7" y2="12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
        <line x1="2" y1="7" x2="12" y2="7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
      </svg>
      <span>{timer.tertiaryIncrementLabel}</span>
    </button>
  </div>

  <!-- Time adjust -->
  <div class="adjust-row">
    <button
      class="adjust-btn"
      onclick={() => timer.handleRemoveTime()}
      disabled={timer.timerState === "idle" || timer.isFinished}
      aria-label="Remove time"
    >
      <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
        <line x1="2" y1="7" x2="12" y2="7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
      </svg>
      <span>{timer.incrementLabel}</span>
    </button>

    <!-- Primary action -->
    {#if timer.timerState === "idle"}
      <button class="primary-btn" onclick={() => timer.handleStart()} aria-label="Start timer">
        <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
          <polygon points="5,3 17,10 5,17"/>
        </svg>
      </button>
    {:else if timer.timerState === "running"}
      <button class="primary-btn active" onclick={() => timer.handlePause()} aria-label="Pause timer">
        <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
          <rect x="4" y="3" width="4" height="14" rx="1"/>
          <rect x="12" y="3" width="4" height="14" rx="1"/>
        </svg>
      </button>
    {:else if timer.timerState === "paused"}
      <button class="primary-btn" onclick={() => timer.handleResume()} aria-label="Resume timer">
        <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
          <polygon points="5,3 17,10 5,17"/>
        </svg>
      </button>
    {:else}
      <button class="primary-btn" onclick={() => timer.handleReset()} aria-label="Reset timer">
        <svg width="18" height="18" viewBox="0 0 18 18" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
          <path d="M3 9a6 6 0 1 1 1.5 4"/>
          <polyline points="1,7 3,9.5 5.5,7.5"/>
        </svg>
      </button>
    {/if}

    <button
      class="adjust-btn"
      onclick={() => timer.handleAddTime()}
      disabled={timer.timerState === "idle" || timer.isFinished}
      aria-label="Add time"
    >
      <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
        <line x1="7" y1="2" x2="7" y2="12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
        <line x1="2" y1="7" x2="12" y2="7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
      </svg>
      <span>{timer.incrementLabel}</span>
    </button>
  </div>

  <!-- Secondary time adjust -->
  <div class="adjust-row secondary">
    <button
      class="adjust-btn secondary"
      onclick={() => timer.handleRemoveTimeSecondary()}
      disabled={timer.timerState === "idle" || timer.isFinished}
      aria-label="Remove time (secondary)"
    >
      <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
        <line x1="2" y1="7" x2="12" y2="7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
      </svg>
      <span>{timer.secondaryIncrementLabel}</span>
    </button>

    <button
      class="adjust-btn secondary"
      onclick={() => timer.handleAddTimeSecondary()}
      disabled={timer.timerState === "idle" || timer.isFinished}
      aria-label="Add time (secondary)"
    >
      <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
        <line x1="7" y1="2" x2="7" y2="12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
        <line x1="2" y1="7" x2="12" y2="7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
      </svg>
      <span>{timer.secondaryIncrementLabel}</span>
    </button>
  </div>

  <!-- Remove timer (only when multiple) -->
  {#if timer.hasMultipleTimers}
    <button class="remove-timer-btn" onclick={() => timer.handleRemoveTimer(timer.activeIndex)}>
      Remove timer
    </button>
  {/if}

  <!-- Reset -->
  {#if timer.timerState !== "idle" && !timer.isFinished}
    <button class="reset-btn" onclick={() => timer.handleReset()}>
      Reset
    </button>
  {/if}
</div>
