import { posts, projects, skills, stats } from '@/lib/content';
import type { Post, Project } from '@/lib/content';

const apiBaseUrl = process.env.DEVFOLIO_API_URL?.replace(/\/$/, '');

async function fetchJson<T>(path: string, fallback: T): Promise<T> {
  if (!apiBaseUrl) {
    return fallback;
  }

  try {
    const response = await fetch(`${apiBaseUrl}${path}`, { cache: 'no-store' });
    if (!response.ok) {
      return fallback;
    }
    return (await response.json()) as T;
  } catch {
    return fallback;
  }
}

export async function getProjects() {
  return fetchJson<Project[]>('/api/v1/projects?status=published', projects.filter((project) => project.status === 'published'));
}

export async function getFeaturedProjects() {
  return fetchJson<Project[]>('/api/v1/projects?featured=true&status=published', projects.filter((project) => project.featured && project.status === 'published'));
}

export async function getPosts() {
  return fetchJson<Post[]>('/api/v1/posts', posts);
}

export async function getSiteStats() {
  return fetchJson('/api/v1/stats', stats);
}

export async function getSkills() {
  return fetchJson<string[]>('/api/v1/skills', skills);
}

export async function getProjectBySlug(slug: string) {
  return fetchJson<Project | null>(`/api/v1/projects/${slug}`, projects.find((project) => project.slug === slug) ?? null);
}

export async function getPostBySlug(slug: string) {
  return fetchJson<Post | null>(`/api/v1/posts/${slug}`, posts.find((post) => post.slug === slug) ?? null);
}

function authHeader(): Record<string, string> {
  if (typeof window === 'undefined') return {};
  const token = localStorage.getItem('devfolio_token');
  return token ? { Authorization: `Bearer ${token}` } : {};
}

async function mutateJson<T>(path: string, method: string, body?: any): Promise<T | null> {
  const base = apiBaseUrl ?? '';
  try {
    const res = await fetch(`${base}${path}`, {
      method,
      headers: {
        'Content-Type': 'application/json',
        ...authHeader(),
      },
      body: body ? JSON.stringify(body) : undefined,
    });
    if (!res.ok) return null;
    return (await res.json()) as T;
  } catch {
    return null;
  }
}

export async function adminLogin(credentials: { username?: string; password?: string; token?: string }) {
  const base = apiBaseUrl ?? '';
  const res = await fetch(`${base}/api/v1/admin/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(credentials),
  });
  if (!res.ok) return null;
  const data = await res.json();
  if (data?.token) {
    localStorage.setItem('devfolio_token', data.token);
    return data.token;
  }
  return null;
}

export async function adminMe() {
  const base = apiBaseUrl ?? '';
  try {
    const res = await fetch(`${base}/api/v1/admin/me`, { headers: authHeader() });
    if (!res.ok) return null;
    return await res.json();
  } catch {
    return null;
  }
}

export async function createProject(project: Partial<Project>) {
  return mutateJson<Project>('/api/v1/projects', 'POST', project);
}

export async function updateProject(slug: string, project: Partial<Project>) {
  return mutateJson<Project>(`/api/v1/projects/${slug}`, 'PUT', project);
}

export async function deleteProject(slug: string) {
  return fetch(`${apiBaseUrl ?? ''}/api/v1/projects/${slug}`, { method: 'DELETE', headers: authHeader() }).then((r) => r.ok);
}

export async function createPost(post: Partial<Post>) {
  return mutateJson<Post>('/api/v1/posts', 'POST', post);
}

export async function updatePost(slug: string, post: Partial<Post>) {
  return mutateJson<Post>(`/api/v1/posts/${slug}`, 'PUT', post);
}

export async function deletePost(slug: string) {
  return fetch(`${apiBaseUrl ?? ''}/api/v1/posts/${slug}`, { method: 'DELETE', headers: authHeader() }).then((r) => r.ok);
}

export async function uploadFile(file: File): Promise<string | null> {
  const base = apiBaseUrl ?? '';
  const formData = new FormData();
  formData.append('file', file);
  
  try {
    const res = await fetch(`${base}/api/v1/upload`, {
      method: 'POST',
      headers: authHeader(),
      body: formData,
    });
    if (!res.ok) return null;
    const data = (await res.json()) as { url: string };
    return data.url;
  } catch {
    return null;
  }
}

export async function submitContact(name: string, email: string, message: string): Promise<boolean> {
  const base = apiBaseUrl ?? '';
  try {
    const res = await fetch(`${base}/api/v1/contact`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, email, message }),
    });
    return res.ok;
  } catch {
    return false;
  }
}

export async function getContactSubmissions(): Promise<Array<{ name: string; email: string; message: string; createdAt?: string }>> {
  const base = apiBaseUrl ?? '';
  try {
    const res = await fetch(`${base}/api/v1/contact`, { headers: authHeader() });
    if (!res.ok) return [];
    return (await res.json()) as Array<{ name: string; email: string; message: string; createdAt?: string }>;
  } catch {
    return [];
  }
}