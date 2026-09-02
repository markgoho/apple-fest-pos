// Redraws the sales screen from /api/admin/sales every few seconds. A meta
// refresh did this before, but that is a navigation, and a navigation drops
// the tablet out of full screen every time it fires.
"use strict";

const POLL_MS = 5000;

const clockElement = document.getElementById("clock");
const businessDateElement = document.getElementById("business-date");
const orderCountElement = document.getElementById("order-count");
const revenueElement = document.getElementById("revenue");
const printFailuresElement = document.getElementById("print-failures");
const printFailuresCardElement = document.getElementById("print-failures-card");
const itemsBodyElement = document.getElementById("items-body");
const ordersBodyElement = document.getElementById("orders-body");

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

function formatClock(timestamp) {
  const date = new Date(timestamp);
  if (Number.isNaN(date.getTime())) {
    return timestamp;
  }
  return date.toLocaleTimeString([], { hour: "numeric", minute: "2-digit" });
}

function formatCents(cents) {
  const dollars = Math.floor(cents / 100);
  const rest = cents % 100;
  return rest === 0 ? "$" + dollars : "$" + dollars + "." + String(rest).padStart(2, "0");
}

function drawEmpty(container, message) {
  container.innerHTML = "";
  container.append(element("p", "empty", message));
}

function drawItems(items) {
  if (items.length === 0) {
    drawEmpty(itemsBodyElement, "No items sold yet today.");
    return;
  }
  const table = element("table");
  const thead = element("thead");
  const headRow = element("tr");
  headRow.append(element("th", "", "Item"), element("th", "num", "Qty"), element("th", "num", "Revenue"));
  thead.append(headRow);
  table.append(thead);

  const tbody = element("tbody");
  for (const item of items) {
    const row = element("tr");
    row.append(
      element("td", "", item.name),
      element("td", "num", String(item.quantity)),
      element("td", "num", formatCents(item.revenueCents)),
    );
    tbody.append(row);
  }
  table.append(tbody);

  itemsBodyElement.innerHTML = "";
  itemsBodyElement.append(table);
}

function orderItemsSummary(items) {
  return items.map((line) => line.quantity + "× " + line.name).join(", ");
}

function drawOrders(orders) {
  if (orders.length === 0) {
    drawEmpty(ordersBodyElement, "No orders yet today.");
    return;
  }
  const table = element("table");
  const thead = element("thead");
  const headRow = element("tr");
  headRow.append(
    element("th", "num", "#"),
    element("th", "", "Time"),
    element("th", "", "Items"),
    element("th", "num", "Total"),
    element("th", "", "Print"),
  );
  thead.append(headRow);
  table.append(thead);

  const tbody = element("tbody");
  for (const order of orders) {
    const row = element("tr", order.status === "print_failed" ? "fail" : "");
    row.append(
      element("td", "num", String(order.orderNumber)),
      element("td", "", formatClock(order.createdAt)),
      element("td", "", orderItemsSummary(order.items)),
      element("td", "num", formatCents(order.totalCents)),
      element(
        "td",
        "",
        order.status === "print_failed"
          ? "failed (customer " + order.customerPrintStatus + ", kitchen " + order.kitchenPrintStatus + ")"
          : order.status,
      ),
    );
    tbody.append(row);
  }
  table.append(tbody);

  ordersBodyElement.innerHTML = "";
  ordersBodyElement.append(table);
}

function draw(sales) {
  clockElement.textContent = formatClock(sales.serverTime);
  businessDateElement.textContent = sales.businessDate;
  orderCountElement.textContent = String(sales.summary.orderCount);
  revenueElement.textContent = formatCents(sales.summary.totalCents);
  printFailuresElement.textContent = String(sales.summary.printFailures);
  printFailuresCardElement.classList.toggle("warn", sales.summary.printFailures > 0);
  drawItems(sales.items);
  drawOrders(sales.orders);
}

async function poll() {
  try {
    const response = await fetch("/api/admin/sales" + window.location.search);
    if (!response.ok) return;
    draw(await response.json());
  } catch {
    // Keep the stale screen up; the next tick tries again.
  }
}

poll();
setInterval(poll, POLL_MS);
