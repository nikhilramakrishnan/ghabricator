<script lang="ts">
  import { S } from '$lib/strings';

  let { navActive = '' }: { navActive?: string } = $props();
</script>

<aside class="sidebar">
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
</aside>

<style>
  .sidebar {
    width: 200px;
    flex-shrink: 0;
    background: var(--bg-card);
    border-right: 1px solid var(--border);
    padding: 12px 0;
    overflow-y: auto;
    height: 100%;
  }

  .sidebar-section {
    margin-bottom: 16px;
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
</style>
