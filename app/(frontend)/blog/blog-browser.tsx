'use client';

import { useMemo, useState } from 'react';
import type { Post } from '@/lib/content';
import { PostCard } from '@/components/ui/post-card';

type BlogBrowserProps = {
  posts: Post[];
};

export function BlogBrowser({ posts }: BlogBrowserProps) {
  const [query, setQuery] = useState('');
  const [category, setCategory] = useState('all');

  const categories = useMemo(() => ['all', ...Array.from(new Set(posts.map((post) => post.category)))], [posts]);

  const filteredPosts = useMemo(() => {
    const normalizedQuery = query.trim().toLowerCase();

    return posts.filter((post) => {
      const matchesCategory = category === 'all' || post.category === category;
      const matchesQuery =
        normalizedQuery.length === 0 ||
        [post.title, post.excerpt, post.category, post.readTime, post.publishedAt, ...post.tags].join(' ').toLowerCase().includes(normalizedQuery);

      return matchesCategory && matchesQuery;
    });
  }, [posts, query, category]);

  return (
    <div className="space-y-6">
      <div className="grid gap-4 rounded-3xl border border-[var(--border)] bg-[var(--surface)] p-5 md:grid-cols-[minmax(0,1fr)_auto] md:items-end">
        <label className="block">
          <span className="text-xs font-semibold uppercase tracking-[0.2em] text-[var(--text-secondary)]">Search</span>
          <input
            type="search"
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search posts, tags, or dates"
            className="mt-2 w-full rounded-2xl border border-[var(--border)] bg-[var(--bg)] px-4 py-3 text-sm text-[var(--text-primary)] outline-none transition focus:border-[var(--accent)]"
          />
        </label>

        <div className="flex flex-wrap gap-2 md:justify-end">
          {categories.map((item) => (
            <button
              key={item}
              type="button"
              onClick={() => setCategory(item)}
              className={`rounded-full border px-4 py-2 text-sm font-medium transition ${
                category === item
                  ? 'border-[var(--accent)] bg-[var(--accent)] text-white'
                  : 'border-[var(--border)] bg-[var(--surface)] text-[var(--text-secondary)] hover:text-[var(--text-primary)]'
              }`}
            >
              {item === 'all' ? 'All' : item}
            </button>
          ))}
        </div>
      </div>

      <div className="flex items-center justify-between text-sm text-[var(--text-secondary)]">
        <span>
          Showing <strong className="text-[var(--text-primary)]">{filteredPosts.length}</strong> of {posts.length} posts
        </span>
        {(query || category !== 'all') && (
          <button
            type="button"
            onClick={() => {
              setQuery('');
              setCategory('all');
            }}
            className="font-medium text-[var(--accent)] transition hover:text-[var(--accent-strong)]"
          >
            Clear filters
          </button>
        )}
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {filteredPosts.map((post) => (
          <PostCard key={post.slug} post={post} />
        ))}
      </div>

      {filteredPosts.length === 0 && (
        <div className="rounded-3xl border border-[var(--border)] bg-[var(--surface)] p-8 text-sm text-[var(--text-secondary)]">
          No posts match the current filters.
        </div>
      )}
    </div>
  );
}