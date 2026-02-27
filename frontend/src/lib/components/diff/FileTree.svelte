<script lang="ts">
  import type { APIChangeset } from './DiffTable.svelte';

  let {
    changesets,
    activeFile = ''
  }: {
    changesets: APIChangeset[];
    activeFile?: string;
  } = $props();

  let filter = $state('');

  let filtered = $derived(
    filter
      ? changesets.filter((cs) =>
          cs.displayPath.toLowerCase().includes(filter.toLowerCase())
        )
      : changesets
  );

  function splitPath(path: string): { dir: string; file: string } {
    const idx = path.lastIndexOf('/');
    if (idx < 0) return { dir: '', file: path };
    return { dir: path.slice(0, idx + 1), file: path.slice(idx + 1) };
  }

  function statusColor(cs: APIChangeset): string {
    if (cs.isNew) return 'var(--green)';
    if (cs.isDeleted) return 'var(--red)';
    return 'var(--orange)';
  }
</script>

<div class="file-tree">
  {#if changesets.length > 5}
    <div class="filter-box">
      <input
        type="text"
        placeholder="Filter files..."
        bind:value={filter}
      />
    </div>
  {/if}

  {#each filtered as cs}
    {@const { dir, file } = splitPath(cs.displayPath)}
    <div
      class="tree-item"
      class:active={cs.displayPath === activeFile}
    >
      <div class="tree-row">
        <div class="tree-name">
          <a href="#C{cs.id}">
            <span class="status-dot" style="color:{statusColor(cs)}">&#x25CF;</span>
            {#if dir}
              <span class="dir-prefix">{dir}</span>
            {/if}
            {file}
          </a>
        </div>
        {#if cs.linesAdded > 0 || cs.linesRemoved > 0}
          <div class="tree-stats">
            {#if cs.linesAdded > 0}
              <span class="stat-add">+{cs.linesAdded}</span>
            {/if}
            {#if cs.linesRemoved > 0}
              <span class="stat-del">-{cs.linesRemoved}</span>
            {/if}
          </div>
        {/if}
      </div>
    </div>
  {/each}
</div>

<style>
  .file-tree {
    font-size: 12px;
  }

  .filter-box {
    padding: 8px;
  }

  .filter-box input {
    width: 100%;
    padding: 4px 8px;
    border: 1px solid var(--border);
    border-radius: 3px;
    font-size: 12px;
    outline: none;
    box-sizing: border-box;
    background: var(--bg-card);
    color: var(--text);
  }
  .filter-box input:focus {
    border-color: var(--text-link);
  }

  .tree-item {
    padding: 4px 8px;
    cursor: pointer;
  }
  .tree-item:hover {
    background: var(--bg-hover);
  }
  .tree-item.active {
    background: rgba(19, 108, 178, 0.08);
  }

  .tree-row {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .tree-name {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .tree-name a {
    text-decoration: none;
    color: var(--text);
  }

  .status-dot {
    margin-right: 4px;
  }

  .dir-prefix {
    opacity: 0.65;
  }

  .tree-stats {
    display: flex;
    gap: 4px;
    flex-shrink: 0;
  }

  .stat-add {
    color: var(--green);
    background: rgba(46, 168, 107, 0.1);
    padding: 0 4px;
    border-radius: 2px;
    font-size: 11px;
  }

  .stat-del {
    color: var(--red);
    background: rgba(192, 57, 43, 0.1);
    padding: 0 4px;
    border-radius: 2px;
    font-size: 11px;
  }
</style>
