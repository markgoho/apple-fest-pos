<script lang="ts">
  let { children } = $props();
</script>

{@render children()}

<style>
  :global(:root) {
    --apple-red: #b91c1c;
    --apple-red-dark: #7f1d1d;
    --leaf-green: #2f6b2f;
    --leaf-green-dark: #225022;
    --cider-gold: #d97706;
    --cream: #fff7df;
    --parchment: #f8e7bc;
    --orchard-brown: #5f3217;
  }

  :global(*) {
    box-sizing: border-box;
  }

  :global(html),
  :global(body) {
    height: 100%;
  }

  :global(body) {
    margin: 0;
    min-height: 100dvh;
    background:
      radial-gradient(circle at 12% 8%, #fef3c7 0 10rem, transparent 10.5rem),
      linear-gradient(135deg, var(--cream), #f6d28d 55%, #e8b15f);
    color: var(--orchard-brown);
    font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  }

  :global(a) {
    color: inherit;
  }
</style>
