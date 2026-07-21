import React, { useState } from "react";
import { Terminal, Shield, Cpu, Zap, ArrowRight, Check, Loader2 } from "lucide-react";

interface OnboardingProps {
  onComplete: (userSession?: any) => void;
}

export default function Onboarding({ onComplete }: OnboardingProps) {
  const [loadingMode, setLoadingMode] = useState<"google" | "github" | "skip" | null>(null);

  const handleGoogleLogin = () => {
    setLoadingMode("google");
    setTimeout(() => {
      const user = {
        name: "Google User",
        email: "user@gmail.com",
        provider: "google",
      };
      localStorage.setItem("vpsm_user_session", JSON.stringify(user));
      localStorage.setItem("vpsm_onboarding_completed", "true");
      onComplete(user);
    }, 500);
  };

  const handleGithubLogin = () => {
    setLoadingMode("github");
    setTimeout(() => {
      const user = {
        name: "GitHub Developer",
        email: "dev@github.com",
        provider: "github",
      };
      localStorage.setItem("vpsm_user_session", JSON.stringify(user));
      localStorage.setItem("vpsm_onboarding_completed", "true");
      onComplete(user);
    }, 500);
  };

  const handleSkip = () => {
    setLoadingMode("skip");
    setTimeout(() => {
      localStorage.setItem("vpsm_onboarding_completed", "true");
      onComplete(null);
    }, 300);
  };

  return (
    <div className="fixed inset-0 bg-[#07090e] text-slate-100 z-50 flex items-center justify-center p-6 select-none overflow-hidden font-sans">
      {/* Dynamic Background Gradients */}
      <div className="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[350px] bg-gradient-to-tr from-cyan-500/10 via-indigo-500/10 to-purple-500/10 rounded-full blur-[140px] pointer-events-none" />

      <div className="w-full max-w-4xl grid grid-cols-1 md:grid-cols-12 gap-8 items-center relative z-10">
        {/* Left Column: Product Showcase / Terminal Preview */}
        <div className="md:col-span-6 space-y-6">
          <div className="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-slate-900/80 border border-slate-800 text-[11px] font-medium text-cyan-400">
            <span className="flex h-2 w-2 rounded-full bg-cyan-400 animate-pulse" />
            <span>v0.2.0 Release Ready</span>
          </div>

          <div>
            <h1 className="text-3xl md:text-4xl font-extrabold tracking-tight text-white leading-tight">
              Master Your VPS Infrastructure.
            </h1>
            <p className="text-slate-400 text-sm mt-3 leading-relaxed">
              Unified real-time SSH inventory management, specifications collector, and session logging for modern engineering teams.
            </p>
          </div>

          {/* Interactive Shell Preview Box */}
          <div className="bg-slate-950/80 border border-slate-800/80 rounded-2xl p-4 font-mono text-[11px] shadow-2xl shadow-black/60 relative overflow-hidden group">
            <div className="flex items-center space-x-1.5 pb-3 border-b border-slate-900 mb-3">
              <div className="w-2.5 h-2.5 rounded-full bg-rose-500/80" />
              <div className="w-2.5 h-2.5 rounded-full bg-amber-500/80" />
              <div className="w-2.5 h-2.5 rounded-full bg-emerald-500/80" />
              <span className="text-[10px] text-slate-500 ml-2 font-sans">vpsm-daemon ~ inventory audit</span>
            </div>
            <div className="space-y-1.5 text-slate-350">
              <div className="flex items-center space-x-2 text-cyan-400 font-bold">
                <span>$</span>
                <span>vpsm audit --host app.vpsm.io</span>
              </div>
              <div className="text-slate-400">[info] Connecting over SSH (rsa-2048)...</div>
              <div className="text-emerald-400">[success] Discovered Ubuntu 24.04 LTS (x86_64)</div>
              <div className="text-indigo-300">[specs] 8 Cores • 31.0GB RAM • Docker & Nginx</div>
              <div className="text-slate-400">[sync] Inventory synced with Cloud REST API</div>
            </div>
          </div>

          {/* Core Highlights List */}
          <div className="space-y-2 pt-1 text-xs text-slate-400">
            <div className="flex items-center space-x-2.5">
              <Check className="h-4 w-4 text-cyan-400 shrink-0" />
              <span>Decoupled REST API Integration (Zero local direct queries)</span>
            </div>
            <div className="flex items-center space-x-2.5">
              <Check className="h-4 w-4 text-cyan-400 shrink-0" />
              <span>Background SSH specs collector & software application auditor</span>
            </div>
          </div>
        </div>

        {/* Right Column: Authentication & Sign-in Box */}
        <div className="md:col-span-6">
          <div className="bg-slate-900/60 border border-slate-800 rounded-3xl p-8 shadow-2xl backdrop-blur-md relative overflow-hidden">
            {/* Subtle Top Accent */}
            <div className="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-cyan-500 via-indigo-500 to-purple-500" />

            <div className="mb-6">
              <h2 className="text-xl font-bold text-white tracking-tight">Get Started</h2>
              <p className="text-slate-400 text-xs mt-1">
                Choose your preferred login option or continue directly as a guest.
              </p>
            </div>

            <div className="space-y-3.5">
              {/* Google Button */}
              <button
                onClick={handleGoogleLogin}
                disabled={loadingMode !== null}
                className="w-full py-3 px-4 bg-slate-950/80 hover:bg-slate-950 text-slate-200 hover:text-white rounded-xl border border-slate-800 hover:border-slate-700 flex items-center justify-center space-x-3 transition-all duration-200 text-xs font-semibold shadow-sm group disabled:opacity-50"
              >
                {loadingMode === "google" ? (
                  <Loader2 className="h-4 w-4 animate-spin text-cyan-400" />
                ) : (
                  <svg className="h-4 w-4" viewBox="0 0 24 24">
                    <path
                      fill="#4285F4"
                      d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                    />
                    <path
                      fill="#34A853"
                      d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                    />
                    <path
                      fill="#FBBC05"
                      d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.06H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.94l2.85-2.22.81-.63z"
                    />
                    <path
                      fill="#EA4335"
                      d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.06l3.66 2.84c.87-2.6 3.3-4.52 6.16-4.52z"
                    />
                  </svg>
                )}
                <span>{loadingMode === "google" ? "Connecting Google..." : "Continue with Google"}</span>
              </button>

              {/* GitHub Button */}
              <button
                onClick={handleGithubLogin}
                disabled={loadingMode !== null}
                className="w-full py-3 px-4 bg-slate-950/80 hover:bg-slate-950 text-slate-200 hover:text-white rounded-xl border border-slate-800 hover:border-slate-700 flex items-center justify-center space-x-3 transition-all duration-200 text-xs font-semibold shadow-sm group disabled:opacity-50"
              >
                {loadingMode === "github" ? (
                  <Loader2 className="h-4 w-4 animate-spin text-cyan-400" />
                ) : (
                  <svg className="h-4 w-4 fill-current text-white" viewBox="0 0 24 24">
                    <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z" />
                  </svg>
                )}
                <span>{loadingMode === "github" ? "Connecting GitHub..." : "Continue with GitHub"}</span>
              </button>
            </div>

            {/* Divider */}
            <div className="relative my-6 flex items-center justify-center">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-slate-800/80" />
              </div>
              <span className="relative bg-slate-900 px-3 text-[10px] uppercase font-bold text-slate-550 tracking-wider">
                Or Continue As Guest
              </span>
            </div>

            {/* Skip Button */}
            <button
              onClick={handleSkip}
              disabled={loadingMode !== null}
              className="w-full py-2.5 px-4 bg-transparent hover:bg-slate-800/50 text-slate-400 hover:text-slate-200 rounded-xl border border-slate-800/60 transition-all duration-200 text-xs font-semibold flex items-center justify-center space-x-1.5 disabled:opacity-50"
            >
              {loadingMode === "skip" ? (
                <Loader2 className="h-3.5 w-3.5 animate-spin text-slate-400" />
              ) : (
                <>
                  <span>Skip for now</span>
                  <ArrowRight className="h-3.5 w-3.5" />
                </>
              )}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
