# SmartHome Hub - Monorepo Single Executable Setup

## Quick Start

This guide walks you through creating a single executable that combines your Go backend and Angular frontend with auto-update capability.

## Current Project Structure

```
SmartHome.Hub/
├── backend/                    # Go service
│   ├── cmd/launcher/           # New: Launcher entry point
│   ├── internal/config/        # New: Configuration management
│   ├── pkg/hub/                
│   │   ├── static.go          # New: Frontend serving
│   │   ├── updater.go         # New: Update checking
│   │   └── frontend/          # New: Embedded Angular files go here
│   ├── .goreleaser.yaml       # Updated: Build config
│   └── Makefile               # New: Build helpers
├── frontend/                   # Angular app (create if needed)
│   ├── src/
│   ├── dist/                  # Build output
│   └── package.json
└── DEPLOYMENT.md              # New: Full deployment guide
```

## Step 1: Prepare Your Angular Frontend

### Option A: Existing Angular Project

If you already have an Angular frontend:

```bash
# Copy it to workspace root
cp -r /path/to/your/frontend ./frontend
```

### Option B: Create New Angular Project

```bash
cd SmartHome.Hub
ng new frontend --skip-git --package-manager=npm
cd frontend
```

### Build the Frontend

```bash
cd SmartHome.Hub/frontend
npm install
npm run build
```

## Step 2: Embed Frontend in Go Binary

After building the Angular frontend:

```bash
cd SmartHome.Hub/backend

# Create frontend directory
mkdir -p pkg/hub/frontend

# Copy built frontend files
cp -r ../frontend/dist/* pkg/hub/frontend/

# Verify files are there
ls pkg/hub/frontend/
```

## Step 3: Build the Single Executable

### Using Make (Recommended)

```bash
cd backend

# Build for your current platform
make build
# Output: ./build/smarthome-hub

# Build for all platforms
make build-all
# Output: ./build/{linux,windows,macos}/smarthome-hub*

# Or use specific targets
make build-linux
make build-windows
make build-macos
```

### Manual Build

```bash
cd backend

# Build for Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build \
  -ldflags="-X main.version=1.0.0" \
  -o build/smarthome-hub-linux \
  ./cmd/launcher

# Build for macOS
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build \
  -ldflags="-X main.version=1.0.0" \
  -o build/smarthome-hub-macos \
  ./cmd/launcher

# Build for Windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build \
  -ldflags="-X main.version=1.0.0" \
  -o build/smarthome-hub.exe \
  ./cmd/launcher
```

## Step 4: Run the Application

```bash
# Navigate to where the executable is
cd backend/build

# Run it
./smarthome-hub

# You should see:
# - Update check
# - Server starting on :8080
# - Open http://localhost:8080 in browser
```

## Architecture

### What Happens When You Run It

1. **Initialization**
   ```
   Executable starts
   ├─ Load/Create config.json
   ├─ Initialize logger
   ├─ Check for updates
   ├─ Apply updates if available
   └─ Start API server
   
   Server starts
   ├─ Initialize database
   ├─ Set up authentication
   ├─ Register API routes
   └─ Serve static frontend
   ```

2. **Port 8080 (Single Port)**
   ```
   /                    → Angular frontend (static files)
   /api/*              → Go backend API endpoints
   /swagger/*          → API documentation
   /auth/*             → Authentication endpoints
   ```

### Single Executable Contains

```
smarthome-hub (executable)
│
├─ Go Runtime
├─ Embedded Frontend Files (Angular dist/)
├─ SQLite Database (created on demand)
└─ Configuration (config.json)
```

## Configuration

On first run, `config.json` is created with defaults:

```json
{
  "addr": ":8080",
  "databasePath": "./data/hub.db",
  "updateOwner": "Project-SmartHome",
  "updateRepo": "SmartHome.Hub",
  "enableUpdateCheck": true,
  "updateAutoApply": true,
  "currentVersion": "1.0.0",
  "updateGitHubToken": "",
  "appEnv": "production",
  "debugMode": false,
  "frontendFS": true
}
```

**Key Options:**
- `addr`: Server address (change to `:3000` to use port 3000)
- `updateAutoApply`: `true` to auto-update, `false` to require manual update
- `updateOwner/updateRepo`: GitHub repo for checking updates
- `debugMode`: `true` for verbose logging

## Automatic Updates

### How It Works

1. On startup, checks GitHub releases for new versions
2. If found and `updateAutoApply: true`, downloads changed files only (differential)
3. Applies updates and prompts restart

### Setting Up GitHub Releases

1. **Create Release** on GitHub with tag `v1.0.0`
2. **Build your app** with version flag: `make VERSION=1.0.0 build-all`
3. **Generate manifest** (see below)
4. **Upload to release**:
   - Executable files
   - `update-manifest.json` file

### Generate Update Manifest

```bash
# After building, generate SHA256 hashes
cd backend/build

# Generate manifest
cat > update-manifest.json << 'EOF'
{
  "version": "1.0.0",
  "baseVersion": "0.9.0",
  "baseURL": "https://github.com/Project-SmartHome/SmartHome.Hub/releases/download/v1.0.0/",
  "files": [
    {
      "path": "smarthome-hub",
      "sha256": "$(shasum -a 256 linux/smarthome-hub-amd64 | cut -d' ' -f1)",
      "size": $(stat -f%z linux/smarthome-hub-amd64 2>/dev/null || stat -c%s linux/smarthome-hub-amd64),
      "mode": 493,
      "goos": "linux",
      "goarch": "amd64"
    }
  ]
}
EOF
```

Or use GitHub Actions (configured in `.github/workflows/build-release.yml`):

```bash
# GitHub Actions will:
# 1. Build frontend
# 2. Embed in binary
# 3. Build for all platforms
# 4. Generate hashes
# 5. Upload everything
```

## Deployment Options

### Docker

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY backend .
RUN go build -o smarthome-hub ./cmd/launcher

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/smarthome-hub .
COPY config.json .
EXPOSE 8080
CMD ["./smarthome-hub"]
```

```bash
docker build -t smarthome-hub .
docker run -p 8080:8080 smarthome-hub
```

### Linux Service (systemd)

```ini
# /etc/systemd/system/smarthome-hub.service
[Unit]
Description=SmartHome Hub Service
After=network.target

[Service]
Type=simple
User=smarthome
WorkingDirectory=/opt/smarthome-hub
ExecStart=/opt/smarthome-hub/smarthome-hub
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable smarthome-hub
sudo systemctl start smarthome-hub
sudo systemctl status smarthome-hub
```

### Windows Task Scheduler

1. Save executable to `C:\Program Files\SmartHome.Hub\smarthome-hub.exe`
2. Create task: Point to executable
3. Set to run at startup

### Cloud (AWS, Azure, Google Cloud)

1. Build for Linux: `make build-linux`
2. Upload to cloud storage
3. Deploy to VM/Container service
4. Update config with your domain

## Troubleshooting

### Build Fails - "CGO_ENABLED=1"

**Issue**: Cross-platform builds need CGO dependencies

**Solution**:
```bash
# On macOS for Linux:
brew install mingw-w64

# On Linux for Windows:
apt-get install mingw-w64

# Or disable CGO (not recommended for SQLite):
GOOS=linux CGO_ENABLED=0 go build ...
```

### Frontend Not Loading

**Issue**: Blank page or 404 on http://localhost:8080

**Solution**:
1. Verify frontend files exist: `ls backend/pkg/hub/frontend/`
2. Check if Angular built correctly: `ls frontend/dist/`
3. Rebuild: `make clean && make build`

### Update Check Fails

**Issue**: "Failed to check for updates"

**Solution**:
1. Check config.json has correct owner/repo
2. Verify internet connection
3. Check GitHub API rate limits: `curl -I https://api.github.com`
4. Use token if private: set `updateGitHubToken` in config.json

### Database Issues

**Issue**: "permission denied" on hub.db

**Solution**:
```bash
# Ensure data directory is writable
mkdir -p data
chmod 755 data

# Or change path in config.json to writable location
```

## Development Workflow

### Local Development

```bash
# Terminal 1: Backend API
cd backend
make dev

# Terminal 2: Frontend (if doing frontend work)
cd frontend
npm start
```

Set API proxy in `frontend/angular.json`:
```json
"serve": {
  "options": {
    "proxyConfig": "proxy.conf.json"
  }
}
```

### Building for Release

```bash
cd backend

# 1. Update version in code or use VERSION
make VERSION=1.0.1 build-all

# 2. Generate and upload manifest to GitHub
# 3. Create GitHub release with binaries
```

## Environment Variables

```bash
# Override database path
export DATABASE_PATH=/custom/path/hub.db

# Enable debug logging
export DEBUG_MODE=true

# Set GitHub token for updates (private repos)
export GITHUB_TOKEN=your_token_here
```

## Performance Tips

1. **Database**: Consider database indexing for large datasets
2. **Frontend**: Angular production build is already optimized
3. **Updates**: Differential updates download only changed files
4. **Resources**: Single executable is ~50MB typical size

## Security Considerations

1. **Secrets**: Never commit secrets to Git
   - Use environment variables
   - Store GitHub token separately
   
2. **HTTPS**: For production, use reverse proxy:
   - Nginx with SSL
   - Caddy (automatic HTTPS)
   
3. **Database**: 
   - Regular backups of `hub.db`
   - Restrict file permissions: `chmod 600 hub.db`

4. **Updates**: 
   - Verify checksums are downloaded over HTTPS
   - Keep GitHub token secure (use Actions secrets)

## Next Steps

1. ✅ Create frontend directory or copy existing
2. ✅ Build frontend: `npm run build`
3. ✅ Embed in binary: Copy to `backend/pkg/hub/frontend/`
4. ✅ Build executable: `make build`
5. ✅ Test locally: `./build/smarthome-hub`
6. ✅ Set up GitHub releases for auto-updates
7. ✅ Deploy to your server/cloud

## Support

- Check logs: Look for output in terminal
- Debug mode: Set `debugMode: true` in config.json
- Rebuild: `make clean && make build`
- Check DEPLOYMENT.md for more detailed guide

---

For full deployment guide, see [DEPLOYMENT.md](DEPLOYMENT.md)
