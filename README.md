# 🔍 goPort

A fast, concurrent port scanner written in Go. Supports TCP banner grabbing and UDP service detection with protocol-specific probes.

## Features

- **TCP scanning** – connects and grabs service banners (SSH, HTTP, FTP, SMTP, etc.)
- **UDP scanning** – protocol-specific probes for accurate service detection
- **Concurrent** – worker pool architecture for high-speed scanning
- **Rate limiting** – built-in rate limiter to avoid network flooding
- **Banner grabbing** – automatically detects service versions

## Installation

### 🐧 Linux / macOS

```bash
git clone https://github.com/Slashas632/goPort
cd goPort
go build -o goPort ./cmd/app
```

### 🪟 Windows

```powershell
git clone https://github.com/Slashas632/goPort
cd goPort
go build -o goPort.exe ./cmd/app
```

## Usage

### 🐧 Linux / macOS

```bash
./goPort [flags]
```

### 🪟 Windows

```powershell
.\goPort.exe [flags]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-tcp` | false | Enable TCP scanning |
| `-udp` | false | Enable UDP scanning |
| `-ip` | 127.0.0.1 | Target IP address |
| `-p` | 65535 | Port or port range (e.g. `80` or `0-65535`) |
| `-w` | 1000 | Number of workers |

### Examples

**🐧 Linux / macOS**
```bash
# TCP scan common ports
./goPort -tcp -ip 10.0.0.1 -p 0-1024

# UDP scan all ports
./goPort -udp -ip 10.0.0.1 -p 0-65535

# TCP + UDP full scan
./goPort -tcp -udp -ip 10.0.0.1 -p 0-65535

# Custom worker count
./goPort -tcp -ip 10.0.0.1 -p 0-65535 -w 500
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
```

### Example Output

```
STATUS     IP                   PORT     BANNER
────────────────────────────────────────────────────────────────────────────────
[OPEN]     10.0.0.1             22       SSH-2.0-OpenSSH_9.2p1 Debian-2+deb12u3
[OPEN]     10.0.0.1             25       220 mail.example.com ESMTP Postfix
[OPEN]     10.0.0.1             53       (DNS)
[OPEN]     10.0.0.1             80       HTTP/1.1 200 OK
[OPEN]     10.0.0.1             110      +OK Dovecot ready
[OPEN]     10.0.0.1             143      * OK Dovecot ready
Work finished.
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
port-scanner/
├── cmd/
│   └── app/
│       └── main.go
└── internal/
    ├── cli/
    │   └── args.go
    ├── display/
    │   └── table.go
    ├── protocols/
    │   ├── TCP/
    │   │   └── TCP.go
    │   └── UDP/
    │       ├── UDP.go
    │       └── udp_probes.go
    ├── ratelimit/
    │   └── ratelimit.go
    └── scanner/
        └── engine.go
```

## ⚠️ Legal Disclaimer

This tool is intended for use on networks and systems you own or have explicit permission to scan. Unauthorized port scanning may be illegal in your jurisdiction.

## License

MIT