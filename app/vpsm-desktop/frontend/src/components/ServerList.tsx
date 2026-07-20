import React from "react";
import { Server, Search, RefreshCw, Trash2, Shield, HardDrive, Cpu, Layers } from "lucide-react";

interface ServerView {
  id: number;
  name: string;
  host: string;
  port: number;
  username: string;
  auth_type: string;
  provider: string;
  tags?: { id: number; name: string }[];
}

interface ServerListProps {
  servers: ServerView[];
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  loading: boolean;
  scanningId: number | null;
  fetchServers: () => void;
  handleScanServer: (id: number) => void;
  handleDeleteServer: (id: number, name: string) => void;
  handleSelectServer: (s: ServerView) => void;
  setActiveTab: (tab: "servers" | "add-server" | "settings") => void;
}

export default function ServerList({
  servers,
  searchQuery,
  setSearchQuery,
  loading,
  scanningId,
  fetchServers,
  handleScanServer,
  handleDeleteServer,
  handleSelectServer,
  setActiveTab
}: ServerListProps) {
  const filteredServers = servers.filter(
    (s) =>
      s.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      s.host.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (s.provider && s.provider.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  const getProviderColor = (provider: string) => {
    const p = (provider || "").toLowerCase();
    if (p.includes("aws")) return "from-amber-500/10 to-orange-500/5 text-orange-400";
    if (p.includes("digitalocean") || p.includes("do")) return "from-blue-500/10 to-indigo-500/5 text-blue-400";
    if (p.includes("gcp") || p.includes("google")) return "from-emerald-500/10 to-teal-500/5 text-emerald-400";
    if (p.includes("vultr")) return "from-sky-500/10 to-blue-500/5 text-sky-400";
    if (p.includes("linode") || p.includes("akamai")) return "from-green-500/10 to-emerald-500/5 text-green-400";
    return "from-slate-500/10 to-slate-600/5 text-slate-400";
  };

  return (
    <div className="flex-1 flex flex-col p-8 overflow-y-auto">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between pb-6 gap-4">
        <div>
          <h1 className="text-2xl font-extrabold tracking-tight bg-gradient-to-r from-white via-slate-100 to-slate-300 bg-clip-text text-transparent">
            Server Inventory
          </h1>
          <p className="text-slate-400 text-sm mt-1">Manage remote machines, retrieve specifications, and audit packages.</p>
        </div>
        <div className="flex items-center space-x-3">
          <div className="relative">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
            <input
              type="text"
              placeholder="Search by name, IP, provider..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9 pr-4 py-2 w-64 rounded-xl text-sm placeholder-slate-500 focus:outline-none glass-input text-slate-200"
            />
          </div>
          <button
            onClick={fetchServers}
            className="p-2.5 bg-slate-900/60 hover:bg-slate-900 rounded-xl transition-all duration-300"
            title="Refresh Inventory"
          >
            <RefreshCw className={`h-4 w-4 ${loading ? "animate-spin text-cyan-400" : "text-slate-400"}`} />
          </button>
        </div>
      </div>

      {/* Stats Dashboard Summary */}
      {servers.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 pt-6">
          <div className="glass-card rounded-2xl p-5 flex items-center justify-between">
            <div>
              <span className="text-[10px] font-extrabold uppercase tracking-wider text-slate-500">Total Nodes</span>
              <h2 className="text-3xl font-extrabold text-white mt-1">{servers.length}</h2>
            </div>
            <div className="p-3.5 bg-cyan-950/20 rounded-2xl text-cyan-400 shadow-[0_0_15px_-3px_rgba(6,182,212,0.2)]">
              <Layers className="h-5 w-5" />
            </div>
          </div>

          <div className="glass-card rounded-2xl p-5 flex items-center justify-between">
            <div>
              <span className="text-[10px] font-extrabold uppercase tracking-wider text-slate-500">AWS / GCP Nodes</span>
              <h2 className="text-3xl font-extrabold text-white mt-1">
                {servers.filter(s => { const p = (s.provider||"").toLowerCase(); return p.includes("aws") || p.includes("gcp") || p.includes("google"); }).length}
              </h2>
            </div>
            <div className="p-3.5 bg-indigo-950/20 rounded-2xl text-indigo-400 shadow-[0_0_15px_-3px_rgba(99,102,241,0.2)]">
              <Shield className="h-5 w-5" />
            </div>
          </div>

          <div className="glass-card rounded-2xl p-5 flex items-center justify-between">
            <div>
              <span className="text-[10px] font-extrabold uppercase tracking-wider text-slate-500">SSH Key Access</span>
              <h2 className="text-3xl font-extrabold text-white mt-1">
                {servers.filter(s => s.auth_type === "key").length}
              </h2>
            </div>
            <div className="p-3.5 bg-emerald-950/20 rounded-2xl text-emerald-400 shadow-[0_0_15px_-3px_rgba(16,185,129,0.2)]">
              <HardDrive className="h-5 w-5" />
            </div>
          </div>
        </div>
      )}

      {/* List */}
      {filteredServers.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center p-12 text-center">
          <Server className="h-16 w-16 text-slate-700 mb-4 stroke-1" />
          <h3 className="text-lg font-semibold text-slate-300">No servers found</h3>
          <p className="text-slate-500 text-sm max-w-sm mt-1">
            {searchQuery ? "Try refining your search query." : "Add your first VPS to get started monitoring specifications."}
          </p>
          {!searchQuery && (
            <button
              onClick={() => setActiveTab("add-server")}
              className="mt-4 px-5 py-2.5 bg-gradient-to-r from-cyan-600 to-indigo-600 hover:from-cyan-500 hover:to-indigo-500 text-white text-sm font-bold rounded-xl shadow-lg shadow-cyan-950/30 transition-all duration-200"
            >
              Add Server
            </button>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 pt-8">
          {filteredServers.map((s) => (
            <div
              key={s.id}
              onClick={() => handleSelectServer(s)}
              className="glass-card rounded-2xl p-5 flex flex-col justify-between group relative overflow-hidden cursor-pointer"
            >
              {/* Quick actions top-right */}
              <div className="absolute top-4 right-4 flex space-x-1 opacity-0 group-hover:opacity-100 transition-opacity z-10">
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    handleScanServer(s.id);
                  }}
                  disabled={scanningId === s.id}
                  className="p-2 bg-slate-900 rounded-lg text-slate-400 hover:text-cyan-400 transition-colors"
                  title="Audit Specifications"
                >
                  <RefreshCw className={`h-3.5 w-3.5 ${scanningId === s.id ? "animate-spin text-cyan-400" : ""}`} />
                </button>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    handleDeleteServer(s.id, s.name);
                  }}
                  className="p-2 bg-slate-900 rounded-lg text-slate-400 hover:text-rose-400 transition-colors"
                  title="Delete VPS"
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </button>
              </div>

              <div>
                {/* Custom Styled Provider Badge */}
                <span className={`inline-flex items-center px-2 py-0.5 rounded-lg text-[9px] font-extrabold uppercase tracking-wider bg-gradient-to-r ${getProviderColor(s.provider)} mb-3`}>
                  {s.provider || "Generic VPS"}
                </span>

                <h3 className="text-md font-extrabold text-slate-100 group-hover:text-cyan-400 transition-all duration-300">
                  {s.name}
                </h3>
                <p className="text-slate-500 text-xs font-mono mt-1 group-hover:text-slate-400 transition-colors">
                  {s.username}@{s.host}:{s.port}
                </p>

                {/* Tags */}
                {s.tags && s.tags.length > 0 && (
                  <div className="flex flex-wrap gap-1.5 mt-4">
                    {s.tags.map((t) => (
                      <span key={t.id} className="text-[9px] bg-cyan-500/10 text-cyan-300 px-2.5 py-0.5 rounded-full font-bold uppercase tracking-wide">
                        {t.name}
                      </span>
                    ))}
                  </div>
                )}
              </div>

              <div className="mt-6 pt-4 flex items-center justify-between">
                <span className="text-[10px] text-slate-500 font-semibold uppercase tracking-wider">
                  Access: <span className="text-slate-400 font-bold">{s.auth_type}</span>
                </span>
                <span className="text-xs font-bold text-cyan-400 flex items-center space-x-1 group-hover:translate-x-1 transition-transform">
                  <span>Open Detail</span>
                  <span>&rarr;</span>
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
