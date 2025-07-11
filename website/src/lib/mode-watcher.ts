import { writable } from 'svelte/store';

// Get initial mode from localStorage or system preference
function getInitialMode(): 'light' | 'dark' {
  if (typeof localStorage !== 'undefined') {
    const stored = localStorage.getItem('theme-mode');
    if (stored === 'light' || stored === 'dark') return stored;
  }
  if (typeof window !== 'undefined' && window.matchMedia) {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }
  return 'light';
}

export const mode = writable<'light' | 'dark'>(getInitialMode());

mode.subscribe((value) => {
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('theme-mode', value);
  }
  if (typeof document !== 'undefined') {
    document.documentElement.classList.toggle('dark', value === 'dark');
  }
});

export function toggleMode() {
  mode.update((current) => (current === 'light' ? 'dark' : 'light'));
} 