<script lang="ts">
  import { timer, formatMinutes, formatSeconds } from "../lib/timerStore.svelte";
  import { barColor } from "../lib/tints";

  const total = $derived(
    timer.timerList.reduce((s, e) => s + e.original_duration_secs, 0) || 1
  );

  function hexToRgba(hex: string, alpha: number): string {
    const clean = hex.replace("#", "");
    if (clean.length !== 6) return `rgba(140,140,140,${alpha})`;
    const r = parseInt(clean.slice(0, 2), 16);
    const g = parseInt(clean.slice(2, 4), 16);
    const b = parseInt(clean.slice(4, 6), 16);
    return `rgba(${r},${g},${b},${alpha})`;
  }
</script>

<div class="stacked-bar">
  {#each timer.timerList as entry, index (entry.index)}
    {@const fillPct = entry.original_duration_secs > 0
      ? Math.min(1, entry.remaining_secs / entry.original_duration_secs) * 100
      : 0}
    {@const color = barColor(entry, index)}
    <button
      class="bar-lane"
      class:active={entry.is_active}
      class:break-source={entry.is_break_timer}
      style={`flex-grow:${entry.original_duration_secs / total};background:${hexToRgba(color, 0.12)}`}
      onclick={() => timer.handleSwitchTimer(entry.index)}
      aria-label={`Switch to ${entry.name || "Timer " + (entry.index + 1)}`}
    >
      {#if entry.is_break_timer}
        <span class="bar-break-badge">B</span>
      {/if}
      <div
        class="bar-fill"
        style={`height:${fillPct}%;background:${color};opacity:${entry.is_active ? 1 : 0.7}`}
      ></div>
      <div class="bar-tooltip">
        {entry.name || `Timer ${entry.index + 1}`} · {formatMinutes(entry.remaining_secs)}:{formatSeconds(entry.remaining_secs)}
      </div>
    </button>
  {/each}
</div>
