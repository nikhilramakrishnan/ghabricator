<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import { Button } from '$lib/components/phui';
  import { apiPost } from '$lib/api';
  import { goto } from '$app/navigation';
  import type { APIPasteCreateResponse } from '$lib/types';

  const crumbs = [
    { name: 'Home', href: '/' },
    { name: 'Paste', href: '/paste' },
    { name: 'Create' }
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

<PageShell title="Create Paste" icon="fa-plus">
  {#snippet breadcrumbs()}
    <Breadcrumbs {crumbs} />
  {/snippet}

  <div class="phui-box phui-box-border phui-object-box">
    <form onsubmit={handleSubmit}>
      <div class="phui-form-view" style="padding:16px">
        <div class="aphront-form-control" style="margin-bottom:12px">
          <label class="aphront-form-label" style="display:block;font-weight:bold;margin-bottom:4px;font-size:13px">Title</label>
          <input type="text" bind:value={title} placeholder="Paste title" class="aphront-form-input" style="width:100%;max-width:460px;padding:6px 8px;border:1px solid #c7ccd9;border-radius:3px" />
        </div>

        <div class="aphront-form-control" style="margin-bottom:12px">
          <label class="aphront-form-label" style="display:block;font-weight:bold;margin-bottom:4px;font-size:13px">Language</label>
          <select bind:value={language} class="aphront-form-input" style="padding:6px 8px;border:1px solid #c7ccd9;border-radius:3px">
            {#each languages as lang}
              <option value={lang.ext}>{lang.name}</option>
            {/each}
          </select>
        </div>

        <div class="aphront-form-control" style="margin-bottom:12px">
          <label class="aphront-form-label" style="display:block;font-weight:bold;margin-bottom:4px;font-size:13px">Content</label>
          <textarea
            bind:value={content}
            rows="20"
            required
            placeholder="Paste content here..."
            class="PhabricatorMonospaced aphront-textarea-very-tall"
            style="width:100%;padding:8px;border:1px solid #c7ccd9;border-radius:3px;font-family:monospace;font-size:12px;resize:vertical"
          ></textarea>
        </div>

        <div class="aphront-form-control" style="margin-bottom:16px">
          <label class="aphront-form-label" style="display:block;font-weight:bold;margin-bottom:4px;font-size:13px">Visibility</label>
          <select bind:value={visibility} class="aphront-form-input" style="padding:6px 8px;border:1px solid #c7ccd9;border-radius:3px">
            <option value="secret">Secret</option>
            <option value="public">Public</option>
          </select>
        </div>

        <Button type="submit" color="green" icon="fa-save" disabled={submitting}>
          Create Paste
        </Button>
      </div>
    </form>
  </div>
</PageShell>
