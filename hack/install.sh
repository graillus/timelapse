#!/usr/bin/env bash
#
# Install Timelapse Agent
#

# Install tlagent binary
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(arch)
TLAGENT_VERSION="0.1.0"
TLAGENT_URL="https://github.com/graillus/timelapse/releases/download/${TLAGENT_VERSION}/tlagent_${TLAGENT_VERSION}_${OS}_${ARCH}.tar.gz"

curl -sL "$TLAGENT_URL" | tar xzf - > /usr/bin/tlagent
chmod +x /usr/bin/tlagent

# Setup configuration
mkdir -p /etc/tlagent
cat > /etc/tlagent/config <<EOF
interval: 30s
serverUrl: http://timelapse:8990
EOF

# Install systemd unit
cp ./systemd/tlagent.service /lib/systemd/system
chmod 644 /lib/systemd/system/tlagent.service

systemctl daemon-reload
systemctl enable tlagent.service
