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
