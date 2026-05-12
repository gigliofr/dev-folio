# Production Build & Deployment Guide

## Overview
dev-folio is a full-stack portfolio platform with:
- **Backend**: Go REST API (stdlib HTTP)
- **Frontend**: Next.js 14 (React 18, App Router, TypeScript)
- **Database**: Optional MongoDB persistence; in-memory store by default
- **Auth**: JWT-based admin authentication

---

## Prerequisites

### System Requirements
- **Node.js**: ≥18.17.0 (strongly recommend 20.20.2 LTS or later)
- **Go**: ≥1.22
- **MongoDB** (optional): If using persistent storage

### Environment Variables
Create a `.env.local` file in the root directory:
```bash
# JWT secret for admin token signing/validation
DEVFOLIO_JWT_SECRET=your-secure-secret-here-min-32-chars

# MongoDB connection string (optional; omit for in-memory store)
MONGODB_URI=mongodb+srv://user:pass@cluster.mongodb.net/devfolio

# Admin credentials (optional; used for seeding or legacy token auth)
DEVFOLIO_ADMIN_TOKEN=optional-legacy-token

# Frontend API base (for dev/production)
NEXT_PUBLIC_API_URL=https://api.example.com
```

---

## Frontend Build & Deployment

### Local Development
```bash
# Install dependencies
npm install

# Start dev server (uses portable Node or system Node ≥18.17.0)
npm run dev
# Opens at http://localhost:3000
```

### Production Build
```bash
# Install dependencies
npm install --production

# Build for production
npm run build

# Validate TypeScript
npm run typecheck

# Start production server
npm start
# Listens on http://localhost:3000 (or PORT env var)
```

### Deployment to Vercel / Node Hosting
1. **Push to Git** and connect repository to Vercel
2. **Build Command**: `npm run build`
3. **Start Command**: `npm start`
4. **Environment**: Set `DEVFOLIO_JWT_SECRET`, `NEXT_PUBLIC_API_URL` in Vercel dashboard
5. Deploy automatically on push

### Deployment to Docker
```dockerfile
FROM node:20-alpine as builder
WORKDIR /app
COPY . .
RUN npm install && npm run build

FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json
ENV NODE_ENV=production
EXPOSE 3000
CMD ["npm", "start"]
```

Build and run:
```bash
docker build -t devfolio-frontend .
docker run -e DEVFOLIO_JWT_SECRET=your-secret -p 3000:3000 devfolio-frontend
```

---

## Backend Build & Deployment

### Local Development
```bash
cd backend

# Build the API
go build ./cmd/devfolio-api

# Run the API
./devfolio-api
# Listens on http://localhost:8080
```

### Production Build
```bash
cd backend

# Build optimized binary
go build -ldflags="-s -w" ./cmd/devfolio-api

# Run with production environment
DEVFOLIO_JWT_SECRET=your-secret-key PORT=8080 ./devfolio-api
```

### Docker Deployment
```dockerfile
FROM golang:1.22-alpine as builder
WORKDIR /app
COPY backend/ .
RUN go build -ldflags="-s -w" ./cmd/devfolio-api

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/devfolio-api .
COPY --from=builder /app/internal/data/seed.json ./internal/data/seed.json
EXPOSE 8080
CMD ["./devfolio-api"]
```

Build and run:
```bash
docker build -f Dockerfile.backend -t devfolio-backend .
docker run -e DEVFOLIO_JWT_SECRET=your-secret -e MONGODB_URI=mongodb://... -p 8080:8080 devfolio-backend
```

### Kubernetes Deployment (Optional)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: devfolio-backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: devfolio-backend
  template:
    metadata:
      labels:
        app: devfolio-backend
    spec:
      containers:
      - name: api
        image: your-registry/devfolio-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: DEVFOLIO_JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: devfolio-secrets
              key: jwt-secret
        - name: MONGODB_URI
          valueFrom:
            secretKeyRef:
              name: devfolio-secrets
              key: mongodb-uri
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

---

## API Endpoints

### Health & Stats
- `GET /health` — API health check
- `GET /api/v1/stats` — Site statistics (public)

### Content (Public)
- `GET /api/v1/site` — Site metadata
- `GET /api/v1/projects` — List projects
- `GET /api/v1/posts` — List blog posts
- `GET /api/v1/skills` — List skills

### Admin (Protected by JWT)
- `POST /api/v1/login` — Admin login → returns JWT token
- `GET /api/v1/projects` (admin) — Manage projects
- `POST /api/v1/projects` — Create project
- `PUT /api/v1/projects/:slug` — Update project
- `DELETE /api/v1/projects/:slug` — Delete project
- `GET /api/v1/posts` (admin) — Manage posts
- `POST /api/v1/posts` — Create post
- `PUT /api/v1/posts/:slug` — Update post
- `DELETE /api/v1/posts/:slug` — Delete post
- `POST /api/v1/upload` — Upload media (returns `{url}`)
- `GET /api/v1/contact` (admin) — List contact submissions
- `POST /api/v1/contact` (public) — Submit contact form

### Admin Authentication
Include JWT token in `Authorization` header:
```bash
curl -H "Authorization: Bearer your-jwt-token" https://api.example.com/api/v1/projects
```

---

## Testing

### Frontend Tests
```bash
# TypeScript type checking
npm run typecheck

# Run Jest/Vitest (if configured)
npm run test
```

### Backend Tests
```bash
cd backend

# Run all tests
go test ./...

# Run specific test suite
go test ./internal/server -v -run TestContactSubmission

# With coverage
go test -cover ./...
```

---

## Performance Optimization

### Frontend
- **Image Optimization**: Next.js automatically optimizes images; use `<Image>` component
- **Code Splitting**: App Router enables automatic code splitting per route
- **Caching**: Set cache headers in `next.config.js`:
  ```js
  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          { key: 'Cache-Control', value: 'public, max-age=3600' }
        ]
      }
    ]
  }
  ```

### Backend
- **Connection Pooling**: MongoDB driver handles pooling automatically
- **Rate Limiting**: Contact form limited to 5 submissions/email/hour (configurable in `ratelimit.go`)
- **CORS**: Enabled for all origins in development; configure for production

---

## Monitoring & Logging

### Health Checks
```bash
# Frontend
curl https://your-domain.com/

# Backend
curl https://api.your-domain.com/health
```

### Logs
- **Frontend**: Browser console (`window.console`), server logs via Node
- **Backend**: Stdout/stderr; integrate with ELK/Datadog as needed

---

## Security Checklist

- [ ] Set strong `DEVFOLIO_JWT_SECRET` (32+ chars, random)
- [ ] Use HTTPS in production for all frontend and API endpoints
- [ ] Set CORS headers correctly in backend
- [ ] Validate all user inputs (contact form, admin submissions)
- [ ] Rate limit contact submissions (already implemented: 5/hour/email)
- [ ] Use environment variables; never commit secrets
- [ ] Enable MongoDB authentication and IP whitelisting (if used)
- [ ] Run `npm audit` and `go mod tidy` to check for vulnerabilities
- [ ] Enable HSTS header for frontend
- [ ] Use Content Security Policy headers

---

## Scaling Considerations

### Horizontal Scaling
- **Frontend**: Stateless; deploy multiple replicas behind a load balancer
- **Backend**: Stateless (rate limiter is in-memory per instance); use Redis for shared rate limits at scale

### Database Scaling
- **In-Memory**: Single instance; fine for ≤10K submissions
- **MongoDB**: Horizontal scaling via sharding; enable replica sets for HA

### Caching Layer
- Add Redis or Memcached for frequently-accessed data (projects, posts)
- Cache contact list retrieval with TTL

---

## Troubleshooting

### Node.js Version Issues
- Verify: `node --version` must be ≥18.17.0
- Use `.nvmrc` file (included) with nvm:
  ```bash
  nvm use  # switches to 20.20.2
  npm run dev
  ```

### Build Errors
```bash
# Clear cache and reinstall
rm -rf node_modules .next && npm install && npm run build

# Verify types
npm run typecheck
```

### API Errors
```bash
# Check if backend is running
curl http://localhost:8080/health

# Verify JWT secret is set
echo $DEVFOLIO_JWT_SECRET  # should not be empty
```

### MongoDB Connection
```bash
# Test connection
mongosh "mongodb+srv://user:pass@cluster.mongodb.net/devfolio"
```

---

## Support & Maintenance

- **Dependencies**: Run `npm update` and `go get -u ./...` quarterly
- **Security Updates**: Monitor `npm audit` and `go mod tidy` output
- **Backups**: Regular MongoDB backups if using persistent storage
- **Monitoring**: Set up uptime monitoring for API and frontend health endpoints
