<script lang="ts">
  export type TimelineEvent = {
    author: { login: string; avatarURL?: string };
    action: string;
    body?: string;
    createdAt: string;
    iconClass: string;
    iconColor: string;
  };

  let { events }: { events: TimelineEvent[] } = $props();

  function iconBg(color: string): string {
    switch (color) {
      case 'green': return '#139543';
      case 'red': return '#c0392b';
      case 'blue': return '#136cb2';
      case 'violet': return '#6e5494';
      default: return '#6b748c';
    }
  }
</script>

<div class="phui-box phui-box-border phui-object-box" style="padding:10px 12px; margin:12px 0 0 0;">
  <div class="phui-timeline-view">
    {#each events as ev, i}
      {@const isMajor = !!ev.body}
      {@const isLast = i === events.length - 1}
      <div
        class="phui-timeline-event-view"
        style="display:flex; gap:10px; padding:10px 0;{isLast ? '' : ' border-bottom:1px solid #e3e4e8;'}"
      >
        {#if isMajor && ev.author.avatarURL}
          <img
            src={ev.author.avatarURL}
            alt=""
            style="width:32px; height:32px; border-radius:50%; flex-shrink:0;"
          />
        {:else}
          <div style="width:32px; height:32px; border-radius:50%; background:{iconBg(ev.iconColor)}; display:flex; align-items:center; justify-content:center; flex-shrink:0;">
            <i class="fa {ev.iconClass}" style="color:#fff; font-size:14px;"></i>
          </div>
        {/if}

        <div style="{isMajor ? 'flex:1;' : ''}min-width:0;">
          <div class="phui-timeline-title" style="font-size:13px; display:flex; align-items:baseline; gap:8px;">
            <span><strong>{ev.author.login}</strong> {ev.action}</span>
            <span style="font-size:12px; color:#6b748c; margin-left:auto; white-space:nowrap;">{ev.createdAt}</span>
          </div>
          {#if isMajor}
            <div style="background:#f6f8fa; border:1px solid #e3e4e8; border-radius:4px; padding:12px; font-size:13px; line-height:1.5; overflow-wrap:break-word; word-break:break-word; overflow-x:auto; max-width:100%; margin-top:8px;">
              {@html ev.body}
            </div>
          {/if}
        </div>
      </div>
    {/each}
  </div>
</div>
