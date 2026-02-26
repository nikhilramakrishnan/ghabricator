<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    title,
    href = '#',
    barColor = '',
    imageUrl = '',
    icon = '',
    tags,
    attributes,
    handleIcons
  }: {
    title: string;
    href?: string;
    barColor?: string;
    imageUrl?: string;
    icon?: string;
    tags?: Snippet;
    attributes?: Snippet;
    handleIcons?: Snippet;
  } = $props();

  let imageClass = $derived(imageUrl ? 'phui-oi-with-image' : 'phui-oi-with-image-icon');
  let barClass = $derived(barColor ? `phui-oi-bar-color-${barColor}` : '');
</script>

<li class="phui-oi {barClass} {imageClass}">
  <div class="phui-oi-frame">
    <div class="phui-oi-frame-content">
      {#if imageUrl}
        <div class="phui-oi-image" style="background-image:url({imageUrl})"></div>
      {:else if icon}
        <div class="phui-oi-image-icon">
          <span class="phui-icon-view phui-font-fa {icon}"></span>
        </div>
      {/if}
      <div class="phui-oi-content-box">
        <div class="phui-oi-table">
          <div class="phui-oi-table-row">
            <div class="phui-oi-col1">
              <div class="phui-oi-name">
                <a {href} class="phui-oi-link">{title}</a>
                {#if tags}
                  {' '}{@render tags()}
                {/if}
              </div>
              {#if attributes}
                <div class="phui-oi-content">
                  <ul class="phui-oi-attributes">
                    {@render attributes()}
                  </ul>
                </div>
              {/if}
            </div>
          </div>
        </div>
      </div>
    </div>
    {#if handleIcons}
      <div class="phui-oi-handle-icons">
        {@render handleIcons()}
      </div>
    {/if}
  </div>
</li>
