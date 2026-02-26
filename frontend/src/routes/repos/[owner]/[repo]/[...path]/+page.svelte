<script lang="ts">
  import { Breadcrumbs, CurtainLayout } from '$lib/components/layout';
  import { CurtainBox, PropertyList, Tag, ActionList, InfoView } from '$lib/components/phui';
  import type { APIRepoEntry, APIRepoInfo } from '$lib/types';

  let { data } = $props();
  let owner = $derived(data.owner);
  let repo = $derived(data.repo);
  let path = $derived(data.path);
  let ref = $derived(data.ref);
  let mode = $derived(data.mode);

  // Build breadcrumbs from path segments
  let crumbs = $derived.by(() => {
    const items: { name: string; href?: string }[] = [
      { name: 'Home', href: '/' },
      { name: 'Repositories', href: '/repos' },
      { name: `${owner}/${repo}`, href: `/repos/${owner}/${repo}` }
    ];
    if (path) {
      const segments = path.split('/');
      for (let i = 0; i < segments.length; i++) {
        const partial = segments.slice(0, i + 1).join('/');
        if (i < segments.length - 1) {
          const qs = ref ? `?ref=${ref}` : '';
          items.push({ name: segments[i], href: `/repos/${owner}/${repo}/${partial}${qs}` });
        } else {
          items.push({ name: segments[i] });
        }
      }
    }
    return items;
  });

  // Helpers
  function entryHref(entry: APIRepoEntry): string {
    const qs = ref ? `?ref=${ref}` : '';
    return `/repos/${owner}/${repo}/${entry.path}${qs}`;
  }

  function parentHref(): string {
    if (!path) return '';
    const segments = path.split('/');
    const parent = segments.slice(0, -1).join('/');
    const qs = ref ? `?ref=${ref}` : '';
    return parent ? `/repos/${owner}/${repo}/${parent}${qs}` : `/repos/${owner}/${repo}${qs}`;
  }

  function formatSize(bytes: number): string {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  }

  function isImage(name: string): boolean {
    return /\.(png|jpg|jpeg|gif|webp|svg|ico|bmp)$/i.test(name);
  }

  // Curtain for sidebar
  function buildCurtainProps(info: APIRepoInfo) {
    return [
      { label: 'Visibility', value: info.private ? 'Private' : 'Public' },
      { label: 'Stars', value: String(info.stars) },
      { label: 'Forks', value: String(info.forks) }
    ];
  }
</script>

<div class="phui-two-column-view">
  <div class="phui-two-column-container">
    <Breadcrumbs {crumbs} />
    <div class="phui-two-column-content">

{#if mode === 'error'}
  <InfoView severity="warning" icon="fa-exclamation-triangle">
    Could not load path. It may not exist or you may not have access.
  </InfoView>

{:else if mode === 'tree' && data.tree}
  {@const tree = data.tree}
  {@const info = tree.repoInfo}
  <CurtainLayout>
    <table class="phui-oi-list-view" style="width:100%;border-collapse:collapse;font-size:13px">
      <thead>
        <tr style="border-bottom:2px solid #c7ccd9;text-align:left">
          <th style="padding:8px 12px">Name</th>
          <th style="padding:8px 12px;width:80px">Type</th>
          <th style="padding:8px 12px;width:100px;text-align:right">Size</th>
        </tr>
      </thead>
      <tbody>
        {#if path}
          <tr style="border-bottom:1px solid #e3e4e8">
            <td style="padding:6px 12px">
              <a href={parentHref()} style="text-decoration:none;color:#136CB2">
                <span class="phui-icon-view phui-font-fa fa-level-up mrs"></span>..
              </a>
            </td>
            <td style="padding:6px 12px;color:#6b748c"></td>
            <td style="padding:6px 12px"></td>
          </tr>
        {/if}
        {#each tree.entries as entry}
          <tr style="border-bottom:1px solid #e3e4e8">
            <td style="padding:6px 12px">
              <a href={entryHref(entry)} style="text-decoration:none;color:#136CB2">
                <span class="phui-icon-view phui-font-fa {entry.type === 'dir' ? 'fa-folder' : 'fa-file-o'} mrs"
                  style="color:{entry.type === 'dir' ? '#8C6E00' : '#6b748c'}"
                ></span>
                {entry.name}
              </a>
            </td>
            <td style="padding:6px 12px;color:#6b748c">{entry.type === 'dir' ? 'Directory' : 'File'}</td>
            <td style="padding:6px 12px;text-align:right;color:#6b748c">
              {entry.type === 'file' ? formatSize(entry.size) : ''}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>

    {#snippet curtain()}
      <CurtainBox title="Actions">
        <ActionList actions={[
          { label: 'View on GitHub', href: info.htmlURL, icon: 'fa-github' }
        ]} />
      </CurtainBox>
      <CurtainBox title="Details">
        <PropertyList items={[
          ...buildCurtainProps(info),
          ...(info.description ? [{ label: 'About', value: info.description }] : [])
        ]} />
      </CurtainBox>
    {/snippet}
  </CurtainLayout>

{:else if mode === 'file' && data.file}
  {@const fileResp = data.file}
  {@const file = fileResp.file}
  {@const info = fileResp.repoInfo}
  <CurtainLayout>
    <!-- File header -->
    <div class="phui-box phui-box-border phui-object-box">
      <div class="phui-header-shell">
        <div class="phui-header-view" style="display:flex;align-items:center;justify-content:space-between">
          <h1 class="phui-header-header">
            <span class="phui-icon-view phui-font-fa {isImage(file.name) ? 'fa-file-image-o' : 'fa-file-code-o'}"></span>
            {file.name}
            <span style="font-weight:normal;font-size:12px;color:#6b748c;margin-left:8px">{formatSize(file.size)}</span>
          </h1>
          {#if file.htmlURL}
            <a href={file.htmlURL} target="_blank" rel="noopener" class="mood-btn mood-btn-default" style="font-size:12px;padding:4px 10px">
              <span class="phui-icon-view phui-font-fa fa-github mrs"></span>GitHub
            </a>
          {/if}
        </div>
      </div>

      {#if file.rawURL && isImage(file.name)}
        <div style="padding:16px;text-align:center">
          <img src={file.rawURL} alt={file.name} style="max-width:100%;border:1px solid #e3e4e8;border-radius:4px" />
        </div>
      {:else if file.lines}
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
      {/if}
    </div>

    {#snippet curtain()}
      <CurtainBox title="Actions">
        <ActionList actions={[
          { label: 'View on GitHub', href: info.htmlURL, icon: 'fa-github' }
        ]} />
      </CurtainBox>
      <CurtainBox title="Details">
        <PropertyList items={buildCurtainProps(info)} />
      </CurtainBox>
    {/snippet}
  </CurtainLayout>
{/if}

    </div>
  </div>
</div>
