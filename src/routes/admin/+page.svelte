<script lang="ts">
  import { onMount } from "svelte";
  import { formatCents } from "#lib/utils/currency";
  import type { AdminSalesResponse } from "#lib/types/api";
  import type { PageProps } from "./$types";

  let { data }: PageProps = $props();

  const refreshIntervalMs = 5000;

  let sales = $derived<AdminSalesResponse>(data.sales);
  let error = $state<string | null>(null);
  let lastUpdated = $state<Date | null>(null);

  async function refresh() {
    try {
      const response = await fetch(`/api/admin/sales?date=${encodeURIComponent(sales.businessDate)}`);
      if (!response.ok) {
        error = `Request failed with ${response.status}`;
        return;
      }

      sales = (await response.json()) as AdminSalesResponse;
      error = null;
      lastUpdated = new Date();
    } catch {
      error = "Could not reach server";
    }
  }

  onMount(() => {
    const interval = setInterval(refresh, refreshIntervalMs);
    return () => clearInterval(interval);
  });

  function formatTime(iso: string): string {
    return new Date(iso).toLocaleTimeString();
  }
</script>

<svelte:head>
  <title>Admin · Sales</title>
</svelte:head>

<main>
  <header>
    <a href="/">Apple Fest POS</a>
    <div class="status">
      {#if error}
        <span class="error">{error}</span>
      {:else if lastUpdated}
        <span>Updated {lastUpdated.toLocaleTimeString()}</span>
      {:else}
        <span>Live</span>
      {/if}
    </div>
  </header>

    <section class="summary">
      <h1>Sales · {sales.businessDate}</h1>
      <div class="cards">
        <div class="card">
          <p class="label">Orders</p>
          <p class="value">{sales.summary.orderCount}</p>
        </div>
        <div class="card">
          <p class="label">Revenue</p>
          <p class="value">{formatCents(sales.summary.totalCents)}</p>
        </div>
        <div class="card" class:warn={sales.summary.printFailures > 0}>
          <p class="label">Print failures</p>
          <p class="value">{sales.summary.printFailures}</p>
        </div>
      </div>
    </section>

    <section class="panel">
      <h2>Items sold</h2>
      {#if sales.items.length === 0}
        <p class="empty">No items sold yet today.</p>
      {:else}
        <table>
          <thead>
            <tr>
              <th>Item</th>
              <th class="num">Qty</th>
              <th class="num">Revenue</th>
            </tr>
          </thead>
          <tbody>
            {#each sales.items as item (item.menuItemId)}
              <tr>
                <td>{item.name}</td>
                <td class="num">{item.quantity}</td>
                <td class="num">{formatCents(item.revenueCents)}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    </section>

    <section class="panel">
      <h2>Recent orders</h2>
      {#if sales.orders.length === 0}
        <p class="empty">No orders yet today.</p>
      {:else}
        <table>
          <thead>
            <tr>
              <th class="num">#</th>
              <th>Time</th>
              <th>Items</th>
              <th class="num">Total</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            {#each sales.orders as order (order.id)}
              <tr class:fail={order.status === "print_failed"}>
                <td class="num">{order.orderNumber}</td>
                <td>{formatTime(order.createdAt)}</td>
                <td>
                  {#each order.items as line, i (line.menuItemId)}
                    {#if i > 0}, {/if}{line.quantity}× {line.name}
                  {/each}
                </td>
                <td class="num">{formatCents(order.totalCents)}</td>
                <td>
                  {order.status}
                  {#if order.status === "print_failed"}
                    <small>(c:{order.customerPrintStatus} k:{order.kitchenPrintStatus})</small>
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    </section>
</main>

<style>
  main {
    min-height: 100vh;
    padding: 1.5rem;
    display: grid;
    gap: 1.5rem;
    max-width: 80rem;
    margin: 0 auto;
  }

  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 800;
  }

  .status .error {
    color: #b91c1c;
  }

  .summary h1 {
    margin: 0 0 1rem;
    font-size: clamp(1.75rem, 4vw, 2.5rem);
  }

  .cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(12rem, 1fr));
    gap: 1rem;
  }

  .card {
    background: white;
    border-radius: 1rem;
    padding: 1.25rem;
    box-shadow: 0 0.5rem 1.5rem #7c2d1218;
  }

  .card.warn {
    background: #fee2e2;
  }

  .card .label {
    margin: 0;
    font-size: 0.875rem;
    color: #7c2d12;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    font-weight: 800;
  }

  .card .value {
    margin: 0.25rem 0 0;
    font-size: 2.5rem;
    font-weight: 900;
    line-height: 1;
  }

  .panel {
    background: white;
    border-radius: 1.5rem;
    padding: 1.5rem;
    box-shadow: 0 0.5rem 1.5rem #7c2d1218;
  }

  .panel h2 {
    margin: 0 0 1rem;
    font-size: 1.5rem;
  }

  .empty {
    color: #6b7280;
    margin: 0;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th,
  td {
    padding: 0.5rem 0.75rem;
    text-align: left;
    border-bottom: 1px solid #f3e8d8;
  }

  th.num,
  td.num {
    text-align: right;
    font-variant-numeric: tabular-nums;
  }

  tr.fail {
    background: #fef2f2;
  }

  small {
    color: #7c2d12;
  }
</style>
