import { writable, derived } from 'svelte/store';

export interface DraftComment {
  path: string;
  line: number;
  side: string;
  body: string;
  inReplyTo?: number;
}

export const drafts = writable<DraftComment[]>([]);

export const pendingCount = derived(drafts, ($drafts) =>
  $drafts.filter((d) => d.body.trim().length > 0).length
);

export function addDraft(path: string, line: number, side: string) {
  drafts.update((d) => [...d, { path, line, side, body: '' }]);
}

export function addReplyDraft(path: string, line: number, side: string, inReplyTo: number) {
  drafts.update((d) => [...d, { path, line, side, body: '', inReplyTo }]);
}

export function removeDraft(path: string, line: number, side: string) {
  drafts.update((d) =>
    d.filter((x) => !(x.path === path && x.line === line && x.side === side))
  );
}

export function updateDraft(path: string, line: number, side: string, body: string) {
  drafts.update((d) =>
    d.map((x) =>
      x.path === path && x.line === line && x.side === side ? { ...x, body } : x
    )
  );
}

export function clearDrafts() {
  drafts.set([]);
}
