"use client";

import { useEffect, useState } from 'react';
import { adminLogin, adminMe, getProjects, getPosts, createProject, updateProject, deleteProject, createPost, updatePost, deletePost, uploadFile, getContactSubmissions, clearAuthToken, getAuthToken } from '@/lib/backend';
import type { Project, Post } from '@/lib/content';

type Tab = 'projects' | 'posts' | 'contacts';
type EditingProject = (Project & { _isNew?: boolean }) | null;
type EditingPost = (Post & { _isNew?: boolean }) | null;

export default function AdminPage() {
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<any>(null);
  const [tab, setTab] = useState<Tab>('projects');
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');
  
  // Auth form
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  // Projects
  const [projects, setProjects] = useState<Project[]>([]);
  const [editingProject, setEditingProject] = useState<EditingProject>(null);
  const [projectForm, setProjectForm] = useState<Partial<Project>>({});

  // Posts
  const [posts, setPosts] = useState<Post[]>([]);
  const [editingPost, setEditingPost] = useState<EditingPost>(null);
  const [postForm, setPostForm] = useState<Partial<Post>>({});
  // Contacts
  const [contacts, setContacts] = useState<Array<{ name: string; email: string; message: string; createdAt?: string }>>([]);
  const [sortKey, setSortKey] = useState<'createdAt' | 'name' | 'email'>('createdAt');
  const [sortDir, setSortDir] = useState<'asc' | 'desc'>('desc');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  // Computed values for contacts pagination
  const sortedContacts = [...contacts].sort((a, b) => {
    let aVal: any = a[sortKey] ?? '';
    let bVal: any = b[sortKey] ?? '';
    if (sortKey === 'createdAt') {
      aVal = new Date(aVal as string).getTime();
      bVal = new Date(bVal as string).getTime();
    }
    const cmp = aVal < bVal ? -1 : aVal > bVal ? 1 : 0;
    return sortDir === 'asc' ? cmp : -cmp;
  });
  const totalContacts = sortedContacts.length;
  const paginatedContacts = sortedContacts.slice((page - 1) * pageSize, page * pageSize);
  const maxPage = Math.ceil(totalContacts / pageSize) || 1;

  useEffect(() => {
    const t = getAuthToken();
    if (t) {
      setToken(t);
      adminMe().then((me) => setUser(me));
    }
    refreshData();
  }, []);

  async function refreshData() {
    setProjects(await getProjects() || []);
    setPosts(await getPosts() || []);
    try {
      const token = getAuthToken();
      if (token) {
        setContacts(await getContactSubmissions() || []);
      }
    } catch {
      // ignore
    }
  }

  function showMessage(msg: string) {
    setMessage(msg);
    setTimeout(() => setMessage(''), 3000);
  }

  // Auth
  async function handleLogin(e: any) {
    e.preventDefault();
    setLoading(true);
    const t = await adminLogin({ username, password });
    if (t) {
      setToken(t);
      setUser(await adminMe());
      setUsername('');
      setPassword('');
      await refreshData();
      showMessage('✓ Logged in');
    } else {
      showMessage('✗ Login failed');
    }
    setLoading(false);
  }

  function handleLogout() {
    setToken(null);
    setUser(null);
    clearAuthToken();
  }

  // Projects
  function newProject() {
    setEditingProject({ title: '', slug: '', descriptionShort: '', descriptionLong: '', technologies: [], status: 'draft', featured: false, year: new Date().getFullYear().toString(), image: '', _isNew: true });
    setProjectForm({});
  }

  async function saveProject() {
    if (!editingProject?.slug || !editingProject?.title) {
      showMessage('✗ Title and slug required');
      return;
    }
    setLoading(true);
    const payload = { ...editingProject, _isNew: undefined };
    const result = editingProject._isNew ? await createProject(payload) : await updateProject(editingProject.slug, payload);
    if (result) {
      showMessage(`✓ Project ${editingProject._isNew ? 'created' : 'updated'}`);
      refreshData();
      setEditingProject(null);
    } else {
      showMessage('✗ Error saving project');
    }
    setLoading(false);
  }

  async function removeProject(slug: string) {
    if (!confirm('Delete this project?')) return;
    setLoading(true);
    const ok = await deleteProject(slug);
    if (ok) {
      showMessage('✓ Project deleted');
      refreshData();
    } else {
      showMessage('✗ Error deleting project');
    }
    setLoading(false);
  }

  async function handleImageUpload(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    setLoading(true);
    const url = await uploadFile(file);
    if (url && editingProject) {
      setEditingProject({ ...editingProject, image: url });
      showMessage('✓ Image uploaded');
    } else {
      showMessage('✗ Upload failed');
    }
    setLoading(false);
  }

  // Posts
  function newPost() {
    setEditingPost({ title: '', slug: '', excerpt: '', category: '', readTime: '', publishedAt: new Date().toISOString().split('T')[0], tags: [], _isNew: true });
  }

  async function savePost() {
    if (!editingPost?.slug || !editingPost?.title) {
      showMessage('✗ Title and slug required');
      return;
    }
    setLoading(true);
    const payload = { ...editingPost, _isNew: undefined };
    const result = editingPost._isNew ? await createPost(payload) : await updatePost(editingPost.slug, payload);
    if (result) {
      showMessage(`✓ Post ${editingPost._isNew ? 'created' : 'updated'}`);
      refreshData();
      setEditingPost(null);
    } else {
      showMessage('✗ Error saving post');
    }
    setLoading(false);
  }

  async function removePost(slug: string) {
    if (!confirm('Delete this post?')) return;
    setLoading(true);
    const ok = await deletePost(slug);
    if (ok) {
      showMessage('✓ Post deleted');
      refreshData();
    } else {
      showMessage('✗ Error deleting post');
    }
    setLoading(false);
  }

  return (
    <section className="container-shell py-16 md:py-24">
      <div className="max-w-5xl rounded-3xl border border-[var(--border)] bg-[var(--surface)] p-8">
        <div className="flex items-end justify-between">
          <div>
            <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[var(--accent)]">Admin</p>
            <h1 className="mt-4 text-2xl font-semibold tracking-tight">Admin Console</h1>
          </div>
          {token && <button className="btn-ghost text-sm" onClick={handleLogout}>Sign out</button>}
        </div>

        {message && (
          <div className="mt-6 rounded-lg border border-[var(--border)] bg-[var(--surface-alt)] px-4 py-2 text-sm">
            {message}
          </div>
        )}

        {!token ? (
          <form onSubmit={handleLogin} className="mt-6 grid max-w-sm gap-3">
            <input placeholder="Username" value={username} onChange={(e) => setUsername(e.target.value)} className="input" />
            <input placeholder="Password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} className="input" />
            <button className="btn-primary" disabled={loading}>{loading ? 'Signing in...' : 'Sign in'}</button>
          </form>
        ) : (
          <div className="mt-8">
            <p className="text-sm text-[var(--text-secondary)]">Signed in as <strong>{user?.description ?? 'admin'}</strong></p>

            {/* Tabs */}
            <div className="mt-8 flex gap-4 border-b border-[var(--border)]">
              <button
                onClick={() => setTab('projects')}
                aria-pressed={tab === 'projects'}
                className={`pb-3 text-sm font-medium ${tab === 'projects' ? 'border-b-2 border-[var(--accent)] text-[var(--accent)]' : 'text-[var(--text-secondary)]'}`}
              >
                Projects ({projects.length})
              </button>
              <button
                onClick={() => setTab('posts')}
                aria-pressed={tab === 'posts'}
                className={`pb-3 text-sm font-medium ${tab === 'posts' ? 'border-b-2 border-[var(--accent)] text-[var(--accent)]' : 'text-[var(--text-secondary)]'}`}
              >
                Posts ({posts.length})
              </button>
              <button
                onClick={() => setTab('contacts')}
                aria-pressed={tab === 'contacts'}
                className={`pb-3 text-sm font-medium ${tab === 'contacts' ? 'border-b-2 border-[var(--accent)] text-[var(--accent)]' : 'text-[var(--text-secondary)]'}`}
              >
                Contacts ({contacts.length})
              </button>
            </div>

            {/* Projects Tab */}
            {tab === 'projects' && (
              <div className="mt-8">
                {editingProject ? (
                  <form className="max-w-2xl space-y-4" onSubmit={(e) => { e.preventDefault(); saveProject(); }}>
                    <h3 className="text-lg font-medium">{editingProject._isNew ? 'New Project' : 'Edit Project'}</h3>
                    <div>
                      <label className="text-sm font-medium">Title *</label>
                      <input
                        type="text"
                        value={editingProject.title ?? ''}
                        onChange={(e) => setEditingProject({ ...editingProject, title: e.target.value })}
                        className="input mt-1 w-full"
                        placeholder="Project title"
                      />
                    </div>
                    <div>
                      <label className="text-sm font-medium">Slug *</label>
                      <input
                        type="text"
                        value={editingProject.slug ?? ''}
                        onChange={(e) => setEditingProject({ ...editingProject, slug: e.target.value })}
                        className="input mt-1 w-full"
                        placeholder="project-slug"
                      />
                    </div>
                    <div>
                      <label className="text-sm font-medium">Short Description</label>
                      <textarea
                        value={editingProject.descriptionShort ?? ''}
                        onChange={(e) => setEditingProject({ ...editingProject, descriptionShort: e.target.value })}
                        className="input mt-1 w-full resize-none"
                        rows={2}
                        placeholder="Brief description"
                      />
                    </div>
                    <div>
                      <label className="text-sm font-medium">Long Description</label>
                      <textarea
                        value={editingProject.descriptionLong ?? ''}
                        onChange={(e) => setEditingProject({ ...editingProject, descriptionLong: e.target.value })}
                        className="input mt-1 w-full resize-none"
                        rows={4}
                        placeholder="Full description"
                      />
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="text-sm font-medium">Status</label>
                        <select value={editingProject.status ?? 'draft'} onChange={(e) => setEditingProject({ ...editingProject, status: e.target.value as 'draft' | 'published' | 'archived' })} className="input mt-1 w-full">
                          <option value="draft">Draft</option>
                          <option value="published">Published</option>
                          <option value="archived">Archived</option>
                        </select>
                      </div>
                      <div>
                        <label className="text-sm font-medium">Year</label>
                        <input
                          type="text"
                          value={editingProject.year ?? ''}
                          onChange={(e) => setEditingProject({ ...editingProject, year: e.target.value })}
                          className="input mt-1 w-full"
                        />
                      </div>
                    </div>
                    <div>
                      <label className="flex gap-2">
                        <input
                          type="checkbox"
                          checked={editingProject.featured ?? false}
                          onChange={(e) => setEditingProject({ ...editingProject, featured: e.target.checked })}
                        />
                        <span className="text-sm">Featured</span>
                      </label>
                    </div>
                    <div>
                      <label className="text-sm font-medium">Technologies (comma-separated)</label>
                      <input
                        type="text"
                        value={(editingProject.technologies ?? []).join(', ')}
                        onChange={(e) => setEditingProject({ ...editingProject, technologies: e.target.value.split(',').map(t => t.trim()) })}
                        className="input mt-1 w-full"
                        placeholder="React, Node.js, PostgreSQL"
                      />
                    </div>
                    <div>
                      <label className="text-sm font-medium">Image URL</label>
                      <input
                        type="text"
                        value={editingProject.image ?? ''}
                        onChange={(e) => setEditingProject({ ...editingProject, image: e.target.value })}
                        className="input mt-1 w-full"
                        placeholder="https://..."
                      />
                      <div className="mt-3">
                        <label className="text-sm font-medium">Or upload image</label>
                        <input
                          type="file"
                          accept="image/*"
                          onChange={handleImageUpload}
                          disabled={loading}
                          className="input mt-1 w-full"
                        />
                      </div>
                    </div>
                    <div className="flex gap-2 pt-4">
                      <button type="submit" className="btn-primary" disabled={loading}>{loading ? 'Saving...' : 'Save'}</button>
                      <button type="button" className="btn-ghost" onClick={() => setEditingProject(null)}>Cancel</button>
                    </div>
                  </form>
                ) : (
                  <>
                    <button className="btn-primary mb-6" onClick={newProject} disabled={loading}>+ New Project</button>
                    <div className="space-y-2">
                      {projects.map((p) => (
                        <div key={p.slug} className="flex items-center justify-between rounded-lg border border-[var(--border)] p-4">
                          <div className="flex-1">
                            <p className="font-medium">{p.title}</p>
                            <p className="text-xs text-[var(--text-secondary)]">{p.slug} • {p.status}</p>
                          </div>
                          <div className="flex gap-2">
                            <button className="btn-ghost text-sm" onClick={() => setEditingProject(p)}>Edit</button>
                            <button className="btn-ghost text-sm text-red-500" onClick={() => removeProject(p.slug)}>Delete</button>
                          </div>
                        </div>
                      ))}
                    </div>
                  </>
                )}
              </div>
            )}

            {/* Posts Tab */}
            {tab === 'posts' && (
              <div className="mt-8">
                {editingPost ? (
                  <form className="max-w-2xl space-y-4" onSubmit={(e) => { e.preventDefault(); savePost(); }}>
                    <h3 className="text-lg font-medium">{editingPost._isNew ? 'New Post' : 'Edit Post'}</h3>
                    <div>
                      <label className="text-sm font-medium">Title *</label>
                      <input
                        type="text"
                        value={editingPost.title ?? ''}
                        onChange={(e) => setEditingPost({ ...editingPost, title: e.target.value })}
                        className="input mt-1 w-full"
                        placeholder="Post title"
                      />
                    </div>
                    <div>
                      <label className="text-sm font-medium">Slug *</label>
                      <input
                        type="text"
                        value={editingPost.slug ?? ''}
                        onChange={(e) => setEditingPost({ ...editingPost, slug: e.target.value })}
                        className="input mt-1 w-full"
                        placeholder="post-slug"
                      />
                    </div>
                    <div>
                      <label className="text-sm font-medium">Excerpt</label>
                      <textarea
                        value={editingPost.excerpt ?? ''}
                        onChange={(e) => setEditingPost({ ...editingPost, excerpt: e.target.value })}
                        className="input mt-1 w-full resize-none"
                        rows={2}
                        placeholder="Brief excerpt"
                      />
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="text-sm font-medium">Category</label>
                        <input
                          type="text"
                          value={editingPost.category ?? ''}
                          onChange={(e) => setEditingPost({ ...editingPost, category: e.target.value })}
                          className="input mt-1 w-full"
                          placeholder="e.g., Tutorial"
                        />
                      </div>
                      <div>
                        <label className="text-sm font-medium">Read Time</label>
                        <input
                          type="text"
                          value={editingPost.readTime ?? ''}
                          onChange={(e) => setEditingPost({ ...editingPost, readTime: e.target.value })}
                          className="input mt-1 w-full"
                          placeholder="e.g., 5 min"
                        />
                      </div>
                    </div>
                    <div>
                      <label className="text-sm font-medium">Published Date</label>
                      <input
                        type="date"
                        value={editingPost.publishedAt ?? ''}
                        onChange={(e) => setEditingPost({ ...editingPost, publishedAt: e.target.value })}
                        className="input mt-1 w-full"
                      />
                    </div>
                    <div>
                      <label className="text-sm font-medium">Tags (comma-separated)</label>
                      <input
                        type="text"
                        value={(editingPost.tags ?? []).join(', ')}
                        onChange={(e) => setEditingPost({ ...editingPost, tags: e.target.value.split(',').map(t => t.trim()) })}
                        className="input mt-1 w-full"
                        placeholder="React, Next.js, API"
                      />
                    </div>
                    <div className="flex gap-2 pt-4">
                      <button type="submit" className="btn-primary" disabled={loading}>{loading ? 'Saving...' : 'Save'}</button>
                      <button type="button" className="btn-ghost" onClick={() => setEditingPost(null)}>Cancel</button>
                    </div>
                  </form>
                ) : (
                  <>
                    <button className="btn-primary mb-6" onClick={newPost} disabled={loading}>+ New Post</button>
                    <div className="space-y-2">
                      {posts.map((p) => (
                        <div key={p.slug} className="flex items-center justify-between rounded-lg border border-[var(--border)] p-4">
                          <div className="flex-1">
                            <p className="font-medium">{p.title}</p>
                            <p className="text-xs text-[var(--text-secondary)]">{p.slug} • {p.category} • {p.readTime}</p>
                          </div>
                          <div className="flex gap-2">
                            <button className="btn-ghost text-sm" onClick={() => setEditingPost(p)}>Edit</button>
                            <button className="btn-ghost text-sm text-red-500" onClick={() => removePost(p.slug)}>Delete</button>
                          </div>
                        </div>
                      ))}
                    </div>
                  </>
                )}
              </div>
            )}
            {/* Contacts Tab */}
            {tab === 'contacts' && (
              <div className="mt-8">
                {/* Sort & Pagination Controls */}
                <div className="mb-6 flex flex-wrap items-center gap-4">
                  <div className="flex items-center gap-2">
                    <label className="text-sm font-medium">Sort by:</label>
                    <select
                      value={sortKey}
                      onChange={(e) => { setSortKey(e.target.value as any); setPage(1); }}
                      className="input px-2 py-1 text-sm"
                      aria-label="Sort contacts"
                    >
                      <option value="createdAt">Date</option>
                      <option value="name">Name</option>
                      <option value="email">Email</option>
                    </select>
                    <button
                      onClick={() => setSortDir(sortDir === 'asc' ? 'desc' : 'asc')}
                      className="btn-ghost text-sm px-2 py-1"
                      aria-pressed={sortDir === 'asc'}
                    >
                      {sortDir === 'asc' ? '↑ Asc' : '↓ Desc'}
                    </button>
                  </div>
                  <div className="ml-auto flex items-center gap-2">
                    <label className="text-sm font-medium">Per page:</label>
                    <select
                      value={pageSize}
                      onChange={(e) => { setPageSize(Number(e.target.value)); setPage(1); }}
                      className="input px-2 py-1 text-sm"
                      aria-label="Contacts per page"
                    >
                      <option value={5}>5</option>
                      <option value={10}>10</option>
                      <option value={25}>25</option>
                    </select>
                  </div>
                </div>

                {/* Contacts List */}
                {totalContacts === 0 ? (
                  <p className="text-sm text-[var(--text-secondary)]">No contact submissions found.</p>
                ) : (
                  <>
                    <div className="space-y-3">
                      {paginatedContacts.map((c, idx) => (
                        <div key={idx} className="rounded-lg border border-[var(--border)] p-4">
                          <div className="flex items-center justify-between">
                            <div className="flex-1">
                              <p className="font-medium">{c.name} <span className="text-xs text-[var(--text-secondary)]">• {c.email}</span></p>
                              <p className="text-xs text-[var(--text-secondary)]">{c.createdAt ? new Date(c.createdAt).toLocaleString() : ''}</p>
                            </div>
                          </div>
                          <div className="mt-3 text-sm text-[var(--text-primary)]">{c.message}</div>
                        </div>
                      ))}
                    </div>

                    {/* Pagination Info & Controls */}
                    <div className="mt-6 flex items-center justify-between border-t border-[var(--border)] pt-4">
                      <div className="text-sm text-[var(--text-secondary)]">
                        Showing {totalContacts === 0 ? 0 : (page - 1) * pageSize + 1} – {Math.min(page * pageSize, totalContacts)} of {totalContacts}
                      </div>
                      <div className="flex items-center gap-2">
                        <button
                          onClick={() => setPage((p) => Math.max(1, p - 1))}
                          disabled={page === 1}
                          className="btn-ghost text-sm px-3 py-1 disabled:opacity-50"
                        >
                          ← Prev
                        </button>
                        <div className="text-sm font-medium">Page {page} of {maxPage}</div>
                        <button
                          onClick={() => setPage((p) => Math.min(maxPage, p + 1))}
                          disabled={page >= maxPage}
                          className="btn-ghost text-sm px-3 py-1 disabled:opacity-50"
                        >
                          Next →
                        </button>
                      </div>
                    </div>
                  </>
                )}
              </div>
            )}
          </div>
        )}
      </div>
    </section>
  );
}