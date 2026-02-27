<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import { Button, Box } from '$lib/components/phui';
  import { apiPost } from '$lib/api';
  import { goto } from '$app/navigation';
  import type { APIPasteCreateResponse } from '$lib/types';
  import { S } from '$lib/strings';

  const crumbs = [
    { name: S.crumb.home, href: '/' },
    { name: S.paste.title, href: '/paste' },
    { name: S.paste.createPaste }
  ];

  const languages = [
    { ext: 'txt', name: 'Plain Text' },
    { ext: 'go', name: 'Go' },
    { ext: 'py', name: 'Python' },
    { ext: 'js', name: 'JavaScript' },
    { ext: 'ts', name: 'TypeScript' },
    { ext: 'rs', name: 'Rust' },
    { ext: 'c', name: 'C' },
    { ext: 'cpp', name: 'C++' },
    { ext: 'java', name: 'Java' },
    { ext: 'rb', name: 'Ruby' },
    { ext: 'php', name: 'PHP' },
    { ext: 'sh', name: 'Shell' },
    { ext: 'sql', name: 'SQL' },
    { ext: 'html', name: 'HTML' },
    { ext: 'css', name: 'CSS' },
    { ext: 'json', name: 'JSON' },
    { ext: 'yaml', name: 'YAML' },
    { ext: 'md', name: 'Markdown' },
    { ext: 'xml', name: 'XML' },
    { ext: 'diff', name: 'Diff' },
  ];

  let title = $state('');
  let language = $state('txt');
  let content = $state('');
  let visibility = $state('secret');
  let submitting = $state(false);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (!content.trim() || submitting) return;
    submitting = true;
    try {
      const resp = await apiPost<APIPasteCreateResponse>('/api/paste', {
        title: title || 'Untitled Paste',
        language,
        content,
        public: visibility === 'public'
      });
      goto(resp.url);
    } catch (err: unknown) {
      alert(err instanceof Error ? err.message : 'Failed to create paste');
    } finally {
      submitting = false;
    }
  }
</script>

<PageShell title={S.paste.createPaste} icon="fa-plus">
  {#snippet breadcrumbs()}
    <Breadcrumbs {crumbs} />
  {/snippet}

  <Box border>
    <form onsubmit={handleSubmit}>
      <div class="form-body">
        <div class="form-group">
          <label class="form-label" for="paste-title">Title</label>
          <input id="paste-title" type="text" bind:value={title} placeholder="Paste title" class="form-input" />
        </div>

        <div class="form-group">
          <label class="form-label" for="paste-lang">Language</label>
          <select id="paste-lang" bind:value={language} class="form-input form-select">
            {#each languages as lang}
              <option value={lang.ext}>{lang.name}</option>
            {/each}
          </select>
        </div>

        <div class="form-group">
          <label class="form-label" for="paste-content">Content</label>
          <textarea
            id="paste-content"
            bind:value={content}
            rows="20"
            required
            placeholder="Paste content here..."
            class="form-input form-textarea"
          ></textarea>
        </div>

        <div class="form-group">
          <label class="form-label" for="paste-vis">Visibility</label>
          <select id="paste-vis" bind:value={visibility} class="form-input form-select">
            <option value="secret">Secret</option>
            <option value="public">Public</option>
          </select>
        </div>

        <Button type="submit" color="green" icon="fa-save" disabled={submitting}>
          Create Paste
        </Button>
      </div>
    </form>
  </Box>
</PageShell>

<style>
  .form-body {
    padding: 16px;
  }

  .form-group {
    margin-bottom: 12px;
  }

  .form-label {
    display: block;
    font-weight: bold;
    margin-bottom: 4px;
    font-size: 13px;
    color: var(--text);
  }

  .form-input {
    padding: 6px 8px;
    border: 1px solid var(--border);
    border-radius: 3px;
    font-size: 13px;
    color: var(--text);
    background: var(--bg-card);
  }

  .form-input:focus {
    outline: none;
    border-color: var(--blue);
  }

  input.form-input {
    width: 100%;
    max-width: 460px;
  }

  .form-textarea {
    width: 100%;
    font-family: var(--font-mono);
    font-size: 12px;
    resize: vertical;
    padding: 8px;
  }
</style>
