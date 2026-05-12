import type { MetadataRoute } from 'next';
import { posts, projects } from '@/lib/content';

export default function sitemap(): MetadataRoute.Sitemap {
  const baseUrl = 'https://devfolio.local';

  const staticRoutes = ['/', '/projects', '/blog', '/about', '/contact'];
  const dynamicRoutes = [
    ...projects.map((project) => `/projects/${project.slug}`),
    ...posts.map((post) => `/blog/${post.slug}`)
  ];

  return [...staticRoutes, ...dynamicRoutes].map((path) => ({
    url: `${baseUrl}${path}`,
    lastModified: new Date()
  }));
}