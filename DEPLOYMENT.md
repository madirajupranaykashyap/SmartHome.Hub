# SmartHome Hub - Single Executable Deployment Guide

## Overview

The SmartHome Hub now builds as a single, self-contained executable that includes:
- ✅ Go backend service with API
- ✅ Embedded Angular frontend (static assets)
- ✅ Automatic update checking and installation
- ✅ Single configuration file (`config.json`)

## Architecture

```
smarthome-hub (executable)
├── Go Backend (runs on :8080)
│   ├── API endpoints (/api/*)
│   ├── Authentication
│   ├── Database
│   └── Update checker
└── Angular Frontend (embedded assets)
    ├── Served as static files
    └── Loads from / path
```

## Building the Single Executable

### Prerequisites

- Go 1.25+
- (Optional) Node.js 18+ if building Angular frontend

### Step 1: Prepare Angular Frontend

If you have an Angular frontend, place it in the workspace root:

```bash
# Create frontend directory at workspace root
mkdir -p ../frontend
cd ../frontend

# Initialize or copy your Angular project
ng new . --skip-git
# or copy existing project files
```

Build the frontend:

```bash
cd ../frontend
npm install
npm run build
# Output should be in dist/your-app-name
```

### Step 2: Embed Frontend in Go Binary

Copy the built frontend to the `backend/pkg/hub/frontend` directory:

```bash
# From backend directory
mkdir -p pkg/hub/frontend
cp -r ../../frontend/dist/* pkg/hub/frontend/
```

### Step 3: Build the Executable

```bash
# Navigate to backend directory
cd backend

# Build for current platform
make build

# Or build for all platforms
make build-all

# Binary will be in ./build directory
```

## Configuration

The application creates a `config.json` file on first run:

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

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `addr` | Server address and port | `:8080` |
| `databasePath` | Path to SQLite database | `./data/hub.db` |
| `updateOwner` | GitHub owner for updates | `Project-SmartHome` |
| `updateRepo` | GitHub repository for updates | `SmartHome.Hub` |
| `enableUpdateCheck` | Enable update checking | `true` |
| `updateAutoApply` | Automatically apply updates | `true` |
| `currentVersion` | Current version (auto-updated) | `1.0.0` |
| `updateGitHubToken` | GitHub token (for private repos) | `` |
| `appEnv` | Environment (production/development) | `production` |
| `debugMode` | Enable debug logging | `false` |
| `frontendFS` | Serve embedded frontend | `true` |

## Running the Application

### Standalone

```bash
./build/smarthome-hub
```

The application will:
1. Check for updates (configurable)
2. Apply auto-updates if available
3. Start the Go backend server
4. Serve Angular frontend at `http://localhost:8080`

### With Custom Configuration

```bash
# Set config path via environment variable
export CONFIG_PATH="./config.json"
./build/smarthome-hub
```

### Development Mode

```bash
make dev
```

## Update Mechanism

### How It Works

1. **Check**: On startup, checks GitHub releases for `update-manifest.json`
2. **Download**: Downloads only files that have changed (differential updates)
3. **Apply**: Extracts and replaces files in the installation directory
4. **Restart**: Prompts user to restart application

### GitHub Release Setup

Create an `update-manifest.json` in your GitHub release:

```json
{
  "version": "1.0.1",
  "baseVersion": "1.0.0",
  "baseURL": "https://github.com/Project-SmartHome/SmartHome.Hub/releases/download/v1.0.1/",
  "files": [
    {
      "path": "smarthome-hub",
      "sha256": "hash_here",
      "size": 12345678,
      "mode": 493
    }
  ],
  "delete": []
}
```

## Docker Deployment

### Dockerfile

```dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /build
COPY backend .

RUN go build -ldflags="-s -w -X main.version=docker" \
    -o smarthome-hub ./cmd/launcher

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/smarthome-hub .
COPY config.json .

EXPOSE 8080

CMD ["./smarthome-hub"]
```

Build and run:

```bash
docker build -t smarthome-hub .
docker run -p 8080:8080 -v /path/to/data:/app/data smarthome-hub
```

## Nginx Reverse Proxy (Optional)

If running behind Nginx:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Troubleshooting

### Update Check Fails
- Check internet connection
- Verify GitHub repository name in config
- Check rate limits: `curl -I https://api.github.com`

### Frontend Not Loading
- Verify `pkg/hub/frontend` contains built Angular files
- Check web browser console for 404 errors
- Ensure frontend assets were embedded correctly

### Database Permission Issues
- Ensure `./data` directory is writable
- On Linux: `chmod 755 ./data`

### Build Issues
```bash
# Clean and rebuild
make clean
make build

# Check Go version
go version

# Install missing dependencies
go mod download
go mod tidy
```

## Production Deployment

### Best Practices

1. **Configuration Management**
   - Use environment-specific config files
   - Never commit secrets to version control
   - Use GitHub token in environment for private repos

2. **Updates**
   - Test updates in staging first
   - Set `updateAutoApply: false` and manually control updates
   - Monitor logs for update failures

3. **Monitoring**
   - Monitor application logs
   - Set up health check endpoint
   - Track database size

4. **Backup**
   - Backup `./data/hub.db` regularly
   - Keep configuration backed up
   - Version control your config templates

## Next Steps

1. Integrate your Angular frontend
2. Test the single executable locally
3. Set up GitHub Actions for automated builds
4. Create releases with update manifests
5. Deploy to your server/cloud platform

For more details, see the [main README](../readme.md).
