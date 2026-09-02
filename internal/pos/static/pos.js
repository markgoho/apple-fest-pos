// The cart of the cashier screen. It is client-side state: a failed submit
// keeps the cart in memory, and the same clientOrderId goes out again, so the
// server replays the order instead of selling it twice.
"use strict";

const linesElement = document.getElementById("lines");
const emptyElement = document.getElementById("empty");
const totalElement = document.getElementById("total");
const noticeElement = document.getElementById("notice");
const submitButton = document.getElementById("submit");
const clearButton = document.getElementById("clear");

let cart = [];
let clientOrderId = newId();
let submitting = false;

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

function showNotice(message, bad) {
  noticeElement.textContent = message;
  noticeElement.classList.toggle("bad", Boolean(bad));
}

// addItem starts a new line for an item that carries Sides, because each one
// gets its own choice. Every other item merges into its existing line.
function addItem(item) {
  if (item.sides.length === 0) {
    const existing = cart.find((line) => line.menuItemId === item.id);
    if (existing) {
      existing.quantity += 1;
      draw();
      return;
    }
  }
  cart.push({ key: newId(), menuItemId: item.id, name: item.name, priceCents: item.priceCents, sides: item.sides, side: "", quantity: 1 });
  draw();
}

function changeQuantity(key, quantity) {
  cart = cart.filter((line) => {
    if (line.key !== key) {
      return true;
    }
    line.quantity = quantity;
    return quantity > 0;
  });
  draw();
}

function chooseSide(line, side) {
  line.side = line.side === side ? "" : side;
  draw();
}

function element(tag, className, text) {
  const node = document.createElement(tag);
  if (className) {
    node.className = className;
  }
  if (text !== undefined) {
    node.textContent = text;
  }
  return node;
}

function drawLine(line) {
  const article = element("article", "line");
  const main = element("div", "line-main");
  main.append(element("strong", "", line.name));

  if (line.sides.length > 0) {
    const remove = element("button", "remove", "Remove");
    remove.type = "button";
    remove.setAttribute("aria-label", "Remove " + line.name);
    remove.addEventListener("click", () => changeQuantity(line.key, 0));
    main.append(remove);
  } else {
    const quantity = element("div", "quantity");
    const less = element("button", "", "−");
    less.type = "button";
    less.setAttribute("aria-label", "Remove one " + line.name);
    less.addEventListener("click", () => changeQuantity(line.key, line.quantity - 1));

    const count = element("output", "", String(line.quantity));
    count.setAttribute("aria-label", line.name + " quantity");

    const more = element("button", "add", "+");
    more.type = "button";
    more.setAttribute("aria-label", "Add one " + line.name);
    more.addEventListener("click", () => changeQuantity(line.key, line.quantity + 1));

    quantity.append(less, count, more);
    main.append(quantity);
  }
  article.append(main);

  if (line.sides.length > 0) {
    const sides = element("div", "sides");
    sides.setAttribute("aria-label", "Side for " + line.name);
    for (const side of line.sides) {
      const choice = element("button", "", side);
      choice.type = "button";
      choice.setAttribute("aria-pressed", String(line.side === side));
      choice.addEventListener("click", () => chooseSide(line, side));
      sides.append(choice);
    }
    article.append(sides);
  }
  return article;
}

function draw() {
  linesElement.replaceChildren(...cart.map(drawLine));
  emptyElement.hidden = cart.length > 0;
  totalElement.textContent = formatCents(totalCents());
  submitButton.disabled = cart.length === 0 || submitting;
  clearButton.disabled = cart.length === 0 || submitting;
  submitButton.textContent = submitting ? "Submitting..." : "Submit order";
}

function clearCart() {
  cart = [];
  clientOrderId = newId();
  draw();
}

async function submitOrder() {
  submitting = true;
  showNotice("");
  draw();

  const body = {
    clientOrderId: clientOrderId,
    deviceId: deviceId(),
    payment: { method: "cash" },
    items: cart.map((line) => ({ menuItemId: line.menuItemId, quantity: line.quantity, notes: line.side || undefined }))
  };

  try {
    const response = await fetch("/api/orders", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body)
    });
    const answer = await response.json().catch(() => null);

    if (!response.ok) {
      showNotice(answer && answer.error ? answer.error : "Request failed with " + response.status, true);
      return;
    }

    showNotice("Order #" + answer.order.orderNumber + " accepted · " + formatCents(answer.order.totalCents) + " · kitchen print " + answer.print.kitchen);
    clearCart();
  } catch {
    // The cart and the clientOrderId stay, so the next tap is the same order.
    showNotice("No connection. Try again.", true);
  } finally {
    submitting = false;
    draw();
  }
}

for (const button of document.querySelectorAll(".menu-item")) {
  const item = {
    id: button.dataset.menuItemId,
    name: button.dataset.name,
    priceCents: Number(button.dataset.priceCents),
    sides: button.dataset.sides ? button.dataset.sides.split("|") : []
  };
  button.addEventListener("click", () => addItem(item));
}

submitButton.addEventListener("click", submitOrder);
clearButton.addEventListener("click", () => {
  clearCart();
  showNotice("");
});

draw();
