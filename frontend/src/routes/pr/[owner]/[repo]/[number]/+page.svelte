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
      color = cr.status === 'in_progress' ? '#c69026' : '#6b748c';
    } else {
      switch (cr.conclusion) {
        case 'success': icon = 'fa-check-circle'; color = '#139543'; break;
        case 'failure': icon = 'fa-times-circle'; color = '#c0392b'; break;
        case 'cancelled': icon = 'fa-ban'; color = '#6b748c'; break;
        case 'skipped': icon = 'fa-minus-circle'; color = '#6b748c'; break;
        case 'timed_out': icon = 'fa-clock-o'; color = '#c0392b'; break;
        case 'action_required': icon = 'fa-exclamation-circle'; color = '#c69026'; break;
        default: icon = 'fa-question-circle'; color = '#6b748c';
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
    { name: 'Home', href: '/' },
    { name: 'Dashboard', href: '/dashboard' },
    { name: `${owner}/${repo}`, href: `/repos/${owner}/${repo}` },
    { name: `D${number}` }
  ]);
</script>

<div class="phui-two-column-view">
  <div class="phui-two-column-container">
    <Breadcrumbs {crumbs} />
    <div class="phui-two-column-header">
      <div class="phui-header-view">
        <div class="phui-header-shell">
          <h1 class="phui-header-header">
            <Tag shade={status.color}>{status.text}</Tag>
            {' '}
            <span class="phui-header-subheader" style="color:rgba(55,55,55,.6);font-weight:normal">D{number}</span>
            {' '}{pr.title}
          </h1>
        </div>
      </div>
    </div>
  </div>
</div>

<FormationView>
  {#snippet filetree()}
    <FileTree {changesets} activeFile={activePath} />
  {/snippet}

  <!-- Main content -->
  {#if pr.body?.trim()}
    <Box border>
      <HeaderView title="Summary" icon="fa-file-text-o" />
      <div style="padding:10px 12px">
        <div class="phabricator-remarkup" style="font-size:13px;line-height:1.5;overflow-wrap:break-word;word-break:break-word;">
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
    <CurtainBox title="Summary">
      <PropertyList items={[
        { label: 'Author', value: pr.author.login },
        { label: 'Status', value: status.text },
        { label: 'Created', value: pr.createdAt },
        { label: 'Updated', value: pr.updatedAt },
        { label: 'Base', value: pr.base.ref },
        { label: 'Head', value: pr.head.ref },
        { label: 'Changes', value: `+${pr.additions} / -${pr.deletions} in ${pr.changedFiles} files` }
      ]} />
    </CurtainBox>

    <!-- Reviewers -->
    <CurtainBox title="Reviewers">
      {#if pr.reviewers?.length}
        {#each pr.reviewers as reviewer}
          {@const rs = reviewStateForUser(reviewer.login)}
          {@const display = reviewStateDisplay(rs)}
          <div style="display:flex;align-items:center;gap:8px;margin-bottom:8px">
            {#if reviewer.avatarURL}
              <img src={reviewer.avatarURL} alt="" style="width:24px;height:24px;border-radius:3px" />
            {/if}
            <span style="font-size:13px">{reviewer.login}</span>
            <Tag shade={display.shade} icon={display.icon}>{display.text}</Tag>
          </div>
        {/each}
      {:else}
        <div style="font-size:13px;color:#6b748c">None assigned</div>
      {/if}
    </CurtainBox>

    <!-- Buildables -->
    {#if checkRuns.length > 0}
      <CurtainBox title="Buildables">
        <StatusList items={checkRuns.map((cr) => {
          const d = checkRunDisplay(cr);
          return { name: d.name, icon: d.icon, color: d.color, href: cr.detailsURL || undefined };
        })} />
      </CurtainBox>
    {/if}

    <!-- Herald -->
    <CurtainBox title="Herald">
      {#if heraldMatches.length > 0}
        {#each heraldMatches as match}
          <div style="font-size:12px;margin-bottom:6px">
            <span class="phui-icon-view phui-font-fa fa-check mrs" style="color:#139543"></span>
            <a href="/herald/{match.ruleId}"><strong>{match.ruleName}</strong></a>
            {#if match.actions?.length}
              <div style="font-size:11px;color:#6b748c;margin-left:20px">
                {#each match.actions as action}
                  <div>
                    <span class="phui-icon-view phui-font-fa {action.type === 'add_reviewer' ? 'fa-user-plus' : action.type === 'add_label' ? 'fa-tag' : action.type === 'post_comment' ? 'fa-comment' : 'fa-bolt'} mrs"></span>
                    {action.type === 'add_reviewer' ? 'Add reviewer' : action.type === 'add_label' ? 'Add label' : action.type === 'post_comment' ? 'Post comment' : action.type}: {action.value}
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/each}
      {:else}
        <div style="font-size:12px;color:#6b748c">No rules matched.</div>
      {/if}
      <div style="margin-top:8px;font-size:12px">
        <a href="/herald"><span class="phui-icon-view phui-font-fa fa-list mrs"></span>Manage Rules</a>
      </div>
    </CurtainBox>

    <!-- Labels -->
    {#if pr.labels?.length}
      <CurtainBox title="Labels">
        {#each pr.labels as label}
          <Tag shade={labelShade(label.color)}>{label.name}</Tag>
          {' '}
        {/each}
      </CurtainBox>
    {/if}

    <!-- Actions -->
    {#if !pr.merged}
      <div class="mood-curtain-actions">
        {#if pr.state !== 'closed'}
          <div style="display:flex;gap:8px;align-items:center;margin-bottom:8px">
            <select bind:value={mergeMethod} style="font-size:12px;padding:4px 6px;border:1px solid #c7ccd9;border-radius:3px;background:#fff;flex:1">
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
