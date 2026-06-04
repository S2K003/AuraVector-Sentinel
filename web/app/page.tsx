"use client";

import { useEffect, useState } from "react";

// 1. We expand the interface to include unique React IDs and timestamps
interface Alert {
  id: string;
  ip: string;
  path: string;
  distance: number;
  report: string;
  timestamp: string;
}

export default function SOCDashboard() {
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const [isConnected, setIsConnected] = useState(false);

  // 2. ON MOUNT: Recover existing alerts from LocalStorage so they survive page refreshes
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

  // 3. ON UPDATE: Sync the state to LocalStorage every time a new anomaly is trapped
  useEffect(() => {
    localStorage.setItem("aura-sentinel-logs", JSON.stringify(alerts));
  }, [alerts]);

  // 4. THE SSE BROKER CONNECTION
  useEffect(() => {
    const eventSource = new EventSource("http://localhost:8081/soc/stream");

    eventSource.onopen = () => setIsConnected(true);
    eventSource.onerror = () => setIsConnected(false);

    eventSource.onmessage = (event) => {
      try {
        const rawAlert = JSON.parse(event.data);
        
        // 5. Generate a unique ID and Timestamp for perfect React state reconciliation
        const newAlert: Alert = {
          ...rawAlert,
          id: crypto.randomUUID(), 
          timestamp: new Date().toLocaleTimeString(),
        };

        // Push newest to the top
        setAlerts((prev) => [newAlert, ...prev]);
      } catch (error) {
        console.error("Failed to parse incoming telemetry:", error);
      }
    };

    return () => {
      eventSource.close();
    };
  }, []);

  // Utility to wipe the database if it gets too cluttered
  const clearLogs = () => {
    setAlerts([]);
    localStorage.removeItem("aura-sentinel-logs");
  };

  return (
    <div className="min-h-screen bg-gray-950 text-gray-300 font-mono p-8 selection:bg-red-900 selection:text-white">
      {/* HEADER & CONNECTION STATUS */}
      <header className="border-b border-gray-800 pb-4 mb-8 flex flex-col md:flex-row justify-between items-start md:items-end gap-4">
        <div>
          <h1 className="text-3xl font-bold text-white tracking-widest uppercase">
            🛡️ Aura-Sentinel SOC
          </h1>
          <p className="text-gray-500 mt-2 text-sm">Real-Time Zero-Day Telemetry Dashboard</p>
        </div>
        
        <div className="flex items-center gap-6">
          <button 
            onClick={clearLogs}
            className="text-xs text-gray-500 hover:text-red-400 border border-gray-800 hover:border-red-900 px-3 py-1 rounded transition-colors"
          >
            PURGE LOGS
          </button>

          <div className="flex items-center gap-2 bg-gray-900 px-4 py-2 rounded-md border border-gray-800">
            <span className="relative flex h-3 w-3">
              {isConnected && (
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
              )}
              <span className={`relative inline-flex rounded-full h-3 w-3 ${isConnected ? 'bg-green-500' : 'bg-red-600'}`}></span>
            </span>
            <span className={`text-sm font-bold tracking-widest uppercase ${isConnected ? 'text-green-500' : 'text-red-500'}`}>
              {isConnected ? "System Armed" : "Disconnected"}
            </span>
          </div>
        </div>
      </header>

      {/* ANOMALY FEED */}
      <main className="space-y-6">
        {alerts.length === 0 ? (
          <div className="text-gray-600 italic border border-dashed border-gray-800 p-8 text-center rounded-lg">
            Awaiting inbound telemetry... no anomalies detected in the current session.
          </div>
        ) : (
          alerts.map((alert) => (
            // CRITICAL FIX: We now use the unique crypto ID instead of the array index!
            <div
              key={alert.id}
              className="bg-[#110a0a] border-l-4 border-red-600 p-6 rounded-md shadow-2xl animate-in slide-in-from-top-4 fade-in duration-500"
            >
              <div className="flex justify-between items-start mb-4">
                <h2 className="text-red-500 font-bold text-xl flex items-center gap-2">
                  ⚠️ ZERO-DAY TRAPPED & IP BLOCKED
                </h2>
                <span className="text-gray-600 text-xs">{alert.timestamp}</span>
              </div>
              
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6 text-sm bg-black/40 p-4 rounded border border-red-900/30">
                <div>
                  <span className="text-red-900 font-bold block mb-1 uppercase tracking-wider text-[10px]">Target Path</span>
                  <span className="text-white font-semibold">{alert.path}</span>
                </div>
                <div>
                  <span className="text-red-900 font-bold block mb-1 uppercase tracking-wider text-[10px]">Attacker IP</span>
                  <span className="text-white font-semibold">{alert.ip}</span>
                </div>
                <div>
                  <span className="text-red-900 font-bold block mb-1 uppercase tracking-wider text-[10px]">Math Deviation</span>
                  <span className="text-white font-semibold">{alert.distance.toFixed(2)}</span>
                </div>
              </div>

              <div>
                <span className="text-amber-500/80 font-bold text-xs uppercase tracking-widest mb-2 block">
                  [ Automated AI Incident Report ]
                </span>
                <pre className="bg-[#0a0f14] border border-blue-900/30 p-5 rounded text-blue-400 whitespace-pre-wrap text-sm leading-relaxed font-mono shadow-inner">
                  {alert.report}
                </pre>
              </div>
            </div>
          ))
        )}
      </main>
    </div>
  );
}