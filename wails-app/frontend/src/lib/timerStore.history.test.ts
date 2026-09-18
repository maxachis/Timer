import { describe, it, expect, vi, beforeEach } from "vitest";

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

describe("history", () => {
  beforeEach(() => {
    invokeMock.mockReset();
    timer.showHistory = false;
    timer.historyEntries = [];
  });

  it("openHistory fetches history, stores entries newest-first, and opens the view", async () => {
    const oldestFirst = [
      { at: "2026-09-18T09:00:00-04:00", kind: "created", timer_name: "Work" },
      { at: "2026-09-18T09:05:00-04:00", kind: "started", timer_name: "Work" },
      { at: "2026-09-18T09:10:00-04:00", kind: "finished", timer_name: "Work" },
    ];
    invokeMock.mockImplementation((cmd: string) => {
      if (cmd === "get_history") return Promise.resolve(oldestFirst);
      return Promise.resolve(undefined);
    });

    await timer.openHistory();

    expect(invokeMock).toHaveBeenCalledWith("get_history");
    expect(timer.historyEntries.map((e) => e.kind)).toEqual([
      "finished",
      "started",
      "created",
    ]);
    expect(timer.showHistory).toBe(true);
  });

  it("closeHistory clears the open flag", () => {
    timer.showHistory = true;

    timer.closeHistory();

    expect(timer.showHistory).toBe(false);
  });
});
