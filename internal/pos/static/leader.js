// Figures/Orders tabs are pure client state (the whole day is already in the
// rendered page); Void and Comp each need one confirm before they fire, per
// CONTEXT.md's Voided and Comped entries: the confirm text shows the order
// number so the Leader matches paper to screen (issue #45, ADR-0012).
"use strict";

const tabs = {
  figures: { button: document.getElementById("tab-figures"), panel: document.getElementById("figures-panel") },
  orders: { button: document.getElementById("tab-orders"), panel: document.getElementById("orders-panel") },
};

function selectTab(name) {
  for (const [tabName, tab] of Object.entries(tabs)) {
    tab.panel.hidden = tabName !== name;
    tab.button.setAttribute("aria-selected", String(tabName === name));
  }
}

for (const [name, tab] of Object.entries(tabs)) {
  tab.button.addEventListener("click", () => selectTab(name));
}

function confirmBeforeSubmit(selector, question) {
  for (const form of document.querySelectorAll(selector)) {
    form.addEventListener("submit", (event) => {
      if (!window.confirm(question(form.dataset.orderNumber))) {
        event.preventDefault();
      }
    });
  }
}

confirmBeforeSubmit(
  ".void-form",
  (orderNumber) =>
    "Void order #" + orderNumber + "? This only corrects the sales figures — settle the till and the paper by hand.",
);
confirmBeforeSubmit(
  ".comp-form",
  (orderNumber) =>
    "Comp order #" + orderNumber + "? The meal is free for a worker: the food still counts as served, and no money counts for it.",
);
confirmBeforeSubmit(
  ".reprint-form",
  (orderNumber) =>
    "Reprint order #" + orderNumber + "? The receipt and the kitchen ticket print again, marked REPRINT.",
);
