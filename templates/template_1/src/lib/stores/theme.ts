import { browser } from '$app/environment';
import { writable } from 'svelte/store';

// Determine initial theme: localStorage, system, or default to 'light'
function getInitialTheme(): 'light' | 'dark' {
  if (!browser) return 'light';
  const stored = localStorage.getItem('theme');
  if (stored === 'light' || stored === 'dark') {
    return stored as 'light' | 'dark';
  }
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

// Theme store
export const theme = writable<'light' | 'dark'>(getInitialTheme());

// Subscribe to store changes to update document class and localStorage
theme.subscribe((value) => {
  if (!browser) return;
  if (value === 'dark') {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
  localStorage.setItem('theme', value);
});

// Toggle between light and dark
export function toggleTheme() {
  theme.update((current) => (current === 'light' ? 'dark' : 'light'));
} 