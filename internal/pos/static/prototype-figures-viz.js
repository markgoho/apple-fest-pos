// PROTOTYPE data + rendering for issue #11. Throwaway — see the HTML file
// for the full disclaimer. Colors reuse pos.css's apple-fest palette.

const ITEMS = [
  { id: "potato-pancake", name: "Potato Pancake", priceCents: 1000, color: "#b91c1c" },
  { id: "og-toastie", name: "OG Toastie", priceCents: 500, color: "#d9a441" },
  { id: "pizza-toastie", name: "Pizza Toastie", priceCents: 600, color: "#2f6b2f" },
  { id: "harvest-toastie", name: "Harvest Toastie", priceCents: 800, color: "#5f3217" },
];

const HOURS = [8, 9, 10, 11, 12, 13, 14, 15, 16, 17];

// Fake units sold per item per hour, one row per item (same order as ITEMS).
const QTY = [
  [2, 4, 7, 10, 12, 9, 4, 3, 6, 8],
  [1, 3, 5, 8, 9, 6, 3, 2, 4, 5],
  [1, 2, 4, 7, 8, 6, 3, 2, 5, 6],
  [0, 2, 3, 5, 6, 5, 2, 1, 3, 4],
];

const ORDER_COUNT = 104;

function hourLabel(h) {
  const period = h < 12 ? "a" : "p";
  const display = h % 12 === 0 ? 12 : h % 12;
  return `${display}${period}`;
}

// item -> { qty, revenueCents }
const itemTotals = ITEMS.map((item, i) => {
  const qty = QTY[i].reduce((a, b) => a + b, 0);
  return { ...item, qty, revenueCents: qty * item.priceCents };
});

// hour index -> total revenue across items, for that hour
const hourlyRevenue = HOURS.map((_, h) =>
  ITEMS.reduce((sum, item, i) => sum + QTY[i][h] * item.priceCents, 0)
);

const totalRevenueCents = itemTotals.reduce((sum, item) => sum + item.revenueCents, 0);

function money(cents) {
  return `$${(cents / 100).toFixed(2)}`;
}

function legendHtml(items) {
  return `<div class="legend">${items
    .map((item) => `<span><i style="background:${item.color}"></i>${item.name}</span>`)
    .join("")}</div>`;
}

// One hourly stacked bar per hour, dollars by item, tallest bar full height.
function stackedBarChart(width, height) {
  const maxRevenue = Math.max(...hourlyRevenue);
  const barGap = 4;
  const barWidth = (width - barGap * (HOURS.length - 1)) / HOURS.length;
  const bars = HOURS.map((hour, h) => {
    const x = h * (barWidth + barGap);
    let y = height;
    let segments = "";
    ITEMS.forEach((item, i) => {
      const cents = QTY[i][h] * item.priceCents;
      const segHeight = (cents / maxRevenue) * height;
      y -= segHeight;
      segments += `<rect x="${x.toFixed(1)}" y="${y.toFixed(1)}" width="${barWidth.toFixed(1)}" height="${segHeight.toFixed(1)}" fill="${item.color}" />`;
    });
    const label = `<text x="${(x + barWidth / 2).toFixed(1)}" y="${height + 12}" font-size="7" text-anchor="middle" fill="#5f3217">${hourLabel(hour)}</text>`;
    return segments + label;
  });
  return `<svg viewBox="0 0 ${width} ${height + 14}" role="img" aria-label="Revenue by hour, by item">${bars.join("")}</svg>`;
}

// One horizontal bar, split into four segments by each item's share of the
// day's revenue — a second way to read the same numbers as a whole.
function shareBar(width, height) {
  let x = 0;
  const segments = itemTotals
    .map((item) => {
      const segWidth = (item.revenueCents / totalRevenueCents) * width;
      const rect = `<rect x="${x.toFixed(1)}" y="0" width="${segWidth.toFixed(1)}" height="${height}" fill="${item.color}" />`;
      x += segWidth;
      return rect;
    })
    .join("");
  return `<svg viewBox="0 0 ${width} ${height}" role="img" aria-label="Share of the day's revenue, by item">${segments}</svg>`;
}

// Four small per-item charts, each scaled to its own max so a slow item
// isn't flattened by Potato Pancake's volume.
function smallMultiples(width, height) {
  const cellGap = 10;
  const cellWidth = (width - cellGap) / 2;
  const cellHeight = (height - cellGap) / 2;
  const chartHeight = cellHeight - 16;
  const barGap = 2;
  const barWidth = (cellWidth - barGap * (HOURS.length - 1)) / HOURS.length;

  const cells = ITEMS.map((item, i) => {
    const maxQty = Math.max(...QTY[i]);
    const col = i % 2;
    const row = Math.floor(i / 2);
    const cellX = col * (cellWidth + cellGap);
    const cellY = row * (cellHeight + cellGap);
    const bars = HOURS.map((_, h) => {
      const barX = cellX + h * (barWidth + barGap);
      const barHeight = maxQty === 0 ? 0 : (QTY[i][h] / maxQty) * chartHeight;
      const barY = cellY + 12 + (chartHeight - barHeight);
      return `<rect x="${barX.toFixed(1)}" y="${barY.toFixed(1)}" width="${barWidth.toFixed(1)}" height="${barHeight.toFixed(1)}" fill="${item.color}" />`;
    }).join("");
    const label = `<text x="${cellX}" y="${cellY + 9}" font-size="8" font-weight="700" fill="${item.color}">${item.name}</text>`;
    return label + bars;
  }).join("");

  return `<svg viewBox="0 0 ${width} ${height}" role="img" aria-label="Units sold by hour, per item">${cells}</svg>`;
}

// A single line tracing money raised so far, hour over hour.
function cumulativeLine(width, height) {
  let running = 0;
  const points = hourlyRevenue.map((revenue, h) => {
    running += revenue;
    const x = (h / (HOURS.length - 1)) * width;
    const y = height - (running / totalRevenueCents) * height;
    return `${x.toFixed(1)},${y.toFixed(1)}`;
  });
  const path = `<polyline points="${points.join(" ")}" fill="none" stroke="#b91c1c" stroke-width="2.5" />`;
  const dots = points
    .map((p) => {
      const [x, y] = p.split(",");
      return `<circle cx="${x}" cy="${y}" r="2.5" fill="#7f1d1d" />`;
    })
    .join("");
  return `<svg viewBox="0 0 ${width} ${height}" role="img" aria-label="Cumulative revenue across the day">${path}${dots}</svg>`;
}

// --- Shell pieces shared by every variant, matching leader.html's markup ---

function headerHtml() {
  return `
    <header class="bar"><a class="home" href="#">Apple Fest POS</a></header>
    <h1>Leader &middot; <span>2026-10-03</span></h1>
    <div class="day-toggle" role="tablist" aria-label="Event day">
      <button type="button" role="tab" aria-selected="true">Sat</button>
      <button type="button" role="tab" aria-selected="false">Sun</button>
      <button type="button" role="tab" aria-selected="false">Event</button>
    </div>
    <div class="tabs" role="tablist">
      <button type="button" role="tab" aria-selected="true">Figures</button>
      <button type="button" role="tab" aria-selected="false">Orders</button>
    </div>`;
}

function cardsHtml() {
  return `
    <section class="cards">
      <div class="card">
        <p class="label">Orders</p>
        <p class="value">${ORDER_COUNT}</p>
      </div>
      <div class="card">
        <p class="label">Revenue</p>
        <p class="value">${money(totalRevenueCents)}</p>
      </div>
    </section>`;
}

function tableHtml() {
  const rows = itemTotals
    .map(
      (item) =>
        `<tr><td>${item.name}</td><td class="num">${item.qty}</td><td class="num">${money(item.revenueCents)}</td></tr>`
    )
    .join("");
  return `
    <h2>Items sold</h2>
    <table>
      <thead><tr><th>Item</th><th class="num">Qty</th><th class="num">Revenue</th></tr></thead>
      <tbody>${rows}</tbody>
    </table>`;
}

// --- The three candidate visualization sets ---

const VARIANTS = {
  A: {
    name: "Chart first: hourly stack only",
    build() {
      const chart = `
        <div class="chart-block">
          <h3>Revenue by hour</h3>
          ${stackedBarChart(300, 110)}
          ${legendHtml(ITEMS)}
        </div>`;
      return `<section class="panel">${chart}${cardsHtml()}${tableHtml()}</section>`;
    },
  },
  B: {
    name: "Stack + share bar, between cards and table",
    build() {
      const charts = `
        <div class="chart-block">
          <h3>Revenue by hour</h3>
          ${stackedBarChart(300, 90)}
          <h3>Share of the day</h3>
          ${shareBar(300, 22)}
          ${legendHtml(ITEMS)}
        </div>`;
      return `<section class="panel">${cardsHtml()}${charts}${tableHtml()}</section>`;
    },
  },
  C: {
    name: "Small multiples + running total, below the table",
    build() {
      const charts = `
        <div class="chart-block">
          <h3>Units sold by hour, per item</h3>
          ${smallMultiples(300, 150)}
          <h3>Money raised so far</h3>
          ${cumulativeLine(300, 60)}
        </div>`;
      return `<section class="panel">${cardsHtml()}${tableHtml()}${charts}</section>`;
    },
  },
};

// --- Switcher + routing ---

const VARIANT_KEYS = Object.keys(VARIANTS);

function currentVariant() {
  const key = new URLSearchParams(location.search).get("variant");
  return VARIANT_KEYS.includes(key) ? key : VARIANT_KEYS[0];
}

function setVariant(key) {
  const url = new URL(location.href);
  url.searchParams.set("variant", key);
  history.replaceState(null, "", url);
  render();
}

function cycle(step) {
  const i = VARIANT_KEYS.indexOf(currentVariant());
  const next = VARIANT_KEYS[(i + step + VARIANT_KEYS.length) % VARIANT_KEYS.length];
  setVariant(next);
}

function render() {
  const key = currentVariant();
  document.getElementById("app").innerHTML = headerHtml() + VARIANTS[key].build();
  document.getElementById("switcher").innerHTML = `
    <div class="pill">
      <button type="button" id="prev" aria-label="Previous variant">&larr;</button>
      <span>${key} &middot; ${VARIANTS[key].name}</span>
      <button type="button" id="next" aria-label="Next variant">&rarr;</button>
    </div>`;
  document.getElementById("prev").addEventListener("click", () => cycle(-1));
  document.getElementById("next").addEventListener("click", () => cycle(1));
}

document.addEventListener("keydown", (event) => {
  const tag = document.activeElement.tagName;
  if (tag === "INPUT" || tag === "TEXTAREA" || document.activeElement.isContentEditable) return;
  if (event.key === "ArrowLeft") cycle(-1);
  if (event.key === "ArrowRight") cycle(1);
});

render();
