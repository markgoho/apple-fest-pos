// The lock screen of a PIN-gated page: an on-screen pad in place of the
// tablet's keyboard. The page works without this script (a typed input and an
// Unlock button); the script takes those over only after it has run, so a
// fault here leaves the typed input in place. Every PIN is as long as the
// input's maxlength, so the last digit sends the form and there is no Unlock.
"use strict";

const pinInput = document.getElementById("pin");
const pinPad = document.querySelector(".pin-pad");

if (pinInput && pinPad) {
  const form = pinInput.form;
  const pinLength = pinInput.maxLength;

  const press = (key) => {
    if (pinPad.disabled) {
      return;
    }
    if (key === "clear") {
      pinInput.value = "";
    } else if (key === "delete") {
      pinInput.value = pinInput.value.slice(0, -1);
    } else if (pinInput.value.length < pinLength) {
      pinInput.value += key;
    }
    if (pinInput.value.length === pinLength) {
      pinPad.disabled = true;
      form.submit();
    }
  };

  pinPad.addEventListener("click", (event) => {
    const button = event.target.closest("button");
    if (button) {
      press(button.dataset.pinKey || button.textContent);
    }
  });

  // A hardware keyboard still works, for the laptop.
  document.addEventListener("keydown", (event) => {
    if (/^[0-9]$/.test(event.key)) {
      press(event.key);
    } else if (event.key === "Backspace") {
      press("delete");
    }
  });

  // Back can restore this page as it was left: full and disabled.
  window.addEventListener("pageshow", () => {
    pinInput.value = "";
    pinPad.disabled = false;
  });

  pinInput.readOnly = true;
  pinInput.inputMode = "none";
  pinInput.blur();
  form.querySelector("button[type=submit]").hidden = true;
  pinPad.hidden = false;
}
