<script lang="ts">
  import { user } from '$lib/stores/auth';
  import { theme, toggleTheme } from '$lib/stores/theme';
  import { goto } from '$app/navigation';

  let { navActive = '' }: { navActive?: string } = $props();

  let searchQuery = $state('');

  function handleSearch(e: Event) {
    e.preventDefault();
    if (searchQuery.trim()) {
      goto(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
    }
  }
</script>

<nav class="topbar">
  <a class="brand" href="/">
    <i class="fa fa-cog"></i>
    <span>Ghabricator</span>
  </a>

  <form class="search-form" onsubmit={handleSearch}>
    <i class="fa fa-search search-icon"></i>
    <input
      type="text"
      bind:value={searchQuery}
      placeholder="Search..."
      class="search-input"
    />
  </form>

  <div class="spacer"></div>

  {#if $user}
    <a href="/api/auth/logout" title="Log out ({$user.login})" class="user-menu">
      {#if $user.avatarURL}
        <img src={$user.avatarURL} alt="" />
      {:else}
        <i class="fa fa-user"></i>
      {/if}
    </a>
  {:else}
    <a href="/api/auth/github" class="login-btn">
      <i class="fa fa-github"></i> Sign in
    </a>
  {/if}
  <button class="theme-toggle" title="Toggle theme" onclick={toggleTheme} type="button">
    <i class="fa {$theme === 'dark' ? 'fa-sun-o' : 'fa-moon-o'}" title="{$theme === 'dark' ? 'Switch to light' : 'Switch to dark'}"></i>
  </button>
</nav>

<style>
  .topbar {
    display: flex;
    align-items: center;
    height: 44px;
    padding: 0 16px;
    gap: 12px;
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
    text-decoration: none;
    flex-shrink: 0;
    width: 184px; /* 200px sidebar - 16px padding */
  }
  .brand:hover { text-decoration: none; }
  .brand .fa { font-size: 18px; opacity: 0.8; }

  .search-form {
    position: relative;
    flex: 0 1 320px;
  }
  .search-icon {
    position: absolute;
    left: 10px;
    top: 50%;
    transform: translateY(-50%);
    color: rgba(255,255,255,0.4);
    font-size: 12px;
    pointer-events: none;
  }
  .search-input {
    width: 100%;
    padding: 6px 10px 6px 30px;
    border: 1px solid rgba(255,255,255,0.15);
    border-radius: 4px;
    background: rgba(255,255,255,0.08);
    color: #fff;
    font-size: 13px;
    outline: none;
  }
  .search-input::placeholder {
    color: rgba(255,255,255,0.4);
  }
  .search-input:focus {
    border-color: rgba(255,255,255,0.3);
    background: rgba(255,255,255,0.12);
  }

  .spacer { flex: 1; }

  .user-menu {
    display: flex;
    align-items: center;
    color: rgba(255,255,255,0.8);
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
  }
  .theme-toggle:hover { color: rgba(255,255,255,0.8); }
</style>
