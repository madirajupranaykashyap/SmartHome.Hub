#!/bin/bash

set -e

BINARY="/usr/local/bin/smarthome-hub"

echo "Installing SmartHome Hub..."

# Ensure executable permissions
if [ -f "$BINARY" ]; then
    chmod +x "$BINARY"
fi

# Create systemd service
cat > /etc/systemd/system/smarthome-hub.service <<EOF
[Unit]
Description=SmartHome Hub
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/smarthome-hub
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd
systemctl daemon-reload

# Enable service on boot
systemctl enable smarthome-hub.service

# Start service
systemctl restart smarthome-hub.service || true

echo "SmartHome Hub installed successfully."