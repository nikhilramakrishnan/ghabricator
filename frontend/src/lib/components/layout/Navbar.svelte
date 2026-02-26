<script lang="ts">
  import { user } from '$lib/stores/auth';
  import { theme, toggleTheme } from '$lib/stores/theme';

  let { navActive = '' }: { navActive?: string } = $props();
</script>

<div class="phabricator-main-menu phabricator-main-menu-background" data-sigil="phabricator-main-menu">
  <a class="phabricator-main-menu-brand" href="/">
    <span class="brand-icon phui-icon-view phui-font-fa fa-cog"></span>
    <span>Ghabricator</span>
  </a>
  {#if $user}
    <a href="/dashboard" class="phabricator-main-menu-href" class:nav-active={navActive === 'revisions'}>
      <span class="phui-icon-view phui-font-fa fa-code-fork"></span> Revisions
    </a>
    <a href="/repos" class="phabricator-main-menu-href" class:nav-active={navActive === 'repos'}>
      <span class="phui-icon-view phui-font-fa fa-database"></span> Repos
    </a>
    <a href="/paste" class="phabricator-main-menu-href" class:nav-active={navActive === 'paste'}>
      <span class="phui-icon-view phui-font-fa fa-clipboard"></span> Paste
    </a>
    <a href="/herald" class="phabricator-main-menu-href" class:nav-active={navActive === 'herald'}>
      <span class="phui-icon-view phui-font-fa fa-shield"></span> Herald
    </a>
    <a href="/search" class="phabricator-main-menu-href" class:nav-active={navActive === 'search'}>
      <span class="phui-icon-view phui-font-fa fa-search"></span> Search
    </a>
    <div class="nav-spacer"></div>
    <a href="/api/auth/logout" title="Log out ({$user.login})" class="phabricator-core-user-menu">
      {#if $user.avatarURL}
        <img src={$user.avatarURL} alt="" />
      {:else}
        <span class="phui-icon-view phui-font-fa fa-user"></span>
      {/if}
      <span>{$user.login}</span>
    </a>
  {:else}
    <div class="nav-spacer"></div>
    <a href="/api/auth/github" class="phabricator-core-login-button">
      <span class="phui-icon-view phui-font-fa fa-github"></span> Sign in
    </a>
  {/if}
  <div class="phabricator-main-menu-alerts">
    <button class="alert-notifications" title="Toggle theme" onclick={toggleTheme} type="button">
      <span class="phui-icon-view phui-font-fa {$theme === 'dark' ? 'fa-sun-o' : 'fa-moon-o'}"></span>
    </button>
  </div>
</div>
