import React, { useState } from "react";
import { Database, Terminal, Monitor, CheckCircle, XCircle, RefreshCw, AlertTriangle, ShieldCheck } from "lucide-react";
import { TestDatabaseConnection, TestAPIConnection } from "../../wailsjs/go/main/App";

interface SettingsProps {
  config: any;
  setConfig: (cfg: any) => void;
  onSave: (e: React.FormEvent) => Promise<void>;
}

export default function Settings({ config, setConfig, onSave }: SettingsProps) {
  const [dbTest, setDbTest] = useState<{ tested: boolean; success: boolean; loading: boolean; message: string }>({
    tested: false,
    success: false,
    loading: false,
    message: "",
  });

  const [apiTest, setApiTest] = useState<{ tested: boolean; success: boolean; loading: boolean; message: string }>({
    tested: false,
    success: false,
    loading: false,
    message: "",
  });

  const [bypassTest, setBypassTest] = useState(false);
  const [saveError, setSaveError] = useState("");
  const [isSaving, setIsSaving] = useState(false);

  if (!config) return null;

  const handleTestDB = async () => {
    setDbTest({ tested: false, success: false, loading: true, message: "Testing database connection..." });
    setSaveError("");
    try {
      const driver = config.Database?.Driver || "sqlite";
      const path = config.Database?.Path || "";
      const uri = config.Database?.URI || "";
      const name = config.Database?.Name || "vpsm";

      const res = await TestDatabaseConnection(driver, path, uri, name);
      setDbTest({
        tested: true,
        success: !!res?.success,
        loading: false,
        message: res?.message || (res?.success ? "Database connection verified!" : "Connection failed"),
      });
      return !!res?.success;
    } catch (err: any) {
      setDbTest({
        tested: true,
        success: false,
        loading: false,
        message: err.message || "Failed to execute database connection test",
      });
      return false;
    }
  };

  const handleTestAPI = async () => {
    setApiTest({ tested: false, success: false, loading: true, message: "Testing API server endpoint..." });
    setSaveError("");
    try {
      const apiURL = config.API?.GlobalURL || "http://127.0.0.1:8080";
      const res = await TestAPIConnection(apiURL);
      setApiTest({
        tested: true,
        success: !!res?.success,
        loading: false,
        message: res?.message || (res?.success ? "API server connection verified!" : "Connection failed"),
      });
      return !!res?.success;
    } catch (err: any) {
      setApiTest({
        tested: true,
        success: false,
        loading: false,
        message: err.message || "Failed to execute API connection test",
      });
      return false;
    }
  };

  const handleSaveSubmitted = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaveError("");
    setIsSaving(true);

    try {
      if (!bypassTest) {
        let dbOk = dbTest.tested && dbTest.success;
        if (!dbOk) {
          dbOk = await handleTestDB();
        }

        if (!dbOk) {
          setSaveError("Database connection test failed. Fix database URI/path or enable 'Bypass Connection Test' to force save.");
          setIsSaving(false);
          return;
        }
      }

      await onSave(e);
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="flex-1 p-8 overflow-y-auto w-full max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="border-b border-slate-900 pb-6">
        <h1 className="text-2xl font-black tracking-tight bg-linear-to-r from-white via-slate-200 to-slate-400 bg-clip-text text-transparent">
          Database & Engine Settings
        </h1>
        <p className="text-slate-400 text-xs mt-1 font-medium">
          Configure database storage drivers, custom database connection URIs, and central API server endpoints.
        </p>
      </div>

      {saveError && (
        <div className="bg-rose-950/40 border border-rose-800/60 rounded-xl p-4 text-xs font-semibold text-rose-300 flex items-center space-x-3 shadow-lg">
          <AlertTriangle className="h-5 w-5 text-rose-400 shrink-0" />
          <span>{saveError}</span>
        </div>
      )}

      <div className="bg-slate-900/35 border border-slate-900 rounded-2xl p-6 shadow-xl shadow-black/40 backdrop-blur-md">
        <form onSubmit={handleSaveSubmitted} className="space-y-6">
          {/* Database Section */}
          <div className="pb-6 border-b border-slate-900">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-xs font-extrabold uppercase tracking-wider text-slate-300 flex items-center space-x-2">
                <Database className="h-4 w-4 text-cyan-400" />
                <span>Database Connection & Storage Driver</span>
              </h3>
              <button
                type="button"
                onClick={handleTestDB}
                disabled={dbTest.loading}
                className="px-3 py-1.5 bg-cyan-950/50 hover:bg-cyan-900/50 text-cyan-300 hover:text-white border border-cyan-800/50 text-xs font-extrabold rounded-lg transition-all flex items-center space-x-1.5"
              >
                {dbTest.loading ? (
                  <RefreshCw className="h-3.5 w-3.5 animate-spin text-cyan-400" />
                ) : (
                  <ShieldCheck className="h-3.5 w-3.5 text-cyan-400" />
                )}
                <span>{dbTest.loading ? "Testing..." : "Test Database Connection"}</span>
              </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                  Database Driver
                </label>
                <select
                  value={config.Database?.Driver || "sqlite"}
                  onChange={(e) => {
                    const updated = { ...config };
                    if (!updated.Database) updated.Database = {};
                    updated.Database.Driver = e.target.value;
                    setConfig(updated);
                    setDbTest({ tested: false, success: false, loading: false, message: "" });
                  }}
                  className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input font-medium"
                >
                  <option value="sqlite">SQLite (Local File Storage)</option>
                  <option value="mongodb">MongoDB (Custom Database Connection URL)</option>
                </select>
              </div>

              {config.Database?.Driver === "sqlite" ? (
                <div>
                  <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                    SQLite Database File Path
                  </label>
                  <input
                    type="text"
                    value={config.Database?.Path || ""}
                    onChange={(e) => {
                      const updated = { ...config };
                      if (!updated.Database) updated.Database = {};
                      updated.Database.Path = e.target.value;
                      setConfig(updated);
                      setDbTest({ tested: false, success: false, loading: false, message: "" });
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
                    value={config.Database?.Name || "vpsm"}
                    onChange={(e) => {
                      const updated = { ...config };
                      if (!updated.Database) updated.Database = {};
                      updated.Database.Name = e.target.value;
                      setConfig(updated);
                      setDbTest({ tested: false, success: false, loading: false, message: "" });
                    }}
                    className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input font-medium"
                  />
                </div>
              )}
            </div>

            {config.Database?.Driver === "mongodb" && (
              <div className="mt-4">
                <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                  MongoDB Connection URL / URI
                </label>
                <input
                  type="text"
                  placeholder="mongodb://127.0.0.1:27017 or mongodb+srv://user:pass@cluster.mongodb.net"
                  value={config.Database?.URI || ""}
                  onChange={(e) => {
                    const updated = { ...config };
                    if (!updated.Database) updated.Database = {};
                    updated.Database.URI = e.target.value;
                    setConfig(updated);
                    setDbTest({ tested: false, success: false, loading: false, message: "" });
                  }}
                  className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 font-mono focus:outline-none glass-input"
                />
              </div>
            )}

            {/* Test Results Badge */}
            {dbTest.tested && (
              <div
                className={`mt-4 p-3 rounded-xl border text-xs font-semibold flex items-center space-x-2.5 ${
                  dbTest.success
                    ? "bg-emerald-950/40 border-emerald-800/60 text-emerald-300"
                    : "bg-rose-950/40 border-rose-800/60 text-rose-300"
                }`}
              >
                {dbTest.success ? (
                  <CheckCircle className="h-4 w-4 text-emerald-400 shrink-0" />
                ) : (
                  <XCircle className="h-4 w-4 text-rose-400 shrink-0" />
                )}
                <span>{dbTest.message}</span>
              </div>
            )}
          </div>

          {/* REST API Endpoint Settings */}
          <div className="pb-6 border-b border-slate-900">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-xs font-extrabold uppercase tracking-wider text-slate-300 flex items-center space-x-2">
                <Terminal className="h-4 w-4 text-cyan-400" />
                <span>API Server & Engine Endpoints</span>
              </h3>
              <button
                type="button"
                onClick={handleTestAPI}
                disabled={apiTest.loading}
                className="px-3 py-1.5 bg-indigo-950/50 hover:bg-indigo-900/50 text-indigo-300 hover:text-white border border-indigo-800/50 text-xs font-extrabold rounded-lg transition-all flex items-center space-x-1.5"
              >
                {apiTest.loading ? (
                  <RefreshCw className="h-3.5 w-3.5 animate-spin text-indigo-400" />
                ) : (
                  <ShieldCheck className="h-3.5 w-3.5 text-indigo-400" />
                )}
                <span>{apiTest.loading ? "Testing..." : "Test API Server"}</span>
              </button>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                  Central API Server Endpoint URL
                </label>
                <input
                  type="text"
                  placeholder="http://127.0.0.1:8080 or https://your-custom-api.com"
                  value={config.API?.GlobalURL || ""}
                  onChange={(e) => {
                    const updated = { ...config };
                    if (!updated.API) updated.API = {};
                    updated.API.GlobalURL = e.target.value;
                    setConfig(updated);
                    setApiTest({ tested: false, success: false, loading: false, message: "" });
                  }}
                  className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 font-mono focus:outline-none glass-input"
                />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                    API Port
                  </label>
                  <input
                    type="number"
                    value={config.API?.Port || 8080}
                    onChange={(e) => {
                      const updated = { ...config };
                      if (!updated.API) updated.API = {};
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
                    value={config.Collector?.Workers || 5}
                    onChange={(e) => {
                      const updated = { ...config };
                      if (!updated.Collector) updated.Collector = {};
                      updated.Collector.Workers = Number(e.target.value);
                      setConfig(updated);
                    }}
                    className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input font-mono"
                  />
                </div>
              </div>
            </div>

            {/* Test Results Badge */}
            {apiTest.tested && (
              <div
                className={`mt-4 p-3 rounded-xl border text-xs font-semibold flex items-center space-x-2.5 ${
                  apiTest.success
                    ? "bg-emerald-950/40 border-emerald-800/60 text-emerald-300"
                    : "bg-rose-950/40 border-rose-800/60 text-rose-300"
                }`}
              >
                {apiTest.success ? (
                  <CheckCircle className="h-4 w-4 text-emerald-400 shrink-0" />
                ) : (
                  <XCircle className="h-4 w-4 text-rose-400 shrink-0" />
                )}
                <span>{apiTest.message}</span>
              </div>
            )}
          </div>

          <div className="pt-2 flex flex-col md:flex-row justify-between items-center gap-4">
            <label className="flex items-center space-x-2 text-xs text-slate-400 font-medium cursor-pointer select-none">
              <input
                type="checkbox"
                checked={bypassTest}
                onChange={(e) => setBypassTest(e.target.checked)}
                className="rounded border-slate-800 text-cyan-600 focus:ring-cyan-500 bg-slate-900"
              />
              <span>Bypass connection test verification (Force Save)</span>
            </label>

            <button
              type="submit"
              disabled={isSaving}
              className="w-full md:w-auto px-6 py-2.5 bg-linear-to-r from-cyan-600 to-indigo-600 hover:from-cyan-500 hover:to-indigo-500 text-white text-xs font-extrabold rounded-xl shadow-lg shadow-cyan-950/30 transition-all duration-200 border border-cyan-500/20 flex items-center justify-center space-x-2"
            >
              {isSaving ? (
                <RefreshCw className="h-4 w-4 animate-spin" />
              ) : (
                <ShieldCheck className="h-4 w-4" />
              )}
              <span>Save & Verify Configuration</span>
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
