<script lang="ts">
  import InlineComment from './InlineComment.svelte';
  import InlineEditor from './InlineEditor.svelte';
  import ContextExpander from './ContextExpander.svelte';
  import { drafts, addDraft, removeDraft, updateDraft } from '$lib/stores/inline';

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

  export interface APIReviewComment {
    id: number;
    author: string;
    avatarURL: string;
    body: string;
    path: string;
    line: number;
    side: string;
    createdAt: string;
  }

  let {
    changeset,
    comments = [],
    onNewComment
  }: {
    changeset: APIChangeset;
    comments?: APIReviewComment[];
    onNewComment?: (path: string, line: number, side: string) => void;
  } = $props();

  // Index comments by line+side for rendering after rows
  let commentsByKey = $derived.by(() => {
    const map = new Map<string, APIReviewComment[]>();
    for (const c of comments) {
      const key = `${c.line}:${c.side}`;
      if (!map.has(key)) map.set(key, []);
      map.get(key)!.push(c);
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

  function getCommentsForRow(
    row: APIDiffRow
  ): { line: number; side: string; comments: APIReviewComment[] }[] {
    const result: { line: number; side: string; comments: APIReviewComment[] }[] = [];
    if (row.newNum > 0) {
      const key = `${row.newNum}:RIGHT`;
      const cmts = commentsByKey.get(key);
      if (cmts) result.push({ line: row.newNum, side: 'RIGHT', comments: cmts });
    }
    if (row.oldNum > 0) {
      const key = `${row.oldNum}:LEFT`;
      const cmts = commentsByKey.get(key);
      if (cmts) result.push({ line: row.oldNum, side: 'LEFT', comments: cmts });
    }
    return result;
  }
</script>

<div class="changeset-view-content">
  <table
    class="differential-diff remarkup-code PhabricatorMonospaced diff-2up"
  >
    <colgroup>
      <col class="num" style="width:4em" />
      <col class="left" />
      <col class="num" style="width:4em" />
      <col class="copy" />
      <col class="right" />
      <col class="cov" />
    </colgroup>
    <tbody>
      {#each changeset.rows as row, i}
        <tr>
          <!-- Old line number -->
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

          <!-- Old content -->
          {#if row.oldClass}
            <td class={row.oldClass} data-copy-mode="copy-l">{@html row.oldContent}</td>
          {:else}
            <td data-copy-mode="copy-l">{@html row.oldContent}</td>
          {/if}

          <!-- New line number -->
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

          <!-- Copy column -->
          <td class="copy"></td>

          <!-- New content -->
          {#if row.newClass}
            <td class={row.newClass} colspan="2" data-copy-mode="copy-r"
              >{@html row.newContent}</td
            >
          {:else}
            <td colspan="2" data-copy-mode="copy-r">{@html row.newContent}</td>
          {/if}
        </tr>

        <!-- Inline comments for this row -->
        {#each getCommentsForRow(row) as group}
          {#each group.comments as comment}
            <tr class="inline">
              <td colspan="6">
                <InlineComment {comment} />
              </td>
            </tr>
          {/each}
        {/each}

        <!-- Draft editors for this row -->
        {#each $drafts.filter((d) => d.path === changeset.displayPath && ((row.newNum > 0 && d.line === row.newNum && d.side === 'RIGHT') || (row.oldNum > 0 && d.line === row.oldNum && d.side === 'LEFT'))) as draft}
          <tr class="inline">
            <td colspan="6">
              <InlineEditor
                path={draft.path}
                line={draft.line}
                side={draft.side}
                onSave={(body) => {
                  updateDraft(draft.path, draft.line, draft.side, body);
                }}
                onCancel={() => removeDraft(draft.path, draft.line, draft.side)}
              />
            </td>
          </tr>
        {/each}
      {/each}
    </tbody>
  </table>
</div>

<style>
  .changeset-view-content {
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-family: ui-monospace, 'Cascadia Code', 'SF Mono', Menlo, Consolas, monospace;
    font-size: 12px;
    line-height: 1.5;
    table-layout: fixed;
  }

  td {
    padding: 1px 8px;
    white-space: pre;
    border: none;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  td.n {
    color: #6b748c;
    text-align: right;
    user-select: none;
    vertical-align: top;
    padding: 1px 4px;
  }

  td.copy {
    width: 12px;
    padding: 0;
  }

  /* Line number buttons — invisible until hover */
  .line-btn {
    all: unset;
    cursor: pointer;
    color: inherit;
    display: block;
    width: 100%;
    text-align: right;
  }
  .line-btn:hover {
    color: #136cb2;
    text-decoration: underline;
  }

  /* Change highlighting */
  :global(td.old),
  :global(td.old-full) {
    background: rgba(251, 175, 175, 0.3);
  }
  :global(td.new),
  :global(td.new-full) {
    background: rgba(151, 234, 151, 0.3);
  }
  :global(td.old.n),
  :global(td.old-full.n) {
    background: rgba(251, 175, 175, 0.2);
  }
  :global(td.new.n),
  :global(td.new-full.n) {
    background: rgba(151, 234, 151, 0.2);
  }

  /* Show more row */
  :global(tr.show-more td) {
    text-align: center;
    padding: 6px;
    background: rgba(55, 55, 55, 0.04);
    color: #6b748c;
    font-size: 12px;
    cursor: pointer;
  }

  /* Inline comment row */
  tr.inline td {
    padding: 4px 0;
    white-space: normal;
  }
</style>
