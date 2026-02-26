<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import {
    ObjectItemList, ObjectItem, Attribute, Tag
  } from '$lib/components/phui';
  import type { APIRepoSummary } from '$lib/types';

  let { data } = $props();
  let repos: APIRepoSummary[] = $derived(data.repos ?? []);

  const crumbs = [
    { name: 'Home', href: '/' },
    { name: 'Repositories' }
  ];

  // Language dot colors
  const langColors: Record<string, string> = {
    Go: '#00ADD8',
    TypeScript: '#3178C6',
    JavaScript: '#F7DF1E',
    Python: '#3572A5',
    Rust: '#DEA584',
    Java: '#B07219',
    Ruby: '#701516',
    PHP: '#4F5D95',
    'C++': '#f34b7d',
    C: '#555555',
    Shell: '#89e051',
    HTML: '#e34c26',
    CSS: '#563d7c',
  };

  function barColor(repo: APIRepoSummary): string {
    if (repo.archived) return 'grey';
    if (repo.fork) return 'violet';
    return 'blue';
  }
</script>

<PageShell title="Your Repositories" icon="fa-database">
  {#snippet breadcrumbs()}
    <Breadcrumbs {crumbs} />
  {/snippet}

  <ObjectItemList>
    {#each repos as repo}
      {@const parts = repo.fullName.split('/')}
      <ObjectItem
        title={repo.fullName}
        href="/repos/{parts[0]}/{parts[1]}"
        imageUrl={repo.avatarURL}
        barColor={barColor(repo)}
      >
        {#snippet tags()}
          {#if repo.private}
            <Tag shade="yellow" icon="fa-lock">Private</Tag>
          {/if}
          {#if repo.fork}
            <Tag shade="violet" icon="fa-code-fork">Fork</Tag>
          {/if}
          {#if repo.archived}
            <Tag shade="grey" icon="fa-archive">Archived</Tag>
          {/if}
        {/snippet}
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
          <Attribute icon="fa-code-fork">{repo.forks}</Attribute>
          <Attribute icon="fa-clock-o">{repo.updatedAt}</Attribute>
        {/snippet}
      </ObjectItem>
    {/each}
  </ObjectItemList>
</PageShell>
