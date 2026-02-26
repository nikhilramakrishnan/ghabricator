<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    title,
    children
  }: {
    title: string;
    children: Snippet;
  } = $props();
</script>

<div class="mood-curtain-box">
  <div class="mood-curtain-title">{title}</div>
  {@render children()}
</div>
