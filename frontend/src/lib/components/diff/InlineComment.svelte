<script lang="ts">
  import type { APIReviewComment } from './DiffTable.svelte';

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

<div class="mood-inline-comment">
  <div class="mood-inline-header">
    {#if comment.avatarURL}
      <img src={comment.avatarURL} alt="" class="avatar" />
    {/if}
    <strong>{comment.author}</strong>
    {#if comment.createdAt}
      <span class="time">{comment.createdAt}</span>
    {/if}
  </div>
  <div class="mood-inline-body">
    <div class="phabricator-remarkup">
      {@html comment.body}
    </div>
  </div>
  <div class="mood-inline-actions">
    {#if onReply}
      <button class="inline-action" onclick={onReply}>
        <span class="phui-icon-view phui-font-fa fa-reply"></span> Reply
      </button>
    {/if}
    {#if onDone}
      <button class="inline-action" onclick={onDone}>
        <span class="phui-icon-view phui-font-fa fa-check"></span> Done
      </button>
    {/if}
  </div>
</div>

<style>
  .mood-inline-comment {
    margin: 8px 0 8px 60px;
    border: 1px solid #c7ccd9;
    border-radius: 4px;
    background: #fff;
    overflow: hidden;
  }

  .mood-inline-header {
    background: #f8f9fc;
    padding: 6px 12px;
    font-size: 12px;
    border-bottom: 1px solid #c7ccd9;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .mood-inline-header strong {
    color: #464c5c;
  }

  .mood-inline-header .time {
    color: #6b748c;
    margin-left: auto;
  }

  .avatar {
    width: 20px;
    height: 20px;
    border-radius: 3px;
  }

  .mood-inline-body {
    padding: 8px 12px;
    font-size: 13px;
    color: #464c5c;
    line-height: 1.5;
  }

  .mood-inline-actions {
    padding: 6px 12px;
    border-top: 1px solid #e3e4e8;
    display: flex;
    gap: 12px;
  }

  .inline-action {
    all: unset;
    font-size: 11px;
    color: #6b748c;
    cursor: pointer;
  }

  .inline-action:hover {
    color: #136cb2;
  }

  /* Dark mode */
  :global(.phui-theme-dark) .mood-inline-comment {
    background: #26374c;
    border-color: rgba(255, 255, 255, 0.3);
  }
  :global(.phui-theme-dark) .mood-inline-header {
    background: #1c293b;
    border-color: rgba(255, 255, 255, 0.3);
  }
  :global(.phui-theme-dark) .mood-inline-header strong {
    color: rgba(255, 255, 255, 0.9);
  }
  :global(.phui-theme-dark) .mood-inline-body {
    color: rgba(255, 255, 255, 0.8);
  }
  :global(.phui-theme-dark) .mood-inline-actions {
    border-color: rgba(255, 255, 255, 0.3);
  }
  :global(.phui-theme-dark) .inline-action {
    color: rgba(255, 255, 255, 0.6);
  }
</style>
