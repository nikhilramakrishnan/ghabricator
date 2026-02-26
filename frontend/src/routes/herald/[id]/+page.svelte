<script lang="ts">
  import { Breadcrumbs, CurtainLayout } from '$lib/components/layout';
  import { Box, HeaderView, Tag, CurtainBox, PropertyList, ActionList } from '$lib/components/phui';
  import type { HeraldRule } from '$lib/types';

  let { data } = $props();
  let rule: HeraldRule = $derived(data.rule);

  let crumbs = $derived([
    { name: 'Home', href: '/' },
    { name: 'Herald', href: '/herald' },
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

<div class="phui-two-column-view">
  <div class="phui-two-column-container">
    <Breadcrumbs {crumbs} />

    <div class="phui-two-column-header">
      <div class="phui-header-view">
        <div class="phui-header-shell" style="display:flex;align-items:center;justify-content:space-between">
          <h1 class="phui-header-header">
            <span class="phui-header-icon phui-icon-view phui-font-fa fa-bullhorn"></span>
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
      </div>
    </div>

    <div class="phui-two-column-content">
      <CurtainLayout>
        <Box border>
          <HeaderView title="Conditions" icon="fa-filter" />
          <div style="padding:0 16px 16px">
            <p style="font-size:13px;color:#6b748c;margin:8px 0">
              {rule.must_match_all ? 'When all of these conditions are met:' : 'When any of these conditions are met:'}
            </p>
            {#if rule.conditions.length === 0}
              <p style="font-size:13px;color:#92969d">No conditions configured.</p>
            {:else}
              <ul style="list-style:none;padding:0;margin:0">
                {#each rule.conditions as cond}
                  <li style="padding:6px 0;border-bottom:1px solid #e3e4e8;font-size:13px;display:flex;align-items:center;gap:8px">
                    <span class="phui-icon-view phui-font-fa fa-chevron-right" style="color:#6b748c;font-size:10px"></span>
                    <strong>{conditionLabels[cond.type] ?? cond.type}</strong>
                    <code style="background:#f7f7f7;padding:2px 6px;border-radius:3px;font-size:12px">{cond.value}</code>
                  </li>
                {/each}
              </ul>
            {/if}
          </div>
        </Box>

        <Box border>
          <HeaderView title="Actions" icon="fa-bolt" />
          <div style="padding:0 16px 16px">
            {#if rule.actions.length === 0}
              <p style="font-size:13px;color:#92969d">No actions configured.</p>
            {:else}
              <ul style="list-style:none;padding:0;margin:0">
                {#each rule.actions as action}
                  <li style="padding:6px 0;border-bottom:1px solid #e3e4e8;font-size:13px;display:flex;align-items:center;gap:8px">
                    <span class="phui-icon-view phui-font-fa fa-arrow-right" style="color:#6b748c;font-size:10px"></span>
                    <strong>{actionLabels[action.type] ?? action.type}</strong>
                    <code style="background:#f7f7f7;padding:2px 6px;border-radius:3px;font-size:12px">{action.value}</code>
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
              { label: 'Edit Rule', href: `/herald/new?edit=${rule.id}`, icon: 'fa-pencil' }
            ]} />
          </CurtainBox>
        {/snippet}
      </CurtainLayout>
    </div>
  </div>
</div>
