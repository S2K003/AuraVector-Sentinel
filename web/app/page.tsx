"use client";

import { useEffect, useState } from "react";

// ─── Types ────────────────────────────────────────────────────────────────────

interface Alert {
  id: string;
  ip: string;
  path: string;
  distance: number;
  report: string;
  timestamp: string;
}

// ─── Component ────────────────────────────────────────────────────────────────

export default function SOCDashboard() {
  const [alerts, setAlerts]           = useState<Alert[]>([]);
  const [isConnected, setIsConnected] = useState(false);

  // 1. ON MOUNT: Recover existing alerts from LocalStorage so they survive page refreshes
  useEffect(() => {
    const savedLogs = localStorage.getItem("aura-sentinel-logs");
    if (savedLogs) {
      try {
        setAlerts(JSON.parse(savedLogs));
      } catch (e) {
        console.error("Failed to parse saved logs");
      }
    }
  }, []);

  // 2. ON UPDATE: Sync state to LocalStorage every time a new anomaly is trapped
  useEffect(() => {
    localStorage.setItem("aura-sentinel-logs", JSON.stringify(alerts));
  }, [alerts]);

  // 3. THE SSE BROKER CONNECTION
  useEffect(() => {
    const eventSource = new EventSource("http://localhost:8081/soc/stream");

    eventSource.onopen  = () => setIsConnected(true);
    eventSource.onerror = () => setIsConnected(false);

    eventSource.onmessage = (event) => {
      try {
        const rawAlert = JSON.parse(event.data);

        // 4. Generate a unique ID and timestamp for React state reconciliation
        const newAlert: Alert = {
          ...rawAlert,
          id:        crypto.randomUUID(),
          timestamp: new Date().toLocaleTimeString(),
        };

        // Push newest to the top
        setAlerts((prev) => [newAlert, ...prev]);
      } catch (error) {
        console.error("Failed to parse incoming telemetry:", error);
      }
    };

    return () => eventSource.close();
  }, []);

  // Utility to wipe the database if it gets too cluttered
  const clearLogs = () => {
    setAlerts([]);
    localStorage.removeItem("aura-sentinel-logs");
  };

  // ─── Render ─────────────────────────────────────────────────────────────────
  return (
    <>
      {/* ── Scoped styles ──────────────────────────────────────────────────── */}
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:ital,wght@0,300;0,400;0,500;1,400&family=IBM+Plex+Sans:wght@300;400;500&display=swap');

        /* ── Design tokens ── */
        .soc-shell {
          --ink:        #08080a;
          --ink-2:      #0f0f13;
          --ink-3:      #17171d;
          --rule:       #22222a;
          --dim:        #32323e;
          --muted:      #4e4e60;
          --soft:       #7a7a90;
          --text:       #c8c8d8;
          --text-hi:    #e8e8f0;
          --red-hi:     #ff4444;
          --red-mid:    #cc2222;
          --red-lo:     #3d1010;
          --red-xlo:    #1a0808;
          --amber:      #d4820a;
          --amber-lo:   #2a1a04;
          --blue-hi:    #5b9cf6;
          --blue-lo:    #091428;
          --green-hi:   #22c55e;
          --green-lo:   #052010;
          --mono:       'IBM Plex Mono', 'Fira Code', monospace;
          --sans:       'IBM Plex Sans', system-ui, sans-serif;
          --border:     0.5px solid var(--rule);
        }

        /* ── Reset ── */
        .soc-shell *, .soc-shell *::before, .soc-shell *::after {
          box-sizing: border-box;
          margin: 0; padding: 0;
        }

        /* ── Shell ── */
        .soc-shell {
          background: var(--ink);
          color: var(--text);
          font-family: var(--mono);
          min-height: 100vh;
          display: flex;
          flex-direction: column;
        }

        /* ── Topbar ── */
        .soc-topbar {
          position: sticky;
          top: 0;
          z-index: 50;
          background: var(--ink-2);
          border-bottom: var(--border);
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 0 32px;
          height: 58px;
          flex-shrink: 0;
        }

        .soc-topbar-left { display: flex; align-items: center; gap: 20px; }

        .soc-wordmark {
          font-family: var(--mono);
          font-weight: 500;
          font-size: 13px;
          letter-spacing: 0.18em;
          text-transform: uppercase;
          color: var(--text-hi);
        }
        .soc-wordmark em { color: var(--red-hi); font-style: normal; }

        .soc-tagline {
          font-size: 10px;
          letter-spacing: 0.1em;
          color: var(--muted);
          text-transform: uppercase;
        }

        .soc-topbar-right { display: flex; align-items: center; gap: 16px; }

        /* Purge button */
        .soc-purge-btn {
          font-family: var(--mono);
          font-size: 10px;
          letter-spacing: 0.12em;
          text-transform: uppercase;
          color: var(--muted);
          background: transparent;
          border: 0.5px solid var(--dim);
          padding: 5px 12px;
          border-radius: 3px;
          cursor: pointer;
          transition: color 0.15s, border-color 0.15s, background 0.15s;
        }
        .soc-purge-btn:hover {
          color: var(--red-hi);
          border-color: var(--red-mid);
          background: var(--red-xlo);
        }

        /* Connection badge */
        .soc-conn-badge {
          display: flex;
          align-items: center;
          gap: 9px;
          background: var(--ink-3);
          border: var(--border);
          padding: 6px 14px;
          border-radius: 4px;
        }

        .soc-led-wrap {
          position: relative;
          width: 8px; height: 8px;
          flex-shrink: 0;
        }

        .soc-led-ping {
          position: absolute; inset: 0;
          border-radius: 50%;
          background: var(--green-hi);
          opacity: 0.6;
          animation: soc-ping 1.4s cubic-bezier(0,0,0.2,1) infinite;
        }

        .soc-led-core {
          position: absolute; inset: 1px;
          border-radius: 50%;
        }

        .soc-conn-label {
          font-family: var(--mono);
          font-size: 10px;
          font-weight: 500;
          letter-spacing: 0.14em;
          text-transform: uppercase;
        }

        /* Alert counter badge */
        .soc-alert-count {
          font-family: var(--mono);
          font-size: 10px;
          letter-spacing: 0.08em;
          color: var(--muted);
          background: var(--ink-3);
          border: var(--border);
          padding: 4px 10px;
          border-radius: 3px;
        }
        .soc-alert-count span { color: var(--red-hi); font-weight: 500; }

        /* ── Column ruler ── */
        .soc-ruler {
          display: flex;
          align-items: center;
          padding: 0 32px;
          height: 30px;
          border-bottom: var(--border);
          background: var(--ink-2);
          flex-shrink: 0;
        }
        .soc-ruler-col {
          font-size: 8px;
          letter-spacing: 0.14em;
          text-transform: uppercase;
          color: var(--muted);
        }

        /* ── Main feed ── */
        .soc-feed {
          flex: 1;
          padding: 24px 32px 40px;
          display: flex;
          flex-direction: column;
          gap: 12px;
        }

        /* ── Empty state ── */
        .soc-empty {
          flex: 1;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          gap: 14px;
          padding: 80px 0;
          color: var(--dim);
        }

        .soc-empty-icon {
          width: 48px; height: 48px;
          border: 0.5px dashed var(--dim);
          border-radius: 6px;
          display: flex; align-items: center; justify-content: center;
          color: var(--muted);
        }

        .soc-empty-label {
          font-size: 11px;
          letter-spacing: 0.1em;
          text-align: center;
          line-height: 1.7;
          max-width: 340px;
        }

        /* ── Alert card ── */
        .soc-card {
          background: var(--ink-2);
          border: 0.5px solid var(--rule);
          border-left: 2px solid var(--red-mid);
          border-radius: 5px;
          overflow: hidden;
          animation: soc-slide-in 0.35s ease both;
        }

        @keyframes soc-slide-in {
          from { opacity: 0; transform: translateY(-10px); }
          to   { opacity: 1; transform: translateY(0); }
        }

        .soc-card-head {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 10px 16px;
          border-bottom: var(--border);
          background: var(--red-xlo);
        }

        .soc-card-title {
          display: flex;
          align-items: center;
          gap: 10px;
          font-size: 11px;
          font-weight: 500;
          letter-spacing: 0.12em;
          text-transform: uppercase;
          color: var(--red-hi);
        }

        .soc-card-threat-icon {
          width: 18px; height: 18px;
          border: 0.5px solid var(--red-mid);
          border-radius: 3px;
          display: flex; align-items: center; justify-content: center;
          color: var(--red-hi);
          flex-shrink: 0;
        }

        .soc-card-time {
          font-size: 10px;
          letter-spacing: 0.06em;
          color: var(--muted);
        }

        /* ── Telemetry grid ── */
        .soc-tele-grid {
          display: grid;
          grid-template-columns: repeat(3, 1fr);
          gap: 0;
          border-bottom: var(--border);
        }

        .soc-tele-cell {
          padding: 12px 16px;
          border-right: var(--border);
        }
        .soc-tele-cell:last-child { border-right: none; }

        .soc-tele-key {
          font-size: 8px;
          letter-spacing: 0.18em;
          text-transform: uppercase;
          color: var(--muted);
          margin-bottom: 5px;
          display: flex;
          align-items: center;
          gap: 5px;
        }

        .soc-tele-key-dot {
          width: 4px; height: 4px;
          border-radius: 50%;
          background: var(--red-mid);
          flex-shrink: 0;
        }

        .soc-tele-val {
          font-family: var(--mono);
          font-size: 12px;
          font-weight: 500;
          color: var(--text-hi);
          letter-spacing: 0.04em;
          word-break: break-all;
        }

        .soc-tele-val.deviation {
          color: var(--red-hi);
          font-size: 15px;
        }

        /* ── Report block ── */
        .soc-report-wrap { padding: 14px 16px; }

        .soc-report-label {
          font-size: 8px;
          letter-spacing: 0.18em;
          text-transform: uppercase;
          color: var(--amber);
          margin-bottom: 8px;
          display: flex;
          align-items: center;
          gap: 7px;
        }

        .soc-report-label::after {
          content: '';
          flex: 1;
          height: 0.5px;
          background: var(--amber-lo);
        }

        .soc-report-pre {
          background: var(--blue-lo);
          border: 0.5px solid var(--dim);
          border-radius: 4px;
          padding: 14px 16px;
          font-family: var(--mono);
          font-size: 11px;
          line-height: 1.72;
          color: var(--blue-hi);
          white-space: pre-wrap;
          overflow-wrap: break-word;
        }

        /* ── Animations ── */
        @keyframes soc-ping {
          75%, 100% { transform: scale(2.2); opacity: 0; }
        }

        /* ── Scrollbar ── */
        .soc-shell::-webkit-scrollbar { width: 4px; }
        .soc-shell::-webkit-scrollbar-track { background: transparent; }
        .soc-shell::-webkit-scrollbar-thumb { background: var(--dim); border-radius: 2px; }

        /* ── Responsive ── */
        @media (max-width: 640px) {
          .soc-topbar    { padding: 0 16px; }
          .soc-tagline   { display: none; }
          .soc-ruler     { padding: 0 16px; }
          .soc-feed      { padding: 16px; }
          .soc-tele-grid { grid-template-columns: 1fr; }
          .soc-tele-cell { border-right: none; border-bottom: var(--border); }
          .soc-tele-cell:last-child { border-bottom: none; }
        }
      `}</style>

      {/* ── Shell ─────────────────────────────────────────────────────────── */}
      <div className="soc-shell">

        {/* ── Topbar ──────────────────────────────────────────────────────── */}
        <header className="soc-topbar">
          <div className="soc-topbar-left">
            <div className="soc-wordmark">
              <em>Aura</em>-Sentinel
            </div>
            <div className="soc-tagline">Real-Time Zero-Day Telemetry</div>
          </div>

          <div className="soc-topbar-right">
            {/* Alert count */}
            {alerts.length > 0 && (
              <div className="soc-alert-count">
                <span>{alerts.length}</span> event{alerts.length !== 1 ? 's' : ''} trapped
              </div>
            )}

            {/* Purge button */}
            <button className="soc-purge-btn" onClick={clearLogs}>
              Purge Logs
            </button>

            {/* Connection badge */}
            <div className="soc-conn-badge">
              <div className="soc-led-wrap">
                {isConnected && <div className="soc-led-ping" />}
                <div
                  className="soc-led-core"
                  style={{ background: isConnected ? 'var(--green-hi)' : 'var(--red-hi)' }}
                />
              </div>
              <div
                className="soc-conn-label"
                style={{ color: isConnected ? 'var(--green-hi)' : 'var(--red-hi)' }}
              >
                {isConnected ? 'System Armed' : 'Disconnected'}
              </div>
            </div>
          </div>
        </header>

        {/* ── Column ruler ──────────────────────────────────────────────── */}
        <div className="soc-ruler">
          <div className="soc-ruler-col" style={{ width: '34%' }}>Target Path</div>
          <div className="soc-ruler-col" style={{ width: '34%' }}>Attacker IP</div>
          <div className="soc-ruler-col" style={{ width: '32%' }}>Math Deviation</div>
        </div>

        {/* ── Alert feed ────────────────────────────────────────────────── */}
        <main className="soc-feed">

          {/* Empty state */}
          {alerts.length === 0 && (
            <div className="soc-empty">
              <div className="soc-empty-icon">
                <svg width="22" height="22" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="1.2"
                    d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
                </svg>
              </div>
              <div className="soc-empty-label">
                Awaiting inbound telemetry…<br />
                No anomalies detected in the current session.
              </div>
            </div>
          )}

          {/* Alert cards */}
          {alerts.map((alert) => (
            <div key={alert.id} className="soc-card">

              {/* Card header */}
              <div className="soc-card-head">
                <div className="soc-card-title">
                  <div className="soc-card-threat-icon">
                    <svg width="10" height="10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2.5"
                        d="M12 9v4m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" />
                    </svg>
                  </div>
                  Zero-Day Trapped &amp; IP Blocked
                </div>
                <div className="soc-card-time">{alert.timestamp}</div>
              </div>

              {/* Telemetry grid */}
              <div className="soc-tele-grid">
                <div className="soc-tele-cell">
                  <div className="soc-tele-key">
                    <div className="soc-tele-key-dot" />
                    Target Path
                  </div>
                  <div className="soc-tele-val">{alert.path}</div>
                </div>
                <div className="soc-tele-cell">
                  <div className="soc-tele-key">
                    <div className="soc-tele-key-dot" />
                    Attacker IP
                  </div>
                  <div className="soc-tele-val">{alert.ip}</div>
                </div>
                <div className="soc-tele-cell">
                  <div className="soc-tele-key">
                    <div className="soc-tele-key-dot" />
                    Math Deviation
                  </div>
                  <div className="soc-tele-val deviation">
                    {alert.distance.toFixed(2)}
                  </div>
                </div>
              </div>

              {/* AI incident report */}
              <div className="soc-report-wrap">
                <div className="soc-report-label">
                  Automated AI Incident Report
                </div>
                <pre className="soc-report-pre">{alert.report}</pre>
              </div>

            </div>
          ))}
        </main>
      </div>
    </>
  );
}