import React from "react";
import { Server, Settings, PlusCircle, Activity, Database } from "lucide-react";

interface SidebarProps {
  activeTab: "servers" | "add-server" | "settings";
  setActiveTab: (tab: "servers" | "add-server" | "settings") => void;
  driver?: string;
}

export default function Sidebar({ activeTab, setActiveTab, driver }: SidebarProps) {
  return (
    <aside className="w-64 bg-slate-950 border-r border-slate-900 flex flex-col justify-between z-10 shadow-lg shadow-black/40">
      <div>
        {/* Brand */}
        <div className="p-6 flex items-center space-x-3 bg-slate-950/20 border-b border-slate-900/60">
          <div className="p-2 bg-linear-to-tr from-cyan-500/10 to-indigo-500/10 rounded-xl border border-cyan-500/10">
            <Activity className="h-5 w-5 text-cyan-400 animate-pulse" />
          </div>
          <span className="font-extrabold text-md tracking-wider bg-linear-to-r from-cyan-400 via-sky-400 to-indigo-400 bg-clip-text text-transparent uppercase">
            VPCM Panel
          </span>
        </div>

        {/* Navigation */}
        <nav className="p-4 space-y-2 mt-4">
          <button
            onClick={() => setActiveTab("servers")}
            className={`w-full flex items-center space-x-3 px-4 py-3 rounded-xl text-sm font-bold transition-all duration-300 relative border ${
              activeTab === "servers"
                ? "bg-slate-900/60 border-slate-800 text-cyan-400 shadow-[0_4px_12px_rgba(6,182,212,0.08)] before:absolute before:left-0 before:top-3 before:bottom-3 before:w-1 before:bg-cyan-400 before:rounded-r"
                : "text-slate-400 hover:bg-slate-900/35 hover:text-slate-200 border-transparent"
            }`}
          >
            <Server className="h-4 w-4" />
            <span>Server Inventory</span>
          </button>
          <button
            onClick={() => setActiveTab("add-server")}
            className={`w-full flex items-center space-x-3 px-4 py-3 rounded-xl text-sm font-bold transition-all duration-300 relative border ${
              activeTab === "add-server"
                ? "bg-slate-900/60 border-slate-800 text-cyan-400 shadow-[0_4px_12px_rgba(6,182,212,0.08)] before:absolute before:left-0 before:top-3 before:bottom-3 before:w-1 before:bg-cyan-400 before:rounded-r"
                : "text-slate-400 hover:bg-slate-900/35 hover:text-slate-200 border-transparent"
            }`}
          >
            <PlusCircle className="h-4 w-4" />
            <span>Add New Server</span>
          </button>
          <button
            onClick={() => setActiveTab("settings")}
            className={`w-full flex items-center space-x-3 px-4 py-3 rounded-xl text-sm font-bold transition-all duration-300 relative border ${
              activeTab === "settings"
                ? "bg-slate-900/60 border-slate-800 text-cyan-400 shadow-[0_4px_12px_rgba(6,182,212,0.08)] before:absolute before:left-0 before:top-3 before:bottom-3 before:w-1 before:bg-cyan-400 before:rounded-r"
                : "text-slate-400 hover:bg-slate-900/35 hover:text-slate-200 border-transparent"
            }`}
          >
            <Settings className="h-4 w-4" />
            <span>Settings & DB</span>
          </button>
        </nav>
      </div>

      {/* Footer info */}
      <div className="p-5 text-xs text-slate-500 flex justify-between items-center bg-slate-950/20 border-t border-slate-900/60">
        <span className="font-semibold tracking-wider text-slate-600">BUILD v1.0.0</span>
        {driver && (
          <span className="flex items-center space-x-1.5 uppercase text-cyan-400 font-bold bg-cyan-950/20 border border-cyan-800/20 px-2 py-0.5 rounded-lg text-[10px]">
            <Database className="h-3 w-3" />
            <span>{driver}</span>
          </span>
        )}
      </div>
    </aside>
  );
}
