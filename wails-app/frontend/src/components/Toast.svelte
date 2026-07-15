<script lang="ts">
  import { timer } from "../lib/timerStore.svelte";
</script>

<!-- Toast notification -->
{#if timer.toastMessage}
  <div class="toast">{timer.toastMessage}</div>
{/if}
