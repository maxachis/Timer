import { timer } from "./timerStore.svelte";

export function handleNameKeydown(event: KeyboardEvent) {
  if (event.key === "Enter") {
    (event.target as HTMLInputElement)?.blur();
  } else if (event.key === "Escape") {
    timer.editingName = false;
  }
}

export function handleKeydown(event: KeyboardEvent) {
  if (timer.showSettings) return;
  if (timer.editingName) return;
  switch (event.key) {
    case " ":
      if (event.altKey) {
        event.preventDefault();
        timer.toggleBreak();
        break;
      }
      event.preventDefault();
      if (timer.timerState === "idle") timer.handleStart();
      else if (timer.timerState === "running") timer.handlePause();
      else if (timer.timerState === "paused") timer.handleResume();
      else if (timer.timerState === "finished") timer.handleReset();
      break;
    case "ArrowDown":
      if (timer.timerState !== "idle") {
        event.preventDefault();
        timer.handleReset();
      }
      break;
    case "+":
    case "=":
      if (timer.timerState === "running" || timer.timerState === "paused") {
        event.preventDefault();
        timer.handleAddTime();
      }
      break;
    case "ArrowRight":
      if (event.altKey) {
        event.preventDefault();
        if (event.ctrlKey) {
          timer.addBreakTimeTertiary();
        } else if (event.shiftKey) {
          timer.addBreakTimeSecondary();
        } else {
          timer.addBreakTime();
        }
        break;
      }
      if (timer.timerState === "running" || timer.timerState === "paused") {
        event.preventDefault();
        if (event.ctrlKey) {
          timer.handleAddTimeTertiary();
        } else if (event.shiftKey) {
          timer.handleAddTimeSecondary();
        } else {
          timer.handleAddTime();
        }
      }
      break;
    case "-":
      if (timer.timerState === "running" || timer.timerState === "paused") {
        event.preventDefault();
        timer.handleRemoveTime();
      }
      break;
    case "ArrowLeft":
      if (event.altKey) {
        event.preventDefault();
        if (event.ctrlKey) {
          timer.removeBreakTimeTertiary();
        } else if (event.shiftKey) {
          timer.removeBreakTimeSecondary();
        } else {
          timer.removeBreakTime();
        }
        break;
      }
      if (timer.timerState === "running" || timer.timerState === "paused") {
        event.preventDefault();
        if (event.ctrlKey) {
          timer.handleRemoveTimeTertiary();
        } else if (event.shiftKey) {
          timer.handleRemoveTimeSecondary();
        } else {
          timer.handleRemoveTime();
        }
      }
      break;
    case "Escape":
      if (timer.timerState === "finished") {
        event.preventDefault();
        timer.handleReset();
      }
      break;
    case "t":
    case "T":
    case "p":
    case "P":
      event.preventDefault();
      timer.toggleAlwaysOnTop();
      break;
    case "e":
    case "E":
      if (timer.hasMultipleTimers) {
        event.preventDefault();
        timer.handleSwitchNext();
      }
      break;
    case "q":
    case "Q":
      if (timer.hasMultipleTimers) {
        event.preventDefault();
        timer.handleSwitchPrev();
      }
      break;
    case "m":
    case "M":
      event.preventDefault();
      timer.toggleMiniMode();
      break;
    case "b":
    case "B":
      event.preventDefault();
      timer.toggleBreakDesignation();
      break;
  }
}
