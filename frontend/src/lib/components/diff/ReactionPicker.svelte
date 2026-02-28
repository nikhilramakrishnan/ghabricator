<script lang="ts">
  let {
    onPick,
    onClose
  }: {
    onPick: (emoji: string) => void;
    onClose: () => void;
  } = $props();

  const reactions = [
    { emoji: '+1', icon: 'fa-thumbs-up', label: 'Thumbs up' },
    { emoji: '-1', icon: 'fa-thumbs-down', label: 'Thumbs down' },
    { emoji: 'laugh', icon: 'fa-smile-o', label: 'Laugh' },
    { emoji: 'confused', icon: 'fa-question', label: 'Confused' },
    { emoji: 'heart', icon: 'fa-heart', label: 'Heart' },
    { emoji: 'star', icon: 'fa-star', label: 'Star' },
    { emoji: 'rocket', icon: 'fa-rocket', label: 'Rocket' },
    { emoji: 'eyes', icon: 'fa-eye', label: 'Eyes' }
  ];

  function handleClick(emoji: string) {
    onPick(emoji);
    onClose();
  }

  function handleClickOutside(e: MouseEvent) {
    const target = e.target as HTMLElement;
    if (!target.closest('.reaction-picker')) {
      onClose();
    }
  }

  function stopProp(e: MouseEvent) {
    e.stopPropagation();
  }
</script>

<svelte:window onclick={handleClickOutside} />

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="reaction-picker" onclick={stopProp}>
  {#each reactions as r}
    <button class="rp-btn" title={r.label} onclick={() => handleClick(r.emoji)}>
      <i class="fa {r.icon}"></i>
    </button>
  {/each}
</div>

<style>
  .reaction-picker {
    position: absolute;
    bottom: calc(100% + 4px);
    left: 0;
    display: flex;
    gap: 2px;
    padding: 4px 6px;
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 6px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    z-index: 10;
    white-space: nowrap;
  }

  .rp-btn {
    all: unset;
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    color: var(--text-muted);
    transition: background 0.1s, color 0.1s;
  }
  .rp-btn:hover {
    background: var(--bg-hover);
    color: var(--text-link);
  }
</style>
