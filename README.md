# Aura-Sentinel

> **AI-Driven Zero-Day Web Application Firewall (WAF) & Active Intrusion Prevention System**

Aura-Sentinel is a next-generation, highly concurrent WAF/IPS built entirely in Go. Unlike legacy signature-based rule engines that rely on static regex matching, Sentinel uses local LLMs to mathematically model normal web traffic. This allows for intercepting zero-day attacks through real-time Euclidean distance anomaly detection, deceiving attackers with dynamically generated honeypots, and streaming telemetry to a live Next.js SOC dashboard.

---

## Features

### Semantic Vector Trap: Zero-Day Detection
Integrates with Ollama (`nomic-embed-text`) to translate raw HTTP metadata into dense 768-dimensional embedding vectors. If a payload's spatial geometry deviates beyond a strict Euclidean threshold (e.g., `> 12.0`), the request is intercepted.

### Adaptive O(1) Edge-Blocking: The Fast Path
Thread-safe in-memory caching (`sync.Map`) ensures that once an attacker is identified, all subsequent requests from their IP address are dropped in < 0.001ms at the network edge, bypassing heavy AI inference entirely.

### Active AI Tarpit: Dynamic Honeypot
Instead of instantly dropping anomalous connections, Sentinel routes attackers to an AI-generated tarpit. Powered by `llama3.2`, it generates highly realistic fake vulnerabilities (such as simulated SQL errors or HTML administration panels) and forces the connection to hang, effectively draining the attacker's compute threads.

### Real-Time SOC Dashboard
A zero-dependency Server-Sent Events (SSE) broker in Go streams zero-day telemetry to a modern Next.js frontend. This features automatic React state reconciliation, cryptographic incident ID generation, and `localStorage` persistence.

---

## Technology Stack

| Layer | Technology |
|---|---|
| **Core Engine** | Go (Golang), Goroutines, Channels, `sync.Map`, Mutexes |
| **AI / Embeddings** | Ollama, `nomic-embed-text` |
| **AI / Generative** | Ollama, `llama3.2` |
| **Frontend SOC** | Next.js (React), Tailwind CSS, Server-Sent Events |

---

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- Ollama running locally on port `11434`

### 1. Clone and Setup

```bash
git clone https://github.com/yourusername/Aura-Sentinel.git
cd Aura-Sentinel

# Pull required AI models
ollama pull nomic-embed-text
ollama pull llama3.2

# Sync Go dependencies
go mod tidy
```

### 2. Boot the Infrastructure

Aura-Sentinel requires three active processes. Open three separate terminal windows:

**Terminal 1: WAF Backend (Go)**
```bash
go run ./cmd/sentinel/main.go
# WAF listens on :8080 | Telemetry stream on :8081
```

**Terminal 2: SOC Dashboard (Next.js)**
```bash
cd ui
npm install
npm run dev
# Dashboard available at http://localhost:3000
```

---

## Testing and Simulation

Use a third terminal window to simulate real-world traffic and trigger the anomaly detection pipeline.

---

### Test Case A: Normal Traffic

Simulate a standard user requesting the homepage:

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/" -UseBasicParsing
```

**Expected Result:**
- Euclidean distance calculated as low (e.g., `4.5`), remaining within the normal threshold.
- The request is proxied transparently to the backend.
- The terminal outputs: `[OK] Traffic Normal`.

---

### Test Case B: Zero-Day Attack (Hydra Brute Force)

Simulate an attacker probing for an administrative configuration file:

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/admin/config.php" `
  -Method POST `
  -UserAgent "Kali-Linux-Hydra" `
  -UseBasicParsing
```

**Expected Result:**
1. Mathematical deviation breaches the anomaly threshold (`> 12.0`).
2. The WAF routes the connection into the AI Tarpit, stalling it for approximately 3 seconds.
3. `llama3.2` generates a simulated HTML response to deceive the attacker.
4. The SOC Dashboard displays an alert with an automated Incident Report.

---

### Test Case C: Fast-Path Edge Block

Execute the exact same payload from Test Case B a second time:

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/admin/config.php" `
  -Method POST `
  -UserAgent "Kali-Linux-Hydra" `
  -UseBasicParsing
```

**Expected Result:**
- The attacker's IP address is recognized instantly in the `sync.Map` cache.
- The connection is terminated with HTTP 403 in O(1) time, bypassing the AI engine entirely.
- The terminal outputs: `[X] Blocked IP at the Edge`.

---

## Security Disclaimer

This project is built for educational and research purposes. The honeypot and tarpit modules are designed to study attacker behavior in controlled environments. Do not deploy against production traffic without a thorough security review.

---

## License

MIT License — see `LICENSE` for details.