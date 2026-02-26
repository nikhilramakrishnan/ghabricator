<script lang="ts">
  import { apiPost } from '$lib/api';
  import { pendingCount, drafts, clearDrafts } from '$lib/stores/inline';

  let {
    owner,
    repo,
    number
  }: {
    owner: string;
    repo: string;
    number: number;
  } = $props();

  let body = $state('');
  let submitting = $state(false);

  async function submit(event: 'COMMENT' | 'APPROVE' | 'REQUEST_CHANGES') {
    if (submitting) return;
    submitting = true;

    let inlineDrafts: { path: string; line: number; side: string; body: string }[] = [];
    drafts.subscribe((d) => (inlineDrafts = d.filter((x) => x.body.trim())))();

    try {
      await apiPost(`/api/v2/review`, {
        owner,
        repo,
        number,
        body: body.trim(),
        event,
        comments: inlineDrafts
      });
      body = '';
      clearDrafts();
    } catch {
      // TODO: surface error
    } finally {
      submitting = false;
    }
  }
</script>

<div class="mood-review-form">
  <textarea
    bind:value={body}
    placeholder="Leave a review comment..."
    rows="4"
    disabled={submitting}
  ></textarea>
  <div class="form-footer">
    {#if $pendingCount > 0}
      <span class="pending">{$pendingCount} pending comment{$pendingCount === 1 ? '' : 's'}</span>
    {/if}
    <button
      class="mood-btn mood-btn-default"
      disabled={submitting}
      onclick={() => submit('COMMENT')}
    >
      Comment
    </button>
    <button
      class="mood-btn mood-btn-green"
      disabled={submitting}
      onclick={() => submit('APPROVE')}
    >
      <span class="phui-icon-view phui-font-fa fa-check mrs"></span>
      Accept
    </button>
    <button
      class="mood-btn mood-btn-red"
      disabled={submitting}
      onclick={() => submit('REQUEST_CHANGES')}
    >
      <span class="phui-icon-view phui-font-fa fa-times mrs"></span>
      Request Changes
    </button>
  </div>
</div>

<style>
  .mood-review-form {
    border: 1px solid #c7ccd9;
    border-radius: 4px;
    overflow: hidden;
    background: #fff;
  }

  textarea {
    width: 100%;
    border: none;
    padding: 12px;
    font-size: 13px;
    font-family: inherit;
    min-height: 80px;
    resize: vertical;
    outline: none;
    box-sizing: border-box;
    background: transparent;
    color: #292e36;
  }

  .form-footer {
    background: #f6f8fa;
    padding: 8px 12px;
    display: flex;
    align-items: center;
    gap: 8px;
    border-top: 1px solid #e3e4e8;
  }

  .pending {
    font-size: 12px;
    color: #6b748c;
    margin-right: auto;
  }

  .mood-btn {
    padding: 6px 14px;
    border: 1px solid #c7ccd9;
    border-radius: 3px;
    font-size: 12px;
    cursor: pointer;
    font-weight: 600;
    background: #fff;
    color: #464c5c;
  }
  .mood-btn:hover {
    background: #f0f0f0;
  }
  .mood-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .mood-btn-green {
    background: #139543;
    color: #fff;
    border-color: #139543;
  }
  .mood-btn-green:hover {
    background: #117a38;
  }

  .mood-btn-red {
    background: #c0392b;
    color: #fff;
    border-color: #c0392b;
  }
  .mood-btn-red:hover {
    background: #a33025;
  }

  /* Dark mode */
  :global(.phui-theme-dark) .mood-review-form {
    background: #26374c;
    border-color: rgba(255, 255, 255, 0.3);
  }
  :global(.phui-theme-dark) textarea {
    color: rgba(255, 255, 255, 0.9);
  }
  :global(.phui-theme-dark) .form-footer {
    background: #1c293b;
    border-color: rgba(255, 255, 255, 0.3);
  }
  :global(.phui-theme-dark) .mood-btn-default {
    background: #1c293b;
    color: rgba(255, 255, 255, 0.8);
    border-color: rgba(255, 255, 255, 0.3);
  }
  :global(.phui-theme-dark) .mood-btn-default:hover {
    background: #26374c;
  }
</style>
