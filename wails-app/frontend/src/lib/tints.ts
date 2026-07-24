export interface TimerStatus {
  remaining_secs: number;
  state: string;
  is_finished: boolean;
  active_index: number;
  active_name: string;
  timer_count: number;
  on_break: boolean;
  break_index: number;
}

export interface TimerListEntry {
  index: number;
  name: string;
  state: string;
  remaining_secs: number;
  is_active: boolean;
  tint: string;
  original_duration_secs: number;
  is_break_timer: boolean;
  break_accrued_secs: number;
}

export const TINTS = [
  { id: "none",     label: "None",     swatch: "transparent", wash: "transparent" },
  { id: "rose",     label: "Rose",     swatch: "#d98a9c",     wash: "rgba(201,100,120,0.22)" },
  { id: "sage",     label: "Sage",     swatch: "#9cb890",     wash: "rgba(126,154,116,0.24)" },
  { id: "sky",      label: "Sky",      swatch: "#8db4dc",     wash: "rgba(96,146,194,0.22)" },
  { id: "lavender", label: "Lavender", swatch: "#b6a0dc",     wash: "rgba(150,122,196,0.22)" },
  { id: "amber",    label: "Amber",    swatch: "#e0a866",     wash: "rgba(212,140,60,0.24)" },
  { id: "teal",     label: "Teal",     swatch: "#8fbfb4",     wash: "rgba(110,170,158,0.22)" },
  { id: "clay",     label: "Clay",     swatch: "#d69a7c",     wash: "rgba(198,128,92,0.22)" },
];

export function tintWash(id: string): string {
  return TINTS.find((t) => t.id === id)?.wash ?? "transparent";
}

export function tintSwatch(id: string): string {
  return TINTS.find((t) => t.id === id)?.swatch ?? "transparent";
}

// NEUTRAL_BAR_COLOR is the stacked-bar lane color for a timer with no tint
// (unset or explicitly "none"), so an untinted timer never shows a phantom
// color borrowed from its position in the list.
export const NEUTRAL_BAR_COLOR = "#b8b0a4";

export function barColor(entry: TimerListEntry): string {
  if (entry.tint && entry.tint !== "none") {
    return tintSwatch(entry.tint);
  }
  return NEUTRAL_BAR_COLOR;
}
