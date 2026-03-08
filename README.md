# 🔥 HotReload Engine

A robust, cross-platform CLI tool built in Go that automatically watches, builds, and restarts your backend server on file changes.

## 🚀 Features

* **Sub-Second Debouncing:** Intelligently groups rapid `Save All` events to prevent unnecessary rebuilds.
* **Zombie Process Annihilation:** Uses OS-specific syscalls (Process Groups on Unix, `taskkill` on Windows) to guarantee child processes are terminated and ports are freed.
* **Dynamic Directory Watching:** Automatically detects and watches new folders created while the tool is running.
* **Crash Resilience:** Gracefully catches syntax errors and failed builds without crashing the watcher loop.
* **Standard Library Only:** Zero external dependencies for logging or process management, utilizing idiomatic Go's `log/slog`.

## 🛠️ Quick Start

To see the hot reload in action, run the demo command. This will spin up the tool and a dummy background worker server.

```bash
# For Unix / Mac
make demo

# For Windows PowerShell
.\demo.ps1
