import Link from 'next/link';
import type { Metadata } from 'next';
import { getFeaturedProjects, getPosts, getSkills, getSiteStats } from '@/lib/backend';
import { SectionHeading } from '@/components/ui/section-heading';
import { ProjectCard } from '@/components/ui/project-card';
import { PostCard } from '@/components/ui/post-card';
import { site } from '@/lib/site';

export const metadata: Metadata = {
  title: `${site.name} | ${site.tagline}`,
  description: site.description
};

export default async function HomePage() {
  const [featuredProjects, posts, skills, stats] = await Promise.all([
    getFeaturedProjects(),
    getPosts(),
    getSkills(),
    getSiteStats()
  ]);

  return (
    <div>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{
          __html: JSON.stringify({
            '@context': 'https://schema.org',
            '@type': 'WebSite',
            name: site.name,
            url: site.url,
            description: site.description
          })
        }}
      />
      <section className="container-shell py-16 md:py-24">
        <div className="grid items-center gap-12 lg:grid-cols-[1.15fr_0.85fr]">
          <div className="space-y-8">
            <div className="inline-flex items-center gap-2 rounded-full border border-[var(--border)] bg-[var(--surface)] px-4 py-2 text-sm text-[var(--text-secondary)] shadow-[0_10px_40px_rgba(15,23,42,0.05)]">
              <span className="h-2.5 w-2.5 rounded-full bg-[var(--success)]" />
              Phase 0 scaffold ready
            </div>
            <div className="space-y-5">
              <h1 className="max-w-3xl text-5xl font-semibold tracking-tight text-balance md:text-7xl">
                A portfolio system designed to be rebranded, not rebuilt.
              </h1>
              <p className="max-w-2xl text-lg leading-8 text-[var(--text-secondary)]">
                DevFolio is built as a modular foundation for developer portfolios, with a clear structure for content, brand settings, and future CMS integration.
              </p>
            </div>
            <div className="flex flex-wrap gap-4">
              <Link href="/projects" className="rounded-full bg-[var(--accent)] px-6 py-3 text-sm font-semibold text-white transition hover:-translate-y-0.5 hover:bg-[var(--accent-strong)]">
                View projects
              </Link>
              <Link href="/contact" className="rounded-full border border-[var(--border)] bg-[var(--surface)] px-6 py-3 text-sm font-semibold text-[var(--text-primary)] transition hover:-translate-y-0.5 hover:shadow-glow">
                Contact me
              </Link>
            </div>
            <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
              {stats.map((stat) => (
                <div key={stat.label} className="glass-panel rounded-3xl p-5">
                  <p className="text-3xl font-semibold tracking-tight">{stat.value}</p>
                  <p className="mt-2 text-sm text-[var(--text-secondary)]">{stat.label}</p>
                </div>
              ))}
            </div>
          </div>

          <div className="relative">
            <div className="absolute inset-0 -z-10 rounded-[2rem] bg-[var(--accent-soft)] blur-3xl" />
            <div className="hero-grid glass-panel rounded-[2rem] p-8">
              <div className="space-y-6 rounded-[1.6rem] border border-[var(--border)] bg-[var(--surface-strong)] p-6">
                <p className="text-sm font-semibold uppercase tracking-[0.24em] text-[var(--text-secondary)]">Brand snapshot</p>
                <div className="space-y-3">
                  <h2 className="text-3xl font-semibold tracking-tight">Francesco Giglio</h2>
                  <p className="max-w-sm text-sm leading-6 text-[var(--text-secondary)]">
                    Full-stack developer focused on clear interfaces, maintainable systems, and fast shipping.
                  </p>
                </div>
                <div className="flex flex-wrap gap-2">
                  {skills.slice(0, 6).map((skill) => (
                    <span key={skill} className="rounded-full bg-[var(--accent-soft)] px-3 py-1 text-xs font-semibold text-[var(--accent-strong)]">
                      {skill}
                    </span>
                  ))}
                </div>
                <div className="grid gap-3 sm:grid-cols-2">
                  <div className="rounded-2xl bg-[var(--surface-muted)] p-4">
                    <p className="text-xs uppercase tracking-[0.2em] text-[var(--text-secondary)]">Focus</p>
                    <p className="mt-2 font-medium">Portfolio architecture</p>
                  </div>
                  <div className="rounded-2xl bg-[var(--surface-muted)] p-4">
                    <p className="text-xs uppercase tracking-[0.2em] text-[var(--text-secondary)]">Stack</p>
                    <p className="mt-2 font-medium">Next.js + Go API</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="container-shell py-10 md:py-16">
        <SectionHeading eyebrow="Featured work" title="Selected projects built around a reusable content model." description="The first phase ships a live portfolio shell with structured data so each later phase can replace mock content incrementally." />
        <div className="mt-10 grid gap-6 lg:grid-cols-3">
          {featuredProjects.map((project) => (
            <ProjectCard key={project.slug} project={project} />
          ))}
        </div>
      </section>

      <section className="container-shell py-10 md:py-16">
        <SectionHeading eyebrow="Editorial" title="Latest articles prepared for SEO-friendly publishing." description="The blog surface is already structured for later MDX and CMS integration." />
        <div className="mt-10 grid gap-6 lg:grid-cols-3">
          {posts.map((post) => (
            <PostCard key={post.slug} post={post} />
          ))}
        </div>
      </section>

      <section className="container-shell py-10 md:py-16">
        <SectionHeading
          eyebrow="Latest highlights"
          title="Quick access to the newest content and strongest project signals."
          description="This section keeps the front page moving by mixing editorial updates with active project entry points."
        />
        <div className="mt-10 grid gap-6 lg:grid-cols-2">
          <div className="rounded-[2rem] border border-[var(--border)] bg-[var(--surface)] p-6">
            <div className="flex items-center justify-between gap-3">
              <h3 className="text-xl font-semibold tracking-tight">Recent posts</h3>
              <Link href="/blog" className="text-sm font-medium text-[var(--accent)] transition hover:text-[var(--accent-strong)]">
                View all
              </Link>
            </div>
            <div className="mt-5 space-y-4">
              {posts.slice(0, 2).map((post) => (
                <Link
                  key={post.slug}
                  href={`/blog/${post.slug}`}
                  className="block rounded-2xl border border-[var(--border)] bg-[var(--surface-strong)] p-4 transition hover:-translate-y-0.5 hover:shadow-glow"
                >
                  <div className="flex items-center justify-between gap-3 text-xs font-semibold uppercase tracking-[0.22em] text-[var(--text-secondary)]">
                    <span>{post.category}</span>
                    <span>{post.readTime}</span>
                  </div>
                  <p className="mt-3 text-base font-semibold tracking-tight">{post.title}</p>
                  <p className="mt-2 text-sm leading-6 text-[var(--text-secondary)]">{post.excerpt}</p>
                </Link>
              ))}
            </div>
          </div>

          <div className="rounded-[2rem] border border-[var(--border)] bg-[var(--surface)] p-6">
            <div className="flex items-center justify-between gap-3">
              <h3 className="text-xl font-semibold tracking-tight">Featured projects</h3>
              <Link href="/projects" className="text-sm font-medium text-[var(--accent)] transition hover:text-[var(--accent-strong)]">
                View all
              </Link>
            </div>
            <div className="mt-5 space-y-4">
              {featuredProjects.slice(0, 2).map((project) => (
                <Link
                  key={project.slug}
                  href={`/projects/${project.slug}`}
                  className="block rounded-2xl border border-[var(--border)] bg-[var(--surface-strong)] p-4 transition hover:-translate-y-0.5 hover:shadow-glow"
                >
                  <div className="flex items-center justify-between gap-3 text-xs font-semibold uppercase tracking-[0.22em] text-[var(--text-secondary)]">
                    <span>{project.year}</span>
                    <span>{project.status}</span>
                  </div>
                  <p className="mt-3 text-base font-semibold tracking-tight">{project.title}</p>
                  <p className="mt-2 text-sm leading-6 text-[var(--text-secondary)]">{project.descriptionShort}</p>
                </Link>
              ))}
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}