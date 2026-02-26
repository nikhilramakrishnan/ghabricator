import type { PageLoad } from './$types';
import { apiFetch } from '$lib/api';
import type { HeraldRule } from '$lib/types';

export const load: PageLoad = async () => {
  const rules = await apiFetch<HeraldRule[]>('/api/herald');
  return { rules };
};
