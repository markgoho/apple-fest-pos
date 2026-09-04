// Figures/Orders tabs are pure client state (the whole day is already in the
// rendered page); Void needs one confirm before it fires, per CONTEXT.md's
// Voided entry: the confirm text shows the order number so the Leader
// matches paper to screen (issue #45).
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

for (const form of document.querySelectorAll(".void-form")) {
  form.addEventListener("submit", (event) => {
    const orderNumber = form.dataset.orderNumber;
    const confirmed = window.confirm(
      "Void order #" + orderNumber + "? This only corrects the sales figures — settle the till and the paper by hand.",
    );
    if (!confirmed) {
      event.preventDefault();
    }
  });
}
