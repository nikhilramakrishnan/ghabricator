<script lang="ts">
  import { Breadcrumbs, PageShell } from '$lib/components/layout';
  import { Button } from '$lib/components/phui';
  import { apiPost } from '$lib/api';
  import { goto } from '$app/navigation';
  import type { HeraldRule, HeraldCondition, HeraldAction } from '$lib/types';

  const crumbs = [
    { name: 'Home', href: '/' },
    { name: 'Herald', href: '/herald' },
    { name: 'New Rule' }
  ];

  const conditionTypes = [
    { value: '', label: '-- select --' },
    { value: 'file_path', label: 'File path matches' },
    { value: 'author', label: 'Author is' },
    { value: 'title', label: 'Title contains' },
    { value: 'label', label: 'Label is' },
    { value: 'base_branch', label: 'Base branch is' },
  ];

  const actionTypes = [
    { value: '', label: '-- select --' },
    { value: 'add_reviewer', label: 'Add reviewer' },
    { value: 'add_label', label: 'Add label' },
    { value: 'post_comment', label: 'Post comment' },
  ];

  let name = $state('');
  let matchMode = $state('all');
  let conditions = $state<HeraldCondition[]>([{ type: '', value: '' }]);
  let actions = $state<HeraldAction[]>([{ type: '', value: '' }]);
  let submitting = $state(false);

  function addCondition() {
    conditions = [...conditions, { type: '', value: '' }];
  }

  function removeCondition(idx: number) {
    conditions = conditions.filter((_, i) => i !== idx);
  }

  function addAction() {
    actions = [...actions, { type: '', value: '' }];
  }

  function removeAction(idx: number) {
    actions = actions.filter((_, i) => i !== idx);
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (!name.trim() || submitting) return;
    submitting = true;
    try {
      const rule = await apiPost<HeraldRule>('/api/herald', {
        name,
        must_match_all: matchMode === 'all',
        conditions: conditions.filter((c) => c.type && c.value),
        actions: actions.filter((a) => a.type && a.value),
        disabled: false
      });
      goto(`/herald/${rule.id}`);
    } catch (err: unknown) {
      alert(err instanceof Error ? err.message : 'Failed to save rule');
    } finally {
      submitting = false;
    }
  }
</script>

<PageShell title="Create Rule" icon="fa-bullhorn">
  {#snippet breadcrumbs()}
    <Breadcrumbs {crumbs} />
  {/snippet}

  <div class="phui-box phui-box-border phui-object-box herald-rule-form">
    <form onsubmit={handleSubmit}>
      <div class="phui-form-view" style="padding:16px">
        <!-- Name -->
        <div style="margin-bottom:12px">
          <label class="herald-label" style="display:block;font-weight:bold;margin-bottom:4px;font-size:13px">Rule Name</label>
          <input type="text" bind:value={name} required class="aphront-form-input" style="width:100%;max-width:460px;padding:6px 8px;border:1px solid #c7ccd9;border-radius:3px" />
        </div>

        <!-- Match mode -->
        <div style="margin-bottom:12px">
          <label class="herald-label" style="display:block;font-weight:bold;margin-bottom:4px;font-size:13px">Match Mode</label>
          <select bind:value={matchMode} class="aphront-form-input" style="padding:4px 8px;border:1px solid #c7ccd9;border-radius:3px">
            <option value="all">All conditions must match</option>
            <option value="any">Any condition must match</option>
          </select>
        </div>

        <!-- Conditions -->
        <div style="margin-bottom:12px">
          <label class="herald-label" style="display:block;font-weight:bold;margin-bottom:4px;font-size:13px">Conditions</label>
          <table class="herald-condition-table">
            {#each conditions as cond, i}
              <tr>
                <td style="padding:4px">
                  <select bind:value={cond.type} class="aphront-form-input" style="width:160px;padding:2px 4px;border:1px solid #c7ccd9;border-radius:3px">
                    {#each conditionTypes as ct}
                      <option value={ct.value}>{ct.label}</option>
                    {/each}
                  </select>
                </td>
                <td style="padding:4px">
                  <input type="text" bind:value={cond.value} placeholder="value" class="aphront-form-input" style="width:95%;max-width:460px;padding:2px 4px;border:1px solid #c7ccd9;border-radius:3px" />
                </td>
                <td style="padding:4px">
                  <button type="button" onclick={() => removeCondition(i)} style="border:none;background:none;cursor:pointer;color:#92969d;font-size:14px">
                    <span class="phui-icon-view phui-font-fa fa-times"></span>
                  </button>
                </td>
              </tr>
            {/each}
          </table>
          <button type="button" onclick={addCondition} class="mood-btn mood-btn-default" style="margin-top:4px;padding:2px 8px;font-size:12px">
            <span class="phui-icon-view phui-font-fa fa-plus mrs"></span>Add Condition
          </button>
        </div>

        <!-- Actions -->
        <div style="margin-bottom:16px">
          <label class="herald-label" style="display:block;font-weight:bold;margin-bottom:4px;font-size:13px">Actions</label>
          <table class="herald-action-table">
            {#each actions as act, i}
              <tr>
                <td style="padding:4px">
                  <select bind:value={act.type} class="aphront-form-input" style="width:160px;padding:2px 4px;border:1px solid #c7ccd9;border-radius:3px">
                    {#each actionTypes as at}
                      <option value={at.value}>{at.label}</option>
                    {/each}
                  </select>
                </td>
                <td style="padding:4px">
                  <input type="text" bind:value={act.value} placeholder="value" class="aphront-form-input" style="width:95%;max-width:460px;padding:2px 4px;border:1px solid #c7ccd9;border-radius:3px" />
                </td>
                <td style="padding:4px">
                  <button type="button" onclick={() => removeAction(i)} style="border:none;background:none;cursor:pointer;color:#92969d;font-size:14px">
                    <span class="phui-icon-view phui-font-fa fa-times"></span>
                  </button>
                </td>
              </tr>
            {/each}
          </table>
          <button type="button" onclick={addAction} class="mood-btn mood-btn-default" style="margin-top:4px;padding:2px 8px;font-size:12px">
            <span class="phui-icon-view phui-font-fa fa-plus mrs"></span>Add Action
          </button>
        </div>

        <Button type="submit" color="green" icon="fa-check" disabled={submitting}>
          Save Rule
        </Button>
      </div>
    </form>
  </div>
</PageShell>
