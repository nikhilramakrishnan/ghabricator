<script lang="ts">
  import type { APIReviewComment } from './DiffTable.svelte';
  import ReactionPicker from './ReactionPicker.svelte';
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
    isReply = false,
    onReply,
    onDone,
    onReaction
  }: {
    comment: APIReviewComment;
    isReply?: boolean;
    onReply?: () => void;
    onDone?: () => void;
    onReaction?: (emoji: string) => void;
  } = $props();

  let pickerOpen = $state(false);
</script>

<div class="inline-comment" class:is-reply={isReply}>
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
  {#if comment.reactions?.length}
    <div class="reaction-pills">
      {#each comment.reactions as r}
        <button
          class="reaction-pill"
          onclick={() => onReaction?.(r.emoji)}
          title={r.emoji}
        >
          <i class="fa {EMOJI_ICONS[r.emoji] ?? 'fa-smile-o'}"></i>
          <span class="rp-count">{r.count}</span>
        </button>
      {/each}
    </div>
  {/if}
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
    {#if onReaction}
      <div class="picker-anchor">
        <button class="action-btn" title="Add reaction" onclick={() => pickerOpen = !pickerOpen}>
          <i class="fa fa-plus mrs"></i>
        </button>
        {#if pickerOpen}
          <ReactionPicker
            onPick={(emoji) => onReaction(emoji)}
            onClose={() => pickerOpen = false}
          />
        {/if}
      </div>
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

  .inline-comment.is-reply {
    margin-left: 32px;
    border-left: 3px solid var(--blue);
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

  .reaction-pills {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    padding: 0 12px 6px;
  }

  .reaction-pill {
    all: unset;
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 2px 8px;
    border-radius: 10px;
    font-size: 11px;
    background: var(--bg-subtle);
    border: 1px solid var(--border-subtle);
    color: var(--text-muted);
    cursor: pointer;
    transition: background 0.1s;
  }
  .reaction-pill:hover {
    background: var(--bg-hover);
    color: var(--text);
  }
  .reaction-pill.reacted {
    background: var(--tag-blue-bg);
    border-color: var(--blue);
    color: var(--blue);
  }
  .rp-count {
    font-weight: 600;
  }

  .inline-actions {
    padding: 6px 12px;
    border-top: 1px solid var(--border-subtle);
    display: flex;
    gap: 12px;
    align-items: center;
  }

  .picker-anchor {
    position: relative;
    margin-left: auto;
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
