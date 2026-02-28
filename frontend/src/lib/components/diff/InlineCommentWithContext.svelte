<script lang="ts">
  import type { APIReviewComment } from '$lib/types';
  import type { APIDiffRow } from './DiffTable.svelte';
  import { formatTimestamp } from '$lib/time';

  let {
    comment,
    contextRows = [],
    side = 'RIGHT',
    onNavigate
  }: {
    comment: APIReviewComment;
    contextRows?: APIDiffRow[];
    side?: string;
    onNavigate?: (path: string, line: number) => void;
  } = $props();

  function lineNum(row: APIDiffRow): number {
    return side === 'RIGHT' ? row.newNum : row.oldNum;
  }

  function lineContent(row: APIDiffRow): string {
    return side === 'RIGHT' ? row.newContent : row.oldContent;
  }

  function lineClass(row: APIDiffRow): string {
    return side === 'RIGHT' ? row.newClass : row.oldClass;
  }

  function isTarget(row: APIDiffRow): boolean {
    return lineNum(row) === comment.line;
  }
</script>

<div class="icwc">
  <button class="icwc-file" onclick={() => onNavigate?.(comment.path, comment.line)}>
    <i class="fa fa-file-code-o"></i>
    <span class="icwc-path">{comment.path}</span><span class="icwc-line">:{comment.line}</span>
  </button>

  {#if contextRows.length > 0}
    <div class="icwc-code">
      {#each contextRows as row}
        {@const ln = lineNum(row)}
        {@const cls = lineClass(row)}
        <div class="icwc-row" class:target={isTarget(row)} class:old={cls === 'old' || cls === 'old-full'} class:new={cls === 'new' || cls === 'new-full'}>
          <span class="icwc-ln">{ln > 0 ? ln : ''}</span>
          <span class="icwc-content">{@html lineContent(row)}</span>
        </div>
      {/each}
    </div>
  {/if}

  <div class="icwc-comment">
    <div class="icwc-meta">
      {#if comment.author.avatarURL}
        <img src={comment.author.avatarURL} alt="" class="icwc-avatar" />
      {/if}
      <strong>{comment.author.login}</strong>
      <span class="icwc-time">{formatTimestamp(comment.createdAt)}</span>
    </div>
    <div class="icwc-body">{@html comment.body}</div>
  </div>
</div>

<style>
  .icwc {
    border: 1px solid var(--border);
    border-radius: 4px;
    overflow: hidden;
    background: var(--bg-card);
  }

  .icwc-file {
    all: unset;
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 12px;
    font-size: 12px;
    font-weight: 600;
    color: var(--text-link);
    background: var(--bg-card-header);
    border-bottom: 1px solid var(--border-subtle);
    cursor: pointer;
    width: 100%;
    box-sizing: border-box;
  }
  .icwc-file:hover {
    text-decoration: underline;
  }
  .icwc-file i {
    opacity: 0.5;
    color: var(--text-muted);
  }
  .icwc-line {
    opacity: 0.6;
  }

  .icwc-code {
    font-family: var(--font-mono);
    font-size: 11px;
    line-height: 1.5;
    border-bottom: 1px solid var(--border-subtle);
    overflow-x: auto;
  }

  .icwc-row {
    display: flex;
    white-space: pre;
  }
  .icwc-row.target {
    background: var(--bg-hover);
    outline: 1px solid var(--border);
    outline-offset: -1px;
  }
  .icwc-row.old {
    background: var(--diff-del-bg);
  }
  .icwc-row.new {
    background: var(--diff-add-bg);
  }
  .icwc-row.target.old {
    background: var(--diff-del-bg);
    outline: 1px solid var(--border);
  }
  .icwc-row.target.new {
    background: var(--diff-add-bg);
    outline: 1px solid var(--border);
  }

  .icwc-ln {
    width: 4em;
    text-align: right;
    padding: 1px 6px 1px 0;
    color: var(--text-muted);
    user-select: none;
    flex-shrink: 0;
  }
  .icwc-content {
    padding: 1px 8px;
    flex: 1;
  }

  .icwc-comment {
    padding: 8px 12px;
  }
  .icwc-meta {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    margin-bottom: 6px;
  }
  .icwc-meta strong {
    color: var(--text);
  }
  .icwc-avatar {
    width: 18px;
    height: 18px;
    border-radius: 3px;
  }
  .icwc-time {
    color: var(--text-muted);
    margin-left: auto;
  }
  .icwc-body {
    font-size: 13px;
    color: var(--text);
    line-height: 1.5;
  }
</style>
