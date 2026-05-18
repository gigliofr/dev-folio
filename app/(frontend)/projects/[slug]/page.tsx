import { notFound } from 'next/navigation';
import type { Metadata } from 'next';
import Image from 'next/image';
import { getProjectBySlug } from '@/lib/backend';
import { projects } from '@/lib/content';
import { site } from '@/lib/site';

type ProjectDetailsPageProps = {
  params: {
    slug: string;
  };
};

export async function generateStaticParams() {
  return projects.map((project) => ({ slug: project.slug }));
}

export async function generateMetadata({ params }: ProjectDetailsPageProps): Promise<Metadata> {
  const project = await getProjectBySlug(params.slug);
  if (!project) {
    return { title: 'Projects' };
  }

  return {
    title: project.title,
    description: project.descriptionShort,
    openGraph: {
      title: `${project.title} | ${site.name}`,
      description: project.descriptionShort
    }
  };
}

export default async function ProjectDetailsPage({ params }: ProjectDetailsPageProps) {
  const project = await getProjectBySlug(params.slug);

  if (!project) {
    notFound();
  }

  return (
    <section className="container-shell py-16 md:py-24">
      <div className="grid gap-10 lg:grid-cols-[minmax(0,1.15fr)_minmax(280px,0.85fr)]">
        <article className="max-w-3xl space-y-6">
          <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[var(--accent)]">{project.year}</p>
          <h1 className="text-4xl font-semibold tracking-tight md:text-6xl">{project.title}</h1>
          <p className="text-lg leading-8 text-[var(--text-secondary)]">{project.descriptionLong}</p>

          <div className="overflow-hidden rounded-[2rem] border border-[var(--border)] bg-[var(--surface)]">
            <Image src={project.image} alt={project.title} width={1200} height={600} className="h-80 w-full object-cover" />
          </div>

          <div className="flex flex-wrap gap-2">
            {project.technologies.map((technology) => (
              <span key={technology} className="rounded-full bg-[var(--accent-soft)] px-3 py-1 text-xs font-semibold text-[var(--accent-strong)]">
                {technology}
              </span>
            ))}
          </div>

          {(project.liveUrl || project.githubUrl) && (
            <div className="flex flex-wrap gap-3 pt-2">
              {project.liveUrl && (
                <a href={project.liveUrl} target="_blank" rel="noreferrer" className="rounded-full bg-[var(--accent)] px-5 py-3 text-sm font-semibold text-white transition hover:-translate-y-0.5 hover:bg-[var(--accent-strong)]">
                  Live demo
                </a>
              )}
              {project.githubUrl && (
                <a href={project.githubUrl} target="_blank" rel="noreferrer" className="rounded-full border border-[var(--border)] bg-[var(--surface)] px-5 py-3 text-sm font-semibold text-[var(--text-primary)] transition hover:-translate-y-0.5 hover:shadow-glow">
                  Source code
                </a>
              )}
            </div>
          )}
        </article>

        <aside className="space-y-4 lg:sticky lg:top-24 lg:self-start">
          <div className="rounded-3xl border border-[var(--border)] bg-[var(--surface)] p-6">
            <p className="text-xs font-semibold uppercase tracking-[0.24em] text-[var(--text-secondary)]">Project snapshot</p>
            <dl className="mt-4 space-y-3 text-sm">
              <div className="flex items-center justify-between gap-4">
                <dt className="text-[var(--text-secondary)]">Status</dt>
                <dd className="font-medium capitalize">{project.status}</dd>
              </div>
              <div className="flex items-center justify-between gap-4">
                <dt className="text-[var(--text-secondary)]">Year</dt>
                <dd className="font-medium">{project.year}</dd>
              </div>
              <div>
                <dt className="text-[var(--text-secondary)]">Stack</dt>
                <dd className="mt-2 flex flex-wrap gap-2">
                  {project.technologies.map((technology) => (
                    <span key={technology} className="rounded-full border border-[var(--border)] px-3 py-1 text-xs text-[var(--text-secondary)]">
                      {technology}
                    </span>
                  ))}
                </dd>
              </div>
            </dl>
          </div>

          <div className="rounded-3xl border border-[var(--border)] bg-[var(--surface)] p-6">
            <p className="text-xs font-semibold uppercase tracking-[0.24em] text-[var(--text-secondary)]">Why it matters</p>
            <p className="mt-3 text-sm leading-7 text-[var(--text-secondary)]">
              A portfolio project page should show the implementation and the decision-making behind it, so the archive doubles as proof of process and technical depth.
            </p>
          </div>
        </aside>
      </div>
    </section>
  );
}