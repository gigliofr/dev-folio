"use client";

import { useTheme } from './theme-provider';

export function ThemeToggle() {
  const { theme, toggleTheme } = useTheme();

  return (
    <button
      type="button"
      onClick={toggleTheme}
      className="inline-flex items-center gap-2 rounded-full border border-[var(--border)] bg-[var(--surface-strong)] px-4 py-2 text-sm font-medium text-[var(--text-primary)] transition hover:-translate-y-0.5 hover:shadow-glow"
      aria-label="Toggle color theme"
    >
      <span className="h-2.5 w-2.5 rounded-full bg-[var(--accent)]" aria-hidden="true" />
      {theme === 'light' ? 'Dark mode' : 'Light mode'}
    </button>
  );
}