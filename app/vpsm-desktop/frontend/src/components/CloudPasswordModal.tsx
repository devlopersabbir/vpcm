import React, { useState } from "react";
import { Lock, ShieldAlert, KeyRound, X, Loader2 } from "lucide-react";
import { VerifyCloudPassword } from "../../wailsjs/go/main/App";

interface CloudPasswordModalProps {
  isOpen: boolean;
  targetURL?: string;
  onClose: () => void;
  onSuccess: () => void;
}

export default function CloudPasswordModal({
  isOpen,
  targetURL = "http://127.0.0.1:8080",
  onClose,
  onSuccess,
}: CloudPasswordModalProps) {
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!password.trim()) {
      setError("Please enter the cloud guard password.");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      // Call backend Wails function to verify cloud password
      await VerifyCloudPassword(targetURL, password.trim());
      setLoading(false);
      setPassword("");
      onSuccess();
    } catch (err: any) {
      setLoading(false);
      setError(err?.toString() || "Invalid cloud access password");
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm animate-in fade-in duration-200">
      <div className="relative w-full max-w-md bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-2xl text-slate-100">
        {/* Header */}
        <div className="flex justify-between items-start mb-4">
          <div className="flex items-center space-x-3">
            <div className="p-2.5 rounded-xl bg-cyan-950/80 border border-cyan-700/50 text-cyan-400">
              <KeyRound className="h-6 w-6" />
            </div>
            <div>
              <h2 className="text-lg font-bold tracking-tight text-white flex items-center gap-2">
                Cloud Access Guard
                <Lock className="h-4 w-4 text-emerald-400" />
              </h2>
              <p className="text-xs text-slate-400 font-medium">
                Authorization required to connect with Cloud Server & DB
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="text-slate-400 hover:text-white p-1 rounded-lg hover:bg-slate-800 transition-colors"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        {/* Warning Banner */}
        <div className="mb-4 p-3 rounded-xl bg-cyan-950/30 border border-cyan-800/40 text-xs text-cyan-300/90 leading-relaxed font-medium">
          Accessing cloud resources requires administrator verification password configured on deployment.
        </div>

        {error && (
          <div className="mb-4 p-3 rounded-xl bg-rose-950/50 border border-rose-800/50 text-rose-300 text-xs flex items-center space-x-2">
            <ShieldAlert className="h-4 w-4 shrink-0" />
            <span>{error}</span>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
              Cloud Access Secret Password
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter cloud password"
              autoFocus
              className="w-full rounded-xl px-4 py-2.5 text-sm text-slate-100 bg-slate-950/80 border border-slate-800 focus:outline-none focus:border-cyan-500 transition-colors font-mono"
            />
          </div>

          <div className="flex justify-end space-x-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-xs font-semibold text-slate-400 hover:text-white hover:bg-slate-800 rounded-xl transition-all"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="px-5 py-2 text-xs font-bold text-slate-950 bg-cyan-400 hover:bg-cyan-300 active:scale-95 rounded-xl transition-all flex items-center space-x-1.5 shadow-md shadow-cyan-950 disabled:opacity-50"
            >
              {loading ? (
                <>
                  <Loader2 className="h-4 w-4 animate-spin" />
                  <span>Verifying...</span>
                </>
              ) : (
                <span>Authenticate & Access</span>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
