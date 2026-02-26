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
  <div class="phui-two-column-view">
    <div class="phui-two-column-container">
      <div class="phui-two-column-content">
        <div class="phui-two-column-row grouped">
          <div class="phui-main-column">
            <div class="phui-box phui-box-border phui-object-box" style="text-align:center;padding:40px 20px">
              <div style="margin-bottom:24px">
                <span class="phui-icon-view phui-font-fa fa-code-fork" style="font-size:48px;color:#6b748c"></span>
              </div>
              <h2 style="margin-bottom:12px;font-size:24px">Welcome to Ghabricator</h2>
              <p style="color:#6b748c;font-size:15px;margin-bottom:24px;max-width:480px;margin-left:auto;margin-right:auto">
                A modern code review experience powered by GitHub. Review pull requests with
                Phabricator's legendary diff viewer, inline comments, and Herald automation.
              </p>
              <a href="/api/auth/github" class="phui-button-view button-green" style="font-size:16px;padding:10px 28px;display:inline-flex;align-items:center;gap:8px">
                <span class="phui-icon-view phui-font-fa fa-github"></span> Sign in with GitHub
              </a>
              <div style="margin-top:32px;color:#92969d;font-size:13px">
                <div style="display:flex;justify-content:center;gap:32px;flex-wrap:wrap">
                  <span><span class="phui-icon-view phui-font-fa fa-search mrs"></span> Side-by-side diffs</span>
                  <span><span class="phui-icon-view phui-font-fa fa-commenting-o mrs"></span> Inline comments</span>
                  <span><span class="phui-icon-view phui-font-fa fa-bullhorn mrs"></span> Herald rules</span>
                  <span><span class="phui-icon-view phui-font-fa fa-clipboard mrs"></span> Paste bin</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
{/if}
