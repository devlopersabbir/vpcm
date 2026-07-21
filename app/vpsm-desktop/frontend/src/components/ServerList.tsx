import React, { useState } from "react";
import { Server, Search, RefreshCw, Trash2, Shield, HardDrive, Layers, Star, Clock } from "lucide-react";

interface ServerView {
  id: number;
  name: string;
  host: string;
  port: number;
  username: string;
  auth_type: string;
  provider: string;
  is_favorite: boolean;
  last_seen?: string;
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
  handleToggleFavorite: (id: number) => void;
  setActiveTab: (tab: "servers" | "history" | "add-server" | "settings") => void;
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
  handleToggleFavorite,
  setActiveTab
}: ServerListProps) {
  const [filterTab, setFilterTab] = useState<"all" | "favorites" | "recents">("all");

  const getProviderColor = (provider: string) => {
    const p = (provider || "").toLowerCase();
    if (p.includes("aws")) return "from-orange-500/10 to-orange-500/5 text-orange-400 border border-orange-500/20 shadow-[0_2px_10px_-3px_rgba(249,115,22,0.15)]";
    if (p.includes("digitalocean") || p.includes("do")) return "from-blue-500/10 to-indigo-500/5 text-blue-400 border border-blue-500/20 shadow-[0_2px_10px_-3px_rgba(59,130,246,0.15)]";
    if (p.includes("gcp") || p.includes("google")) return "from-emerald-500/10 to-teal-500/5 text-emerald-400 border border-emerald-500/20 shadow-[0_2px_10px_-3px_rgba(16,185,129,0.15)]";
    if (p.includes("vultr")) return "from-sky-500/10 to-blue-500/5 text-sky-400 border border-sky-500/20 shadow-[0_2px_10px_-3px_rgba(56,189,248,0.15)]";
    if (p.includes("linode") || p.includes("akamai")) return "from-green-500/10 to-emerald-500/5 text-green-400 border border-green-500/20 shadow-[0_2px_10px_-3px_rgba(34,197,94,0.15)]";
    return "from-slate-500/10 to-slate-600/5 text-slate-400 border border-slate-700/25";
  };

  // Search Filter
  let displayServers = servers.filter(
    (s) =>
      s.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      s.host.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (s.provider && s.provider.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  // Tab Filter
  if (filterTab === "favorites") {
    displayServers = displayServers.filter((s) => s.is_favorite);
  } else if (filterTab === "recents") {
    displayServers = displayServers
      .filter((s) => s.last_seen)
      .sort((a, b) => new Date(b.last_seen!).getTime() - new Date(a.last_seen!).getTime());
  }

  // Top 3 Recent Connected
  const recentServers = servers
    .filter((s) => s.last_seen)
    .sort((a, b) => new Date(b.last_seen!).getTime() - new Date(a.last_seen!).getTime())
    .slice(0, 3);

  const formatLastSeen = (ts?: string) => {
    if (!ts) return "";
    const d = new Date(ts);
    if (isNaN(d.getTime())) return "";
    return d.toLocaleDateString() + " " + d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  return (
    <div className="flex-1 flex flex-col p-8 overflow-y-auto">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between pb-6 gap-4 border-b border-slate-900/60">
        <div>
          <h1 className="text-2xl font-extrabold tracking-tight bg-linear-to-r from-white via-slate-100 to-slate-300 bg-clip-text text-transparent">
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
              className="pl-9 pr-4 py-2 w-64 rounded-xl text-sm placeholder-slate-500 focus:outline-none glass-input text-slate-200 shadow-md shadow-black/20"
            />
          </div>
          <button
            onClick={fetchServers}
            className="p-2.5 bg-slate-900/60 hover:bg-slate-900 border border-slate-800 rounded-xl transition-all duration-300 shadow-md shadow-black/20 hover:border-slate-700"
            title="Refresh Inventory"
          >
            <RefreshCw className={`h-4 w-4 ${loading ? "animate-spin text-cyan-400" : "text-slate-400"}`} />
          </button>
        </div>
      </div>

      {/* Stats Dashboard Summary */}
      {servers.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 pt-6">
          <div className="glass-card rounded-2xl p-5 flex items-center justify-between border border-slate-800/60 shadow-md shadow-black/30">
            <div>
              <span className="text-[10px] font-extrabold uppercase tracking-wider text-slate-500">Total Nodes</span>
              <h2 className="text-3xl font-extrabold text-white mt-1">{servers.length}</h2>
            </div>
            <div className="p-3.5 bg-cyan-950/20 border border-cyan-800/20 rounded-2xl text-cyan-400 shadow-[0_0_15px_-3px_rgba(6,182,212,0.2)]">
              <Layers className="h-5 w-5" />
            </div>
          </div>

          <div className="glass-card rounded-2xl p-5 flex items-center justify-between border border-slate-800/60 shadow-md shadow-black/30">
            <div>
              <span className="text-[10px] font-extrabold uppercase tracking-wider text-slate-500">Starred Nodes</span>
              <h2 className="text-3xl font-extrabold text-white mt-1">
                {servers.filter((s) => s.is_favorite).length}
              </h2>
            </div>
            <div className="p-3.5 bg-indigo-950/20 border border-indigo-800/20 rounded-2xl text-indigo-400 shadow-[0_0_15px_-3px_rgba(99,102,241,0.2)]">
              <Star className="h-5 w-5 fill-indigo-400/20" />
            </div>
          </div>

          <div className="glass-card rounded-2xl p-5 flex items-center justify-between border border-slate-800/60 shadow-md shadow-black/30">
            <div>
              <span className="text-[10px] font-extrabold uppercase tracking-wider text-slate-500">Active SSH Keys</span>
              <h2 className="text-3xl font-extrabold text-white mt-1">
                {servers.filter((s) => s.auth_type === "key").length}
              </h2>
            </div>
            <div className="p-3.5 bg-emerald-950/20 border border-emerald-800/20 rounded-2xl text-emerald-400 shadow-[0_0_15px_-3px_rgba(16,185,129,0.2)]">
              <HardDrive className="h-5 w-5" />
            </div>
          </div>
        </div>
      )}

      {/* Horizontal row of 2-3 recently connected servers */}
      {recentServers.length > 0 && (
        <div className="mt-8">
          <h3 className="text-xs font-black uppercase tracking-wider text-slate-500 mb-3 flex items-center space-x-1.5">
            <Clock className="h-3.5 w-3.5" />
            <span>Recently Connected</span>
          </h3>
          <div className="flex flex-wrap gap-4">
            {recentServers.map((s) => (
              <div
                key={s.id}
                onClick={() => handleSelectServer(s)}
                className="glass-card flex-1 min-w-[200px] max-w-[280px] p-4 rounded-xl flex items-center justify-between border border-slate-900 bg-slate-950/40 hover:bg-slate-900/60 cursor-pointer transition-all duration-200"
              >
                <div className="truncate pr-4">
                  <h4 className="text-sm font-extrabold text-slate-200 truncate flex items-center space-x-1.5">
                    <span>{s.name}</span>
                    {s.is_favorite && <Star className="h-3 w-3 text-amber-400 fill-amber-400" />}
                  </h4>
                  <p className="text-[10px] text-slate-500 font-mono truncate mt-0.5">{s.username}@{s.host}</p>
                </div>
                <span className={`inline-flex px-2 py-0.5 rounded-lg text-[9px] font-extrabold uppercase tracking-wider bg-linear-to-r ${getProviderColor(s.provider)}`}>
                  {s.provider || "VPS"}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Tabs Filter Bar */}
      <div className="flex space-x-2 mt-8 border-b border-slate-900/60 pb-3">
        <button
          onClick={() => setFilterTab("all")}
          className={`px-4 py-1.5 rounded-lg text-xs font-extrabold tracking-wide uppercase transition-all duration-300 ${
            filterTab === "all"
              ? "bg-cyan-500/10 text-cyan-400 shadow-[0_0_15px_-3px_rgba(6,182,212,0.15)]"
              : "text-slate-400 hover:bg-slate-900/35 hover:text-slate-200"
          }`}
        >
          All Nodes
        </button>
        <button
          onClick={() => setFilterTab("favorites")}
          className={`px-4 py-1.5 rounded-lg text-xs font-extrabold tracking-wide uppercase transition-all duration-300 flex items-center space-x-1.5 ${
            filterTab === "favorites"
              ? "bg-amber-500/10 text-amber-400 shadow-[0_0_15px_-3px_rgba(245,158,11,0.15)]"
              : "text-slate-400 hover:bg-slate-900/35 hover:text-slate-200"
          }`}
        >
          <Star className="h-3.5 w-3.5" />
          <span>Starred</span>
        </button>
        <button
          onClick={() => setFilterTab("recents")}
          className={`px-4 py-1.5 rounded-lg text-xs font-extrabold tracking-wide uppercase transition-all duration-300 flex items-center space-x-1.5 ${
            filterTab === "recents"
              ? "bg-indigo-500/10 text-indigo-400 shadow-[0_0_15px_-3px_rgba(99,102,241,0.15)]"
              : "text-slate-400 hover:bg-slate-900/35 hover:text-slate-200"
          }`}
        >
          <Clock className="h-3.5 w-3.5" />
          <span>Recents</span>
        </button>
      </div>

      {/* List */}
      {displayServers.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center p-12 text-center">
          <Server className="h-16 w-16 text-slate-700 mb-4 stroke-1" />
          <h3 className="text-lg font-semibold text-slate-300">No servers found</h3>
          <p className="text-slate-500 text-sm max-w-sm mt-1">
            {searchQuery ? "Try refining your search query." : "Add a VPS or choose a different tab."}
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 pt-6">
          {displayServers.map((s) => (
            <div
              key={s.id}
              onClick={() => handleSelectServer(s)}
              className="glass-card rounded-2xl p-5 flex flex-col justify-between group relative overflow-hidden cursor-pointer border border-slate-800/60 shadow-md shadow-black/40"
            >
              {/* Quick actions top-right */}
              <div className="absolute top-4 right-4 flex space-x-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200 z-10">
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    handleToggleFavorite(s.id);
                  }}
                  className={`p-2 bg-slate-900 border border-slate-800 rounded-lg transition-colors ${
                    s.is_favorite ? "text-amber-400 hover:text-slate-400" : "text-slate-400 hover:text-amber-400"
                  }`}
                  title={s.is_favorite ? "Remove from Starred" : "Star Server"}
                >
                  <Star className={`h-3.5 w-3.5 ${s.is_favorite ? "fill-amber-400" : ""}`} />
                </button>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    handleScanServer(s.id);
                  }}
                  disabled={scanningId === s.id}
                  className="p-2 bg-slate-900 border border-slate-800 rounded-lg text-slate-400 hover:text-cyan-400 hover:border-cyan-500/30 transition-colors"
                  title="Audit Specifications"
                >
                  <RefreshCw className={`h-3.5 w-3.5 ${scanningId === s.id ? "animate-spin text-cyan-400" : ""}`} />
                </button>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    handleDeleteServer(s.id, s.name);
                  }}
                  className="p-2 bg-slate-900 border border-slate-800 rounded-lg text-slate-400 hover:text-rose-400 hover:border-rose-500/30 transition-colors"
                  title="Delete VPS"
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </button>
              </div>

              <div>
                {/* Custom Styled Provider Badge */}
                <div className="flex items-center justify-between mb-3 pr-20">
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-lg text-[9px] font-extrabold uppercase tracking-wider bg-linear-to-r ${getProviderColor(s.provider)}`}>
                    {s.provider || "Generic VPS"}
                  </span>
                </div>

                <h3 className="text-md font-extrabold text-slate-100 group-hover:text-cyan-400 transition-all duration-300 flex items-center space-x-1.5">
                  <span>{s.name}</span>
                  {s.is_favorite && <Star className="h-3.5 w-3.5 text-amber-400 fill-amber-400" />}
                </h3>
                <p className="text-slate-550 text-xs font-mono mt-1 group-hover:text-slate-400 transition-colors">
                  {s.username}@{s.host}:{s.port}
                </p>

                {/* Last Seen timestamp */}
                {s.last_seen && (
                  <p className="text-[10px] text-slate-500 font-mono mt-1.5 flex items-center space-x-1">
                    <Clock className="h-3 w-3" />
                    <span>Seen: {formatLastSeen(s.last_seen)}</span>
                  </p>
                )}

                {/* Tags */}
                {s.tags && s.tags.length > 0 && (
                  <div className="flex flex-wrap gap-1.5 mt-4">
                    {s.tags.map((t) => (
                      <span key={t.id} className="text-[9px] bg-cyan-500/10 text-cyan-300 px-2.5 py-0.5 rounded-full font-bold uppercase tracking-wide border border-cyan-400/10 shadow-[0_2px_8px_rgba(6,182,212,0.05)]">
                        {t.name}
                      </span>
                    ))}
                  </div>
                )}
              </div>

              <div className="mt-6 pt-4 flex items-center justify-between border-t border-slate-900/60">
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
