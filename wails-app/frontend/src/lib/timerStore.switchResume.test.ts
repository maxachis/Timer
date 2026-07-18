import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";

// The store imports the Wails bridge, audio, and notifications. None of those
// should load in a unit test, so replace them with controllable fakes.
const invokeMock = vi.hoisted(() => vi.fn());

vi.mock("../ipc", () => ({
  invoke: invokeMock,
  windowApi: {
    setAlwaysOnTop: vi.fn(),
    setSize: vi.fn(),
    getSize: vi.fn(async () => ({ w: 300, h: 300 })),
    requestUserAttention: vi.fn(),
  },
}));
vi.mock("./audio", () => ({
  playAlertSound: vi.fn(),
  playStartBlip: vi.fn(),
  playStopBlip: vi.fn(),
}));
vi.mock("./notifications", () => ({ sendCompletionNotification: vi.fn() }));

import { timer } from "./timerStore.svelte";

// A status snapshot for the landed timer, reflecting whatever `landedState`
// currently is. active_index/name are fixed stand-ins — the tests only care
// about state transitions, not which timer is displayed.
function landedStatus(state: string) {
  return {
    remaining_secs: 300,
    state,
    is_finished: false,
    active_index: 1,
    active_name: "Timer 2",
    timer_count: 3,
    on_break: false,
    break_index: -1,
  };
}

describe("switch auto-resume", () => {
  let landedState: string;

  beforeEach(() => {
    vi.useFakeTimers();
    invokeMock.mockReset();
    landedState = "paused";
    invokeMock.mockImplementation((cmd: string) => {
      if (cmd === "get_timer_status") return Promise.resolve(landedStatus(landedState));
      if (cmd === "get_timer_list") return Promise.resolve([]);
      if (cmd === "resume_timer" || cmd === "start_timer") {
        landedState = "running";
        return Promise.resolve(undefined);
      }
      return Promise.resolve(undefined);
    });
    if (timer.pollInterval) {
      clearInterval(timer.pollInterval);
      timer.pollInterval = null;
    }
    timer.cancelPendingResume();
  });

  afterEach(() => {
    if (timer.pollInterval) {
      clearInterval(timer.pollInterval);
      timer.pollInterval = null;
    }
    timer.cancelPendingResume();
    vi.useRealTimers();
  });

  it("schedules and fires resume when leaving a running timer", async () => {
    timer.timerState = "running";
    landedState = "paused";

    await timer.handleSwitchTimer(1);

    expect(timer.pendingResumeTimeout).not.toBeNull();
    expect(invokeMock).not.toHaveBeenCalledWith("resume_timer");

    await vi.advanceTimersByTimeAsync(1600);

    expect(invokeMock).toHaveBeenCalledWith("resume_timer");
    expect(timer.pollInterval).not.toBeNull();
  });

  it("does NOT schedule resume when leaving a non-running timer", async () => {
    timer.timerState = "paused";
    landedState = "paused";

    await timer.handleSwitchTimer(1);

    expect(timer.pendingResumeTimeout).toBeNull();

    await vi.advanceTimersByTimeAsync(1600);

    expect(invokeMock).not.toHaveBeenCalledWith("resume_timer");
  });

  it("rotation coalesces to a single resume of the final timer", async () => {
    timer.timerState = "running";
    landedState = "paused";

    await timer.handleSwitchTimer(1);
    await vi.advanceTimersByTimeAsync(500);

    await timer.handleSwitchTimer(2);
    await vi.advanceTimersByTimeAsync(1600);

    const resumeCalls = invokeMock.mock.calls.filter(([cmd]) => cmd === "resume_timer");
    expect(resumeCalls.length).toBe(1);
  });
});
