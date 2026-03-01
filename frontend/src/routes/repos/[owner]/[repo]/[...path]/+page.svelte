<script lang="ts">
  import { Breadcrumbs, CurtainLayout } from '$lib/components/layout';
  import { CurtainBox, PropertyList, ActionList, InfoView } from '$lib/components/phui';
  import type { APIRepoEntry, APIRepoInfo, APIBlameRange } from '$lib/types';
  import { apiFetch } from '$lib/api';
  import { S } from '$lib/strings';

  let { data } = $props();
  let owner = $derived(data.owner);
  let repo = $derived(data.repo);
  let path = $derived(data.path);
  let ref = $derived(data.ref);
  let mode = $derived(data.mode);

  // Blame state
  let blameEnabled = $state(false);
  let blameRanges: APIBlameRange[] = $state([]);
  let blameLoading = $state(false);

  // Build a line→range lookup when blame is loaded
  let blameByLine = $derived.by(() => {
    const map = new Map<number, APIBlameRange>();
    for (const r of blameRanges) {
      for (let i = r.startLine; i <= r.endLine; i++) {
        map.set(i, r);
      }
    }
    return map;
  });

  // Track which lines are the first in their blame range (show annotation there)
  let blameFirstLines = $derived.by(() => {
    const set = new Set<number>();
    for (const r of blameRanges) {
      set.add(r.startLine);
    }
    return set;
  });

  async function toggleBlame() {
    if (blameEnabled) {
      blameEnabled = false;
      return;
    }
    if (blameRanges.length > 0) {
      blameEnabled = true;
      return;
    }
    blameLoading = true;
    try {
      const qs = new URLSearchParams();
      if (ref) qs.set('ref', ref);
      if (path) qs.set('path', path);
      const query = qs.toString() ? `?${qs.toString()}` : '';
      const resp = await apiFetch<{ ranges: APIBlameRange[] }>(`/api/repo/${owner}/${repo}/blame${query}`);
      blameRanges = resp.ranges;
      blameEnabled = true;
    } catch (e) {
      console.error('Blame fetch failed:', e);
    } finally {
      blameLoading = false;
    }
  }

  // Breadcrumbs
  let crumbs = $derived.by(() => {
    const items: { name: string; href?: string }[] = [
      { name: S.crumb.home, href: '/' },
      { name: S.repos.title, href: '/repos' },
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

  function buildCurtainProps(info: APIRepoInfo) {
    return [
      { label: S.repos.visibility, value: info.private ? 'Private' : 'Public' },
      { label: S.repos.stars, value: String(info.stars) },
      { label: S.repos.forks, value: String(info.forks) }
    ];
  }

  function timeAgo(dateStr: string): string {
    const date = new Date(dateStr);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const days = Math.floor(diffMs / (1000 * 60 * 60 * 24));
    if (days === 0) return 'today';
    if (days === 1) return 'yesterday';
    if (days < 30) return `${days}d ago`;
    if (days < 365) return `${Math.floor(days / 30)}mo ago`;
    return `${Math.floor(days / 365)}y ago`;
  }

  // Assign alternating colors to commits for visual grouping
  let commitColorMap = $derived.by(() => {
    const map = new Map<string, number>();
    let idx = 0;
    for (const r of blameRanges) {
      if (!map.has(r.commitOID)) {
        map.set(r.commitOID, idx++ % 6);
      }
    }
    return map;
  });
</script>

<div class="page-wrapper">
  <Breadcrumbs {crumbs} />

{#if mode === 'error'}
  <InfoView severity="warning" icon="fa-exclamation-triangle">
    Could not load path. It may not exist or you may not have access.
  </InfoView>

{:else if mode === 'tree' && data.tree}
  {@const tree = data.tree}
  {@const info = tree.repoInfo}
  <CurtainLayout>
    <table class="file-table">
      <thead>
        <tr class="file-table-head">
          <th class="file-th">Name</th>
          <th class="file-th file-th-type">Type</th>
          <th class="file-th file-th-size">Size</th>
        </tr>
      </thead>
      <tbody>
        {#if path}
          <tr class="file-row">
            <td class="file-td">
              <a href={parentHref()} class="file-link">
                <i class="fa fa-level-up mrs"></i>..
              </a>
            </td>
            <td class="file-td file-td-muted"></td>
            <td class="file-td"></td>
          </tr>
        {/if}
        {#each tree.entries as entry}
          <tr class="file-row">
            <td class="file-td">
              <a href={entryHref(entry)} class="file-link">
                <i class="fa {entry.type === 'dir' ? 'fa-folder' : 'fa-file-o'} mrs"
                  class:icon-folder={entry.type === 'dir'}
                  class:icon-file={entry.type !== 'dir'}
                ></i>
                {entry.name}
              </a>
            </td>
            <td class="file-td file-td-muted">{entry.type === 'dir' ? 'Directory' : 'File'}</td>
            <td class="file-td file-td-size">
              {entry.type === 'file' ? formatSize(entry.size) : ''}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>

    {#snippet curtain()}
      <CurtainBox title={S.common.actions}>
        <ActionList actions={[
          { label: S.repos.viewOnGitHub, href: info.htmlURL, icon: 'fa-github' }
        ]} />
      </CurtainBox>
      <CurtainBox title={S.common.details}>
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
    <div class="file-viewer-box">
      <div class="file-header">
        <h1 class="file-title">
          <i class="fa {isImage(file.name) ? 'fa-file-image-o' : 'fa-file-code-o'} mrs"></i>
          {file.name}
          <span class="file-size">{formatSize(file.size)}</span>
        </h1>
        <div class="file-header-actions">
          {#if file.lines && !isImage(file.name)}
            <button
              class="header-btn"
              class:header-btn-active={blameEnabled}
              onclick={toggleBlame}
              disabled={blameLoading}
              title="Toggle blame annotations"
            >
              {#if blameLoading}
                <i class="fa fa-spinner fa-spin"></i>
              {:else}
                <i class="fa fa-history mrs"></i>Blame
              {/if}
            </button>
          {/if}
          {#if file.htmlURL}
            <a href={file.htmlURL} target="_blank" rel="noopener" class="header-btn">
              <i class="fa fa-github mrs"></i>GitHub
            </a>
          {/if}
        </div>
      </div>

      {#if file.rawURL && isImage(file.name)}
        <div class="image-preview">
          <img src={file.rawURL} alt={file.name} class="preview-img" />
        </div>
      {:else if file.lines}
        <div class="source-container">
          <table class="source-table" class:blame-active={blameEnabled}>
            <tbody>
              {#each file.lines as line, i}
                {@const lineNum = i + 1}
                {@const blameRange = blameEnabled ? blameByLine.get(lineNum) : null}
                {@const isFirst = blameEnabled && blameFirstLines.has(lineNum)}
                {@const colorIdx = blameRange ? commitColorMap.get(blameRange.commitOID) ?? 0 : 0}
                <tr class:blame-boundary={isFirst && lineNum > 1}>
                  {#if blameEnabled}
                    <td class="blame-gutter blame-color-{colorIdx}" class:blame-first={isFirst}>
                      {#if isFirst && blameRange}
                        <div class="blame-info">
                          <span class="blame-commit" title={blameRange.message}>{blameRange.commitShort}</span>
                          <span class="blame-author" title={blameRange.authorLogin || blameRange.authorName}>
                            {blameRange.authorLogin || blameRange.authorName}
                          </span>
                          <span class="blame-date">{timeAgo(blameRange.authoredDate)}</span>
                        </div>
                      {/if}
                    </td>
                  {/if}
                  <th class="line-number"><span>{lineNum}</span></th>
                  <td class="line-code">{@html line}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </div>

    {#snippet curtain()}
      <CurtainBox title={S.common.actions}>
        <ActionList actions={[
          { label: S.repos.viewOnGitHub, href: info.htmlURL, icon: 'fa-github' }
        ]} />
      </CurtainBox>
      <CurtainBox title={S.common.details}>
        <PropertyList items={buildCurtainProps(info)} />
      </CurtainBox>
    {/snippet}
  </CurtainLayout>
{/if}
</div>

<style>
  .page-wrapper {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 16px;
  }

  /* File tree table */
  .file-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 13px;
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 4px;
  }

  .file-table-head {
    border-bottom: 2px solid var(--border);
    text-align: left;
  }

  .file-th {
    padding: 8px 12px;
    color: var(--text);
    font-weight: 600;
  }

  .file-th-type {
    width: 80px;
  }

  .file-th-size {
    width: 100px;
    text-align: right;
  }

  .file-row {
    border-bottom: 1px solid var(--border-subtle);
  }

  .file-row:hover {
    background: var(--bg-hover);
  }

  .file-td {
    padding: 6px 12px;
  }

  .file-td-muted {
    color: var(--text-muted);
  }

  .file-td-size {
    text-align: right;
    color: var(--text-muted);
  }

  .file-link {
    text-decoration: none;
  }

  .icon-folder {
    color: var(--orange);
  }

  .icon-file {
    color: var(--text-muted);
  }

  /* File viewer */
  .file-viewer-box {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 4px;
  }

  .file-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 16px;
    border-bottom: 1px solid var(--border-subtle);
    background: var(--bg-card-header);
  }

  .file-title {
    font-size: 15px;
    font-weight: 600;
    color: var(--text);
    margin: 0;
  }

  .file-size {
    font-weight: normal;
    font-size: 12px;
    color: var(--text-muted);
    margin-left: 8px;
  }

  .file-header-actions {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .header-btn {
    font-size: 12px;
    padding: 4px 10px;
    border: 1px solid var(--border);
    border-radius: 3px;
    color: var(--text);
    text-decoration: none;
    background: var(--bg-card);
    cursor: pointer;
    white-space: nowrap;
  }

  .header-btn:hover {
    background: var(--bg-hover);
    text-decoration: none;
  }

  .header-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .header-btn-active {
    background: var(--blue);
    color: #fff;
    border-color: var(--blue);
  }

  .header-btn-active:hover {
    background: var(--blue);
    opacity: 0.9;
  }

  .image-preview {
    padding: 16px;
    text-align: center;
  }

  .preview-img {
    max-width: 100%;
    border: 1px solid var(--border-subtle);
    border-radius: 4px;
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

  /* Blame gutter */
  .blame-gutter {
    width: 1%;
    min-width: 200px;
    max-width: 240px;
    padding: 0;
    vertical-align: top;
    user-select: none;
    background: var(--bg-subtle);
    border-right: 1px solid var(--border-subtle);
    overflow: hidden;
  }

  .blame-boundary td,
  .blame-boundary th {
    border-top: 1px solid var(--border);
  }

  .blame-first .blame-info {
    display: flex;
    align-items: baseline;
    gap: 6px;
    padding: 0 8px;
    font-family: var(--font-mono);
    font-size: 11px;
    line-height: 1.6;
    white-space: nowrap;
    overflow: hidden;
  }

  .blame-commit {
    color: var(--blue);
    font-weight: 600;
    flex-shrink: 0;
  }

  .blame-author {
    color: var(--text);
    overflow: hidden;
    text-overflow: ellipsis;
    flex-shrink: 1;
    min-width: 0;
  }

  .blame-date {
    color: var(--text-muted);
    flex-shrink: 0;
    margin-left: auto;
  }

  /* Alternating commit colors via left border */
  .blame-color-0 { border-left: 3px solid var(--blue); }
  .blame-color-1 { border-left: 3px solid var(--green); }
  .blame-color-2 { border-left: 3px solid var(--orange); }
  .blame-color-3 { border-left: 3px solid var(--violet); }
  .blame-color-4 { border-left: 3px solid var(--red); }
  .blame-color-5 { border-left: 3px solid var(--yellow); }
</style>
