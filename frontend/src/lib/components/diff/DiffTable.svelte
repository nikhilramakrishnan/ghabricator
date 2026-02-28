<script lang="ts">
  import InlineComment from './InlineComment.svelte';
  import InlineEditor from './InlineEditor.svelte';
  import ContextExpander from './ContextExpander.svelte';
  import { drafts, addDraft, addReplyDraft, removeDraft, updateDraft } from '$lib/stores/inline';
  import { apiPost } from '$lib/api';

  export interface APIChangeset {
    id: number;
    oldName: string;
    newName: string;
    displayPath: string;
    linesAdded: number;
    linesRemoved: number;
    isNew: boolean;
    isDeleted: boolean;
    isRenamed: boolean;
    isBinary: boolean;
    rows: APIDiffRow[];
  }

  export interface APIDiffRow {
    oldNum: number;
    newNum: number;
    oldClass: string;
    newClass: string;
    oldContent: string;
    newContent: string;
    isContext: boolean;
  }

  export interface APIReaction {
    emoji: string;
    count: number;
  }

  export interface APIReviewComment {
    id: number;
    author: string;
    avatarURL: string;
    body: string;
    path: string;
    line: number;
    side: string;
    createdAt: string;
    inReplyTo?: number;
    reactions?: APIReaction[];
  }

  let {
    changeset,
    comments = [],
    owner = '',
    repo = '',
    onNewComment
  }: {
    changeset: APIChangeset;
    comments?: APIReviewComment[];
    owner?: string;
    repo?: string;
    onNewComment?: (path: string, line: number, side: string) => void;
  } = $props();

  let fullWidth = $derived(changeset.isNew || changeset.isDeleted);
  let colSpan = $derived(fullWidth ? 2 : 6);

  // Build threads: group replies under their root comment
  interface CommentThread {
    root: APIReviewComment;
    replies: APIReviewComment[];
  }

  let threadsByKey = $derived.by(() => {
    const map = new Map<string, CommentThread[]>();
    // Index all comments by id for reply-to lookup
    const byId = new Map<number, APIReviewComment>();
    for (const c of comments) byId.set(c.id, c);

    // Separate roots from replies
    const roots: APIReviewComment[] = [];
    const replies: APIReviewComment[] = [];
    for (const c of comments) {
      if (c.inReplyTo && byId.has(c.inReplyTo)) {
        replies.push(c);
      } else {
        roots.push(c);
      }
    }

    // Build thread map keyed by root comment id
    const threadMap = new Map<number, CommentThread>();
    for (const r of roots) {
      threadMap.set(r.id, { root: r, replies: [] });
    }
    for (const r of replies) {
      // Walk up to find root
      let parentId = r.inReplyTo!;
      let parent = byId.get(parentId);
      while (parent && parent.inReplyTo && byId.has(parent.inReplyTo)) {
        parent = byId.get(parent.inReplyTo);
        if (parent) parentId = parent.id;
      }
      const thread = threadMap.get(parentId);
      if (thread) {
        thread.replies.push(r);
      }
    }

    // Sort replies chronologically within each thread
    for (const thread of threadMap.values()) {
      thread.replies.sort((a, b) => a.createdAt.localeCompare(b.createdAt));
    }

    // Group threads by line+side key
    for (const thread of threadMap.values()) {
      const key = `${thread.root.line}:${thread.root.side}`;
      if (!map.has(key)) map.set(key, []);
      map.get(key)!.push(thread);
    }

    return map;
  });

  function lineClick(line: number, side: string) {
    if (onNewComment) {
      onNewComment(changeset.displayPath, line, side);
    }
    addDraft(changeset.displayPath, line, side);
  }

  function hasDraft(path: string, line: number, side: string): boolean {
    let d: import('$lib/stores/inline').DraftComment[] = [];
    drafts.subscribe((v) => (d = v))();
    return d.some((x) => x.path === path && x.line === line && x.side === side);
  }

  function getThreadsForRow(
    row: APIDiffRow
  ): { line: number; side: string; threads: CommentThread[] }[] {
    const result: { line: number; side: string; threads: CommentThread[] }[] = [];
    if (row.newNum > 0) {
      const key = `${row.newNum}:RIGHT`;
      const threads = threadsByKey.get(key);
      if (threads) result.push({ line: row.newNum, side: 'RIGHT', threads });
    }
    if (row.oldNum > 0) {
      const key = `${row.oldNum}:LEFT`;
      const threads = threadsByKey.get(key);
      if (threads) result.push({ line: row.oldNum, side: 'LEFT', threads });
    }
    return result;
  }

  function handleReply(comment: APIReviewComment) {
    addReplyDraft(changeset.displayPath, comment.line, comment.side, comment.id);
  }

  async function handleReaction(commentId: number, emoji: string) {
    try {
      await apiPost('/api/v2/reaction', {
        owner,
        repo,
        commentId,
        emoji
      });
    } catch {
      // ignore errors silently for now
    }
  }
</script>

<div class="diff-wrap">
  <table class="diff-table">
    {#if fullWidth}
      <colgroup>
        <col class="num" style="width:4em" />
        <col class="full" />
      </colgroup>
    {:else}
      <colgroup>
        <col class="num" style="width:4em" />
        <col class="left" />
        <col class="num" style="width:4em" />
        <col class="copy" />
        <col class="right" />
        <col class="cov" />
      </colgroup>
    {/if}
    <tbody>
      {#each changeset.rows as row, i}
        {#if fullWidth}
          {@const lineNum = changeset.isNew ? row.newNum : row.oldNum}
          {@const cls = changeset.isNew ? 'new new-full' : 'old old-full'}
          <tr>
            {#if lineNum > 0}
              <td class="{cls} n" data-n={lineNum}>
                <button class="line-btn" onclick={() => lineClick(lineNum, changeset.isNew ? 'RIGHT' : 'LEFT')}>
                  {lineNum}
                </button>
              </td>
            {:else}
              <td class="n"></td>
            {/if}
            <td class={cls}>{@html changeset.isNew ? row.newContent : row.oldContent}</td>
          </tr>
        {:else}
          <tr>
            {#if row.oldNum > 0}
              <td
                class="{row.oldClass ? row.oldClass + ' n' : 'n'}"
                data-n={row.oldNum}
              >
                <button class="line-btn" onclick={() => lineClick(row.oldNum, 'LEFT')}>
                  {row.oldNum}
                </button>
              </td>
            {:else}
              <td class="n"></td>
            {/if}

            {#if row.oldClass}
              <td class={row.oldClass} data-copy-mode="copy-l">{@html row.oldContent}</td>
            {:else}
              <td data-copy-mode="copy-l">{@html row.oldContent}</td>
            {/if}

            {#if row.newNum > 0}
              <td
                class="{row.newClass ? row.newClass + ' n' : 'n'}"
                data-n={row.newNum}
              >
                <button class="line-btn" onclick={() => lineClick(row.newNum, 'RIGHT')}>
                  {row.newNum}
                </button>
              </td>
            {:else}
              <td class="n"></td>
            {/if}

            <td class="copy"></td>

            {#if row.newClass}
              <td class={row.newClass} colspan="2" data-copy-mode="copy-r"
                >{@html row.newContent}</td
              >
            {:else}
              <td colspan="2" data-copy-mode="copy-r">{@html row.newContent}</td>
            {/if}
          </tr>
        {/if}

        {#each getThreadsForRow(row) as group}
          {#each group.threads as thread}
            <tr class="inline" id="ic-{thread.root.id}">
              {#if fullWidth}
                <td colspan="2">
                  <InlineComment comment={thread.root} onReply={() => handleReply(thread.root)} onReaction={(emoji) => handleReaction(thread.root.id, emoji)} />
                  {#each thread.replies as reply, ri}
                    <InlineComment comment={reply} isReply onReply={ri === thread.replies.length - 1 ? () => handleReply(reply) : undefined} onReaction={(emoji) => handleReaction(reply.id, emoji)} />
                  {/each}
                </td>
              {:else if group.side === 'RIGHT'}
                <td colspan="2"></td>
                <td colspan="4">
                  <InlineComment comment={thread.root} onReply={() => handleReply(thread.root)} onReaction={(emoji) => handleReaction(thread.root.id, emoji)} />
                  {#each thread.replies as reply, ri}
                    <InlineComment comment={reply} isReply onReply={ri === thread.replies.length - 1 ? () => handleReply(reply) : undefined} onReaction={(emoji) => handleReaction(reply.id, emoji)} />
                  {/each}
                </td>
              {:else}
                <td colspan="2">
                  <InlineComment comment={thread.root} onReply={() => handleReply(thread.root)} onReaction={(emoji) => handleReaction(thread.root.id, emoji)} />
                  {#each thread.replies as reply, ri}
                    <InlineComment comment={reply} isReply onReply={ri === thread.replies.length - 1 ? () => handleReply(reply) : undefined} onReaction={(emoji) => handleReaction(reply.id, emoji)} />
                  {/each}
                </td>
                <td colspan="4"></td>
              {/if}
            </tr>
          {/each}
        {/each}

        {#each $drafts.filter((d) => d.path === changeset.displayPath && ((row.newNum > 0 && d.line === row.newNum && d.side === 'RIGHT') || (row.oldNum > 0 && d.line === row.oldNum && d.side === 'LEFT'))) as draft}
          <tr class="inline">
            {#if fullWidth}
              <td colspan="2">
                <InlineEditor
                  path={draft.path} line={draft.line} side={draft.side}
                  onSave={(body) => updateDraft(draft.path, draft.line, draft.side, body)}
                  onCancel={() => removeDraft(draft.path, draft.line, draft.side)}
                />
              </td>
            {:else if draft.side === 'RIGHT'}
              <td colspan="2"></td>
              <td colspan="4">
                <InlineEditor
                  path={draft.path} line={draft.line} side={draft.side}
                  onSave={(body) => updateDraft(draft.path, draft.line, draft.side, body)}
                  onCancel={() => removeDraft(draft.path, draft.line, draft.side)}
                />
              </td>
            {:else}
              <td colspan="2">
                <InlineEditor
                  path={draft.path} line={draft.line} side={draft.side}
                  onSave={(body) => updateDraft(draft.path, draft.line, draft.side, body)}
                  onCancel={() => removeDraft(draft.path, draft.line, draft.side)}
                />
              </td>
              <td colspan="4"></td>
            {/if}
          </tr>
        {/each}
      {/each}
    </tbody>
  </table>
</div>

<style>
  .diff-wrap {
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-family: var(--font-mono);
    font-size: 11px;
    line-height: 1.5;
  }

  td {
    padding: 1px 8px;
    white-space: pre;
    border: none;
  }

  td.n {
    color: var(--text-muted);
    text-align: right;
    user-select: none;
    vertical-align: top;
    padding: 1px 4px;
  }

  td.copy {
    width: 12px;
    padding: 0;
  }

  /* Line number buttons -- invisible until hover */
  .line-btn {
    all: unset;
    cursor: pointer;
    color: inherit;
    display: block;
    width: 100%;
    text-align: right;
  }
  .line-btn:hover {
    color: var(--text-link);
    text-decoration: underline;
  }

  /* Change highlighting */
  :global(td.old),
  :global(td.old-full) {
    background: var(--diff-del-bg);
  }
  :global(td.new),
  :global(td.new-full) {
    background: var(--diff-add-bg);
  }
  :global(td.old.n),
  :global(td.old-full.n) {
    background: var(--diff-del-num-bg);
  }
  :global(td.new.n),
  :global(td.new-full.n) {
    background: var(--diff-add-num-bg);
  }

  /* Show more row */
  :global(tr.show-more td) {
    text-align: center;
    padding: 6px;
    background: var(--bg-hover);
    color: var(--text-muted);
    font-size: 12px;
    cursor: pointer;
  }

  /* Inline comment row */
  tr.inline td {
    padding: 4px 0;
    white-space: normal;
  }
</style>
