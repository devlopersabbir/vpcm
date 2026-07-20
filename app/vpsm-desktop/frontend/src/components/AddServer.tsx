import React, { useState } from "react";
import { Server, ShieldAlert } from "lucide-react";

interface AddServerProps {
  onAdd: (name: string, host: string, port: number, username: string, authType: string, authSecret: string) => Promise<void>;
}

export default function AddServer({ onAdd }: AddServerProps) {
  const [formName, setFormName] = useState("");
  const [formHost, setFormHost] = useState("");
  const [formPort, setFormPort] = useState(22);
  const [formUsername, setFormUsername] = useState("root");
  const [formAuthType, setFormAuthType] = useState("key");
  const [formAuthSecret, setFormAuthSecret] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onAdd(formName, formHost, formPort, formUsername, formAuthType, formAuthSecret);
    setFormName("");
    setFormHost("");
    setFormPort(22);
    setFormUsername("root");
    setFormAuthSecret("");
  };

  return (
    <div className="flex-1 p-8 overflow-y-auto max-w-2xl mx-auto w-full">
      <h1 className="text-2xl font-black tracking-tight bg-gradient-to-r from-white to-slate-350 bg-clip-text text-transparent">
        Add New Server
      </h1>
      <p className="text-slate-400 text-sm mt-1">Register a remote node in your secure inventory database.</p>

      <form onSubmit={handleSubmit} className="mt-8 space-y-6 bg-slate-900/35 rounded-2xl p-6 backdrop-blur-md">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-450 mb-1.5">Server Name *</label>
            <input
              type="text"
              required
              placeholder="e.g. production-db-1"
              value={formName}
              onChange={(e) => setFormName(e.target.value)}
              className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input"
            />
          </div>
          <div>
            <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-450 mb-1.5">Host / IP *</label>
            <input
              type="text"
              required
              placeholder="e.g. 192.168.1.100"
              value={formHost}
              onChange={(e) => setFormHost(e.target.value)}
              className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input"
            />
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-450 mb-1.5">SSH Port *</label>
            <input
              type="number"
              required
              value={formPort}
              onChange={(e) => setFormPort(Number(e.target.value))}
              className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input"
            />
          </div>
          <div>
            <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-450 mb-1.5">SSH Username *</label>
            <input
              type="text"
              required
              value={formUsername}
              onChange={(e) => setFormUsername(e.target.value)}
              className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input"
            />
          </div>
          <div>
            <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-450 mb-1.5">Auth Type *</label>
            <select
              value={formAuthType}
              onChange={(e) => setFormAuthType(e.target.value)}
              className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input"
            >
              <option value="key">SSH Key (recommended)</option>
              <option value="password">Password</option>
            </select>
          </div>
        </div>

        <div>
          <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-450 mb-1.5">
            {formAuthType === "key" ? "SSH Private Key Path or Content" : "SSH Password"}
          </label>
          <textarea
            placeholder={formAuthType === "key" ? "e.g. ~/.ssh/id_rsa or paste BEGIN OPENSSH PRIVATE KEY..." : "Password"}
            value={formAuthSecret}
            onChange={(e) => setFormAuthSecret(e.target.value)}
            rows={4}
            className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 font-mono focus:outline-none glass-input"
          />
        </div>

        <div className="pt-4 flex justify-end">
          <button
            type="submit"
            className="px-6 py-2.5 bg-gradient-to-r from-cyan-600 to-indigo-600 hover:from-cyan-500 hover:to-indigo-500 text-white text-sm font-bold rounded-xl shadow-lg shadow-cyan-950/20 transition-all duration-200"
          >
            Register Server
          </button>
        </div>
      </form>
    </div>
  );
}
