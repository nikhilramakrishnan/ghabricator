import { writable } from 'svelte/store';
import type { APIChangeset } from '$lib/components/diff/DiffTable.svelte';

interface FileTreeState {
  changesets: APIChangeset[];
  activeFile: string;
  commentCounts: Record<string, number>;
}

// Set by the PR page when it has changeset data, cleared on navigation away
export const fileTreeData = writable<FileTreeState | null>(null);

// Toggled by the sidebar
export const fileTreeOpen = writable(false);
