<script lang="ts">
  let {
    path,
    line,
    side,
    initialBody = '',
    onSave,
    onCancel
  }: {
    path: string;
    line: number;
    side: string;
    initialBody?: string;
    onSave: (body: string) => void;
    onCancel: () => void;
  } = $props();

  // We intentionally capture initialBody once — the editor owns its own draft state.
  let body = $state(initialBody);
</script>

<div class="inline-editor">
  <div class="editor-header">
    <span class="line-ref">Line {line} ({side === 'LEFT' ? 'old' : 'new'})</span>
  </div>
  <textarea
    bind:value={body}
    placeholder="Leave an inline comment..."
    rows="3"
  ></textarea>
  <div class="editor-actions">
    <button class="btn btn-green" onclick={() => onSave(body)} disabled={!body.trim()}>
      Save
    </button>
    <button class="btn" onclick={onCancel}>
      Cancel
    </button>
  </div>
</div>

<style>
  .inline-editor {
    margin: 8px 0 8px 60px;
    border: 1px solid var(--blue);
    border-radius: 4px;
    background: var(--bg-card);
    overflow: hidden;
  }

  .editor-header {
    background: var(--bg-card-header);
    padding: 4px 12px;
    font-size: 11px;
    color: var(--text-muted);
    border-bottom: 1px solid var(--border-subtle);
  }

  textarea {
    width: 100%;
    border: none;
    padding: 8px 12px;
    font-size: 13px;
    font-family: inherit;
    min-height: 60px;
    resize: vertical;
    outline: none;
    box-sizing: border-box;
    background: transparent;
    color: var(--text);
  }

  .editor-actions {
    padding: 6px 12px;
    border-top: 1px solid var(--border-subtle);
    display: flex;
    gap: 8px;
    justify-content: flex-end;
  }

  .btn {
    padding: 4px 12px;
    border: 1px solid var(--border);
    border-radius: 3px;
    font-size: 12px;
    cursor: pointer;
    font-weight: 600;
    background: var(--bg-card);
    color: var(--text);
  }
  .btn:hover {
    background: var(--bg-subtle);
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
</style>
