import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";

// The store imports the Wails bridge, audio, and notifications. None of those
// should load in a unit test, so replace them with controllable fakes. invoke
// is the only one the break flow exercises.
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

// A status snapshot for the active (break) timer in the given run state.
function breakStatus(state: string) {
  return {
    remaining_secs: 700,
    state,
    is_finished: false,
    active_index: 3,
    active_name: "Break",
    timer_count: 7,
    on_break: true,
    break_index: 3,
  };
}

function mockInvokeReturning(state: string) {
  invokeMock.mockImplementation((cmd: string) => {
    if (cmd === "get_timer_status") return Promise.resolve(breakStatus(state));
    if (cmd === "get_timer_list") return Promise.resolve([]);
    return Promise.resolve(undefined);
  });
}

describe("toggleBreak polling", () => {
  beforeEach(() => {
    invokeMock.mockReset();
    if (timer.pollInterval) {
      clearInterval(timer.pollInterval);
      timer.pollInterval = null;
    }
    // Simulate a work timer that was NOT running when the break began, so no
    // poll loop is already alive — the case where the bug surfaced.
    timer.timerState = "paused";
    timer.prevState = "paused";
  });

  afterEach(() => {
    if (timer.pollInterval) {
      clearInterval(timer.pollInterval);
      timer.pollInterval = null;
    }
  });

  it("starts the poll loop when a break begins from a non-running timer", async () => {
    mockInvokeReturning("running");

    await timer.toggleBreak();

    // The break timer is now running server-side; the UI must poll so the
    // displayed digits actually tick down.
    expect(timer.timerState).toBe("running");
    expect(timer.pollInterval).not.toBeNull();
  });

  it("does not start polling when a break ends on a paused timer", async () => {
    mockInvokeReturning("paused");

    await timer.toggleBreak();

    expect(timer.timerState).toBe("paused");
    expect(timer.pollInterval).toBeNull();
  });
});
