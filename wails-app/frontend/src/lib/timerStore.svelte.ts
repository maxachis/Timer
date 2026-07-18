import { invoke, windowApi } from "../ipc";
import { playAlertSound, playStartBlip, playStopBlip } from "./audio";
import { sendCompletionNotification } from "./notifications";
import { TINTS, type TimerListEntry, type TimerStatus } from "./tints";

const getCurrentWindow = () => windowApi;

// SVG circle math
export const RADIUS = 130;
export const CIRCUMFERENCE = 2 * Math.PI * RADIUS;

export function formatMinutes(secs: number): string {
  return String(Math.floor(secs / 60)).padStart(2, "0");
}

export function formatSeconds(secs: number): string {
  return String(Math.floor(secs % 60)).padStart(2, "0");
}

class TimerStore {
  remainingSecs = $state(300);
  timerState = $state("idle");
  isFinished = $state(false);
  pollInterval: ReturnType<typeof setInterval> | null = $state(null);
  totalDuration = $state(300);
  prevState = $state("idle");
  alwaysOnTop = $state(false);
  showSettings = $state(false);
  settingsDurationMin = $state(5);
  settingsIncrementMin = $state(5);
  settingsSecondaryIncrementMin = $state(1);
  settingsTertiaryIncrementMin = $state(60);
  settingsMiniWindowWidth = $state(70);
  settingsMiniWindowHeight = $state(70);
  incrementSecs = $state(300);
  secondaryIncrementSecs = $state(60);
  tertiaryIncrementSecs = $state(3600);
  miniWindowWidth = $state(70);
  miniWindowHeight = $state(70);

  // Multi-timer state
  activeIndex = $state(0);
  activeName = $state("");
  timerCount = $state(1);
  timerList: TimerListEntry[] = $state([]);
  toastMessage = $state("");
  toastTimeout: ReturnType<typeof setTimeout> | null = $state(null);
  editingName = $state(false);
  editNameValue = $state("");
  showTintMenu = $state(false);

  // Break mode
  onBreak = $state(false);
  breakIndex = $state(-1);

  preMiniSize: { w: number; h: number } | null = $state(null);

  get hasMultipleTimers() {
    return this.timerCount > 1;
  }

  get activeTintId() {
    return this.timerList.find((t) => t.is_active)?.tint || "none";
  }

  get activeTintWash() {
    return TINTS.find((t) => t.id === this.activeTintId)?.wash ?? "transparent";
  }

  get incrementLabel() {
    return Math.round(this.incrementSecs / 60) + " min";
  }

  get secondaryIncrementLabel() {
    return Math.round(this.secondaryIncrementSecs / 60) + " min";
  }

  get tertiaryIncrementLabel() {
    return this.tertiaryIncrementSecs >= 3600
      ? Math.round(this.tertiaryIncrementSecs / 3600) + " hr"
      : Math.round(this.tertiaryIncrementSecs / 60) + " min";
  }

  // Track total duration for the progress ring
  get progress() {
    return this.totalDuration > 0 ? this.remainingSecs / this.totalDuration : 0;
  }

  get strokeOffset() {
    return CIRCUMFERENCE * (1 - this.progress);
  }

  showToast(msg: string) {
    this.toastMessage = msg;
    if (this.toastTimeout) clearTimeout(this.toastTimeout);
    this.toastTimeout = setTimeout(() => {
      this.toastMessage = "";
      this.toastTimeout = null;
    }, 3000);
  }

  async loadSettings() {
    const s = await invoke<{
      default_duration_secs: number;
      default_increment_secs: number;
      secondary_increment_secs: number;
      tertiary_increment_secs: number;
      mini_window_width: number;
      mini_window_height: number;
    }>("get_settings");
    this.settingsDurationMin = Math.round(s.default_duration_secs / 60);
    this.settingsIncrementMin = Math.round(s.default_increment_secs / 60);
    this.settingsSecondaryIncrementMin = Math.round(s.secondary_increment_secs / 60);
    this.settingsTertiaryIncrementMin = Math.round(s.tertiary_increment_secs / 60);
    this.settingsMiniWindowWidth = s.mini_window_width;
    this.settingsMiniWindowHeight = s.mini_window_height;
    this.incrementSecs = s.default_increment_secs;
    this.secondaryIncrementSecs = s.secondary_increment_secs;
    this.tertiaryIncrementSecs = s.tertiary_increment_secs;
    this.miniWindowWidth = s.mini_window_width;
    this.miniWindowHeight = s.mini_window_height;
  }

  async openSettings() {
    await this.loadSettings();
    this.showSettings = true;
  }

  closeSettings() {
    this.showSettings = false;
  }

  async saveSettings() {
    const durationMin = Math.max(1, Math.min(180, this.settingsDurationMin));
    const incrementMin = Math.max(1, Math.min(60, this.settingsIncrementMin));
    const secondaryIncrementMin = Math.max(1, Math.min(60, this.settingsSecondaryIncrementMin));
    const tertiaryIncrementMin = Math.max(1, Math.min(1440, this.settingsTertiaryIncrementMin));
    const miniWidth = Math.max(50, Math.min(400, this.settingsMiniWindowWidth));
    const miniHeight = Math.max(50, Math.min(400, this.settingsMiniWindowHeight));
    const newSettings = {
      default_duration_secs: durationMin * 60,
      default_increment_secs: incrementMin * 60,
      secondary_increment_secs: secondaryIncrementMin * 60,
      tertiary_increment_secs: tertiaryIncrementMin * 60,
      mini_window_width: miniWidth,
      mini_window_height: miniHeight,
    };
    await invoke("save_settings", { newSettings });
    this.incrementSecs = newSettings.default_increment_secs;
    this.secondaryIncrementSecs = newSettings.secondary_increment_secs;
    this.tertiaryIncrementSecs = newSettings.tertiary_increment_secs;
    this.miniWindowWidth = newSettings.mini_window_width;
    this.miniWindowHeight = newSettings.mini_window_height;
    this.showSettings = false;
    if (this.timerState === "idle") {
      await invoke("create_timer", { durationSecs: newSettings.default_duration_secs });
      await this.fetchStatus();
    }
  }

  async toggleAlwaysOnTop() {
    this.alwaysOnTop = !this.alwaysOnTop;
    await getCurrentWindow().setAlwaysOnTop(this.alwaysOnTop);
  }

  async toggleMiniMode() {
    if (this.preMiniSize) {
      const { w, h } = this.preMiniSize;
      this.preMiniSize = null;
      await getCurrentWindow().setSize(w, h);
      return;
    }
    const current = await getCurrentWindow().getSize();
    this.preMiniSize = { w: current.w, h: current.h };
    await getCurrentWindow().setSize(this.miniWindowWidth, this.miniWindowHeight);
  }

  async fetchStatus() {
    const status = await invoke<TimerStatus>("get_timer_status");
    this.remainingSecs = status.remaining_secs;
    this.timerState = status.state;
    this.isFinished = status.is_finished;
    this.activeIndex = status.active_index;
    this.activeName = status.active_name;
    this.timerCount = status.timer_count;
    this.onBreak = status.on_break;
    this.breakIndex = status.break_index;

    // Detect transition to finished state
    if (this.prevState !== "finished" && status.state === "finished") {
      playAlertSound();
      sendCompletionNotification(this.activeName);
      getCurrentWindow().requestUserAttention("critical");
    }
    this.prevState = status.state;

    // Auto-end break mode when the break timer hits 0. toggleBreak() ends the
    // break server-side and re-fetches fresh status (restoring/resuming the work
    // timer), so return here rather than let this stale snapshot's poll-clear
    // logic below wrongly stop polling a just-resumed running work timer.
    if (this.onBreak && status.is_finished) {
      await this.toggleBreak();
      return;
    }

    // Update total duration when idle (after reset or create)
    if (status.state === "idle") {
      this.totalDuration = status.remaining_secs;
    }

    if (status.state !== "running" && this.pollInterval) {
      clearInterval(this.pollInterval);
      this.pollInterval = null;
    }
  }

  async fetchTimerList() {
    this.timerList = await invoke<TimerListEntry[]>("get_timer_list");
  }

  startPolling() {
    if (!this.pollInterval) {
      this.pollInterval = setInterval(() => {
        this.fetchStatus();
        this.fetchTimerList();
      }, 80);
    }
  }

  async handleStart() {
    playStartBlip();
    await invoke("start_timer");
    await this.fetchStatus();
    await this.fetchTimerList();
    this.startPolling();
  }

  async handlePause() {
    playStopBlip();
    await invoke("pause_timer");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async handleResume() {
    playStartBlip();
    await invoke("resume_timer");
    await this.fetchStatus();
    await this.fetchTimerList();
    this.startPolling();
  }

  async handleReset() {
    playStopBlip();
    await invoke("reset_timer");
    getCurrentWindow().requestUserAttention(null);
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async handleAddTime() {
    await invoke("add_time");
    this.totalDuration += this.incrementSecs;
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async handleRemoveTime() {
    await invoke("remove_time");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async handleAddTimeSecondary() {
    await invoke("add_time_secondary");
    this.totalDuration += this.secondaryIncrementSecs;
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async handleRemoveTimeSecondary() {
    await invoke("remove_time_secondary");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async handleAddTimeTertiary() {
    await invoke("add_time_tertiary");
    this.totalDuration += this.tertiaryIncrementSecs;
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async handleRemoveTimeTertiary() {
    await invoke("remove_time_tertiary");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  // Multi-timer actions
  async handleAddTimer() {
    try {
      const defaultName = `Timer ${this.timerCount + 1}`;
      await invoke("add_new_timer", { name: defaultName });
      await this.fetchStatus();
      await this.fetchTimerList();
    } catch (e) {
      this.showToast("Timer limit reached — remove one first");
    }
  }

  async handleRemoveTimer(index: number) {
    try {
      await invoke("remove_existing_timer", { index });
      await this.fetchStatus();
      await this.fetchTimerList();
    } catch (_e) {
      // Can't remove last timer
    }
  }

  async handleSetTint(tintId: string) {
    this.showTintMenu = false;
    await invoke("set_timer_tint", { index: this.activeIndex, tint: tintId });
    await this.fetchTimerList();
  }

  async handleSwitchTimer(index: number) {
    if (this.timerState === "running") playStopBlip();
    if (this.pollInterval) {
      clearInterval(this.pollInterval);
      this.pollInterval = null;
    }
    await invoke("switch_timer", { index });
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async handleSwitchNext() {
    if (this.timerCount <= 1) return;
    if (this.timerState === "running") playStopBlip();
    if (this.pollInterval) {
      clearInterval(this.pollInterval);
      this.pollInterval = null;
    }
    await invoke("switch_timer_next");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async handleSwitchPrev() {
    if (this.timerCount <= 1) return;
    if (this.timerState === "running") playStopBlip();
    if (this.pollInterval) {
      clearInterval(this.pollInterval);
      this.pollInterval = null;
    }
    await invoke("switch_timer_prev");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  startEditName() {
    this.editNameValue = this.activeName;
    this.editingName = true;
  }

  async finishEditName() {
    this.editingName = false;
    const trimmed = this.editNameValue.trim();
    if (trimmed !== this.activeName) {
      await invoke("rename_timer", { index: this.activeIndex, name: trimmed });
      await this.fetchStatus();
      await this.fetchTimerList();
    }
  }

  // Break mode actions
  async toggleBreakDesignation() {
    await invoke("toggle_break_timer");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async toggleBreak() {
    await invoke("toggle_break");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async setMaxTime() {
    await invoke("set_max_time");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async addBreakTime() {
    await invoke("add_break_time");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async removeBreakTime() {
    await invoke("remove_break_time");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async addBreakTimeSecondary() {
    await invoke("add_break_time_secondary");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async removeBreakTimeSecondary() {
    await invoke("remove_break_time_secondary");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async addBreakTimeTertiary() {
    await invoke("add_break_time_tertiary");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  async removeBreakTimeTertiary() {
    await invoke("remove_break_time_tertiary");
    await this.fetchStatus();
    await this.fetchTimerList();
  }

  handleMinimize() {
    getCurrentWindow().minimize();
  }

  handleClose() {
    getCurrentWindow().close();
  }
}

export const timer = new TimerStore();
