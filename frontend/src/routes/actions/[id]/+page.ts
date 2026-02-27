import type { PageLoad } from './$types';
import { apiFetch } from '$lib/api';
import type { HeraldRule } from '$lib/types';

export const load: PageLoad = async ({ params }) => {
  const rule = await apiFetch<HeraldRule>(`/api/herald/${params.id}`);
  return { rule };
};
