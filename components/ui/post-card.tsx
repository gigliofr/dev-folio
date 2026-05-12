import Link from 'next/link';
import type { Post } from '@/lib/content';

type PostCardProps = {
  post: Post;
};

export function PostCard({ post }: PostCardProps) {
  return (
    <article className="rounded-3xl border border-[var(--border)] bg-[var(--surface)] p-6 transition hover:-translate-y-1 hover:shadow-glow">
      <div className="flex items-center justify-between gap-3 text-xs font-semibold uppercase tracking-[0.22em] text-[var(--text-secondary)]">
        <span>{post.category}</span>
        <span>{post.readTime}</span>
      </div>
      <h3 className="mt-4 text-xl font-semibold tracking-tight">{post.title}</h3>
      <p className="mt-3 text-sm leading-6 text-[var(--text-secondary)]">{post.excerpt}</p>
      <div className="mt-5 flex flex-wrap gap-2">
        {post.tags.map((tag) => (
          <span key={tag} className="rounded-full border border-[var(--border)] px-3 py-1 text-xs text-[var(--text-secondary)]">
            {tag}
          </span>
        ))}
      </div>
      <div className="mt-6 flex items-center justify-between text-sm text-[var(--text-secondary)]">
        <span>{post.publishedAt}</span>
        <Link href={`/blog/${post.slug}`} className="font-medium text-[var(--accent)] transition hover:text-[var(--accent-strong)]">
          Read article
        </Link>
      </div>
    </article>
  );
}