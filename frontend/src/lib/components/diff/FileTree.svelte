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
    if (cs.isNew) return '#2EA86B';
    if (cs.isDeleted) return '#C0392B';
    return '#E8A83E';
  }
</script>

<div class="diff-tree-view">
  {#if changesets.length > 5}
    <div class="diff-tree-filter">
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
      class="diff-tree-path diff-tree-path-changeset"
      class:active={cs.displayPath === activeFile}
    >
      <div class="diff-tree-path-indent">
        <div class="diff-tree-path-name">
          <a href="#C{cs.id}">
            <span class="status-dot" style="color:{statusColor(cs)}">&#x25CF;</span>
            {#if dir}
              <span class="dir-prefix">{dir}</span>
            {/if}
            {file}
          </a>
        </div>
        {#if cs.linesAdded > 0 || cs.linesRemoved > 0}
          <div class="diff-tree-path-inlines">
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
  .diff-tree-view {
    font-size: 12px;
  }

  .diff-tree-filter {
    padding: 8px;
  }

  .diff-tree-filter input {
    width: 100%;
    padding: 4px 8px;
    border: 1px solid #c7ccd9;
    border-radius: 3px;
    font-size: 12px;
    outline: none;
    box-sizing: border-box;
  }
  .diff-tree-filter input:focus {
    border-color: #136cb2;
  }

  .diff-tree-path {
    padding: 4px 8px;
    cursor: pointer;
  }
  .diff-tree-path:hover {
    background: rgba(55, 55, 55, 0.04);
  }
  .diff-tree-path.active {
    background: rgba(19, 108, 178, 0.08);
  }

  .diff-tree-path-indent {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .diff-tree-path-name {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .diff-tree-path-name a {
    text-decoration: none;
    color: inherit;
  }

  .status-dot {
    margin-right: 4px;
  }

  .dir-prefix {
    opacity: 0.65;
  }

  .diff-tree-path-inlines {
    display: flex;
    gap: 4px;
    flex-shrink: 0;
  }

  .stat-add {
    color: #2ea86b;
    background: rgba(46, 168, 107, 0.1);
    padding: 0 4px;
    border-radius: 2px;
    font-size: 11px;
  }

  .stat-del {
    color: #c0392b;
    background: rgba(192, 57, 43, 0.1);
    padding: 0 4px;
    border-radius: 2px;
    font-size: 11px;
  }

  /* Dark mode */
  :global(.phui-theme-dark) .diff-tree-filter input {
    background: #1b2028;
    color: #c8d1db;
    border-color: #464c5c;
  }
  :global(.phui-theme-dark) .diff-tree-path:hover {
    background: rgba(255, 255, 255, 0.05);
  }
  :global(.phui-theme-dark) .diff-tree-path.active {
    background: rgba(255, 255, 255, 0.1);
  }
  :global(.phui-theme-dark) .diff-tree-path-name a {
    color: #c8d1db;
  }
</style>
