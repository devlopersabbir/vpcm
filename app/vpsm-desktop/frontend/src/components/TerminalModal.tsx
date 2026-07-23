import React, { useEffect, useRef, useState } from 'react';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import '@xterm/xterm/css/xterm.css';
import {
  X,
  Maximize2,
  Minimize2,
  RefreshCw,
  Terminal as TerminalIcon,
  ShieldCheck,
  AlertCircle,
  AppWindow
} from 'lucide-react';

// Wails runtime bindings declarations
declare global {
  interface Window {
    runtime: {
      EventsOn: (eventName: string, callback: (data: any) => void) => () => void;
      EventsOff: (eventName: string, ...additional: string[]) => void;
    };
    go: {
      main: {
        App: {
          StartSSHTerminal: (params: any) => Promise<string>;
          SendSSHTerminalInput: (sessionId: string, data: string) => Promise<void>;
          ResizeSSHTerminal: (sessionId: string, rows: number, cols: number) => Promise<void>;
          CloseSSHTerminal: (sessionId: string) => Promise<void>;
          OpenStandaloneTerminalWindow?: (serverID: number, params: any) => Promise<void>;
          GetTerminalInitialParams?: () => Promise<any>;
        };
      };
    };
  }
}

interface ServerInfo {
  id: number;
  name: string;
  host: string;
  port: number;
  username: string;
  auth_type?: string;
  auth_secret?: string;
}

interface TerminalModalProps {
  server: ServerInfo;
  onClose: () => void;
  isStandaloneWindow?: boolean;
}

export const TerminalModal: React.FC<TerminalModalProps> = ({ server, onClose, isStandaloneWindow = false }) => {
  const terminalRef = useRef<HTMLDivElement>(null);
  const xtermRef = useRef<Terminal | null>(null);
  const fitAddonRef = useRef<FitAddon | null>(null);
  const sessionIdRef = useRef<string | null>(null);

  const [status, setStatus] = useState<'connecting' | 'connected' | 'error' | 'disconnected'>('connecting');
  const [errorMessage, setErrorMessage] = useState<string>('');
  const [isMaximized, setIsMaximized] = useState<boolean>(isStandaloneWindow);

  // Custom password prompt state if credentials missing
  const [authSecret, setAuthSecret] = useState<string>(server.auth_secret || '');
  const [requiresAuthInput, setRequiresAuthInput] = useState<boolean>(!server.auth_secret);

  const syncFit = () => {
    if (fitAddonRef.current && xtermRef.current && sessionIdRef.current) {
      fitAddonRef.current.fit();
      const rows = xtermRef.current.rows;
      const cols = xtermRef.current.cols;
      window.go?.main?.App?.ResizeSSHTerminal(sessionIdRef.current, rows, cols).catch(() => {});
    }
  };

  const startSession = async (secretToUse: string) => {
    setStatus('connecting');
    setErrorMessage('');

    try {
      const app = window.go?.main?.App;
      if (!app) {
        throw new Error('Wails runtime backend is not available.');
      }

      // Initial rows and cols
      const cols = Math.floor((terminalRef.current?.clientWidth || 800) / 9) || 80;
      const rows = Math.floor((terminalRef.current?.clientHeight || 450) / 18) || 24;

      const sid = await app.StartSSHTerminal({
        host: server.host,
        port: server.port || 22,
        username: server.username || 'root',
        auth_type: server.auth_type || 'password',
        auth_secret: secretToUse,
        rows: rows,
        cols: cols,
      });

      sessionIdRef.current = sid;
      setStatus('connected');

      // Listen for stdout/stderr data from Go backend
      window.runtime?.EventsOn(`ssh:data:${sid}`, (data: string) => {
        if (xtermRef.current) {
          xtermRef.current.write(data);
        }
      });

      // Listen for error event
      window.runtime?.EventsOn(`ssh:error:${sid}`, (errData: string) => {
        setErrorMessage(errData);
        setStatus('error');
      });

      // Listen for close event
      window.runtime?.EventsOn(`ssh:closed:${sid}`, () => {
        setStatus('disconnected');
        if (xtermRef.current) {
          xtermRef.current.write('\r\n\x1b[33m[SSH Connection Closed]\x1b[0m\r\n');
        }
      });

      // Fit layout after connect
      setTimeout(() => {
        syncFit();
      }, 150);

    } catch (err: any) {
      console.error('SSH Connection failed:', err);
      setStatus('error');
      setErrorMessage(err.message || 'Failed to establish SSH connection');
      setRequiresAuthInput(true);
    }
  };

  useEffect(() => {
    if (!terminalRef.current) return;

    // Initialize xterm.js
    const term = new Terminal({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'Menlo, Monaco, "Courier New", monospace',
      theme: {
        background: '#0d1117',
        foreground: '#c9d1d9',
        cursor: '#58a6ff',
        selectionBackground: '#264f78',
        black: '#484f58',
        red: '#ff7b72',
        green: '#3fb950',
        yellow: '#d29922',
        blue: '#58a6ff',
        magenta: '#bc8cff',
        cyan: '#39c5cf',
        white: '#b1bac4',
        brightBlack: '#6e7681',
        brightRed: '#ffa198',
        brightGreen: '#56d364',
        brightYellow: '#e3b341',
        brightBlue: '#79c0ff',
        brightMagenta: '#d2a8ff',
        brightCyan: '#56d4dd',
        brightWhite: '#f0f6fc',
      },
      convertEol: true,
    });

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);

    term.open(terminalRef.current);
    fitAddon.fit();

    xtermRef.current = term;
    fitAddonRef.current = fitAddon;

    // Handle user keystrokes in terminal
    term.onData((data) => {
      if (sessionIdRef.current && window.go?.main?.App) {
        window.go.main.App.SendSSHTerminalInput(sessionIdRef.current, data).catch((err) => {
          console.error('Failed to send terminal data:', err);
        });
      }
    });

    // Handle window resize
    const handleResize = () => {
      syncFit();
    };

    window.addEventListener('resize', handleResize);

    // Listen for menu zoom events from Wails native menu
    const zoomInOff = window.runtime?.EventsOn('terminal:zoom_in', () => {
      if (xtermRef.current) {
        const curSize = xtermRef.current.options.fontSize || 14;
        if (curSize < 32) {
          xtermRef.current.options.fontSize = curSize + 2;
          syncFit();
        }
      }
    });

    const zoomOutOff = window.runtime?.EventsOn('terminal:zoom_out', () => {
      if (xtermRef.current) {
        const curSize = xtermRef.current.options.fontSize || 14;
        if (curSize > 8) {
          xtermRef.current.options.fontSize = curSize - 2;
          syncFit();
        }
      }
    });

    const zoomResetOff = window.runtime?.EventsOn('terminal:zoom_reset', () => {
      if (xtermRef.current) {
        xtermRef.current.options.fontSize = 14;
        syncFit();
      }
    });

    // Auto start session immediately
    startSession(server.auth_secret || '');

    return () => {
      window.removeEventListener('resize', handleResize);
      zoomInOff?.();
      zoomOutOff?.();
      zoomResetOff?.();
      if (sessionIdRef.current && window.go?.main?.App) {
        window.go.main.App.CloseSSHTerminal(sessionIdRef.current).catch(() => {});
      }
      term.dispose();
    };
  }, []);

  const handleDisconnect = () => {
    if (sessionIdRef.current && window.go?.main?.App) {
      window.go.main.App.CloseSSHTerminal(sessionIdRef.current).catch(() => {});
    }
    onClose();
  };

  const handleReconnect = () => {
    if (xtermRef.current) {
      xtermRef.current.clear();
    }
    startSession(authSecret);
  };

  const toggleMaximize = () => {
    const nextState = !isMaximized;
    setIsMaximized(nextState);

    // Refit after transition
    syncFit();
    setTimeout(syncFit, 100);
    setTimeout(syncFit, 200);
    setTimeout(syncFit, 350);
  };

  const handleOpenStandaloneWindow = () => {
    if (window.go?.main?.App?.OpenStandaloneTerminalWindow) {
      window.go.main.App.OpenStandaloneTerminalWindow(
        server.id || 0,
        {
          Host: server.host,
          Port: server.port || 22,
          Username: server.username || 'root',
          AuthType: server.auth_type || 'password',
          AuthSecret: server.auth_secret || '',
        }
      ).catch((err) => {
        console.error('Failed to open standalone terminal window:', err);
      });
    }
  };

  if (isStandaloneWindow) {
    return (
      <div className="w-screen h-screen bg-[#0d1117] p-0 m-0 overflow-hidden relative">
        {status === 'error' && errorMessage && (
          <div className="absolute top-2 left-1/2 -translate-x-1/2 z-20 px-4 py-1.5 bg-red-500/20 border border-red-500/40 text-red-300 rounded-md text-xs font-mono flex items-center gap-2 shadow-lg">
            <AlertCircle className="w-4 h-4 text-red-400 shrink-0" />
            <span>{errorMessage}</span>
          </div>
        )}
        <div ref={terminalRef} className="w-full h-full text-left" />
      </div>
    );
  }

  return (
    <div
      className={`fixed inset-0 z-50 flex items-center justify-center transition-all duration-300 ${
        isStandaloneWindow
          ? 'p-0 bg-[#0d1117]'
          : isMaximized
          ? 'p-2 bg-black/90 backdrop-blur-md'
          : 'p-4 bg-black/80 backdrop-blur-sm'
      }`}
    >
      <div
        className={`flex flex-col bg-[#0d1117] transition-all duration-300 ${
          isStandaloneWindow || isMaximized
            ? 'w-full h-full rounded-none border-none'
            : 'w-[95%] max-w-5xl h-[80vh] rounded-xl border border-gray-800 shadow-2xl overflow-hidden'
        }`}
      >
        {/* Header Bar */}
        <div className="flex items-center justify-between px-4 py-3 bg-[#161b22] border-b border-gray-800 select-none">
          <div className="flex items-center gap-3">
            <div className="p-1.5 bg-blue-500/10 text-blue-400 rounded-lg border border-blue-500/20">
              <TerminalIcon className="w-4 h-4" />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <h3 className="font-semibold text-sm text-gray-200">{server.name}</h3>
                <span className="text-xs text-gray-400 font-mono">
                  {server.username}@{server.host}:{server.port || 22}
                </span>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-2">
            {/* Status Badge */}
            {status === 'connecting' && (
              <span className="flex items-center gap-1.5 text-xs px-2.5 py-1 bg-yellow-500/10 text-yellow-400 border border-yellow-500/20 rounded-full">
                <RefreshCw className="w-3 h-3 animate-spin" /> Connecting...
              </span>
            )}
            {status === 'connected' && (
              <span className="flex items-center gap-1.5 text-xs px-2.5 py-1 bg-green-500/10 text-green-400 border border-green-500/20 rounded-full">
                <ShieldCheck className="w-3 h-3" /> SSH Active
              </span>
            )}
            {status === 'error' && (
              <span className="flex items-center gap-1.5 text-xs px-2.5 py-1 bg-red-500/10 text-red-400 border border-red-500/20 rounded-full">
                <AlertCircle className="w-3 h-3" /> Connection Failed
              </span>
            )}
            {status === 'disconnected' && (
              <span className="flex items-center gap-1.5 text-xs px-2.5 py-1 bg-gray-500/10 text-gray-400 border border-gray-500/20 rounded-full">
                Disconnected
              </span>
            )}

            {/* Reconnect button */}
            {(status === 'error' || status === 'disconnected') && !requiresAuthInput && (
              <button
                onClick={handleReconnect}
                className="p-1.5 text-gray-400 hover:text-gray-200 hover:bg-gray-800 rounded-lg transition"
                title="Reconnect SSH"
              >
                <RefreshCw className="w-4 h-4" />
              </button>
            )}

            {/* Standalone Window Pop-Out Button */}
            {!isStandaloneWindow && (
              <button
                onClick={handleOpenStandaloneWindow}
                className="flex items-center gap-1.5 px-2.5 py-1 bg-cyan-500/10 hover:bg-cyan-500 text-cyan-400 hover:text-slate-950 border border-cyan-500/20 rounded-lg text-xs font-black transition-all shadow-sm"
                title="Open as Standalone Desktop Window (Alacritty Style)"
              >
                <AppWindow className="w-3.5 h-3.5" />
                <span>Window</span>
              </button>
            )}

            {/* Maximize / Restore Toggle */}
            {!isStandaloneWindow && (
              <button
                onClick={toggleMaximize}
                className="p-1.5 text-gray-400 hover:text-gray-200 hover:bg-gray-800 rounded-lg transition"
                title={isMaximized ? 'Restore Modal' : 'Maximize Window'}
              >
                {isMaximized ? <Minimize2 className="w-4 h-4" /> : <Maximize2 className="w-4 h-4" />}
              </button>
            )}

            {/* Close Terminal Modal */}
            {!isStandaloneWindow && (
              <button
                onClick={handleDisconnect}
                className="p-1.5 text-gray-400 hover:text-red-400 hover:bg-red-500/10 rounded-lg transition"
                title="Close Terminal Modal"
              >
                <X className="w-4 h-4" />
              </button>
            )}
          </div>
        </div>

        {/* Content Body */}
        <div className="relative flex-1 bg-[#0d1117] p-2 overflow-hidden">
          {status === 'error' && errorMessage && (
            <div className="absolute top-4 left-1/2 -translate-x-1/2 z-10 px-4 py-2 bg-red-500/10 border border-red-500/30 text-red-400 rounded-lg text-xs flex items-center gap-2 shadow-lg">
              <AlertCircle className="w-4 h-4 shrink-0" />
              <span>{errorMessage}</span>
            </div>
          )}

          {/* Terminal Canvas Container */}
          <div ref={terminalRef} className="w-full h-full text-left" />
        </div>
      </div>
    </div>
  );
};
