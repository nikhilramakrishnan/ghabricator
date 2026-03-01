// Centralized UX strings — no hardcoded labels in components/pages.

export const S = {
  // App
  appName: 'Ghabricator',

  // Navigation / Sidebar
  nav: {
    revisions: 'Revisions',
    repositories: 'Repositories',
    search: 'Search',
    paste: 'Paste',
    actions: 'Actions',
    tools: 'Tools',
  },

  // Dashboard (home)
  dashboard: {
    title: 'Dashboard',
    needsReview: 'Needs Review',
    authored: 'Authored',
    recentBuilds: 'Recent Builds',
    noReviews: 'No reviews requested.',
    noPRs: 'No open pull requests.',
    noBuilds: 'No recent builds.',
    viewAllBuilds: 'View all builds',
  },

  // Revisions page (/dashboard)
  revisions: {
    title: 'Revisions',
    authored: 'Authored',
    needsReview: 'Needs Review',
    noPRs: 'No open pull requests.',
  },

  // PR detail
  pr: {
    revisionContents: 'Revision Contents',
    summary: 'Summary',
    reviewers: 'Reviewers',
    buildables: 'Buildables',
    changeset: 'Changeset',
    herald: 'Herald',
    labels: 'Labels',
    author: 'Author',
    repository: 'Repository',
    status: 'Status',
    created: 'Created',
    updated: 'Updated',
    base: 'Base',
    head: 'Head',
    changes: 'Changes',
    landRevision: 'Land Revision',
    close: 'Close',
    reopen: 'Reopen',
    mergeSquash: 'Squash',
    mergeMerge: 'Merge',
    mergeRebase: 'Rebase',
    mergeFailed: 'Merge failed',
    actionFailed: 'Failed',
    // Status badges
    statusOpen: 'Open',
    statusClosed: 'Closed',
    statusMerged: 'Merged',
    statusDraft: 'Draft',
    // Review states
    reviewAccepted: 'Accepted',
    reviewChangesRequested: 'Changes Requested',
    reviewCommented: 'Commented',
    reviewDismissed: 'Dismissed',
    reviewWaiting: 'Waiting',
    // Interdiff
    showingChanges: 'Showing changes',
    loadingDiff: 'Loading diff...',
  },

  // Diff / Commit history
  diff: {
    title: 'Diff History',
    showChangesFrom: 'Show changes from',
    to: 'to',
    baseSuffix: '(base)',
    latest: 'Latest',
    colSha: 'SHA',
    colAuthor: 'Author',
    colMessage: 'Message',
    colDate: 'Date',
    justNow: 'just now',
    yesterday: 'yesterday',
  },

  // Actions
  actions: {
    title: 'Actions',
    builds: 'Builds',
    rules: 'Rules',
    newRule: 'New Rule',
    noRules: 'No Herald rules configured.',
    noRuns: 'No workflow runs found.',
    conditions: 'Conditions',
    actions: 'Actions',
  },

  // Repos
  repos: {
    title: 'Repositories',
    viewOnGitHub: 'View on GitHub',
    visibility: 'Visibility',
    stars: 'Stars',
    forks: 'Forks',
    details: 'Details',
  },

  // Paste
  paste: {
    title: 'Paste',
    recentPastes: 'Recent Pastes',
    createPaste: 'Create Paste',
    visibility: 'Visibility',
  },

  // Search
  search: {
    title: 'Search',
    prs: 'Pull requests',
    issues: 'Issues',
    code: 'Code',
    repos: 'Repositories',
    filterBy: 'Filter by',
    state: 'State',
    languages: 'Languages',
    open: 'Open',
    closed: 'Closed',
    results: 'results',
    sortBy: 'Sort by:',
    sortBestMatch: 'Best match',
    sortNewest: 'Most recently created',
    sortOldest: 'Least recently created',
    sortMostCommented: 'Most commented',
    sortRecentlyUpdated: 'Most recently updated',
    noResults: 'No results found.',
    matches: 'matches',
  },

  // Breadcrumbs
  crumb: {
    home: 'Home',
  },

  // Landing page
  landing: {
    title: 'Welcome to Ghabricator',
    desc: 'A modern code review experience powered by GitHub. Review pull requests with Phabricator\'s legendary diff viewer, inline comments, and Herald automation.',
    signIn: 'Sign in with GitHub',
    sideDiffs: 'Side-by-side diffs',
    inlineComments: 'Inline comments',
    heraldRules: 'Herald rules',
    pasteBin: 'Paste bin',
  },

  // Common
  common: {
    details: 'Details',
    actions: 'Actions',
    author: 'Author',
    created: 'Created',
    updated: 'Updated',
    public: 'Public',
    secret: 'Secret',
    active: 'Active',
    disabled: 'Disabled',
    draft: 'Draft',
    remove: 'Remove',
    editRule: 'Edit Rule',
  },
} as const;
