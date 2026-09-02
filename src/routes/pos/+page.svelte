<script lang="ts">
  import { onMount } from "svelte";
  import CartPanel from "#lib/components/CartPanel.svelte";
  import MenuGrid from "#lib/components/MenuGrid.svelte";
  import {
    addMenuItem,
    buildOrderRequest,
    createEmptyCart,
    getTotalCents,
    updateLineNotes,
    updateQuantity,
    type CartState
  } from "#lib/stores/cart";
  import type { PlaceOrderResponse } from "#lib/types/api";
  import { formatCents } from "#lib/utils/currency";
  import { getDeviceId } from "#lib/utils/device";
  import { createId } from "#lib/utils/id";
  import type { PageProps } from "./$types";

  let { data }: PageProps = $props();

  const menuItems = $derived(data.menuItems);

  let cart = $state<CartState>(createEmptyCart());
  let submitting = $state(false);
  let confirmation = $state<PlaceOrderResponse | null>(null);
  let notice = $state("");
  let deviceId = "";

  const totalCents = $derived(getTotalCents(cart));

  onMount(() => {
    deviceId = getDeviceId();
  });

  async function handleSubmit() {
    submitting = true;
    notice = "";
    confirmation = null;

    const clientOrderId = createId();
    const order = buildOrderRequest(cart, clientOrderId, deviceId);

    try {
      const response = await fetch("/api/orders", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(order)
      });

      const body = await response.json().catch(() => null);

      if (!response.ok) {
        notice = typeof body?.error === "string" ? body.error : `Request failed with ${response.status}`;
        return;
      }

      confirmation = body as PlaceOrderResponse;
      notice = `Order #${confirmation.order.orderNumber} accepted.`;
      cart = createEmptyCart();
    } catch {
      notice = "No connection. Try again.";
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head>
  <title>Cashier POS</title>
</svelte:head>

<main>
  <section class="workspace">
    <header>
      <a href="/">Apple Fest POS</a>
    </header>

    {#if notice}
      <div class="notice" role="status">{notice}</div>
    {/if}

    {#if confirmation}
      <div class="confirmation" role="status">
        <strong>Order #{confirmation.order.orderNumber}</strong>
        <span>Total {formatCents(confirmation.order.totalCents)}</span>
        <span>Payment cash exact</span>
        <span>Kitchen print {confirmation.print.kitchen}</span>
      </div>
    {/if}

    {#if menuItems.length === 0}
      <p>No menu items are available.</p>
    {:else}
      <MenuGrid items={menuItems} onAdd={(item) => (cart = addMenuItem(cart, item))} />
    {/if}
  </section>

  <CartPanel
    {cart}
    {totalCents}
    {submitting}
    onQuantityChange={(cartItemId, quantity) => (cart = updateQuantity(cart, cartItemId, quantity))}
    onLineNotesChange={(cartItemId, notes) => (cart = updateLineNotes(cart, cartItemId, notes))}
    onSubmit={handleSubmit}
    onClear={() => {
      cart = createEmptyCart();
    }}
  />
</main>

<style>
  main {
    height: 100dvh;
    overflow: hidden;
    display: grid;
    grid-template-columns: 1fr minmax(22rem, 30rem);
  }

  .workspace {
    min-height: 0;
    overflow: auto;
    padding: 1rem;
    display: grid;
    gap: 1rem;
    align-content: start;
  }

  header {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 1rem;
  }

  a {
    color: var(--leaf-green-dark);
    font-weight: 900;
  }

  .notice,
  .confirmation {
    border: 2px solid #e4b95c;
    border-radius: 1rem;
    padding: 1rem;
    background: #fffaf0;
    box-shadow: 0 0.5rem 1.5rem #5f32171f;
  }

  .confirmation {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    color: var(--leaf-green);
    font-weight: 800;
  }

  @media (max-width: 840px), (orientation: portrait) {
    main {
      grid-template-columns: 1fr;
      grid-template-rows: auto minmax(0, 1fr);
    }

    .workspace {
      overflow: visible;
      padding-bottom: 0.75rem;
    }
  }
</style>
