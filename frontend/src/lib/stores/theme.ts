import { writable } from 'svelte/store';

function getInitialTheme(): string {
  if (typeof document === 'undefined') return '';
  const match = document.cookie.match(/(?:^|;\s*)theme=(\w*)/);
  return match?.[1] === 'dark' ? 'dark' : '';
}

export const theme = writable(getInitialTheme());

export function toggleTheme() {
  theme.update((t) => {
    const next = t === 'dark' ? '' : 'dark';
    document.cookie = `theme=${next};path=/;max-age=31536000;SameSite=Lax`;
    return next;
  });
}
