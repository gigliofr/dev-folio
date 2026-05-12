import type { Metadata } from 'next';
import { getSkills } from '@/lib/backend';
import { SectionHeading } from '@/components/ui/section-heading';
import { site } from '@/lib/site';

export const metadata: Metadata = {
  title: 'About',
  description: 'A short professional bio, career timeline, and the skills that matter for this stack.',
  openGraph: {
    title: `About | ${site.name}`,
    description: 'A short professional bio, career timeline, and the skills that matter for this stack.'
  }
};

const timeline = [
  { year: '2026', title: 'DevFolio', text: 'Starting the modular portfolio platform and shaping its content model.' },
  { year: '2025', title: 'Frontend systems', text: 'Focused on design systems, accessibility, and product-oriented UI delivery.' },
  { year: '2024', title: 'Full-stack growth', text: 'Expanded work on APIs, data modeling, and maintainable content-driven sites.' }
];

export default async function AboutPage() {
  const skills = await getSkills();

  return (
    <section className="container-shell py-16 md:py-24">
      <SectionHeading
        eyebrow="About"
        title="A short professional bio, career timeline, and the skills that matter for this stack."
        description="This page is already aligned with the spec so it can grow into a richer profile without changing its structure."
      />

      <div className="mt-10 grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
        <div className="rounded-3xl border border-[var(--border)] bg-[var(--surface)] p-8">
          <p className="text-base leading-8 text-[var(--text-secondary)]">
            I build interfaces and content systems that stay readable as they grow. The immediate goal for DevFolio is to establish a clean foundation for the future admin panel, project archive, and editorial workflow.
          </p>
          <div className="mt-8 space-y-4">
            {timeline.map((item) => (
              <div key={item.year} className="flex gap-4 rounded-2xl border border-[var(--border)] bg-[var(--surface-strong)] p-4">
                <div className="w-16 shrink-0 text-sm font-semibold text-[var(--accent)]">{item.year}</div>
                <div>
                  <p className="font-semibold">{item.title}</p>
                  <p className="mt-1 text-sm leading-6 text-[var(--text-secondary)]">{item.text}</p>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="rounded-3xl border border-[var(--border)] bg-[var(--surface)] p-8">
          <p className="text-sm font-semibold uppercase tracking-[0.2em] text-[var(--text-secondary)]">Skills</p>
          <div className="mt-5 flex flex-wrap gap-2">
            {skills.map((skill) => (
              <span key={skill} className="rounded-full bg-[var(--accent-soft)] px-3 py-1 text-xs font-semibold text-[var(--accent-strong)]">
                {skill}
              </span>
            ))}
          </div>
          <div className="mt-8 rounded-2xl bg-[var(--surface-muted)] p-5">
            <p className="text-sm font-semibold uppercase tracking-[0.18em] text-[var(--text-secondary)]">What I do</p>
            <p className="mt-3 text-sm leading-6 text-[var(--text-secondary)]">
              Product-minded frontend work, content-first architecture, and portable design systems that are ready for white-label adaptation.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}