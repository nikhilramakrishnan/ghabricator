<script lang="ts">
  import { S } from '$lib/strings';
  import { fileTreeData, fileTreeOpen } from '$lib/stores/filetree';
  import { FileTree } from '$lib/components/diff';

  let { navActive = '' }: { navActive?: string } = $props();
</script>

<aside class="sidebar" class:tree-open={$fileTreeOpen && $fileTreeData}>
  <div class="sidebar-nav">
    <div class="sidebar-section">
      <a href="/dashboard" class="sidebar-link" class:active={navActive === 'revisions'}>
        <i class="fa fa-code-fork"></i>
        <span>{S.nav.revisions}</span>
      </a>
      <a href="/repos" class="sidebar-link" class:active={navActive === 'repos'}>
        <i class="fa fa-database"></i>
        <span>{S.nav.repositories}</span>
      </a>
      <a href="/search" class="sidebar-link" class:active={navActive === 'search'}>
        <i class="fa fa-search"></i>
        <span>{S.nav.search}</span>
      </a>
    </div>

    <div class="sidebar-section">
      <div class="section-label">{S.nav.tools}</div>
      <a href="/paste" class="sidebar-link" class:active={navActive === 'paste'}>
        <i class="fa fa-clipboard"></i>
        <span>{S.nav.paste}</span>
      </a>
      <a href="/actions" class="sidebar-link" class:active={navActive === 'actions'}>
        <i class="fa fa-bullhorn"></i>
        <span>{S.nav.actions}</span>
      </a>
    </div>

    {#if $fileTreeData}
      <div class="sidebar-section tool-section">
        <div class="section-label">PANELS</div>
        <button
          class="sidebar-link tool-toggle"
          class:active={$fileTreeOpen}
          onclick={() => fileTreeOpen.update(v => !v)}
        >
          <i class="fa fa-sitemap"></i>
          <span>Files</span>
          <span class="tool-badge">{$fileTreeData.changesets.length}</span>
        </button>
      </div>
    {/if}
  </div>

  {#if $fileTreeOpen && $fileTreeData}
    <div class="sidebar-panel">
      <div class="panel-header">
        <i class="fa fa-sitemap"></i>
        <span>Files</span>
        <button class="panel-close" onclick={() => fileTreeOpen.set(false)}>
          <i class="fa fa-times"></i>
        </button>
      </div>
      <div class="panel-body">
        <FileTree
          changesets={$fileTreeData.changesets}
          activeFile={$fileTreeData.activeFile}
          commentCounts={$fileTreeData.commentCounts}
        />
      </div>
    </div>
  {/if}
</aside>

<style>
  .sidebar {
    width: 200px;
    flex-shrink: 0;
    background: var(--bg-card);
    border-right: 1px solid var(--border);
    display: flex;
    height: 100%;
    transition: width 0.15s ease;
  }

  .sidebar.tree-open {
    width: 480px;
  }

  .sidebar-nav {
    width: 200px;
    flex-shrink: 0;
    padding: 12px 0;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
  }

  .sidebar-section {
    margin-bottom: 16px;
  }

  .tool-section {
    margin-top: auto;
    margin-bottom: 0;
    border-top: 1px solid var(--border-subtle);
    padding-top: 8px;
  }

  .section-label {
    font-size: 11px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--text-muted);
    padding: 4px 16px 6px;
  }

  .sidebar-link {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 16px;
    font-size: 13px;
    color: var(--text);
    text-decoration: none;
  }

  .sidebar-link:hover {
    background: var(--bg-subtle);
    text-decoration: none;
  }

  .sidebar-link.active {
    background: var(--bg-subtle);
    color: var(--text-link);
    font-weight: 600;
  }

  .sidebar-link .fa {
    width: 16px;
    text-align: center;
    color: var(--text-muted);
    font-size: 13px;
  }

  .sidebar-link.active .fa {
    color: var(--text-link);
  }

  .tool-toggle {
    all: unset;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 16px;
    font-size: 13px;
    color: var(--text);
    cursor: pointer;
    width: 100%;
    box-sizing: border-box;
  }
  .tool-toggle:hover {
    background: var(--bg-subtle);
  }
  .tool-toggle.active {
    background: var(--bg-subtle);
    color: var(--text-link);
    font-weight: 600;
  }
  .tool-toggle .fa {
    width: 16px;
    text-align: center;
    color: var(--text-muted);
    font-size: 13px;
  }
  .tool-toggle.active .fa {
    color: var(--text-link);
  }

  .tool-badge {
    margin-left: auto;
    background: var(--bg-subtle);
    color: var(--text-muted);
    font-size: 10px;
    font-weight: 600;
    padding: 1px 6px;
    border-radius: 8px;
  }

  /* File tree panel */
  .sidebar-panel {
    flex: 1;
    border-left: 1px solid var(--border-subtle);
    display: flex;
    flex-direction: column;
    min-width: 0;
    overflow: hidden;
  }

  .panel-header {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 12px;
    font-size: 12px;
    font-weight: 600;
    color: var(--text);
    border-bottom: 1px solid var(--border-subtle);
    flex-shrink: 0;
  }
  .panel-header .fa {
    color: var(--text-muted);
    font-size: 11px;
  }

  .panel-close {
    all: unset;
    margin-left: auto;
    cursor: pointer;
    color: var(--text-muted);
    padding: 2px 4px;
    font-size: 11px;
  }
  .panel-close:hover {
    color: var(--text);
  }

  .panel-body {
    flex: 1;
    overflow-y: auto;
  }
</style>
