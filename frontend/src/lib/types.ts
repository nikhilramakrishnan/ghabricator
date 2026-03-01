// Shared API response types — mirrors internal/server/api_types.go

export interface APIUser {
  login: string;
  avatarURL: string;
}

export interface APILabel {
  name: string;
  color: string;
}

// --- Dashboard ---

export interface APIPRSummary {
  number: number;
  title: string;
  owner: string;
  repo: string;
  author: APIUser;
  draft: boolean;
  labels?: APILabel[];
  reviewers?: APIUser[];
  updatedAt: string;
}

export interface APIDashboardResponse {
  authored: APIPRSummary[];
  reviewRequested: APIPRSummary[];
}

// --- PR Detail ---

export interface APIRef {
  ref: string;
  sha: string;
  repo: string;
}

export interface APIPRDetail {
  number: number;
  title: string;
  body: string;
  bodyRaw?: string;
  state: string;
  draft: boolean;
  merged: boolean;
  author: APIUser;
  createdAt: string;
  updatedAt: string;
  labels: APILabel[];
  reviewers: APIUser[];
  head: APIRef;
  base: APIRef;
  additions: number;
  deletions: number;
  changedFiles: number;
}

export interface APIDiffRow {
  oldNum: number;
  newNum: number;
  oldClass: string;
  newClass: string;
  oldContent: string;
  newContent: string;
  isContext: boolean;
}

export interface APIChangeset {
  id: number;
  oldName: string;
  newName: string;
  displayPath: string;
  linesAdded: number;
  linesRemoved: number;
  isNew: boolean;
  isDeleted: boolean;
  isRenamed: boolean;
  isBinary: boolean;
  rows: APIDiffRow[];
}

export interface APIReaction {
  emoji: string;
  count: number;
}

export interface APIReviewComment {
  id: number;
  author: APIUser;
  body: string;
  path: string;
  line: number;
  side: string;
  createdAt: string;
  inReplyTo?: number;
  reactions?: APIReaction[];
}

export interface APIReview {
  id: number;
  author: APIUser;
  state: string;
  body: string;
  createdAt: string;
}

export interface APIIssueComment {
  id: number;
  author: APIUser;
  body: string;
  createdAt: string;
}

export interface APICheckRun {
  name: string;
  status: string;
  conclusion: string;
  detailsURL: string;
  appName: string;
  startedAt: string;
  completedAt: string;
}

export interface APITimelineEvent {
  author: APIUser;
  action: string;
  body?: string;
  bodyRaw?: string;
  createdAt: string;
  iconClass: string;
  iconColor: string;
  commentID?: number;
  commentType?: string;
  reactions?: APIReaction[];
}

export interface APIHeraldAction {
  type: string;
  value: string;
}

export interface APIHeraldMatch {
  ruleId: string;
  ruleName: string;
  actions: APIHeraldAction[];
}

export interface APICommit {
  sha: string;
  message: string;
  author: APIUser;
  date: string;
}

export interface APIPRDetailResponse {
  pr: APIPRDetail;
  changesets: APIChangeset[];
  commentsByPath: Record<string, APIReviewComment[]>;
  reviews: APIReview[];
  issueComments: APIIssueComment[];
  checkRuns: APICheckRun[];
  timeline: APITimelineEvent[];
  heraldMatches?: APIHeraldMatch[];
  commits: APICommit[];
  viewerPermission: string; // ADMIN, MAINTAIN, WRITE, TRIAGE, READ, or ""
}

// --- Repos ---

export interface APIRepoSummary {
  name: string;
  fullName: string;
  description: string;
  language: string;
  stars: number;
  forks: number;
  private: boolean;
  fork: boolean;
  archived: boolean;
  avatarURL: string;
  updatedAt: string;
}

export interface APIRepoInfo {
  fullName: string;
  description: string;
  defaultBranch: string;
  private: boolean;
  htmlURL: string;
  stars: number;
  forks: number;
}

export interface APIRepoEntry {
  name: string;
  path: string;
  type: string; // "file" | "dir"
  size: number;
}

export interface APIRepoTreeResponse {
  entries: APIRepoEntry[];
  repoInfo: APIRepoInfo;
}

export interface APIRepoFile {
  name: string;
  path: string;
  size: number;
  lines?: string[];
  rawURL?: string;
  htmlURL: string;
}

export interface APIRepoFileResponse {
  file: APIRepoFile;
  repoInfo: APIRepoInfo;
}

// --- Blame ---

export interface APIBlameRange {
  startLine: number;
  endLine: number;
  commitOID: string;
  commitShort: string;
  message: string;
  authorLogin: string;
  authorAvatarURL: string;
  authorName: string;
  authoredDate: string;
}

export interface APIBlameResponse {
  ranges: APIBlameRange[];
}

// --- Paste ---

export interface APIPasteSummary {
  id: string;
  title: string;
  language: string;
  public: boolean;
  createdAt: string;
}

export interface APIPasteFile {
  filename: string;
  language: string;
  size: number;
  lines: string[];
}

export interface APIPasteDetail {
  id: string;
  title: string;
  public: boolean;
  owner: APIUser;
  htmlURL: string;
  files: APIPasteFile[];
  createdAt: string;
  updatedAt: string;
}

export interface APIPasteCreateResponse {
  id: string;
  url: string;
}

// --- Herald ---

export interface HeraldCondition {
  type: string;
  value: string;
}

export interface HeraldAction {
  type: string;
  value: string;
}

export interface HeraldRule {
  id: string;
  name: string;
  author_login: string;
  conditions: HeraldCondition[];
  actions: HeraldAction[];
  must_match_all: boolean;
  disabled: boolean;
  created_at: string;
  updated_at: string;
}

// --- Actions / Workflow Runs ---

export interface APIWorkflowRun {
  id: number;
  name: string;
  displayTitle: string;
  status: string;
  conclusion: string;
  branch: string;
  event: string;
  actor: APIUser;
  repoOwner: string;
  repoName: string;
  durationMs: number;
  htmlURL: string;
  createdAt: string;
}

export interface APIWorkflowRunsResponse {
  runs: APIWorkflowRun[];
}

// --- Search ---

export interface APISearchPR {
  number: number;
  title: string;
  repo: string;
  state: string;
  author: string;
  avatarURL: string;
  labels: APILabel[];
  body: string;
  commentsCount: number;
  createdAt: string;
  updatedAt: string;
  draft: boolean;
  url: string;
}

export interface APISearchIssue {
  number: number;
  title: string;
  repo: string;
  state: string;
  author: string;
  avatarURL: string;
  labels: APILabel[];
  body: string;
  commentsCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface APISearchCodeResult {
  repo: string;
  path: string;
  fragment: string;
  language: string;
  matchCount: number;
  htmlURL: string;
}

export interface APISearchRepoResult {
  fullName: string;
  description: string;
  stars: number;
  forks: number;
  language: string;
  updatedAt: string;
  topics: string[];
  avatarURL: string;
}

export interface APISearchResponse {
  counts?: Record<string, number>;
  prs?: APISearchPR[];
  issues?: APISearchIssue[];
  code?: APISearchCodeResult[];
  repos?: APISearchRepoResult[];
}
