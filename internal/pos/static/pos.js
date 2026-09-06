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
const sidesSheet = document.getElementById("sides-sheet");
const sidesTitle = document.getElementById("sides-title");
const sidesChoices = document.getElementById("sides-choices");
const sidesAdd = document.getElementById("sides-add");
const sidesCancel = document.getElementById("sides-cancel");
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

// lineKey identifies a cart line by what it is, not by when it was added, so
// the same pancake with the same toppings is one line however many times the
// Operator adds it. The chosen ids are already in menu order.
function lineKey(menuItemId, chosen) {
  return menuItemId + "|" + chosen.join(",");
}

// addToCart adds one of an item with its toppings already chosen, merging into
// the matching line when there is one.
function addToCart(item, chosen) {
  disarm();
  const key = lineKey(item.menuItemId, chosen);
  const existing = cart.find((line) => line.key === key);
  if (existing) {
    existing.quantity += 1;
  } else {
    cart.push({ key, ...item, chosen, quantity: 1 });
  }
  draw();
}

// toggleSide adds or removes one Side of a cart line that is already in the
// cart, so a topping asked for late is one tap and not a delete and a redo.
// The Sides stay in menu order, so the cart line, the receipt and the kitchen
// ticket all read the same way whatever order the Operator tapped.
function toggleSide(key, sideId) {
  disarm();
  const line = cart.find((candidate) => candidate.key === key);
  if (!line) return;

  const chosen = new Set(line.chosen);
  if (chosen.has(sideId)) {
    chosen.delete(sideId);
  } else {
    chosen.add(sideId);
  }
  line.chosen = line.sides.filter((side) => chosen.has(side.id)).map((side) => side.id);
  line.key = lineKey(line.menuItemId, line.chosen);

  // The edit can make this line the twin of another one. Fold it in, or the
  // cart shows the same pancake twice and the kitchen ticket prints it twice.
  const twin = cart.find((candidate) => candidate !== line && candidate.key === line.key);
  if (twin) {
    twin.quantity += line.quantity;
    cart = cart.filter((candidate) => candidate !== line);
  }
  draw();
}

// drawSides draws the topping row of a cart line. Every Side is one gold
// toggle, resting or chosen, because a Side costs nothing and is reversible
// (ADR-0004); a chosen Side is marked by aria-pressed and a checkmark, never
// by a second colour.
function drawSides(line) {
  const row = document.createElement("div");
  row.className = "sides";
  for (const side of line.sides) {
    const chosen = line.chosen.includes(side.id);
    const toggle = document.createElement("button");
    toggle.type = "button";
    toggle.className = "side-toggle";
    toggle.textContent = side.label;
    toggle.setAttribute("aria-pressed", chosen ? "true" : "false");
    toggle.addEventListener("click", () => toggleSide(line.key, side.id));
    row.append(toggle);
  }
  return row;
}

function drawLine(line) {
  const item = document.createElement("li");
  item.className = "cart-line";

  const name = document.createElement("p");
  name.className = "name";
  name.textContent = line.name;
  // A line whose toggles are all resting is a Plain one, and the dialog made
  // that a deliberate choice. Say so, or it reads as a pancake nobody was
  // asked about.
  if (line.sides.length > 0 && line.chosen.length === 0) {
    const tag = document.createElement("small");
    tag.textContent = "Plain";
    name.append(tag);
  }
  item.append(name);

  const quantity = document.createElement("div");
  quantity.className = "stepper";

  const less = document.createElement("button");
  less.type = "button";
  less.textContent = "−";
  less.setAttribute("aria-label", "Remove one " + line.name);
  less.addEventListener("click", () => changeQuantity(line.key, line.quantity - 1));

  const count = document.createElement("output");
  count.textContent = String(line.quantity);
  count.setAttribute("aria-label", line.name + " quantity");

  const more = document.createElement("button");
  more.type = "button";
  more.textContent = "+";
  more.setAttribute("aria-label", "Add one " + line.name);
  more.addEventListener("click", () => changeQuantity(line.key, line.quantity + 1));

  quantity.append(less, count, more);
  item.append(quantity);
  if (line.sides.length > 0) {
    item.append(drawSides(line));
  }
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
    items: cart.map((line) => ({
      menuItemId: line.menuItemId,
      quantity: line.quantity,
      sides: line.chosen.length > 0 ? line.chosen : undefined
    }))
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

// The toppings dialog. ADR-0005 turned a dialog down for the Place order
// checkpoint, because it would cover the cart lines at the moment they are
// read back. This one opens at add time and is gone before the read-back, so
// the checkpoint it protects is untouched.
let pending = null;

function setSidesOpen(open) {
  sidesSheet.hidden = !open;
  scrim.hidden = !open && moreSheet.hidden;
  if (!open) pending = null;
}

// choice draws one button of the dialog. Plain and the toppings are the same
// kind of tap and look the same; they differ only in that Plain means "none of
// these" and so clears them, and any topping clears Plain.
function choice(label, chosen, onPick) {
  const toggle = document.createElement("button");
  toggle.type = "button";
  toggle.className = "side-toggle";
  toggle.textContent = label;
  toggle.setAttribute("aria-pressed", chosen ? "true" : "false");
  toggle.addEventListener("click", onPick);
  return toggle;
}

function drawSidesChoices() {
  const plain = choice("Plain", pending.plain, () => {
    pending.plain = !pending.plain;
    if (pending.plain) pending.chosen = [];
    drawSidesChoices();
  });
  plain.classList.add("side-plain");

  const toppings = pending.item.sides.map((side) =>
    choice(side.label, pending.chosen.includes(side.id), () => {
      const picked = new Set(pending.chosen);
      if (picked.has(side.id)) {
        picked.delete(side.id);
      } else {
        picked.add(side.id);
      }
      pending.chosen = pending.item.sides.filter((one) => picked.has(one.id)).map((one) => one.id);
      if (pending.chosen.length > 0) pending.plain = false;
      drawSidesChoices();
    })
  );

  sidesChoices.replaceChildren(plain, ...toppings);
  // Add stays inert until the Operator has said which it is. A pancake with no
  // toppings and a pancake nobody has been asked about look the same in the
  // cart, and only one of them is an order.
  sidesAdd.disabled = !pending.plain && pending.chosen.length === 0;
}

function openSides(item) {
  setMoreOpen(false);
  pending = { item, chosen: [], plain: false };
  sidesTitle.textContent = item.name;
  drawSidesChoices();
  setSidesOpen(true);
}

// readSides unpacks the tile's "id:Label|id:Label" attribute. See
// menuTile.SidesAttribute in pages.go.
function readSides(packed) {
  if (!packed) return [];
  return packed.split("|").map((pair) => {
    const cut = pair.indexOf(":");
    return { id: pair.slice(0, cut), label: pair.slice(cut + 1) };
  });
}

for (const tile of menuElement.querySelectorAll(".tile")) {
  const item = {
    menuItemId: tile.dataset.menuItemId,
    name: tile.dataset.name,
    priceCents: Number(tile.dataset.priceCents),
    sides: readSides(tile.dataset.sides)
  };
  tile.addEventListener("click", () => {
    // An item with Sides asks for them first: the toppings are part of what is
    // being added, and a cart that has grown long is a poor place to hunt for
    // the line that has just appeared.
    if (item.sides.length > 0) {
      openSides(item);
      return;
    }
    addToCart(item, []);
  });
}

sidesAdd.addEventListener("click", () => {
  const { item, chosen } = pending;
  setSidesOpen(false);
  addToCart(item, chosen);
});
sidesCancel.addEventListener("click", () => setSidesOpen(false));

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
scrim.addEventListener("click", () => {
  setMoreOpen(false);
  setSidesOpen(false);
});

draw();
