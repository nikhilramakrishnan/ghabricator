<script lang="ts">
  import { Breadcrumbs } from '$lib/components/layout';
  import { Box, HeaderView, Tag } from '$lib/components/phui';
  import { DiffTable, ChangesetHeader, CommitHistory, FileTree, InlineCommentWithContext, ReactionPicker } from '$lib/components/diff';
  import type { APIReviewComment as DiffComment, APIDiffRow } from '$lib/components/diff';
  import { ReviewForm } from '$lib/components/review';
  import { apiFetch, apiPost } from '$lib/api';
  import { S } from '$lib/strings';
  import { formatTimestamp } from '$lib/time';
  import { addDraft, addReplyDraft } from '$lib/stores/inline';
  import { fileTreeData } from '$lib/stores/filetree';
  import { user } from '$lib/stores/auth';
  import { MarkdownEditor } from '$lib/components/editor';
  import { marked } from 'marked';
  import { onDestroy } from 'svelte';
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
  let changesetCollapsed = $state(false);
  let commentsCollapsed = $state(false);

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

  // Feed file tree data to sidebar store
  $effect(() => {
    fileTreeData.set({
      changesets: displayChangesets,
      activeFile: activePath,
      commentCounts: Object.fromEntries(Object.entries(commentsByPath).map(([k, v]) => [k, v.length]))
    });
  });
  onDestroy(() => fileTreeData.set(null));

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
      createdAt: c.createdAt,
      inReplyTo: c.inReplyTo,
      reactions: c.reactions
    }));
  }

  // Check run icon/color mapping
  function checkRunDisplay(cr: APICheckRun): { icon: string; color: string; name: string; duration: string } {
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
    // Duration
    let duration = '';
    if (cr.startedAt && cr.completedAt) {
      const start = new Date(cr.startedAt).getTime();
      const end = new Date(cr.completedAt).getTime();
      const secs = Math.floor((end - start) / 1000);
      if (secs >= 60) {
        const m = Math.floor(secs / 60);
        const s = secs % 60;
        duration = s > 0 ? `${m}m ${s}s` : `${m}m`;
      } else if (secs > 0) {
        duration = `${secs}s`;
      }
    }
    return { icon, color, name, duration };
  }

  // Is the PR approved? Check if any reviewer's latest review is APPROVED
  let isApproved = $derived.by(() => {
    const latest = new Map<string, { state: string; at: string }>();
    for (const r of reviews) {
      const prev = latest.get(r.author.login);
      if (!prev || r.createdAt > prev.at) {
        latest.set(r.author.login, { state: r.state, at: r.createdAt });
      }
    }
    for (const { state } of latest.values()) {
      if (state === 'APPROVED') return true;
    }
    return false;
  });

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
      iconColor: ev.iconColor,
      commentID: ev.commentID,
      commentType: ev.commentType,
      reactions: ev.reactions
    }))
  );

  // Unified comment stream: inline comments with code context + timeline events
  type StreamItem =
    | { type: 'inline'; comment: APIReviewComment; contextRows: APIDiffRow[]; createdAt: string }
    | { type: 'timeline'; event: TimelineEvent; createdAt: string };

  function getContextRows(comment: APIReviewComment): APIDiffRow[] {
    const cs = displayChangesets.find(c => c.displayPath === comment.path);
    if (!cs) return [];
    const idx = cs.rows.findIndex(r =>
      comment.side === 'RIGHT' ? r.newNum === comment.line : r.oldNum === comment.line
    );
    if (idx < 0) return [];
    return cs.rows.slice(Math.max(0, idx - 2), idx + 3);
  }

  let _commentStreamData = $derived.by(() => {
    const items: StreamItem[] = [];

    // Build id lookup for all comments across all paths
    const allComments: APIReviewComment[] = [];
    for (const comments of Object.values(commentsByPath)) {
      for (const c of comments) allComments.push(c);
    }
    const byId = new Map<number, APIReviewComment>();
    for (const c of allComments) byId.set(c.id, c);

    // Identify reply IDs so we can filter them from top-level
    const replyIds = new Set<number>();
    for (const c of allComments) {
      if (c.inReplyTo && byId.has(c.inReplyTo)) {
        replyIds.add(c.id);
      }
    }

    // Group replies by root
    const repliesByRoot = new Map<number, APIReviewComment[]>();
    for (const c of allComments) {
      if (!replyIds.has(c.id)) continue;
      // Walk up to find root
      let rootId = c.inReplyTo!;
      let parent = byId.get(rootId);
      while (parent && parent.inReplyTo && byId.has(parent.inReplyTo)) {
        rootId = parent.inReplyTo;
        parent = byId.get(rootId);
      }
      if (!repliesByRoot.has(rootId)) repliesByRoot.set(rootId, []);
      repliesByRoot.get(rootId)!.push(c);
    }
    // Sort replies chronologically
    for (const replies of repliesByRoot.values()) {
      replies.sort((a, b) => a.createdAt.localeCompare(b.createdAt));
    }

    // Only add root comments (non-replies) to the stream
    for (const c of allComments) {
      if (replyIds.has(c.id)) continue;
      items.push({
        type: 'inline',
        comment: c,
        contextRows: getContextRows(c),
        createdAt: c.createdAt
      });
    }

    // Add timeline events
    for (const ev of timelineEvents) {
      items.push({ type: 'timeline', event: ev, createdAt: ev.createdAt });
    }

    // Sort chronologically
    items.sort((a, b) => a.createdAt.localeCompare(b.createdAt));

    return { items, repliesByRoot };
  });

  let commentStream: StreamItem[] = $derived(_commentStreamData.items);
  let repliesByRoot = $derived(_commentStreamData.repliesByRoot);

  function navigateToInline(path: string, line: number) {
    // Expand changeset box if collapsed
    if (changesetCollapsed) changesetCollapsed = false;

    // Find and expand the file
    const cs = displayChangesets.find(c => c.displayPath === path);
    if (cs && collapsedFiles.has(cs.id)) {
      collapsedFiles.delete(cs.id);
      collapsedFiles = new Set(collapsedFiles);
    }

    // Wait for DOM to update, then scroll to the exact line
    requestAnimationFrame(() => {
      requestAnimationFrame(() => {
        if (!cs) return;
        const container = document.getElementById(`C${cs.id}`);
        if (!container) return;

        // Find the td with matching data-n inside this changeset
        const cell = container.querySelector(`td[data-n="${line}"]`);
        const row = cell?.closest('tr');
        if (row) {
          row.scrollIntoView({ behavior: 'smooth', block: 'center' });
          row.classList.add('flash-highlight');
          setTimeout(() => row.classList.remove('flash-highlight'), 2000);
        } else {
          container.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
      });
    });
  }

  // Merge/close actions
  function handleNewComment(path: string, line: number, side: string) {
    addDraft(path, line, side);
  }

  const EMOJI_ICONS: Record<string, string> = {
    '+1': 'fa-thumbs-up',
    '-1': 'fa-thumbs-down',
    'laugh': 'fa-smile-o',
    'confused': 'fa-question',
    'heart': 'fa-heart',
    'star': 'fa-star',
    'rocket': 'fa-rocket',
    'eyes': 'fa-eye'
  };

  let reactionPickerOpen = $state<number | null>(null);
  let localReactions = $state<Map<number, import('$lib/types').APIReaction[]>>(new Map());

  function getReactions(ev: TimelineEvent): import('$lib/types').APIReaction[] {
    if (ev.commentID && localReactions.has(ev.commentID)) {
      return localReactions.get(ev.commentID)!;
    }
    return ev.reactions ?? [];
  }

  async function handleTimelineReaction(ev: TimelineEvent, emoji: string) {
    if (!ev.commentID || !ev.commentType) return;
    try {
      await apiPost('/api/v2/reaction', {
        owner,
        repo,
        commentID: ev.commentID,
        content: emoji,
        commentType: ev.commentType
      });
      // Optimistic update in local state
      const current = getReactions(ev).map(r => ({ ...r }));
      const existing = current.find(r => r.emoji === emoji);
      if (existing) {
        existing.count++;
      } else {
        current.push({ emoji, count: 1 });
      }
      localReactions.set(ev.commentID, current);
      localReactions = new Map(localReactions);
    } catch {
      // silent fail
    }
    reactionPickerOpen = null;
  }

  // --- Edit state ---
  let editingSummary = $state(false);
  let summaryDraft = $state('');
  let savingSummary = $state(false);
  let editingCommentId = $state<number | null>(null);
  let commentDraft = $state('');
  let savingComment = $state(false);

  function startEditSummary() {
    summaryDraft = pr.bodyRaw ?? '';
    editingSummary = true;
  }

  async function saveSummary() {
    savingSummary = true;
    try {
      await apiPost('/api/v2/edit-pr', { owner, repo, number, body: summaryDraft });
      pr.body = marked.parse(summaryDraft, { async: false }) as string;
      pr.bodyRaw = summaryDraft;
      editingSummary = false;
    } catch { /* silent */ }
    savingSummary = false;
  }

  function startEditComment(ev: TimelineEvent) {
    commentDraft = ev.bodyRaw ?? '';
    editingCommentId = ev.commentID ?? null;
  }

  async function saveComment(ev: TimelineEvent) {
    if (!ev.commentID) return;
    savingComment = true;
    try {
      await apiPost('/api/v2/edit-comment', { owner, repo, commentID: ev.commentID, body: commentDraft });
      ev.body = marked.parse(commentDraft, { async: false }) as string;
      ev.bodyRaw = commentDraft;
      editingCommentId = null;
    } catch { /* silent */ }
    savingComment = false;
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

<div class="pr-content">
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

    {#if editingSummary}
      <div class="summary-section">
        <MarkdownEditor bind:value={summaryDraft} minRows={6} autofocus />
        <div class="edit-actions">
          <button class="edit-btn save" onclick={saveSummary} disabled={savingSummary}>
            <i class="fa fa-check"></i> {savingSummary ? 'Saving...' : 'Save'}
          </button>
          <button class="edit-btn cancel" onclick={() => editingSummary = false}>
            <i class="fa fa-times"></i> Cancel
          </button>
        </div>
      </div>
    {:else}
      <div class="summary-section">
        {#if pr.body?.trim()}
          <div class="remarkup-content">
            {@html pr.body}
          </div>
        {/if}
        <div class="summary-edit-row">
          <button class="summary-edit-link" onclick={startEditSummary}>
            <i class="fa fa-pencil"></i> Edit Summary
          </button>
        </div>
      </div>
    {/if}

  </Box>

  {#if checkRuns.length > 0}
    <Box border>
      <HeaderView title={S.pr.buildables} icon="fa-cogs" />
      <div class="buildables-list">
        {#each checkRuns as cr}
          {@const d = checkRunDisplay(cr)}
          {#if cr.detailsURL}
            <a href={cr.detailsURL} target="_blank" rel="noopener" class="buildable-item">
              <i class="fa {d.icon}" style="color:{d.color}"></i>
              <span class="buildable-name">{d.name}</span>
              {#if d.duration}<span class="buildable-duration"><i class="fa fa-clock-o"></i> {d.duration}</span>{/if}
            </a>
          {:else}
            <div class="buildable-item">
              <i class="fa {d.icon}" style="color:{d.color}"></i>
              <span class="buildable-name">{d.name}</span>
              {#if d.duration}<span class="buildable-duration"><i class="fa fa-clock-o"></i> {d.duration}</span>{/if}
            </div>
          {/if}
        {/each}
      </div>
    </Box>
  {/if}

  {#if commentStream.length > 0}
    <Box border>
      <HeaderView title="Comments" icon="fa-comments" count={commentStream.length} collapsible collapsed={commentsCollapsed} onToggle={() => commentsCollapsed = !commentsCollapsed} />
      {#if !commentsCollapsed}
        <div class="comment-stream">
          {#each commentStream as item, idx}
            {#if item.type === 'inline'}
              <InlineCommentWithContext
                comment={item.comment}
                contextRows={item.contextRows}
                side={item.comment.side}
                replies={repliesByRoot.get(item.comment.id) ?? []}
                onNavigate={navigateToInline}
                onReply={(c) => addReplyDraft(c.path, c.line, c.side, c.id)}
              />
            {:else if item.event.body}
              {@const reactions = getReactions(item.event)}
              <div class="stream-card">
                <div class="stream-card-meta">
                  {#if item.event.author.avatarURL}
                    <img src={item.event.author.avatarURL} alt="" class="stream-card-avatar" />
                  {/if}
                  <strong>{item.event.author.login}</strong>
                  <span class="stream-card-action">{item.event.action}</span>
                  <span class="stream-card-time">{formatTimestamp(item.event.createdAt)}</span>
                </div>
                {#if editingCommentId === item.event.commentID}
                  <div class="stream-card-edit">
                    <MarkdownEditor bind:value={commentDraft} minRows={4} autofocus />
                    <div class="edit-actions">
                      <button class="edit-btn save" onclick={() => saveComment(item.event)} disabled={savingComment}>
                        <i class="fa fa-check"></i> {savingComment ? 'Saving...' : 'Save'}
                      </button>
                      <button class="edit-btn cancel" onclick={() => editingCommentId = null}>
                        <i class="fa fa-times"></i> Cancel
                      </button>
                    </div>
                  </div>
                {:else}
                  <div class="stream-card-body">
                    {@html item.event.body}
                  </div>
                  {#if reactions.length > 0}
                    <div class="stream-card-reactions">
                      {#each reactions as r}
                        <button
                          class="stream-card-pill"
                          onclick={() => handleTimelineReaction(item.event, r.emoji)}
                          title={r.emoji}
                        >
                          <i class="fa {EMOJI_ICONS[r.emoji] ?? 'fa-smile-o'}"></i>
                          <span class="pill-count">{r.count}</span>
                        </button>
                      {/each}
                    </div>
                  {/if}
                {/if}
                <div class="stream-card-actions">
                  <button class="stream-action-btn" onclick={() => document.querySelector('.review-form-anchor')?.scrollIntoView({ behavior: 'smooth' })}>
                    <i class="fa fa-reply"></i> Reply
                  </button>
                  {#if item.event.commentID && item.event.bodyRaw !== undefined}
                    <button class="stream-action-btn" onclick={() => startEditComment(item.event)}>
                      <i class="fa fa-pencil"></i> Edit
                    </button>
                  {/if}
                  {#if item.event.commentID}
                    <div class="picker-anchor">
                      <button class="stream-action-btn" title="Add reaction" onclick={() => reactionPickerOpen = reactionPickerOpen === item.event.commentID ? null : (item.event.commentID ?? null)}>
                        <i class="fa fa-smile-o"></i>
                      </button>
                      {#if reactionPickerOpen === item.event.commentID}
                        <ReactionPicker
                          onPick={(emoji) => handleTimelineReaction(item.event, emoji)}
                          onClose={() => reactionPickerOpen = null}
                        />
                      {/if}
                    </div>
                  {/if}
                </div>
              </div>
            {:else}
              <div class="stream-status">
                <div class="stream-icon" style="background:{item.event.iconColor === 'green' ? 'var(--green)' : item.event.iconColor === 'red' ? 'var(--red)' : item.event.iconColor === 'blue' ? 'var(--blue)' : item.event.iconColor === 'violet' ? 'var(--violet)' : 'var(--text-muted)'}">
                  <i class="fa {item.event.iconClass}"></i>
                </div>
                <div class="stream-status-body">
                  <span><strong>{item.event.author.login}</strong> {item.event.action}</span>
                  <span class="stream-status-time">{formatTimestamp(item.event.createdAt)}</span>
                </div>
              </div>
            {/if}
          {/each}
        </div>
      {/if}
    </Box>
  {/if}

  <Box border>
    <HeaderView title={S.pr.changeset} icon="fa-files-o" count={displayChangesets.length} collapsible collapsed={changesetCollapsed} onToggle={() => changesetCollapsed = !changesetCollapsed} />
    {#if !changesetCollapsed}
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
    {/if}
  </Box>

  {#if commits.length > 0}
    <CommitHistory
      {commits}
      baseBranch={pr.base.ref}
      onRangeChange={handleRangeChange}
    />
  {/if}

  <div class="review-form-anchor"></div>
  <ReviewForm {owner} {repo} {number} merged={pr.merged} prState={pr.state} authorLogin={pr.author.login} approved={isApproved} />
</div>

<style>
  .pr-content {
    padding: 0 16px;
  }

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

  /* Summary body — below property list, separated by border */
  .summary-section {
    padding: 12px 16px;
    border-top: 1px solid var(--border-subtle);
  }

  /* Buildables */
  .buildables-list {
    padding: 4px 16px 8px;
  }

  .buildable-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 5px 0;
    font-size: 12px;
    color: var(--text-muted);
    text-decoration: none;
    border-bottom: 1px solid var(--border-subtle);
  }
  .buildable-item:last-child {
    border-bottom: none;
  }
  .buildable-name {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .buildable-duration {
    font-size: 11px;
    color: var(--text-muted);
    font-family: var(--font-mono);
  }
  a.buildable-item:hover {
    color: var(--text);
  }

  /* Comment stream */
  .comment-stream {
    padding: 8px;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  /* Card-style timeline events (with body) */
  .stream-card {
    border: 1px solid var(--border);
    border-radius: 4px;
    overflow: hidden;
    background: var(--bg-card);
  }

  .stream-card-meta {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 12px;
    font-size: 12px;
    background: var(--bg-card-header);
    border-bottom: 1px solid var(--border-subtle);
  }
  .stream-card-meta strong {
    color: var(--text);
  }
  .stream-card-avatar {
    width: 20px;
    height: 20px;
    border-radius: 3px;
  }
  .stream-card-action {
    color: var(--text-muted);
  }
  .stream-card-time {
    color: var(--text-muted);
    margin-left: auto;
    white-space: nowrap;
  }

  .stream-card-body {
    padding: 4px 12px;
    font-size: 13px;
    color: var(--text);
    line-height: 1.5;
    overflow-wrap: break-word;
  }

  .stream-card-reactions {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    padding: 0 12px 8px;
  }

  .stream-card-pill {
    all: unset;
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 2px 8px;
    border-radius: 10px;
    font-size: 11px;
    background: var(--bg-subtle);
    border: 1px solid var(--border-subtle);
    color: var(--text-muted);
    cursor: pointer;
    transition: background 0.1s;
  }
  .stream-card-pill:hover {
    background: var(--bg-hover);
    color: var(--text);
  }
  .pill-count {
    font-weight: 600;
  }

  .stream-card-actions {
    padding: 6px 12px;
    border-top: 1px solid var(--border-subtle);
    display: flex;
    gap: 12px;
    align-items: center;
  }

  .picker-anchor {
    position: relative;
  }

  .stream-action-btn {
    all: unset;
    font-size: 11px;
    color: var(--text-muted);
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 4px;
  }
  .stream-action-btn:hover {
    color: var(--text-link);
  }

  .stream-card-edit {
    padding: 8px 12px;
  }

  .edit-actions {
    display: flex;
    gap: 8px;
    margin-top: 8px;
  }

  .edit-btn {
    all: unset;
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 4px 12px;
    font-size: 12px;
    border-radius: 3px;
    cursor: pointer;
    font-weight: 500;
  }
  .edit-btn.save {
    background: var(--blue);
    color: #fff;
  }
  .edit-btn.save:hover { opacity: 0.9; }
  .edit-btn.save:disabled { opacity: 0.5; cursor: default; }
  .edit-btn.cancel {
    color: var(--text-muted);
  }
  .edit-btn.cancel:hover { color: var(--text); }

  .summary-edit-row {
    padding: 4px 12px 8px;
  }
  .summary-edit-link {
    all: unset;
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 12px;
    color: var(--text-muted);
    cursor: pointer;
  }
  .summary-edit-link:hover {
    color: var(--text-link);
  }

  /* Compact status events (no body) */
  .stream-status {
    display: flex;
    gap: 10px;
    align-items: center;
    padding: 4px 0;
  }

  .stream-icon {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    color: #fff;
    font-size: 11px;
  }

  .stream-status-body {
    flex: 1;
    font-size: 13px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .stream-status-time {
    font-size: 12px;
    color: var(--text-muted);
    margin-left: auto;
    white-space: nowrap;
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
