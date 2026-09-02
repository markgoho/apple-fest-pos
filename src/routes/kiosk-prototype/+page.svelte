<!--
  PROTOTYPE - throwaway. Ticket #6: how close to a kiosk can browser APIs get
  on a Pixel tablet in Chrome, with no device settings changed?
  Open http://<pi-or-dev-host>:<port>/kiosk-prototype on the tablet.
-->
<script lang="ts">
  import { onMount } from "svelte";
  import { beforeNavigate, pushState } from "$app/navigation";

  let guardBack = $state(true);
  let blockMenu = $state(true);
  let confirmUnload = $state(true);

  let started = $state(false);
  let isFullscreen = $state(false);
  let wakeLockSupported = $state(false);
  let wakeLockActive = $state(false);
  let secureContext = $state(false);
  let standalone = $state(false);
  let backCaught = $state(0);
  let menuBlocked = $state(0);
  let entries = $state(0);
  let log = $state<{ id: number; text: string }[]>([]);
  let logId = 0;

  let wakeVideo: HTMLVideoElement | undefined = $state();
  let lock: WakeLockSentinel | null = null;

  function note(message: string) {
    const stamp = new Date().toLocaleTimeString();
    logId += 1;
    log = [{ id: logId, text: `${stamp}  ${message}` }, ...log].slice(0, 40);
  }

  async function takeWakeLock() {
    if (!("wakeLock" in navigator)) {
      note("wake lock: API missing");
      return;
    }
    try {
      lock = await navigator.wakeLock.request("screen");
      wakeLockActive = true;
      note("wake lock: held");
      lock.addEventListener("release", () => {
        wakeLockActive = false;
        note("wake lock: released by the system");
      });
    } catch (error) {
      note(`wake lock: refused (${(error as Error).name}) - needs HTTPS or localhost`);
    }
  }

  function pushSentinel() {
    pushState("", { kiosk: true });
    entries += 1;
  }

  async function startShift() {
    try {
      await document.documentElement.requestFullscreen({ navigationUI: "hide" });
      note("full screen: entered");
    } catch (error) {
      note(`full screen: refused (${(error as Error).name})`);
    }
    await takeWakeLock();
    try {
      if (wakeVideo && !wakeVideo.srcObject) {
        const canvas = document.createElement("canvas");
        canvas.width = 2;
        canvas.height = 2;
        canvas.getContext("2d")?.fillRect(0, 0, 2, 2);
        wakeVideo.srcObject = canvas.captureStream(1);
      }
      await wakeVideo?.play();
      note("silent video: playing (wake-lock fallback)");
    } catch {
      note("silent video: refused");
    }
    pushSentinel();
    started = true;
  }

  async function resumeFullscreen() {
    try {
      await document.documentElement.requestFullscreen({ navigationUI: "hide" });
    } catch (error) {
      note(`full screen: refused (${(error as Error).name})`);
    }
  }

  beforeNavigate((navigation) => {
    if (guardBack && navigation.type === "popstate") {
      navigation.cancel();
      backCaught += 1;
      note("back gesture: cancelled by the router");
      pushSentinel();
    }
  });

  onMount(() => {
    secureContext = window.isSecureContext;
    wakeLockSupported = "wakeLock" in navigator;
    standalone = window.matchMedia("(display-mode: standalone), (display-mode: fullscreen)").matches;
    document.documentElement.classList.add("kiosk");
    note(`secure context: ${secureContext ? "yes" : "no"}`);

    const onFullscreen = () => {
      isFullscreen = document.fullscreenElement !== null;
      note(`full screen: ${isFullscreen ? "on" : "off"}`);
    };
    const onVisibility = () => {
      note(`page: ${document.visibilityState}`);
      if (document.visibilityState === "visible" && started && !wakeLockActive) takeWakeLock();
    };
    const onContextMenu = (event: Event) => {
      if (!blockMenu) return;
      event.preventDefault();
      menuBlocked += 1;
      note("long press: menu blocked");
    };
    const onUnload = (event: BeforeUnloadEvent) => {
      if (!confirmUnload) return;
      event.preventDefault();
      note("unload: confirm asked");
    };

    document.addEventListener("fullscreenchange", onFullscreen);
    document.addEventListener("visibilitychange", onVisibility);
    document.addEventListener("contextmenu", onContextMenu);
    window.addEventListener("beforeunload", onUnload);

    return () => {
      document.documentElement.classList.remove("kiosk");
      document.removeEventListener("fullscreenchange", onFullscreen);
      document.removeEventListener("visibilitychange", onVisibility);
      document.removeEventListener("contextmenu", onContextMenu);
      window.removeEventListener("beforeunload", onUnload);
      lock?.release();
    };
  });
</script>

<svelte:head><title>Kiosk prototype</title></svelte:head>

<video bind:this={wakeVideo} class="wake" muted loop playsinline></video>

{#if !started}
  <button class="start" onclick={startShift}>Start shift</button>
{:else}
  <main>
    <h1>Kiosk prototype</h1>
    {#if !isFullscreen}
      <button class="resume" onclick={resumeFullscreen}>Resume full screen</button>
    {/if}

    <dl>
      <div><dt>Secure context</dt><dd class:bad={!secureContext}>{secureContext ? "yes" : "no (HTTP)"}</dd></div>
      <div><dt>Wake lock API</dt><dd class:bad={!wakeLockSupported}>{wakeLockSupported ? "present" : "missing"}</dd></div>
      <div><dt>Wake lock held</dt><dd class:bad={!wakeLockActive}>{wakeLockActive ? "yes" : "no"}</dd></div>
      <div><dt>Full screen</dt><dd class:bad={!isFullscreen}>{isFullscreen ? "yes" : "no"}</dd></div>
      <div><dt>Installed look</dt><dd>{standalone ? "standalone" : "browser tab"}</dd></div>
      <div><dt>Back gestures caught</dt><dd>{backCaught}</dd></div>
      <div><dt>Long presses blocked</dt><dd>{menuBlocked}</dd></div>
      <div><dt>History sentinels</dt><dd>{entries}</dd></div>
    </dl>

    <p class="hint">Try: pull the notification shade, swipe back, long press this text, pull down to refresh, then leave the tablet alone for 10 minutes.</p>

    <div class="toggles">
      <label><input type="checkbox" bind:checked={guardBack} /> Guard the back gesture</label>
      <label><input type="checkbox" bind:checked={blockMenu} /> Block the long-press menu</label>
      <label><input type="checkbox" bind:checked={confirmUnload} /> Confirm before unload</label>
    </div>

    <p class="hint">Selection test: try to select this sentence with a long press and drag.</p>
    <p class="hint"><a href="/kiosk-prototype?again=1">Client-side link</a> - then swipe back.</p>

    <ul class="log">{#each log as line (line.id)}<li>{line.text}</li>{/each}</ul>
  </main>
{/if}

<style>
  :global(html.kiosk),
  :global(html.kiosk body) {
    overscroll-behavior: none;
    touch-action: manipulation;
    user-select: none;
    -webkit-user-select: none;
    -webkit-touch-callout: none;
  }

  :global(html.kiosk input) {
    user-select: text;
    -webkit-user-select: text;
  }

  .wake {
    position: fixed;
    width: 1px;
    height: 1px;
    opacity: 0;
    pointer-events: none;
  }

  .start {
    position: fixed;
    inset: 0;
    width: 100%;
    height: 100%;
    border: 0;
    font-size: 3rem;
    font-weight: 900;
    color: white;
    background: var(--apple-red, #b91c1c);
  }

  main {
    padding: 1.5rem;
    font-size: 1.1rem;
  }

  h1 {
    margin: 0 0 1rem;
  }

  .resume {
    display: block;
    width: 100%;
    min-height: 4rem;
    margin-bottom: 1rem;
    border: 0;
    border-radius: 1rem;
    font-size: 1.5rem;
    font-weight: 900;
    color: white;
    background: var(--leaf-green, #2f6b2f);
  }

  dl {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(14rem, 1fr));
    gap: 0.5rem;
    margin: 0 0 1rem;
  }

  dl div {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    padding: 0.6rem 0.9rem;
    border-radius: 0.75rem;
    background: #ffffffcc;
  }

  dt {
    font-weight: 700;
  }

  dd {
    margin: 0;
    font-variant-numeric: tabular-nums;
  }

  dd.bad {
    color: #b91c1c;
    font-weight: 900;
  }

  .toggles {
    display: grid;
    gap: 0.5rem;
    margin: 1rem 0;
  }

  .toggles label {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .toggles input {
    width: 1.5rem;
    height: 1.5rem;
  }

  .hint {
    margin: 0.5rem 0;
  }

  .log {
    max-height: 40vh;
    overflow: auto;
    margin: 1rem 0 0;
    padding: 0.75rem 0.75rem 0.75rem 1.5rem;
    border-radius: 0.75rem;
    background: #ffffffcc;
    font-family: ui-monospace, monospace;
    font-size: 0.9rem;
  }
</style>
