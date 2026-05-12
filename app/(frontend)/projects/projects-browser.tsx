'use client';

import { useMemo, useState } from 'react';
import type { Project } from '@/lib/content';
import { ProjectCard } from '@/components/ui/project-card';

type ProjectsBrowserProps = {
  projects: Project[];
};

const statusOptions = ['all', 'featured', 'published', 'draft', 'archived'] as const;

export function ProjectsBrowser({ projects }: ProjectsBrowserProps) {
  const [query, setQuery] = useState('');
  const [status, setStatus] = useState<(typeof statusOptions)[number]>('all');

  const filteredProjects = useMemo(() => {
    const normalizedQuery = query.trim().toLowerCase();

    return projects.filter((project) => {
      const matchesQuery =
        normalizedQuery.length === 0 ||
        [project.title, project.descriptionShort, project.descriptionLong, project.year, project.slug, ...project.technologies]
          .join(' ')
          .toLowerCase()
          .includes(normalizedQuery);

      const matchesStatus =
        status === 'all' ||
        (status === 'featured' && project.featured) ||
        project.status === status;

      return matchesQuery && matchesStatus;
    });
  }, [projects, query, status]);

  return (
    <div className="space-y-6">
      <div className="grid gap-4 rounded-3xl border border-[var(--border)] bg-[var(--surface)] p-5 md:grid-cols-[minmax(0,1fr)_auto] md:items-end">
        <label className="block">
          <span className="text-xs font-semibold uppercase tracking-[0.2em] text-[var(--text-secondary)]">Search</span>
          <input
            type="search"
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search projects, technologies, or years"
            className="mt-2 w-full rounded-2xl border border-[var(--border)] bg-[var(--bg)] px-4 py-3 text-sm text-[var(--text-primary)] outline-none transition focus:border-[var(--accent)]"
          />
        </label>

        <div className="flex flex-wrap gap-2 md:justify-end">
          {statusOptions.map((option) => (
            <button
              key={option}
              type="button"
              onClick={() => setStatus(option)}
              className={`rounded-full border px-4 py-2 text-sm font-medium transition ${
                status === option
                  ? 'border-[var(--accent)] bg-[var(--accent)] text-white'
                  : 'border-[var(--border)] bg-[var(--surface)] text-[var(--text-secondary)] hover:text-[var(--text-primary)]'
              }`}
            >
              {option === 'all' ? 'All' : option === 'featured' ? 'Featured' : option.charAt(0).toUpperCase() + option.slice(1)}
            </button>
          ))}
        </div>
      </div>

      <div className="flex items-center justify-between text-sm text-[var(--text-secondary)]">
        <span>
          Showing <strong className="text-[var(--text-primary)]">{filteredProjects.length}</strong> of {projects.length} projects
        </span>
        {(query || status !== 'all') && (
          <button
            type="button"
            onClick={() => {
              setQuery('');
              setStatus('all');
            }}
            className="font-medium text-[var(--accent)] transition hover:text-[var(--accent-strong)]"
          >
            Clear filters
          </button>
        )}
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {filteredProjects.map((project) => (
          <ProjectCard key={project.slug} project={project} compact />
        ))}
      </div>

      {filteredProjects.length === 0 && (
        <div className="rounded-3xl border border-[var(--border)] bg-[var(--surface)] p-8 text-sm text-[var(--text-secondary)]">
          No projects match the current filters.
        </div>
      )}
    </div>
  );
}