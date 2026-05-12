import type { Metadata } from 'next';
import { getPosts } from '@/lib/backend';
import { SectionHeading } from '@/components/ui/section-heading';
import { site } from '@/lib/site';
import { BlogBrowser } from './blog-browser';

export const metadata: Metadata = {
  title: 'Blog',
  description: 'Articles that will later be powered by MDX and the CMS.',
  openGraph: {
    title: `Blog | ${site.name}`,
    description: 'Articles that will later be powered by MDX and the CMS.'
  }
};

export default async function BlogPage() {
  const posts = await getPosts();

  return (
    <section className="container-shell py-16 md:py-24">
      <SectionHeading eyebrow="Blog" title="Articles that will later be powered by MDX and the CMS." description="The editorial layer is seeded locally for now. In the next phase it will be replaced by typed content collections and richer article templates." />
      <div className="mt-10">
        <BlogBrowser posts={posts} />
      </div>
    </section>
  );
}