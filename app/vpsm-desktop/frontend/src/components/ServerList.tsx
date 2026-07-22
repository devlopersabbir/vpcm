import React, { useState } from "react";
import { Server, Search, RefreshCw, Trash2, HardDrive, Star, Clock, Key, Shield, ChevronRight, Activity, Terminal } from "lucide-react";
import { maskHost } from "../utils/mask";

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
  os?: {
    os_family?: string;
    os_version?: string;
    architecture?: string;
  };
  hardware?: {
    cpu_model?: string;
    cpu_cores?: number;
    ram_total?: string;
    disk_total?: string;
  };
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
  onOpenTerminal: (srv: ServerView) => void;
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
  setActiveTab,
  onOpenTerminal,
}: ServerListProps) {
  const [filterTab, setFilterTab] = useState<"all" | "favorites" | "recents">("all");

  const getProviderTheme = (provider: string) => {
    const p = (provider || "").toLowerCase();
    if (p.includes("aws")) return {
      badge: "bg-orange-500/10 text-orange-400 border border-orange-500/20",
      accent: "bg-orange-500",
      glow: "shadow-[0_0_20px_rgba(249,115,22,0.15)]"
    };
    if (p.includes("digitalocean") || p.includes("do")) return {
      badge: "bg-blue-500/10 text-blue-405 border border-blue-500/20",
      accent: "bg-blue-500",
      glow: "shadow-[0_0_20px_rgba(59,130,246,0.15)]"
    };
    if (p.includes("gcp") || p.includes("google")) return {
      badge: "bg-emerald-500/10 text-emerald-400 border border-emerald-500/20",
      accent: "bg-emerald-500",
      glow: "shadow-[0_0_20px_rgba(16,185,129,0.15)]"
    };
    if (p.includes("vultr")) return {
      badge: "bg-sky-500/10 text-sky-400 border border-sky-500/20",
      accent: "bg-sky-500",
      glow: "shadow-[0_0_20px_rgba(56,189,248,0.15)]"
    };
    return {
      badge: "bg-slate-500/10 text-slate-400 border border-slate-700/20",
      accent: "bg-slate-600",
      glow: "shadow-none"
    };
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

  return (
    <div className="flex-1 flex flex-col p-8 overflow-y-auto w-full">
      {/* Upper Control Bar (Header) */}
      <div className="flex flex-col md:flex-row md:items-center justify-between pb-6 gap-4">
        <div>
          <div className="flex items-center space-x-3">
            <h1 className="text-xl font-black tracking-tight text-white uppercase">
              Server Inventory
            </h1>
            <div className="flex items-center space-x-1.5 bg-slate-900/60 px-2.5 py-1 rounded-lg border border-slate-800 text-[10px] text-slate-400 font-bold uppercase">
              <span>{servers.length} Total</span>
              <span className="text-slate-700">•</span>
              <span>{servers.filter(s => s.is_favorite).length} Starred</span>
            </div>
          </div>
          <p className="text-slate-400 text-xs mt-1.5">Manage connections, audit system libraries, and view specs.</p>
        </div>

        <div className="flex items-center space-x-3">
          <div className="relative">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
            <input
              type="text"
              placeholder="Search server name, IP, provider..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9 pr-4 py-2 w-64 rounded-xl text-xs placeholder-slate-500 focus:outline-none glass-input text-slate-200"
            />
          </div>
          <button
            onClick={fetchServers}
            className="p-2.5 bg-slate-900/40 hover:bg-slate-900 border border-slate-800/80 rounded-xl transition-all duration-200"
            title="Refresh Inventory"
          >
            <RefreshCw className={`h-4 w-4 ${loading ? "animate-spin text-cyan-400" : "text-slate-450"}`} />
          </button>
        </div>
      </div>

      {/* Modern Dashboard row (Horizontal, minimal, no nested boxes) */}
      {recentServers.length > 0 && (
        <div className="py-4 mb-6">
          <h3 className="text-[10px] font-black uppercase tracking-wider text-slate-500 mb-3 flex items-center space-x-1.5">
            <Clock className="h-3.5 w-3.5" />
            <span>Recently Connected</span>
          </h3>
          <div className="flex flex-wrap gap-4">
            {recentServers.map((s) => {
              const theme = getProviderTheme(s.provider);
              return (
                <div
                  key={s.id}
                  onClick={() => handleSelectServer(s)}
                  className={`flex-1 min-w-50 max-w-70 p-3.5 rounded-2xl flex items-center justify-between border border-slate-900 bg-slate-950/20 hover:bg-slate-900/30 cursor-pointer transition-all duration-200 ${theme.glow}`}
                >
                  <div className="truncate pr-3">
                    <h4 className="text-xs font-black text-slate-200 truncate flex items-center space-x-1">
                      <span>{s.name}</span>
                      {s.is_favorite && <Star className="h-2.5 w-2.5 text-amber-400 fill-amber-400" />}
                    </h4>
                    <p className="text-[9px] text-slate-500 font-mono truncate mt-0.5">{s.username}@{s.host}</p>
                  </div>
                  <span className={`inline-flex items-center px-2 py-0.5 rounded text-[8px] font-extrabold uppercase tracking-wide ${theme.badge}`}>
                    {s.provider ? s.provider.substring(0, 3) : "VPS"}
                  </span>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Styled Segmented Tabs Bar */}
      <div className="flex items-center space-x-1 bg-slate-950/40 p-1 rounded-xl border border-slate-900/80 self-start mb-6">
        <button
          onClick={() => setFilterTab("all")}
          className={`px-4 py-1.5 rounded-lg text-[10px] font-black tracking-wider uppercase transition-all duration-200 ${
            filterTab === "all" ? "bg-slate-900 text-cyan-400 shadow-sm" : "text-slate-400 hover:text-slate-200"
          }`}
        >
          All Nodes
        </button>
        <button
          onClick={() => setFilterTab("favorites")}
          className={`px-4 py-1.5 rounded-lg text-[10px] font-black tracking-wider uppercase transition-all duration-200 flex items-center space-x-1.5 ${
            filterTab === "favorites" ? "bg-slate-900 text-amber-400 shadow-sm" : "text-slate-400 hover:text-slate-200"
          }`}
        >
          <Star className="h-3 w-3" />
          <span>Starred</span>
        </button>
        <button
          onClick={() => setFilterTab("recents")}
          className={`px-4 py-1.5 rounded-lg text-[10px] font-black tracking-wider uppercase transition-all duration-200 flex items-center space-x-1.5 ${
            filterTab === "recents" ? "bg-slate-900 text-indigo-400 shadow-sm" : "text-slate-400 hover:text-slate-200"
          }`}
        >
          <Clock className="h-3 w-3" />
          <span>Recents</span>
        </button>
      </div>

      {/* Main Grid */}
      {displayServers.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center p-16 text-center">
          <Server className="h-12 w-12 text-slate-800 mb-4 stroke-1" />
          <h3 className="text-sm font-bold text-slate-400">No matching servers found</h3>
          <p className="text-slate-550 text-xs max-w-sm mt-1">
            {searchQuery ? "No results match your query." : "Add a new server from the sidebar to populate your inventory."}
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {displayServers.map((s) => {
            const theme = getProviderTheme(s.provider);
            const isScanning = scanningId === s.id;
            return (
              <div
                key={s.id}
                onClick={() => handleSelectServer(s)}
                className={`glass-card rounded-2xl flex flex-col justify-between group relative overflow-hidden cursor-pointer border border-slate-900 bg-slate-950/10 shadow-lg hover:border-slate-850/80 ${theme.glow}`}
              >
                {/* Highlight left line */}
                <div className={`absolute left-0 top-0 bottom-0 w-1 ${theme.accent}`} />

                {/* Main Card Body */}
                <div className="p-5 pl-6 flex-1 flex flex-col justify-between">
                  <div>
                    {/* Top Row: Provider & Inline Star Toggle */}
                    <div className="flex items-center justify-between mb-3.5">
                      <span className={`inline-flex items-center px-2 py-0.5 rounded text-[8px] font-extrabold uppercase tracking-wider ${theme.badge}`}>
                        {s.provider || "Generic VPS"}
                      </span>
                      <div className="flex items-center space-x-1.5">
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            onOpenTerminal(s);
                          }}
                          className="p-1.5 bg-blue-500/10 text-blue-400 hover:bg-blue-500/20 border border-blue-500/20 rounded-lg transition-colors"
                          title="Open SSH Terminal"
                        >
                          <Terminal className="h-3.5 w-3.5" />
                        </button>
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleToggleFavorite(s.id);
                          }}
                          className={`p-1 hover:text-amber-400 transition-colors ${
                            s.is_favorite ? "text-amber-400" : "text-slate-600"
                          }`}
                          title="Star Node"
                        >
                          <Star className={`h-4.5 w-4.5 ${s.is_favorite ? "fill-amber-400" : ""}`} />
                        </button>
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleDeleteServer(s.id, s.name);
                          }}
                          className="p-1 text-slate-600 hover:text-rose-450 transition-colors"
                          title="Delete Server"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </div>

                    {/* Server Name & Host */}
                    <h3 className="text-sm font-black text-slate-200 group-hover:text-cyan-400 transition-all duration-200">
                      {s.name}
                    </h3>
                    <p className="text-slate-500 text-[10px] font-mono mt-0.5">
                      {s.username}@{maskHost(s.host)}:{s.port}
                    </p>
                  </div>

                  {/* Inline Metrics (Clean Premium UX detail) */}
                  {s.os?.os_family || s.hardware?.ram_total ? (
                    <div className="mt-4 flex flex-col space-y-1.5 text-[10px] text-slate-400 border-t border-slate-900/60 pt-3.5">
                      {s.os?.os_family && (
                        <div className="flex justify-between items-center">
                          <span className="text-slate-550 font-bold uppercase text-[9px]">OS</span>
                          <span className="font-semibold text-slate-350">{s.os.os_family} {s.os.os_version || ""}</span>
                        </div>
                      )}
                      {s.hardware?.ram_total && (
                        <div className="flex justify-between items-center">
                          <span className="text-slate-550 font-bold uppercase text-[9px]">Specs</span>
                          <span className="font-semibold text-slate-350">
                            {s.hardware.cpu_cores ? `${s.hardware.cpu_cores} Cores • ` : ""}
                            {s.hardware.ram_total}
                          </span>
                        </div>
                      )}
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleScanServer(s.id);
                        }}
                        disabled={isScanning}
                        className="w-full flex items-center justify-center space-x-1.5 py-1.5 mt-3.5 rounded-xl border border-dashed border-slate-800/80 text-[10px] text-slate-500 hover:text-cyan-400 hover:border-cyan-500/30 transition-all duration-200"
                      >
                        <RefreshCw className={`h-3.5 w-3.5 ${isScanning ? "animate-spin text-cyan-400" : ""}`} />
                        <span>{isScanning ? "Re-collecting Specs..." : "Refresh Specifications"}</span>
                      </button>
                    </div>
                  ) : (
                    <div className="mt-4 border-t border-slate-900/60 pt-3.5">
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleScanServer(s.id);
                        }}
                        disabled={isScanning}
                        className="w-full flex items-center justify-center space-x-1.5 py-1.5 rounded-xl border border-dashed border-slate-800 text-[10px] text-slate-500 hover:text-cyan-400 hover:border-cyan-500/30 transition-all duration-200"
                      >
                        <Activity className={`h-3.5 w-3.5 ${isScanning ? "animate-spin text-cyan-400" : ""}`} />
                        <span>{isScanning ? "Collecting Specs..." : "Audit Specifications"}</span>
                      </button>
                    </div>
                  )}
                </div>

                {/* Footer Section */}
                <div className="px-5 py-3 bg-slate-950/20 border-t border-slate-900/50 flex items-center justify-between text-[10px]">
                  <span className="text-slate-500 font-bold uppercase tracking-wider flex items-center space-x-1.5">
                    <Shield className="h-3.5 w-3.5 text-slate-500" />
                    <span>{s.auth_type}</span>
                  </span>
                  <span className="text-cyan-400 font-extrabold flex items-center space-x-0.5 group-hover:translate-x-1 transition-transform uppercase tracking-wider">
                    <span>Inspect</span>
                    <ChevronRight className="h-3.5 w-3.5" />
                  </span>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
