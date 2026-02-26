import type { PageLoad } from './$types';
import { apiFetch } from '$lib/api';
import type { APIPRDetailResponse } from '$lib/types';

export const load: PageLoad = async ({ params }) => {
  const { owner, repo, number } = params;
  const data = await apiFetch<APIPRDetailResponse>(`/api/pr/${owner}/${repo}/${number}`);
  return { owner, repo, number: Number(number), data };
};
