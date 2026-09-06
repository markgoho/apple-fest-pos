// ADR-0007: the wipe is a genuine hard delete with no undo, guarded by one
// confirm() dialog, proportionate for a PIN-gated, solo-admin action.
"use strict";

const resetForm = document.getElementById("reset-form");
if (resetForm) {
  resetForm.addEventListener("submit", (event) => {
    if (!window.confirm("Wipe all orders? This cannot be undone.")) {
      event.preventDefault();
    }
  });
}

// ADR-0010: a discovered address fills the matching field in the printer
// assignment form. The admin still presses Save; nothing here writes state.
const foundPrinters = document.getElementById("found-printers");
if (foundPrinters) {
  foundPrinters.addEventListener("click", (event) => {
    const button = event.target.closest("button[data-host]");
    if (!button) {
      return;
    }
    const fieldID = button.classList.contains("use-as-window") ? "window_host" : "kitchen_host";
    document.getElementById(fieldID).value = button.dataset.host;
  });
}
