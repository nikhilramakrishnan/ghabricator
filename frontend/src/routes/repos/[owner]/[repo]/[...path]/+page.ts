import type { PageLoad } from './$types';
import { apiFetch } from '$lib/api';
import type { APIRepoTreeResponse, APIRepoFileResponse } from '$lib/types';

export const load: PageLoad = async ({ params, url }) => {
  const { owner, repo } = params;
  const path = params.path || '';
  const ref = url.searchParams.get('ref') || '';

  const qs = new URLSearchParams();
  if (ref) qs.set('ref', ref);
  if (path) qs.set('path', path);
  const query = qs.toString() ? `?${qs.toString()}` : '';

  // Try tree first, then file
  try {
    const tree = await apiFetch<APIRepoTreeResponse>(`/api/repo/${owner}/${repo}/tree${query}`);
    return { owner, repo, path, ref, mode: 'tree' as const, tree, file: null };
  } catch {
    try {
      const file = await apiFetch<APIRepoFileResponse>(`/api/repo/${owner}/${repo}/file${query}`);
      return { owner, repo, path, ref, mode: 'file' as const, tree: null, file };
    } catch (e) {
      return { owner, repo, path, ref, mode: 'error' as const, tree: null, file: null, error: String(e) };
    }
  }
};
