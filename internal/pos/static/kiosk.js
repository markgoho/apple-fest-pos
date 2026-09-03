// Kiosk lock-down for the Operator's tablets. Ticket #6 proved these browser
// APIs hold on a Pixel tablet with no device setting changed, but the wake
// lock and fullscreen navigationUI hide both need HTTPS to work at all.
"use strict";

const SHIFT_KEY = "apple-fest-pos-shift-started";
const startButton = document.getElementById("kiosk-start");
const resumeButton = document.getElementById("kiosk-resume");

let wakeLock = null;

function shiftStarted() {
  try {
    return localStorage.getItem(SHIFT_KEY) === "1";
  } catch {
    return false;
  }
}

function markShiftStarted() {
  try {
    localStorage.setItem(SHIFT_KEY, "1");
  } catch {
    // Worst case the start screen shows again on the next load.
  }
}

async function takeWakeLock() {
  if (wakeLock || !("wakeLock" in navigator)) return;
  try {
    wakeLock = await navigator.wakeLock.request("screen");
    wakeLock.addEventListener("release", () => {
      wakeLock = null;
    });
  } catch {
    // Refused off-tablet or off-HTTPS; the booth tablets are the target.
  }
}

async function enterFullscreen() {
  try {
    await document.documentElement.requestFullscreen({ navigationUI: "hide" });
  } catch {
    // Fullscreen needs a user gesture; every caller here is a tap.
  }
}

// An installed app launches into fullscreen chrome straight from the
// manifest, with no call to the Fullscreen API and so no
// document.fullscreenElement. Miss that, and the Resume bar wrongly shows
// on a screen that is already fullscreen.
function isImmersive() {
  return (
    document.fullscreenElement !== null ||
    window.matchMedia("(display-mode: fullscreen), (display-mode: standalone)").matches
  );
}

startButton.addEventListener("click", async () => {
  await enterFullscreen();
  await takeWakeLock();
  markShiftStarted();
  startButton.hidden = true;
});

resumeButton.addEventListener("click", () => {
  enterFullscreen();
});

document.addEventListener("fullscreenchange", () => {
  resumeButton.hidden = isImmersive();
});

document.addEventListener("visibilitychange", () => {
  if (document.visibilityState === "visible" && shiftStarted()) {
    takeWakeLock();
  }
});

document.addEventListener("contextmenu", (event) => {
  if (event.target.closest("input, textarea")) return;
  event.preventDefault();
});

if (shiftStarted()) {
  takeWakeLock();
  resumeButton.hidden = isImmersive();
} else {
  startButton.hidden = false;
}
