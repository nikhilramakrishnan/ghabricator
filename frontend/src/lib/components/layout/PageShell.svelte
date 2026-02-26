<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    title = '',
    icon = '',
    headerRight,
    breadcrumbs,
    children
  }: {
    title?: string;
    icon?: string;
    headerRight?: Snippet;
    breadcrumbs?: Snippet;
    children: Snippet;
  } = $props();
</script>

<div class="phui-two-column-view">
  <div class="phui-two-column-container">
    {#if breadcrumbs}
      {@render breadcrumbs()}
    {/if}
    {#if title}
      <div class="phui-two-column-header">
        <div class="phui-header-view">
          <div class="phui-header-shell" style="display:flex;align-items:center;justify-content:space-between;">
            <h1 class="phui-header-header">
              {#if icon}
                <span class="phui-header-icon phui-icon-view phui-font-fa {icon}"></span>
              {/if}
              {title}
            </h1>
            {#if headerRight}
              {@render headerRight()}
            {/if}
          </div>
        </div>
      </div>
    {/if}
    <div class="phui-two-column-content">
      {@render children()}
    </div>
  </div>
</div>
