<script>
  import { onDestroy } from 'svelte';
  import { Browser, Events } from '@wailsio/runtime';
  import { DashboardService } from '../../bindings/github.com/ao-data/albiondata-client/internal/dashboard/index.js';
  import CountersPanel from './CountersPanel.svelte';

  const WIDTH_KEY = 'sidebar-width';
  const MIN_WIDTH = 168;
  const MAX_WIDTH = 360;
  const DEFAULT_WIDTH = 208;

  function loadWidth() {
    const stored = Number(localStorage.getItem(WIDTH_KEY));
    if (Number.isFinite(stored) && stored >= MIN_WIDTH && stored <= MAX_WIDTH) return stored;
    return DEFAULT_WIDTH;
  }

  let width = $state(loadWidth());
  let resizing = $state(false);

  function startResize(e) {
    e.preventDefault();
    resizing = true;
    window.addEventListener('pointermove', onResize);
    window.addEventListener('pointerup', stopResize);
  }

  function onResize(e) {
    width = Math.min(MAX_WIDTH, Math.max(MIN_WIDTH, e.clientX));
  }

  function stopResize() {
    resizing = false;
    localStorage.setItem(WIDTH_KEY, String(width));
    window.removeEventListener('pointermove', onResize);
    window.removeEventListener('pointerup', stopResize);
  }

  onDestroy(() => {
    window.removeEventListener('pointermove', onResize);
    window.removeEventListener('pointerup', stopResize);
  });

  let status = $state({
    Version: '',
    UpdateAvailable: '',
    CaptureRunning: false,
    CaptureError: false,
    ServerID: 0,
    IngestBaseURL: '',
    CustomPublicIngest: false,
    DriverWarning: '',
    DriverHelpURL: '',
    EncryptionStatus: '',
  });

  DashboardService.GetStatus().then((s) => (status = s));

  const unlisten = Events.On('status:changed', (evt) => {
    status = evt.data;
  });

  onDestroy(unlisten);

  const serverNames = { 0: 'Unknown', 1: 'West', 2: 'East', 3: 'Europe' };
  let serverLabel = $derived(
    status.CustomPublicIngest ? 'Private' : (serverNames[status.ServerID] ?? status.ServerID)
  );
  let badgeLabel = $derived(
    status.CaptureError ? 'Error' : status.CaptureRunning ? 'Capturing' : 'Ready'
  );
  let encryptionLabel = $derived(
    status.EncryptionStatus === 'encrypted'
      ? 'Encrypted'
      : status.EncryptionStatus === 'clear'
        ? 'Not Encrypted'
        : 'Encrypted?'
  );

  let startupSupported = $state(false);
  let startupEnabled = $state(false);
  let startupBusy = $state(false);
  let startupError = $state('');

  DashboardService.StartupSupported().then((supported) => {
    startupSupported = supported;
    if (supported) {
      DashboardService.GetStartupEnabled()
        .then((enabled) => (startupEnabled = enabled))
        .catch((err) => (startupError = String(err)));
    }
  });

  async function toggleStartup() {
    if (startupBusy) return;
    startupBusy = true;
    startupError = '';
    try {
      startupEnabled = await DashboardService.SetStartupEnabled(!startupEnabled);
    } catch (err) {
      startupError = String(err);
    } finally {
      startupBusy = false;
    }
  }

  function openDriverHelp(e) {
    e.preventDefault();
    Browser.OpenURL(status.DriverHelpURL);
  }
</script>

<aside class="sidebar" style="width: {width}px">
  <div
    class="resize-handle"
    class:active={resizing}
    onpointerdown={startResize}
    role="separator"
    aria-orientation="vertical"
    aria-label="Resize sidebar"
  ></div>
  <div class="brand">
    <span class="wordmark">Albion Data Client</span>
  </div>

  <div class="group">
    <div class="pill-row">
      <span
        class="status-pill"
        class:running={status.CaptureRunning && !status.CaptureError}
        class:error={status.CaptureError}
      >
        <span class="lantern"><span class="glow"></span></span>
        {badgeLabel}
      </span>

      <span
        class="status-pill"
        class:encrypted={status.EncryptionStatus === 'encrypted'}
        class:clear={status.EncryptionStatus === 'clear'}
      >
        <span class="lantern"><span class="glow"></span></span>
        {encryptionLabel}
      </span>
    </div>

    <div class="field">
      <span class="label">Server</span>
      <span class="value">{serverLabel}</span>
    </div>
  </div>

  {#if status.DriverWarning}
    <div class="group driver-warning-group">
      <div class="driver-warning">
        <span class="driver-warning-text">{status.DriverWarning}</span>
        {#if status.DriverHelpURL}
          <a class="driver-warning-link" href={status.DriverHelpURL} onclick={openDriverHelp}>
            Get Npcap
          </a>
        {/if}
      </div>
    </div>
  {/if}

  <CountersPanel />

  {#if startupSupported}
    <div class="group startup-group">
      <div class="startup-row">
        <div class="startup-copy">
          <span class="label">Windows startup</span>
          <span class="startup-help">Start in the tray when you sign in</span>
        </div>
        <button
          type="button"
          class="startup-toggle"
          class:enabled={startupEnabled}
          aria-pressed={startupEnabled}
          aria-label="Start Albion Data Client with Windows"
          disabled={startupBusy}
          onclick={toggleStartup}
        >{startupBusy ? '...' : startupEnabled ? 'On' : 'Off'}</button>
      </div>
      {#if startupError}
        <p class="startup-error">{startupError}</p>
      {/if}
    </div>
  {/if}

  {#if status.UpdateAvailable}
    <div class="group update-group">
      <span class="update-badge">
        <span class="label">Update available</span>
        <span class="value">{status.UpdateAvailable}</span>
      </span>
    </div>
  {/if}
</aside>

<style>
  .sidebar {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 1.75rem;
    flex-shrink: 0;
    padding: 1.5rem 1.4rem 1rem;
    background: var(--bg-raised);
    border-right: 1px solid var(--border);
    overflow-y: auto;
    overflow-x: hidden;
  }
  .resize-handle {
    position: absolute;
    top: 0;
    right: -3px;
    width: 6px;
    height: 100%;
    cursor: col-resize;
    touch-action: none;
    z-index: 1;
  }
  .resize-handle::after {
    content: '';
    position: absolute;
    top: 0;
    left: 2px;
    width: 2px;
    height: 100%;
    background: transparent;
    transition: background-color 0.15s ease;
  }
  .resize-handle:hover::after,
  .resize-handle.active::after {
    background: var(--blue);
  }
  .brand {
    display: flex;
    min-width: 0;
  }
  .wordmark {
    font-family: var(--font-display);
    font-size: 1.05rem;
    font-weight: 600;
    letter-spacing: 0.01em;
    color: var(--text);
    line-height: 1.2;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .group {
    display: flex;
    flex-direction: column;
    gap: 0.9rem;
    padding-top: 1.15rem;
    border-top: 1px solid var(--border);
  }
  .field {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }
  .startup-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
  }
  .startup-copy {
    display: flex;
    flex-direction: column;
    gap: 0.3rem;
    min-width: 0;
  }
  .startup-help {
    color: var(--text-muted);
    font-size: 0.73rem;
    line-height: 1.35;
  }
  .startup-toggle {
    flex-shrink: 0;
    min-width: 3.1rem;
    padding: 0.35rem 0.5rem;
    border: 1px solid var(--border-strong);
    border-radius: var(--radius-sm);
    background: var(--bg-sunken);
    color: var(--text-muted);
    font: 600 0.73rem var(--font-mono);
    cursor: pointer;
  }
  .startup-toggle.enabled {
    border-color: var(--blue);
    background: var(--blue-soft);
    color: var(--blue-bright);
  }
  .startup-toggle:disabled {
    opacity: 0.6;
    cursor: wait;
  }
  .startup-error {
    margin: 0;
    color: var(--ember);
    font-size: 0.72rem;
    line-height: 1.4;
    overflow-wrap: anywhere;
  }
  .label {
    font-size: 0.66rem;
    font-weight: 600;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: var(--text-faint);
  }
  .value {
    font-family: var(--font-mono);
    font-size: 0.85rem;
    color: var(--text);
  }
  .pill-row {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 0.5rem;
  }
  .status-pill {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    align-self: flex-start;
    padding: 0.3rem 0.7rem 0.3rem 0.55rem;
    border-radius: 999px;
    font-size: 0.78rem;
    font-weight: 500;
    background: var(--ready-soft);
    color: var(--ready);
  }
  .lantern {
    position: relative;
    display: inline-flex;
    width: 7px;
    height: 7px;
  }
  .lantern::before {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: 50%;
    background: currentColor;
  }
  .lantern .glow {
    position: absolute;
    inset: -4px;
    border-radius: 50%;
    background: currentColor;
    opacity: 0;
  }
  .status-pill.running {
    background: var(--blue-soft);
    color: var(--blue-bright);
  }
  .status-pill.running .glow {
    animation: breathe 2.4s ease-in-out infinite;
  }
  .status-pill.error,
  .status-pill.encrypted {
    background: var(--ember-soft);
    color: var(--ember);
  }
  .status-pill.clear {
    background: var(--moss-soft);
    color: var(--moss);
  }
  @keyframes breathe {
    0%, 100% { opacity: 0; transform: scale(0.6); }
    50% { opacity: 0.45; transform: scale(1.4); }
  }
  @media (prefers-reduced-motion: reduce) {
    .status-pill.running .glow {
      animation: none;
    }
  }
  .driver-warning-group {
    border-top: 1px solid var(--border);
  }
  .driver-warning {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    padding: 0.6rem 0.7rem;
    border-radius: var(--radius-sm);
    background: var(--blue-soft);
    border: 1px solid rgba(47, 140, 255, 0.3);
  }
  .driver-warning-text {
    font-size: 0.78rem;
    line-height: 1.5;
    color: var(--blue);
  }
  .driver-warning-link {
    align-self: flex-start;
    font-size: 0.72rem;
    font-weight: 600;
    letter-spacing: 0.02em;
    text-transform: uppercase;
    color: var(--blue);
    text-decoration: none;
    border-bottom: 1px solid rgba(47, 140, 255, 0.5);
  }
  .driver-warning-link:hover {
    border-bottom-color: var(--blue);
  }
  .update-group {
    border-top: 1px solid var(--border);
  }
  .update-badge {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
    padding: 0.6rem 0.7rem;
    border-radius: var(--radius-sm);
    background: var(--blue-soft);
    border: 1px solid rgba(47, 140, 255, 0.25);
  }
  .update-badge .label {
    color: var(--blue);
    opacity: 0.9;
  }
  .update-badge .value {
    font-family: var(--font-body);
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--text);
  }
</style>
