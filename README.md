# DevFolio

DevFolio is a full-stack developer portfolio platform built with **Next.js 14** (React 18, TypeScript, Tailwind CSS) and a **Go REST API**. It includes a complete admin dashboard, contact form with rate limiting, media uploads, and persistent storage (in-memory or MongoDB).

## Features

✅ **Frontend**
- Server & client components (App Router)
- Admin dashboard with Projects, Posts, and Contacts tabs
- Media upload with image validation
- Contact form with rate limiting (5/email/hour)
- Client-side search & filters for Projects and Blog
- SEO metadata & JSON-LD structured data
- Light/dark theme with persistent preference
- Fully typed with TypeScript

✅ **Backend (Go)**
- JWT-based admin authentication
- RESTful API (`/api/v1/*`)
- CRUD operations for projects, posts, skills
- Contact submissions management
- Rate limiting per email
- In-memory store with optional MongoDB persistence
- Comprehensive test coverage

✅ **Admin Console**
- Projects: Create, edit, delete with image uploads
- Posts: Full blog post management with tags & metadata
- Contacts: View, sort, and paginate incoming submissions
- Authentication: JWT token-based access control

---

## Quick Start

### Prerequisites
- **Node.js**: ≥18.17.0 (see [NODE_SETUP.md](NODE_SETUP.md) for installation help)
- **Go**: ≥1.22 (optional; required only for backend development)
- **MongoDB**: Optional; in-memory store used by default

### Development

```bash
# Install dependencies
npm install

# Start dev server (http://localhost:3000)
npm run dev

# In another terminal, start Go backend (http://localhost:8080)
npm run backend:dev

# Type check
npm run typecheck

# Run tests (backend)
cd backend && go test ./...
```

### Production Build

See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed deployment instructions (Docker, Vercel, Kubernetes, etc.).

```bash
npm run build
npm start
```

---

## Architecture

### Frontend Structure
```
app/
├── (frontend)/           # Public pages (Home, Blog, Projects, Contact)
├── (payload)/admin/      # Admin console
├── layout.tsx            # Root layout with metadata & theme
└── ...

lib/
├── backend.ts            # API client functions
├── content.ts            # Type definitions
└── ...
```

### Backend Structure
```
backend/cmd/devfolio-api/
├── main.go               # Entry point
└── internal/
    ├── server/           # HTTP handlers & routing
    ├── store/            # Repository pattern (in-memory, MongoDB)
    ├── domain/           # Domain models
    └── data/             # Seed data
```

---

## Environment Setup

Create `.env.local` in the root:

```bash
# JWT secret for admin auth (generate: openssl rand -hex 32)
DEVFOLIO_JWT_SECRET=your-32-char-secret-here

# MongoDB (optional)
MONGODB_URI=mongodb+srv://user:pass@cluster.mongodb.net/devfolio

# API base URL for frontend
NEXT_PUBLIC_API_URL=http://localhost:8080

# Admin credentials for login UI
DEVFOLIO_ADMIN_USER=admin
DEVFOLIO_ADMIN_PASS=admin-password
```

---

## API Endpoints

### Public Endpoints
- `GET /api/v1/health` — Health check
- `GET /api/v1/site` — Site metadata
- `GET /api/v1/projects` — List projects
- `GET /api/v1/posts` — List blog posts
- `GET /api/v1/skills` — List skills
- `POST /api/v1/contact` — Submit contact form (rate-limited)

### Admin Endpoints (Protected)
- `POST /api/v1/login` — Admin login → JWT token
- `GET /api/v1/projects` — Manage projects (CRUD)
- `POST /api/v1/posts` — Create post
- `PUT /api/v1/projects/:slug` — Update project
- `DELETE /api/v1/projects/:slug` — Delete project
- `POST /api/v1/upload` — Upload media file
- `GET /api/v1/contact` — List contact submissions

---

## Admin Login & Usage

1. Navigate to `/admin`
2. Enter credentials (default: `admin` / `admin-password`)
3. Manage projects, posts, and contact submissions
4. Upload project/post images directly from the UI

### JWT Token Flow
```bash
# 1. Login to get token
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin-password"}'
# Response: {"token": "eyJ..."}

# 2. Use token in subsequent requests
curl http://localhost:8080/api/v1/projects \
  -H "Authorization: Bearer eyJ..."
```

---

## Features in Detail

### Media Upload
- POST `/api/v1/upload` with multipart form data
- Accepts images (image/* MIME types)
- Max 10MB per file
- Returns: `{url: "/uploads/timestamp-filename.ext"}`
- Accessible via `/uploads/` route

### Rate Limiting
- Contact form: 5 submissions per email per hour
- Enforced server-side in `backend/internal/server/ratelimit.go`
- Returns `429 Too Many Requests` if exceeded

### Sorting & Pagination (Admin)
- **Contacts Tab**: Sort by Date, Name, or Email (asc/desc)
- **Pagination**: 5, 10, or 25 items per page
- **Display**: "Showing 1-10 of 47" indicator

### SEO & Metadata
- Per-page metadata via Next.js Metadata API
- JSON-LD WebSite schema on homepage
- OpenGraph tags for social sharing
- Auto-generated sitemap & robots.txt

---

## Testing

### Frontend
```bash
npm run typecheck          # TypeScript validation
```

### Backend
```bash
cd backend
go test ./...              # Run all tests
go test -v ./internal/server -run TestContact*   # Run specific tests
go test -cover ./...       # Coverage report
```

Current test coverage:
- ✅ Contact handler tests (POST, GET, error cases)
- ✅ Rate limiter tests
- ✅ Admin authentication tests

---

## Troubleshooting

### Node.js Issues
- **"not found"**: Verify Node ≥18.17.0 (`node --version`)
- **MSI admin blocker**: Use portable Node or nvm (see [NODE_SETUP.md](NODE_SETUP.md))
- **Port 3000 in use**: `PORT=3001 npm run dev`

### Backend Issues
- **"connection refused"**: Go backend not running; start with `npm run backend:dev`
- **JWT errors**: Verify `DEVFOLIO_JWT_SECRET` is set in `.env.local`
- **Build failures**: Run `go mod tidy && go build ./cmd/devfolio-api`

### TypeScript Errors
```bash
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
npm run typecheck
```

---

## Performance Notes

- Frontend: Deployed on Vercel or self-hosted Node
- Backend: Stateless; horizontally scalable (use Redis for shared rate limits at scale)
- Database: In-memory fine for ≤10K submissions; MongoDB for production scale
- Caching: Implement Redis layer for frequently-accessed projects/posts

---

## Production Checklist

- [ ] Set secure `DEVFOLIO_JWT_SECRET` (32+ random chars)
- [ ] Enable HTTPS for frontend and API
- [ ] Configure MongoDB with authentication & IP whitelist
- [ ] Set `NEXT_PUBLIC_API_URL` to production domain
- [ ] Enable CORS headers correctly
- [ ] Run `npm audit` and `go mod tidy` for vulnerabilities
- [ ] Test rate limiting and contact form flow
- [ ] Set up monitoring & logging
- [ ] Enable database backups
- [ ] Review security checklist in [DEPLOYMENT.md](DEPLOYMENT.md)

---

## Documentation

- **[DEPLOYMENT.md](DEPLOYMENT.md)** — Comprehensive deployment guide (Docker, Vercel, K8s)
- **[NODE_SETUP.md](NODE_SETUP.md)** — Node.js installation and runtime guide
- **[Backend Tests](backend/internal/server/*_test.go)** — Reference implementation examples

---

## License

MIT