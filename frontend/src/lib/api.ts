export async function apiFetch<T>(path: string, opts?: RequestInit & { noRedirect?: boolean }): Promise<T> {
  const { noRedirect, ...fetchOpts } = opts ?? {};
  const res = await fetch(path, { credentials: 'include', ...fetchOpts });
  if (res.status === 401) {
    if (!noRedirect) {
      window.location.href = '/api/auth/github';
    }
    throw new Error('Unauthorized');
  }
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error || res.statusText);
  }
  return res.json();
}

export async function apiPost<T>(path: string, data: unknown): Promise<T> {
  return apiFetch<T>(path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
}
