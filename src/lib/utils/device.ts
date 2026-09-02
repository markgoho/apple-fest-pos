import { createId } from "./id";

const deviceStorageKey = "apple-fest-pos-device-id";

export function getDeviceId(): string {
  let deviceId = localStorage.getItem(deviceStorageKey);
  if (!deviceId) {
    deviceId = createId();
    localStorage.setItem(deviceStorageKey, deviceId);
  }

  return deviceId;
}
