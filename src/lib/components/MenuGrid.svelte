<script lang="ts">
  import type { MenuItem } from "#lib/types/api";
  import { formatCents } from "#lib/utils/currency";

  type Props = {
    items: MenuItem[];
    onAdd: (item: MenuItem) => void;
  };

  let { items, onAdd }: Props = $props();
</script>

<section class="menu-grid" aria-label="Menu items">
  {#each items as item (item.id)}
    <button type="button" class="menu-item" onclick={() => onAdd(item)}>
      <strong>{item.name}</strong>
      <span>{formatCents(item.priceCents)}</span>
    </button>
  {/each}
</section>

<style>
  .menu-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(12rem, 1fr));
    gap: 0.85rem;
  }

  .menu-item {
    min-height: 7.5rem;
    border: 2px solid #f6c56d;
    border-radius: 1.35rem;
    background: linear-gradient(160deg, var(--cream), #ffe8ab 58%, #facd71);
    color: var(--apple-red-dark);
    box-shadow: 0 0.35rem 0 var(--cider-gold), 0 0.9rem 1.5rem #5f32171f;
    font: inherit;
    text-align: left;
    padding: 1rem;
    display: grid;
    align-content: space-between;
    touch-action: manipulation;
    user-select: none;
  }

  .menu-item strong {
    font-size: clamp(1.25rem, 2.4vw, 1.6rem);
  }

  .menu-item span {
    width: fit-content;
    border-radius: 999px;
    background: var(--apple-red);
    color: white;
    font-size: clamp(1.1rem, 2vw, 1.35rem);
    font-weight: 800;
    padding: 0.25rem 0.6rem;
  }
</style>
