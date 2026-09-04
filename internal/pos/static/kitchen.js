// Redraws the kitchen board from /api/kitchen every few seconds. A meta
// refresh did this before, but that is a navigation, and a navigation drops
// the tablet out of full screen every time it fires.
"use strict";

const POLL_MS = 4000;

const clockElement = document.getElementById("clock");
const boardElement = document.getElementById("board-body");

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

function drawTicket(ticket) {
  const item = element("li", "ticket");

  const bar = element("div", "bar");
  bar.append(element("strong", "", "#" + ticket.orderNumber));
  bar.append(element("time", "", formatClock(ticket.createdAt)));
  item.append(bar);

  const lines = element("ul", "ticket-lines");
  for (const line of ticket.lines) {
    const lineItem = element("li");
    lineItem.append(element("span", "qty", line.quantity + "×"));
    lineItem.append(document.createTextNode(" " + line.name));
    if (line.side) {
      lineItem.append(element("em", "", line.side));
    }
    lines.append(lineItem);
  }
  item.append(lines);

  if (ticket.notes) {
    item.append(element("p", "ticket-notes", ticket.notes));
  }
  return item;
}

function draw(board) {
  clockElement.textContent = formatClock(board.serverTime);
  boardElement.innerHTML = "";
  if (board.tickets.length === 0) {
    boardElement.append(element("p", "empty", "No orders yet today."));
    return;
  }
  const list = element("ul", "tickets");
  for (const ticket of board.tickets) {
    list.append(drawTicket(ticket));
  }
  boardElement.append(list);
}

async function poll() {
  try {
    const response = await fetch("/api/kitchen");
    if (!response.ok) return;
    draw(await response.json());
  } catch {
    // Keep the stale board up; the next tick tries again.
  }
}

poll();
setInterval(poll, POLL_MS);
