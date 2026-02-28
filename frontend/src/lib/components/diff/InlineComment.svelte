<script lang="ts">
  import type { APIReviewComment } from './DiffTable.svelte';
  import { formatTimestamp } from '$lib/time';

  let {
    comment,
    onReply,
    onDone
  }: {
    comment: APIReviewComment;
    onReply?: () => void;
    onDone?: () => void;
  } = $props();
</script>

<div class="inline-comment">
  <div class="inline-header">
    {#if comment.avatarURL}
      <img src={comment.avatarURL} alt="" class="avatar" />
    {/if}
    <strong>{comment.author}</strong>
    {#if comment.createdAt}
      <span class="time">{formatTimestamp(comment.createdAt)}</span>
    {/if}
  </div>
  <div class="inline-body">
    <div class="remark">
      {@html comment.body}
    </div>
  </div>
  <div class="inline-actions">
    {#if onReply}
      <button class="action-btn" onclick={onReply}>
        <i class="fa fa-reply mrs"></i> Reply
      </button>
    {/if}
    {#if onDone}
      <button class="action-btn" onclick={onDone}>
        <i class="fa fa-check mrs"></i> Done
      </button>
    {/if}
  </div>
</div>

<style>
  .inline-comment {
    margin: 6px 8px;
    border: 1px solid var(--border);
    border-radius: 4px;
    background: var(--bg-card);
    overflow: hidden;
  }

  .inline-header {
    background: var(--bg-card-header);
    padding: 6px 12px;
    font-size: 12px;
    border-bottom: 1px solid var(--border);
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .inline-header strong {
    color: var(--text);
  }

  .inline-header .time {
    color: var(--text-muted);
    margin-left: auto;
  }

  .avatar {
    width: 20px;
    height: 20px;
    border-radius: 3px;
  }

  .inline-body {
    padding: 8px 12px;
    font-size: 13px;
    color: var(--text);
    line-height: 1.5;
  }

  .inline-actions {
    padding: 6px 12px;
    border-top: 1px solid var(--border-subtle);
    display: flex;
    gap: 12px;
  }

  .action-btn {
    all: unset;
    font-size: 11px;
    color: var(--text-muted);
    cursor: pointer;
  }

  .action-btn:hover {
    color: var(--text-link);
  }
</style>
