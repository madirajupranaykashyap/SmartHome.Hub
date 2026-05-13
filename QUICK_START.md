# Single Executable Quick Reference

**Goal**: One executable that runs both your Go backend and Angular frontend with auto-updates.

## 5-Minute Setup

### Linux/macOS

```bash
# Make setup script executable
chmod +x setup.sh

# Run setup (interactive)
./setup.sh

# Or run it directly
./backend/build/smarthome-hub
# Open: http://localhost:8080
```

### Windows

```powershell
# Run PowerShell script
.\setup.ps1

# Or build specific platform
.\setup.ps1 -Platform windows

# Run executable
.\backend\build\smarthome-hub.exe
# Open: http://localhost:8080
```

## Manual Setup

If you prefer to do it step by step:

```bash
# 1. Build your Angular frontend
cd frontend
npm install
npm run build
cd ..

# 2. Embed frontend in Go
mkdir -p backend/pkg/hub/frontend
cp -r frontend/dist/* backend/pkg/hub/frontend/

# 3. Build executable
cd backend
make build  # or: make build-all for all platforms

# 4. Run it
./build/smarthome-hub

# 5. Access at http://localhost:8080
```

## Project Structure After Setup

```
SmartHome.Hub/
├── backend/
│   ├── build/
│   │   └── smarthome-hub          ← Your single executable
│   ├── cmd/launcher/
│   │   └── main.go                ← Entry point
│   ├── pkg/hub/frontend/          ← Embedded Angular files
│   │   ├── index.html
│   │   ├── styles.css
│   │   └── ...
│   ├── config.json                ← Auto-created on first run
│   ├── data/
│   │   └── hub.db                 ← SQLite database
│   └── Makefile                   ← Build helpers
├── frontend/                      ← Your Angular app
│   └── dist/                      ← Built frontend
├── config.json                    ← App configuration
├── SETUP.md                       ← Detailed setup guide
├── DEPLOYMENT.md                  ← Deployment guide
└── setup.sh / setup.ps1           ← Automated setup scripts
```

## Configuration (config.json)

Created automatically on first run:

```json
{
  "addr": ":8080",                              // Server address
  "databasePath": "./data/hub.db",              // Database location
  "updateOwner": "Project-SmartHome",           // GitHub owner for updates
  "updateRepo": "SmartHome.Hub",                // GitHub repo for updates
  "enableUpdateCheck": true,                    // Check for updates on startup
  "updateAutoApply": true,                      // Auto-apply updates
  "currentVersion": "1.0.0",                    // Current version
  "updateGitHubToken": "",                      // GitHub token (if needed)
  "appEnv": "production",
  "debugMode": false
}
```

### Key Options to Modify

| Option | Use Case |
|--------|----------|
| `addr: ":3000"` | Use different port |
| `updateAutoApply: false` | Manual update control |
| `debugMode: true` | Verbose logging |
| `databasePath: "/custom/path"` | Custom database location |
| `updateGitHubToken: "your_token"` | Private repo access |

## Commands

```bash
cd backend

# Build
make build              # Current platform
make build-all          # All platforms (Linux, Windows, macOS)
make build-linux        # Linux only
make build-windows      # Windows only
make build-macos        # macOS only

# Development
make dev                # Run in development mode

# Utilities
make test               # Run tests
make fmt                # Format code
make clean              # Clean build artifacts
make help               # Show all commands
```

## Accessing the App

Once running:

- **Frontend**: http://localhost:8080
- **API Docs**: http://localhost:8080/swagger
- **API Base**: http://localhost:8080/api
- **Database**: `./data/hub.db` (automatic)

## What's Included

The single executable contains:

1. **Go Backend**
   - RESTful API
   - SQLite database
   - Authentication (JWT)
   - Swagger documentation

2. **Angular Frontend**
   - All static assets (HTML, CSS, JS)
   - Production-optimized build
   - Served from `/` path

3. **Auto-Update System**
   - Checks GitHub releases on startup
   - Downloads only changed files (differential)
   - Applies updates automatically or manually
   - Tracks current version

4. **Configuration Management**
   - Automatic config file creation
   - Persistent settings
   - Environment-specific configs

## Deployment

### Docker

```bash
# Build image
docker build -t smarthome-hub .

# Run container
docker run -p 8080:8080 smarthome-hub
```

### Linux Systemd

Copy executable to `/opt/smarthome-hub/`:

```bash
sudo cp backend/build/smarthome-hub /opt/smarthome-hub/
sudo chmod +x /opt/smarthome-hub/smarthome-hub

# Create service file (see DEPLOYMENT.md)
# Then:
sudo systemctl enable smarthome-hub
sudo systemctl start smarthome-hub
```

### Windows Task Scheduler

1. Copy `backend/build/smarthome-hub.exe` to program directory
2. Create Task Scheduler task pointing to executable
3. Set to run at startup

### Cloud (AWS/Azure/GCP)

1. Build for Linux: `make build-linux`
2. Upload to cloud storage or container registry
3. Deploy to VM or container service
4. Configure reverse proxy (Nginx/Caddy) for HTTPS

See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed instructions.

## Automatic Updates

### For Users

When new version available (on next startup):
- Checks GitHub releases
- Downloads changed files only
- Applies update automatically (if configured)
- Restarts required

### For Developers

1. Make changes to code/frontend
2. Build everything: `make build-all`
3. Upload binaries + manifest to GitHub release
4. Users auto-update on next run

Create `update-manifest.json`:

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
      "mode": 493,
      "goos": "linux",
      "goarch": "amd64"
    }
  ]
}
```

See [DEPLOYMENT.md](DEPLOYMENT.md) for full update setup.

## Troubleshooting

### Frontend Not Loading

```bash
# Check frontend files exist
ls backend/pkg/hub/frontend/

# Rebuild frontend
cd frontend && npm run build
cp -r dist/* ../backend/pkg/hub/frontend/

# Rebuild executable
cd ../backend && make clean && make build
```

### Update Check Fails

```bash
# Check config
cat config.json

# Verify GitHub access
curl -I https://api.github.com

# Use token for private repos
# Edit config.json: "updateGitHubToken": "your_token"
```

### Database Issues

```bash
# Ensure data directory is writable
mkdir -p data
chmod 755 data

# Or change path in config.json
```

### Build Fails

```bash
# Clean and retry
make clean && make build

# Check Go version
go version

# Update dependencies
go mod download
go mod tidy
```

## File Locations

After running the executable:

- `config.json` - Configuration
- `data/hub.db` - SQLite database
- `data/hub.db-shm` - Database shared memory
- `data/hub.db-wal` - Database write-ahead log

All in the same directory as the executable.

## Performance

- **Binary Size**: ~50-100 MB (depending on platform)
- **Memory Usage**: ~50-100 MB typical
- **Startup Time**: 1-2 seconds
- **API Response**: <100ms typical
- **Updates**: Only downloads changed files

## Security Tips

1. **Change default credentials** (if applicable)
2. **Use HTTPS** via reverse proxy (Nginx/Caddy)
3. **Backup database** regularly
4. **Keep tokens out of config** (use environment variables)
5. **Restrict file permissions**: `chmod 600 hub.db`

## Next Steps

1. ✅ Run setup: `./setup.sh` or `./setup.ps1`
2. ✅ Test locally: `./backend/build/smarthome-hub`
3. ✅ Access frontend: http://localhost:8080
4. ✅ Configure settings: Edit `config.json`
5. ✅ Set up GitHub releases (optional)
6. ✅ Deploy to your server

## Getting Help

- **Detailed setup**: See [SETUP.md](SETUP.md)
- **Deployment options**: See [DEPLOYMENT.md](DEPLOYMENT.md)
- **Build commands**: `cd backend && make help`
- **Development**: Run in dev mode with `make dev`

---

**Happy deploying!** 🚀

For detailed guides, see SETUP.md and DEPLOYMENT.md in the project root.
