// Kiosk lock-down for the Operator's tablets. Ticket #6 proved these browser
// APIs hold on a Pixel tablet with no device setting changed, but the wake
// lock and fullscreen navigationUI hide both need HTTPS to work at all.
"use strict";

const resumeButton = document.getElementById("kiosk-resume");

let wakeLock = null;

// 100dvh measures larger than the real visible screen when installed and
// running fullscreen on Android Chrome, which pushed fixed bottom elements
// like the checkout bar below the fold. window.innerHeight is the real
// number; every full-page layout keys off --app-height instead of dvh.
function setAppHeight() {
  document.documentElement.style.setProperty("--app-height", `${window.innerHeight}px`);
}
setAppHeight();
window.addEventListener("resize", setAppHeight);

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

resumeButton.addEventListener("click", () => {
  enterFullscreen();
});

// Only present on the mode-chooser screen, not on the transaction screens,
// so a stray tap mid-shift cannot force a reload or drop full screen.
const refreshButton = document.getElementById("kiosk-refresh");

// A reload on the Pi takes a beat, and until it lands the screen is unchanged.
// Without this the Operator cannot tell a slow reload from a tap that missed,
// so they tap again. The button says what it is doing and stops taking taps.
function markBusy(button, label) {
  button.disabled = true;
  button.textContent = label;
  button.classList.add("is-busy");
}

refreshButton?.addEventListener("click", () => {
  markBusy(refreshButton, "Refreshing…");
  location.reload();
});

// Cashier POS is a plain link to a full page load, which
// can take a few seconds on a slow Pi or thin Wi-Fi. A link gives no feedback
// of its own between the tap and the new page painting, so a slow load reads
// exactly like a tap that missed and the Operator taps again. The link still
// navigates itself; this only holds the pressed look and, if the page never
// actually leaves, says so instead of leaving the tablet looking stuck.
const modeFault = document.getElementById("mode-fault");
document.querySelectorAll(".modes a").forEach((link) => {
  link.addEventListener("click", () => {
    if (modeFault) modeFault.textContent = "";
    link.classList.add("is-busy");
    const label = link.textContent;
    link.textContent = "Opening…";
    setTimeout(() => {
      link.classList.remove("is-busy");
      link.textContent = label;
      if (modeFault) modeFault.textContent = "Could not open " + label + ". Check the connection and try again.";
    }, 8000);
  });
});

document.addEventListener("fullscreenchange", () => {
  resumeButton.hidden = isImmersive();
});

document.addEventListener("visibilitychange", () => {
  if (document.visibilityState === "visible") {
    takeWakeLock();
  }
});

document.addEventListener("contextmenu", (event) => {
  if (event.target.closest("input, textarea")) return;
  event.preventDefault();
});

// The wake lock needs no tap, so every page load takes it. Full screen does
// need one: the installed app gets it from the manifest, and a browser tab
// gets the Resume bar.
takeWakeLock();
resumeButton.hidden = isImmersive();
