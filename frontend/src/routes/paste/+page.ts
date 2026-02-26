import type { PageLoad } from './$types';
import { apiFetch } from '$lib/api';
import type { APIPasteSummary } from '$lib/types';

export const load: PageLoad = async () => {
  const pastes = await apiFetch<APIPasteSummary[]>('/api/paste');
  return { pastes };
};
