# Node.js Runtime Setup Guide

## Minimum Requirements
- **Node.js**: ≥18.17.0 (engine requirement in `package.json`)
- **npm**: ≥9.0 (comes with Node)

---

## Installation Options

### Option 1: System-Wide Installation (Recommended for Production)

#### Windows (MSI Installer)
1. Download LTS installer from [nodejs.org](https://nodejs.org/)
2. Run the `.msi` installer as Administrator
3. Accept default installation path (`C:\Program Files\nodejs`)
4. Verify installation:
   ```bash
   node --version    # should show v20.x.x or higher
   npm --version     # should show v10.x.x or higher
   ```

#### macOS (Homebrew)
```bash
brew install node@20
```

#### Linux (apt, yum, etc.)
```bash
# Ubuntu/Debian
sudo apt update && sudo apt install nodejs npm

# CentOS/RHEL
sudo yum install nodejs npm
```

### Option 2: Portable Node (For Development Without Admin Rights)

Used when system Node is incompatible or admin installation unavailable.

#### Windows (Portable Extraction)

1. Download portable Node.js from [nodejs.org](https://nodejs.org/en/download/)
   - Choose the `.zip` file (not `.msi`)
   
2. Extract to user-accessible directory (e.g., `C:\Users\YourUsername\tools\node-v20.20.2`)

3. Add to PATH temporarily (for current terminal session):
   ```powershell
   $env:PATH = "C:\Users\YourUsername\tools\node-v20.20.2;$env:PATH"
   node --version  # verify
   ```

4. Or add permanently via Environment Variables:
   - Press `Win+X` → `System` → `Advanced System Settings`
   - Click `Environment Variables`
   - Under `User Variables`, select `Path` → `Edit`
   - Add: `C:\Users\YourUsername\tools\node-v20.20.2`
   - Click OK and restart terminal/IDE

### Option 3: Version Manager (For Multiple Node Versions)

#### Windows + nvm-windows
```powershell
# Install nvm-windows (one-time setup)
# https://github.com/coreybutler/nvm-windows/releases
# Download and run the installer

# Then use:
nvm list available         # see available versions
nvm install 20.20.2        # install specific version
nvm use 20.20.2            # switch to version
nvm default 20.20.2        # set as default
node --version             # verify
```

#### macOS/Linux + nvm
```bash
# Install nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc  # or ~/.zshrc

# Use
nvm install 20.20.2
nvm use 20.20.2
nvm alias default 20.20.2  # set as default
```

---

## Verifying Installation

```bash
# Check Node version (should be ≥18.17.0)
node --version

# Check npm version (should be ≥9.0)
npm --version

# Check that npm can find Node
npm config get prefix
```

---

## .nvmrc File

The repository includes `.nvmrc` specifying the recommended Node version:
```
20.20.2
```

**Using nvm** to auto-switch:
```bash
cd /path/to/dev-folio
nvm use    # switches to 20.20.2 if installed
```

---

## Development Server Setup

### Starting the Dev Server

#### Using System Node (if ≥18.17.0)
```bash
cd /path/to/dev-folio
npm install
npm run dev
# Opens http://localhost:3000
```

#### Using Portable Node (Windows Example)
```powershell
# Open PowerShell or cmd

# Option A: Set PATH for this session only
$env:PATH = "C:\Users\YourUsername\tools\node-v20.20.2;$env:PATH"
node --version  # verify

# Then run dev
cd C:\Users\gigli\GoWs\dev-folio
npm run dev
```

#### Using nvm
```bash
nvm use 20.20.2
npm run dev
```

### Dev Server Health Check
```bash
# Should return 200 and page content
curl http://localhost:3000

# Or open in browser
# http://localhost:3000
```

---

## Production Deployment

### Node.js Runtime Environment

#### Environment Variables
```bash
# Set in .env or deployment platform
export NODE_ENV=production
export NEXT_PUBLIC_API_URL=https://api.your-domain.com
export DEVFOLIO_JWT_SECRET=your-secret-here
```

#### Starting Production Server
```bash
npm install --production   # install only prod dependencies
npm run build              # build optimized bundle
npm start                  # start prod server on PORT 3000 (or custom)
```

#### Custom Port
```bash
PORT=8000 npm start
# or
NODE_OPTIONS="--port=8000" npm start
```

### Docker (Recommended)
```dockerfile
FROM node:20-alpine

WORKDIR /app
COPY package*.json ./
RUN npm ci --omit=dev

COPY . .
RUN npm run build

ENV NODE_ENV=production
EXPOSE 3000

CMD ["npm", "start"]
```

Run:
```bash
docker build -t devfolio-frontend .
docker run -e NODE_ENV=production -e DEVFOLIO_JWT_SECRET=secret -p 3000:3000 devfolio-frontend
```

---

## Troubleshooting

### Node Command Not Found
```bash
# Verify PATH
echo $PATH  # on Unix/macOS/Linux
echo %PATH% # on Windows

# If Node path missing, add it (see Option 2 above)

# Verify installation
where node      # Windows
which node      # Unix/macOS/Linux
```

### npm ERR! code ERESOLVE
```bash
# Common with Node 18+
# Solution 1: Clear cache
npm cache clean --force

# Solution 2: Use legacy dependency resolution
npm install --legacy-peer-deps

# Solution 3: Upgrade npm
npm install -g npm@latest
```

### Port 3000 Already in Use
```bash
# Find process using port 3000
lsof -i :3000                # Unix/macOS/Linux
netstat -ano | findstr 3000  # Windows

# Use different port
PORT=3001 npm run dev
```

### Slow npm Install
```bash
# Use npm mirror
npm config set registry https://registry.npmmirror.com

# Or use yarn/pnpm (faster alternatives)
yarn install
```

---

## Performance Tips

### Reduce Node Memory Usage
```bash
node --max-old-space-size=512 node_modules/.bin/next build
# Limits heap to 512MB (adjust as needed)
```

### Enable Node Clustering (Production)
Use PM2 or similar:
```bash
npm install -g pm2
pm2 start npm --name "devfolio" -- start
pm2 save      # persist across reboots
```

### Verify Node Performance
```bash
node -v && npm -v           # Show versions
node --version              # Alternative
npm run build --timing      # Show build timing
```

---

## Updating Node.js

### System Installation
```bash
# Windows: Download and run new .msi installer
# macOS: brew upgrade node
# Linux: sudo apt update && sudo apt upgrade nodejs
```

### Portable Installation
1. Download new portable version
2. Extract to new folder
3. Update PATH or nvm alias

### After Update
```bash
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
```

---

## Additional Resources
- [Node.js Official Docs](https://nodejs.org/docs/)
- [npm Documentation](https://docs.npmjs.com/)
- [nvm-windows](https://github.com/coreybutler/nvm-windows)
- [Next.js Deployment](https://nextjs.org/docs/deployment)
