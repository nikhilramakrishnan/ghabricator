<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import {
    ObjectItemList, ObjectItem, Attribute, Tag, Button, InfoView
  } from '$lib/components/phui';
  import { apiFetch } from '$lib/api';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import type { APISearchResponse, APISearchPR, APISearchCodeResult, APISearchRepoResult } from '$lib/types';

  const crumbs = [
    { name: 'Home', href: '/' },
    { name: 'Search' }
  ];

  let query = $state($page.url.searchParams.get('q') ?? '');
  let searchType = $state($page.url.searchParams.get('type') ?? 'pr');
  let results = $state<APISearchResponse | null>(null);
  let loading = $state(false);
  let searched = $state(false);

  // Auto-search if URL has query params
  $effect(() => {
    const q = $page.url.searchParams.get('q');
    const t = $page.url.searchParams.get('type');
    if (q) {
      query = q;
      if (t) searchType = t;
      doSearch();
    }
  });

  async function doSearch() {
    if (!query.trim()) return;
    loading = true;
    searched = true;
    try {
      results = await apiFetch<APISearchResponse>(`/api/search?q=${encodeURIComponent(query)}&type=${searchType}`);
    } catch {
      results = null;
    } finally {
      loading = false;
    }
  }

  function handleSubmit(e: Event) {
    e.preventDefault();
    goto(`/search?q=${encodeURIComponent(query)}&type=${searchType}`, { replaceState: true });
    doSearch();
  }

  let hasResults = $derived(
    results && ((results.prs?.length ?? 0) > 0 || (results.code?.length ?? 0) > 0 || (results.repos?.length ?? 0) > 0)
  );
</script>

<PageShell title="Search" icon="fa-search">
  {#snippet breadcrumbs()}
    <Breadcrumbs {crumbs} />
  {/snippet}

  <div class="phui-box phui-box-border phui-object-box" style="padding:16px">
    <form onsubmit={handleSubmit} style="display:flex;gap:8px;align-items:center;flex-wrap:wrap">
      <input
        type="text"
        bind:value={query}
        placeholder="Search..."
        class="aphront-form-input"
        style="flex:1;min-width:200px;padding:8px 12px;border:1px solid #c7ccd9;border-radius:3px;font-size:14px"
      />
      <select bind:value={searchType} style="padding:8px;border:1px solid #c7ccd9;border-radius:3px;font-size:13px">
        <option value="pr">Pull Requests</option>
        <option value="code">Code</option>
        <option value="repo">Repositories</option>
      </select>
      <Button type="submit" color="green" icon="fa-search">Search</Button>
    </form>
  </div>

  {#if loading}
    <InfoView icon="fa-spinner">Searching...</InfoView>
  {:else if searched && !hasResults}
    <InfoView icon="fa-inbox">No results found.</InfoView>
  {:else if results}

    <!-- PR results -->
    {#if results.prs?.length}
      <div class="phui-box phui-box-border phui-object-box" style="margin-top:12px">
        <div class="phui-header-shell">
          <div class="phui-header-view">
            <h1 class="phui-header-header">
              <span class="phui-header-icon phui-icon-view phui-font-fa fa-code-fork"></span>
              Pull Requests
            </h1>
          </div>
        </div>
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
                <Attribute icon="fa-user">{pr.author}</Attribute>
                <Attribute icon="fa-github">{pr.repo}#{pr.number}</Attribute>
                <Attribute icon="fa-clock-o">{pr.updatedAt}</Attribute>
              {/snippet}
            </ObjectItem>
          {/each}
        </ObjectItemList>
      </div>
    {/if}

    <!-- Code results -->
    {#if results.code?.length}
      <div class="phui-box phui-box-border phui-object-box" style="margin-top:12px">
        <div class="phui-header-shell">
          <div class="phui-header-view">
            <h1 class="phui-header-header">
              <span class="phui-header-icon phui-icon-view phui-font-fa fa-file-code-o"></span>
              Code
            </h1>
          </div>
        </div>
        <ObjectItemList>
          {#each results.code as code}
            <ObjectItem
              title={code.path}
              href="https://github.com/{code.repo}/blob/HEAD/{code.path}"
              icon="fa-file-code-o"
            >
              {#snippet attributes()}
                <Attribute icon="fa-github">{code.repo}</Attribute>
              {/snippet}
            </ObjectItem>
            {#if code.fragment}
              <div style="padding:4px 12px 12px 56px">
                <pre style="background:#f6f8fa;border:1px solid #e3e4e8;border-radius:4px;padding:8px;font-size:12px;overflow-x:auto;margin:0">{code.fragment}</pre>
              </div>
            {/if}
          {/each}
        </ObjectItemList>
      </div>
    {/if}

    <!-- Repo results -->
    {#if results.repos?.length}
      <div class="phui-box phui-box-border phui-object-box" style="margin-top:12px">
        <div class="phui-header-shell">
          <div class="phui-header-view">
            <h1 class="phui-header-header">
              <span class="phui-header-icon phui-icon-view phui-font-fa fa-database"></span>
              Repositories
            </h1>
          </div>
        </div>
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
                  <Attribute icon="fa-code">{repo.language}</Attribute>
                {/if}
                <Attribute icon="fa-star">{repo.stars}</Attribute>
              {/snippet}
            </ObjectItem>
          {/each}
        </ObjectItemList>
      </div>
    {/if}
  {/if}
</PageShell>
