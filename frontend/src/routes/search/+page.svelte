<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import { Tag, InfoView } from '$lib/components/phui';
  import { apiFetch } from '$lib/api';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import type { APISearchResponse } from '$lib/types';
  import { S } from '$lib/strings';

  const crumbs = [
    { name: S.crumb.home, href: '/' },
    { name: S.search.title }
  ];

  type SearchType = 'prs' | 'issues' | 'code' | 'repos';

  const searchTypes: { value: SearchType; label: string; icon: string }[] = [
    { value: 'code', label: S.search.code, icon: 'fa-code' },
    { value: 'repos', label: S.search.repos, icon: 'fa-database' },
    { value: 'issues', label: S.search.issues, icon: 'fa-circle-o' },
    { value: 'prs', label: S.search.prs, icon: 'fa-code-fork' },
  ];

  const sortOptions = [
    { value: '', label: S.search.sortBestMatch },
    { value: 'created-desc', label: S.search.sortNewest },
    { value: 'created-asc', label: S.search.sortOldest },
    { value: 'comments-desc', label: S.search.sortMostCommented },
    { value: 'updated-desc', label: S.search.sortRecentlyUpdated },
  ];

  const langColors: Record<string, string> = {
    Go: '#00ADD8', TypeScript: '#3178C6', JavaScript: '#F7DF1E',
    Python: '#3572A5', Rust: '#DEA584', Java: '#B07219',
    Ruby: '#701516', PHP: '#4F5D95', 'C++': '#f34b7d',
    C: '#555555', Shell: '#89e051', Swift: '#F05138',
    Kotlin: '#A97BFF', Dart: '#00B4AB', Scala: '#c22d40',
    HTML: '#e34c26', CSS: '#563d7c', Vue: '#41b883',
    Svelte: '#ff3e00',
  };

  let query = $state($page.url.searchParams.get('q') ?? '');
  let searchType = $state<SearchType>(($page.url.searchParams.get('type') as SearchType) ?? 'prs');
  let sortBy = $state($page.url.searchParams.get('sort') ?? '');
  let results: APISearchResponse | null = $state(null);
  let loading = $state(false);
  let searched = $state(!!$page.url.searchParams.get('q'));
  let collapsedCode: Record<number, boolean> = $state({});

  if ($page.url.searchParams.get('q')) {
    doSearch();
  }

  async function doSearch() {
    if (!query.trim()) return;
    loading = true;
    searched = true;
    try {
      let url = `/api/search?type=${searchType}&q=${encodeURIComponent(query)}`;
      if (sortBy) url += `&sort=${sortBy}`;
      results = await apiFetch<APISearchResponse>(url);
    } catch {
      results = null;
    } finally {
      loading = false;
    }
  }

  function updateUrl() {
    let url = `/search?q=${encodeURIComponent(query)}&type=${searchType}`;
    if (sortBy) url += `&sort=${sortBy}`;
    goto(url, { replaceState: true, keepFocus: true });
  }

  function handleSubmit(e: Event) {
    e.preventDefault();
    updateUrl();
    doSearch();
  }

  function selectType(t: SearchType) {
    searchType = t;
    if (query.trim()) {
      updateUrl();
      doSearch();
    }
  }

  function setSort(s: string) {
    sortBy = s;
    if (query.trim()) {
      updateUrl();
      doSearch();
    }
  }

  function appendStateFilter(state: string) {
    // Remove existing is:open / is:closed from query
    let q = query.replace(/\bis:(open|closed)\b/g, '').trim();
    q += ` is:${state}`;
    query = q.trim();
    updateUrl();
    doSearch();
  }

  let activeState = $derived(
    query.includes('is:open') ? 'open' : query.includes('is:closed') ? 'closed' : null
  );

  function formatCount(n: number): string {
    if (n >= 1_000_000) return `${Math.round(n / 1_000_000)}M`;
    if (n >= 1_000) return `${Math.round(n / 1_000)}k`;
    return String(n);
  }

  function relativeDate(iso: string): string {
    const d = new Date(iso);
    const now = new Date();
    const diffMs = now.getTime() - d.getTime();
    const diffSec = Math.floor(diffMs / 1000);
    const diffMin = Math.floor(diffSec / 60);
    const diffHr = Math.floor(diffMin / 60);
    const diffDay = Math.floor(diffHr / 24);
    const diffMon = Math.floor(diffDay / 30);
    const diffYr = Math.floor(diffDay / 365);
    if (diffSec < 60) return 'just now';
    if (diffMin < 60) return `${diffMin}m ago`;
    if (diffHr < 24) return `${diffHr}h ago`;
    if (diffDay < 30) return `${diffDay}d ago`;
    if (diffMon < 12) return `${diffMon}mo ago`;
    return `${diffYr}y ago`;
  }

  function truncate(s: string, max: number): string {
    if (!s) return '';
    return s.length > max ? s.slice(0, max) + '...' : s;
  }

  function prHref(repo: string, num: number): string {
    const [owner, name] = repo.split('/');
    return `/pr/${owner}/${name}/${num}`;
  }

  function repoHref(fullName: string): string {
    const [owner, name] = fullName.split('/');
    return `/repos/${owner}/${name}`;
  }

  function toggleCodeCollapse(idx: number) {
    collapsedCode[idx] = !collapsedCode[idx];
  }

  function labelTextColor(hex: string): string {
    const r = parseInt(hex.slice(0, 2), 16);
    const g = parseInt(hex.slice(2, 4), 16);
    const b = parseInt(hex.slice(4, 6), 16);
    const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;
    return luminance > 0.5 ? '#24292e' : '#fff';
  }

  let counts = $derived(results?.counts ?? {});
  let totalCount = $derived((results?.counts ?? {})[searchType] ?? 0);
</script>

<PageShell title={S.search.title} icon="fa-search">
  {#snippet breadcrumbs()}
    <Breadcrumbs {crumbs} />
  {/snippet}

  <div class="search-bar">
    <form onsubmit={handleSubmit}>
      <div class="search-row">
        <input
          type="text"
          bind:value={query}
          placeholder="Search..."
          class="search-input"
        />
        <button type="submit" class="search-btn" disabled={loading}>
          <i class="fa fa-search"></i>
        </button>
      </div>
    </form>
  </div>

  {#if searched}
    <div class="search-layout">
      <!-- Left sidebar -->
      <aside class="search-sidebar">
        <div class="sidebar-section">
          <h3 class="sidebar-heading">{S.search.filterBy}</h3>
          <ul class="type-list">
            {#each searchTypes as st}
              <li>
                <button
                  type="button"
                  class="type-item"
                  class:active={searchType === st.value}
                  onclick={() => selectType(st.value)}
                >
                  <span class="type-item-left">
                    <i class="fa {st.icon} type-icon"></i>
                    <span class="type-label">{st.label}</span>
                  </span>
                  {#if counts[st.value] !== undefined}
                    <span class="count-badge">{formatCount(counts[st.value])}</span>
                  {/if}
                </button>
              </li>
            {/each}
          </ul>
        </div>

        {#if searchType === 'prs' || searchType === 'issues'}
          <div class="sidebar-section">
            <h3 class="sidebar-heading">{S.search.state}</h3>
            <ul class="state-list">
              <li>
                <button
                  type="button"
                  class="state-item"
                  class:active={activeState === 'open'}
                  onclick={() => appendStateFilter('open')}
                >
                  <i class="fa fa-circle-o state-icon"></i>
                  {S.search.open}
                </button>
              </li>
              <li>
                <button
                  type="button"
                  class="state-item"
                  class:active={activeState === 'closed'}
                  onclick={() => appendStateFilter('closed')}
                >
                  <i class="fa fa-check-circle-o state-icon"></i>
                  {S.search.closed}
                </button>
              </li>
            </ul>
          </div>
        {/if}

        {#if searchType === 'code'}
          {@const codeResults = results?.code ?? []}
          {@const langs = [...new Set(codeResults.map(c => c.language).filter(Boolean))]}
          {#if langs.length > 0}
            <div class="sidebar-section">
              <h3 class="sidebar-heading">{S.search.languages}</h3>
              <ul class="lang-list">
                {#each langs as lang}
                  <li class="lang-item">
                    <span class="lang-dot" style="background:{langColors[lang] ?? 'var(--text-muted)'}"></span>
                    {lang}
                  </li>
                {/each}
              </ul>
            </div>
          {/if}
        {/if}
      </aside>

      <!-- Main results -->
      <div class="search-main">
        {#if loading}
          <div class="loading-state">
            <i class="fa fa-spinner fa-spin mrs"></i>
            Searching...
          </div>
        {:else}
          <div class="results-header">
            <span class="results-count">
              {formatCount(totalCount)} {S.search.results}
            </span>
            <div class="sort-control">
              <span class="sort-label">{S.search.sortBy}</span>
              <select class="sort-select" bind:value={sortBy} onchange={() => setSort(sortBy)}>
                {#each sortOptions as opt}
                  <option value={opt.value}>{opt.label}</option>
                {/each}
              </select>
            </div>
          </div>

          <!-- PR Results -->
          {#if searchType === 'prs' && results?.prs}
            {#if results.prs.length === 0}
              <InfoView icon="fa-inbox">{S.search.noResults}</InfoView>
            {:else}
              {#each results.prs as pr}
                <div class="result-card">
                  <div class="result-repo">{pr.repo}</div>
                  <div class="result-title-row">
                    {#if pr.state === 'merged'}
                      <i class="fa fa-code-fork result-icon merged"></i>
                    {:else if pr.state === 'closed'}
                      <i class="fa fa-times-circle result-icon closed"></i>
                    {:else}
                      <i class="fa fa-circle-o result-icon open"></i>
                    {/if}
                    <a href={prHref(pr.repo, pr.number)} class="result-title">{pr.title}</a>
                    {#if pr.commentsCount > 0}
                      <span class="comment-count">
                        <i class="fa fa-comment-o"></i>
                        {pr.commentsCount}
                      </span>
                    {/if}
                  </div>
                  {#if pr.body}
                    <div class="result-body">{truncate(pr.body, 200)}</div>
                  {/if}
                  <div class="result-meta">
                    <img src={pr.avatarURL} alt="" class="meta-avatar" />
                    <span class="meta-author">{pr.author}</span>
                    <span class="meta-sep">opened {relativeDate(pr.createdAt)}</span>
                    <span class="meta-sep">#{pr.number}</span>
                  </div>
                  {#if pr.labels && pr.labels.length > 0}
                    <div class="result-labels">
                      {#each pr.labels as label}
                        <span class="label-tag" style="background:#{label.color}; color:{labelTextColor(label.color)}">
                          {label.name}
                        </span>
                      {/each}
                    </div>
                  {/if}
                </div>
              {/each}
            {/if}

          <!-- Issue Results -->
          {:else if searchType === 'issues' && results?.issues}
            {#if results.issues.length === 0}
              <InfoView icon="fa-inbox">{S.search.noResults}</InfoView>
            {:else}
              {#each results.issues as issue}
                <div class="result-card">
                  <div class="result-repo">{issue.repo}</div>
                  <div class="result-title-row">
                    {#if issue.state === 'closed'}
                      <i class="fa fa-check-circle result-icon closed"></i>
                    {:else}
                      <i class="fa fa-circle-o result-icon open"></i>
                    {/if}
                    <span class="result-title">{issue.title}</span>
                    {#if issue.commentsCount > 0}
                      <span class="comment-count">
                        <i class="fa fa-comment-o"></i>
                        {issue.commentsCount}
                      </span>
                    {/if}
                  </div>
                  {#if issue.body}
                    <div class="result-body">{truncate(issue.body, 200)}</div>
                  {/if}
                  <div class="result-meta">
                    <img src={issue.avatarURL} alt="" class="meta-avatar" />
                    <span class="meta-author">{issue.author}</span>
                    <span class="meta-sep">opened {relativeDate(issue.createdAt)}</span>
                    <span class="meta-sep">#{issue.number}</span>
                  </div>
                  {#if issue.labels && issue.labels.length > 0}
                    <div class="result-labels">
                      {#each issue.labels as label}
                        <span class="label-tag" style="background:#{label.color}; color:{labelTextColor(label.color)}">
                          {label.name}
                        </span>
                      {/each}
                    </div>
                  {/if}
                </div>
              {/each}
            {/if}

          <!-- Code Results -->
          {:else if searchType === 'code' && results?.code}
            {#if results.code.length === 0}
              <InfoView icon="fa-inbox">{S.search.noResults}</InfoView>
            {:else}
              {#each results.code as result, idx}
                <div class="result-card code-card">
                  <button type="button" class="code-header" onclick={() => toggleCodeCollapse(idx)}>
                    <div class="code-header-left">
                      <i class="fa fa-chevron-{collapsedCode[idx] ? 'right' : 'down'} chevron-icon"></i>
                      <span class="code-repo">{result.repo}</span>
                      <span class="code-sep">&middot;</span>
                      <span class="code-path">{result.path}</span>
                    </div>
                    <div class="code-header-right">
                      {#if result.language}
                        <Tag shade="grey">{result.language}</Tag>
                      {/if}
                      {#if result.matchCount > 0}
                        <span class="match-count">{result.matchCount} {S.search.matches}</span>
                      {/if}
                    </div>
                  </button>
                  {#if !collapsedCode[idx]}
                    <div class="code-snippet">
                      <pre class="code-fragment">{@html result.fragment}</pre>
                    </div>
                  {/if}
                </div>
              {/each}
            {/if}

          <!-- Repo Results -->
          {:else if searchType === 'repos' && results?.repos}
            {#if results.repos.length === 0}
              <InfoView icon="fa-inbox">{S.search.noResults}</InfoView>
            {:else}
              {#each results.repos as repo}
                <div class="result-card">
                  <a href={repoHref(repo.fullName)} class="repo-name">{repo.fullName}</a>
                  {#if repo.description}
                    <div class="repo-desc">{repo.description}</div>
                  {/if}
                  <div class="repo-meta">
                    {#if repo.language}
                      <span class="repo-meta-item">
                        <span class="lang-dot" style="background:{langColors[repo.language] ?? 'var(--text-muted)'}"></span>
                        {repo.language}
                      </span>
                    {/if}
                    <span class="repo-meta-item">
                      <i class="fa fa-star"></i>
                      {formatCount(repo.stars)}
                    </span>
                    <span class="repo-meta-item">
                      <i class="fa fa-code-fork"></i>
                      {formatCount(repo.forks)}
                    </span>
                    <span class="repo-meta-item">
                      Updated {relativeDate(repo.updatedAt)}
                    </span>
                  </div>
                </div>
              {/each}
            {/if}
          {/if}
        {/if}
      </div>
    </div>
  {/if}
</PageShell>


<style>
  /* Search bar */
  .search-bar {
    margin-bottom: 12px;
  }

  .search-row {
    display: flex;
    gap: 0;
  }

  .search-input {
    flex: 1;
    padding: 8px 12px;
    border: 1px solid var(--border);
    border-right: none;
    border-radius: 4px 0 0 4px;
    font-size: 14px;
    color: var(--text);
    background: var(--bg-card);
  }

  .search-input:focus {
    outline: none;
    border-color: var(--blue);
  }

  .search-btn {
    padding: 8px 14px;
    background: var(--blue);
    color: var(--text-on-dark);
    border: 1px solid var(--blue);
    border-radius: 0 4px 4px 0;
    font-size: 14px;
    cursor: pointer;
  }

  .search-btn:hover {
    opacity: 0.9;
  }

  .search-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  /* Two-column layout */
  .search-layout {
    display: flex;
    gap: 16px;
    align-items: flex-start;
  }

  /* Sidebar */
  .search-sidebar {
    width: 250px;
    flex-shrink: 0;
  }

  .sidebar-section {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 4px;
    margin-bottom: 12px;
    overflow: hidden;
  }

  .sidebar-heading {
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--text-muted);
    padding: 10px 12px 6px;
    margin: 0;
  }

  .type-list, .state-list, .lang-list {
    list-style: none;
    margin: 0;
    padding: 0 0 4px;
  }

  .type-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    padding: 7px 12px;
    border: none;
    border-left: 3px solid transparent;
    background: none;
    color: var(--text);
    font-size: 13px;
    cursor: pointer;
    text-align: left;
  }

  .type-item:hover {
    background: var(--bg-hover);
  }

  .type-item.active {
    border-left-color: var(--blue);
    font-weight: 600;
    background: var(--bg-hover);
  }

  .type-item-left {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .type-icon {
    width: 16px;
    text-align: center;
    color: var(--text-muted);
  }

  .type-item.active .type-icon {
    color: var(--text);
  }

  .count-badge {
    font-size: 11px;
    color: var(--text-muted);
    background: var(--bg-subtle);
    padding: 1px 6px;
    border-radius: 10px;
    font-weight: 600;
  }

  .state-item {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    padding: 6px 12px;
    border: none;
    background: none;
    color: var(--text);
    font-size: 13px;
    cursor: pointer;
    text-align: left;
  }

  .state-item:hover {
    background: var(--bg-hover);
  }

  .state-item.active {
    font-weight: 600;
    color: var(--blue);
  }

  .state-icon {
    color: var(--text-muted);
  }

  .state-item.active .state-icon {
    color: var(--blue);
  }

  .lang-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 5px 12px;
    font-size: 13px;
    color: var(--text);
  }

  .lang-dot {
    display: inline-block;
    width: 10px;
    height: 10px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  /* Main results area */
  .search-main {
    flex: 1;
    min-width: 0;
  }

  .loading-state {
    padding: 24px;
    text-align: center;
    color: var(--text-muted);
  }

  .results-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 0;
    margin-bottom: 8px;
  }

  .results-count {
    font-size: 14px;
    font-weight: 600;
    color: var(--text);
  }

  .sort-control {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .sort-label {
    font-size: 13px;
    color: var(--text-muted);
  }

  .sort-select {
    padding: 4px 8px;
    border: 1px solid var(--border);
    border-radius: 4px;
    background: var(--bg-card);
    color: var(--text);
    font-size: 13px;
    cursor: pointer;
  }

  /* Result cards */
  .result-card {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 4px;
    padding: 14px 16px;
    margin-bottom: 8px;
  }

  .result-repo {
    font-size: 12px;
    color: var(--text-muted);
    margin-bottom: 4px;
  }

  .result-title-row {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-bottom: 4px;
  }

  .result-icon {
    font-size: 14px;
    flex-shrink: 0;
  }

  .result-icon.open {
    color: var(--green);
  }

  .result-icon.closed {
    color: var(--red);
  }

  .result-icon.merged {
    color: var(--violet);
  }

  .result-title {
    font-size: 15px;
    font-weight: 600;
    color: var(--text-link);
    text-decoration: none;
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .result-title:hover {
    text-decoration: underline;
  }

  .comment-count {
    display: flex;
    align-items: center;
    gap: 3px;
    font-size: 12px;
    color: var(--text-muted);
    flex-shrink: 0;
  }

  .result-body {
    font-size: 13px;
    color: var(--text-muted);
    line-height: 1.4;
    margin-bottom: 6px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .result-meta {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    color: var(--text-muted);
  }

  .meta-avatar {
    width: 16px;
    height: 16px;
    border-radius: 50%;
  }

  .meta-author {
    font-weight: 600;
    color: var(--text);
  }

  .meta-sep::before {
    content: '\00b7';
    margin-right: 6px;
  }

  .result-labels {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    margin-top: 6px;
  }

  .label-tag {
    display: inline-block;
    padding: 1px 7px;
    border-radius: 12px;
    font-size: 11px;
    font-weight: 600;
    line-height: 18px;
  }

  /* Code results */
  .code-card {
    padding: 0;
    overflow: hidden;
  }

  .code-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    padding: 10px 14px;
    border: none;
    background: var(--bg-card-header);
    color: var(--text);
    font-size: 13px;
    cursor: pointer;
    text-align: left;
  }

  .code-header:hover {
    background: var(--bg-hover);
  }

  .code-header-left {
    display: flex;
    align-items: center;
    gap: 6px;
    min-width: 0;
  }

  .chevron-icon {
    font-size: 10px;
    color: var(--text-muted);
    width: 12px;
  }

  .code-repo {
    font-weight: 600;
  }

  .code-sep {
    color: var(--text-muted);
  }

  .code-path {
    color: var(--text-link);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .code-header-right {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-shrink: 0;
  }

  .match-count {
    font-size: 11px;
    color: var(--text-muted);
  }

  .code-snippet {
    border-top: 1px solid var(--border);
    max-height: 250px;
    overflow: auto;
  }

  .code-fragment {
    margin: 0;
    padding: 8px 12px;
    font-family: var(--font-mono);
    font-size: 12px;
    line-height: 1.5;
    overflow-x: auto;
  }

  /* Repo results */
  .repo-name {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-link);
  }

  .repo-desc {
    font-size: 13px;
    color: var(--text-muted);
    margin-top: 4px;
    line-height: 1.4;
  }

  .repo-meta {
    display: flex;
    align-items: center;
    gap: 14px;
    margin-top: 8px;
    font-size: 12px;
    color: var(--text-muted);
  }

  .repo-meta-item {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .search-layout {
      flex-direction: column;
    }
    .search-sidebar {
      width: 100%;
    }
  }
</style>
