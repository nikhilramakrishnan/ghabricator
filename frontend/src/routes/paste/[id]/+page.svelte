<script lang="ts">
  import { Breadcrumbs, CurtainLayout } from '$lib/components/layout';
  import { Box, HeaderView, Tag, CurtainBox, PropertyList, ActionList } from '$lib/components/phui';
  import type { APIPasteDetail } from '$lib/types';
  import { S } from '$lib/strings';

  let { data } = $props();
  let paste: APIPasteDetail = $derived(data.paste);

  let crumbs = $derived([
    { name: S.crumb.home, href: '/' },
    { name: S.paste.title, href: '/paste' },
    { name: `P${paste.id.slice(0, 8)}` }
  ]);
</script>

<div class="page-wrapper">
  <Breadcrumbs {crumbs} />

  <div class="page-header">
    <h1 class="page-title">
      <i class="fa fa-clipboard mrs"></i>
      {paste.title}
    </h1>
  </div>

  <CurtainLayout>
    {#each paste.files as file}
      <Box border>
        <HeaderView title={file.filename} icon="fa-file-code-o" />
        {#if file.language}
          <div class="file-lang">
            <Tag shade="blue">{file.language}</Tag>
          </div>
        {/if}
        <div class="source-container">
          <table class="source-table">
            <tbody>
              {#each file.lines as line, i}
                <tr>
                  <th class="line-number"><span>{i + 1}</span></th>
                  <td class="line-code">{@html line}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </Box>
    {/each}

    {#snippet curtain()}
      {#if paste.owner?.login}
        <CurtainBox title={S.common.author}>
          <div class="author-info">
            {#if paste.owner.avatarURL}
              <img src={paste.owner.avatarURL} alt="" class="author-avatar" />
            {/if}
            {paste.owner.login}
          </div>
        </CurtainBox>
      {/if}

      <CurtainBox title={S.common.details}>
        <PropertyList items={[
          { label: S.paste.visibility, value: paste.public ? S.common.public : S.common.secret },
          { label: S.common.created, value: paste.createdAt },
          ...(paste.updatedAt !== paste.createdAt ? [{ label: S.common.updated, value: paste.updatedAt }] : [])
        ]} />
      </CurtainBox>

      <CurtainBox title={S.common.actions}>
        <ActionList actions={[
          { label: S.repos.viewOnGitHub, href: paste.htmlURL, icon: 'fa-github' }
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
    padding: 12px 0;
    margin-bottom: 8px;
  }

  .page-title {
    font-size: 20px;
    font-weight: 600;
    color: var(--text);
    margin: 0;
  }

  .file-lang {
    padding: 4px 12px;
  }

  .source-container {
    overflow-x: auto;
  }

  .source-table {
    width: 100%;
    border-collapse: collapse;
    font-family: var(--font-mono);
    font-size: 12px;
    line-height: 1.6;
  }

  .line-number {
    width: 1%;
    min-width: 44px;
    padding: 0 8px;
    text-align: right;
    color: var(--text-muted);
    user-select: none;
    white-space: nowrap;
    vertical-align: top;
    background: var(--bg-subtle);
    border-right: 1px solid var(--border-subtle);
  }

  .line-code {
    padding: 0 12px;
    white-space: pre;
    color: var(--text);
  }

  .author-info {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
  }

  .author-avatar {
    width: 24px;
    height: 24px;
    border-radius: 3px;
  }
</style>
