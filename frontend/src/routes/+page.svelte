<script lang="ts">
  import { user } from '$lib/stores/auth';
  import { PageShell } from '$lib/components/layout';
  import {
    Box, HeaderView, ObjectItemList, ObjectItem, Attribute,
    Tag, InfoView
  } from '$lib/components/phui';
  import type { APIPRSummary, APIWorkflowRun } from '$lib/types';
  import { S } from '$lib/strings';

  let { data } = $props();
  let authored: APIPRSummary[] = $derived(data.authored ?? []);
  let reviewRequested: APIPRSummary[] = $derived(data.reviewRequested ?? []);
  let recentRuns: APIWorkflowRun[] = $derived(data.recentRuns ?? []);

  function runBarColor(run: APIWorkflowRun): string {
    if (run.conclusion === 'success') return 'green';
    if (run.conclusion === 'failure') return 'red';
    if (run.status === 'in_progress' || run.status === 'queued' || run.status === 'pending') return 'yellow';
    return 'grey';
  }

  function runStatusShade(run: APIWorkflowRun): string {
    return runBarColor(run);
  }

  function runStatusLabel(run: APIWorkflowRun): string {
    if (run.conclusion) return run.conclusion;
    return run.status;
  }
</script>

{#if !$user}
  <div class="landing">
    <!-- Hero -->
    <section class="hero">
      <h1 class="hero-title">Code review that doesn't suck.</h1>
      <p class="hero-sub">Phabricator's legendary review workflow, powered by your GitHub repos. No migration, no lock-in.</p>
      <a href="/api/auth/github" class="cta">
        <i class="fa fa-github"></i> Sign in with GitHub
      </a>
    </section>

    <!-- Diff preview -->
    <section class="preview">
      <div class="preview-window">
        <div class="preview-bar">
          <span class="dot red"></span><span class="dot yellow"></span><span class="dot green"></span>
          <span class="preview-path">src/auth/session.go</span>
        </div>
        <table class="diff-mock">
          <tbody>
            <tr class="ctx"><td class="ln">14</td><td class="ln">14</td><td class="code">func (s *Store) Get(id string) *Session {'{'}</td></tr>
            <tr class="del"><td class="ln">15</td><td class="ln"></td><td class="code">-   sess, ok := s.cache[id]</td></tr>
            <tr class="add"><td class="ln"></td><td class="ln">15</td><td class="code">+   sess, ok := s.sessions[id]</td></tr>
            <tr class="ctx"><td class="ln">16</td><td class="ln">16</td><td class="code">    if !ok {'{'}</td></tr>
            <tr class="ctx"><td class="ln">17</td><td class="ln">17</td><td class="code">        return nil</td></tr>
            <tr class="ctx"><td class="ln">18</td><td class="ln">18</td><td class="code">    {'}'}</td></tr>
            <tr class="add"><td class="ln"></td><td class="ln">19</td><td class="code">+   if time.Since(sess.CreatedAt) > sessionTTL {'{'}</td></tr>
            <tr class="add"><td class="ln"></td><td class="ln">20</td><td class="code">+       return nil</td></tr>
            <tr class="add"><td class="ln"></td><td class="ln">21</td><td class="code">+   {'}'}</td></tr>
            <tr class="ctx"><td class="ln">19</td><td class="ln">22</td><td class="code">    return sess</td></tr>
            <tr class="ctx"><td class="ln">20</td><td class="ln">23</td><td class="code">{'}'}</td></tr>
          </tbody>
        </table>
        <div class="inline-mock">
          <div class="inline-avatar">NR</div>
          <div class="inline-body">
            <span class="inline-author">nikhilr</span>
            <span class="inline-text">Should we also evict expired sessions on write? This only catches reads.</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Features -->
    <section class="features">
      <div class="feature">
        <div class="feature-icon"><i class="fa fa-columns"></i></div>
        <h3>Side-by-side diffs</h3>
        <p>Syntax-highlighted, two-column diff viewer with context expansion. The way diffs should be read.</p>
      </div>
      <div class="feature">
        <div class="feature-icon"><i class="fa fa-commenting-o"></i></div>
        <h3>Inline comments</h3>
        <p>Click any line to comment. Draft comments, threaded replies, edit history. Ship reviews as a batch.</p>
      </div>
      <div class="feature">
        <div class="feature-icon"><i class="fa fa-bullhorn"></i></div>
        <h3>Herald automation</h3>
        <p>Rules that auto-assign reviewers, add labels, or post comments based on file paths, authors, or branches.</p>
      </div>
      <div class="feature">
        <div class="feature-icon"><i class="fa fa-search"></i></div>
        <h3>Search everything</h3>
        <p>Search PRs, issues, code, and repos across all your GitHub organizations. One search bar, instant results.</p>
      </div>
      <div class="feature">
        <div class="feature-icon"><i class="fa fa-code"></i></div>
        <h3>Repository browser</h3>
        <p>Browse files, view blame, syntax-highlighted source — without leaving the review tool.</p>
      </div>
      <div class="feature">
        <div class="feature-icon"><i class="fa fa-github"></i></div>
        <h3>GitHub is the source of truth</h3>
        <p>No database, no migration. Your repos, PRs, and comments live on GitHub. Ghabricator is just a better lens.</p>
      </div>
    </section>

    <!-- Bottom CTA -->
    <section class="bottom-cta">
      <a href="/api/auth/github" class="cta">
        <i class="fa fa-github"></i> Get started
      </a>
      <p class="bottom-note">Connects to your existing GitHub account.</p>
    </section>
  </div>
{:else}
  <PageShell title={S.dashboard.title} icon="fa-home">
    <div class="dashboard-grid">
      <div class="dash-col">
        <Box border>
          <HeaderView title={S.dashboard.needsReview} icon="fa-eye" count={reviewRequested.length} />
          {#if reviewRequested.length === 0}
            <InfoView icon="fa-inbox">{S.dashboard.noReviews}</InfoView>
          {:else}
            <ObjectItemList>
              {#each reviewRequested as pr}
                <ObjectItem
                  title={pr.title}
                  href="/pr/{pr.owner}/{pr.repo}/{pr.number}"
                  imageUrl={pr.author.avatarURL}
                  barColor={pr.draft ? 'grey' : 'blue'}
                >
                  {#snippet tags()}
                    {#if pr.draft}<Tag shade="grey">Draft</Tag>{/if}
                    {#if pr.labels}
                      {#each pr.labels as label}
                        <Tag shade="blue">{label.name}</Tag>
                      {/each}
                    {/if}
                  {/snippet}
                  {#snippet attributes()}
                    <Attribute icon="fa-user">{pr.author.login}</Attribute>
                    <Attribute icon="fa-github">{pr.owner}/{pr.repo}#{pr.number}</Attribute>
                  {/snippet}
                </ObjectItem>
              {/each}
            </ObjectItemList>
          {/if}
        </Box>

        <Box border>
          <HeaderView title={S.dashboard.authored} icon="fa-pencil" count={authored.length} />
          {#if authored.length === 0}
            <InfoView icon="fa-inbox">{S.dashboard.noPRs}</InfoView>
          {:else}
            <ObjectItemList>
              {#each authored as pr}
                <ObjectItem
                  title={pr.title}
                  href="/pr/{pr.owner}/{pr.repo}/{pr.number}"
                  imageUrl={pr.author.avatarURL}
                  barColor={pr.draft ? 'grey' : 'blue'}
                >
                  {#snippet tags()}
                    {#if pr.draft}<Tag shade="grey">Draft</Tag>{/if}
                    {#if pr.labels}
                      {#each pr.labels as label}
                        <Tag shade="blue">{label.name}</Tag>
                      {/each}
                    {/if}
                  {/snippet}
                  {#snippet attributes()}
                    <Attribute icon="fa-user">{pr.author.login}</Attribute>
                    <Attribute icon="fa-github">{pr.owner}/{pr.repo}#{pr.number}</Attribute>
                  {/snippet}
                </ObjectItem>
              {/each}
            </ObjectItemList>
          {/if}
        </Box>
      </div>

      <div class="dash-col dash-col-narrow">
        <Box border>
          <HeaderView title={S.dashboard.recentBuilds} icon="fa-cog" />
          {#if recentRuns.length === 0}
            <InfoView icon="fa-inbox">{S.dashboard.noBuilds}</InfoView>
          {:else}
            <ObjectItemList>
              {#each recentRuns as run}
                <ObjectItem
                  title={run.name}
                  href={run.htmlURL}
                  imageUrl={run.actor.avatarURL}
                  barColor={runBarColor(run)}
                >
                  {#snippet tags()}
                    <Tag shade={runStatusShade(run)}>{runStatusLabel(run)}</Tag>
                  {/snippet}
                  {#snippet attributes()}
                    <Attribute icon="fa-code-fork">{run.branch}</Attribute>
                    <Attribute icon="fa-database">{run.repoOwner}/{run.repoName}</Attribute>
                  {/snippet}
                </ObjectItem>
              {/each}
            </ObjectItemList>
            <a href="/actions" class="view-all">{S.dashboard.viewAllBuilds} <i class="fa fa-arrow-right"></i></a>
          {/if}
        </Box>
      </div>
    </div>
  </PageShell>
{/if}

<style>
  /* ===== Landing page ===== */
  .landing {
    max-width: 960px;
    margin: 0 auto;
    padding: 0 24px 80px;
  }

  /* Hero */
  .hero {
    text-align: center;
    padding: 72px 0 48px;
  }
  .hero-title {
    font-size: 40px;
    font-weight: 800;
    margin: 0 0 16px;
    color: var(--text);
    letter-spacing: -0.5px;
    line-height: 1.15;
  }
  .hero-sub {
    font-size: 17px;
    color: var(--text-muted);
    margin: 0 auto 32px;
    max-width: 540px;
    line-height: 1.6;
  }
  .cta {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    background: var(--green);
    color: var(--text-on-dark);
    font-size: 16px;
    font-weight: 700;
    padding: 12px 32px;
    border-radius: 6px;
    text-decoration: none;
    transition: background 0.15s;
  }
  .cta:hover {
    background: var(--green-hover);
    text-decoration: none;
  }

  /* Diff preview */
  .preview {
    margin: 0 auto 64px;
    max-width: 680px;
  }
  .preview-window {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 4px 24px rgba(0,0,0,0.08);
  }
  .preview-bar {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 10px 14px;
    background: var(--bg-card-header);
    border-bottom: 1px solid var(--border-subtle);
  }
  .dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
  }
  .dot.red { background: #ff5f57; }
  .dot.yellow { background: #febc2e; }
  .dot.green { background: #28c840; }
  .preview-path {
    margin-left: 8px;
    font-size: 12px;
    font-family: var(--font-mono);
    color: var(--text-muted);
  }

  .diff-mock {
    width: 100%;
    border-collapse: collapse;
    font-family: var(--font-mono);
    font-size: 12px;
    line-height: 1.6;
  }
  .diff-mock td {
    padding: 0 10px;
    white-space: pre;
  }
  .diff-mock .ln {
    width: 36px;
    text-align: right;
    color: var(--text-muted);
    opacity: 0.5;
    user-select: none;
    padding: 0 8px;
  }
  .diff-mock .code {
    width: 100%;
  }
  .diff-mock .ctx { background: transparent; }
  .diff-mock .add { background: var(--diff-add-bg); }
  .diff-mock .add .ln { background: var(--diff-add-num-bg); }
  .diff-mock .del { background: var(--diff-del-bg); }
  .diff-mock .del .ln { background: var(--diff-del-num-bg); }

  .inline-mock {
    display: flex;
    gap: 10px;
    padding: 12px 14px;
    border-top: 1px solid var(--border-subtle);
    background: var(--bg-subtle);
  }
  .inline-avatar {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    background: var(--blue);
    color: #fff;
    font-size: 11px;
    font-weight: 700;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }
  .inline-body {
    font-size: 13px;
    line-height: 1.5;
    color: var(--text);
  }
  .inline-author {
    font-weight: 700;
    margin-right: 6px;
  }
  .inline-text {
    color: var(--text-muted);
  }

  /* Features grid */
  .features {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 32px;
    margin-bottom: 64px;
  }
  .feature {
    text-align: center;
  }
  .feature-icon {
    font-size: 28px;
    color: var(--blue);
    margin-bottom: 12px;
  }
  .feature h3 {
    font-size: 15px;
    font-weight: 700;
    margin: 0 0 8px;
    color: var(--text);
  }
  .feature p {
    font-size: 13px;
    color: var(--text-muted);
    line-height: 1.6;
    margin: 0;
  }

  /* Bottom CTA */
  .bottom-cta {
    text-align: center;
    padding: 24px 0;
  }
  .bottom-note {
    margin: 16px 0 0;
    font-size: 13px;
    color: var(--text-muted);
  }

  .dashboard-grid {
    display: flex;
    gap: 16px;
  }

  .dash-col {
    flex: 3;
    min-width: 0;
  }

  .dash-col-narrow {
    flex: 2;
    min-width: 0;
  }

  .view-all {
    display: block;
    text-align: center;
    padding: 8px;
    font-size: 12px;
    color: var(--text-link);
    text-decoration: none;
  }

  .view-all:hover {
    text-decoration: underline;
  }
</style>
