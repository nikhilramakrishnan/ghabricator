<script lang="ts">
  import type { APIReviewComment } from '$lib/types';
  import type { APIDiffRow } from './DiffTable.svelte';
  import { formatTimestamp } from '$lib/time';

  const EMOJI_ICONS: Record<string, string> = {
    '+1': 'fa-thumbs-up',
    '-1': 'fa-thumbs-down',
    'laugh': 'fa-smile-o',
    'confused': 'fa-question',
    'heart': 'fa-heart',
    'star': 'fa-star',
    'rocket': 'fa-rocket',
    'eyes': 'fa-eye'
  };

  let {
    comment,
    contextRows = [],
    side = 'RIGHT',
    replies = [],
    onNavigate
  }: {
    comment: APIReviewComment;
    contextRows?: APIDiffRow[];
    side?: string;
    replies?: APIReviewComment[];
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

  function isOld(row: APIDiffRow): boolean {
    return lineClass(row).includes('old');
  }

  function isNew(row: APIDiffRow): boolean {
    return lineClass(row).includes('new');
  }

  function changeMarker(row: APIDiffRow): string {
    if (isOld(row)) return '\u2212';
    if (isNew(row)) return '+';
    return ' ';
  }

  let bodyEl: HTMLDivElement | undefined = $state();
  let overflows = $state(false);
  let expanded = $state(false);

  $effect(() => {
    if (bodyEl) {
      overflows = bodyEl.scrollHeight > 120;
    }
  });
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
        <div class="icwc-row" class:target={isTarget(row)} class:old={isOld(row)} class:new={isNew(row)}>
          <span class="icwc-ln">{ln > 0 ? ln : ''}</span>
          <span class="icwc-marker" class:add={isNew(row)} class:del={isOld(row)}>{changeMarker(row)}</span>
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
    <div class="icwc-body" class:collapsed={overflows && !expanded} bind:this={bodyEl}>
      {@html comment.body}
    </div>
    {#if overflows}
      <button class="icwc-toggle" onclick={() => expanded = !expanded}>
        <i class="fa {expanded ? 'fa-chevron-up' : 'fa-chevron-down'}"></i>
        {expanded ? 'Show Less' : 'Show More'}
      </button>
    {/if}
    {#if comment.reactions?.length}
      <div class="icwc-reactions">
        {#each comment.reactions as r}
          <span class="icwc-pill" >
            <i class="fa {EMOJI_ICONS[r.emoji] ?? 'fa-smile-o'}"></i>
            <span>{r.count}</span>
          </span>
        {/each}
      </div>
    {/if}
    <div class="icwc-actions">
      <button class="icwc-action" onclick={() => onNavigate?.(comment.path, comment.line)}>
        <i class="fa fa-reply"></i> Reply
      </button>
    </div>
  </div>

  {#if replies.length > 0}
    {#each replies as reply}
      <div class="icwc-reply">
        <div class="icwc-meta">
          {#if reply.author.avatarURL}
            <img src={reply.author.avatarURL} alt="" class="icwc-avatar" />
          {/if}
          <strong>{reply.author.login}</strong>
          <span class="icwc-time">{formatTimestamp(reply.createdAt)}</span>
        </div>
        <div class="icwc-body">
          {@html reply.body}
        </div>
        {#if reply.reactions?.length}
          <div class="icwc-reactions">
            {#each reply.reactions as r}
              <span class="icwc-pill" >
                <i class="fa {EMOJI_ICONS[r.emoji] ?? 'fa-smile-o'}"></i>
                <span>{r.count}</span>
              </span>
            {/each}
          </div>
        {/if}
      </div>
    {/each}
  {/if}
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
  .icwc-row.old {
    background: var(--diff-del-bg);
  }
  .icwc-row.new {
    background: var(--diff-add-bg);
  }
  .icwc-row.target {
    background: rgba(255, 212, 59, 0.25);
    box-shadow: inset 3px 0 0 var(--yellow);
  }
  .icwc-row.target.old {
    background: rgba(255, 212, 59, 0.25);
    box-shadow: inset 3px 0 0 var(--red);
  }
  .icwc-row.target.new {
    background: rgba(255, 212, 59, 0.25);
    box-shadow: inset 3px 0 0 var(--green);
  }

  .icwc-ln {
    width: 4em;
    text-align: right;
    padding: 1px 6px 1px 0;
    color: var(--text-muted);
    user-select: none;
    flex-shrink: 0;
  }
  .icwc-marker {
    width: 1.2em;
    text-align: center;
    flex-shrink: 0;
    color: var(--text-muted);
    user-select: none;
  }
  .icwc-marker.add {
    color: var(--green);
    font-weight: 700;
  }
  .icwc-marker.del {
    color: var(--red);
    font-weight: 700;
  }
  .icwc-content {
    padding: 1px 8px;
    flex: 1;
  }

  .icwc-comment {
    border-top: 1px solid var(--border-subtle);
  }
  .icwc-meta {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 12px;
    font-size: 12px;
    background: var(--bg-card-header);
    border-bottom: 1px solid var(--border-subtle);
  }
  .icwc-meta strong {
    color: var(--text);
  }
  .icwc-avatar {
    width: 20px;
    height: 20px;
    border-radius: 3px;
  }
  .icwc-time {
    color: var(--text-muted);
    margin-left: auto;
  }
  .icwc-body {
    padding: 8px 12px;
    font-size: 13px;
    color: var(--text);
    line-height: 1.5;
  }
  .icwc-body.collapsed {
    max-height: 120px;
    overflow: hidden;
    -webkit-mask-image: linear-gradient(to bottom, black 60%, transparent 100%);
    mask-image: linear-gradient(to bottom, black 60%, transparent 100%);
  }
  .icwc-toggle {
    all: unset;
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 4px 12px 8px;
    font-size: 11px;
    color: var(--text-link);
    cursor: pointer;
  }
  .icwc-toggle:hover {
    text-decoration: underline;
  }

  .icwc-reactions {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    padding: 0 12px 8px;
  }

  .icwc-pill {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 2px 8px;
    border-radius: 10px;
    font-size: 11px;
    background: var(--bg-subtle);
    border: 1px solid var(--border-subtle);
    color: var(--text-muted);
  }
  .icwc-pill.reacted {
    background: var(--tag-blue-bg);
    border-color: var(--blue);
    color: var(--blue);
  }

  .icwc-actions {
    padding: 6px 12px;
    border-top: 1px solid var(--border-subtle);
    display: flex;
    gap: 12px;
  }
  .icwc-action {
    all: unset;
    font-size: 11px;
    color: var(--text-muted);
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 4px;
  }
  .icwc-action:hover {
    color: var(--text-link);
  }

  .icwc-reply {
    border-top: 1px solid var(--border-subtle);
    margin-left: 24px;
    border-left: 3px solid var(--blue);
  }
</style>
