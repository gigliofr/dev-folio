import Image from 'next/image';
import Link from 'next/link';
import type { Project } from '@/lib/content';

type ProjectCardProps = {
  project: Project;
  compact?: boolean;
};

export function ProjectCard({ project, compact = false }: ProjectCardProps) {
  return (
    <article className="overflow-hidden rounded-3xl border border-[var(--border)] bg-[var(--surface)] shadow-[0_20px_80px_rgba(15,23,42,0.08)] transition hover:-translate-y-1 hover:shadow-glow">
      <div className="relative aspect-[16/10]">
        <Image src={project.image} alt={project.title} fill className="object-cover" sizes="(max-width: 768px) 100vw, 50vw" />
      </div>
      <div className={`space-y-4 p-6 ${compact ? 'p-5' : ''}`}>
        <div className="flex items-center justify-between gap-3 text-xs font-semibold uppercase tracking-[0.22em] text-[var(--text-secondary)]">
          <span>{project.year}</span>
          <span>{project.status}</span>
        </div>
        <div className="space-y-2">
          <h3 className="text-xl font-semibold tracking-tight">{project.title}</h3>
          <p className="text-sm leading-6 text-[var(--text-secondary)]">{project.descriptionShort}</p>
        </div>
        <div className="flex flex-wrap gap-2">
          {project.technologies.map((technology) => (
            <span key={technology} className="rounded-full bg-[var(--accent-soft)] px-3 py-1 text-xs font-semibold text-[var(--accent-strong)]">
              {technology}
            </span>
          ))}
        </div>
        <div className="flex items-center gap-3 pt-2 text-sm font-medium">
          <Link href={`/projects/${project.slug}`} className="text-[var(--accent)] transition hover:text-[var(--accent-strong)]">
            View details
          </Link>
          {project.liveUrl ? (
            <a href={project.liveUrl} className="text-[var(--text-secondary)] transition hover:text-[var(--text-primary)]">
              Live demo
            </a>
          ) : null}
        </div>
      </div>
    </article>
  );
}