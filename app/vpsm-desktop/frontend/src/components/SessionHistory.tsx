import React, { useState, useEffect } from "react";
import { Clock, Search, RefreshCw, AlertCircle } from "lucide-react";
import { GetConnectionHistory } from "../../wailsjs/go/main/App";
import { maskHost } from "../utils/mask";

export default function SessionHistory() {
  const [logs, setLogs] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");

  const fetchLogs = async () => {
    setLoading(true);
    try {
      const res = await GetConnectionHistory(0);
      if (res) setLogs(res);
    } catch (err) {
      console.error("Failed to fetch connection logs", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLogs();
  }, []);

  const filteredLogs = logs.filter(
    (l) =>
      l.server_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      l.host.toLowerCase().includes(searchQuery.toLowerCase()) ||
      l.username.toLowerCase().includes(searchQuery.toLowerCase()) ||
      l.status.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="flex-1 flex flex-col p-8 overflow-y-auto">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between pb-6 gap-4">
        <div>
          <h1 className="text-2xl font-black tracking-tight bg-linear-to-r from-white to-slate-350 bg-clip-text text-transparent">
            Session Connection History
          </h1>
          <p className="text-slate-400 text-sm mt-1">Audit log of all SSH connections, durations, and session events.</p>
        </div>
        <div className="flex items-center space-x-3">
          <div className="relative">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
            <input
              type="text"
              placeholder="Search by server, user, host..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9 pr-4 py-2 w-64 rounded-xl text-sm placeholder-slate-550 focus:outline-none glass-input text-slate-200"
            />
          </div>
          <button
            onClick={fetchLogs}
            className="p-2.5 bg-slate-900/60 hover:bg-slate-900 rounded-xl transition-all duration-300"
            title="Refresh History Logs"
          >
            <RefreshCw className={`h-4 w-4 ${loading ? "animate-spin text-cyan-400" : "text-slate-400"}`} />
          </button>
        </div>
      </div>

      {/* Connection logs table */}
      {filteredLogs.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center p-12 text-center">
          <Clock className="h-16 w-16 text-slate-800 mb-4 stroke-1" />
          <h3 className="text-lg font-semibold text-slate-300">No session logs found</h3>
          <p className="text-slate-500 text-sm max-w-sm mt-1">
            {searchQuery ? "Try adjusting your search query." : "Connection logs will appear once you connect to a server."}
          </p>
        </div>
      ) : (
        <div className="bg-slate-900/35 border border-slate-800/80 rounded-2xl overflow-hidden shadow-xl shadow-black/30 mt-6">
          <table className="w-full text-left border-collapse text-xs">
            <thead>
              <tr className="bg-slate-950/60 text-slate-400 border-b border-slate-900">
                <th className="px-6 py-4 font-extrabold uppercase tracking-wider text-[9px]">Server</th>
                <th className="px-6 py-4 font-extrabold uppercase tracking-wider text-[9px]">User & Host</th>
                <th className="px-6 py-4 font-extrabold uppercase tracking-wider text-[9px]">Login Time</th>
                <th className="px-6 py-4 font-extrabold uppercase tracking-wider text-[9px]">Duration</th>
                <th className="px-6 py-4 font-extrabold uppercase tracking-wider text-[9px]">Status</th>
                <th className="px-6 py-4 font-extrabold uppercase tracking-wider text-[9px]">Error / Event Detail</th>
              </tr>
            </thead>
            <tbody>
              {filteredLogs.map((log) => (
                <tr key={log.id} className="hover:bg-slate-900/30 border-b border-slate-900/30 last:border-0 transition-all duration-150">
                  <td className="px-6 py-4">
                    <span className="font-extrabold text-slate-200">{log.server_name}</span>
                  </td>
                  <td className="px-6 py-4 font-mono text-slate-400">
                    {log.username}@{maskHost(log.host)}
                  </td>
                  <td className="px-6 py-4 text-slate-400">
                    {new Date(log.logged_in_at).toLocaleString()}
                  </td>
                  <td className="px-6 py-4 text-slate-400 font-mono">
                    {log.duration || "Active Session"}
                  </td>
                  <td className="px-6 py-4">
                    <span
                      className={`inline-flex px-2 py-0.5 rounded-md text-[9px] font-extrabold uppercase tracking-wide border ${
                        log.status === "success"
                          ? "bg-emerald-950/60 text-emerald-400 border-emerald-800/30"
                          : log.status === "active"
                          ? "bg-cyan-950/60 text-cyan-400 border-cyan-800/30 animate-pulse"
                          : "bg-rose-950/60 text-rose-400 border-rose-800/30"
                      }`}
                    >
                      {log.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-rose-400 font-mono truncate max-w-xs" title={log.error_message}>
                    {log.error_message ? (
                      <span className="flex items-center space-x-1">
                        <AlertCircle className="h-3.5 w-3.5" />
                        <span>{log.error_message}</span>
                      </span>
                    ) : (
                      "—"
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
