# 🛡️ Aura-Sentinel

> **AI-Driven Zero-Day Web Application Firewall (WAF) & Active Intrusion Prevention System**

Aura-Sentinel is a next-generation, highly concurrent WAF/IPS built entirely in Go. Unlike legacy signature-based rule engines that rely on static regex matching, Sentinel uses local LLMs to mathematically model normal web traffic — intercepting zero-day attacks through real-time Euclidean distance anomaly detection, deceiving attackers with dynamically generated honeypots, and streaming telemetry to a live Next.js SOC dashboard.

---

## ✨ Features

### 🔬 Semantic Vector Trap — Zero-Day Detection
Integrates with Ollama (`nomic-embed-text`) to translate raw HTTP metadata into dense **768-dimensional embedding vectors**. If a payload's spatial geometry deviates beyond a strict Euclidean threshold (e.g., `> 12.0`), the trap snaps shut.

### ⚡ Adaptive O(1) Edge-Blocking — The Fast Path
Thread-safe in-memory caching (`sync.Map`) ensures that once an attacker is identified, all subsequent requests from their IP are dropped in **< 0.001ms** at the network edge — bypassing heavy AI inference entirely.

### 🕳️ Active AI Tarpit — Dynamic Honeypot
Instead of instantly dropping anomalous connections, Sentinel routes attackers to an AI-generated tarpit. Powered by `llama3.2`, it generates highly realistic fake vulnerabilities (fake SQL errors, HTML admin panels) and **forces the connection to hang**, draining the attacker's compute threads.

### 📡 Real-Time SOC Dashboard
A zero-dependency **Server-Sent Events (SSE)** broker in Go streams zero-day telemetry to a modern Next.js frontend — featuring automatic React state reconciliation, cryptographic incident ID generation, and `localStorage` persistence.

---

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| **Core Engine** | Go (Golang), Goroutines, Channels, `sync.Map`, Mutexes |
| **AI / Embeddings** | Ollama, `nomic-embed-text` |
| **AI / Generative** | Ollama, `llama3.2` |
| **Frontend SOC** | Next.js (React), Tailwind CSS, Server-Sent Events |

---

## 🚀 Getting Started

### Prerequisites

- **Go 1.21+**
- **Node.js 18+**
- **Ollama** running locally on port `11434`

### 1. Clone & Setup

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

Aura-Sentinel requires **three active processes**. Open three separate terminals:

**Terminal 1 — WAF Backend (Go)**
```bash
go run ./cmd/sentinel/main.go
# WAF listens on :8080 | Telemetry stream on :8081
```

**Terminal 2 — SOC Dashboard (Next.js)**
```bash
cd ui
npm install
npm run dev
# Dashboard available at http://localhost:3000
```

---

## 🧪 Testing & Simulation

Use **Terminal 3** to simulate real-world traffic and trigger the anomaly detection pipeline.

---

### ✅ Test Case A — Normal Traffic

Simulate a standard user requesting the homepage:

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/" -UseBasicParsing
```

**Expected result:**
- Euclidean distance calculated as low (e.g., `4.5`) — within normal threshold
- Request is proxied transparently to the backend
- Terminal prints: `[OK] Traffic Normal`

---

### 🔴 Test Case B — Zero-Day Attack (Hydra Brute Force)

Simulate an attacker probing for an admin configuration file:

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/admin/config.php" `
  -Method POST `
  -UserAgent "Kali-Linux-Hydra" `
  -UseBasicParsing
```

**Expected result:**
1. Mathematical deviation breaches the anomaly threshold (`> 12.0`)
2. WAF routes the connection into the **AI Tarpit**, stalling it for ~3 seconds
3. `llama3.2` generates a fake HTML response to deceive the attacker
4. SOC Dashboard flashes a **red alert** with an automated Incident Report

---

### ⛔ Test Case C — Fast-Path Edge Block

Run the **exact same payload** from Test Case B a second time:

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/admin/config.php" `
  -Method POST `
  -UserAgent "Kali-Linux-Hydra" `
  -UseBasicParsing
```

**Expected result:**
- Attacker's IP recognized instantly in the `sync.Map` cache
- Connection terminated with **HTTP 403** in O(1) time — AI engine bypassed entirely
- Terminal prints: `[X] Blocked IP at the Edge`

---

## 📁 Project Structure

```
Aura-Sentinel/
├── cmd/
│   └── sentinel/
│       └── main.go        # Entrypoint
├── internal/
│   ├── waf/               # Core WAF middleware & routing logic
│   ├── detector/          # Embedding + Euclidean anomaly engine
│   ├── tarpit/            # LLaMA-powered honeypot generator
│   ├── blocklist/         # sync.Map edge-blocking cache
│   └── telemetry/         # SSE broker & event streaming
├── ui/                    # Next.js SOC Dashboard
│   ├── app/
│   └── components/
├── go.mod
└── README.md
```

---

## 🔒 Security & Disclaimer

This project is built for **educational and research purposes**. The honeypot and tarpit modules are designed to study attacker behaviour in controlled environments. Do not deploy against production traffic without a thorough security review.

---

## 📄 License

MIT License — see [`LICENSE`](./LICENSE) for details.