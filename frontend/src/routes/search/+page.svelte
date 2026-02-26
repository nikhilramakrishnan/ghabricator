<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import {
    Box, ObjectItemList, ObjectItem, Attribute, Tag, InfoView
  } from '$lib/components/phui';
  import { apiFetch } from '$lib/api';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import type { APISearchResponse } from '$lib/types';

  const crumbs = [
    { name: 'Home', href: '/' },
    { name: 'Search' }
  ];

  const searchTypes = [
    { value: 'prs', label: 'Pull Requests', icon: 'fa-code-fork' },
    { value: 'code', label: 'Code', icon: 'fa-code' },
    { value: 'repos', label: 'Repositories', icon: 'fa-database' },
  ];

  let query = $state($page.url.searchParams.get('q') ?? '');
  let searchType = $state($page.url.searchParams.get('type') ?? 'prs');
  let results: APISearchResponse | null = $state(null);
  let loading = $state(false);
  let searched = $state(!!$page.url.searchParams.get('q'));

  // Run search on initial load if query params present
  if (query) {
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

<PageShell title="Search" icon="fa-search">
  {#snippet breadcrumbs()}
    <Breadcrumbs {crumbs} />
  {/snippet}

  <Box border>
    <div style="padding:16px">
      <form onsubmit={handleSubmit}>
        <div style="display:flex;gap:8px;margin-bottom:12px">
          <input
            type="text"
            bind:value={query}
            placeholder="Search..."
            class="aphront-form-input"
            style="flex:1;padding:8px 12px;border:1px solid #c7ccd9;border-radius:3px;font-size:14px"
          />
          <button type="submit" class="mood-btn mood-btn-green" disabled={loading}>
            <span class="phui-icon-view phui-font-fa fa-search mrs"></span>
            Search
          </button>
        </div>
      </form>

      <div style="display:flex;gap:4px">
        {#each searchTypes as st}
          <button
            type="button"
            class="mood-btn {searchType === st.value ? 'mood-btn-blue' : 'mood-btn-default'}"
            style="font-size:12px;padding:4px 12px"
            onclick={() => selectType(st.value)}
          >
            <span class="phui-icon-view phui-font-fa {st.icon} mrs"></span>
            {st.label}
          </button>
        {/each}
      </div>
    </div>
  </Box>

  {#if loading}
    <div style="padding:24px;text-align:center;color:#6b748c">
      <span class="phui-icon-view phui-font-fa fa-spinner fa-spin" style="margin-right:8px"></span>
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
            <div style="padding:12px 16px">
              <div style="font-size:13px;margin-bottom:6px">
                <span class="phui-icon-view phui-font-fa fa-github mrs" style="color:#6b748c"></span>
                <strong>{result.repo}</strong>
                <span style="color:#6b748c;margin:0 4px">/</span>
                <span style="color:#136CB2">{result.path}</span>
              </div>
              <div class="phabricator-source-code-container" style="max-height:200px;overflow:auto">
                <pre class="PhabricatorMonospaced" style="margin:0;padding:8px;background:#f7f7f7;border-radius:3px;font-size:12px;overflow-x:auto">{@html result.fragment}</pre>
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
                    <span style="display:inline-block;width:10px;height:10px;border-radius:50%;background:{langColors[repo.language] ?? '#6b748c'};margin-right:4px;vertical-align:middle"></span>
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
