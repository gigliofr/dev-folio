import Link from 'next/link';
import { navItems, socialLinks, site } from '@/lib/site';

export function SiteFooter() {
  return (
    <footer className="mt-20 border-t border-[var(--border)] bg-[var(--surface)]">
      <div className="container-shell grid gap-8 py-10 md:grid-cols-[1.5fr_1fr_1fr]">
        <div className="space-y-4">
          <p className="text-lg font-semibold">{site.name}</p>
          <p className="max-w-md text-sm leading-6 text-[var(--text-secondary)]">{site.description}</p>
        </div>

        <div>
          <p className="mb-3 text-sm font-semibold uppercase tracking-[0.2em] text-[var(--text-secondary)]">Explore</p>
          <div className="flex flex-col gap-2 text-sm">
            {navItems.map((item) => (
              <Link key={item.href} href={item.href} className="text-[var(--text-secondary)] transition hover:text-[var(--text-primary)]">
                {item.label}
              </Link>
            ))}
          </div>
        </div>

        <div>
          <p className="mb-3 text-sm font-semibold uppercase tracking-[0.2em] text-[var(--text-secondary)]">Social</p>
          <div className="flex flex-col gap-2 text-sm">
            {socialLinks.map((item) => (
              <a key={item.href} href={item.href} className="text-[var(--text-secondary)] transition hover:text-[var(--text-primary)]">
                {item.label}
              </a>
            ))}
          </div>
        </div>
      </div>
    </footer>
  );
}