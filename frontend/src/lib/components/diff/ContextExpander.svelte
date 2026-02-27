<script lang="ts">
  import { apiFetch } from '$lib/api';

  let {
    path,
    start,
    end,
    csId,
    onExpand
  }: {
    path: string;
    start: number;
    end: number;
    csId: number;
    onExpand: (html: string) => void;
  } = $props();

  let loading = $state(false);

  async function expand() {
    if (loading) return;
    loading = true;
    try {
      const params = new URLSearchParams({
        path,
        start: String(start),
        end: String(end),
        cs: String(csId)
      });
      const result = await apiFetch<{ html: string }>(`/api/context?${params}`);
      onExpand(result.html);
    } catch {
      // silently fail
    } finally {
      loading = false;
    }
  }
</script>

<tr class="show-more">
  <th class="num"></th>
  <td
    class="show-more-content"
    colspan="5"
    role="button"
    tabindex="0"
    onclick={expand}
    onkeydown={(e) => e.key === 'Enter' && expand()}
  >
    {#if loading}
      <i class="fa fa-spinner fa-spin mrs"></i>
      Loading...
    {:else}
      <i class="fa fa-ellipsis-h mrs"></i>
      Show {end - start + 1} more lines
    {/if}
  </td>
</tr>

<style>
  tr.show-more td {
    text-align: center;
    padding: 6px;
    background: var(--bg-hover);
    color: var(--text-muted);
    font-size: 12px;
    cursor: pointer;
  }
  tr.show-more td:hover {
    background: var(--bg-subtle);
  }

  th.num {
    width: 4em;
  }
</style>
