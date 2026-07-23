import React, { useState } from "react";
import {
  Server,
  Search,
  RefreshCw,
  Trash2,
  HardDrive,
  Star,
  Clock,
  Shield,
  ChevronRight,
  Activity,
  Terminal,
  LayoutGrid,
  List,
  Copy,
  Check,
  Plus,
  Cpu,
  Cloud,
  X,
  Sparkles
} from "lucide-react";
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
  const [providerFilter, setProviderFilter] = useState<string>("all");
  const [viewMode, setViewMode] = useState<"grid" | "table">("grid");
  const [copiedId, setCopiedId] = useState<number | null>(null);

  const getProviderTheme = (provider: string) => {
    const p = (provider || "").toLowerCase();
    if (p.includes("aws") || p.includes("amazon"))
      return {
        badge: "bg-orange-500/15 text-orange-400 border border-orange-500/30",
        accent: "bg-gradient-to-r from-orange-500 to-amber-500",
        glow: "hover:shadow-[0_0_25px_rgba(249,115,22,0.25)]",
        iconColor: "text-orange-400",
        dot: "bg-orange-400",
      };
    if (p.includes("digitalocean") || p.includes("do"))
      return {
        badge: "bg-blue-500/15 text-blue-400 border border-blue-500/30",
        accent: "bg-gradient-to-r from-blue-500 to-cyan-500",
        glow: "hover:shadow-[0_0_25px_rgba(59,130,246,0.25)]",
        iconColor: "text-blue-400",
        dot: "bg-blue-400",
      };
    if (p.includes("gcp") || p.includes("google"))
      return {
        badge: "bg-emerald-500/15 text-emerald-400 border border-emerald-500/30",
        accent: "bg-gradient-to-r from-emerald-500 to-teal-500",
        glow: "hover:shadow-[0_0_25px_rgba(16,185,129,0.25)]",
        iconColor: "text-emerald-400",
        dot: "bg-emerald-400",
      };
    if (p.includes("vultr"))
      return {
        badge: "bg-sky-500/15 text-sky-400 border border-sky-500/30",
        accent: "bg-gradient-to-r from-sky-500 to-indigo-500",
        glow: "hover:shadow-[0_0_25px_rgba(56,189,248,0.25)]",
        iconColor: "text-sky-400",
        dot: "bg-sky-400",
      };
    if (p.includes("hetzner"))
      return {
        badge: "bg-rose-500/15 text-rose-400 border border-rose-500/30",
        accent: "bg-gradient-to-r from-rose-500 to-red-500",
        glow: "hover:shadow-[0_0_25px_rgba(244,63,94,0.25)]",
        iconColor: "text-rose-400",
        dot: "bg-rose-400",
      };
    return {
      badge: "bg-slate-800/60 text-slate-300 border border-slate-700/50",
      accent: "bg-gradient-to-r from-cyan-500 to-blue-500",
      glow: "hover:shadow-[0_0_20px_rgba(6,182,212,0.15)]",
      iconColor: "text-cyan-400",
      dot: "bg-cyan-400",
    };
  };

  const handleCopySSH = (e: React.MouseEvent, s: ServerView) => {
    e.stopPropagation();
    const cmd = `ssh ${s.username}@${s.host} -p ${s.port}`;
    navigator.clipboard.writeText(cmd);
    setCopiedId(s.id);
    setTimeout(() => setCopiedId(null), 2000);
  };

  // Get unique providers list
  const uniqueProviders = Array.from(
    new Set(servers.map((s) => s.provider).filter(Boolean))
  );

  // Search Filter
  let displayServers = servers.filter(
    (s) =>
      s.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      s.host.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (s.provider && s.provider.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  // Provider Filter
  if (providerFilter !== "all") {
    displayServers = displayServers.filter(
      (s) => (s.provider || "").toLowerCase() === providerFilter.toLowerCase()
    );
  }

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
    .slice(0, 4);

  // Stats calculation
  const totalCores = servers.reduce(
    (acc, s) => acc + (s.hardware?.cpu_cores || 0),
    0
  );
  const starredCount = servers.filter((s) => s.is_favorite).length;

  return (
    <div className="flex-1 flex flex-col p-8 overflow-y-auto w-full space-y-6">
      {/* ── UPPER CONTROL BAR & HEADER ── */}
      <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4 pb-2 border-b border-slate-900/80">
        <div>
          <div className="flex items-center space-x-3">
            <div className="p-2.5 bg-gradient-to-br from-cyan-500/20 to-blue-600/10 rounded-2xl border border-cyan-500/30 shadow-[0_0_20px_rgba(6,182,212,0.2)]">
              <Server className="h-6 w-6 text-cyan-400" />
            </div>
            <div>
              <div className="flex items-center space-x-2">
                <h1 className="text-xl font-black tracking-tight text-white uppercase">
                  Server Inventory
                </h1>
                <span className="flex h-2 w-2 relative">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-cyan-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-2 w-2 bg-cyan-500"></span>
                </span>
              </div>
              <p className="text-slate-400 text-xs mt-0.5">
                Centralized node lobby • Monitor system specs & launch terminal sessions
              </p>
            </div>
          </div>
        </div>

        {/* Right actions: Search, Refresh, Add Server */}
        <div className="flex items-center space-x-3">
          <div className="relative">
            <Search className="absolute left-3.5 top-3 h-4 w-4 text-slate-500" />
            <input
              type="text"
              placeholder="Search server name, IP, provider..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10 pr-9 py-2.5 w-64 md:w-80 rounded-xl text-xs placeholder-slate-500 focus:outline-none glass-input text-slate-200"
            />
            {searchQuery && (
              <button
                onClick={() => setSearchQuery("")}
                className="absolute right-3 top-3 text-slate-500 hover:text-slate-300"
              >
                <X className="h-3.5 w-3.5" />
              </button>
            )}
          </div>

          <button
            onClick={fetchServers}
            className="p-2.5 bg-slate-900/60 hover:bg-slate-900 border border-slate-800 rounded-xl transition-all duration-200 text-slate-400 hover:text-cyan-400"
            title="Refresh Inventory"
          >
            <RefreshCw className={`h-4 w-4 ${loading ? "animate-spin text-cyan-400" : ""}`} />
          </button>

          <button
            onClick={() => setActiveTab("add-server")}
            className="flex items-center space-x-2 px-4 py-2.5 bg-gradient-to-r from-cyan-500 to-blue-600 hover:from-cyan-400 hover:to-blue-500 text-slate-950 font-black rounded-xl text-xs shadow-[0_0_20px_rgba(6,182,212,0.3)] transition-all duration-200 active:scale-95"
          >
            <Plus className="h-4 w-4 stroke-[3]" />
            <span>Add Server</span>
          </button>
        </div>
      </div>

      {/* ── HIGH IMPACT SUMMARY DASHBOARD WIDGETS ── */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {/* Total Nodes */}
        <div className="p-4 rounded-2xl bg-gradient-to-br from-slate-900/80 to-slate-950/60 border border-slate-800/80 flex items-center justify-between shadow-lg">
          <div>
            <p className="text-[10px] font-black uppercase text-slate-500 tracking-wider">Total Managed Nodes</p>
            <h3 className="text-2xl font-black text-white mt-1">{servers.length}</h3>
          </div>
          <div className="p-3 bg-cyan-500/10 rounded-xl border border-cyan-500/20 text-cyan-400">
            <Server className="h-5 w-5" />
          </div>
        </div>

        {/* Starred Nodes */}
        <div className="p-4 rounded-2xl bg-gradient-to-br from-slate-900/80 to-slate-950/60 border border-slate-800/80 flex items-center justify-between shadow-lg">
          <div>
            <p className="text-[10px] font-black uppercase text-slate-500 tracking-wider">Starred Nodes</p>
            <h3 className="text-2xl font-black text-amber-400 mt-1">{starredCount}</h3>
          </div>
          <div className="p-3 bg-amber-500/10 rounded-xl border border-amber-500/20 text-amber-400">
            <Star className="h-5 w-5 fill-amber-400/20" />
          </div>
        </div>

        {/* Total Managed Cores */}
        <div className="p-4 rounded-2xl bg-gradient-to-br from-slate-900/80 to-slate-950/60 border border-slate-800/80 flex items-center justify-between shadow-lg">
          <div>
            <p className="text-[10px] font-black uppercase text-slate-500 tracking-wider">Total CPU Cores</p>
            <h3 className="text-2xl font-black text-indigo-400 mt-1">{totalCores > 0 ? totalCores : "—"}</h3>
          </div>
          <div className="p-3 bg-indigo-500/10 rounded-xl border border-indigo-500/20 text-indigo-400">
            <Cpu className="h-5 w-5" />
          </div>
        </div>

        {/* Cloud Providers */}
        <div className="p-4 rounded-2xl bg-gradient-to-br from-slate-900/80 to-slate-950/60 border border-slate-800/80 flex items-center justify-between shadow-lg">
          <div>
            <p className="text-[10px] font-black uppercase text-slate-500 tracking-wider">Cloud Providers</p>
            <h3 className="text-2xl font-black text-emerald-400 mt-1">{uniqueProviders.length || 1}</h3>
          </div>
          <div className="p-3 bg-emerald-500/10 rounded-xl border border-emerald-500/20 text-emerald-400">
            <Cloud className="h-5 w-5" />
          </div>
        </div>
      </div>

      {/* ── RECENTLY CONNECTED QUICK LAUNCH BAR ── */}
      {recentServers.length > 0 && (
        <div className="p-4.5 rounded-2xl bg-slate-950/80 border border-slate-900 shadow-xl">
          <h3 className="text-[10px] font-black uppercase tracking-wider text-slate-400 mb-3 flex items-center space-x-1.5">
            <Sparkles className="h-3.5 w-3.5 text-cyan-400" />
            <span>Quick Launch — Recently Connected</span>
          </h3>
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-3">
            {recentServers.map((s) => {
              const theme = getProviderTheme(s.provider);
              return (
                <div
                  key={s.id}
                  onClick={() => onOpenTerminal(s)}
                  className={`p-3 rounded-xl border border-slate-900 bg-slate-900/60 hover:bg-slate-900 hover:border-slate-800 flex items-center justify-between cursor-pointer transition-all duration-200 group ${theme.glow}`}
                >
                  <div className="truncate pr-2">
                    <h4 className="text-xs font-black text-slate-200 group-hover:text-cyan-400 truncate flex items-center space-x-1">
                      <span>{s.name}</span>
                      {s.is_favorite && <Star className="h-2.5 w-2.5 text-amber-400 fill-amber-400" />}
                    </h4>
                    <p className="text-[9px] text-slate-500 font-mono truncate mt-0.5">{s.username}@{s.host}</p>
                  </div>
                  <div className="p-2 bg-cyan-500/10 text-cyan-400 rounded-lg group-hover:bg-cyan-500 group-hover:text-slate-950 transition-all">
                    <Terminal className="h-3.5 w-3.5" />
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* ── FILTER & VIEW MODE BAR ── */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        {/* Left: Tab segmented filter */}
        <div className="flex items-center space-x-1 bg-slate-950/60 p-1 rounded-xl border border-slate-900">
          <button
            onClick={() => setFilterTab("all")}
            className={`px-4 py-1.5 rounded-lg text-[10px] font-black tracking-wider uppercase transition-all duration-200 ${
              filterTab === "all" ? "bg-slate-900 text-cyan-400 shadow-sm" : "text-slate-400 hover:text-slate-200"
            }`}
          >
            All ({servers.length})
          </button>
          <button
            onClick={() => setFilterTab("favorites")}
            className={`px-4 py-1.5 rounded-lg text-[10px] font-black tracking-wider uppercase transition-all duration-200 flex items-center space-x-1.5 ${
              filterTab === "favorites" ? "bg-slate-900 text-amber-400 shadow-sm" : "text-slate-400 hover:text-slate-200"
            }`}
          >
            <Star className="h-3 w-3" />
            <span>Starred ({starredCount})</span>
          </button>
          <button
            onClick={() => setFilterTab("recents")}
            className={`px-4 py-1.5 rounded-lg text-[10px] font-black tracking-wider uppercase transition-all duration-200 flex items-center space-x-1.5 ${
              filterTab === "recents" ? "bg-slate-900 text-indigo-400 shadow-sm" : "text-slate-400 hover:text-slate-200"
            }`}
          >
            <Clock className="h-3 w-3" />
            <span>Recent</span>
          </button>
        </div>

        {/* Right: Provider filter & Grid/Table view toggles */}
        <div className="flex items-center space-x-3 self-end sm:self-auto">
          {uniqueProviders.length > 0 && (
            <select
              value={providerFilter}
              onChange={(e) => setProviderFilter(e.target.value)}
              className="bg-slate-950/60 border border-slate-900 rounded-xl px-3 py-1.5 text-[10px] font-bold text-slate-300 focus:outline-none focus:border-cyan-500/40"
            >
              <option value="all">All Providers</option>
              {uniqueProviders.map((p) => (
                <option key={p} value={p}>
                  {p}
                </option>
              ))}
            </select>
          )}

          <div className="flex items-center bg-slate-950/60 p-1 rounded-xl border border-slate-900 space-x-1">
            <button
              onClick={() => setViewMode("grid")}
              className={`p-1.5 rounded-lg transition-all ${
                viewMode === "grid" ? "bg-slate-900 text-cyan-400" : "text-slate-500 hover:text-slate-300"
              }`}
              title="Grid View"
            >
              <LayoutGrid className="h-4 w-4" />
            </button>
            <button
              onClick={() => setViewMode("table")}
              className={`p-1.5 rounded-lg transition-all ${
                viewMode === "table" ? "bg-slate-900 text-cyan-400" : "text-slate-500 hover:text-slate-300"
              }`}
              title="Table View"
            >
              <List className="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>

      {/* ── MAIN INVENTORY CONTENT (GRID vs TABLE) ── */}
      {displayServers.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center p-16 text-center border border-dashed border-slate-900 rounded-3xl bg-slate-950/20">
          <div className="p-4 bg-slate-900/60 rounded-full mb-4 border border-slate-800">
            <Server className="h-10 w-10 text-slate-600 stroke-[1.5]" />
          </div>
          <h3 className="text-base font-bold text-slate-300">No matching servers in inventory</h3>
          <p className="text-slate-500 text-xs max-w-sm mt-1 mb-6">
            {searchQuery ? "No results match your search query." : "Start by registering your first server node."}
          </p>
          <button
            onClick={() => setActiveTab("add-server")}
            className="flex items-center space-x-2 px-4 py-2 bg-cyan-500 hover:bg-cyan-400 text-slate-950 font-black rounded-xl text-xs transition-all"
          >
            <Plus className="h-4 w-4 stroke-[3]" />
            <span>Add Server Node</span>
          </button>
        </div>
      ) : viewMode === "grid" ? (
        /* GRID VIEW */
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {displayServers.map((s) => {
            const theme = getProviderTheme(s.provider);
            const isScanning = scanningId === s.id;
            const isCopied = copiedId === s.id;

            return (
              <div
                key={s.id}
                onClick={() => handleSelectServer(s)}
                className={`glass-card rounded-2xl flex flex-col justify-between group relative overflow-hidden cursor-pointer border border-slate-900/90 bg-slate-950/20 shadow-xl hover:border-slate-800 transition-all duration-300 ${theme.glow}`}
              >
                {/* Left accent border line */}
                <div className={`absolute left-0 top-0 bottom-0 w-1 ${theme.accent}`} />

                {/* Main Card Body */}
                <div className="p-5 pl-6 flex-1 flex flex-col justify-between">
                  <div>
                    {/* Header Row: Provider Badge & Action Icons */}
                    <div className="flex items-center justify-between mb-3.5">
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-lg text-[9px] font-black uppercase tracking-wider ${theme.badge}`}>
                        <span className={`w-1.5 h-1.5 rounded-full ${theme.dot} mr-1.5`} />
                        {s.provider || "Generic VPS"}
                      </span>

                      <div className="flex items-center space-x-1">
                        {/* Copy SSH Command */}
                        <button
                          onClick={(e) => handleCopySSH(e, s)}
                          className="p-1.5 text-slate-500 hover:text-slate-200 hover:bg-slate-900 rounded-lg transition-colors"
                          title="Copy SSH Command"
                        >
                          {isCopied ? <Check className="h-3.5 w-3.5 text-emerald-400" /> : <Copy className="h-3.5 w-3.5" />}
                        </button>

                        {/* Direct SSH Launch */}
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            onOpenTerminal(s);
                          }}
                          className="p-1.5 bg-cyan-500/10 text-cyan-400 hover:bg-cyan-500 hover:text-slate-950 border border-cyan-500/20 rounded-lg transition-all"
                          title="Launch SSH Terminal"
                        >
                          <Terminal className="h-3.5 w-3.5" />
                        </button>

                        {/* Favorite Star */}
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleToggleFavorite(s.id);
                          }}
                          className={`p-1.5 hover:text-amber-400 transition-colors ${
                            s.is_favorite ? "text-amber-400" : "text-slate-600"
                          }`}
                          title="Toggle Star"
                        >
                          <Star className={`h-4 w-4 ${s.is_favorite ? "fill-amber-400" : ""}`} />
                        </button>

                        {/* Delete */}
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleDeleteServer(s.id, s.name);
                          }}
                          className="p-1.5 text-slate-600 hover:text-rose-400 transition-colors"
                          title="Delete Server"
                        >
                          <Trash2 className="h-3.5 w-3.5" />
                        </button>
                      </div>
                    </div>

                    {/* Server Name & Host Address */}
                    <h3 className="text-sm font-black text-slate-100 group-hover:text-cyan-400 transition-all duration-200 truncate">
                      {s.name}
                    </h3>
                    <p className="text-slate-400 text-[11px] font-mono mt-0.5 truncate">
                      {s.username}@{maskHost(s.host)}:{s.port}
                    </p>
                  </div>

                  {/* System Specs or Scan Action */}
                  {s.os?.os_family || s.hardware?.ram_total ? (
                    <div className="mt-4 flex flex-col space-y-1.5 text-[10px] text-slate-400 border-t border-slate-900/80 pt-3">
                      {s.os?.os_family && (
                        <div className="flex justify-between items-center">
                          <span className="text-slate-500 font-bold uppercase text-[9px]">OS</span>
                          <span className="font-semibold text-slate-300">{s.os.os_family} {s.os.os_version || ""}</span>
                        </div>
                      )}
                      {s.hardware?.ram_total && (
                        <div className="flex justify-between items-center">
                          <span className="text-slate-500 font-bold uppercase text-[9px]">Specs</span>
                          <span className="font-semibold text-slate-300">
                            {s.hardware.cpu_cores ? `${s.hardware.cpu_cores} Cores • ` : ""}
                            {s.hardware.ram_total} RAM
                          </span>
                        </div>
                      )}
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleScanServer(s.id);
                        }}
                        disabled={isScanning}
                        className="w-full flex items-center justify-center space-x-1.5 py-1.5 mt-2 rounded-xl border border-dashed border-slate-800 text-[10px] text-slate-500 hover:text-cyan-400 hover:border-cyan-500/30 transition-all duration-200"
                      >
                        <RefreshCw className={`h-3 w-3 ${isScanning ? "animate-spin text-cyan-400" : ""}`} />
                        <span>{isScanning ? "Scanning..." : "Re-collect Specs"}</span>
                      </button>
                    </div>
                  ) : (
                    <div className="mt-4 border-t border-slate-900/80 pt-3">
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleScanServer(s.id);
                        }}
                        disabled={isScanning}
                        className="w-full flex items-center justify-center space-x-1.5 py-1.5 rounded-xl border border-dashed border-slate-800 text-[10px] text-slate-400 hover:text-cyan-400 hover:border-cyan-500/30 transition-all duration-200"
                      >
                        <Activity className={`h-3.5 w-3.5 ${isScanning ? "animate-spin text-cyan-400" : ""}`} />
                        <span>{isScanning ? "Scanning Specs..." : "Audit Specifications"}</span>
                      </button>
                    </div>
                  )}
                </div>

                {/* Footer Section */}
                <div className="px-5 py-2.5 bg-slate-950/40 border-t border-slate-900/60 flex items-center justify-between text-[10px]">
                  <span className="text-slate-500 font-bold uppercase tracking-wider flex items-center space-x-1">
                    <Shield className="h-3 w-3 text-slate-500" />
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
      ) : (
        /* TABLE VIEW (HIGH DENSITY) */
        <div className="rounded-2xl border border-slate-900 bg-slate-950/20 overflow-hidden shadow-xl">
          <table className="w-full text-left text-xs text-slate-300">
            <thead className="bg-slate-900/60 text-[10px] font-black uppercase text-slate-400 tracking-wider border-b border-slate-900">
              <tr>
                <th className="p-3.5 pl-5">Server Name</th>
                <th className="p-3.5">Host / Address</th>
                <th className="p-3.5">Provider</th>
                <th className="p-3.5">Auth Type</th>
                <th className="p-3.5">Specs & OS</th>
                <th className="p-3.5 text-right pr-5">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-900/60">
              {displayServers.map((s) => {
                const theme = getProviderTheme(s.provider);
                const isCopied = copiedId === s.id;

                return (
                  <tr
                    key={s.id}
                    onClick={() => handleSelectServer(s)}
                    className="hover:bg-slate-900/40 cursor-pointer transition-colors"
                  >
                    <td className="p-3.5 pl-5 font-bold text-slate-100 flex items-center space-x-2">
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleToggleFavorite(s.id);
                        }}
                        className="text-slate-600 hover:text-amber-400"
                      >
                        <Star className={`h-4 w-4 ${s.is_favorite ? "text-amber-400 fill-amber-400" : ""}`} />
                      </button>
                      <span className="hover:text-cyan-400 transition-colors">{s.name}</span>
                    </td>
                    <td className="p-3.5 font-mono text-slate-400 text-[11px]">
                      {s.username}@{maskHost(s.host)}:{s.port}
                    </td>
                    <td className="p-3.5">
                      <span className={`inline-flex items-center px-2 py-0.5 rounded text-[9px] font-extrabold uppercase ${theme.badge}`}>
                        {s.provider || "Generic VPS"}
                      </span>
                    </td>
                    <td className="p-3.5 text-slate-400 font-medium uppercase text-[10px]">{s.auth_type}</td>
                    <td className="p-3.5 text-slate-400 text-[11px]">
                      {s.hardware?.ram_total ? (
                        <span>
                          {s.hardware.cpu_cores ? `${s.hardware.cpu_cores}C • ` : ""}
                          {s.hardware.ram_total}
                        </span>
                      ) : (
                        <span className="text-slate-600">—</span>
                      )}
                    </td>
                    <td className="p-3.5 text-right pr-5">
                      <div className="flex items-center justify-end space-x-1.5">
                        <button
                          onClick={(e) => handleCopySSH(e, s)}
                          className="p-1.5 text-slate-500 hover:text-slate-200 rounded-lg hover:bg-slate-800"
                          title="Copy SSH Command"
                        >
                          {isCopied ? <Check className="h-3.5 w-3.5 text-emerald-400" /> : <Copy className="h-3.5 w-3.5" />}
                        </button>
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            onOpenTerminal(s);
                          }}
                          className="px-2.5 py-1 bg-cyan-500/10 hover:bg-cyan-500 text-cyan-400 hover:text-slate-950 border border-cyan-500/20 rounded-lg font-bold text-[10px] flex items-center space-x-1 transition-all"
                        >
                          <Terminal className="h-3 w-3" />
                          <span>SSH</span>
                        </button>
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleDeleteServer(s.id, s.name);
                          }}
                          className="p-1.5 text-slate-600 hover:text-rose-400 rounded-lg hover:bg-slate-800"
                          title="Delete Server"
                        >
                          <Trash2 className="h-3.5 w-3.5" />
                        </button>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
