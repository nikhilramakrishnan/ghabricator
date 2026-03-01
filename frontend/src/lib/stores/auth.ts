import { writable } from 'svelte/store';
import { apiFetch } from '$lib/api';

export interface AuthUser {
  login: string;
  avatarURL: string;
}

export const user = writable<AuthUser | null>(null);
export const authLoading = writable(true);

export async function checkAuth() {
  try {
    const data = await apiFetch<AuthUser>('/api/auth/me', { noRedirect: true });
    user.set(data);
  } catch {
    user.set(null);
  } finally {
    authLoading.set(false);
  }
}
