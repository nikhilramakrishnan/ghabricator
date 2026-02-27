import type { PageLoad } from './$types';
import { apiFetch } from '$lib/api';
import type { HeraldRule, APIWorkflowRunsResponse } from '$lib/types';

export const load: PageLoad = async () => {
  const [rules, runsResp] = await Promise.all([
    apiFetch<HeraldRule[]>('/api/herald'),
    apiFetch<APIWorkflowRunsResponse>('/api/actions/runs'),
  ]);
  return { rules, runs: runsResp.runs ?? [] };
};
