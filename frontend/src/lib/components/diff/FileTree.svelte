<script lang="ts">
  import type { APIChangeset } from './DiffTable.svelte';

  type TreeNode = {
    name: string;
    path: string;       // full path prefix for dirs, full path for files
    isDir: boolean;
    children: TreeNode[];
    changeset?: APIChangeset;
    totalAdded: number;
    totalRemoved: number;
    commentCount: number;
  };

  let {
    changesets,
    activeFile = '',
    commentCounts = {}
  }: {
    changesets: APIChangeset[];
    activeFile?: string;
    commentCounts?: Record<string, number>;
  } = $props();

  let filter = $state('');
  let collapsedDirs = $state(new Set<string>());
  let viewedFiles = $state(new Set<string>());

  let filtered = $derived(
    filter
      ? changesets.filter((cs) =>
          cs.displayPath.toLowerCase().includes(filter.toLowerCase())
        )
      : changesets
  );

  // Build tree structure from flat file paths
  let tree: TreeNode[] = $derived.by(() => {
    const root: TreeNode[] = [];

    for (const cs of filtered) {
      const parts = cs.displayPath.split('/');
      let current = root;

      // Build directory nodes
      for (let i = 0; i < parts.length - 1; i++) {
        const dirName = parts[i];
        const dirPath = parts.slice(0, i + 1).join('/');
        let dirNode = current.find(n => n.isDir && n.name === dirName);
        if (!dirNode) {
          dirNode = {
            name: dirName,
            path: dirPath,
            isDir: true,
            children: [],
            totalAdded: 0,
            totalRemoved: 0,
            commentCount: 0
          };
          current.push(dirNode);
        }
        dirNode.totalAdded += cs.linesAdded;
        dirNode.totalRemoved += cs.linesRemoved;
        dirNode.commentCount += (commentCounts[cs.displayPath] ?? 0);
        current = dirNode.children;
      }

      // Add file node
      const fileName = parts[parts.length - 1];
      current.push({
        name: fileName,
        path: cs.displayPath,
        isDir: false,
        children: [],
        changeset: cs,
        totalAdded: cs.linesAdded,
        totalRemoved: cs.linesRemoved,
        commentCount: commentCounts[cs.displayPath] ?? 0
      });
    }

    // Sort: dirs first, then alphabetical
    function sortNodes(nodes: TreeNode[]) {
      nodes.sort((a, b) => {
        if (a.isDir !== b.isDir) return a.isDir ? -1 : 1;
        return a.name.localeCompare(b.name);
      });
      for (const n of nodes) {
        if (n.isDir) sortNodes(n.children);
      }
    }
    sortNodes(root);

    // Auto-collapse single-child dir chains (compact paths like "src/services/")
    function compact(nodes: TreeNode[]): TreeNode[] {
      return nodes.map(n => {
        if (n.isDir && n.children.length === 1 && n.children[0].isDir) {
          const child = n.children[0];
          return {
            ...child,
            name: n.name + '/' + child.name,
            path: child.path,
            children: compact(child.children)
          };
        }
        if (n.isDir) {
          return { ...n, children: compact(n.children) };
        }
        return n;
      });
    }

    return compact(root);
  });

  function toggleDir(path: string) {
    const next = new Set(collapsedDirs);
    if (next.has(path)) next.delete(path);
    else next.add(path);
    collapsedDirs = next;
  }

  function toggleViewed(path: string, e: Event) {
    e.preventDefault();
    e.stopPropagation();
    const next = new Set(viewedFiles);
    if (next.has(path)) next.delete(path);
    else next.add(path);
    viewedFiles = next;
  }

  function statusColor(cs: APIChangeset): string {
    if (cs.isNew) return 'var(--green)';
    if (cs.isDeleted) return 'var(--red)';
    return 'var(--blue)';
  }

  function fileIcon(cs: APIChangeset): string {
    if (cs.isNew) return 'fa-plus-circle';
    if (cs.isDeleted) return 'fa-minus-circle';
    return 'fa-pencil';
  }
</script>

<div class="file-tree">
  {#if changesets.length > 8}
    <div class="filter-box">
      <i class="fa fa-search filter-icon"></i>
      <input
        type="text"
        placeholder="Filter files..."
        bind:value={filter}
      />
    </div>
  {/if}

  <div class="tree-nodes">
    {#snippet renderNodes(nodes: TreeNode[], depth: number)}
      {#each nodes as node}
        {#if node.isDir}
          {@const collapsed = collapsedDirs.has(node.path)}
          <button
            class="tree-dir"
            style="padding-left:{8 + depth * 12}px"
            onclick={() => toggleDir(node.path)}
          >
            <i class="fa {collapsed ? 'fa-caret-right' : 'fa-caret-down'} caret"></i>
            <i class="fa fa-folder{collapsed ? '' : '-open'} folder-icon"></i>
            <span class="dir-name">{node.name}</span>
            {#if node.commentCount > 0}
              <span class="badge comment-badge" title="{node.commentCount} comments">
                <i class="fa fa-comment-o"></i> {node.commentCount}
              </span>
            {/if}
          </button>
          {#if !collapsed}
            {@render renderNodes(node.children, depth + 1)}
          {/if}
        {:else if node.changeset}
          {@const cs = node.changeset}
          {@const viewed = viewedFiles.has(cs.displayPath)}
          <a
            href="#C{cs.id}"
            class="tree-file"
            class:active={cs.displayPath === activeFile}
            class:viewed
            style="padding-left:{8 + depth * 12}px"
          >
            <i class="fa {fileIcon(cs)} file-status-icon" style="color:{statusColor(cs)}"></i>
            <span class="file-name" class:viewed>{node.name}</span>
            <span class="file-meta">
              {#if node.commentCount > 0}
                <span class="badge comment-badge" title="{node.commentCount} comments">
                  <i class="fa fa-comment-o"></i> {node.commentCount}
                </span>
              {/if}
              {#if cs.linesAdded > 0}
                <span class="stat-add">+{cs.linesAdded}</span>
              {/if}
              {#if cs.linesRemoved > 0}
                <span class="stat-del">-{cs.linesRemoved}</span>
              {/if}
              <button
                class="view-toggle"
                class:viewed
                title={viewed ? 'Mark as unviewed' : 'Mark as viewed'}
                onclick={(e) => toggleViewed(cs.displayPath, e)}
              >
                <i class="fa {viewed ? 'fa-check-circle' : 'fa-circle-o'}"></i>
              </button>
            </span>
          </a>
        {/if}
      {/each}
    {/snippet}
    {@render renderNodes(tree, 0)}
  </div>

  {#if changesets.length > 0}
    <div class="tree-summary">
      {changesets.length} files
      <span class="stat-add">+{changesets.reduce((s, c) => s + c.linesAdded, 0)}</span>
      <span class="stat-del">-{changesets.reduce((s, c) => s + c.linesRemoved, 0)}</span>
    </div>
  {/if}
</div>

<style>
  .file-tree {
    font-size: 12px;
    user-select: none;
  }

  .filter-box {
    padding: 6px 8px;
    position: relative;
  }
  .filter-icon {
    position: absolute;
    left: 16px;
    top: 50%;
    transform: translateY(-50%);
    color: var(--text-muted);
    font-size: 10px;
    pointer-events: none;
  }
  .filter-box input {
    width: 100%;
    padding: 4px 8px 4px 24px;
    border: 1px solid var(--border);
    border-radius: 3px;
    font-size: 11px;
    outline: none;
    box-sizing: border-box;
    background: var(--bg-card);
    color: var(--text);
  }
  .filter-box input:focus {
    border-color: var(--text-link);
  }

  .tree-nodes {
    display: flex;
    flex-direction: column;
  }

  /* Directory row */
  .tree-dir {
    all: unset;
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 3px 8px;
    cursor: pointer;
    color: var(--text);
    white-space: nowrap;
  }
  .tree-dir:hover {
    background: var(--bg-hover);
  }
  .caret {
    width: 10px;
    text-align: center;
    font-size: 11px;
    color: var(--text-muted);
    flex-shrink: 0;
  }
  .folder-icon {
    font-size: 11px;
    color: var(--text-link);
    flex-shrink: 0;
  }
  .dir-name {
    font-weight: 600;
    color: var(--text);
  }

  /* File row */
  .tree-file {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 3px 8px;
    cursor: pointer;
    text-decoration: none;
    color: var(--text);
    white-space: nowrap;
  }
  .tree-file:hover {
    background: var(--bg-hover);
  }
  .tree-file.active {
    background: rgba(19, 108, 178, 0.1);
    border-right: 2px solid var(--blue);
  }

  .file-status-icon {
    width: 10px;
    text-align: center;
    font-size: 10px;
    flex-shrink: 0;
  }

  .file-name {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .file-name.viewed {
    opacity: 0.5;
    text-decoration: line-through;
  }

  .file-meta {
    display: flex;
    align-items: center;
    gap: 4px;
    flex-shrink: 0;
    margin-left: auto;
  }

  /* Stats */
  .stat-add {
    color: var(--green);
    font-size: 10px;
    font-weight: 600;
  }
  .stat-del {
    color: var(--red);
    font-size: 10px;
    font-weight: 600;
  }

  /* Badges */
  .badge {
    display: inline-flex;
    align-items: center;
    gap: 2px;
    font-size: 10px;
    padding: 0 4px;
    border-radius: 8px;
  }
  .comment-badge {
    color: var(--text-link);
    background: rgba(19, 108, 178, 0.08);
  }

  /* Viewed toggle */
  .view-toggle {
    all: unset;
    cursor: pointer;
    font-size: 11px;
    color: var(--text-muted);
    opacity: 0;
    transition: opacity 0.1s;
  }
  .tree-file:hover .view-toggle {
    opacity: 1;
  }
  .view-toggle.viewed {
    color: var(--green);
    opacity: 1;
  }

  /* Summary footer */
  .tree-summary {
    padding: 6px 8px;
    border-top: 1px solid var(--border-subtle);
    font-size: 11px;
    color: var(--text-muted);
    display: flex;
    gap: 6px;
    align-items: center;
  }
</style>
