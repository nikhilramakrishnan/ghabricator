<script lang="ts">
  import { Breadcrumbs, FormationView } from '$lib/components/layout';
  import { Box, HeaderView, Tag, Timeline, Button } from '$lib/components/phui';
  import { DiffTable, ChangesetHeader, CommitHistory, FileTree } from '$lib/components/diff';
  import type { APIReviewComment as DiffComment } from '$lib/components/diff';
  import { ReviewForm } from '$lib/components/review';
  import { apiFetch, apiPost } from '$lib/api';
  import { S } from '$lib/strings';
  import { addDraft } from '$lib/stores/inline';
  import type {
    APIPRDetailResponse, APIChangeset, APIReviewComment,
    APICheckRun, APIHeraldMatch, APICommit
  } from '$lib/types';
  import type { TimelineEvent } from '$lib/components/phui';

  let { data } = $props();
  let owner: string = $derived(data.owner);
  let repo: string = $derived(data.repo);
  let number: number = $derived(data.number);
  let resp: APIPRDetailResponse = $derived(data.data);
  let pr = $derived(resp.pr);
  let changesets: APIChangeset[] = $derived(resp.changesets ?? []);
  let commentsByPath = $derived(resp.commentsByPath ?? {});
  let reviews = $derived(resp.reviews ?? []);
  let checkRuns: APICheckRun[] = $derived(resp.checkRuns ?? []);
  let timeline = $derived(resp.timeline ?? []);
  let heraldMatches: APIHeraldMatch[] = $derived(resp.heraldMatches ?? []);
  let commits: APICommit[] = $derived(resp.commits ?? []);

  // Interdiff state
  let compareBase = $state<string | null>(null);
  let compareHead = $state<string | null>(null);
  let interdiffChangesets = $state<APIChangeset[] | null>(null);
  let interdiffLoading = $state(false);

  let displayChangesets = $derived(interdiffChangesets ?? changesets);

  async function handleRangeChange(base: string | null, head: string | null) {
    compareBase = base;
    compareHead = head;
    if (base === null && head === null) {
      interdiffChangesets = null;
      return;
    }
    interdiffLoading = true;
    try {
      const params = new URLSearchParams();
      if (base) params.set('base', base);
      if (head) params.set('head', head);
      const resp = await apiFetch<{ changesets: APIChangeset[] }>(
        `/api/pr/${owner}/${repo}/${number}/compare?${params}`
      );
      interdiffChangesets = resp.changesets ?? [];
    } catch {
      interdiffChangesets = null;
      compareBase = null;
      compareHead = null;
    } finally {
      interdiffLoading = false;
    }
  }

  // Collapsed state — auto-fold changesets > 70 rows
  let collapsedFiles = $state(
    new Set(changesets.filter(cs => cs.rows.length > 70).map(cs => cs.id))
  );
  function toggleCollapse(id: number) {
    if (collapsedFiles.has(id)) {
      collapsedFiles.delete(id);
    } else {
      collapsedFiles.add(id);
    }
    collapsedFiles = new Set(collapsedFiles);
  }

  // Active file path for file tree highlighting
  let activePath = $state('');

  function scrollToChangeset(path: string) {
    activePath = path;
    const cs = displayChangesets.find((c) => c.displayPath === path);
    if (cs) {
      const el = document.getElementById(`C${cs.id}`);
      if (el) {
        el.scrollIntoView({ behavior: 'smooth', block: 'start' });
        if (collapsedFiles.has(cs.id)) {
          collapsedFiles.delete(cs.id);
          collapsedFiles = new Set(collapsedFiles);
        }
      }
    }
  }

  // Status badge
  function statusBadge(p: typeof pr): { text: string; color: string } {
    if (p.merged) return { text: S.pr.statusMerged, color: 'violet' };
    if (p.state === 'closed') return { text: S.pr.statusClosed, color: 'red' };
    if (p.draft) return { text: S.pr.statusDraft, color: 'grey' };
    return { text: S.pr.statusOpen, color: 'green' };
  }

  let status = $derived(statusBadge(pr));

  // Transform API comments to DiffTable format
  function flattenComments(comments: APIReviewComment[]): DiffComment[] {
    return (comments ?? []).map((c) => ({
      id: c.id,
      author: c.author.login,
      avatarURL: c.author.avatarURL,
      body: c.body,
      path: c.path,
      line: c.line,
      side: c.side,
      createdAt: c.createdAt
    }));
  }

  // Check run icon/color mapping
  function checkRunDisplay(cr: APICheckRun): { icon: string; color: string; name: string } {
    let icon: string;
    let color: string;
    if (cr.status !== 'completed') {
      icon = cr.status === 'in_progress' ? 'fa-circle-o-notch' : 'fa-clock-o';
      color = cr.status === 'in_progress' ? 'var(--orange)' : 'var(--text-muted)';
    } else {
      switch (cr.conclusion) {
        case 'success': icon = 'fa-check-circle'; color = 'var(--green)'; break;
        case 'failure': icon = 'fa-times-circle'; color = 'var(--red)'; break;
        case 'cancelled': icon = 'fa-ban'; color = 'var(--text-muted)'; break;
        case 'skipped': icon = 'fa-minus-circle'; color = 'var(--text-muted)'; break;
        case 'timed_out': icon = 'fa-clock-o'; color = 'var(--red)'; break;
        case 'action_required': icon = 'fa-exclamation-circle'; color = 'var(--orange)'; break;
        default: icon = 'fa-question-circle'; color = 'var(--text-muted)';
      }
    }
    let name = cr.name;
    if (cr.appName && cr.appName !== 'GitHub Actions') {
      name = cr.appName + ' / ' + name;
    }
    return { icon, color, name };
  }

  // Review state helpers
  function reviewStateForUser(login: string): string {
    let latest: { state: string; createdAt: string } | null = null;
    for (const r of reviews) {
      if (r.author.login === login) {
        if (!latest || r.createdAt > latest.createdAt) {
          latest = r;
        }
      }
    }
    return latest?.state ?? '';
  }

  function reviewStateDisplay(state: string): { text: string; icon: string; shade: string } {
    switch (state) {
      case 'APPROVED': return { text: S.pr.reviewAccepted, icon: 'fa-check', shade: 'green' };
      case 'CHANGES_REQUESTED': return { text: S.pr.reviewChangesRequested, icon: 'fa-times', shade: 'red' };
      case 'COMMENTED': return { text: S.pr.reviewCommented, icon: 'fa-comment', shade: 'blue' };
      case 'DISMISSED': return { text: S.pr.reviewDismissed, icon: 'fa-ban', shade: 'grey' };
      default: return { text: S.pr.reviewWaiting, icon: 'fa-clock-o', shade: 'orange' };
    }
  }

  // Label shade from hex color
  function labelShade(hex: string): string {
    if (!hex || hex.length < 6) return 'blue';
    const r = parseInt(hex.slice(0, 2), 16);
    const g = parseInt(hex.slice(2, 4), 16);
    const b = parseInt(hex.slice(4, 6), 16);
    if (r > 180 && g < 120 && b < 120) return 'red';
    if (g > 180 && r < 120) return 'green';
    if (b > 180 && r < 150) return 'blue';
    if (r > 180 && g > 180 && b < 120) return 'yellow';
    if (r > 150 && g > 80 && g < 170 && b < 100) return 'orange';
    if (r > 80 && b > 150 && g < 120) return 'violet';
    if (r > 200 && g > 200 && b > 200) return 'grey';
    return 'blue';
  }

  // Timeline events
  let timelineEvents: TimelineEvent[] = $derived(
    timeline.map((ev) => ({
      author: { login: ev.author.login, avatarURL: ev.author.avatarURL },
      action: ev.action,
      body: ev.body,
      createdAt: ev.createdAt,
      iconClass: ev.iconClass,
      iconColor: ev.iconColor
    }))
  );

  // Merge/close actions
  let actionLoading = $state(false);

  async function handleMerge() {
    if (actionLoading) return;
    actionLoading = true;
    try {
      await apiPost('/api/v2/merge', { owner, repo, number, mergeMethod: 'squash' });
      window.location.reload();
    } catch (e: unknown) {
      alert(e instanceof Error ? e.message : S.pr.mergeFailed);
    } finally {
      actionLoading = false;
    }
  }

  async function handleClose(newState: string) {
    if (actionLoading) return;
    actionLoading = true;
    try {
      await apiPost('/api/v2/close', { owner, repo, number, state: newState });
      window.location.reload();
    } catch (e: unknown) {
      alert(e instanceof Error ? e.message : S.pr.actionFailed);
    } finally {
      actionLoading = false;
    }
  }

  function handleNewComment(path: string, line: number, side: string) {
    addDraft(path, line, side);
  }

  let crumbs = $derived([
    { name: S.crumb.home, href: '/' },
    { name: S.revisions.title, href: '/dashboard' },
    { name: `${owner}/${repo}`, href: `/repos/${owner}/${repo}` },
    { name: `D${number}` }
  ]);
</script>

<div class="pr-page-header">
  <Breadcrumbs {crumbs} />
  <h1 class="pr-title">
    <Tag shade={status.color}>{status.text}</Tag>
    {pr.title}
  </h1>
</div>

<FormationView>
  {#snippet filetree()}
    <FileTree changesets={displayChangesets} activeFile={activePath} />
  {/snippet}

  <!-- Revision Contents — Phabricator-style property card -->
  <Box border>
    <HeaderView title={S.pr.revisionContents} icon="fa-file-text-o" />
    <div class="plist">
      <div class="plist-row">
        <span class="plist-key">{S.pr.author}</span>
        <span class="plist-val">
          {#if pr.author.avatarURL}
            <img src={pr.author.avatarURL} alt="" class="plist-avatar" />
          {/if}
          {pr.author.login}
        </span>
      </div>
      {#if pr.reviewers?.length}
        <div class="plist-row">
          <span class="plist-key">{S.pr.reviewers}</span>
          <span class="plist-val">
            {#each pr.reviewers as reviewer, i}
              {#if i > 0}<span class="plist-comma">,</span>{/if}
              {@const rs = reviewStateForUser(reviewer.login)}
              {@const display = reviewStateDisplay(rs)}
              <span class="reviewer-chip">
                {#if reviewer.avatarURL}
                  <img src={reviewer.avatarURL} alt="" class="plist-avatar" />
                {/if}
                {reviewer.login}
                <Tag shade={display.shade} icon={display.icon}>{display.text}</Tag>
              </span>
            {/each}
          </span>
        </div>
      {/if}
      <div class="plist-row">
        <span class="plist-key">{S.pr.repository}</span>
        <span class="plist-val">
          <span class="prop-ref">{pr.base.ref}</span>
          <i class="fa fa-long-arrow-left prop-arrow"></i>
          <span class="prop-ref">{pr.head.ref}</span>
        </span>
      </div>
      <div class="plist-row">
        <span class="plist-key">{S.pr.changes}</span>
        <span class="plist-val">
          <span class="prop-add">+{pr.additions}</span>{' '}
          <span class="prop-del">&minus;{pr.deletions}</span>{' '}
          in {pr.changedFiles} files
        </span>
      </div>
      {#if checkRuns.length > 0}
        <div class="plist-row">
          <span class="plist-key">{S.pr.buildables}</span>
          <span class="plist-val checks-val">
            {#each checkRuns as cr}
              {@const d = checkRunDisplay(cr)}
              {#if cr.detailsURL}
                <a href={cr.detailsURL} target="_blank" rel="noopener" class="check-chip">
                  <i class="fa {d.icon}" style="color:{d.color}"></i>
                  {d.name}
                </a>
              {:else}
                <span class="check-chip">
                  <i class="fa {d.icon}" style="color:{d.color}"></i>
                  {d.name}
                </span>
              {/if}
            {/each}
          </span>
        </div>
      {/if}
      {#if pr.labels?.length}
        <div class="plist-row">
          <span class="plist-key">{S.pr.labels}</span>
          <span class="plist-val">
            {#each pr.labels as label}
              <Tag shade={labelShade(label.color)}>{label.name}</Tag>
              {' '}
            {/each}
          </span>
        </div>
      {/if}
    </div>

    {#if pr.body?.trim()}
      <div class="summary-section">
        <div class="remarkup-content">
          {@html pr.body}
        </div>
      </div>
    {/if}
  </Box>

  {#if compareBase || compareHead}
    <div class="interdiff-indicator">
      <i class="fa fa-exchange"></i>
      {S.pr.showingChanges} {compareBase ? compareBase.slice(0, 7) : pr.base.ref}..{compareHead ? compareHead.slice(0, 7) : S.diff.latest.toLowerCase()}
    </div>
  {/if}

  {#if interdiffLoading}
    <div class="interdiff-loading">
      <i class="fa fa-circle-o-notch fa-spin"></i> {S.pr.loadingDiff}
    </div>
  {/if}

  <!-- Diffs -->
  {#each displayChangesets as cs (cs.id)}
    {@const collapsed = collapsedFiles.has(cs.id)}
    <div id="C{cs.id}">
      <ChangesetHeader changeset={cs} {collapsed} onToggle={() => toggleCollapse(cs.id)} />
      {#if !collapsed}
        <DiffTable
          changeset={cs}
          comments={flattenComments(commentsByPath[cs.displayPath] ?? [])}
          onNewComment={handleNewComment}
        />
      {/if}
    </div>
  {/each}

  {#if commits.length > 0}
    <CommitHistory
      {commits}
      baseBranch={pr.base.ref}
      onRangeChange={handleRangeChange}
    />
  {/if}

  {#if timelineEvents.length > 0}
    <Timeline events={timelineEvents} />
  {/if}

  <ReviewForm {owner} {repo} {number} />

  {#if !pr.merged}
    <div class="action-row">
      {#if pr.state !== 'closed'}
        <Button color="green" icon="fa-check-circle" disabled={actionLoading} onclick={handleMerge}>
          {S.pr.landRevision}
        </Button>
        <Button color="default" icon="fa-times-circle" disabled={actionLoading} onclick={() => handleClose('closed')}>
          {S.pr.close}
        </Button>
      {:else}
        <Button color="green" icon="fa-refresh" disabled={actionLoading} onclick={() => handleClose('open')}>
          {S.pr.reopen}
        </Button>
      {/if}
    </div>
  {/if}
</FormationView>

<style>
  .pr-page-header {
    padding: 0 16px;
  }

  .pr-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--text);
    margin: 0;
    padding: 6px 0 4px;
    line-height: 1.4;
    display: flex;
    align-items: baseline;
    gap: 8px;
    flex-wrap: wrap;
    min-width: 0;
  }

  /* Property list — Phabricator-style key/value rows */
  .plist {
    padding: 8px 16px;
  }

  .plist-row {
    display: flex;
    align-items: baseline;
    padding: 4px 0;
    font-size: 13px;
    border-bottom: 1px solid var(--border-subtle);
  }
  .plist-row:last-child {
    border-bottom: none;
  }

  .plist-key {
    width: 100px;
    flex-shrink: 0;
    font-size: 12px;
    font-weight: 600;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.03em;
  }

  .plist-val {
    flex: 1;
    color: var(--text);
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 6px;
  }

  .plist-avatar {
    width: 20px;
    height: 20px;
    border-radius: 3px;
    vertical-align: middle;
  }

  .plist-comma {
    margin-right: 4px;
  }

  .reviewer-chip {
    display: inline-flex;
    align-items: center;
    gap: 5px;
  }

  .prop-ref {
    font-family: var(--font-mono);
    font-size: 12px;
    background: var(--bg-subtle);
    padding: 2px 6px;
    border-radius: 3px;
    color: var(--text);
  }

  .prop-arrow {
    font-size: 10px;
    margin: 0 4px;
    color: var(--text-muted);
  }

  .prop-add {
    color: var(--green);
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
  }
  .prop-del {
    color: var(--red);
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
  }

  .checks-val {
    gap: 10px;
  }

  .check-chip {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    color: var(--text-muted);
    text-decoration: none;
    font-size: 13px;
  }
  a.check-chip:hover {
    color: var(--text);
    text-decoration: none;
  }

  /* Summary body — below property list, separated by border */
  .summary-section {
    padding: 12px 16px;
    border-top: 1px solid var(--border-subtle);
  }

  /* Action buttons — outside the card */
  .action-row {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 0;
  }

  .interdiff-indicator {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 12px;
    margin-bottom: 8px;
    font-size: 12px;
    font-family: var(--font-mono);
    color: var(--blue);
    background: var(--tag-blue-bg);
    border: 1px solid var(--border);
    border-radius: 4px;
  }

  .interdiff-loading {
    padding: 12px;
    text-align: center;
    font-size: 13px;
    color: var(--text-muted);
  }
</style>
