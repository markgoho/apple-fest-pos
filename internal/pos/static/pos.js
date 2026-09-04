// The cart of the cashier screen. It is client-side state: a failed submit
// keeps the cart in memory, and the same clientOrderId goes out again, so the
// server replays the order instead of selling it twice.
"use strict";

const menuElement = document.getElementById("menu");
const linesElement = document.getElementById("lines");
const totalElement = document.getElementById("total");
const faultElement = document.getElementById("fault");
const lastOrderElement = document.getElementById("last-order");
const submitButton = document.getElementById("submit");
const clearButton = document.getElementById("clear");
const moreButton = document.getElementById("more");
const moreSheet = document.getElementById("more-sheet");
const moreClose = document.getElementById("more-close");
const scrim = document.getElementById("scrim");
const readBackElement = document.getElementById("read-back");
const reprintButton = document.getElementById("reprint");

let cart = [];
let clientOrderId = newId();
let submitting = false;
let lastOrderId = null;

// ADR-0005: placing an order takes two taps. The first arms the control and
// commits nothing; the second sends the order. Between them the control is
// inert for half a second, which absorbs a fumbled double tap. Any change to
// the cart disarms, so an Operator can never place an order that differs from
// the one they just read back.
const ARMING_MS = 500;
const PLACE_LABELS = { review: "Review order", checking: "Check the order…", armed: "Place order" };
let placeState = "review";
let armingTimer = null;

function disarm() {
  if (armingTimer !== null) {
    clearTimeout(armingTimer);
    armingTimer = null;
  }
  placeState = "review";
}

function armPlace() {
  placeState = "checking";
  draw();
  armingTimer = setTimeout(() => {
    armingTimer = null;
    placeState = "armed";
    draw();
  }, ARMING_MS);
}

function newId() {
  if (globalThis.crypto && globalThis.crypto.randomUUID) {
    return globalThis.crypto.randomUUID();
  }
  return Date.now().toString(36) + "-" + Math.random().toString(36).slice(2);
}

function deviceId() {
  const key = "apple-fest-pos-device-id";
  try {
    let stored = localStorage.getItem(key);
    if (!stored) {
      stored = newId();
      localStorage.setItem(key, stored);
    }
    return stored;
  } catch {
    return "unknown-device";
  }
}

function formatCents(cents) {
  const dollars = Math.floor(cents / 100);
  const rest = cents % 100;
  return rest === 0 ? "$" + dollars : "$" + dollars + "." + String(rest).padStart(2, "0");
}

function totalCents() {
  return cart.reduce((total, line) => total + line.priceCents * line.quantity, 0);
}

function showFault(message) {
  faultElement.textContent = message;
}

function changeQuantity(key, quantity) {
  disarm();
  cart = cart.filter((line) => {
    if (line.key !== key) {
      return true;
    }
    line.quantity = quantity;
    return quantity > 0;
  });
  draw();
}

function drawLine(line) {
  const item = document.createElement("li");
  item.className = "cart-line";

  const name = document.createElement("p");
  name.className = "name";
  name.textContent = line.name;
  if (line.sideLabel) {
    const tag = document.createElement("small");
    tag.textContent = line.sideLabel;
    name.append(tag);
  }
  item.append(name);

  const quantity = document.createElement("div");
  quantity.className = "stepper";

  const less = document.createElement("button");
  less.type = "button";
  less.textContent = "−";
  less.setAttribute("aria-label", "Remove one " + line.name + (line.sideLabel ? " " + line.sideLabel : ""));
  less.addEventListener("click", () => changeQuantity(line.key, line.quantity - 1));

  const count = document.createElement("output");
  count.textContent = String(line.quantity);
  count.setAttribute("aria-label", line.name + " quantity");

  const more = document.createElement("button");
  more.type = "button";
  more.textContent = "+";
  more.setAttribute("aria-label", "Add one " + line.name + (line.sideLabel ? " " + line.sideLabel : ""));
  more.addEventListener("click", () => changeQuantity(line.key, line.quantity + 1));

  quantity.append(less, count, more);
  item.append(quantity);
  return item;
}

function draw() {
  linesElement.replaceChildren(...cart.map(drawLine));
  totalElement.textContent = formatCents(totalCents());
  const empty = cart.length === 0;
  if (empty) disarm();
  submitButton.disabled = empty || submitting || placeState === "checking";
  clearButton.disabled = empty || submitting;
  submitButton.textContent = submitting ? "Sending…" : PLACE_LABELS[placeState];
  submitButton.classList.toggle("is-armed", placeState === "armed" && !submitting);
  readBackElement.hidden = empty || submitting || placeState === "review";
}

function clearCart() {
  disarm();
  cart = [];
  clientOrderId = newId();
  draw();
}

// printFault names the physical printer, not the document, because that is
// what the Operator can act on.
function printFault(print) {
  const faults = [];
  if (print.customer === "failed") faults.push("Window Printer did not print.");
  if (print.kitchen === "failed") faults.push("Kitchen Printer did not print.");
  return faults.join(" ");
}

async function submitOrder() {
  submitting = true;
  showFault("");
  draw();

  const body = {
    clientOrderId: clientOrderId,
    deviceId: deviceId(),
    payment: { method: "cash" },
    items: cart.map((line) => ({ menuItemId: line.menuItemId, quantity: line.quantity, side: line.sideId || undefined }))
  };

  try {
    const response = await fetch("/api/orders", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body)
    });
    const answer = await response.json().catch(() => null);

    if (!response.ok) {
      showFault(answer && answer.error ? answer.error : "Request failed with " + response.status);
      return;
    }

    lastOrderId = answer.order.id;
    lastOrderElement.textContent = "#" + answer.order.orderNumber;
    reprintButton.disabled = false;
    showFault(printFault(answer.print));
    clearCart();
  } catch {
    // The cart and the clientOrderId stay, so the next tap is the same order.
    // The control stays armed, so the retry is one tap and not three.
    showFault("Not sent. Tap Place order again.");
  } finally {
    submitting = false;
    draw();
  }
}

async function reprintOrder() {
  if (!lastOrderId || submitting) return;
  submitting = true;
  reprintButton.disabled = true;
  reprintButton.textContent = "Sending…";
  showFault("");
  draw();

  try {
    const response = await fetch("/api/orders/" + lastOrderId + "/reprint", { method: "POST" });
    const answer = await response.json().catch(() => null);
    if (!response.ok) {
      showFault(answer && answer.error ? answer.error : "Reprint failed with " + response.status);
    } else {
      showFault(printFault(answer.print));
    }
  } catch {
    showFault("Not sent. Tap Reprint again.");
  } finally {
    submitting = false;
    reprintButton.disabled = false;
    reprintButton.textContent = "Reprint";
    draw();
  }
}

for (const tile of menuElement.querySelectorAll(".tile")) {
  const item = {
    menuItemId: tile.dataset.menuItemId,
    name: tile.dataset.name,
    priceCents: Number(tile.dataset.priceCents),
    sideId: tile.dataset.sideId || "",
    sideLabel: tile.dataset.sideLabel || ""
  };
  tile.addEventListener("click", () => {
    disarm();
    const key = item.menuItemId + "|" + item.sideId;
    const existing = cart.find((line) => line.key === key);
    if (existing) {
      existing.quantity += 1;
    } else {
      cart.push({ key, ...item, quantity: 1 });
    }
    draw();
  });
}

// ADR-0005: Clear lives behind the More sheet, so the only destructive control
// on /pos needs two deliberate taps and is never on the screen during a sale.
function setMoreOpen(open) {
  moreSheet.hidden = !open;
  scrim.hidden = !open;
  moreButton.setAttribute("aria-expanded", open ? "true" : "false");
}

submitButton.addEventListener("click", () => {
  if (placeState === "armed") {
    submitOrder();
  } else if (placeState === "review") {
    armPlace();
  }
});
clearButton.addEventListener("click", () => {
  clearCart();
  setMoreOpen(false);
});
reprintButton.addEventListener("click", reprintOrder);
moreButton.addEventListener("click", () => setMoreOpen(moreSheet.hidden));
moreClose.addEventListener("click", () => setMoreOpen(false));
scrim.addEventListener("click", () => setMoreOpen(false));

draw();
