import React, { useState, useEffect } from "react";
import { Sliders, CheckCircle2, Sparkles } from "lucide-react";

export default function TerminalSetup() {
  const [termPref, setTermPref] = useState<{
    font_size: number;
    font_family: string;
    background: string;
    foreground: string;
    opacity: number;
    blur: number;
    window_width: number;
    window_height: number;
    cursor_style: string;
    cursor_blink: boolean;
  }>({
    font_size: 14,
    font_family: 'Menlo, Monaco, "Courier New", monospace',
    background: "#0d1117",
    foreground: "#c9d1d9",
    opacity: 0.85,
    blur: 12,
    window_width: 900,
    window_height: 600,
    cursor_style: "block",
    cursor_blink: true,
  });

  const [prefLoading, setPrefLoading] = useState(false);
  const [prefSavedMsg, setPrefSavedMsg] = useState(false);

  useEffect(() => {
    if ((window as any).go?.main?.App?.GetTerminalPreference) {
      (window as any).go.main.App.GetTerminalPreference()
        .then((p: any) => {
          if (p) {
            setTermPref({
              font_size: p.font_size || 14,
              font_family:
                p.font_family || 'Menlo, Monaco, "Courier New", monospace',
              background: p.background || "#0d1117",
              foreground: p.foreground || "#c9d1d9",
              opacity: p.opacity !== undefined ? p.opacity : 0.85,
              blur: p.blur !== undefined ? p.blur : 12,
              window_width: p.window_width || 900,
              window_height: p.window_height || 600,
              cursor_style: p.cursor_style || "block",
              cursor_blink:
                p.cursor_blink !== undefined ? p.cursor_blink : true,
            });
          }
        })
        .catch(() => {});
    }
  }, []);

  const [toast, setToast] = useState<{ type: "success" | "error"; text: string } | null>(null);

  const showToast = (type: "success" | "error", text: string) => {
    setToast({ type, text });
    setTimeout(() => setToast(null), 4000);
  };

  const handleSaveTermPref = async (e: React.FormEvent) => {
    e.preventDefault();
    setPrefLoading(true);
    try {
      if ((window as any).go?.main?.App?.SaveTerminalPreference) {
        await (window as any).go.main.App.SaveTerminalPreference(termPref);
        setPrefSavedMsg(true);
        showToast("success", "Terminal settings saved successfully!");
        setTimeout(() => setPrefSavedMsg(false), 3000);
      }
    } catch (err: any) {
      console.error("Failed to save terminal preferences:", err);
      showToast("error", "Failed to save terminal settings: " + err);
    } finally {
      setPrefLoading(false);
    }
  };

  return (
    <div className="flex-1 p-8 overflow-y-auto w-full max-w-7xl mx-auto space-y-6 relative">
      {/* Toast Notification */}
      {toast && (
        <div
          className={`fixed top-6 right-6 z-50 px-4 py-3 rounded-xl shadow-2xl flex items-center space-x-2.5 border backdrop-blur-md transition-all duration-300 ${
            toast.type === "success"
              ? "bg-emerald-950/90 border-emerald-500/50 text-emerald-300 shadow-emerald-950/50"
              : "bg-rose-950/90 border-rose-500/50 text-rose-300 shadow-rose-950/50"
          }`}
        >
          <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0" />
          <span className="text-xs font-bold font-mono">{toast.text}</span>
        </div>
      )}

      {/* Header */}
      <div className="border-b border-slate-900 pb-6">
        <h1 className="text-2xl font-black tracking-tight bg-linear-to-r from-white via-slate-200 to-slate-400 bg-clip-text text-transparent">
          Terminal Customization & Aesthetics
        </h1>
        <p className="text-slate-400 text-xs mt-1 font-medium">
          Configure font size, window size, background opacity, blur effects,
          and cursor styling for your custom terminal.
        </p>
      </div>

      {/* Split Screen Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
        {/* Left Column: Controls Form */}
        <div className="lg:col-span-7 bg-slate-900/35 border border-slate-900 rounded-2xl p-6 shadow-xl shadow-black/40 backdrop-blur-md space-y-6">
          <div className="flex items-center justify-between border-b border-slate-900 pb-4">
            <h3 className="text-xs font-extrabold uppercase tracking-wider text-slate-300 flex items-center space-x-2">
              <Sliders className="h-4 w-4 text-cyan-400" />
              <span>Terminal Preferences</span>
            </h3>
            <span className="text-[10px] font-bold text-cyan-400 bg-cyan-950/40 border border-cyan-800/40 px-2 py-0.5 rounded-md">
              Realtime Sync
            </span>
          </div>

          <form onSubmit={handleSaveTermPref} className="space-y-6">
            {/* Font Size & Window Dimensions */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <div>
                <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                  Font Size (px)
                </label>
                <input
                  type="number"
                  min={8}
                  max={36}
                  value={termPref.font_size}
                  onChange={(e) =>
                    setTermPref({
                      ...termPref,
                      font_size: Number(e.target.value),
                    })
                  }
                  className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input font-mono"
                />
              </div>
              <div>
                <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                  Width (px)
                </label>
                <input
                  type="number"
                  min={600}
                  max={2560}
                  value={termPref.window_width}
                  onChange={(e) =>
                    setTermPref({
                      ...termPref,
                      window_width: Number(e.target.value),
                    })
                  }
                  className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input font-mono"
                />
              </div>
              <div>
                <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                  Height (px)
                </label>
                <input
                  type="number"
                  min={400}
                  max={1600}
                  value={termPref.window_height}
                  onChange={(e) =>
                    setTermPref({
                      ...termPref,
                      window_height: Number(e.target.value),
                    })
                  }
                  className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input font-mono"
                />
              </div>
            </div>

            {/* Sliders for Opacity & Blur */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 bg-slate-950/40 p-4 rounded-xl border border-slate-900">
              <div>
                <div className="flex justify-between items-center mb-2">
                  <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400">
                    Background Opacity
                  </label>
                  <span className="text-xs font-mono text-cyan-400 font-extrabold">
                    {Math.round(termPref.opacity * 100)}%
                  </span>
                </div>
                <input
                  type="range"
                  min={0.1}
                  max={1.0}
                  step={0.05}
                  value={termPref.opacity}
                  onChange={(e) =>
                    setTermPref({
                      ...termPref,
                      opacity: Number(e.target.value),
                    })
                  }
                  className="w-full h-2 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-cyan-500"
                />
              </div>

              <div>
                <div className="flex justify-between items-center mb-2">
                  <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400">
                    Backdrop Blur Effect
                  </label>
                  <span className="text-xs font-mono text-cyan-400 font-extrabold">
                    {termPref.blur}px
                  </span>
                </div>
                <input
                  type="range"
                  min={0}
                  max={40}
                  step={2}
                  value={termPref.blur}
                  onChange={(e) =>
                    setTermPref({ ...termPref, blur: Number(e.target.value) })
                  }
                  className="w-full h-2 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-cyan-500"
                />
              </div>
            </div>

            {/* Font Family & Cursor */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                  Font Family
                </label>
                <input
                  type="text"
                  value={termPref.font_family}
                  onChange={(e) =>
                    setTermPref({ ...termPref, font_family: e.target.value })
                  }
                  className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 font-mono focus:outline-none glass-input"
                />
              </div>
              <div>
                <label className="block text-[10px] font-extrabold uppercase tracking-wider text-slate-400 mb-1.5">
                  Cursor Style
                </label>
                <select
                  value={termPref.cursor_style}
                  onChange={(e) =>
                    setTermPref({ ...termPref, cursor_style: e.target.value })
                  }
                  className="w-full rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none glass-input"
                >
                  <option value="block">Solid Block</option>
                  <option value="underline">Underline</option>
                  <option value="bar font-bold">Vertical Bar</option>
                </select>
              </div>
            </div>

            <div className="flex items-center justify-between pt-4 border-t border-slate-900">
              {prefSavedMsg ? (
                <span className="text-xs font-bold text-emerald-400 flex items-center gap-1.5">
                  <CheckCircle2 className="w-4 h-4 text-emerald-400" />
                  Terminal preferences saved!
                </span>
              ) : (
                <span className="text-[11px] text-slate-500 font-medium">
                  Applied to all newly opened terminal windows.
                </span>
              )}
              <button
                type="submit"
                disabled={prefLoading}
                className="px-6 py-2.5 bg-linear-to-r from-cyan-600 to-indigo-600 hover:from-cyan-500 hover:to-indigo-500 text-white text-xs font-extrabold rounded-xl shadow-lg shadow-cyan-950/30 transition-all duration-200 border border-cyan-500/30"
              >
                {prefLoading ? "Saving..." : "Save Terminal Settings"}
              </button>
            </div>
          </form>
        </div>

        {/* Right Column: Interactive Live Preview Box */}
        <div className="lg:col-span-5 bg-slate-900/35 border border-slate-900 rounded-2xl p-6 shadow-xl shadow-black/40 backdrop-blur-md space-y-4">
          <div className="flex items-center justify-between border-b border-slate-900 pb-3">
            <h3 className="text-xs font-extrabold uppercase tracking-wider text-slate-300 flex items-center space-x-2">
              <Sparkles className="h-4 w-4 text-cyan-400" />
              <span>Live Aesthetic Preview</span>
            </h3>
            <span className="text-[10px] font-mono text-slate-400">
              {termPref.window_width} × {termPref.window_height}px
            </span>
          </div>

          {/* Simulated Glass Window */}
          <div
            className="w-full rounded-xl border border-slate-800 overflow-hidden shadow-2xl transition-all duration-300 relative"
            style={{
              backgroundColor: `rgba(13, 17, 23, ${termPref.opacity})`,
              backdropFilter: `blur(${termPref.blur}px)`,
              WebkitBackdropFilter: `blur(${termPref.blur}px)`,
              minHeight: "280px",
            }}
          >
            {/* Window Header */}
            <div className="px-4 py-2.5 bg-slate-950/60 border-b border-slate-900 flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <div className="w-2.5 h-2.5 rounded-full bg-rose-500/80" />
                <div className="w-2.5 h-2.5 rounded-full bg-amber-500/80" />
                <div className="w-2.5 h-2.5 rounded-full bg-emerald-500/80" />
              </div>
              <span className="text-[10px] font-mono text-slate-400 font-bold">
                root@vpsm-server:~
              </span>
            </div>

            {/* Terminal Canvas Content */}
            <div
              className="p-4 font-mono leading-relaxed"
              style={{
                fontFamily: termPref.font_family,
                fontSize: `${termPref.font_size}px`,
              }}
            >
              <p className="text-emerald-400">
                root@vpsm-node:~#{" "}
                <span className="text-slate-200 font-normal">vpsm status</span>
              </p>
              <p className="text-slate-400 text-xs mt-1.5">
                [✓] System operational
              </p>
              <p className="text-slate-400 text-xs">
                [✓] SSH key authenticated
              </p>
              <p className="text-cyan-400 text-xs mt-1.5">
                Linux 6.8.0-110-generic x86_64
              </p>
              <div className="flex items-center space-x-1 mt-3 text-slate-200">
                <span>root@vpsm-node:~# </span>
                <span
                  className={`inline-block bg-cyan-400 ${
                    termPref.cursor_style === "block"
                      ? "w-2 h-4"
                      : termPref.cursor_style === "underline"
                        ? "w-2 h-0.5 self-end"
                        : "w-0.5 h-4"
                  } ${termPref.cursor_blink ? "animate-pulse" : ""}`}
                />
              </div>
            </div>
          </div>

          <p className="text-[11px] text-slate-500 text-center font-medium">
            Changes to sliders and inputs above reflect immediately in this
            preview container.
          </p>
        </div>
      </div>
    </div>
  );
}
