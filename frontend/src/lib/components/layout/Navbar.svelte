<script lang="ts">
  import { user } from '$lib/stores/auth';
  import { theme, toggleTheme } from '$lib/stores/theme';

  let { navActive = '' }: { navActive?: string } = $props();
</script>

<nav class="navbar">
  <a class="brand" href="/">
    <i class="fa fa-cog"></i>
    <span>Ghabricator</span>
  </a>
  {#if $user}
    <a href="/dashboard" class="nav-link" class:active={navActive === 'revisions'}>
      <i class="fa fa-code-fork"></i> Revisions
    </a>
    <a href="/repos" class="nav-link" class:active={navActive === 'repos'}>
      <i class="fa fa-database"></i> Repos
    </a>
    <a href="/paste" class="nav-link" class:active={navActive === 'paste'}>
      <i class="fa fa-clipboard"></i> Paste
    </a>
    <a href="/herald" class="nav-link" class:active={navActive === 'herald'}>
      <i class="fa fa-bullhorn"></i> Herald
    </a>
    <a href="/search" class="nav-link" class:active={navActive === 'search'}>
      <i class="fa fa-search"></i> Search
    </a>
    <div class="spacer"></div>
    <a href="/api/auth/logout" title="Log out ({$user.login})" class="user-menu">
      {#if $user.avatarURL}
        <img src={$user.avatarURL} alt="" />
      {:else}
        <i class="fa fa-user"></i>
      {/if}
      <span>{$user.login}</span>
    </a>
  {:else}
    <div class="spacer"></div>
    <a href="/api/auth/github" class="login-btn">
      <i class="fa fa-github"></i> Sign in
    </a>
  {/if}
  <button class="theme-toggle" title="Toggle theme" onclick={toggleTheme} type="button">
    <i class="fa {$theme === 'dark' ? 'fa-sun-o' : 'fa-moon-o'}"></i>
  </button>
</nav>

<style>
  .navbar {
    display: flex;
    align-items: center;
    height: 44px;
    padding: 0 16px;
    gap: 4px;
    background: var(--nav-bg);
    color: var(--text-on-dark);
  }
  .brand {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #fff;
    font-weight: 700;
    font-size: 14px;
    margin-right: 16px;
    text-decoration: none;
  }
  .brand:hover { text-decoration: none; }
  .brand .fa { font-size: 18px; opacity: 0.8; }
  .nav-link {
    color: rgba(255,255,255,0.7);
    font-size: 13px;
    padding: 8px 12px;
    text-decoration: none;
    border-radius: 4px;
    display: flex;
    align-items: center;
    gap: 6px;
  }
  .nav-link:hover { color: #fff; background: rgba(255,255,255,0.08); text-decoration: none; }
  .nav-link.active { color: #fff; background: rgba(255,255,255,0.12); }
  .spacer { flex: 1; }
  .user-menu {
    display: flex;
    align-items: center;
    gap: 8px;
    color: rgba(255,255,255,0.8);
    font-size: 13px;
    text-decoration: none;
  }
  .user-menu:hover { text-decoration: none; }
  .user-menu img { width: 28px; height: 28px; border-radius: 3px; }
  .login-btn {
    color: rgba(255,255,255,0.8);
    font-size: 13px;
    text-decoration: none;
    display: flex;
    align-items: center;
    gap: 6px;
  }
  .login-btn:hover { text-decoration: none; color: #fff; }
  .theme-toggle {
    color: rgba(255,255,255,0.5);
    font-size: 16px;
    padding: 4px 8px;
    cursor: pointer;
    background: none;
    border: none;
    margin-left: 8px;
  }
  .theme-toggle:hover { color: rgba(255,255,255,0.8); }
</style>
