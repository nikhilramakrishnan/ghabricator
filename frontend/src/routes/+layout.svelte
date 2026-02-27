<script lang="ts">
  import '../app.css';
  import Navbar from '$lib/components/layout/Navbar.svelte';
  import Sidebar from '$lib/components/layout/Sidebar.svelte';
  import { checkAuth, authLoading, user } from '$lib/stores/auth';
  import { theme } from '$lib/stores/theme';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';

  let { children } = $props();

  onMount(() => {
    checkAuth();
    const unsubTheme = theme.subscribe((t) => {
      document.body.classList.toggle('dark', t === 'dark');
    });
    return unsubTheme;
  });

  function getNavActive(path: string): string {
    if (path.startsWith('/pr/') || path === '/dashboard') return 'revisions';
    if (path.startsWith('/repo')) return 'repos';
    if (path.startsWith('/paste')) return 'paste';
    if (path.startsWith('/actions')) return 'actions';
    if (path.startsWith('/search')) return 'search';
    return '';
  }

  let navActive = $derived(getNavActive($page.url.pathname));
</script>

<div class="app-layout">
  <Navbar {navActive} />
  <div class="app-body">
    {#if $user}
      <Sidebar {navActive} />
    {/if}
    <main class="app-main">
      {#if $authLoading}
        <!-- loading -->
      {:else}
        {@render children()}
      {/if}
    </main>
  </div>
</div>

<style>
  .app-layout {
    display: flex;
    flex-direction: column;
    height: 100vh;
  }

  .app-body {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  .app-main {
    flex: 1;
    overflow-y: auto;
    background: var(--bg);
  }
</style>
