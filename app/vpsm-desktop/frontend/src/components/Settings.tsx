import React from "react";
import { Database, Terminal, Shield, Monitor } from "lucide-react";

interface SettingsProps {
  config: any;
  setConfig: (cfg: any) => void;
  onSave: (e: React.FormEvent) => Promise<void>;
}

export default function Settings({ config, setConfig, onSave }: SettingsProps) {
  if (!config) return null;

  const maskCloudHost = (val: string) => {
    if (!val) return "";
    return val.replace(/187\.77\.151\.75/g, "***.***.***.***");
  };

  return (
    <div className="flex-1 p-8 overflow-y-auto w-full max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="border-b border-slate-900 pb-6">
        <h1 className="text-2xl font-black tracking-tight bg-linear-to-r from-white via-slate-200 to-slate-400 bg-clip-text text-transparent">
          Database & Engine Settings
        </h1>
        <p className="text-slate-400 text-xs mt-1 font-medium">
          Configure database storage drivers, central API endpoint URLs, and collector thread workers.
        </p>
      </div>

      <div className="bg-slate-900/35 border border-slate-900 rounded-2xl p-6 shadow-xl shadow-black/40 backdrop-blur-md">
        <form onSubmit={onSave} className="space-y-6">
          {/* Database Section */}
          <div className="pb-6 border-b border-slate-900">
            <h3 className="text-xs font-extrabold uppercase tracking-wider text-slate-300 flex items-center space-x-2 mb-4">
              <Database className="h-4 w-4 text-cyan-400" />
              <span>Database Connection & Driver</span>
            </h3>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                  Database Driver
                </label>
                <select
                  value={config.Database.Driver}
                  onChange={(e) => {
                    const updated = { ...config };
                    updated.Database.Driver = e.target.value;
                    setConfig(updated);
                  }}
                  className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input font-medium"
                >
                  <option value="sqlite">SQLite (Local File Storage)</option>
                  <option value="mongodb">MongoDB (Cloud / Remote Server)</option>
                </select>
              </div>

              {config.Database.Driver === "sqlite" ? (
                <div>
                  <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                    SQLite File Location
                  </label>
                  <input
                    type="text"
                    value={config.Database.Path}
                    onChange={(e) => {
                      const updated = { ...config };
                      updated.Database.Path = e.target.value;
                      setConfig(updated);
                    }}
                    className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 font-mono focus:outline-none glass-input"
                  />
                </div>
              ) : (
                <div>
                  <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                    Database Name
                  </label>
                  <input
                    type="text"
                    value={config.Database.Name}
                    onChange={(e) => {
                      const updated = { ...config };
                      updated.Database.Name = e.target.value;
                      setConfig(updated);
                    }}
                    className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input font-medium"
                  />
                </div>
              )}
            </div>

            {config.Database.Driver === "mongodb" && (
              <div className="mt-4">
                <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                  MongoDB Connection URI
                </label>
                <input
                  type="text"
                  placeholder="mongodb://localhost:27017"
                  value={maskCloudHost(config.Database.URI)}
                  onChange={(e) => {
                    const updated = { ...config };
                    updated.Database.URI = e.target.value;
                    setConfig(updated);
                  }}
                  className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 font-mono focus:outline-none glass-input"
                />
              </div>
            )}
          </div>

          {/* REST API Endpoint Settings */}
          <div>
            <h3 className="text-xs font-extrabold uppercase tracking-wider text-slate-300 flex items-center space-x-2 mb-4">
              <Terminal className="h-4 w-4 text-cyan-400" />
              <span>API Server & Collector Engine</span>
            </h3>

            <div className="space-y-4">
              <div>
                <div className="flex justify-between items-center mb-1.5">
                  <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400">
                    Central API Server Endpoint URL
                  </label>
                  <div className="flex space-x-2">
                    <button
                      type="button"
                      onClick={() => {
                        const updated = { ...config };
                        updated.API.GlobalURL = "http://localhost:8080";
                        setConfig(updated);
                      }}
                      className="text-[10px] font-bold text-cyan-400 hover:text-cyan-300 bg-cyan-950/40 border border-cyan-800/40 px-2.5 py-0.5 rounded-md transition-all"
                    >
                      Set Local Host
                    </button>
                    <button
                      type="button"
                      onClick={() => {
                        const updated = { ...config };
                        updated.API.GlobalURL = "http://187.77.151.75:8080";
                        setConfig(updated);
                      }}
                      className="text-[10px] font-bold text-indigo-400 hover:text-indigo-300 bg-indigo-950/40 border border-indigo-800/40 px-2.5 py-0.5 rounded-md transition-all"
                    >
                      Set Cloud Default
                    </button>
                  </div>
                </div>
                <input
                  type="text"
                  placeholder="http://localhost:8080"
                  value={maskCloudHost(config.API?.GlobalURL || "")}
                  onChange={(e) => {
                    const updated = { ...config };
                    updated.API.GlobalURL = e.target.value;
                    setConfig(updated);
                  }}
                  className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 font-mono focus:outline-none glass-input"
                />
                <p className="text-[11px] text-slate-500 mt-1 font-medium">
                  Cloud host IPs in settings preview are masked (`187.***.***.75`) for security.
                </p>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                    API Port
                  </label>
                  <input
                    type="number"
                    value={config.API.Port}
                    onChange={(e) => {
                      const updated = { ...config };
                      updated.API.Port = Number(e.target.value);
                      setConfig(updated);
                    }}
                    className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input font-mono"
                  />
                </div>
                <div>
                  <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                    Collector Worker Threads
                  </label>
                  <input
                    type="number"
                    value={config.Collector.Workers}
                    onChange={(e) => {
                      const updated = { ...config };
                      updated.Collector.Workers = Number(e.target.value);
                      setConfig(updated);
                    }}
                    className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input font-mono"
                  />
                </div>
              </div>
            </div>
          </div>

          <div className="pt-4 flex justify-end">
            <button
              type="submit"
              className="px-6 py-2.5 bg-linear-to-r from-cyan-600 to-indigo-600 hover:from-cyan-500 hover:to-indigo-500 text-white text-xs font-extrabold rounded-xl shadow-lg shadow-cyan-950/30 transition-all duration-200 border border-cyan-500/20"
            >
              Save Engine Configuration
            </button>
          </div>
        </form>
      </div>

      {/* Session & Onboarding Section */}
      <div className="bg-slate-900/35 border border-slate-900 rounded-2xl p-6 shadow-xl shadow-black/40 backdrop-blur-md flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <h3 className="text-xs font-extrabold uppercase tracking-wider text-slate-300 flex items-center space-x-2">
            <Monitor className="h-4 w-4 text-indigo-400" />
            <span>Session Profile & Onboarding State</span>
          </h3>
          <p className="text-slate-400 text-xs mt-1 font-medium">
            {localStorage.getItem("vpsm_user_session")
              ? `Connected as: ${JSON.parse(localStorage.getItem("vpsm_user_session")!).name} (${JSON.parse(localStorage.getItem("vpsm_user_session")!).provider})`
              : "Guest Session (Local Mode Active)"}
          </p>
        </div>
        <button
          type="button"
          onClick={() => {
            localStorage.removeItem("vpsm_onboarding_completed");
            localStorage.removeItem("vpsm_user_session");
            window.location.reload();
          }}
          className="px-4 py-2 bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-white text-xs font-extrabold rounded-xl border border-slate-700/60 transition-all duration-200 shrink-0"
        >
          Re-run Onboarding Wizard
        </button>
      </div>
    </div>
  );
}
