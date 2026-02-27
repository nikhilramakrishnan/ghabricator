import type { PageLoad } from './$types';
import { apiFetch } from '$lib/api';
import type { AuthUser } from '$lib/stores/auth';
import type { APIDashboardResponse, APIWorkflowRunsResponse } from '$lib/types';

export const load: PageLoad = async () => {
  let loggedIn = false;
  try {
    await apiFetch<AuthUser>('/api/auth/me');
    loggedIn = true;
  } catch {
    // not logged in — show landing page
  }
  if (!loggedIn) return { loggedIn: false };

  const [dashboard, actions] = await Promise.all([
    apiFetch<APIDashboardResponse>('/api/dashboard'),
    apiFetch<APIWorkflowRunsResponse>('/api/actions/runs').catch(() => ({ runs: [] })),
  ]);

  return {
    loggedIn: true,
    authored: dashboard.authored ?? [],
    reviewRequested: dashboard.reviewRequested ?? [],
    recentRuns: (actions.runs ?? []).slice(0, 5),
  };
};
