<script lang="ts">
  import { Breadcrumbs, CurtainLayout } from '$lib/components/layout';
  import { Box, HeaderView, Tag, CurtainBox, PropertyList, ActionList } from '$lib/components/phui';
  import type { HeraldRule } from '$lib/types';

  let { data } = $props();
  let rule: HeraldRule = $derived(data.rule);

  let crumbs = $derived([
    { name: 'Home', href: '/' },
    { name: 'Actions', href: '/actions' },
    { name: rule.name }
  ]);

  const conditionLabels: Record<string, string> = {
    title_contains: 'Title contains',
    body_contains: 'Body contains',
    file_matches: 'File matches',
    author_is: 'Author is',
    label_is: 'Label is',
  };

  const actionLabels: Record<string, string> = {
    add_reviewer: 'Add reviewer',
    add_label: 'Add label',
  };
</script>

<div class="page-wrapper">
  <Breadcrumbs {crumbs} />

  <div class="page-header">
    <h1 class="page-title">
      <i class="fa fa-bullhorn mrs"></i>
      {rule.name}
    </h1>
    <div>
      {#if rule.disabled}
        <Tag shade="grey">Disabled</Tag>
      {:else}
        <Tag shade="green">Active</Tag>
      {/if}
    </div>
  </div>

  <CurtainLayout>
    <Box border>
      <HeaderView title="Conditions" icon="fa-filter" />
      <div class="section-body">
        <p class="section-desc">
          {rule.must_match_all ? 'When all of these conditions are met:' : 'When any of these conditions are met:'}
        </p>
        {#if rule.conditions.length === 0}
          <p class="empty-text">No conditions configured.</p>
        {:else}
          <ul class="rule-list">
            {#each rule.conditions as cond}
              <li class="rule-item">
                <i class="fa fa-chevron-right rule-icon"></i>
                <strong>{conditionLabels[cond.type] ?? cond.type}</strong>
                <code class="rule-value">{cond.value}</code>
              </li>
            {/each}
          </ul>
        {/if}
      </div>
    </Box>

    <Box border>
      <HeaderView title="Actions" icon="fa-bolt" />
      <div class="section-body">
        {#if rule.actions.length === 0}
          <p class="empty-text">No actions configured.</p>
        {:else}
          <ul class="rule-list">
            {#each rule.actions as action}
              <li class="rule-item">
                <i class="fa fa-arrow-right rule-icon"></i>
                <strong>{actionLabels[action.type] ?? action.type}</strong>
                <code class="rule-value">{action.value}</code>
              </li>
            {/each}
          </ul>
        {/if}
      </div>
    </Box>

    {#snippet curtain()}
      <CurtainBox title="Details">
        <PropertyList items={[
          { label: 'Author', value: rule.author_login },
          { label: 'Match', value: rule.must_match_all ? 'All Conditions' : 'Any Condition' },
          { label: 'Status', value: rule.disabled ? 'Disabled' : 'Active' },
          { label: 'Created', value: rule.created_at },
          ...(rule.updated_at !== rule.created_at ? [{ label: 'Updated', value: rule.updated_at }] : [])
        ]} />
      </CurtainBox>

      <CurtainBox title="Actions">
        <ActionList actions={[
          { label: 'Edit Rule', href: `/actions/new?edit=${rule.id}`, icon: 'fa-pencil' }
        ]} />
      </CurtainBox>
    {/snippet}
  </CurtainLayout>
</div>

<style>
  .page-wrapper {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 16px;
  }

  .page-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 0;
    margin-bottom: 8px;
  }

  .page-title {
    font-size: 20px;
    font-weight: 600;
    color: var(--text);
    margin: 0;
  }

  .section-body {
    padding: 0 16px 16px;
  }

  .section-desc {
    font-size: 13px;
    color: var(--text-muted);
    margin: 8px 0;
  }

  .empty-text {
    font-size: 13px;
    color: var(--text-muted);
  }

  .rule-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .rule-item {
    padding: 6px 0;
    border-bottom: 1px solid var(--border-subtle);
    font-size: 13px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .rule-icon {
    color: var(--text-muted);
    font-size: 10px;
  }

  .rule-value {
    background: var(--bg-subtle);
    padding: 2px 6px;
    border-radius: 3px;
    font-size: 12px;
  }
</style>
