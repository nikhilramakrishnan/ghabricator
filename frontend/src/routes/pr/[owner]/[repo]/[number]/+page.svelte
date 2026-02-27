<script lang="ts">
  import { Breadcrumbs, FormationView } from '$lib/components/layout';
  import {
    Box, HeaderView, Tag, CurtainBox, PropertyList, StatusList,
    Timeline, Button
  } from '$lib/components/phui';
  import { DiffTable, ChangesetHeader, FileTree } from '$lib/components/diff';
  import type { APIReviewComment as DiffComment } from '$lib/components/diff';
  import { ReviewForm } from '$lib/components/review';
  import { apiPost } from '$lib/api';
  import { S } from '$lib/strings';
  import { addDraft } from '$lib/stores/inline';
  import type {
    APIPRDetailResponse, APIChangeset, APIReviewComment,
    APICheckRun, APIHeraldMatch
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

  // Collapsed state per changeset
  let collapsedFiles = $state(new Set<number>());
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

  // Scroll to changeset when clicking file tree
  function scrollToChangeset(path: string) {
    activePath = path;
    const cs = changesets.find((c) => c.displayPath === path);
    if (cs) {
      const el = document.getElementById(`C${cs.id}`);
      if (el) {
        el.scrollIntoView({ behavior: 'smooth', block: 'start' });
        // Uncollapse if collapsed
        if (collapsedFiles.has(cs.id)) {
          collapsedFiles.delete(cs.id);
          collapsedFiles = new Set(collapsedFiles);
        }
      }
    }
  }

  // Status badge
  function statusBadge(p: typeof pr): { text: string; color: string } {
    if (p.merged) return { text: 'Merged', color: 'violet' };
    if (p.state === 'closed') return { text: 'Closed', color: 'red' };
    if (p.draft) return { text: 'Draft', color: 'grey' };
    return { text: 'Open', color: 'green' };
  }

  let status = $derived(statusBadge(pr));

  // Transform API comments to DiffTable format (flat author/avatarURL)
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
      case 'APPROVED': return { text: 'Accepted', icon: 'fa-check', shade: 'green' };
      case 'CHANGES_REQUESTED': return { text: 'Changes Requested', icon: 'fa-times', shade: 'red' };
      case 'COMMENTED': return { text: 'Commented', icon: 'fa-comment', shade: 'blue' };
      case 'DISMISSED': return { text: 'Dismissed', icon: 'fa-ban', shade: 'grey' };
      default: return { text: 'Waiting', icon: 'fa-clock-o', shade: 'orange' };
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

  // Convert timeline events
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
  let mergeMethod = $state('squash');
  let actionLoading = $state(false);

  async function handleMerge() {
    if (actionLoading) return;
    actionLoading = true;
    try {
      await apiPost('/api/v2/merge', { owner, repo, number, mergeMethod });
      window.location.reload();
    } catch (e: unknown) {
      alert(e instanceof Error ? e.message : 'Merge failed');
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
      alert(e instanceof Error ? e.message : 'Failed');
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

  // Action icon for herald
  function heraldActionIcon(type: string): string {
    switch (type) {
      case 'add_reviewer': return 'fa-user-plus';
      case 'add_label': return 'fa-tag';
      case 'post_comment': return 'fa-comment';
      default: return 'fa-bolt';
    }
  }

  function heraldActionLabel(type: string): string {
    switch (type) {
      case 'add_reviewer': return 'Add reviewer';
      case 'add_label': return 'Add label';
      case 'post_comment': return 'Post comment';
      default: return type;
    }
  }
</script>

<div class="pr-page-header">
  <Breadcrumbs {crumbs} />
  <div class="pr-title-row">
    <h1 class="pr-title">
      <Tag shade={status.color}>{status.text}</Tag>
      {' '}
      <span class="pr-number">D{number}</span>
      {' '}{pr.title}
    </h1>
  </div>
</div>

<FormationView>
  {#snippet filetree()}
    <FileTree {changesets} activeFile={activePath} />
  {/snippet}

  <!-- Main content -->
  {#if pr.body?.trim()}
    <Box border>
      <HeaderView title={S.pr.summary} icon="fa-file-text-o" />
      <div class="summary-body">
        <div class="remarkup-content">
          {@html pr.body}
        </div>
      </div>
    </Box>
  {/if}

  {#each changesets as cs (cs.id)}
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

  {#if timelineEvents.length > 0}
    <Timeline events={timelineEvents} />
  {/if}

  <ReviewForm {owner} {repo} {number} />

  {#snippet curtain()}
    <!-- Summary -->
    <CurtainBox title={S.pr.summary}>
      <PropertyList items={[
        { label: S.pr.author, value: pr.author.login },
        { label: S.pr.status, value: status.text },
        { label: S.pr.created, value: pr.createdAt },
        { label: S.pr.updated, value: pr.updatedAt },
        { label: S.pr.base, value: pr.base.ref },
        { label: S.pr.head, value: pr.head.ref },
        { label: S.pr.changes, value: `+${pr.additions} / -${pr.deletions} in ${pr.changedFiles} files` }
      ]} />
    </CurtainBox>

    <!-- Reviewers -->
    <CurtainBox title={S.pr.reviewers}>
      {#if pr.reviewers?.length}
        {#each pr.reviewers as reviewer}
          {@const rs = reviewStateForUser(reviewer.login)}
          {@const display = reviewStateDisplay(rs)}
          <div class="reviewer-row">
            {#if reviewer.avatarURL}
              <img src={reviewer.avatarURL} alt="" class="reviewer-avatar" />
            {/if}
            <span class="reviewer-name">{reviewer.login}</span>
            <Tag shade={display.shade} icon={display.icon}>{display.text}</Tag>
          </div>
        {/each}
      {:else}
        <div class="empty-curtain">None assigned</div>
      {/if}
    </CurtainBox>

    <!-- Buildables -->
    {#if checkRuns.length > 0}
      <CurtainBox title={S.pr.buildables}>
        <StatusList items={checkRuns.map((cr) => {
          const d = checkRunDisplay(cr);
          return { name: d.name, icon: d.icon, color: d.color, href: cr.detailsURL || undefined };
        })} />
      </CurtainBox>
    {/if}

    <!-- Herald -->
    <CurtainBox title={S.pr.herald}>
      {#if heraldMatches.length > 0}
        {#each heraldMatches as match}
          <div class="herald-match">
            <i class="fa fa-check mrs herald-check"></i>
            <a href="/herald/{match.ruleId}"><strong>{match.ruleName}</strong></a>
            {#if match.actions?.length}
              <div class="herald-actions">
                {#each match.actions as action}
                  <div>
                    <i class="fa {heraldActionIcon(action.type)} mrs"></i>
                    {heraldActionLabel(action.type)}: {action.value}
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/each}
      {:else}
        <div class="empty-curtain">No rules matched.</div>
      {/if}
      <div class="herald-manage">
        <a href="/herald"><i class="fa fa-list mrs"></i>Manage Rules</a>
      </div>
    </CurtainBox>

    <!-- Labels -->
    {#if pr.labels?.length}
      <CurtainBox title={S.pr.labels}>
        {#each pr.labels as label}
          <Tag shade={labelShade(label.color)}>{label.name}</Tag>
          {' '}
        {/each}
      </CurtainBox>
    {/if}

    <!-- Actions -->
    {#if !pr.merged}
      <div class="curtain-actions">
        {#if pr.state !== 'closed'}
          <div class="merge-row">
            <select bind:value={mergeMethod} class="merge-select">
              <option value="squash">Squash</option>
              <option value="merge">Merge</option>
              <option value="rebase">Rebase</option>
            </select>
            <Button color="green" icon="fa-check-circle" disabled={actionLoading} onclick={handleMerge}>
              Land Revision
            </Button>
          </div>
          <Button color="default" icon="fa-times-circle" disabled={actionLoading} onclick={() => handleClose('closed')}>
            Close
          </Button>
        {:else}
          <Button color="green" icon="fa-refresh" disabled={actionLoading} onclick={() => handleClose('open')}>
            Reopen
          </Button>
        {/if}
      </div>
    {/if}
  {/snippet}
</FormationView>

<style>
  .pr-page-header {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 16px;
  }

  .pr-title-row {
    padding: 8px 0 12px;
  }

  .pr-title {
    font-size: 20px;
    font-weight: 600;
    color: var(--text);
    margin: 0;
    line-height: 1.4;
  }

  .pr-number {
    color: var(--text-muted);
    font-weight: normal;
  }

  .summary-body {
    padding: 10px 12px;
  }

  .remarkup-content {
    font-size: 13px;
    line-height: 1.5;
    overflow-wrap: break-word;
    word-break: break-word;
  }

  .reviewer-row {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
  }

  .reviewer-avatar {
    width: 24px;
    height: 24px;
    border-radius: 3px;
  }

  .reviewer-name {
    font-size: 13px;
  }

  .empty-curtain {
    font-size: 13px;
    color: var(--text-muted);
  }

  .herald-match {
    font-size: 12px;
    margin-bottom: 6px;
  }

  .herald-check {
    color: var(--green);
  }

  .herald-actions {
    font-size: 11px;
    color: var(--text-muted);
    margin-left: 20px;
  }

  .herald-manage {
    margin-top: 8px;
    font-size: 12px;
  }

  .curtain-actions {
    padding: 12px;
    border-top: 1px solid var(--border-subtle);
  }

  .merge-row {
    display: flex;
    gap: 8px;
    align-items: center;
    margin-bottom: 8px;
  }

  .merge-select {
    font-size: 12px;
    padding: 4px 6px;
    border: 1px solid var(--border);
    border-radius: 3px;
    background: var(--bg-card);
    color: var(--text);
    flex: 1;
  }
</style>
