import { describe, it, expect } from "vitest";
import {
  groupByDay,
  formatEntryTime,
  KIND_LABELS,
  HISTORY_KINDS,
  type HistoryEntry,
} from "./history";

function entry(at: string, kind: HistoryEntry["kind"] = "started"): HistoryEntry {
  return { at, kind, timer_name: "Work" };
}

describe("groupByDay", () => {
  it("labels today, yesterday, and older days, preserving order", () => {
    const now = new Date("2026-09-18T20:00:00-04:00");
    const entries: HistoryEntry[] = [
      entry("2026-09-18T14:03:05.123-04:00"), // today
      entry("2026-09-18T09:00:00.000-04:00"), // today
      entry("2026-09-17T22:00:00.000-04:00"), // yesterday
      entry("2026-09-15T10:00:00.000-04:00"), // older
    ];

    const groups = groupByDay(entries, now);

    expect(groups.map((g) => g.label)).toEqual([
      "Today",
      "Yesterday",
      new Date("2026-09-15T10:00:00.000-04:00").toLocaleDateString(undefined, {
        weekday: "long",
        month: "long",
        day: "numeric",
      }),
    ]);
    expect(groups[0].entries).toHaveLength(2);
    expect(groups[0].entries[0].at).toBe(entries[0].at);
    expect(groups[0].entries[1].at).toBe(entries[1].at);
    expect(groups[1].entries).toHaveLength(1);
    expect(groups[2].entries).toHaveLength(1);
  });

  it("keeps non-consecutive same-day runs as separate groups", () => {
    const now = new Date("2026-09-18T20:00:00-04:00");
    const entries: HistoryEntry[] = [
      entry("2026-09-18T14:00:00.000-04:00"), // today
      entry("2026-09-17T10:00:00.000-04:00"), // yesterday
      entry("2026-09-18T09:00:00.000-04:00"), // today again, non-consecutive
    ];

    const groups = groupByDay(entries, now);

    expect(groups.map((g) => g.label)).toEqual(["Today", "Yesterday", "Today"]);
  });
});

describe("formatEntryTime", () => {
  it("returns a non-empty string containing the seconds for a fixed timestamp", () => {
    const formatted = formatEntryTime("2026-09-18T14:03:05.123-04:00");
    expect(formatted.length).toBeGreaterThan(0);
    expect(formatted).toContain("05");
  });
});

describe("KIND_LABELS", () => {
  it("covers every HISTORY_KINDS entry", () => {
    for (const kind of HISTORY_KINDS) {
      expect(KIND_LABELS[kind]).toBeTruthy();
    }
    expect(Object.keys(KIND_LABELS).sort()).toEqual([...HISTORY_KINDS].sort());
  });
});
