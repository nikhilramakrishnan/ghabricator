<script lang="ts">
  import { marked } from 'marked';

  let {
    value = $bindable(''),
    placeholder = 'Write a comment...',
    minRows = 4,
    autofocus = false,
  }: {
    value?: string;
    placeholder?: string;
    minRows?: number;
    autofocus?: boolean;
  } = $props();

  let tab = $state<'write' | 'preview'>('write');
  let textarea: HTMLTextAreaElement | undefined = $state();

  let renderedHTML = $derived.by(() => {
    if (tab !== 'preview') return '';
    return marked.parse(value, { async: false }) as string;
  });

  function autoGrow() {
    if (!textarea) return;
    textarea.style.height = 'auto';
    const lineHeight = 20;
    const minHeight = lineHeight * minRows + 16; // 16 = vertical padding
    textarea.style.height = Math.max(textarea.scrollHeight, minHeight) + 'px';
  }

  $effect(() => {
    if (textarea && autofocus && tab === 'write') {
      textarea.focus();
    }
  });

  $effect(() => {
    // Re-run whenever value changes to auto-grow
    void value;
    if (tab === 'write') {
      // Tick so DOM updates first
      requestAnimationFrame(autoGrow);
    }
  });
</script>

<div class="md-editor">
  <div class="tab-bar">
    <button
      class="tab"
      class:active={tab === 'write'}
      onclick={() => tab = 'write'}
    >
      <i class="fa fa-pencil"></i> Write
    </button>
    <button
      class="tab"
      class:active={tab === 'preview'}
      onclick={() => tab = 'preview'}
    >
      <i class="fa fa-eye"></i> Preview
    </button>
  </div>

  {#if tab === 'write'}
    <textarea
      bind:this={textarea}
      bind:value={value}
      {placeholder}
      rows={minRows}
      oninput={autoGrow}
    ></textarea>
  {:else}
    <div class="preview">
      {#if value.trim()}
        {@html renderedHTML}
      {:else}
        <span class="empty">Nothing to preview</span>
      {/if}
    </div>
  {/if}
</div>

<style>
  .md-editor {
    border: 1px solid var(--border);
    border-radius: 4px;
    overflow: hidden;
    background: var(--bg-card);
  }

  .tab-bar {
    display: flex;
    gap: 0;
    border-bottom: 1px solid var(--border-subtle);
    background: var(--bg-card-header);
    padding: 0 8px;
  }

  .tab {
    display: flex;
    align-items: center;
    gap: 5px;
    padding: 6px 12px;
    border: none;
    background: none;
    color: var(--text-muted);
    font: inherit;
    font-size: 12px;
    font-weight: 500;
    cursor: pointer;
    position: relative;
    white-space: nowrap;
  }

  .tab::after {
    content: '';
    position: absolute;
    bottom: -1px;
    left: 0;
    right: 0;
    height: 2px;
    background: transparent;
  }

  .tab:hover {
    color: var(--text);
  }

  .tab.active {
    color: var(--blue);
    font-weight: 600;
  }

  .tab.active::after {
    background: var(--blue);
  }

  textarea {
    width: 100%;
    border: none;
    padding: 8px 12px;
    font-size: 13px;
    font-family: var(--font-mono, monospace);
    line-height: 20px;
    resize: none;
    outline: none;
    box-sizing: border-box;
    background: transparent;
    color: var(--text);
  }

  .preview {
    padding: 12px;
    font-size: 13px;
    line-height: 1.5;
    color: var(--text);
    overflow-wrap: break-word;
    word-break: break-word;
    overflow-x: auto;
    min-height: 80px;
  }

  .preview :global(h1),
  .preview :global(h2),
  .preview :global(h3) {
    margin: 12px 0 6px;
    line-height: 1.3;
  }
  .preview :global(h1) { font-size: 18px; }
  .preview :global(h2) { font-size: 16px; }
  .preview :global(h3) { font-size: 14px; }

  .preview :global(p) {
    margin: 0 0 8px;
  }

  .preview :global(code) {
    background: var(--bg-subtle);
    padding: 1px 4px;
    border-radius: 3px;
    font-family: var(--font-mono, monospace);
    font-size: 12px;
  }

  .preview :global(pre) {
    background: var(--bg-subtle);
    border: 1px solid var(--border-subtle);
    border-radius: 4px;
    padding: 10px 12px;
    overflow-x: auto;
    margin: 8px 0;
  }

  .preview :global(pre code) {
    background: none;
    padding: 0;
  }

  .preview :global(a) {
    color: var(--text-link);
    text-decoration: none;
  }
  .preview :global(a:hover) {
    text-decoration: underline;
  }

  .preview :global(blockquote) {
    border-left: 3px solid var(--border);
    margin: 8px 0;
    padding: 4px 12px;
    color: var(--text-muted);
  }

  .preview :global(ul),
  .preview :global(ol) {
    margin: 4px 0 8px;
    padding-left: 20px;
  }

  .preview :global(img) {
    max-width: 100%;
  }

  .empty {
    color: var(--text-muted);
    font-style: italic;
  }
</style>
