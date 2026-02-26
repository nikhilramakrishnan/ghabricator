<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import {
    ObjectItemList, ObjectItem, Attribute, Button, InfoView
  } from '$lib/components/phui';
  import type { APIPasteSummary } from '$lib/types';

  let { data } = $props();
  let pastes: APIPasteSummary[] = $derived(data.pastes ?? []);

  const crumbs = [
    { name: 'Home', href: '/' },
    { name: 'Paste' }
  ];
</script>

<PageShell title="Recent Pastes" icon="fa-clipboard">
  {#snippet breadcrumbs()}
    <Breadcrumbs {crumbs} />
  {/snippet}
  {#snippet headerRight()}
    <Button color="green" icon="fa-plus" href="/paste/new">Create Paste</Button>
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
