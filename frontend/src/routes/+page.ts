import type { PageLoad } from './$types';
import { apiFetch } from '$lib/api';
import { redirect } from '@sveltejs/kit';
import type { AuthUser } from '$lib/stores/auth';

export const load: PageLoad = async () => {
  try {
    await apiFetch<AuthUser>('/api/auth/me');
    redirect(302, '/dashboard');
  } catch {
    // not logged in — show landing page
  }
};
