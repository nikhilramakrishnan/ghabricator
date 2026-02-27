<script lang="ts">
  import type { APICommit } from '$lib/types';
  import { S } from '$lib/strings';

  let {
    commits,
    baseBranch = 'base',
    onRangeChange
  }: {
    commits: APICommit[];
    baseBranch?: string;
    onRangeChange?: (base: string | null, head: string | null) => void;
  } = $props();

  let collapsed = $state(false);
  let compareBase = $state('__base__');
  let compareHead = $state('__latest__');

  function shortSha(sha: string): string {
    return sha.slice(0, 7);
  }

  function firstLine(msg: string): string {
    const nl = msg.indexOf('\n');
    return nl > 0 ? msg.slice(0, nl) : msg;
  }

  function formatDate(dateStr: string): string {
    const d = new Date(dateStr);
    const now = new Date();
    const diff = now.getTime() - d.getTime();
    const days = Math.floor(diff / 86400000);
    if (days === 0) {
      const hours = Math.floor(diff / 3600000);
      if (hours === 0) {
        const mins = Math.floor(diff / 60000);
        return mins <= 1 ? S.diff.justNow : `${mins}m ago`;
      }
      return `${hours}h ago`;
    }
    if (days === 1) return S.diff.yesterday;
    if (days < 30) return `${days}d ago`;
    return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  }

  function handleRangeChange() {
    if (!onRangeChange) return;
    if (compareBase === '__base__' && compareHead === '__latest__') {
      onRangeChange(null, null);
    } else {
      const base = compareBase === '__base__' ? null : compareBase;
      const head = compareHead === '__latest__' ? null : compareHead;
      onRangeChange(base, head);
    }
  }
</script>

<div class="commit-history">
  <div class="ch-header" role="button" tabindex="0" onclick={() => collapsed = !collapsed} onkeydown={(e) => e.key === 'Enter' && (collapsed = !collapsed)}>
    <i class="fa {collapsed ? 'fa-chevron-right' : 'fa-chevron-down'} ch-toggle"></i>
    <i class="fa fa-code-fork ch-icon"></i>
    <span class="ch-title">{S.diff.title}</span>
    <span class="ch-count">{commits.length}</span>
  </div>

  {#if !collapsed}
    <div class="ch-controls">
      <span class="ch-label">{S.diff.showChangesFrom}</span>
      <select class="ch-select" bind:value={compareBase} onchange={handleRangeChange}>
        <option value="__base__">{baseBranch} {S.diff.baseSuffix}</option>
        {#each commits as c, i}
          <option value={c.sha}>{shortSha(c.sha)} — {firstLine(c.message).slice(0, 40)}</option>
        {/each}
      </select>
      <span class="ch-label">{S.diff.to}</span>
      <select class="ch-select" bind:value={compareHead} onchange={handleRangeChange}>
        <option value="__latest__">{S.diff.latest}</option>
        {#each [...commits].reverse() as c}
          <option value={c.sha}>{shortSha(c.sha)} — {firstLine(c.message).slice(0, 40)}</option>
        {/each}
      </select>
    </div>

    <table class="ch-table">
      <thead>
        <tr>
          <th class="ch-th-sha">{S.diff.colSha}</th>
          <th class="ch-th-author">{S.diff.colAuthor}</th>
          <th class="ch-th-msg">{S.diff.colMessage}</th>
          <th class="ch-th-date">{S.diff.colDate}</th>
        </tr>
      </thead>
      <tbody>
        {#each commits as c (c.sha)}
          <tr class="ch-row">
            <td class="ch-sha">
              <code>{shortSha(c.sha)}</code>
            </td>
            <td class="ch-author">
              {#if c.author.avatarURL}
                <img src={c.author.avatarURL} alt="" class="ch-avatar" />
              {/if}
              {c.author.login}
            </td>
            <td class="ch-msg">{firstLine(c.message)}</td>
            <td class="ch-date">{formatDate(c.date)}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>

<style>
  .commit-history {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 4px;
    margin-bottom: 8px;
    overflow: hidden;
  }

  .ch-header {
    background: var(--bg-card-header);
    padding: 6px 12px;
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    user-select: none;
  }

  .ch-toggle {
    font-size: 12px;
    opacity: 0.6;
  }

  .ch-icon {
    color: var(--text-muted);
    font-size: 14px;
  }

  .ch-title {
    font-size: 15px;
    font-weight: 600;
    color: var(--text);
  }

  .ch-count {
    font-size: 11px;
    font-weight: 600;
    background: var(--bg-subtle);
    color: var(--text-muted);
    padding: 1px 6px;
    border-radius: 10px;
    margin-left: 2px;
  }

  .ch-controls {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 12px;
    border-bottom: 1px solid var(--border);
    font-size: 13px;
    flex-wrap: wrap;
    background: var(--bg-subtle);
  }

  .ch-label {
    color: var(--text-muted);
    font-size: 12px;
  }

  .ch-select {
    font-size: 12px;
    padding: 4px 6px;
    border: 1px solid var(--border);
    border-radius: 3px;
    background: var(--bg-card);
    color: var(--text);
    max-width: 260px;
  }

  .ch-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 12px;
  }

  .ch-table thead {
    background: var(--bg-card-header);
  }

  .ch-table th {
    text-align: left;
    padding: 7px 12px;
    font-size: 11px;
    font-weight: 600;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.03em;
    border-bottom: 1px solid var(--border);
  }

  .ch-th-sha { width: 80px; }
  .ch-th-author { width: 140px; }
  .ch-th-msg { }
  .ch-th-date { width: 90px; }

  .ch-row {
    border-bottom: 1px solid var(--border-subtle);
  }
  .ch-row:last-child {
    border-bottom: none;
  }
  .ch-row:hover {
    background: var(--bg-hover);
  }

  .ch-row td {
    padding: 8px 12px;
    vertical-align: middle;
  }

  .ch-sha code {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--text-link);
  }

  .ch-author {
    display: flex;
    align-items: center;
    gap: 5px;
    color: var(--text);
    white-space: nowrap;
  }

  .ch-avatar {
    width: 18px;
    height: 18px;
    border-radius: 3px;
  }

  .ch-msg {
    color: var(--text);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 0;
    width: 100%;
  }

  .ch-date {
    color: var(--text-muted);
    white-space: nowrap;
  }
</style>
