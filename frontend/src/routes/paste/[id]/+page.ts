import type { PageLoad } from './$types';
import { apiFetch } from '$lib/api';
import type { APIPasteDetail } from '$lib/types';

export const load: PageLoad = async ({ params }) => {
  const paste = await apiFetch<APIPasteDetail>(`/api/paste/${params.id}`);
  return { paste };
};
