FROM mcr.microsoft.com/devcontainers/base:debian

# Install system packages for Tauri development and firewall
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    pkg-config \
    libwebkit2gtk-4.1-dev \
    libgtk-3-dev \
    libsoup-3.0-dev \
    libssl-dev \
    libayatana-appindicator3-dev \
    librsvg2-dev \
    libxdo-dev \
    file \
    iptables \
    ipset \
    iproute2 \
    dnsutils \
    aggregate \
    curl \
    jq \
    sudo \
    && rm -rf /var/lib/apt/lists/*

# Copy firewall script and grant passwordless sudo for it
COPY init-firewall.sh /usr/local/bin/init-firewall.sh
RUN chmod +x /usr/local/bin/init-firewall.sh \
    && echo "vscode ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/vscode-firewall \
    && chmod 0440 /etc/sudoers.d/vscode-firewall

USER vscode
