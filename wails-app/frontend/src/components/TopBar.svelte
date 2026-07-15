<script lang="ts">
  import { timer } from "../lib/timerStore.svelte";
  import { TINTS, tintSwatch } from "../lib/tints";
</script>

<!-- Pin toggle -->
<button
  class="pin-btn"
  class:pinned={timer.alwaysOnTop}
  onclick={() => timer.toggleAlwaysOnTop()}
  aria-label={timer.alwaysOnTop ? "Unpin window" : "Pin window on top"}
>
  <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
    {#if timer.alwaysOnTop}
      <path d="M5.5 2.5L10.5 2.5L11 7L13 8.5V10H8.5V14.5L8 15.5L7.5 14.5V10H3V8.5L5 7Z" fill="currentColor" />
    {:else}
      <path d="M5.5 2.5L10.5 2.5L11 7L13 8.5V10H8.5V14.5L8 15.5L7.5 14.5V10H3V8.5L5 7Z" />
    {/if}
  </svg>
</button>

<!-- Settings toggle -->
<button class="settings-btn" onclick={() => timer.openSettings()} aria-label="Settings">
  <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
    <circle cx="8" cy="8" r="2" />
    <path d="M13.5 8a5.5 5.5 0 0 0-.1-.8l1.3-1-.7-1.2-1.5.5a5.5 5.5 0 0 0-1.2-.7L11 3.3h-1.4l-.3 1.5a5.5 5.5 0 0 0-1.2.7l-1.5-.5-.7 1.2 1.3 1a5.5 5.5 0 0 0 0 1.6l-1.3 1 .7 1.2 1.5-.5a5.5 5.5 0 0 0 1.2.7l.3 1.5H11l.3-1.5a5.5 5.5 0 0 0 1.2-.7l1.5.5.7-1.2-1.3-1a5.5 5.5 0 0 0 .1-.8z" />
  </svg>
</button>

<!-- Mini mode button -->
<button class="mini-btn" onclick={() => timer.toggleMiniMode()} aria-label="Toggle mini mode">
  Mini
</button>

<!-- Add timer button -->
<button class="add-timer-btn" onclick={() => timer.handleAddTimer()} aria-label="Add timer">
  <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
    <line x1="6" y1="1" x2="6" y2="11" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
    <line x1="1" y1="6" x2="11" y2="6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
  </svg>
</button>

<!-- Tint picker -->
<button class="tint-btn" onclick={() => (timer.showTintMenu = !timer.showTintMenu)} aria-label="Background tint">
  <span class="tint-btn-dot" style={`background:${tintSwatch(timer.activeTintId)}`}></span>
</button>
{#if timer.showTintMenu}
  <button class="tint-backdrop" onclick={() => (timer.showTintMenu = false)} aria-label="Close tint menu"></button>
  <div class="tint-popover">
    {#each TINTS as t}
      <button
        class="tint-swatch"
        class:selected={t.id === timer.activeTintId}
        class:none={t.id === "none"}
        style={`background:${t.swatch}`}
        onclick={() => timer.handleSetTint(t.id)}
        aria-label={t.label}
        title={t.label}
      ></button>
    {/each}
  </div>
{/if}

<!-- Minimize button -->
<button class="window-btn minimize-btn" onclick={() => timer.handleMinimize()} aria-label="Minimize">
  <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
    <line x1="2" y1="6" x2="10" y2="6" />
  </svg>
</button>

<!-- Close button -->
<button class="window-btn close-btn" onclick={() => timer.handleClose()} aria-label="Close">
  <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
    <line x1="2" y1="2" x2="10" y2="10" />
    <line x1="10" y1="2" x2="2" y2="10" />
  </svg>
</button>

<!-- Mini pin indicator (compact/tiny mode only, shown when pinned) -->
{#if timer.alwaysOnTop}
  <button class="mini-pin-indicator" onclick={() => timer.toggleAlwaysOnTop()} aria-label="Unpin window">
    <svg width="10" height="10" viewBox="0 0 16 16" fill="currentColor">
      <path d="M5.5 2.5L10.5 2.5L11 7L13 8.5V10H8.5V14.5L8 15.5L7.5 14.5V10H3V8.5L5 7Z" />
    </svg>
  </button>
{/if}
