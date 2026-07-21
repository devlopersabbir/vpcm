import React, { useState } from "react";
import { Database, Terminal, Eye, EyeOff } from "lucide-react";

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
    <div className="flex-1 p-8 overflow-y-auto max-w-2xl mx-auto w-full">
      <h1 className="text-2xl font-black tracking-tight bg-linear-to-r from-white to-slate-350 bg-clip-text text-transparent">
        System Settings
      </h1>
      <p className="text-slate-400 text-sm mt-1">Configure database connections and background collection settings.</p>

      <form onSubmit={onSave} className="mt-8 space-y-6 bg-slate-900/35 border border-slate-800/80 rounded-2xl p-6 shadow-xl shadow-black/30 backdrop-blur-md">
        {/* Database Section */}
        <div className="pb-6 border-b border-slate-900/60">
          <h3 className="text-xs font-extrabold uppercase tracking-wider text-slate-400 flex items-center space-x-1.5 mb-4">
            <Database className="h-4 w-4 text-cyan-400" />
            <span>Database Configuration</span>
          </h3>

          <div className="space-y-4">
            <div>
              <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-450 mb-1.5">Database Driver</label>
              <select
                value={config.Database.Driver}
                onChange={(e) => {
                  const updated = { ...config };
                  updated.Database.Driver = e.target.value;
                  setConfig(updated);
                }}
                className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input"
              >
                <option value="sqlite">SQLite (Local File)</option>
                <option value="mongodb">MongoDB (Cloud/Remote)</option>
              </select>
            </div>

            {config.Database.Driver === "sqlite" ? (
              <div>
                <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-455 mb-1.5">SQLite DB Path</label>
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
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-455 mb-1.5">MongoDB URI</label>
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
                <div>
                  <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-455 mb-1.5">Database Name</label>
                  <input
                    type="text"
                    value={config.Database.Name}
                    onChange={(e) => {
                      const updated = { ...config };
                      updated.Database.Name = e.target.value;
                      setConfig(updated);
                    }}
                    className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input"
                  />
                </div>
              </div>
            )}
          </div>
        </div>

        {/* REST API & Collector settings */}
        <div>
          <h3 className="text-xs font-extrabold uppercase tracking-wider text-slate-400 flex items-center space-x-1.5 mb-4">
            <Terminal className="h-4 w-4 text-cyan-400" />
            <span>General Preferences</span>
          </h3>

          <div className="space-y-4 mb-4">
            <div>
              <div className="flex justify-between items-center mb-1.5">
                <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400">
                  API Server Endpoint URL
                </label>
                <div className="flex space-x-2">
                  <button
                    type="button"
                    onClick={() => {
                      const updated = { ...config };
                      updated.API.GlobalURL = "http://localhost:8080";
                      setConfig(updated);
                    }}
                    className="text-[10px] font-semibold text-cyan-400 hover:text-cyan-300 bg-cyan-950/40 border border-cyan-800/40 px-2 py-0.5 rounded-md"
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
                    className="text-[10px] font-semibold text-indigo-400 hover:text-indigo-300 bg-indigo-950/40 border border-indigo-800/40 px-2 py-0.5 rounded-md"
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
              <p className="text-[11px] text-slate-500 mt-1">
                Cloud host IPs in settings preview are masked (`187.***.***.75`) for privacy.
              </p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-455 mb-1.5">API Port</label>
              <input
                type="number"
                value={config.API.Port}
                onChange={(e) => {
                  const updated = { ...config };
                  updated.API.Port = Number(e.target.value);
                  setConfig(updated);
                }}
                className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input"
              />
            </div>
            <div>
              <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-455 mb-1.5">Collector Workers</label>
              <input
                type="number"
                value={config.Collector.Workers}
                onChange={(e) => {
                  const updated = { ...config };
                  updated.Collector.Workers = Number(e.target.value);
                  setConfig(updated);
                }}
                className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input"
              />
            </div>
          </div>
        </div>

        <div className="pt-4 flex justify-end">
          <button
            type="submit"
            className="px-6 py-2.5 bg-linear-to-r from-cyan-600 to-indigo-600 hover:from-cyan-500 hover:to-indigo-500 text-white text-sm font-bold rounded-xl shadow-lg shadow-cyan-950/30 transition-all duration-200 border border-cyan-500/20 hover:border-cyan-400/40"
          >
            Save Configuration
          </button>
        </div>
      </form>

      {/* Account & Onboarding Reset Section */}
      <div className="mt-6 bg-slate-900/35 border border-slate-800/80 rounded-2xl p-6 shadow-xl shadow-black/30 backdrop-blur-md flex justify-between items-center">
        <div>
          <h3 className="text-xs font-extrabold uppercase tracking-wider text-slate-300 flex items-center space-x-1.5">
            <Terminal className="h-4 w-4 text-indigo-400" />
            <span>Onboarding & Session Profile</span>
          </h3>
          <p className="text-slate-500 text-xs mt-1">
            {localStorage.getItem("vpsm_user_session")
              ? `Connected as: ${JSON.parse(localStorage.getItem("vpsm_user_session")!).name} (${JSON.parse(localStorage.getItem("vpsm_user_session")!).provider})`
              : "Guest Session (No account connected)"}
          </p>
        </div>
        <button
          type="button"
          onClick={() => {
            localStorage.removeItem("vpsm_onboarding_completed");
            localStorage.removeItem("vpsm_user_session");
            window.location.reload();
          }}
          className="px-4 py-2 bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-white text-xs font-bold rounded-xl border border-slate-700/60 transition-all duration-200"
        >
          Re-run Onboarding
        </button>
      </div>
    </div>
  );
}
