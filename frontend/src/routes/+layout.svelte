<script lang="ts">
  import '../app.css';
  import Navbar from '$lib/components/layout/Navbar.svelte';
  import { checkAuth, authLoading } from '$lib/stores/auth';
  import { theme } from '$lib/stores/theme';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';

  let { children } = $props();

  onMount(() => {
    checkAuth();

    // Manage body classes reactively
    const unsubTheme = theme.subscribe((t) => {
      document.body.classList.add('device-desktop', 'phui-theme-blindigo');
      document.body.classList.toggle('phui-theme-dark', t === 'dark');
    });

    return unsubTheme;
  });

  function getNavActive(path: string): string {
    if (path.startsWith('/pr/') || path === '/dashboard') return 'revisions';
    if (path.startsWith('/repo')) return 'repos';
    if (path.startsWith('/paste')) return 'paste';
    if (path.startsWith('/herald')) return 'herald';
    if (path.startsWith('/search')) return 'search';
    return '';
  }

  let navActive = $derived(getNavActive($page.url.pathname));
</script>

<Navbar {navActive} />

{#if $authLoading}
  <!-- loading -->
{:else}
  {@render children()}
{/if}
