<script lang="ts">
  import { timer } from "../lib/timerStore.svelte";
  import { groupByDay, formatEntryTime, KIND_LABELS } from "../lib/history";

  const groups = $derived(groupByDay(timer.historyEntries));
</script>

<button class="history-backdrop" onclick={() => timer.closeHistory()} aria-label="Close history"></button>
<div class="history-modal" role="dialog" aria-modal="true" aria-labelledby="history-title">
  <div class="history-header">
    <h2 id="history-title" class="history-title">History</h2>
    <button class="history-close-btn" onclick={() => timer.closeHistory()} aria-label="Close history">
      <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
        <line x1="2" y1="2" x2="10" y2="10" />
        <line x1="10" y1="2" x2="2" y2="10" />
      </svg>
    </button>
  </div>

  <div class="history-list">
    {#if timer.historyEntries.length === 0}
      <p class="history-empty">No activity yet</p>
    {:else}
      {#each groups as group (group.label)}
        <div class="history-day">{group.label}</div>
        {#each group.entries as entry (entry.at + entry.kind + entry.timer_name)}
          <div class="history-row" data-kind={entry.kind}>
            <span class="history-time">{formatEntryTime(entry.at)}</span>
            <span class="history-kind">{KIND_LABELS[entry.kind]}</span>
            <span class="history-name">{entry.timer_name}</span>
          </div>
        {/each}
      {/each}
    {/if}
  </div>
</div>
