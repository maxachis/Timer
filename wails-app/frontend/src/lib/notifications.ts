import { isPermissionGranted, requestPermission, sendNotification } from "../ipc";

export async function sendCompletionNotification(activeName: string) {
  let permitted = await isPermissionGranted();
  if (!permitted) {
    const permission = await requestPermission();
    permitted = permission === "granted";
  }
  if (permitted) {
    sendNotification({
      title: "Timer Complete",
      body: activeName ? `"${activeName}" has finished.` : "Your timer has finished.",
    });
  }
}
