<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import {
    ObjectItemList, ObjectItem, Attribute, Tag, Button, InfoView, TabView
  } from '$lib/components/phui';
  import type { Tab } from '$lib/components/phui';
  import type { HeraldRule, APIWorkflowRun } from '$lib/types';

  let { data } = $props();
  let rules: HeraldRule[] = $derived(data.rules ?? []);
  let runs: APIWorkflowRun[] = $derived(data.runs ?? []);

  const crumbs = [
    { name: 'Home', href: '/' },
    { name: 'Actions' }
  ];

  let activeTab = $state('builds');

  let tabs: Tab[] = $derived([
    { key: 'builds', label: 'Builds', icon: 'fa-cog', count: runs.length },
    { key: 'rules', label: 'Rules', icon: 'fa-bullhorn', count: rules.length },
  ]);

  function ruleColor(rule: HeraldRule): string {
    if (rule.disabled) return 'grey';
    return 'green';
  }

  function runBarColor(run: APIWorkflowRun): string {
    if (run.conclusion === 'success') return 'green';
    if (run.conclusion === 'failure') return 'red';
    if (run.status === 'in_progress' || run.status === 'queued' || run.status === 'pending') return 'yellow';
    return 'grey';
  }

  function runStatusShade(run: APIWorkflowRun): string {
    if (run.conclusion === 'success') return 'green';
    if (run.conclusion === 'failure') return 'red';
    if (run.status === 'in_progress' || run.status === 'queued' || run.status === 'pending') return 'yellow';
    return 'grey';
  }

  function runStatusLabel(run: APIWorkflowRun): string {
    if (run.conclusion) return run.conclusion;
    return run.status;
  }

  function formatDuration(ms: number): string {
    const secs = Math.round(ms / 1000);
    if (secs < 60) return `${secs}s`;
    const m = Math.floor(secs / 60);
    const s = secs % 60;
    return s > 0 ? `${m}m ${s}s` : `${m}m`;
  }
</script>

<PageShell title="Actions" icon="fa-cog">
  {#snippet breadcrumbs()}
    <Breadcrumbs {crumbs} />
  {/snippet}
  {#snippet headerRight()}
    <Button color="green" icon="fa-plus" href="/actions/new">New Rule</Button>
  {/snippet}

  <TabView {tabs} bind:active={activeTab} />

  {#if activeTab === 'rules'}
    {#if rules.length === 0}
      <InfoView icon="fa-inbox">No Herald rules configured.</InfoView>
    {:else}
      <ObjectItemList>
        {#each rules as rule}
          <ObjectItem
            title={rule.name}
            href="/actions/{rule.id}"
            icon="fa-bullhorn"
            barColor={ruleColor(rule)}
          >
            {#snippet tags()}
              {#if rule.disabled}
                <Tag shade="grey">Disabled</Tag>
              {:else}
                <Tag shade="green">Active</Tag>
              {/if}
              <Tag shade={rule.must_match_all ? 'blue' : 'yellow'}>
                {rule.must_match_all ? 'Match All' : 'Match Any'}
              </Tag>
            {/snippet}
            {#snippet attributes()}
              <Attribute icon="fa-user">{rule.author_login}</Attribute>
              <Attribute icon="fa-filter">{rule.conditions.length} condition{rule.conditions.length !== 1 ? 's' : ''}</Attribute>
              <Attribute icon="fa-bolt">{rule.actions.length} action{rule.actions.length !== 1 ? 's' : ''}</Attribute>
              <Attribute icon="fa-clock-o">{rule.created_at}</Attribute>
            {/snippet}
          </ObjectItem>
        {/each}
      </ObjectItemList>
    {/if}
  {:else if activeTab === 'builds'}
    {#if runs.length === 0}
      <InfoView icon="fa-inbox">No workflow runs found.</InfoView>
    {:else}
      <ObjectItemList>
        {#each runs as run}
          <ObjectItem
            title={run.name}
            href={run.htmlURL}
            imageUrl={run.actor.avatarURL}
            barColor={runBarColor(run)}
          >
            {#snippet tags()}
              <Tag shade={runStatusShade(run)}>{runStatusLabel(run)}</Tag>
            {/snippet}
            {#snippet attributes()}
              <Attribute icon="fa-subtitle">{run.displayTitle}</Attribute>
              <Attribute icon="fa-user">{run.actor.login}</Attribute>
              <Attribute icon="fa-code-fork">{run.branch}</Attribute>
              <Attribute icon="fa-clock-o">{formatDuration(run.durationMs)}</Attribute>
              <Attribute icon="fa-database">{run.repoOwner}/{run.repoName}</Attribute>
            {/snippet}
          </ObjectItem>
        {/each}
      </ObjectItemList>
    {/if}
  {/if}
</PageShell>
