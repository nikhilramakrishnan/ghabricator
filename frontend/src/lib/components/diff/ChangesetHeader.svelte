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

<div class="changeset-header" role="button" tabindex="0" onclick={onToggle} onkeydown={(e) => e.key === 'Enter' && onToggle?.()}>
  <i class="fa {collapsed ? 'fa-chevron-right' : 'fa-chevron-down'} toggle-icon"></i>
  <i class="fa {fileIcon(changeset.displayPath)} file-icon"></i>
  <span class="path-name">{changeset.displayPath}</span>
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
  .changeset-header {
    background: var(--bg-card-header);
    padding: 8px 12px;
    border-bottom: 1px solid var(--border);
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    font-weight: 600;
    color: var(--text);
    position: sticky;
    top: 0;
    z-index: 10;
    cursor: pointer;
    user-select: none;
  }

  .toggle-icon {
    font-size: 12px;
    opacity: 0.6;
  }

  .file-icon {
    opacity: 0.5;
  }

  .stats {
    font-weight: 400;
    color: var(--text-muted);
    margin-left: auto;
  }

  .add-stat {
    color: var(--green);
  }

  .del-stat {
    color: var(--red);
  }
</style>
