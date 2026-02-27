<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import {
    Box, ObjectItemList, ObjectItem, Attribute, Tag, InfoView
  } from '$lib/components/phui';
  import { apiFetch } from '$lib/api';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import type { APISearchResponse } from '$lib/types';
  import { S } from '$lib/strings';

  const crumbs = [
    { name: S.crumb.home, href: '/' },
    { name: S.search.title }
  ];

  const searchTypes = [
    { value: 'prs', label: S.search.prs, icon: 'fa-code-fork' },
    { value: 'code', label: S.search.code, icon: 'fa-code' },
    { value: 'repos', label: S.search.repos, icon: 'fa-database' },
  ];

  let query = $state($page.url.searchParams.get('q') ?? '');
  let searchType = $state($page.url.searchParams.get('type') ?? 'prs');
  let results: APISearchResponse | null = $state(null);
  let loading = $state(false);
  let searched = $state(!!$page.url.searchParams.get('q'));

  // Run search on initial load if query params present
  if ($page.url.searchParams.get('q')) {
    doSearch();
  }

  async function doSearch() {
    if (!query.trim()) return;
    loading = true;
    searched = true;
    try {
      results = await apiFetch<APISearchResponse>(`/api/search?type=${searchType}&q=${encodeURIComponent(query)}`);
    } catch {
      results = null;
    } finally {
      loading = false;
    }
  }

  function handleSubmit(e: Event) {
    e.preventDefault();
    const url = `/search?q=${encodeURIComponent(query)}&type=${searchType}`;
    goto(url, { replaceState: true, keepFocus: true });
    doSearch();
  }

  function selectType(t: string) {
    searchType = t;
    if (query.trim()) {
      const url = `/search?q=${encodeURIComponent(query)}&type=${t}`;
      goto(url, { replaceState: true, keepFocus: true });
      doSearch();
    }
  }

  const langColors: Record<string, string> = {
    Go: '#00ADD8', TypeScript: '#3178C6', JavaScript: '#F7DF1E',
    Python: '#3572A5', Rust: '#DEA584', Java: '#B07219',
    Ruby: '#701516', PHP: '#4F5D95', 'C++': '#f34b7d',
    C: '#555555', Shell: '#89e051',
  };
</script>

<PageShell title={S.search.title} icon="fa-search">
  {#snippet breadcrumbs()}
    <Breadcrumbs {crumbs} />
  {/snippet}

  <Box border>
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
            <i class="fa fa-search mrs"></i>
            Search
          </button>
        </div>
      </form>

      <div class="type-tabs">
        {#each searchTypes as st}
          <button
            type="button"
            class="type-tab"
            class:active={searchType === st.value}
            onclick={() => selectType(st.value)}
          >
            <i class="fa {st.icon} mrs"></i>
            {st.label}
          </button>
        {/each}
      </div>
    </div>
  </Box>

  {#if loading}
    <div class="loading-state">
      <i class="fa fa-spinner fa-spin mrs"></i>
      Searching...
    </div>
  {:else if searched && results}
    {#if searchType === 'prs' && results.prs}
      {#if results.prs.length === 0}
        <InfoView icon="fa-inbox">No pull requests found.</InfoView>
      {:else}
        <ObjectItemList>
          {#each results.prs as pr}
            <ObjectItem
              title={pr.title}
              href={pr.url}
              imageUrl={pr.avatarURL}
              barColor={pr.draft ? 'grey' : 'blue'}
            >
              {#snippet tags()}
                {#if pr.draft}
                  <Tag shade="grey">Draft</Tag>
                {/if}
              {/snippet}
              {#snippet attributes()}
                <Attribute icon="fa-github">{pr.repo}#{pr.number}</Attribute>
                <Attribute icon="fa-user">{pr.author}</Attribute>
                <Attribute icon="fa-clock-o">{pr.updatedAt}</Attribute>
              {/snippet}
            </ObjectItem>
          {/each}
        </ObjectItemList>
      {/if}

    {:else if searchType === 'code' && results.code}
      {#if results.code.length === 0}
        <InfoView icon="fa-inbox">No code results found.</InfoView>
      {:else}
        {#each results.code as result}
          <Box border>
            <div class="code-result">
              <div class="code-result-header">
                <i class="fa fa-github mrs code-icon"></i>
                <strong>{result.repo}</strong>
                <span class="code-sep">/</span>
                <span class="code-path">{result.path}</span>
              </div>
              <div class="code-fragment-container">
                <pre class="code-fragment">{@html result.fragment}</pre>
              </div>
            </div>
          </Box>
        {/each}
      {/if}

    {:else if searchType === 'repos' && results.repos}
      {#if results.repos.length === 0}
        <InfoView icon="fa-inbox">No repositories found.</InfoView>
      {:else}
        <ObjectItemList>
          {#each results.repos as repo}
            {@const parts = repo.fullName.split('/')}
            <ObjectItem
              title={repo.fullName}
              href="/repos/{parts[0]}/{parts[1]}"
              imageUrl={repo.avatarURL}
              barColor="blue"
            >
              {#snippet attributes()}
                {#if repo.description}
                  <Attribute>{repo.description}</Attribute>
                {/if}
                {#if repo.language}
                  <Attribute>
                    <span class="lang-dot" style="background:{langColors[repo.language] ?? 'var(--text-muted)'}"></span>
                    {repo.language}
                  </Attribute>
                {/if}
                <Attribute icon="fa-star">{repo.stars}</Attribute>
              {/snippet}
            </ObjectItem>
          {/each}
        </ObjectItemList>
      {/if}
    {/if}
  {:else if searched}
    <InfoView icon="fa-inbox">No results found.</InfoView>
  {/if}
</PageShell>

<style>
  .search-bar {
    padding: 16px;
  }

  .search-row {
    display: flex;
    gap: 8px;
    margin-bottom: 12px;
  }

  .search-input {
    flex: 1;
    padding: 8px 12px;
    border: 1px solid var(--border);
    border-radius: 3px;
    font-size: 14px;
    color: var(--text);
    background: var(--bg-card);
  }

  .search-input:focus {
    outline: none;
    border-color: var(--blue);
  }

  .search-btn {
    padding: 8px 16px;
    background: var(--green);
    color: var(--text-on-dark);
    border: none;
    border-radius: 3px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
  }

  .search-btn:hover {
    background: var(--green-hover);
  }

  .search-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .type-tabs {
    display: flex;
    gap: 4px;
  }

  .type-tab {
    font-size: 12px;
    padding: 4px 12px;
    border: 1px solid var(--border);
    border-radius: 3px;
    background: var(--bg-card);
    color: var(--text);
    cursor: pointer;
  }

  .type-tab:hover {
    background: var(--bg-hover);
  }

  .type-tab.active {
    background: var(--blue);
    color: var(--text-on-dark);
    border-color: var(--blue);
  }

  .loading-state {
    padding: 24px;
    text-align: center;
    color: var(--text-muted);
  }

  .code-result {
    padding: 12px 16px;
  }

  .code-result-header {
    font-size: 13px;
    margin-bottom: 6px;
  }

  .code-icon {
    color: var(--text-muted);
  }

  .code-path {
    color: var(--text-link);
  }

  .code-sep {
    color: var(--text-muted);
    margin: 0 4px;
  }

  .code-fragment-container {
    max-height: 200px;
    overflow: auto;
  }

  .code-fragment {
    margin: 0;
    padding: 8px;
    background: var(--bg-subtle);
    border-radius: 3px;
    font-family: var(--font-mono);
    font-size: 12px;
    overflow-x: auto;
  }

  .lang-dot {
    display: inline-block;
    width: 10px;
    height: 10px;
    border-radius: 50%;
    margin-right: 4px;
    vertical-align: middle;
  }
</style>
