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
    <div class="hero-card">
      <div class="hero-icon">
        <i class="fa fa-code-fork"></i>
      </div>
      <h2 class="hero-title">{S.landing.title}</h2>
      <p class="hero-desc">{S.landing.desc}</p>
      <a href="/api/auth/github" class="signin-btn">
        <i class="fa fa-github"></i> {S.landing.signIn}
      </a>
      <div class="features">
        <span><i class="fa fa-search mrs"></i> {S.landing.sideDiffs}</span>
        <span><i class="fa fa-commenting-o mrs"></i> {S.landing.inlineComments}</span>
        <span><i class="fa fa-bullhorn mrs"></i> {S.landing.heraldRules}</span>
        <span><i class="fa fa-clipboard mrs"></i> {S.landing.pasteBin}</span>
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

  .dashboard-grid {
    display: flex;
    gap: 16px;
  }

  .dash-col {
    flex: 1;
    min-width: 0;
  }

  .dash-col-narrow {
    flex: 0 0 380px;
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
