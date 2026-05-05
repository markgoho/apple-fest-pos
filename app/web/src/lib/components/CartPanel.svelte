<script lang="ts">
  import type { CartState } from "$lib/stores/cart";
  import { formatCents } from "$lib/utils/currency";

  const potatoPancakeSides = ["Sour cream", "Applesauce", "Ketchup"];

  type Props = {
    cart: CartState;
    totalCents: number;
    submitting: boolean;
    onQuantityChange: (cartItemId: string, quantity: number) => void;
    onLineNotesChange: (cartItemId: string, notes: string) => void;
    onSubmit: () => void;
    onClear: () => void;
  };

  let {
    cart,
    totalCents,
    submitting,
    onQuantityChange,
    onLineNotesChange,
    onSubmit,
    onClear
  }: Props = $props();
</script>

<aside class="cart" aria-label="Current order">
  <div class="cart-header">
    <h2>Order</h2>
    <button type="button" class="ghost" onclick={onClear} disabled={cart.items.length === 0 || submitting}>Clear</button>
  </div>

  {#if cart.items.length === 0}
    <p class="empty">Tap menu items to start an order.</p>
  {:else}
    <div class="lines">
      {#each cart.items as item (item.id)}
        <article class="line">
          <div class="line-main">
            <strong>{item.menuItem.name}</strong>
            {#if item.menuItem.id === "potato-pancake"}
              <button type="button" class="remove" aria-label={`Remove ${item.menuItem.name}`} onclick={() => onQuantityChange(item.id, 0)}>Remove</button>
            {:else}
              <div class="quantity">
                <button type="button" aria-label={`Remove one ${item.menuItem.name}`} onclick={() => onQuantityChange(item.id, item.quantity - 1)}>-</button>
                <span>{item.quantity}</span>
                <button type="button" class="add" aria-label={`Add one ${item.menuItem.name}`} onclick={() => onQuantityChange(item.id, item.quantity + 1)}>+</button>
              </div>
            {/if}
          </div>
          {#if item.menuItem.id === "potato-pancake"}
            <div class="sides" aria-label="Potato pancake side">
              {#each potatoPancakeSides as side (side)}
                <button
                  type="button"
                  class:selected={item.notes === side}
                  onclick={() => onLineNotesChange(item.id, item.notes === side ? "" : side)}
                >
                  {side}
                </button>
              {/each}
            </div>
          {/if}
        </article>
      {/each}
    </div>
  {/if}

  <div class="checkout">
    <strong class="total">{formatCents(totalCents)}</strong>

    <button type="button" class="submit" onclick={onSubmit} disabled={cart.items.length === 0 || submitting}>
      {submitting ? "Submitting..." : "Submit order"}
    </button>
  </div>
</aside>

<style>
  .cart {
    height: 100%;
    min-height: 0;
    background: linear-gradient(180deg, var(--parchment), var(--cream));
    border-left: 3px solid #d59a45;
    box-shadow: inset 0.35rem 0 0 #f1c36a;
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    min-width: min(28rem, 100%);
  }

  .cart-header,
  .line-main {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
  }

  h2,
  p {
    margin: 0;
  }

  h2 {
    color: var(--apple-red-dark);
  }

  .empty {
    color: var(--orchard-brown);
  }

  .lines {
    display: grid;
    gap: 0.65rem;
    overflow: auto;
    padding-right: 0.25rem;
  }

  .line {
    display: grid;
    gap: 0.7rem;
    background: #fffaf0;
    border: 1px solid #e4b95c;
    border-radius: 1rem;
    padding: 0.85rem;
  }

  .line strong {
    display: block;
    color: var(--apple-red-dark);
    flex: 1;
    font-size: 1.3rem;
    line-height: 1.1;
  }

  .line span {
    display: block;
    color: var(--orchard-brown);
    font-size: 0.95rem;
  }

  .sides {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 0.5rem;
  }

  .sides button {
    background: #f9db95;
    color: var(--apple-red-dark);
  }

  .sides .selected {
    background: var(--leaf-green);
    color: white;
  }

  .quantity {
    display: grid;
    grid-template-columns: 3rem 2.5rem 3rem;
    align-items: center;
    gap: 0.35rem;
    text-align: center;
  }

  .quantity span {
    color: var(--apple-red-dark);
    font-size: 1.35rem;
    font-weight: 900;
  }

  .quantity button {
    min-height: 3rem;
    padding: 0;
  }

  .quantity .add {
    background: var(--leaf-green);
  }

  button {
    font: inherit;
    min-height: 3rem;
    border: 0;
    border-radius: 0.8rem;
    padding: 0.5rem 0.9rem;
    background: var(--apple-red);
    color: white;
    font-weight: 800;
    touch-action: manipulation;
    user-select: none;
  }

  button:disabled {
    opacity: 0.55;
  }

  .ghost {
    background: #f9db95;
    color: var(--apple-red-dark);
  }

  .remove {
    background: var(--apple-red);
    color: white;
  }

  .checkout {
    margin-top: auto;
    border-radius: 1rem;
    background: #fffaf0;
    box-shadow: 0 -0.35rem 1rem #5f32171a;
    display: grid;
    grid-template-columns: minmax(9rem, 1fr) 1.3fr;
    align-items: center;
    gap: 0.85rem;
    padding: 0.9rem;
  }

  .total {
    color: var(--apple-red-dark);
    font-size: clamp(2rem, 6vw, 3rem);
    font-weight: 900;
    line-height: 1;
  }

  .submit {
    min-height: 4.75rem;
    background: var(--leaf-green);
    box-shadow: 0 0.35rem 0 var(--leaf-green-dark);
    font-size: 1.45rem;
  }
</style>
