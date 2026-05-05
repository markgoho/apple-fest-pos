<script lang="ts">
  import type { SyncState } from "$lib/stores/sync";

  type Props = {
    sync: SyncState;
  };

  let { sync }: Props = $props();
</script>

<div class:online={sync.online} class:offline={!sync.online} class="banner" role="status">
  <strong>{sync.online ? "Online" : "Offline"}</strong>
  <span>{sync.syncing ? "Syncing queued orders..." : sync.message}</span>
  {#if sync.queuedCount > 0}
    <span>{sync.queuedCount} queued</span>
  {/if}
</div>

<style>
  .banner {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 0.75rem 1rem;
    color: white;
    font-weight: 800;
    box-shadow: 0 0.2rem 0 #5f32172b;
  }

  .online {
    background: var(--leaf-green);
  }

  .offline {
    background: linear-gradient(90deg, var(--cider-gold), var(--apple-red-dark));
  }
</style>
