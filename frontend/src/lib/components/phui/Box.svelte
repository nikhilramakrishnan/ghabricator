<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    border = false,
    children
  }: {
    border?: boolean;
    children: Snippet;
  } = $props();
</script>

<div class="phui-box phui-object-box" class:phui-box-border={border}>
  {@render children()}
</div>
