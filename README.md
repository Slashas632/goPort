# 🔍 goPort

[![Go Reference](https://pkg.go.dev/badge/github.com/Slashas632/goPort.svg)](https://pkg.go.dev/github.com/Slashas632/goPort)
[![Test](https://github.com/Slashas632/goPort/actions/workflows/test.yml/badge.svg)](https://github.com/Slashas632/goPort/actions/workflows/test.yml)
[![License: MIT](https://img.shields.io/github/license/Slashas632/goPort)](LICENSE)

A fast, concurrent TCP/UDP port scanner written in Go, with protocol-specific banner grabbing, JSON export, and a Lua plugin system for extending scan behavior — a lightweight alternative to nmap for quick concurrent scans.

## Features

- **TCP scanning** – connects and grabs service banners (SSH, HTTP, FTP, SMTP, etc.)
- **UDP scanning** – protocol-specific probes for accurate service detection
- **Concurrent** – worker pool architecture for high-speed scanning
- **Rate limiting** – built-in rate limiter to avoid network flooding
- **Banner grabbing** – automatically detects service versions
- **Plugin system** – extend functionality with Lua scripts
- **Json export** - export results to .json format

## Installation

### Arch Linux (AUR)

> ⚠️ **AUR package publishing is currently blocked.** Arch Linux temporarily disabled AUR write access ([announcement](https://www.bleepingcomputer.com/news/security/arch-linux-disables-aur-package-adoption-to-stop-malware-flood/)) following a wave of malicious package takeovers, so the AUR package cannot be updated right now. The AUR is still stuck on **v1.2.0** — for the current release, build from source (below) or grab a prebuilt binary from the [Releases page](https://github.com/Slashas632/goPort/releases/latest).

```bash
yay -S goport   # currently installs v1.2.0 until AUR publishing reopens
```

### 🐧 Linux / macOS (from source)

> Requires [Go 1.22+](https://go.dev/dl/)

```bash
git clone https://github.com/Slashas632/goPort
cd goPort
go build -o goPort ./cmd/app
```

### 🪟 Windows (from source)

> Requires [Go 1.22+](https://go.dev/dl/)

```powershell
git clone https://github.com/Slashas632/goPort
cd goPort
go build -o goPort.exe ./cmd/app
```

## Usage

### 🐧 Linux / macOS

```bash
# AUR or installed to PATH
goPort [flags]

# Built from source
./goPort [flags]
```

### 🪟 Windows

```powershell
# Built from source
.\goPort.exe [flags]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-tcp` | false | Enable TCP scanning |
| `-udp` | false | Enable UDP scanning |
| `-ip` | 127.0.0.1 | Target IP address |
| `-p` | 65535 | Port or port range (e.g. `80` or `0-65535`) |
| `-w` | 500 | Number of workers |
| `-install` | – | Install a Lua plugin |
| `-uninstall` | – | Uninstall a Lua plugin |
| `-json` | - | Export results to json |
### Examples

**🐧 Linux / macOS**
```bash
# TCP scan common ports
goPort -tcp -ip 10.0.0.1 -p 0-1024

# UDP scan all ports
goPort -udp -ip 10.0.0.1 -p 0-65535

# TCP + UDP full scan
goPort -tcp -udp -ip 10.0.0.1 -p 0-65535

# Custom worker count
goPort -tcp -ip 10.0.0.1 -p 0-65535 -w 500

# Export to json
goPort -tcp -ip 10.0.0.1 -p 0-65535 -json output.json
```

**🪟 Windows**
```powershell
# TCP scan common ports
.\goPort.exe -tcp -ip 10.0.0.1 -p 0-1024

# UDP scan all ports
.\goPort.exe -udp -ip 10.0.0.1 -p 0-65535

# TCP + UDP full scan
.\goPort.exe -tcp -udp -ip 10.0.0.1 -p 0-65535

# Custom worker count
.\goPort.exe -tcp -ip 10.0.0.1 -p 0-65535 -w 500

# Export to json
.\goPort.exe -tcp -ip 10.0.0.1 -p 0-65535 -json output.json
```

### Example Output

```
STATUS     IP                   PORT     BANNER
────────────────────────────────────────────────────────────────────────────────
[OPEN]     10.0.0.1             22       SSH-2.0-OpenSSH_9.2p1 Debian-2+deb12u3
[OPEN]     10.0.0.1             25       220 mail.example.com ESMTP Postfix
[OPEN]     10.0.0.1             53       DNS
[OPEN]     10.0.0.1             80       HTTP/1.1 200 OK
[OPEN]     10.0.0.1             110      +OK Dovecot ready
[OPEN]     10.0.0.1             143      * OK Dovecot ready
Work finished.
```

## 🔌 Plugin System

> ⚠️ Plugins use a third-party Lua interpreter ([gopher-lua](https://github.com/yuin/gopher-lua)). Only install plugins from sources you trust.

Plugins are Lua scripts that run after each port is scanned. They receive the IP, port, and banner as arguments.

### Plugin Structure

```lua
function scan(ip, port, banner)
    if string.find(banner, "SSH") then
        print("[SSH] " .. ip .. ":" .. port .. " -> " .. banner)
    end
end
```

### Installing a Plugin

**🐧 Linux / macOS**
```bash
goPort -install /home/user/myplugin.lua
goPort -install ~/myplugin.lua
```

**🪟 Windows**
```powershell
.\goPort.exe -install C:\Users\user\myplugin.lua
.\goPort.exe -install .\myplugin.lua
```

### Uninstalling a Plugin

**🐧 Linux / macOS**
```bash
goPort -uninstall myplugin.lua
```

**🪟 Windows**
```powershell
.\goPort.exe -uninstall myplugin.lua
```

> Note: `-uninstall` takes only the filename, not the full path.

### Plugin Storage

Plugins are stored in:
- **Linux / macOS:** `~/.goPort/plugins/`
- **Windows:** `C:\Users\<user>\.goPort\plugins\`

### Example Plugin

```lua
-- Detects common services and prints alerts
function scan(ip, port, banner)
    if string.find(banner, "SSH") then
        print("[SSH]  " .. ip .. ":" .. port .. " -> " .. banner)
    end

    if string.find(banner, "HTTP") then
        print("[HTTP] " .. ip .. ":" .. port .. " -> " .. banner)
    end

    if string.find(banner, "220") then
        print("[SMTP] " .. ip .. ":" .. port .. " -> " .. banner)
    end
end
```

## UDP Probes

The scanner uses protocol-specific payloads for accurate UDP detection:

| Port | Protocol |
|------|----------|
| 53 | DNS |
| 69 | TFTP |
| 111 | RPC |
| 123 | NTP |
| 137 | NetBIOS |
| 161 | SNMP |
| 443 | QUIC |
| 500 | IKE/VPN |
| 514 | Syslog |
| 1900 | SSDP/UPnP |
| 5353 | mDNS |
| 11211 | Memcached |
| 27015 | Steam |
| 51820 | WireGuard |

Unknown ports fall back to generic probes (NULL, CRLF, HELLO).

## Project Structure

```
goPort/
├── cmd/
│   └── app/
│       └── main.go
└── internal/
    ├── output/
    │   ├── json.go
    │   └── json_test.go
    ├── cli/
    │   ├── args.go
    │   └── args_test.go
    ├── display/
    │   ├── table.go
    │   └── table_test.go
    ├── plugins/
    │   ├── manager.go
    │   ├── manager_test.go
    │   ├── runner.go
    │   └── runner_test.go
    ├── protocols/
    │   ├── TCP/
    │   │   ├── TCP.go
    │   │   └── TCP_test.go
    │   └── UDP/
    │       ├── UDP.go
    │       ├── udp_probes.go
    │       └── udp_test.go
    ├── ratelimit/
    │   ├── ratelimit.go
    │   └── ratelimit_test.go
    └── scanner/
        ├── engine.go
        └── engine_test.go
```

## Testing

The project has an automated test suite covering CLI argument parsing, the rate limiter, JSON output, the plugin manager and Lua runner, and both the TCP and UDP scanning paths (including end-to-end scans against local test servers). CI runs the full suite with the race detector on every push.

```bash
go test ./... -cover
go test ./... -race
```

## ⚠️ Legal Disclaimer

This tool is intended for use on networks and systems you own or have explicit permission to scan. Unauthorized port scanning may be illegal in your jurisdiction.

## License

MIT