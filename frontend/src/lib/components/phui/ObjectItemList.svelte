<script lang="ts">
  import type { Snippet } from 'svelte';

  let { children }: { children: Snippet } = $props();
</script>

<ul class="phui-oi-list-view">
  {@render children()}
</ul>
