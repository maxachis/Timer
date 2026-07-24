import { describe, it, expect } from "vitest";
import { barColor, tintSwatch, NEUTRAL_BAR_COLOR, type TimerListEntry } from "./tints";

function entry(overrides: Partial<TimerListEntry> = {}): TimerListEntry {
  return {
    index: 0,
    name: "Work",
    state: "idle",
    remaining_secs: 300,
    is_active: false,
    tint: "",
    original_duration_secs: 300,
    is_break_timer: false,
    break_accrued_secs: 0,
    ...overrides,
  };
}

describe("barColor", () => {
  it("returns the tint's swatch color when a palette tint is set", () => {
    expect(barColor(entry({ tint: "sky" }))).toBe(tintSwatch("sky"));
  });

  it("returns the neutral color for an unset tint", () => {
    expect(barColor(entry({ tint: "" }))).toBe(NEUTRAL_BAR_COLOR);
  });

  it("returns the neutral color for an explicit 'none' tint", () => {
    expect(barColor(entry({ tint: "none" }))).toBe(NEUTRAL_BAR_COLOR);
  });
});
