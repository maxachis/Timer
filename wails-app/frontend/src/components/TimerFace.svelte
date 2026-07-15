<script lang="ts">
  import { timer, RADIUS, CIRCUMFERENCE, formatMinutes, formatSeconds } from "../lib/timerStore.svelte";
  import { handleNameKeydown } from "../lib/keyboard";
</script>

<!-- Timer face -->
<div class="timer-face" class:has-mini-label={timer.hasMultipleTimers} style="--wails-draggable:drag">
  {#if timer.onBreak}
    <div class="break-label">BREAK</div>
  {/if}
  <svg class="progress-ring" viewBox="0 0 300 300" style="--wails-draggable:drag; pointer-events:none">
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
      stroke-dashoffset={timer.strokeOffset}
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
  <div class="time-display" style="--wails-draggable:drag">
    <span class="digits" style="--wails-draggable:drag">{formatMinutes(timer.remainingSecs)}</span>
    <span class="separator" style="--wails-draggable:drag">:</span>
    <span class="digits" style="--wails-draggable:drag">{formatSeconds(timer.remainingSecs)}</span>
  </div>

  {#if timer.isFinished}
    <div class="finished-label">complete</div>
  {:else if timer.timerState === "paused"}
    <div class="state-label">paused</div>
  {/if}

  <!-- Timer name (only when multiple timers) -->
  {#if timer.hasMultipleTimers}
    <div class="timer-name-area">
      {#if timer.editingName}
        <input
          class="timer-name-input"
          type="text"
          bind:value={timer.editNameValue}
          onblur={() => timer.finishEditName()}
          onkeydown={handleNameKeydown}
          autofocus
        />
      {:else}
        <button class="timer-name-label" onclick={() => timer.startEditName()} aria-label="Rename timer">
          {timer.activeName || `Timer ${timer.activeIndex + 1}`}
        </button>
      {/if}
    </div>

    <!-- Pagination dots -->
    <div class="timer-dots">
      {#each timer.timerList as entry}
        <button
          class="timer-dot"
          class:active={entry.is_active}
          onclick={() => timer.handleSwitchTimer(entry.index)}
          aria-label={`Switch to ${entry.name || "Timer " + (entry.index + 1)}`}
        ></button>
      {/each}
    </div>
  {/if}

  <!-- Compact timer label (shown only in compact/tiny view when multiple timers) -->
  {#if timer.hasMultipleTimers}
    <div class="compact-timer-label">
      {timer.activeName || `T${timer.activeIndex + 1}`}
    </div>
  {/if}
</div>
