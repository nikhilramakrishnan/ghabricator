<script lang="ts">
  export type TimelineEvent = {
    author: { login: string; avatarURL?: string };
    action: string;
    body?: string;
    createdAt: string;
    iconClass: string;
    iconColor: string;
  };

  import { formatTimestamp } from '$lib/time';

  let { events }: { events: TimelineEvent[] } = $props();

  function iconBg(color: string): string {
    switch (color) {
      case 'green': return 'var(--green)';
      case 'red': return 'var(--red)';
      case 'blue': return 'var(--blue)';
      case 'violet': return 'var(--violet)';
      default: return 'var(--text-muted)';
    }
  }
</script>

<div class="timeline-box">
  {#each events as ev, i}
    {@const isMajor = !!ev.body}
    {@const isLast = i === events.length - 1}
    <div class="event" class:last={isLast}>
      {#if isMajor && ev.author.avatarURL}
        <img
          src={ev.author.avatarURL}
          alt=""
          class="event-avatar"
        />
      {:else}
        <div class="event-icon" style="background:{iconBg(ev.iconColor)}">
          <i class="fa {ev.iconClass}"></i>
        </div>
      {/if}

      <div class="event-body" class:major={isMajor}>
        <div class="event-header">
          <span><strong>{ev.author.login}</strong> {ev.action}</span>
          <span class="event-time">{formatTimestamp(ev.createdAt)}</span>
        </div>
        {#if isMajor}
          <div class="event-content">
            {@html ev.body}
          </div>
        {/if}
      </div>
    </div>
  {/each}
</div>

<style>
  .timeline-box {
    padding: 10px 12px;
  }
  .event {
    display: flex;
    gap: 10px;
    padding: 10px 0;
    border-bottom: 1px solid var(--border-subtle);
  }
  .event.last {
    border-bottom: none;
  }
  .event-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    flex-shrink: 0;
  }
  .event-icon {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    color: #fff;
    font-size: 14px;
  }
  .event-body {
    min-width: 0;
  }
  .event-body.major {
    flex: 1;
  }
  .event-header {
    font-size: 13px;
    display: flex;
    align-items: baseline;
    gap: 8px;
  }
  .event-time {
    font-size: 12px;
    color: var(--text-muted);
    margin-left: auto;
    white-space: nowrap;
  }
  .event-content {
    background: var(--bg-subtle);
    border: 1px solid var(--border-subtle);
    border-radius: 4px;
    padding: 12px;
    font-size: 13px;
    line-height: 1.5;
    overflow-wrap: break-word;
    word-break: break-word;
    overflow-x: auto;
    max-width: 100%;
    margin-top: 8px;
  }
</style>
