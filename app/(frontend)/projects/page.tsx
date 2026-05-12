import type { Metadata } from 'next';
import { getProjects } from '@/lib/backend';
import { SectionHeading } from '@/components/ui/section-heading';
import { site } from '@/lib/site';
import { ProjectsBrowser } from './projects-browser';

export const metadata: Metadata = {
  title: 'Projects',
  description: 'A structured project archive ready for filters, featured flags, and future CMS data.',
  openGraph: {
    title: `Projects | ${site.name}`,
    description: 'A structured project archive ready for filters, featured flags, and future CMS data.'
  }
};

export default async function ProjectsPage() {
  const projects = await getProjects();

  return (
    <section className="container-shell py-16 md:py-24">
      <SectionHeading
        eyebrow="Projects"
        title="A structured project archive ready for filters, featured flags, and future CMS data."
        description="This page already matches the content model described in the spec, even though the data currently comes from a local seed layer."
      />
      <div className="mt-10">
        <ProjectsBrowser projects={projects} />
      </div>
    </section>
  );
}