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
      <span class="phui-icon-view phui-font-fa fa-spinner fa-spin mrs"></span>
      Loading...
    {:else}
      <span class="phui-icon-view phui-font-fa fa-ellipsis-h mrs"></span>
      Show {end - start + 1} more lines
    {/if}
  </td>
</tr>

<style>
  tr.show-more td {
    text-align: center;
    padding: 6px;
    background: rgba(55, 55, 55, 0.04);
    color: #6b748c;
    font-size: 12px;
    cursor: pointer;
  }
  tr.show-more td:hover {
    background: rgba(55, 55, 55, 0.08);
  }

  th.num {
    width: 4em;
  }

  /* Dark mode */
  :global(.phui-theme-dark) tr.show-more td {
    background: rgba(255, 255, 255, 0.05);
    color: rgba(255, 255, 255, 0.3);
  }
  :global(.phui-theme-dark) tr.show-more td:hover {
    background: rgba(255, 255, 255, 0.08);
  }
</style>
