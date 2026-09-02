<script lang="ts">
  import { onMount } from "svelte";
  import type { KitchenBoard } from "#lib/types/api";
  import type { PageProps } from "./$types";

  let { data }: PageProps = $props();

  const refreshIntervalMs = 4000;

  let board = $derived<KitchenBoard>(data.board);
  let now = $state(new Date());
  let error = $state<string | null>(null);

  async function refresh() {
    try {
      const response = await fetch("/api/kitchen");
      if (!response.ok) {
        error = `Request failed with ${response.status}`;
        return;
      }

      board = (await response.json()) as KitchenBoard;
      error = null;
    } catch {
      error = "Could not reach server";
    }
  }

  onMount(() => {
    const clock = setInterval(() => (now = new Date()), 1000);
    const poll = setInterval(refresh, refreshIntervalMs);

    return () => {
      clearInterval(clock);
      clearInterval(poll);
    };
  });

  function formatTime(iso: string): string {
    return new Date(iso).toLocaleTimeString();
  }
</script>

<svelte:head>
  <title>Kitchen Display</title>
</svelte:head>

<main>
  <header>
    <a href="/">Apple Fest POS</a>
    {#if error}
      <span class="error">{error}</span>
    {/if}
    <time>{now.toLocaleTimeString()}</time>
  </header>

  {#if board.tickets.length === 0}
    <p class="empty">No orders yet today.</p>
  {:else}
    <ul>
      {#each board.tickets as ticket (ticket.orderNumber)}
        <li>
          <div class="ticket-head">
            <strong>#{ticket.orderNumber}</strong>
            <time>{formatTime(ticket.createdAt)}</time>
          </div>
          <ul class="lines">
            {#each ticket.lines as line (line.menuItemId)}
              <li>
                <span class="qty">{line.quantity}×</span>
                {line.name}
                {#if line.notes}<em>{line.notes}</em>{/if}
              </li>
            {/each}
          </ul>
          {#if ticket.notes}<p class="notes">{ticket.notes}</p>{/if}
        </li>
      {/each}
    </ul>
  {/if}
</main>

<style>
  main {
    min-height: 100vh;
    padding: 1.5rem;
  }

  header {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    gap: 1rem;
    font-weight: 800;
  }

  .error {
    color: #b91c1c;
  }

  .empty {
    margin-top: 2rem;
    font-size: 1.5rem;
  }

  ul {
    list-style: none;
    margin: 1.5rem 0 0;
    padding: 0;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(16rem, 1fr));
    gap: 1rem;
  }

  ul > li {
    background: white;
    border-radius: 1.5rem;
    padding: 1.25rem;
    box-shadow: 0 0.5rem 1.5rem #7c2d1218;
  }

  .ticket-head {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    font-size: 1.75rem;
  }

  .lines {
    display: block;
    margin: 0.75rem 0 0;
    font-size: 1.25rem;
  }

  .lines li {
    padding: 0.25rem 0;
  }

  .qty {
    font-weight: 900;
    color: var(--apple-red);
  }

  em {
    display: block;
    color: var(--leaf-green-dark);
  }

  .notes {
    margin: 0.75rem 0 0;
    color: var(--leaf-green-dark);
  }
</style>
