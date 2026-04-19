import * as App from "../wailsjs/go/main/App";
import {
  Environment,
  Quit,
  WindowGetSize,
  WindowMinimise,
  WindowSetAlwaysOnTop,
  WindowSetSize,
  WindowSetMinSize,
  WindowGetPosition,
  WindowSetPosition,
} from "../wailsjs/runtime/runtime";

let platformPromise: Promise<string> | null = null;
function getPlatform(): Promise<string> {
  if (!platformPromise) platformPromise = Environment().then((e) => e.platform);
  return platformPromise;
}

type Args = Record<string, unknown> | undefined;

// Maps Tauri-style snake_case command names to Wails-generated PascalCase
// methods on the Go App struct. Arg shapes mirror the original Tauri calls.
export async function invoke<T = unknown>(cmd: string, args?: Args): Promise<T> {
  const a = args ?? {};
  switch (cmd) {
    case "get_timer_status":
      return (await App.GetTimerStatus()) as T;
    case "get_timer_list":
      return (await App.GetTimerList()) as T;
    case "start_timer":
      return (await App.StartTimer()) as T;
    case "pause_timer":
      return (await App.PauseTimer()) as T;
    case "resume_timer":
      return (await App.ResumeTimer()) as T;
    case "reset_timer":
      return (await App.ResetTimer()) as T;
    case "add_time":
      return (await App.AddTime()) as T;
    case "remove_time":
      return (await App.RemoveTime()) as T;
    case "add_time_secondary":
      return (await App.AddTimeSecondary()) as T;
    case "remove_time_secondary":
      return (await App.RemoveTimeSecondary()) as T;
    case "add_time_tertiary":
      return (await App.AddTimeTertiary()) as T;
    case "remove_time_tertiary":
      return (await App.RemoveTimeTertiary()) as T;
    case "add_time_custom":
      return (await App.AddTimeCustom((a as any).seconds)) as T;
    case "remove_time_custom":
      return (await App.RemoveTimeCustom((a as any).seconds)) as T;
    case "create_timer":
      return (await App.CreateTimer((a as any).durationSecs)) as T;
    case "add_new_timer":
      return (await App.AddNewTimer((a as any).name)) as T;
    case "remove_existing_timer":
      return (await App.RemoveExistingTimer((a as any).index)) as T;
    case "switch_timer":
      return (await App.SwitchTimer((a as any).index)) as T;
    case "switch_timer_next":
      return (await App.SwitchTimerNext()) as T;
    case "switch_timer_prev":
      return (await App.SwitchTimerPrev()) as T;
    case "rename_timer":
      return (await App.RenameTimer((a as any).index, (a as any).name)) as T;
    case "get_settings":
      return (await App.GetSettings()) as T;
    case "save_settings":
      return (await App.SaveSettings((a as any).newSettings)) as T;
  }
  throw new Error(`Unknown IPC command: ${cmd}`);
}

export const windowApi = {
  setAlwaysOnTop: (on: boolean) => WindowSetAlwaysOnTop(on),
  setSize: async (w: number, h: number) => {
    // Update the min track size so Windows' WM_GETMINMAXINFO returns the new
    // floor; otherwise starting a drag snaps the window back to the startup
    // minimums. Also lets widths/heights below the OS default titled-window
    // minimum (~132px) take effect.
    WindowSetMinSize(w, h);
    // On Windows, WindowSetSize preserves position, and GetPosition/SetPosition
    // have a DPI-scaling mismatch on HiDPI displays that can push the window
    // off-screen. Only save/restore on non-Windows platforms (Linux needs it
    // because WindowSetSize snaps the window to the top of the screen).
    const platform = await getPlatform();
    if (platform === "windows") {
      WindowSetSize(w, h);
      return;
    }
    const pos = await WindowGetPosition();
    WindowSetSize(w, h);
    WindowSetPosition(pos.x, pos.y);
  },
  requestUserAttention: (_: unknown) => {
    /* no Wails equivalent on Linux; noop */
  },
  getSize: () => WindowGetSize(),
  minimize: () => WindowMinimise(),
  close: () => Quit(),
};

// Browser Notification API wrapper matching Tauri plugin-notification shape.
export async function isPermissionGranted(): Promise<boolean> {
  if (!("Notification" in window)) return false;
  return Notification.permission === "granted";
}

export async function requestPermission(): Promise<"granted" | "denied" | "default"> {
  if (!("Notification" in window)) return "denied";
  return await Notification.requestPermission();
}

export function sendNotification(opts: { title: string; body?: string }): void {
  if (!("Notification" in window)) return;
  if (Notification.permission !== "granted") return;
  new Notification(opts.title, { body: opts.body });
}
