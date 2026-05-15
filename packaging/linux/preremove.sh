#!/bin/bash

set -e

echo "Removing SmartHome Hub..."

# Stop service
systemctl stop smarthome-hub.service || true

# Disable service
systemctl disable smarthome-hub.service || true

# Remove service file
rm -f /etc/systemd/system/smarthome-hub.service

# Reload systemd
systemctl daemon-reload || true

echo "SmartHome Hub removed successfully."