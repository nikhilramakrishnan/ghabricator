<script lang="ts">
  type ActionItem = {
    label: string;
    href?: string;
    icon?: string;
    onclick?: () => void;
  };

  let { actions }: { actions: ActionItem[] } = $props();
</script>

<ul class="phabricator-action-list-view">
  {#each actions as action}
    <li class="phabricator-action-view">
      {#if action.href}
        <a href={action.href} class="phabricator-action-view-href">
          {#if action.icon}
            <span class="phui-icon-view phui-font-fa {action.icon} phabricator-action-view-icon"></span>
          {/if}
          <span class="phabricator-action-view-item">{action.label}</span>
        </a>
      {:else}
        <button
          type="button"
          class="phabricator-action-view-href"
          style="background:none;border:none;cursor:pointer;width:100%;text-align:left;padding:0;font:inherit;color:inherit;"
          onclick={action.onclick}
        >
          {#if action.icon}
            <span class="phui-icon-view phui-font-fa {action.icon} phabricator-action-view-icon"></span>
          {/if}
          <span class="phabricator-action-view-item">{action.label}</span>
        </button>
      {/if}
    </li>
  {/each}
</ul>
