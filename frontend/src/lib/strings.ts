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
    summary: 'Summary',
    reviewers: 'Reviewers',
    buildables: 'Buildables',
    herald: 'Herald',
    labels: 'Labels',
    author: 'Author',
    status: 'Status',
    created: 'Created',
    updated: 'Updated',
    base: 'Base',
    head: 'Head',
    changes: 'Changes',
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
    prs: 'Pull Requests',
    code: 'Code',
    repos: 'Repositories',
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
