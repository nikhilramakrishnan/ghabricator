<script lang="ts">
  import { user } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

  $effect(() => {
    if ($user) {
      goto('/dashboard');
    }
  });
</script>

{#if !$user}
  <div class="landing">
    <div class="hero-card">
      <div class="hero-icon">
        <i class="fa fa-code-fork"></i>
      </div>
      <h2 class="hero-title">Welcome to Ghabricator</h2>
      <p class="hero-desc">
        A modern code review experience powered by GitHub. Review pull requests with
        Phabricator's legendary diff viewer, inline comments, and Herald automation.
      </p>
      <a href="/api/auth/github" class="signin-btn">
        <i class="fa fa-github"></i> Sign in with GitHub
      </a>
      <div class="features">
        <span><i class="fa fa-search mrs"></i> Side-by-side diffs</span>
        <span><i class="fa fa-commenting-o mrs"></i> Inline comments</span>
        <span><i class="fa fa-bullhorn mrs"></i> Herald rules</span>
        <span><i class="fa fa-clipboard mrs"></i> Paste bin</span>
      </div>
    </div>
  </div>
{/if}

<style>
  .landing {
    display: flex;
    justify-content: center;
    padding: 60px 20px;
  }

  .hero-card {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 4px;
    text-align: center;
    padding: 40px 20px;
    max-width: 600px;
    width: 100%;
  }

  .hero-icon {
    margin-bottom: 24px;
    font-size: 48px;
    color: var(--text-muted);
  }

  .hero-title {
    margin: 0 0 12px;
    font-size: 24px;
    color: var(--text);
  }

  .hero-desc {
    color: var(--text-muted);
    font-size: 15px;
    margin: 0 auto 24px;
    max-width: 480px;
    line-height: 1.5;
  }

  .signin-btn {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: var(--green);
    color: var(--text-on-dark);
    font-size: 16px;
    padding: 10px 28px;
    border-radius: 4px;
    text-decoration: none;
    font-weight: 600;
  }

  .signin-btn:hover {
    background: var(--green-hover);
    text-decoration: none;
  }

  .features {
    margin-top: 32px;
    color: var(--text-muted);
    font-size: 13px;
    display: flex;
    justify-content: center;
    gap: 32px;
    flex-wrap: wrap;
  }
</style>
