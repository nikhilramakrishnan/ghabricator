import type { PageLoad } from './$types';
import { apiFetch } from '$lib/api';
import type { APIRepoSummary } from '$lib/types';

export const load: PageLoad = async () => {
  const repos = await apiFetch<APIRepoSummary[]>('/api/repos');
  return { repos };
};
