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
</script>

<li class="item" class:has-bar={!!barColor}>
  {#if barColor}
    <div class="bar bar-{barColor}"></div>
  {/if}
  <div class="item-frame">
    {#if imageUrl}
      <div class="avatar" style="background-image:url({imageUrl})"></div>
    {:else if icon}
      <div class="icon-avatar">
        <i class="fa {icon}"></i>
      </div>
    {/if}
    <div class="content">
      <div class="title-row">
        <a {href} class="title-link">{title}</a>
        {#if tags}
          {' '}{@render tags()}
        {/if}
      </div>
      {#if attributes}
        <div class="attrs">
          {@render attributes()}
        </div>
      {/if}
    </div>
    {#if handleIcons}
      <div class="handles">
        {@render handleIcons()}
      </div>
    {/if}
  </div>
</li>

<style>
  .item {
    position: relative;
    border-bottom: 1px solid var(--border-subtle);
    padding: 7px 12px;
  }
  .item:last-child {
    border-bottom: none;
  }
  .has-bar {
    padding-left: 15px;
  }
  .bar {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 3px;
  }
  .bar-blue { background: var(--blue); }
  .bar-green { background: var(--green); }
  .bar-red { background: var(--red); }
  .bar-grey { background: #92969d; }
  .bar-violet { background: var(--violet); }
  .bar-orange { background: var(--orange); }
  .bar-yellow { background: var(--yellow); }
  .item-frame {
    display: flex;
    align-items: flex-start;
    gap: 10px;
  }
  .avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background-size: cover;
    background-position: center;
    flex-shrink: 0;
  }
  .icon-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: var(--bg-subtle);
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    color: var(--text-muted);
    font-size: 16px;
  }
  .content {
    flex: 1;
    min-width: 0;
  }
  .title-row {
    display: flex;
    align-items: center;
    gap: 6px;
    flex-wrap: wrap;
  }
  .title-link {
    color: var(--text-link);
    font-weight: 600;
    font-size: 14px;
    text-decoration: none;
  }
  .title-link:hover {
    text-decoration: underline;
  }
  .attrs {
    margin-top: 4px;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }
  .handles {
    flex-shrink: 0;
    display: flex;
    gap: 4px;
  }
</style>
