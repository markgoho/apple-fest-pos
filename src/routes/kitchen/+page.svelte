<script lang="ts">
  let now = $state(new Date());

  $effect(() => {
    const interval = setInterval(() => {
      now = new Date();
    }, 1000);

    return () => clearInterval(interval);
  });
</script>

<svelte:head>
  <title>Kitchen Display</title>
</svelte:head>

<main>
  <header>
    <a href="/">Apple Fest POS</a>
    <time>{now.toLocaleTimeString()}</time>
  </header>

  <section>
    <h1>Kitchen tickets</h1>
    <p>Order polling will be added after the POS submit flow is stable.</p>
  </section>
</main>

<style>
  main {
    min-height: 100vh;
    padding: 1.5rem;
  }

  header {
    display: flex;
    justify-content: space-between;
    font-weight: 800;
  }

  section {
    margin-top: 2rem;
    background: white;
    border-radius: 1.5rem;
    padding: 2rem;
  }

  h1 {
    margin: 0 0 1rem;
    font-size: clamp(2rem, 5vw, 4rem);
  }
</style>
