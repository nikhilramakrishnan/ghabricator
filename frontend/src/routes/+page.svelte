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
    <!-- Welcome box -->
    <div class="l-box">
      <div class="l-header">
        <i class="fa fa-code-fork l-header-icon"></i>
        <span class="l-header-title">Ghabricator</span>
        <span class="l-tag green">GitHub-powered</span>
      </div>
      <div class="l-welcome">
        <div class="l-welcome-text">
          <h1>Code review that doesn't suck.</h1>
          <p>Phabricator's review workflow on top of your GitHub repos. Side-by-side diffs, inline comments, Herald automation — no migration required.</p>
          <a href="/api/auth/github" class="l-signin">
            <i class="fa fa-github"></i> Sign in with GitHub
          </a>
        </div>
      </div>
    </div>

    <!-- Diff preview box -->
    <div class="l-box">
      <div class="l-header">
        <i class="fa fa-file-code-o l-header-icon"></i>
        <span class="l-header-title">src/auth/session.go</span>
        <span class="l-tag blue">+6 -1</span>
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
      <div class="l-inline">
        <div class="l-inline-bar"></div>
        <div class="l-inline-avatar">NR</div>
        <div class="l-inline-body">
          <span class="l-inline-author">nikhilr</span>
          <span class="l-inline-time">just now</span>
          <p>Should we also evict expired sessions on write? This only catches reads.</p>
        </div>
      </div>
    </div>

    <!-- Features box -->
    <div class="l-box">
      <div class="l-header">
        <i class="fa fa-star l-header-icon"></i>
        <span class="l-header-title">Features</span>
      </div>
      <div class="l-features">
        <div class="l-feat-item">
          <div class="l-feat-bar blue"></div>
          <i class="fa fa-columns l-feat-icon"></i>
          <div class="l-feat-content">
            <div class="l-feat-name">Side-by-side diffs</div>
            <div class="l-feat-desc">Syntax-highlighted, two-column diff viewer with context expansion</div>
          </div>
        </div>
        <div class="l-feat-item">
          <div class="l-feat-bar green"></div>
          <i class="fa fa-commenting-o l-feat-icon"></i>
          <div class="l-feat-content">
            <div class="l-feat-name">Inline comments</div>
            <div class="l-feat-desc">Click any line to comment. Drafts, threaded replies, batch reviews</div>
          </div>
        </div>
        <div class="l-feat-item">
          <div class="l-feat-bar violet"></div>
          <i class="fa fa-bullhorn l-feat-icon"></i>
          <div class="l-feat-content">
            <div class="l-feat-name">Herald automation</div>
            <div class="l-feat-desc">Auto-assign reviewers, add labels, post comments based on rules</div>
          </div>
        </div>
        <div class="l-feat-item">
          <div class="l-feat-bar orange"></div>
          <i class="fa fa-search l-feat-icon"></i>
          <div class="l-feat-content">
            <div class="l-feat-name">Search everything</div>
            <div class="l-feat-desc">PRs, issues, code, and repos across all your GitHub organizations</div>
          </div>
        </div>
        <div class="l-feat-item">
          <div class="l-feat-bar blue"></div>
          <i class="fa fa-code l-feat-icon"></i>
          <div class="l-feat-content">
            <div class="l-feat-name">Repository browser</div>
            <div class="l-feat-desc">Browse files, view blame, syntax highlighting — without leaving</div>
          </div>
        </div>
        <div class="l-feat-item">
          <div class="l-feat-bar green"></div>
          <i class="fa fa-github l-feat-icon"></i>
          <div class="l-feat-content">
            <div class="l-feat-name">GitHub is the source of truth</div>
            <div class="l-feat-desc">No database, no migration. Ghabricator is just a better lens</div>
          </div>
        </div>
      </div>
    </div>
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
  /* ===== Landing page — PHUI style ===== */
  .landing {
    max-width: 780px;
    margin: 0 auto;
    padding: 16px 16px 48px;
  }

  /* Box — mirrors PHUI Box */
  .l-box {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 4px;
    margin-bottom: 8px;
    overflow: hidden;
  }

  /* Header bar — mirrors PHUI HeaderView */
  .l-header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 12px;
    background: var(--bg-card-header);
    border-bottom: 1px solid var(--border-subtle);
  }
  .l-header-icon {
    font-size: 14px;
    color: var(--text-muted);
  }
  .l-header-title {
    font-size: 15px;
    font-weight: 600;
    color: var(--text);
  }

  /* Tag — mirrors PHUI Tag */
  .l-tag {
    font-size: 12px;
    font-weight: 600;
    line-height: 22px;
    padding: 0 8px;
    border-radius: 3px;
    margin-left: auto;
  }
  .l-tag.green {
    background: var(--tag-green-bg);
    color: var(--tag-green-text);
  }
  .l-tag.blue {
    background: var(--tag-blue-bg);
    color: var(--tag-blue-text);
  }

  /* Welcome content */
  .l-welcome {
    padding: 32px 24px;
  }
  .l-welcome h1 {
    font-size: 22px;
    font-weight: 700;
    margin: 0 0 12px;
    color: var(--text);
    line-height: 1.3;
  }
  .l-welcome p {
    font-size: 14px;
    color: var(--text-muted);
    line-height: 1.6;
    margin: 0 0 24px;
    max-width: 520px;
  }
  .l-signin {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: var(--green);
    color: var(--text-on-dark);
    font-size: 13px;
    font-weight: 600;
    padding: 8px 20px;
    border-radius: 3px;
    text-decoration: none;
  }
  .l-signin:hover {
    background: var(--green-hover);
    text-decoration: none;
  }

  /* Diff mock */
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
  .diff-mock .code { width: 100%; }
  .diff-mock .ctx { background: transparent; }
  .diff-mock .add { background: var(--diff-add-bg); }
  .diff-mock .add .ln { background: var(--diff-add-num-bg); }
  .diff-mock .del { background: var(--diff-del-bg); }
  .diff-mock .del .ln { background: var(--diff-del-num-bg); }

  /* Inline comment — mirrors Timeline event */
  .l-inline {
    display: flex;
    gap: 10px;
    padding: 10px 12px;
    border-top: 1px solid var(--border-subtle);
    position: relative;
  }
  .l-inline-bar {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 3px;
    background: var(--blue);
  }
  .l-inline-avatar {
    width: 32px;
    height: 32px;
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
  .l-inline-body {
    font-size: 13px;
    line-height: 1.5;
    min-width: 0;
  }
  .l-inline-author {
    font-weight: 600;
    color: var(--text);
    margin-right: 6px;
  }
  .l-inline-time {
    font-size: 12px;
    color: var(--text-muted);
  }
  .l-inline-body p {
    margin: 4px 0 0;
    color: var(--text-muted);
  }

  /* Features list — mirrors ObjectItemList */
  .l-features {
    /* no extra padding, items handle it */
  }
  .l-feat-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 12px;
    border-bottom: 1px solid var(--border-subtle);
    position: relative;
  }
  .l-feat-item:last-child {
    border-bottom: none;
  }
  .l-feat-bar {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 3px;
  }
  .l-feat-bar.blue { background: var(--blue); }
  .l-feat-bar.green { background: var(--green); }
  .l-feat-bar.violet { background: var(--violet); }
  .l-feat-bar.orange { background: var(--orange); }
  .l-feat-icon {
    font-size: 14px;
    color: var(--text-muted);
    width: 18px;
    text-align: center;
    flex-shrink: 0;
  }
  .l-feat-content {
    min-width: 0;
  }
  .l-feat-name {
    font-size: 14px;
    font-weight: 600;
    color: var(--text-link);
  }
  .l-feat-desc {
    font-size: 12px;
    color: var(--text-muted);
    margin-top: 1px;
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
