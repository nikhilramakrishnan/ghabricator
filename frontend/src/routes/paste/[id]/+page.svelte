<script lang="ts">
  import { Breadcrumbs, CurtainLayout } from '$lib/components/layout';
  import { Box, HeaderView, Tag, CurtainBox, PropertyList, ActionList } from '$lib/components/phui';
  import type { APIPasteDetail } from '$lib/types';

  let { data } = $props();
  let paste: APIPasteDetail = $derived(data.paste);

  let crumbs = $derived([
    { name: 'Home', href: '/' },
    { name: 'Paste', href: '/paste' },
    { name: `P${paste.id.slice(0, 8)}` }
  ]);
</script>

<div class="phui-two-column-view">
  <div class="phui-two-column-container">
    <Breadcrumbs {crumbs} />

    <div class="phui-two-column-header">
      <div class="phui-header-view">
        <div class="phui-header-shell">
          <h1 class="phui-header-header">
            <span class="phui-header-icon phui-icon-view phui-font-fa fa-clipboard"></span>
            {paste.title}
          </h1>
        </div>
      </div>
    </div>

    <div class="phui-two-column-content">
      <CurtainLayout>
        {#each paste.files as file}
          <Box border>
            <HeaderView title={file.filename} icon="fa-file-code-o" />
            {#if file.language}
              <div style="padding:4px 12px">
                <Tag shade="blue">{file.language}</Tag>
              </div>
            {/if}
            <div class="phabricator-source-code-container">
              <table class="phabricator-source-code-view remarkup-code PhabricatorMonospaced chroma">
                {#each file.lines as line, i}
                  <tr>
                    <th class="phabricator-source-line"><span>{i + 1}</span></th>
                    <td class="phabricator-source-code">{@html line}</td>
                  </tr>
                {/each}
              </table>
            </div>
          </Box>
        {/each}

        {#snippet curtain()}
          {#if paste.owner?.login}
            <CurtainBox title="Author">
              <div style="display:flex;align-items:center;gap:8px;font-size:13px">
                {#if paste.owner.avatarURL}
                  <img src={paste.owner.avatarURL} alt="" style="width:24px;height:24px;border-radius:3px" />
                {/if}
                {paste.owner.login}
              </div>
            </CurtainBox>
          {/if}

          <CurtainBox title="Details">
            <PropertyList items={[
              { label: 'Visibility', value: paste.public ? 'Public' : 'Secret' },
              { label: 'Created', value: paste.createdAt },
              ...(paste.updatedAt !== paste.createdAt ? [{ label: 'Updated', value: paste.updatedAt }] : [])
            ]} />
          </CurtainBox>

          <CurtainBox title="Actions">
            <ActionList actions={[
              { label: 'View on GitHub', href: paste.htmlURL, icon: 'fa-github' }
            ]} />
          </CurtainBox>
        {/snippet}
      </CurtainLayout>
    </div>
  </div>
</div>
