import Link from 'next/link';
import { navItems, site } from '@/lib/site';
import { ThemeToggle } from '@/components/theme-toggle';

export function SiteHeader() {
  return (
    <header className="sticky top-0 z-50 border-b border-[var(--border)] bg-[color:rgba(255,255,255,0.72)] backdrop-blur-xl dark:bg-[color:rgba(7,17,31,0.72)]">
      <div className="container-shell flex items-center justify-between gap-4 py-4">
        <Link href="/" className="flex items-center gap-3 font-semibold tracking-tight">
          <span className="grid h-10 w-10 place-items-center rounded-2xl bg-[var(--accent)] text-base font-bold text-white shadow-glow">
            D
          </span>
          <span className="flex flex-col leading-tight">
            <span>{site.name}</span>
            <span className="text-xs font-normal text-[var(--text-secondary)]">{site.tagline}</span>
          </span>
        </Link>

        <nav className="hidden items-center gap-2 md:flex" aria-label="Primary navigation">
          {navItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className="rounded-full px-4 py-2 text-sm font-medium text-[var(--text-secondary)] transition hover:bg-[var(--accent-soft)] hover:text-[var(--text-primary)]"
            >
              {item.label}
            </Link>
          ))}
        </nav>

        <div className="flex items-center gap-3">
          <ThemeToggle />
          <Link
            href="/contact"
            className="hidden rounded-full bg-[var(--accent)] px-4 py-2 text-sm font-semibold text-white transition hover:-translate-y-0.5 hover:bg-[var(--accent-strong)] md:inline-flex"
          >
            Start a project
          </Link>
        </div>
      </div>
    </header>
  );
}