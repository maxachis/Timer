<script lang="ts">
  import { timer } from "./lib/timerStore.svelte";
  import { handleKeydown } from "./lib/keyboard";
  import TimerFace from "./components/TimerFace.svelte";
  import Controls from "./components/Controls.svelte";
  import TopBar from "./components/TopBar.svelte";
  import Settings from "./components/Settings.svelte";
  import History from "./components/History.svelte";
  import Toast from "./components/Toast.svelte";
  import StackedBar from "./components/StackedBar.svelte";

  timer.fetchStatus();
  timer.loadSettings();
  timer.fetchTimerList();
</script>

<svelte:window onkeydown={handleKeydown} />

<svelte:head>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link href="https://fonts.googleapis.com/css2?family=DM+Serif+Display&family=Anybody:wght@300;400;500&display=swap" rel="stylesheet" />
</svelte:head>

<main style={`${timer.showSettings ? "--wails-draggable:no-drag" : "--wails-draggable:drag"};--timer-tint:${timer.activeTintWash}`} class:finished={timer.isFinished} class:running={timer.timerState === "running"} class:paused={timer.timerState === "paused"} class:break={timer.onBreak} class:settings-view={timer.showSettings}>
  {#if timer.showSettings}
    <Settings />
  {:else}
    <TopBar />

    <!-- Subtle texture overlay -->
    <div class="texture"></div>

    <div class="stage">
      <div class="stage-main">
        <TimerFace />
        <Controls />
      </div>
      <StackedBar />
    </div>
  {/if}

  {#if timer.showHistory}
    <History />
  {/if}

  <Toast />
</main>
