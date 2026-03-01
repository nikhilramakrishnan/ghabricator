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
    <div class="l-hero">
      <div class="l-author-note">
        <img class="l-author-avatar" src="https://avatars.githubusercontent.com/nikhilr?s=72" alt="nikhilr" />
        <div class="l-author-bubble">
          <span class="l-author-name">nikhilr</span>
          <p>I missed Phabricator so I made my own.</p>
        </div>
      </div>
      <a href="/api/auth/github" class="l-signin">
        <i class="fa fa-github"></i> Log in to Ghabricator
      </a>
    </div>

    <!-- Diff preview box -->
    <div class="l-box">
      <div class="l-header">
        <i class="fa fa-file-code-o l-header-icon"></i>
        <span class="l-header-title">src/auth/session.go</span>
        <span class="l-tag blue">+6 -1</span>
      </div>
      <table class="diff-mock sbs">
        <colgroup>
          <col class="dm-num" />
          <col class="dm-code" />
          <col class="dm-num" />
          <col class="dm-code" />
        </colgroup>
        <tbody>
          <tr><td class="n">14</td><td class="ctx">func (s *Store) Get(id string) *Session {'{'}</td><td class="n">14</td><td class="ctx">func (s *Store) Get(id string) *Session {'{'}</td></tr>
          <tr><td class="n old">15</td><td class="old">    sess, ok := s.cache[id]</td><td class="n new"></td><td class="new">    sess, ok := s.sessions[id]</td></tr>
          <tr><td class="n">16</td><td class="ctx">    if !ok {'{'}</td><td class="n">16</td><td class="ctx">    if !ok {'{'}</td></tr>
          <tr><td class="n">17</td><td class="ctx">        return nil</td><td class="n">17</td><td class="ctx">        return nil</td></tr>
          <tr><td class="n">18</td><td class="ctx">    {'}'}</td><td class="n">18</td><td class="ctx">    {'}'}</td></tr>
          <tr><td class="n"></td><td class="blank"></td><td class="n new">19</td><td class="new">    if time.Since(sess.CreatedAt) > sessionTTL {'{'}</td></tr>
          <tr><td class="n"></td><td class="blank"></td><td class="n new">20</td><td class="new">        return nil</td></tr>
          <tr><td class="n"></td><td class="blank"></td><td class="n new">21</td><td class="new">    {'}'}</td></tr>
          <tr><td class="n">19</td><td class="ctx">    return sess</td><td class="n">22</td><td class="ctx">    return sess</td></tr>
          <tr><td class="n">20</td><td class="ctx">{'}'}</td><td class="n">23</td><td class="ctx">{'}'}</td></tr>
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
            <div class="l-feat-desc">PRs, issues, code, and repos across all your organizations</div>
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
          <i class="fa fa-bolt l-feat-icon"></i>
          <div class="l-feat-content">
            <div class="l-feat-name">Zero setup</div>
            <div class="l-feat-desc">No database, no migration, no config. Sign in and go</div>
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
    max-width: 860px;
    margin: 0 auto;
    padding: 0 16px 48px;
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

  /* Tag */
  .l-tag {
    font-size: 12px;
    font-weight: 600;
    line-height: 22px;
    padding: 0 8px;
    border-radius: 3px;
    margin-left: auto;
  }
  .l-tag.blue {
    background: var(--tag-blue-bg);
    color: var(--tag-blue-text);
  }

  /* Hero */
  .l-hero {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 28px 24px 16px;
  }
  .l-author-note {
    display: flex;
    align-items: flex-start;
    gap: 10px;
    text-align: left;
    margin-bottom: 20px;
  }
  .l-author-avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    flex-shrink: 0;
  }
  .l-author-bubble {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 8px 12px;
    position: relative;
    font-size: 13px;
    line-height: 1.5;
  }
  .l-author-bubble::before {
    content: '';
    position: absolute;
    left: -6px;
    top: 12px;
    width: 10px;
    height: 10px;
    background: var(--bg-card);
    border-left: 1px solid var(--border);
    border-bottom: 1px solid var(--border);
    transform: rotate(45deg);
  }
  .l-author-name {
    font-weight: 600;
    color: var(--text);
    margin-right: 4px;
  }
  .l-author-bubble p {
    margin: 2px 0 0;
    color: var(--text-muted);
  }
  .l-signin {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: var(--green);
    color: var(--text-on-dark);
    font-size: 14px;
    font-weight: 600;
    padding: 10px 28px;
    border-radius: 3px;
    text-decoration: none;
  }
  .l-signin:hover {
    background: var(--green-hover);
    text-decoration: none;
  }

  /* Diff mock — side-by-side */
  .diff-mock.sbs {
    width: 100%;
    border-collapse: collapse;
    font-family: var(--font-mono);
    font-size: 11px;
    line-height: 1.5;
  }
  .diff-mock.sbs col.dm-num { width: 4em; }
  .diff-mock.sbs col.dm-code { width: calc(50% - 4em); }
  .diff-mock.sbs td {
    padding: 1px 8px;
    white-space: pre;
  }
  .diff-mock.sbs td.n {
    color: var(--text-muted);
    text-align: right;
    user-select: none;
    padding: 1px 4px;
  }
  .diff-mock.sbs td.ctx { background: transparent; }
  .diff-mock.sbs td.old { background: var(--diff-del-bg); }
  .diff-mock.sbs td.n.old { background: var(--diff-del-num-bg); }
  .diff-mock.sbs td.new { background: var(--diff-add-bg); }
  .diff-mock.sbs td.n.new { background: var(--diff-add-num-bg); }
  .diff-mock.sbs td.blank { background: var(--bg-hover); }

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
