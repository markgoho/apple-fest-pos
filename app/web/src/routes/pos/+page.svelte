<script lang="ts">
  import { onMount } from "svelte";
  import CartPanel from "$lib/components/CartPanel.svelte";
  import MenuGrid from "$lib/components/MenuGrid.svelte";
  import SyncBanner from "$lib/components/SyncBanner.svelte";
  import { getMenu } from "$lib/services/api";
  import { listQueuedOrders } from "$lib/services/outbox";
  import {
    addMenuItem,
    buildOrderRequest,
    createEmptyCart,
    getTotalCents,
    updateLineNotes,
    updateQuantity,
    type CartState
  } from "$lib/stores/cart";
  import { createSyncState, refreshReachability, retryQueuedOrders, submitOrQueue, type SyncState } from "$lib/stores/sync";
  import type { MenuItem, PlaceOrderResponse } from "$lib/types/api";
  import { formatCents } from "$lib/utils/currency";
  import { getDeviceId } from "$lib/utils/device";
  import { createId } from "$lib/utils/id";

  let menuItems = $state<MenuItem[]>([]);
  let cart = $state<CartState>(createEmptyCart());
  let sync = $state<SyncState>(createSyncState());
  let loadingMenu = $state(true);
  let submitting = $state(false);
  let confirmation = $state<PlaceOrderResponse | null>(null);
  let notice = $state("");
  let deviceId = "";

  const totalCents = $derived(getTotalCents(cart));

  onMount(() => {
    deviceId = getDeviceId();
    void initialize();

    const interval = setInterval(() => {
      void updateSyncStatus();
    }, 5000);

    return () => clearInterval(interval);
  });

  async function initialize() {
    try {
      menuItems = await getMenu();
    } catch {
      notice = "Menu is unavailable until the server is reachable.";
    } finally {
      loadingMenu = false;
    }

    await updateSyncStatus();
  }

  async function updateSyncStatus() {
    const queuedOrders = await listQueuedOrders();
    const queuedCount = queuedOrders.length;

    sync = await refreshReachability({ ...sync, queuedCount });

    if (!sync.online || queuedCount === 0) {
      return;
    }

    sync = { ...sync, syncing: true, message: "Syncing queued orders..." };

    const syncedCount = await retryQueuedOrders();
    const remainingOrders = await listQueuedOrders();
    const syncedMessage = syncedCount > 0 ? `Synced ${syncedCount} queued order(s)` : sync.message;

    sync = {
      ...sync,
      syncing: false,
      queuedCount: remainingOrders.length,
      message: syncedMessage
    };
  }

  async function handleSubmit() {
    submitting = true;
    notice = "";
    confirmation = null;

    const clientOrderId = createId();
    const order = buildOrderRequest(cart, clientOrderId, deviceId);

    try {
      const result = await submitOrQueue(order);
      if (result.queued) {
        notice = "Server is offline. Order queued and will retry automatically.";
      } else {
        confirmation = result.response;
        notice = `Order #${result.response.order.orderNumber} accepted.`;
      }

      cart = createEmptyCart();
      await updateSyncStatus();
    } catch (error) {
      notice = error instanceof Error ? error.message : "Unable to submit order.";
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head>
  <title>Cashier POS</title>
</svelte:head>

<SyncBanner {sync} />

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

    {#if loadingMenu}
      <p>Loading menu...</p>
    {:else if menuItems.length === 0}
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
    height: calc(100dvh - 3rem);
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
