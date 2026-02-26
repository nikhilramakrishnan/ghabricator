<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import { Button, Box } from '$lib/components/phui';
  import { apiPost } from '$lib/api';
  import { goto } from '$app/navigation';
  import type { HeraldCondition, HeraldAction } from '$lib/types';

  const crumbs = [
    { name: 'Home', href: '/' },
    { name: 'Herald', href: '/herald' },
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
      goto(`/herald/${resp.id}`);
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
      <div class="phui-form-view" style="padding:16px">
        <div class="aphront-form-control" style="margin-bottom:12px">
          <label class="aphront-form-label" style="display:block;font-weight:bold;margin-bottom:4px;font-size:13px">Rule Name</label>
          <input type="text" bind:value={name} required placeholder="e.g. Auto-add reviewers for docs changes" class="aphront-form-input" style="width:100%;max-width:460px;padding:6px 8px;border:1px solid #c7ccd9;border-radius:3px" />
        </div>

        <div class="aphront-form-control" style="margin-bottom:16px">
          <label class="aphront-form-label" style="display:flex;align-items:center;gap:8px;font-weight:bold;font-size:13px;cursor:pointer">
            <input type="checkbox" bind:checked={mustMatchAll} />
            Must match all conditions
          </label>
        </div>

        <div class="aphront-form-control" style="margin-bottom:16px">
          <label class="aphront-form-label" style="display:block;font-weight:bold;margin-bottom:8px;font-size:13px">Conditions</label>
          {#each conditions as cond, i}
            <div style="display:flex;gap:8px;margin-bottom:6px;align-items:center">
              <select bind:value={cond.type} class="aphront-form-input" style="padding:6px 8px;border:1px solid #c7ccd9;border-radius:3px">
                {#each conditionTypes as ct}
                  <option value={ct.value}>{ct.label}</option>
                {/each}
              </select>
              <input type="text" bind:value={cond.value} placeholder="Value" class="aphront-form-input" style="flex:1;padding:6px 8px;border:1px solid #c7ccd9;border-radius:3px" />
              {#if conditions.length > 1}
                <button type="button" class="mood-btn mood-btn-default" style="padding:4px 8px;font-size:12px" onclick={() => removeCondition(i)}>
                  <span class="phui-icon-view phui-font-fa fa-times"></span>
                </button>
              {/if}
            </div>
          {/each}
          <button type="button" class="mood-btn mood-btn-default" style="margin-top:4px;font-size:12px" onclick={addCondition}>
            <span class="phui-icon-view phui-font-fa fa-plus mrs"></span>Add Condition
          </button>
        </div>

        <div class="aphront-form-control" style="margin-bottom:16px">
          <label class="aphront-form-label" style="display:block;font-weight:bold;margin-bottom:8px;font-size:13px">Actions</label>
          {#each actions as act, i}
            <div style="display:flex;gap:8px;margin-bottom:6px;align-items:center">
              <select bind:value={act.type} class="aphront-form-input" style="padding:6px 8px;border:1px solid #c7ccd9;border-radius:3px">
                {#each actionTypes as at}
                  <option value={at.value}>{at.label}</option>
                {/each}
              </select>
              <input type="text" bind:value={act.value} placeholder="Value" class="aphront-form-input" style="flex:1;padding:6px 8px;border:1px solid #c7ccd9;border-radius:3px" />
              {#if actions.length > 1}
                <button type="button" class="mood-btn mood-btn-default" style="padding:4px 8px;font-size:12px" onclick={() => removeAction(i)}>
                  <span class="phui-icon-view phui-font-fa fa-times"></span>
                </button>
              {/if}
            </div>
          {/each}
          <button type="button" class="mood-btn mood-btn-default" style="margin-top:4px;font-size:12px" onclick={addAction}>
            <span class="phui-icon-view phui-font-fa fa-plus mrs"></span>Add Action
          </button>
        </div>

        <Button type="submit" color="green" icon="fa-save" disabled={submitting}>
          Create Rule
        </Button>
      </div>
    </form>
  </Box>
</PageShell>
