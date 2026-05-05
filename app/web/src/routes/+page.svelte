<svelte:head>
  <title>Apple Fest POS</title>
</svelte:head>

<main>
  <section>
    <p class="eyebrow">Apple Fest</p>
    <h1>Booth point of sale</h1>
    <p>Choose a mode for this tablet.</p>
    <div class="actions">
      <a href="/pos">Cashier POS</a>
      <a href="/kitchen">Kitchen display</a>
    </div>
  </section>
</main>

<style>
  main {
    min-height: 100vh;
    display: grid;
    place-items: center;
    padding: 2rem;
  }

  section {
    width: min(42rem, 100%);
    background: white;
    border-radius: 2rem;
    padding: 2rem;
    box-shadow: 0 1rem 3rem #7c2d1220;
  }

  .eyebrow {
    color: #ea580c;
    font-weight: 900;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }

  h1 {
    margin: 0;
    font-size: clamp(2.5rem, 8vw, 5rem);
    line-height: 0.95;
  }

  p {
    font-size: 1.25rem;
  }

  .actions {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(12rem, 1fr));
    gap: 1rem;
    margin-top: 2rem;
  }

  a {
    min-height: 5rem;
    display: grid;
    place-items: center;
    border-radius: 1rem;
    background: #ea580c;
    color: white;
    font-size: 1.25rem;
    font-weight: 900;
    text-decoration: none;
  }
</style>
