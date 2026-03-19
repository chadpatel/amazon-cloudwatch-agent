# Platform & Packaging

Design documents for multi-OS support (Linux, Windows, macOS), RPM/DEB/ZIP/tar packaging, Docker container builds, agent lifecycle management, and workload discovery.

## Designs

| Title | Path | Status | Summary |
|---|---|---|---|
| RPM Packaging | `packaging/linux/amazon-cloudwatch-agent.spec`, `Tools/src/create_rpm.sh` | 🚀 Implemented | RPM spec with %pre/%preun/%posttrans hooks, cwagent user/group creation, systemd/upstart service registration. Linux amd64+arm64. Predates design practice — no design doc available. (Jun 2020) |
| DEB Packaging | `packaging/debian/`, `Tools/src/create_deb.sh` | 🚀 Implemented | Debian package with preinst/postinst/prerm maintainer scripts, conffiles for config preservation, cwagent user/group creation. Linux amd64+arm64. Predates design practice — no design doc available. (Jun 2020) |
| Windows ZIP Packaging | `packaging/windows/`, `Tools/src/create_win.sh` | 🚀 Implemented | Windows ZIP with install.ps1/uninstall.ps1 PowerShell scripts, Windows service registration, ProgramFiles/ProgramData paths. Windows amd64. Predates design practice — no design doc available. (Jun 2020) |
| macOS tar Packaging | `packaging/darwin/`, `Tools/src/create_darwin.sh` | 🚀 Implemented | macOS tar.gz with launchd plist (KeepAlive, 60s throttle), darwin-specific agent-ctl script. macOS amd64+arm64. Predates design practice — no design doc available. (Dec 2020) |
| Docker Multi-Arch Builds | `amazon-cloudwatch-container-insights/cloudwatch-agent-dockerfile/` | 🚀 Implemented | 7 Dockerfile variants (localbin, source, localal2, localdeb, Windows localbin/source, localmsi) for linux/amd64+arm64 and windows/ltsc2019+ltsc2022. Predates design practice — no design doc available. (Jun 2020) |
| Systemd Service | `packaging/dependencies/amazon-cloudwatch-agent.service` | 🚀 Implemented | Systemd unit file (Type=simple, Restart=on-failure, RestartSec=60s) with EnvironmentFile for env-config.json. Predates design practice — no design doc available. (Jun 2020) |
| Agent Control Script | `packaging/dependencies/amazon-cloudwatch-agent-ctl` | 🚀 Implemented | Bash control script supporting start/stop/status/fetch-config/append-config/remove-config/set-log-level actions with SSM/file/default config sources and ec2/onPrem/auto modes. Predates design practice — no design doc available. (Jun 2020) |
