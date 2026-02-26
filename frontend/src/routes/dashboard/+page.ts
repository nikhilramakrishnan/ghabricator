import type { PageLoad } from './$types';
import { apiFetch } from '$lib/api';
import type { APIDashboardResponse } from '$lib/types';

export const load: PageLoad = async () => {
  const data = await apiFetch<APIDashboardResponse>('/api/dashboard');
  return { data };
};
