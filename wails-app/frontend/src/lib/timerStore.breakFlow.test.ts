import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";

// Same fakes as the polling tests: keep the Wails bridge, audio, and
// notifications out of the unit test.
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
import { playAlertSound } from "./audio";

function workStatus() {
  return {
    remaining_secs: 900,
    state: "running",
    is_finished: false,
    active_index: 0,
    active_name: "Work",
    timer_count: 2,
    on_break: false,
    break_index: 1,
  };
}

// A break whose allocation has just run dry: still on break, the break block
// itself is still running with time left, but the countdown reports zero.
function exhaustedBreakStatus() {
  return {
    remaining_secs: 0,
    state: "running",
    is_finished: false,
    active_index: 1,
    active_name: "Break",
    timer_count: 2,
    on_break: true,
    break_index: 1,
  };
}

describe("break allocation flow", () => {
  beforeEach(() => {
    invokeMock.mockReset();
    vi.mocked(playAlertSound).mockReset();
    if (timer.pollInterval) {
      clearInterval(timer.pollInterval);
      timer.pollInterval = null;
    }
    timer.toastMessage = "";
    timer.onBreak = false;
    timer.timerState = "running";
    timer.prevState = "running";
  });

  afterEach(() => {
    if (timer.pollInterval) {
      clearInterval(timer.pollInterval);
      timer.pollInterval = null;
    }
  });

  it("ends the break and alerts when the allocation runs out", async () => {
    let broke = false;
    invokeMock.mockImplementation((cmd: string) => {
      if (cmd === "toggle_break") {
        broke = true;
        return Promise.resolve(undefined);
      }
      if (cmd === "get_timer_status") {
        return Promise.resolve(broke ? workStatus() : exhaustedBreakStatus());
      }
      if (cmd === "get_timer_list") return Promise.resolve([]);
      return Promise.resolve(undefined);
    });

    await timer.fetchStatus();

    expect(invokeMock).toHaveBeenCalledWith("toggle_break");
    expect(timer.onBreak).toBe(false);
    expect(playAlertSound).toHaveBeenCalled();
  });

  it("toasts instead of breaking when the backend refuses to start one", async () => {
    invokeMock.mockImplementation((cmd: string) => {
      if (cmd === "toggle_break") return Promise.reject("No break time accrued");
      if (cmd === "get_timer_status") return Promise.resolve(workStatus());
      if (cmd === "get_timer_list") return Promise.resolve([]);
      return Promise.resolve(undefined);
    });

    await timer.toggleBreak();

    expect(timer.onBreak).toBe(false);
    expect(timer.toastMessage).toBe("No break time accrued");
  });
});
