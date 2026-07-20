import React, { useState } from "react";
import { Info, Cpu, HardDrive, Network as NetIcon, Code, Clock, X, Tag as TagIcon, Server, Database, Activity } from "lucide-react";

interface ServerDetailProps {
  selectedServer: any;
  setSelectedServer: (s: any) => void;
  logs: any[];
}

export default function ServerDetail({ selectedServer, setSelectedServer, logs }: ServerDetailProps) {
  const [detailTab, setDetailTab] = useState<"overview" | "hardware" | "os" | "network" | "software" | "logs">("overview");

  if (!selectedServer) return null;

  return (
    <div className="absolute inset-0 bg-slate-950/60 backdrop-blur-xs z-30 flex justify-end">
      <div className="w-full max-w-3xl bg-slate-900/98 h-full flex flex-col shadow-[-10px_0_30px_rgba(0,0,0,0.6)] border-l border-slate-800/80 relative">
        {/* Close Button */}
        <button
          onClick={() => setSelectedServer(null)}
          className="absolute top-6 right-6 p-2 bg-slate-800/80 hover:bg-slate-700 text-slate-400 hover:text-slate-200 rounded-full border border-slate-700/50 shadow-md shadow-black/30 transition-all duration-200 z-10"
        >
          <X className="h-4 w-4" />
        </button>

        {/* Header */}
        <div className="p-8 bg-slate-950/20 border-b border-slate-900/60">
          <span className="px-2 py-0.5 rounded-lg text-[9px] font-extrabold uppercase tracking-wider bg-cyan-950/30 text-cyan-400 border border-cyan-800/30 shadow-[0_2px_8px_rgba(6,182,212,0.1)]">
            {selectedServer.provider || "Generic VPS"}
          </span>
          <h2 className="text-xl font-black text-white mt-2 tracking-tight">{selectedServer.name}</h2>
          <p className="text-slate-500 font-mono text-xs mt-1">{selectedServer.username}@{selectedServer.host}:{selectedServer.port}</p>
        </div>

        {/* Tab Controls */}
        <div className="flex bg-slate-950/10 px-4 overflow-x-auto border-b border-slate-900/60">
          {[
            { id: "overview", label: "Overview", icon: Info },
            { id: "hardware", label: "Hardware", icon: Cpu },
            { id: "os", label: "OS System", icon: HardDrive },
            { id: "network", label: "Network", icon: NetIcon },
            { id: "software", label: `Apps (${selectedServer.software?.length || 0})`, icon: Code },
            { id: "logs", label: "Sessions", icon: Clock }
          ].map((t) => {
            const Icon = t.icon;
            return (
              <button
                key={t.id}
                onClick={() => setDetailTab(t.id as any)}
                className={`flex items-center space-x-1.5 px-4 py-3.5 text-[10px] font-bold uppercase tracking-wider transition-all duration-200 border-b-2 ${
                  detailTab === t.id
                    ? "text-cyan-400 bg-cyan-955/[0.02] border-cyan-400"
                    : "text-slate-400 hover:text-slate-200 border-transparent"
                }`}
              >
                <Icon className="h-3.5 w-3.5" />
                <span>{t.label}</span>
              </button>
            );
          })}
        </div>

        {/* Scrollable details */}
        <div className="flex-1 overflow-y-auto p-8">
          {/* OVERVIEW */}
          {detailTab === "overview" && (
            <div className="space-y-6">
              <div className="bg-slate-950/50 rounded-2xl p-5 border border-slate-800/60 shadow-md shadow-black/20">
                <h4 className="text-[10px] font-extrabold uppercase tracking-wider text-slate-500 mb-4 flex items-center space-x-1.5">
                  <TagIcon className="h-3.5 w-3.5 text-cyan-400" />
                  <span>VPS Identity Details</span>
                </h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs">
                  <div>
                    <span className="text-slate-500 block uppercase font-bold tracking-wider text-[9px] mb-0.5">UUID</span>
                    <span className="font-mono text-slate-300 font-semibold">{selectedServer.uuid}</span>
                  </div>
                  <div>
                    <span className="text-slate-500 block uppercase font-bold tracking-wider text-[9px] mb-0.5">Credential Method</span>
                    <span className="text-slate-300 font-semibold capitalize">{selectedServer.auth_type}</span>
                  </div>
                  <div>
                    <span className="text-slate-550 block uppercase font-bold tracking-wider text-[9px] mb-0.5">Created At</span>
                    <span className="text-slate-400 font-semibold">{new Date(selectedServer.created_at).toLocaleString()}</span>
                  </div>
                  <div>
                    <span className="text-slate-550 block uppercase font-bold tracking-wider text-[9px] mb-0.5">Last Audited</span>
                    <span className="text-slate-400 font-semibold">
                      {selectedServer.last_seen ? new Date(selectedServer.last_seen).toLocaleString() : "Never"}
                    </span>
                  </div>
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="bg-slate-950/40 rounded-2xl p-5 border border-slate-850 shadow-md shadow-black/20">
                  <span className="text-[9px] font-extrabold uppercase tracking-wider text-slate-500">Distribution</span>
                  <p className="text-md font-extrabold text-slate-200 mt-2">
                    {selectedServer.os?.os_family || "Unknown"} {selectedServer.os?.os_version || ""}
                  </p>
                </div>
                <div className="bg-slate-950/40 rounded-2xl p-5 border border-slate-850 shadow-md shadow-black/20">
                  <span className="text-[9px] font-extrabold uppercase tracking-wider text-slate-500">Resources</span>
                  <p className="text-md font-extrabold text-slate-200 mt-2">
                    {selectedServer.hardware?.ram_total || "—"} / {selectedServer.hardware?.disk_total || "—"}
                  </p>
                </div>
                <div className="bg-slate-950/40 rounded-2xl p-5 border border-slate-850 shadow-md shadow-black/20">
                  <span className="text-[9px] font-extrabold uppercase tracking-wider text-slate-500">IP address</span>
                  <p className="text-md font-extrabold text-cyan-400 mt-2 font-mono truncate" title={selectedServer.network?.public_ip}>
                    {selectedServer.network?.public_ip || selectedServer.host}
                  </p>
                </div>
              </div>
            </div>
          )}

          {/* HARDWARE */}
          {detailTab === "hardware" && (
            <div className="bg-slate-950/40 rounded-2xl p-6 space-y-4 border border-slate-800/60 shadow-md shadow-black/20">
              {!selectedServer.hardware ? (
                <div className="text-center py-8 text-slate-500">No hardware specification audited yet.</div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-y-4 gap-x-6 text-xs">
                  {[
                    { label: "CPU Model", val: selectedServer.hardware.cpu_model },
                    { label: "CPU Cores", val: `${selectedServer.hardware.cpu_cores || "—"} cores` },
                    { label: "Total RAM", val: selectedServer.hardware.ram_total },
                    { label: "Total Swap", val: selectedServer.hardware.swap_total },
                    { label: "Total Storage", val: selectedServer.hardware.disk_total },
                    { label: "Virtualization", val: selectedServer.hardware.virtualization },
                    { label: "Uptime", val: selectedServer.hardware.uptime, highlight: true },
                    { label: "Instance Type", val: selectedServer.hardware.instance_type }
                  ].map((x, i) => (
                    <div key={i} className="pb-3.5">
                      <span className="text-slate-500 block font-bold uppercase tracking-wider text-[9px]">{x.label}</span>
                      <span className={`font-semibold ${x.highlight ? "text-cyan-400" : "text-slate-200"}`}>{x.val || "—"}</span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* OS */}
          {detailTab === "os" && (
            <div className="bg-slate-950/40 rounded-2xl p-6 space-y-4 border border-slate-800/60 shadow-md shadow-black/20">
              {!selectedServer.os ? (
                <div className="text-center py-8 text-slate-500">No OS specification audited yet.</div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-y-4 gap-x-6 text-xs">
                  {[
                    { label: "OS Distribution", val: selectedServer.os.os_family },
                    { label: "Distribution Version", val: selectedServer.os.os_version },
                    { label: "Kernel Version", val: selectedServer.os.kernel_version, mono: true },
                    { label: "Architecture", val: selectedServer.os.architecture },
                    { label: "Timezone", val: selectedServer.os.timezone },
                    { label: "Package Manager", val: selectedServer.os.package_manager, mono: true }
                  ].map((x, i) => (
                    <div key={i} className="pb-3.5">
                      <span className="text-slate-500 block font-bold uppercase tracking-wider text-[9px]">{x.label}</span>
                      <span className={`font-semibold ${x.mono ? "font-mono" : ""} text-slate-200`}>{x.val || "—"}</span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* NETWORK */}
          {detailTab === "network" && (
            <div className="bg-slate-950/40 rounded-2xl p-6 space-y-4 border border-slate-800/60 shadow-md shadow-black/20">
              {!selectedServer.network ? (
                <div className="text-center py-8 text-slate-500">No network configuration audited yet.</div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-y-4 gap-x-6 text-xs">
                  {[
                    { label: "Hostname", val: selectedServer.network.hostname, mono: true },
                    { label: "Public IP", val: selectedServer.network.public_ip, mono: true },
                    { label: "Private IP", val: selectedServer.network.private_ip, mono: true },
                    { label: "MAC Address", val: selectedServer.network.mac_address, mono: true },
                    { label: "Cloud Region", val: selectedServer.network.region },
                    { label: "Availability Zone", val: selectedServer.network.availability_zone }
                  ].map((x, i) => (
                    <div key={i} className="pb-3.5">
                      <span className="text-slate-500 block font-bold uppercase tracking-wider text-[9px]">{x.label}</span>
                      <span className={`font-semibold ${x.mono ? "font-mono text-cyan-400" : ""} text-slate-200`}>{x.val || "—"}</span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* SOFTWARE */}
          {detailTab === "software" && (
            <div className="space-y-4">
              {!selectedServer.software || selectedServer.software.length === 0 ? (
                <div className="bg-slate-950/40 rounded-2xl p-6 text-center text-slate-500 border border-slate-800/60 shadow-md shadow-black/20">
                  No installed software packages scanned.
                </div>
              ) : (
                <div className="bg-slate-950/40 rounded-2xl overflow-hidden border border-slate-800/60 shadow-md shadow-black/20">
                  <table className="w-full text-left border-collapse text-xs">
                    <thead>
                      <tr className="bg-slate-950/60 text-slate-400 border-b border-slate-900">
                        <th className="px-4 py-3 font-extrabold uppercase tracking-wider text-[9px]">Package Name</th>
                        <th className="px-4 py-3 font-extrabold uppercase tracking-wider text-[9px]">Installed Version</th>
                      </tr>
                    </thead>
                    <tbody>
                      {selectedServer.software.map((pkg: any) => (
                        <tr key={pkg.id} className="hover:bg-slate-900/30 border-b border-slate-900/30 transition-colors last:border-0">
                          <td className="px-4 py-3 font-bold text-slate-300">{pkg.name}</td>
                          <td className="px-4 py-3 font-mono text-slate-400">{pkg.version}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          )}

          {/* SESSION LOGS */}
          {detailTab === "logs" && (
            <div className="space-y-4">
              {logs.length === 0 ? (
                <div className="bg-slate-950/40 rounded-2xl p-6 text-center text-slate-500 border border-slate-800/60 shadow-md shadow-black/20">
                  No connection sessions found for this server.
                </div>
              ) : (
                <div className="bg-slate-950/40 rounded-2xl overflow-hidden border border-slate-800/60 shadow-md shadow-black/20">
                  <table className="w-full text-left border-collapse text-[11px]">
                    <thead>
                      <tr className="bg-slate-950/60 text-slate-400 border-b border-slate-900">
                        <th className="px-4 py-3 font-extrabold uppercase text-[9px]">SSH User</th>
                        <th className="px-4 py-3 font-extrabold uppercase text-[9px]">Login Time</th>
                        <th className="px-4 py-3 font-extrabold uppercase text-[9px]">Duration</th>
                        <th className="px-4 py-3 font-extrabold uppercase text-[9px]">Status</th>
                      </tr>
                    </thead>
                    <tbody>
                      {logs.map((log: any) => (
                        <tr key={log.id} className="hover:bg-slate-900/30 border-b border-slate-900/30 transition-colors last:border-0">
                          <td className="px-4 py-3 font-bold text-slate-300 font-mono">{log.username}</td>
                          <td className="px-4 py-3 text-slate-400">{new Date(log.logged_in_at).toLocaleString()}</td>
                          <td className="px-4 py-3 text-slate-400 font-mono">{log.duration || "Active"}</td>
                          <td className="px-4 py-3">
                            <span
                              className={`inline-flex px-2 py-0.5 rounded-md text-[9px] font-extrabold uppercase tracking-wide border ${
                                log.status === "success"
                                  ? "bg-emerald-950/60 text-emerald-400 border-emerald-800/30"
                                  : "bg-rose-950/60 text-rose-400 border-rose-800/30"
                              }`}
                            >
                              {log.status}
                            </span>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
