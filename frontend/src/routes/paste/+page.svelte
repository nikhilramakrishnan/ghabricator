<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import {
    ObjectItemList, ObjectItem, Attribute, Button, InfoView
  } from '$lib/components/phui';
  import type { APIPasteSummary } from '$lib/types';
  import { S } from '$lib/strings';

  let { data } = $props();
  let pastes: APIPasteSummary[] = $derived(data.pastes ?? []);

  const crumbs = [
    { name: S.crumb.home, href: '/' },
    { name: S.paste.title }
  ];
</script>

<PageShell title={S.paste.recentPastes} icon="fa-clipboard">
  {#snippet breadcrumbs()}
    <Breadcrumbs {crumbs} />
  {/snippet}
  {#snippet headerRight()}
    <Button color="green" icon="fa-plus" href="/paste/new">{S.paste.createPaste}</Button>
  {/snippet}

  {#if pastes.length === 0}
    <InfoView icon="fa-inbox">No pastes found. Create one!</InfoView>
  {:else}
    <ObjectItemList>
      {#each pastes as paste}
        <ObjectItem
          title={paste.title}
          href="/paste/{paste.id}"
          icon="fa-file-code-o"
        >
          {#snippet attributes()}
            {#if paste.language}
              <Attribute icon="fa-code">{paste.language}</Attribute>
            {/if}
            <Attribute icon={paste.public ? 'fa-globe' : 'fa-lock'}>
              {paste.public ? 'Public' : 'Secret'}
            </Attribute>
            <Attribute icon="fa-clock-o">{paste.createdAt}</Attribute>
          {/snippet}
        </ObjectItem>
      {/each}
    </ObjectItemList>
  {/if}
</PageShell>
