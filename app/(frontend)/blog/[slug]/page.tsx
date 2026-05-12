import { notFound } from 'next/navigation';
import type { Metadata } from 'next';
import { getPostBySlug } from '@/lib/backend';
import { posts } from '@/lib/content';
import { site } from '@/lib/site';

type BlogDetailsPageProps = {
  params: {
    slug: string;
  };
};

export async function generateStaticParams() {
  return posts.map((post) => ({ slug: post.slug }));
}

export async function generateMetadata({ params }: BlogDetailsPageProps): Promise<Metadata> {
  const post = await getPostBySlug(params.slug);
  if (!post) {
    return { title: 'Blog' };
  }

  return {
    title: post.title,
    description: post.excerpt,
    openGraph: {
      title: `${post.title} | ${site.name}`,
      description: post.excerpt
    }
  };
}

export default async function BlogDetailsPage({ params }: BlogDetailsPageProps) {
  const post = await getPostBySlug(params.slug);

  if (!post) {
    notFound();
  }

  return (
    <section className="container-shell py-16 md:py-24">
      <div className="grid gap-10 lg:grid-cols-[minmax(0,1.2fr)_minmax(280px,0.8fr)]">
        <article className="max-w-3xl space-y-6">
          <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[var(--accent)]">{post.category}</p>
          <h1 className="text-4xl font-semibold tracking-tight md:text-6xl">{post.title}</h1>
          <p className="text-lg leading-8 text-[var(--text-secondary)]">{post.excerpt}</p>
          <div className="flex flex-wrap gap-2">
            {post.tags.map((tag) => (
              <span key={tag} className="rounded-full border border-[var(--border)] px-3 py-1 text-xs text-[var(--text-secondary)]">
                {tag}
              </span>
            ))}
          </div>

          <div className="prose prose-slate max-w-none prose-headings:tracking-tight prose-p:leading-8 prose-li:leading-8 dark:prose-invert">
            {post.sections?.map((section) => (
              <section key={section.heading} className="mt-10">
                <h2>{section.heading}</h2>
                {section.paragraphs.map((paragraph) => (
                  <p key={paragraph}>{paragraph}</p>
                ))}
                {section.bullets && (
                  <ul>
                    {section.bullets.map((bullet) => (
                      <li key={bullet}>{bullet}</li>
                    ))}
                  </ul>
                )}
              </section>
            ))}
          </div>
        </article>

        <aside className="space-y-4 lg:sticky lg:top-24 lg:self-start">
          <div className="rounded-3xl border border-[var(--border)] bg-[var(--surface)] p-6">
            <p className="text-xs font-semibold uppercase tracking-[0.24em] text-[var(--text-secondary)]">Article info</p>
            <dl className="mt-4 space-y-3 text-sm">
              <div className="flex items-center justify-between gap-4">
                <dt className="text-[var(--text-secondary)]">Published</dt>
                <dd className="font-medium">{post.publishedAt}</dd>
              </div>
              <div className="flex items-center justify-between gap-4">
                <dt className="text-[var(--text-secondary)]">Read time</dt>
                <dd className="font-medium">{post.readTime}</dd>
              </div>
              <div className="flex items-center justify-between gap-4">
                <dt className="text-[var(--text-secondary)]">Category</dt>
                <dd className="font-medium">{post.category}</dd>
              </div>
            </dl>
          </div>

          <div className="rounded-3xl border border-[var(--border)] bg-[var(--surface)] p-6">
            <p className="text-xs font-semibold uppercase tracking-[0.24em] text-[var(--text-secondary)]">Quick takeaway</p>
            <p className="mt-3 text-sm leading-7 text-[var(--text-secondary)]">
              The strongest portfolio content model is the one that keeps routing, metadata, and body sections separate, so the editor can grow without rewriting the shell.
            </p>
          </div>
        </aside>
      </div>
    </section>
  );
}