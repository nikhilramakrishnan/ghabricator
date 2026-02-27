<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import {
    Box, HeaderView, ObjectItemList, ObjectItem, Attribute,
    Tag, InfoView
  } from '$lib/components/phui';
  import type { APIPRSummary } from '$lib/types';

  let { data } = $props();
  let authored: APIPRSummary[] = $derived(data.data.authored ?? []);
  let reviewRequested: APIPRSummary[] = $derived(data.data.reviewRequested ?? []);

  const crumbs = [
    { name: 'Home', href: '/' },
    { name: 'Dashboard' }
  ];
</script>

<PageShell title="Active Revisions" icon="fa-columns">
  {#snippet breadcrumbs()}
    <Breadcrumbs {crumbs} />
  {/snippet}

  <Box border>
    <HeaderView title="Authored" icon="fa-pencil" />
    {#if authored.length === 0}
      <InfoView icon="fa-inbox">No open pull requests.</InfoView>
    {:else}
      <ObjectItemList>
        {#each authored as pr}
          <ObjectItem
            title={pr.title}
            href="/pr/{pr.owner}/{pr.repo}/{pr.number}"
            imageUrl={pr.author.avatarURL}
            barColor={pr.draft ? 'grey' : 'blue'}
          >
            {#snippet tags()}
              {#if pr.draft}
                <Tag shade="grey">Draft</Tag>
              {/if}
              {#if pr.labels}
                {#each pr.labels as label}
                  <Tag shade="blue">{label.name}</Tag>
                {/each}
              {/if}
            {/snippet}
            {#snippet attributes()}
              <Attribute icon="fa-user">{pr.author.login}</Attribute>
              <Attribute icon="fa-github">{pr.owner}/{pr.repo}#{pr.number}</Attribute>
              <Attribute icon="fa-clock-o">{pr.updatedAt}</Attribute>
            {/snippet}
            {#snippet handleIcons()}
              {#if pr.reviewers}
                {#each pr.reviewers as reviewer}
                  {#if reviewer.avatarURL}
                    <img
                      class="reviewer-avatar"
                      src={reviewer.avatarURL}
                      alt={reviewer.login}
                      title={reviewer.login}
                    />
                  {/if}
                {/each}
              {/if}
            {/snippet}
          </ObjectItem>
        {/each}
      </ObjectItemList>
    {/if}
  </Box>

  <Box border>
    <HeaderView title="Review Requested" icon="fa-eye" />
    {#if reviewRequested.length === 0}
      <InfoView icon="fa-inbox">No open pull requests.</InfoView>
    {:else}
      <ObjectItemList>
        {#each reviewRequested as pr}
          <ObjectItem
            title={pr.title}
            href="/pr/{pr.owner}/{pr.repo}/{pr.number}"
            imageUrl={pr.author.avatarURL}
            barColor={pr.draft ? 'grey' : 'blue'}
          >
            {#snippet tags()}
              {#if pr.draft}
                <Tag shade="grey">Draft</Tag>
              {/if}
              {#if pr.labels}
                {#each pr.labels as label}
                  <Tag shade="blue">{label.name}</Tag>
                {/each}
              {/if}
            {/snippet}
            {#snippet attributes()}
              <Attribute icon="fa-user">{pr.author.login}</Attribute>
              <Attribute icon="fa-github">{pr.owner}/{pr.repo}#{pr.number}</Attribute>
              <Attribute icon="fa-clock-o">{pr.updatedAt}</Attribute>
            {/snippet}
            {#snippet handleIcons()}
              {#if pr.reviewers}
                {#each pr.reviewers as reviewer}
                  {#if reviewer.avatarURL}
                    <img
                      class="reviewer-avatar"
                      src={reviewer.avatarURL}
                      alt={reviewer.login}
                      title={reviewer.login}
                    />
                  {/if}
                {/each}
              {/if}
            {/snippet}
          </ObjectItem>
        {/each}
      </ObjectItemList>
    {/if}
  </Box>
</PageShell>

<style>
  .reviewer-avatar {
    width: 20px;
    height: 20px;
    border-radius: 3px;
    display: inline-block;
  }
</style>
