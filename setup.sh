#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   SmartHome Hub - Single Executable   ║${NC}"
echo -e "${BLUE}║         Build & Setup Script          ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Function to print colored output
print_step() {
    echo -e "${GREEN}▶${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

# Check dependencies
check_dependencies() {
    print_step "Checking dependencies..."
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        exit 1
    fi
    print_success "Go $(go version | awk '{print $3}')"
    
    if ! command -v node &> /dev/null; then
        print_info "Node.js not found (optional, needed only if building frontend)"
    else
        print_success "Node.js $(node --version)"
    fi
    
    if ! command -v npm &> /dev/null; then
        print_info "npm not found (optional, needed only if building frontend)"
    else
        print_success "npm $(npm --version)"
    fi
}

# Prepare frontend
prepare_frontend() {
    print_step "Preparing Angular frontend..."
    
    if [ ! -d "frontend" ]; then
        print_info "Frontend directory not found"
        read -p "Enter path to your Angular frontend (or press Enter to skip): " frontend_path
        
        if [ -z "$frontend_path" ]; then
            print_info "Skipping frontend setup. Manual embedding required."
            return
        fi
        
        if [ ! -d "$frontend_path" ]; then
            print_error "Frontend path not found: $frontend_path"
            return
        fi
        
        cp -r "$frontend_path" frontend
        print_success "Frontend copied"
    fi
    
    if [ ! -f "frontend/package.json" ]; then
        print_error "frontend/package.json not found"
        return
    fi
    
    print_step "Building Angular frontend..."
    cd frontend
    
    if [ ! -d "node_modules" ]; then
        npm install
    fi
    
    npm run build
    print_success "Frontend built successfully"
    cd ..
}

# Embed frontend
embed_frontend() {
    print_step "Embedding frontend in Go binary..."
    
    cd backend
    
    mkdir -p pkg/hub/frontend
    
    # Find dist directory (could be dist/project-name or just dist)
    if [ -d "../frontend/dist" ]; then
        # Check for subdirectory
        if [ "$(ls -A ../frontend/dist)" ]; then
            if [ -f "../frontend/dist/index.html" ]; then
                cp -r ../frontend/dist/* pkg/hub/frontend/
            else
                # Might be in a subdirectory
                subdirs=(../frontend/dist/*/)
                if [ ${#subdirs[@]} -gt 0 ]; then
                    cp -r ../frontend/dist/*/* pkg/hub/frontend/
                fi
            fi
        fi
    fi
    
    if [ ! -f "pkg/hub/frontend/index.html" ]; then
        print_error "Frontend files not found in dist/"
        print_info "Ensure Angular build output contains index.html"
        cd ..
        return
    fi
    
    print_success "Frontend embedded successfully"
    print_info "Files: $(ls pkg/hub/frontend/ | wc -l) items"
    cd ..
}

# Build executable
build_executable() {
    print_step "Building executable..."
    
    cd backend
    
    read -p "Build for which platform? (1=current, 2=all, 3=linux, 4=macos, 5=windows): " platform_choice
    
    case $platform_choice in
        1)
            print_step "Building for current platform..."
            make build
            print_success "Executable created: ./build/smarthome-hub"
            ;;
        2)
            print_step "Building for all platforms..."
            make build-all
            print_success "Executables created in ./build/"
            ;;
        3)
            print_step "Building for Linux..."
            make build-linux
            print_success "Executables created: ./build/linux/"
            ;;
        4)
            print_step "Building for macOS..."
            make build-macos
            print_success "Executables created: ./build/macos/"
            ;;
        5)
            print_step "Building for Windows..."
            make build-windows
            print_success "Executables created: ./build/windows/"
            ;;
        *)
            print_error "Invalid choice"
            cd ..
            return
            ;;
    esac
    
    cd ..
}

# Test executable
test_executable() {
    print_step "Testing executable..."
    
    cd backend
    
    if [ -f "build/smarthome-hub" ]; then
        chmod +x build/smarthome-hub
        
        print_info "Starting application (press Ctrl+C to stop)..."
        timeout 5 ./build/smarthome-hub || true
        
        if [ -f "config.json" ]; then
            print_success "Configuration created"
            print_info "Config file: config.json"
        fi
    fi
    
    cd ..
}

# Main flow
main() {
    check_dependencies
    echo ""
    
    read -p "Continue with setup? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Setup cancelled"
        exit 0
    fi
    echo ""
    
    prepare_frontend
    echo ""
    
    embed_frontend
    echo ""
    
    build_executable
    echo ""
    
    read -p "Test the application? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        test_executable
        echo ""
    fi
    
    echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}✓ Setup complete!${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${GREEN}Next steps:${NC}"
    echo -e "  1. Test locally: ${YELLOW}./backend/build/smarthome-hub${NC}"
    echo -e "  2. Access UI: ${YELLOW}http://localhost:8080${NC}"
    echo -e "  3. Configure: Edit ${YELLOW}config.json${NC}"
    echo -e "  4. Deploy: Upload executable to your server"
    echo ""
    echo -e "${GREEN}Documentation:${NC}"
    echo -e "  - Setup guide: ${YELLOW}SETUP.md${NC}"
    echo -e "  - Deployment: ${YELLOW}DEPLOYMENT.md${NC}"
    echo -e "  - Build help: ${YELLOW}cd backend && make help${NC}"
}

# Run main function
main
