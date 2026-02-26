<script lang="ts">
  import type { APIChangeset } from './DiffTable.svelte';

  let {
    changeset,
    collapsed = false,
    onToggle
  }: {
    changeset: APIChangeset;
    collapsed?: boolean;
    onToggle?: () => void;
  } = $props();

  function fileIcon(path: string): string {
    if (path.endsWith('.go')) return 'fa-file-code-o';
    if (path.endsWith('.ts') || path.endsWith('.tsx')) return 'fa-file-code-o';
    if (path.endsWith('.js') || path.endsWith('.jsx')) return 'fa-file-code-o';
    if (path.endsWith('.svelte')) return 'fa-file-code-o';
    if (path.endsWith('.css') || path.endsWith('.scss')) return 'fa-css3';
    if (path.endsWith('.html')) return 'fa-html5';
    if (path.endsWith('.json')) return 'fa-file-text-o';
    if (path.endsWith('.md')) return 'fa-file-text-o';
    if (/\.(png|jpg|jpeg|gif|svg|webp)$/.test(path)) return 'fa-file-image-o';
    return 'fa-file-o';
  }
</script>

<div class="mood-changeset-header" role="button" tabindex="0" onclick={onToggle} onkeydown={(e) => e.key === 'Enter' && onToggle?.()}>
  <span class="phui-icon-view phui-font-fa {collapsed ? 'fa-chevron-right' : 'fa-chevron-down'} changeset-collapse-toggle"></span>
  <span class="phui-icon-view phui-font-fa {fileIcon(changeset.displayPath)}" style="opacity:0.5"></span>
  <span class="differential-changeset-path-name">{changeset.displayPath}</span>
  <span class="stats">
    {#if changeset.linesAdded > 0}
      <span class="add-stat">+{changeset.linesAdded}</span>
    {/if}
    {#if changeset.linesAdded > 0 && changeset.linesRemoved > 0}
      {' '}
    {/if}
    {#if changeset.linesRemoved > 0}
      <span class="del-stat">-{changeset.linesRemoved}</span>
    {/if}
  </span>
</div>

<style>
  .mood-changeset-header {
    background: #eceef4;
    padding: 8px 12px;
    border-bottom: 1px solid #c7ccd9;
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    font-weight: 600;
    color: #292e36;
    position: sticky;
    top: 0;
    z-index: 10;
    cursor: pointer;
    user-select: none;
  }

  .stats {
    font-weight: 400;
    color: #6b748c;
    margin-left: auto;
  }

  .add-stat {
    color: #139543;
  }

  .del-stat {
    color: #c0392b;
  }

  .changeset-collapse-toggle {
    font-size: 12px;
    opacity: 0.6;
  }

  /* Dark mode */
  :global(.phui-theme-dark) .mood-changeset-header {
    background: rgba(255, 255, 255, 0.1);
    border-color: rgba(255, 255, 255, 0.3);
    color: rgba(255, 255, 255, 0.9);
  }
  :global(.phui-theme-dark) .stats {
    color: rgba(255, 255, 255, 0.6);
  }
</style>
