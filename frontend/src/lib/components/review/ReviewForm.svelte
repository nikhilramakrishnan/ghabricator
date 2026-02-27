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

<div class="review-form">
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
      class="btn"
      disabled={submitting}
      onclick={() => submit('COMMENT')}
    >
      Comment
    </button>
    <button
      class="btn btn-green"
      disabled={submitting}
      onclick={() => submit('APPROVE')}
    >
      <i class="fa fa-check mrs"></i>
      Accept
    </button>
    <button
      class="btn btn-red"
      disabled={submitting}
      onclick={() => submit('REQUEST_CHANGES')}
    >
      <i class="fa fa-times mrs"></i>
      Request Changes
    </button>
  </div>
</div>

<style>
  .review-form {
    border: 1px solid var(--border);
    border-radius: 4px;
    overflow: hidden;
    background: var(--bg-card);
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
    color: var(--text);
  }

  .form-footer {
    background: var(--bg-subtle);
    padding: 8px 12px;
    display: flex;
    align-items: center;
    gap: 8px;
    border-top: 1px solid var(--border-subtle);
  }

  .pending {
    font-size: 12px;
    color: var(--text-muted);
    margin-right: auto;
  }

  .btn {
    padding: 6px 14px;
    border: 1px solid var(--border);
    border-radius: 3px;
    font-size: 12px;
    cursor: pointer;
    font-weight: 600;
    background: var(--bg-card);
    color: var(--text);
  }
  .btn:hover {
    background: var(--bg-hover);
  }
  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .btn-green {
    background: var(--green);
    color: var(--text-on-dark);
    border-color: var(--green);
  }
  .btn-green:hover {
    background: var(--green-hover);
  }

  .btn-red {
    background: var(--red);
    color: var(--text-on-dark);
    border-color: var(--red);
  }
  .btn-red:hover {
    background: var(--red-hover);
  }
</style>
