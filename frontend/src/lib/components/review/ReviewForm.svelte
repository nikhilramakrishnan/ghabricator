<script lang="ts">
  import { apiPost } from '$lib/api';
  import { S } from '$lib/strings';
  import { pendingCount, drafts, clearDrafts } from '$lib/stores/inline';
  import { user } from '$lib/stores/auth';

  let {
    owner,
    repo,
    number,
    merged = false,
    prState = 'open',
    authorLogin = '',
    approved = false,
  }: {
    owner: string;
    repo: string;
    number: number;
    merged?: boolean;
    prState?: string;
    authorLogin?: string;
    approved?: boolean;
  } = $props();

  let isAuthor = $derived($user?.login === authorLogin);

  let body = $state('');
  let action = $state<'COMMENT' | 'APPROVE' | 'REQUEST_CHANGES'>('COMMENT');
  let submitting = $state(false);

  let submitLabel = $derived(
    action === 'APPROVE' ? 'Accept & Comment'
    : action === 'REQUEST_CHANGES' ? 'Request Changes'
    : 'Comment'
  );
  let submitColor = $derived(
    action === 'APPROVE' ? 'green'
    : action === 'REQUEST_CHANGES' ? 'red'
    : 'default'
  );

  async function handleSubmit() {
    if (submitting) return;
    submitting = true;

    let inlineDrafts: { path: string; line: number; side: string; body: string }[] = [];
    drafts.subscribe((d) => (inlineDrafts = d.filter((x) => x.body.trim())))();

    try {
      await apiPost(`/api/v2/review`, {
        owner, repo, number,
        body: body.trim(),
        event: action,
        comments: inlineDrafts
      });
      body = '';
      action = 'COMMENT';
      clearDrafts();
    } catch {
      // TODO: surface error
    } finally {
      submitting = false;
    }
  }

  async function handleMerge() {
    if (submitting) return;
    submitting = true;
    try {
      await apiPost('/api/v2/merge', { owner, repo, number, mergeMethod: 'squash' });
      window.location.reload();
    } catch (e: unknown) {
      alert(e instanceof Error ? e.message : S.pr.mergeFailed);
    } finally {
      submitting = false;
    }
  }

  async function handleClose(newState: string) {
    if (submitting) return;
    submitting = true;
    try {
      await apiPost('/api/v2/close', { owner, repo, number, state: newState });
      window.location.reload();
    } catch (e: unknown) {
      alert(e instanceof Error ? e.message : S.pr.actionFailed);
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
    {#if !isAuthor}
      <div class="action-picker">
        <label class="action-option" class:selected={action === 'COMMENT'}>
          <input type="radio" bind:group={action} value="COMMENT" />
          <i class="fa fa-comment"></i>
          <span>Comment</span>
        </label>
        <label class="action-option" class:selected={action === 'APPROVE'}>
          <input type="radio" bind:group={action} value="APPROVE" />
          <i class="fa fa-check"></i>
          <span>Accept</span>
        </label>
        <label class="action-option" class:selected={action === 'REQUEST_CHANGES'}>
          <input type="radio" bind:group={action} value="REQUEST_CHANGES" />
          <i class="fa fa-times"></i>
          <span>Reject</span>
        </label>
      </div>
    {/if}
    <div class="footer-right">
      {#if $pendingCount > 0}
        <span class="pending">{$pendingCount} pending</span>
      {/if}
      <button
        class="btn btn-{submitColor}"
        disabled={submitting}
        onclick={handleSubmit}
      >
        {submitLabel}
      </button>
      {#if !merged}
        {#if prState !== 'closed'}
          <button
            class="btn btn-land"
            disabled={submitting || !approved}
            title={approved ? '' : 'Requires approval'}
            onclick={handleMerge}
          >
            <i class="fa fa-check-circle mrs"></i> {S.pr.landRevision}
          </button>
          <button class="btn btn-muted" disabled={submitting} onclick={() => handleClose('closed')}>
            {S.pr.close}
          </button>
        {:else}
          <button class="btn btn-green" disabled={submitting} onclick={() => handleClose('open')}>
            {S.pr.reopen}
          </button>
        {/if}
      {/if}
    </div>
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
    padding: 10px 12px;
    display: flex;
    align-items: center;
    gap: 12px;
    border-top: 1px solid var(--border-subtle);
    flex-wrap: wrap;
  }

  .action-picker {
    display: flex;
    gap: 2px;
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 4px;
    overflow: hidden;
  }

  .action-option {
    display: flex;
    align-items: center;
    gap: 5px;
    padding: 5px 10px;
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;
    color: var(--text-muted);
    user-select: none;
    border-right: 1px solid var(--border-subtle);
  }
  .action-option:last-child {
    border-right: none;
  }
  .action-option input {
    display: none;
  }
  .action-option:hover {
    background: var(--bg-hover);
    color: var(--text);
  }
  .action-option.selected {
    background: var(--bg-hover);
    color: var(--text);
  }

  .footer-right {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-left: auto;
  }

  .pending {
    font-size: 11px;
    color: var(--text-muted);
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
    white-space: nowrap;
  }
  .btn:hover { background: var(--bg-hover); }
  .btn:disabled { opacity: 0.5; cursor: not-allowed; }

  .btn-green {
    background: var(--green);
    color: var(--text-on-dark);
    border-color: var(--green);
  }
  .btn-green:hover { background: var(--green-hover); }

  .btn-red {
    background: var(--red);
    color: var(--text-on-dark);
    border-color: var(--red);
  }
  .btn-red:hover { background: var(--red-hover); }

  .btn-land {
    background: var(--green);
    color: var(--text-on-dark);
    border-color: var(--green);
  }
  .btn-land:hover { background: var(--green-hover); }
  .btn-land:disabled {
    background: var(--bg-card);
    color: var(--text-muted);
    border-color: var(--border);
  }

  .btn-muted {
    color: var(--text-muted);
    border-color: var(--border);
  }
  .btn-muted:hover { color: var(--text); }

  .btn-default {
    /* inherits base .btn styles */
  }
</style>
