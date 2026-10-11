// Pure helpers for the timer history view. No Svelte imports — the store and
// the History component both read from here so the kind list and labels stay
// in one place.

export interface HistoryEntry {
  at: string; // RFC3339 timestamp with offset
  kind: HistoryKind;
  timer_name: string;
  remaining_ms?: number; // clock reading, on paused/resumed entries only
}

export const HISTORY_KINDS = [
  "created",
  "started",
  "paused",
  "resumed",
  "reset",
  "finished",
  "removed",
] as const;

export type HistoryKind = (typeof HISTORY_KINDS)[number];

export const KIND_LABELS: Record<HistoryKind, string> = {
  created: "Created",
  started: "Started",
  paused: "Paused",
  resumed: "Resumed",
  reset: "Reset",
  finished: "Finished",
  removed: "Removed",
};

export function formatEntryTime(at: string): string {
  return new Date(at).toLocaleTimeString(undefined, {
    hour: "numeric",
    minute: "2-digit",
    second: "2-digit",
  });
}

const MS_PER_SECOND = 1000;
const SECONDS_PER_MINUTE = 60;

// Formats a clock reading the way the timer face shows it: whole minutes
// (unbounded) and seconds, partial seconds dropped.
export function formatClockReading(ms: number): string {
  const totalSecs = Math.floor(ms / MS_PER_SECOND);
  const minutes = Math.floor(totalSecs / SECONDS_PER_MINUTE);
  const seconds = totalSecs % SECONDS_PER_MINUTE;
  return `${String(minutes).padStart(2, "0")}:${String(seconds).padStart(2, "0")}`;
}

const TODAY_LABEL = "Today";
const YESTERDAY_LABEL = "Yesterday";

function isSameCalendarDay(a: Date, b: Date): boolean {
  return (
    a.getFullYear() === b.getFullYear() &&
    a.getMonth() === b.getMonth() &&
    a.getDate() === b.getDate()
  );
}

function dayLabel(entryDate: Date, now: Date): string {
  if (isSameCalendarDay(entryDate, now)) return TODAY_LABEL;
  const yesterday = new Date(now);
  yesterday.setDate(now.getDate() - 1);
  if (isSameCalendarDay(entryDate, yesterday)) return YESTERDAY_LABEL;
  return entryDate.toLocaleDateString(undefined, {
    weekday: "long",
    month: "long",
    day: "numeric",
  });
}

// Takes entries NEWEST FIRST and groups consecutive same-calendar-day runs,
// preserving order within and across groups.
export function groupByDay(
  entries: HistoryEntry[],
  now: Date = new Date(),
): { label: string; entries: HistoryEntry[] }[] {
  const groups: { label: string; entries: HistoryEntry[] }[] = [];
  for (const entry of entries) {
    const entryDate = new Date(entry.at);
    const label = dayLabel(entryDate, now);
    const lastGroup = groups[groups.length - 1];
    if (lastGroup && lastGroup.label === label) {
      lastGroup.entries.push(entry);
    } else {
      groups.push({ label, entries: [entry] });
    }
  }
  return groups;
}
