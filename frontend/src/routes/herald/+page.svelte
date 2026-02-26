<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import {
    ObjectItemList, ObjectItem, Attribute, Tag, Button, InfoView
  } from '$lib/components/phui';
  import type { HeraldRule } from '$lib/types';

  let { data } = $props();
  let rules: HeraldRule[] = $derived(data.rules ?? []);

  const crumbs = [
    { name: 'Home', href: '/' },
    { name: 'Herald' }
  ];

  function ruleColor(rule: HeraldRule): string {
    if (rule.disabled) return 'grey';
    return 'green';
  }
</script>

<PageShell title="Herald Rules" icon="fa-bullhorn">
  {#snippet breadcrumbs()}
    <Breadcrumbs {crumbs} />
  {/snippet}
  {#snippet headerRight()}
    <Button color="green" icon="fa-plus" href="/herald/new">New Rule</Button>
  {/snippet}

  {#if rules.length === 0}
    <InfoView icon="fa-inbox">No Herald rules configured.</InfoView>
  {:else}
    <ObjectItemList>
      {#each rules as rule}
        <ObjectItem
          title={rule.name}
          href="/herald/{rule.id}"
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
</PageShell>
