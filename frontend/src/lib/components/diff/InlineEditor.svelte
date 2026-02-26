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
  <div class="inline-editor-header">
    <span class="line-ref">Line {line} ({side === 'LEFT' ? 'old' : 'new'})</span>
  </div>
  <textarea
    bind:value={body}
    placeholder="Leave an inline comment..."
    rows="3"
  ></textarea>
  <div class="inline-editor-actions">
    <button class="mood-btn mood-btn-green" onclick={() => onSave(body)} disabled={!body.trim()}>
      Save
    </button>
    <button class="mood-btn mood-btn-default" onclick={onCancel}>
      Cancel
    </button>
  </div>
</div>

<style>
  .inline-editor {
    margin: 8px 0 8px 60px;
    border: 1px solid #136cb2;
    border-radius: 4px;
    background: #fff;
    overflow: hidden;
  }

  .inline-editor-header {
    background: #f8f9fc;
    padding: 4px 12px;
    font-size: 11px;
    color: #6b748c;
    border-bottom: 1px solid #e3e4e8;
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
    color: #292e36;
  }

  .inline-editor-actions {
    padding: 6px 12px;
    border-top: 1px solid #e3e4e8;
    display: flex;
    gap: 8px;
    justify-content: flex-end;
  }

  .mood-btn {
    padding: 4px 12px;
    border: 1px solid #c7ccd9;
    border-radius: 3px;
    font-size: 12px;
    cursor: pointer;
    font-weight: 600;
    background: #fff;
    color: #464c5c;
  }
  .mood-btn:hover {
    background: #f6f8fa;
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

  /* Dark mode */
  :global(.phui-theme-dark) .inline-editor {
    background: #26374c;
    border-color: #136cb2;
  }
  :global(.phui-theme-dark) .inline-editor-header {
    background: #1c293b;
    border-color: rgba(255, 255, 255, 0.3);
    color: rgba(255, 255, 255, 0.5);
  }
  :global(.phui-theme-dark) textarea {
    color: rgba(255, 255, 255, 0.9);
  }
  :global(.phui-theme-dark) .inline-editor-actions {
    border-color: rgba(255, 255, 255, 0.3);
  }
  :global(.phui-theme-dark) .mood-btn-default {
    background: #1c293b;
    color: rgba(255, 255, 255, 0.8);
    border-color: rgba(255, 255, 255, 0.3);
  }
</style>
