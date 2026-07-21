import React, { useState, useEffect } from "react";
import Sidebar from "./components/Sidebar";
import ServerList from "./components/ServerList";
import AddServer from "./components/AddServer";
import Settings from "./components/Settings";
import ServerDetail from "./components/ServerDetail";
import SessionHistory from "./components/SessionHistory";
import {
  GetServers,
  GetServer,
  AddServer as GoAddServer,
  DeleteServer,
  ScanServer,
  GetConfig,
  SaveConfig,
  GetConnectionHistory,
  ToggleFavorite
} from "../wailsjs/go/main/App";

export default function App() {
  const [activeTab, setActiveTab] = useState<"servers" | "history" | "add-server" | "settings">("servers");
  const [servers, setServers] = useState<any[]>([]);
  const [selectedServer, setSelectedServer] = useState<any | null>(null);
  const [logs, setLogs] = useState<any[]>([]);
  const [config, setConfig] = useState<any>(null);

  const [scanningId, setScanningId] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [statusMessage, setStatusMessage] = useState<{ type: "success" | "error"; text: string } | null>(null);

  useEffect(() => {
    fetchServers();
    fetchConfig();
  }, []);

  const showStatus = (type: "success" | "error", text: string) => {
    setStatusMessage({ type, text });
    setTimeout(() => setStatusMessage(null), 5000);
  };

  const fetchServers = async () => {
    setLoading(true);
    try {
      const res = await GetServers();
      if (res) setServers(res);
    } catch (err: any) {
      showStatus("error", "Failed to fetch servers: " + err);
    } finally {
      setLoading(false);
    }
  };

  const fetchConfig = async () => {
    try {
      const cfg = await GetConfig();
      if (cfg) setConfig(cfg);
    } catch (err: any) {
      showStatus("error", "Failed to load config: " + err);
    }
  };

  const handleAddServer = async (
    name: string,
    host: string,
    port: number,
    username: string,
    authType: string,
    authSecret: string
  ) => {
    try {
      await GoAddServer(name, host, port, username, authType, authSecret);
      showStatus("success", `Server "${name}" registered successfully!`);
      setActiveTab("servers");
      fetchServers();
    } catch (err: any) {
      showStatus("error", "Failed to add server: " + err);
    }
  };

  const handleDeleteServer = async (id: number, name: string) => {
    if (!confirm(`Are you sure you want to delete server "${name}"?`)) return;
    try {
      await DeleteServer(id);
      showStatus("success", `Server "${name}" deleted.`);
      if (selectedServer?.id === id) {
        setSelectedServer(null);
      }
      fetchServers();
    } catch (err: any) {
      showStatus("error", "Failed to delete server: " + err);
    }
  };

  const handleScanServer = async (id: number) => {
    setScanningId(id);
    try {
      await ScanServer(id);
      showStatus("success", "Inventory collection triggered!");
      setTimeout(async () => {
        fetchServers();
        if (selectedServer?.id === id) {
          const updated = await GetServer(id);
          if (updated) setSelectedServer(updated);
        }
        setScanningId(null);
      }, 3000);
    } catch (err: any) {
      showStatus("error", "Scan failed: " + err);
      setScanningId(null);
    }
  };

  const handleSelectServer = async (server: any) => {
    setSelectedServer(server);
    try {
      const detailed = await GetServer(server.id);
      if (detailed) setSelectedServer(detailed);
      const history = await GetConnectionHistory(server.id);
      if (history) setLogs(history);
    } catch (err: any) {
      console.error("Failed to load details/logs", err);
    }
  };

  const handleToggleFavorite = async (id: number) => {
    try {
      await ToggleFavorite(id);
      fetchServers();
    } catch (err: any) {
      showStatus("error", "Failed to toggle favorite: " + err);
    }
  };

  const handleSaveConfig = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!config) return;
    try {
      await SaveConfig(config);
      showStatus("success", "Configuration updated.");
      fetchConfig();
    } catch (err: any) {
      showStatus("error", "Failed to save configuration: " + err);
    }
  };

  return (
    <div className="flex h-screen bg-slate-950 text-slate-100 font-sans overflow-hidden">
      {/* Sidebar Navigation */}
      <Sidebar
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        driver={config?.Database?.Driver}
      />

      {/* Main Content Pane */}
      <main className="flex-1 flex flex-col relative overflow-hidden bg-radial-to-t from-slate-900 via-slate-950 to-slate-950">
        {/* Status Messages */}
        {statusMessage && (
          <div
            className={`absolute top-4 right-4 z-50 px-4 py-3 rounded-lg shadow-lg flex items-center space-x-2 border transition-all duration-300 ${
              statusMessage.type === "success"
                ? "bg-emerald-950/80 border-emerald-800 text-emerald-300"
                : "bg-rose-950/80 border-rose-800 text-rose-300"
            }`}
          >
            <span className="text-sm font-medium">{statusMessage.text}</span>
          </div>
        )}

        {/* Tab Selection */}
        {activeTab === "servers" && (
          <ServerList
            servers={servers}
            searchQuery={searchQuery}
            setSearchQuery={setSearchQuery}
            loading={loading}
            scanningId={scanningId}
            fetchServers={fetchServers}
            handleScanServer={handleScanServer}
            handleDeleteServer={handleDeleteServer}
            handleSelectServer={handleSelectServer}
            handleToggleFavorite={handleToggleFavorite}
            setActiveTab={setActiveTab}
          />
        )}

        {activeTab === "add-server" && <AddServer onAdd={handleAddServer} />}

        {activeTab === "history" && <SessionHistory />}

        {activeTab === "settings" && (
          <Settings
            config={config}
            setConfig={setConfig}
            onSave={handleSaveConfig}
          />
        )}

        {/* Detail Panel */}
        <ServerDetail
          selectedServer={selectedServer}
          setSelectedServer={setSelectedServer}
          logs={logs}
        />
      </main>
    </div>
  );
}
