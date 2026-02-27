<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import { Button, Box } from '$lib/components/phui';
  import { apiPost } from '$lib/api';
  import { goto } from '$app/navigation';
  import type { HeraldCondition, HeraldAction } from '$lib/types';

  const crumbs = [
    { name: 'Home', href: '/' },
    { name: 'Actions', href: '/actions' },
    { name: 'New Rule' }
  ];

  const conditionTypes = [
    { value: 'title_contains', label: 'Title contains' },
    { value: 'body_contains', label: 'Body contains' },
    { value: 'file_matches', label: 'File matches' },
    { value: 'author_is', label: 'Author is' },
    { value: 'label_is', label: 'Label is' },
  ];

  const actionTypes = [
    { value: 'add_reviewer', label: 'Add reviewer' },
    { value: 'add_label', label: 'Add label' },
  ];

  let name = $state('');
  let mustMatchAll = $state(true);
  let conditions: HeraldCondition[] = $state([{ type: 'title_contains', value: '' }]);
  let actions: HeraldAction[] = $state([{ type: 'add_reviewer', value: '' }]);
  let submitting = $state(false);

  function addCondition() {
    conditions = [...conditions, { type: 'title_contains', value: '' }];
  }

  function removeCondition(i: number) {
    conditions = conditions.filter((_, idx) => idx !== i);
  }

  function addAction() {
    actions = [...actions, { type: 'add_reviewer', value: '' }];
  }

  function removeAction(i: number) {
    actions = actions.filter((_, idx) => idx !== i);
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (!name.trim() || submitting) return;
    submitting = true;
    try {
      const resp = await apiPost<{ id: string }>('/api/herald', {
        name,
        must_match_all: mustMatchAll,
        conditions,
        actions,
      });
      goto(`/actions/${resp.id}`);
    } catch (err: unknown) {
      alert(err instanceof Error ? err.message : 'Failed to create rule');
    } finally {
      submitting = false;
    }
  }
</script>

<PageShell title="New Herald Rule" icon="fa-plus">
  {#snippet breadcrumbs()}
    <Breadcrumbs {crumbs} />
  {/snippet}

  <Box border>
    <form onsubmit={handleSubmit}>
      <div class="form-body">
        <div class="form-group">
          <label class="form-label" for="rule-name">Rule Name</label>
          <input id="rule-name" type="text" bind:value={name} required placeholder="e.g. Auto-add reviewers for docs changes" class="form-input" />
        </div>

        <div class="form-group">
          <label class="form-label checkbox-label">
            <input type="checkbox" bind:checked={mustMatchAll} />
            Must match all conditions
          </label>
        </div>

        <div class="form-group">
          <span class="form-label">Conditions</span>
          {#each conditions as cond, i}
            <div class="row-group">
              <select bind:value={cond.type} class="form-input form-select">
                {#each conditionTypes as ct}
                  <option value={ct.value}>{ct.label}</option>
                {/each}
              </select>
              <input type="text" bind:value={cond.value} placeholder="Value" class="form-input row-input" />
              {#if conditions.length > 1}
                <button type="button" class="btn-icon" title="Remove" onclick={() => removeCondition(i)}>
                  <i class="fa fa-times"></i>
                </button>
              {/if}
            </div>
          {/each}
          <button type="button" class="btn-secondary" onclick={addCondition}>
            <i class="fa fa-plus mrs"></i>Add Condition
          </button>
        </div>

        <div class="form-group">
          <span class="form-label">Actions</span>
          {#each actions as act, i}
            <div class="row-group">
              <select bind:value={act.type} class="form-input form-select">
                {#each actionTypes as at}
                  <option value={at.value}>{at.label}</option>
                {/each}
              </select>
              <input type="text" bind:value={act.value} placeholder="Value" class="form-input row-input" />
              {#if actions.length > 1}
                <button type="button" class="btn-icon" title="Remove" onclick={() => removeAction(i)}>
                  <i class="fa fa-times"></i>
                </button>
              {/if}
            </div>
          {/each}
          <button type="button" class="btn-secondary" onclick={addAction}>
            <i class="fa fa-plus mrs"></i>Add Action
          </button>
        </div>

        <Button type="submit" color="green" icon="fa-save" disabled={submitting}>
          Create Rule
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
    margin-bottom: 16px;
  }

  .form-group:nth-child(1),
  .form-group:nth-child(2) {
    margin-bottom: 12px;
  }

  .form-label {
    display: block;
    font-weight: bold;
    margin-bottom: 8px;
    font-size: 13px;
    color: var(--text);
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
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

  input.form-input:not(.row-input) {
    width: 100%;
    max-width: 460px;
  }

  .row-group {
    display: flex;
    gap: 8px;
    margin-bottom: 6px;
    align-items: center;
  }

  .row-input {
    flex: 1;
  }

  .btn-icon {
    padding: 4px 8px;
    font-size: 12px;
    border: 1px solid var(--border);
    border-radius: 3px;
    background: var(--bg-card);
    color: var(--text-muted);
    cursor: pointer;
  }

  .btn-icon:hover {
    background: var(--bg-hover);
  }

  .btn-secondary {
    margin-top: 4px;
    font-size: 12px;
    padding: 4px 10px;
    border: 1px solid var(--border);
    border-radius: 3px;
    background: var(--bg-card);
    color: var(--text);
    cursor: pointer;
  }

  .btn-secondary:hover {
    background: var(--bg-hover);
  }
</style>
