# Custom VPN Implementation in Go

**Author:** Martin Wong
**Purpose:** Learn distributed systems, networking, and backend engineering through building a production-quality VPN from scratch.

## Project Vision

Build a custom VPN that demonstrates deep understanding of:
- Network protocols and packet processing
- Cryptography and secure communications
- Distributed systems patterns
- Concurrent programming in Go
- Production-ready backend systems

## Architecture Overview

```
┌─────────────┐                   ┌─────────────┐
│   Client    │◄────Encrypted────►│   Server    │
│             │      UDP          │             │
│  ┌───────┐  │                   │  ┌───────┐  │
│  │  TUN  │  │                   │  │  TUN  │  │
│  │Device │  │                   │  │Device │  │
│  └───────┘  │                   │  └───────┘  │
└─────────────┘                   └─────────────┘
       │                                 │
       ▼                                 ▼
  Local Apps                         Internet
```

### Core Components

1. **TUN Device Manager** - Virtual network interface creation and management
2. **Packet Processor** - Read/write packets from/to TUN device
3. **Crypto Engine** - Encryption/decryption with authenticated encryption
4. **Connection Manager** - Session lifecycle and state management
5. **Protocol Handler** - Custom wire protocol implementation
6. **Server Router** - Traffic routing and NAT for multiple clients

## Technology Stack

- **Language:** Go 1.21+
- **TUN/TAP:** `github.com/songgao/water`
- **Crypto:** `crypto/aes`, `crypto/cipher`, `golang.org/x/crypto`
- **Networking:** Standard library `net`
- **Logging:** `log/slog`
- **Metrics:** `github.com/prometheus/client_golang`

## Project Structure

```
govpn/
├── README.md                 # This file
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── cmd/
│   ├── server/              # VPN server binary
│   │   └── main.go
│   └── client/              # VPN client binary
│       └── main.go
├── internal/
│   ├── tunnel/              # TUN device management
│   │   ├── tun.go
│   │   └── tun_test.go
│   ├── crypto/              # Encryption/key exchange
│   │   ├── cipher.go
│   │   ├── handshake.go
│   │   └── keys.go
│   ├── protocol/            # Wire protocol definition
│   │   ├── packet.go
│   │   └── messages.go
│   ├── connection/          # Connection management
│   │   ├── session.go
│   │   ├── pool.go
│   │   └── keepalive.go
│   ├── router/              # Packet routing logic
│   │   ├── route.go
│   │   └── nat.go
│   └── config/              # Configuration management
│       └── config.go
├── pkg/
│   └── vpn/                 # Public API (if needed)
│       └── client.go
├── configs/
│   ├── server.yaml          # Server configuration example
│   └── client.yaml          # Client configuration example
├── scripts/
│   ├── setup_tun.sh         # TUN device setup script
│   └── teardown.sh          # Cleanup script
├── docs/
│   ├── protocol.md          # Wire protocol specification
│   ├── architecture.md      # Detailed architecture
│   └── performance.md       # Performance benchmarks
└── Makefile                 # Build and test automation
```

## Implementation Roadmap

### Phase 0: Project Setup (Days 1-2)
**Goal:** Get basic project structure and development environment ready

**Tasks:**
- [X] Initialize Go module: `go mod init github.com/yourusername/govpn`
- [X] Set up project directory structure
- [X] Create Makefile with build/test/run targets
- [X] Add .gitignore for Go projects
- [X] Write basic configuration loading (YAML)
- [X] Set up logging with slog

**Deliverable:** Can build and run hello-world server/client binaries

---

### Phase 1: Basic Tunnel (Week 1, Days 3-9)
**Goal:** Unencrypted UDP tunnel that forwards IP packets

**Learning Focus:**
- TUN device creation and configuration
- IP packet structure (headers, payload)
- UDP socket programming in Go
- Goroutines and channels for concurrent I/O

**Implementation Steps:**

1. **TUN Device Creation** (Day 3-4)
   - [X] Create TUN device using water library
   - [X] Configure IP address and routing
   - [X] Read IP packets from TUN device
   - [X] Write packets back to TUN device
   - [X] Test: Can read/write packets locally

2. **UDP Communication** (Day 5-6)
   - [X] Set up UDP listener on server
   - [X] Set up UDP client connection
   - [X] Send/receive raw packets over UDP
   - [X] Test: Can send packet from client to server

3. **Basic Tunnel Loop** (Day 7-9)
   - [X] Client: TUN → UDP → Server
   - [X] Server: UDP → TUN → Internet
   - [X] Server: Internet → TUN → UDP → Client
   - [X] Client: UDP → TUN → App
   - [X] Test: Can ping through tunnel (unencrypted)

**Key Code Files:**
- `internal/tunnel/tun.go` - TUN device wrapper
- `cmd/client/main.go` - Basic client loop
- `cmd/server/main.go` - Basic server loop

**Success Criteria:**
- Client can ping 8.8.8.8 through VPN server
- Server forwards traffic to internet and back
- No encryption yet (Phase 2)

---

### Phase 2: Encryption & Security (Week 2, Days 10-16)
**Goal:** Add authenticated encryption and secure handshake

**Learning Focus:**
- AEAD ciphers (AES-256-GCM or ChaCha20-Poly1305)
- Key exchange protocols (Diffie-Hellman)
- Nonce generation and replay protection
- Secure session establishment

**Implementation Steps:**

1. **Pre-Shared Key Mode** (Day 10-11)
   - [ ] Define wire protocol (packet format)
   - [ ] Add AES-256-GCM encryption to packets
   - [ ] Generate random nonces per packet
   - [ ] Add packet authentication (AEAD)
   - [ ] Test: Encrypted ping works

2. **Handshake Protocol** (Day 12-14)
   - [ ] Implement 3-way handshake (ClientHello, ServerHello, Finalize)
   - [ ] Add Diffie-Hellman key exchange (X25519)
   - [ ] Derive session keys using HKDF
   - [ ] Add replay protection (sliding window)
   - [ ] Test: Handshake completes, keys derived

3. **Session Management** (Day 15-16)
   - [ ] Track active sessions (map of sessions)
   - [ ] Session timeout and cleanup
   - [ ] Rekey after N packets or time
   - [ ] Test: Multiple clients can connect

**Key Code Files:**
- `internal/crypto/cipher.go` - Encryption wrapper
- `internal/crypto/handshake.go` - Key exchange
- `internal/protocol/packet.go` - Wire format
- `internal/connection/session.go` - Session state

**Success Criteria:**
- All traffic is encrypted end-to-end
- Session keys are derived securely
- Replay attacks are prevented
- Multiple clients can connect simultaneously

---

### Phase 3: Production Features (Week 3, Days 17-23)
**Goal:** Add reliability, observability, and operational features

**Learning Focus:**
- Keepalive and heartbeat mechanisms
- Graceful error handling and recovery
- Metrics collection (Prometheus)
- Structured logging
- Configuration management

**Implementation Steps:**

1. **Connection Reliability** (Day 17-19)
   - [ ] Keepalive packets (every 10s)
   - [ ] Detect dead connections (3 missed keepalives)
   - [ ] Client auto-reconnection with backoff
   - [ ] Handle connection migration (IP change)
   - [ ] Test: Survives network disruption

2. **Observability** (Day 20-21)
   - [ ] Add Prometheus metrics (packets, bytes, latency)
   - [ ] Structured logging with slog (DEBUG, INFO, ERROR)
   - [ ] Log important events (connect, disconnect, errors)
   - [ ] Add metrics endpoint `:9090/metrics`
   - [ ] Test: Can view metrics in browser

3. **Configuration & Deployment** (Day 22-23)
   - [ ] YAML configuration files
   - [ ] Command-line flags (override config)
   - [ ] Systemd service files
   - [ ] Docker containerization
   - [ ] Test: Deploy on Linux VM

**Key Code Files:**
- `internal/connection/keepalive.go` - Heartbeat logic
- `internal/metrics/metrics.go` - Prometheus metrics
- `internal/config/config.go` - Configuration loader

**Success Criteria:**
- VPN stays connected for hours
- Can monitor traffic with metrics
- Easy to deploy and configure
- Logs help debug issues

---

### Phase 4: Advanced Features (Week 4+, Optional)
**Goal:** Performance optimization and advanced capabilities

**Potential Features:**
- [ ] Multi-threaded packet processing (worker pools)
- [ ] Zero-copy optimizations (buffer pooling)
- [ ] NAT traversal (UDP hole punching)
- [ ] Split tunneling (route only specific IPs)
- [ ] IPv6 support
- [ ] Load balancing across multiple servers
- [ ] Web dashboard for monitoring

---

## Wire Protocol Specification

### Packet Format (Version 1)

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|Version|  Type |    Reserved   |          Session ID           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                         Sequence Number                        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                         Payload Length                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                                                               |
+                         Payload (encrypted)                   +
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                                                               |
+                      Authentication Tag (16 bytes)            +
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

**Packet Types:**
- `0x01` - Handshake Request (ClientHello)
- `0x02` - Handshake Response (ServerHello)
- `0x03` - Handshake Complete
- `0x04` - Data (encrypted IP packet)
- `0x05` - Keepalive
- `0x06` - Disconnect

---

## Development Commands

### Build
```bash
make build           # Build client and server
make build-client    # Build client only
make build-server    # Build server only
```

### Run
```bash
# Server (requires root for TUN device)
sudo ./bin/server -config configs/server.yaml

# Client (requires root for TUN device)
sudo ./bin/client -config configs/client.yaml -server 192.168.1.100:8080
```

### Test
```bash
make test           # Run all tests
make test-unit      # Unit tests only
make test-integ     # Integration tests
make bench          # Benchmarks
```

### TUN Device Setup (Linux)
```bash
# Create TUN device
sudo ip tuntap add mode tun dev tun0
sudo ip addr add 10.0.0.1/24 dev tun0
sudo ip link set dev tun0 up

# Enable IP forwarding (server only)
sudo sysctl -w net.ipv4.ip_forward=1
sudo iptables -t nat -A POSTROUTING -s 10.0.0.0/24 -o eth0 -j MASQUERADE
```

---

## Technical Concepts to Understand

### 1. TUN vs TAP Devices
- **TUN:** Layer 3 (IP packets) - We use this
- **TAP:** Layer 2 (Ethernet frames) - More overhead

### 2. Packet Flow (Client)
```
Application → TCP/IP Stack → TUN Device → Read by VPN Client 
→ Encrypt → Encapsulate → Send UDP → Server
```

### 3. Packet Flow (Server)
```
Receive UDP → Decapsulate → Decrypt → Write to TUN Device 
→ Route to Internet → NAT → Return path
```

### 4. Key Exchange (Simplified)
```
Client                              Server
  |                                   |
  |----ClientHello + DH_public------->|
  |                                   | (generates DH keypair)
  |<---ServerHello + DH_public--------|
  |                                   |
  | (both derive shared secret)       |
  | (derive session keys)             | (derive session keys)
  |                                   |
  |----Finalize (encrypted)---------->|
  |                                   |
  |<---Data (encrypted)-------------->|
```

### 5. Concurrency Patterns
- **TUN Reader Goroutine:** Reads packets from TUN device
- **TUN Writer Goroutine:** Writes packets to TUN device
- **UDP Reader Goroutine:** Receives packets from network
- **UDP Writer Goroutine:** Sends packets to network
- **Channels:** Pass packets between goroutines safely

---

## Performance Targets

### Minimum Viable (Phase 1-2)
- Throughput: 100+ Mbps
- Latency: <10ms overhead
- Connections: 10+ concurrent clients

### Production Ready (Phase 3)
- Throughput: 500+ Mbps
- Latency: <5ms overhead
- Connections: 100+ concurrent clients
- Uptime: 99.9% (8 hours downtime/year)

### Optimized (Phase 4)
- Throughput: 1+ Gbps
- Latency: <2ms overhead
- Connections: 1000+ concurrent clients
- CPU: <50% on single core at 1Gbps

---

## Resources for Learning

### Go Networking
- Go standard library: `net`, `net/http`
- Book: "Network Programming with Go" by Jan Newmarch

### Cryptography
- Go x/crypto documentation
- WireGuard whitepaper (simple, modern design)
- Book: "Serious Cryptography" by Jean-Philippe Aumasson

### Networking Fundamentals
- RFC 791: Internet Protocol (IP)
- RFC 768: User Datagram Protocol (UDP)
- Understanding TUN/TAP: kernel.org documentation

### VPN References
- WireGuard: Study its simplicity
- OpenVPN: Traditional approach (more complex)
- Tailscale: Modern mesh VPN architecture

---

## Troubleshooting Guide

### Common Issues

**TUN device creation fails:**
```bash
# Ensure you have permissions
sudo setcap cap_net_admin+ep ./bin/client
sudo setcap cap_net_admin+ep ./bin/server
```

**No internet through tunnel:**
```bash
# Check IP forwarding
cat /proc/sys/net/ipv4/ip_forward  # Should be 1

# Check iptables NAT rule
sudo iptables -t nat -L -v
```

**Connection timeouts:**
- Check firewall allows UDP on VPN port
- Verify server is listening: `netstat -ulnp | grep 8080`
- Check client can reach server: `nc -zvu server_ip 8080`

---

## Context for LLM Reprompting

When asking an LLM for help with this project, provide:

1. **Current Phase:** "I'm on Phase X, working on [specific feature]"
2. **What's Done:** "I've completed X, Y, Z"
3. **Specific Issue:** "I'm stuck on [specific problem]"
4. **Code Context:** Share relevant code snippets
5. **Error Messages:** Full error text and stack traces

**Example Prompt:**
> "I'm building a custom VPN in Go (see README context). I'm on Phase 2, implementing AES-256-GCM encryption. I've successfully created TUN device and UDP tunnel. Now I need help with: [specific crypto question]. Here's my current cipher.go implementation: [code]"# MightyVPN
