import {
  Server,
  Settings,
  PlusCircle,
  Activity,
  Database,
  Clock,
  Terminal,
} from "lucide-react";

interface SidebarProps {
  activeTab: "servers" | "history" | "add-server" | "terminal" | "settings";
  setActiveTab: (
    tab: "servers" | "history" | "add-server" | "terminal" | "settings",
  ) => void;
  driver?: string;
}

export default function Sidebar({
  activeTab,
  setActiveTab,
  driver,
}: SidebarProps) {
  return (
    <aside className="w-64 bg-slate-955/90 border-r border-slate-900/90 shadow-[6px_0_24px_-4px_rgba(0,0,0,0.6)] z-10 flex flex-col justify-between">
      <div>
        {/* Brand */}
        <div className="p-6 flex items-center space-x-3 bg-slate-950/20">
          <div className="p-2 bg-gradient-to-tr from-cyan-500/10 to-indigo-500/10 rounded-xl">
            <Activity className="h-5 w-5 text-cyan-400 animate-pulse" />
          </div>
          <span className="font-extrabold text-md tracking-wider bg-gradient-to-r from-cyan-400 via-sky-400 to-indigo-400 bg-clip-text text-transparent uppercase">
            VPCM Panel
          </span>
        </div>

        {/* Navigation */}
        <nav className="p-4 space-y-2 mt-4">
          <button
            onClick={() => setActiveTab("servers")}
            className={`w-full flex items-center space-x-3 px-4 py-3 rounded-xl text-sm font-bold transition-all duration-300 ${
              activeTab === "servers"
                ? "bg-gradient-to-r from-cyan-500/10 to-transparent text-cyan-400 shadow-[0_0_15px_-3px_rgba(6,182,212,0.15)]"
                : "text-slate-400 hover:bg-slate-900/35 hover:text-slate-200"
            }`}
          >
            <Server className="h-4 w-4" />
            <span>Server Inventory</span>
          </button>
          <button
            onClick={() => setActiveTab("history")}
            className={`w-full flex items-center space-x-3 px-4 py-3 rounded-xl text-sm font-bold transition-all duration-300 ${
              activeTab === "history"
                ? "bg-gradient-to-r from-cyan-500/10 to-transparent text-cyan-400 shadow-[0_0_15px_-3px_rgba(6,182,212,0.15)]"
                : "text-slate-400 hover:bg-slate-900/35 hover:text-slate-200"
            }`}
          >
            <Clock className="h-4 w-4" />
            <span>Session History</span>
          </button>
          <button
            onClick={() => setActiveTab("terminal")}
            className={`w-full flex items-center space-x-3 px-4 py-3 rounded-xl text-sm font-bold transition-all duration-300 ${
              activeTab === "terminal"
                ? "bg-gradient-to-r from-cyan-500/10 to-transparent text-cyan-400 shadow-[0_0_15px_-3px_rgba(6,182,212,0.15)]"
                : "text-slate-400 hover:bg-slate-900/35 hover:text-slate-200"
            }`}
          >
            <Terminal className="h-4 w-4" />
            <span>Terminal Setup</span>
          </button>
          <button
            onClick={() => setActiveTab("add-server")}
            className={`w-full flex items-center space-x-3 px-4 py-3 rounded-xl text-sm font-bold transition-all duration-300 ${
              activeTab === "add-server"
                ? "bg-gradient-to-r from-cyan-500/10 to-transparent text-cyan-400 shadow-[0_0_15px_-3px_rgba(6,182,212,0.15)]"
                : "text-slate-400 hover:bg-slate-900/35 hover:text-slate-200"
            }`}
          >
            <PlusCircle className="h-4 w-4" />
            <span>Add New Server</span>
          </button>
          <button
            onClick={() => setActiveTab("settings")}
            className={`w-full flex items-center space-x-3 px-4 py-3 rounded-xl text-sm font-bold transition-all duration-300 ${
              activeTab === "settings"
                ? "bg-gradient-to-r from-cyan-500/10 to-transparent text-cyan-400 shadow-[0_0_15px_-3px_rgba(6,182,212,0.15)]"
                : "text-slate-400 hover:bg-slate-900/35 hover:text-slate-200"
            }`}
          >
            <Settings className="h-4 w-4" />
            <span>Settings & DB</span>
          </button>
        </nav>
      </div>

      {/* Footer info */}
      <div className="p-5 text-xs text-slate-500 flex justify-between items-center bg-slate-950/10">
        <span className="font-semibold tracking-wider text-slate-600">
          BUILD v1.0.0
        </span>
        {driver && (
          <span className="flex items-center space-x-1.5 uppercase text-cyan-400 font-bold bg-cyan-950/20 px-2 py-0.5 rounded-lg text-[10px]">
            <Database className="h-3 w-3" />
            <span>{driver}</span>
          </span>
        )}
      </div>
    </aside>
  );
}
