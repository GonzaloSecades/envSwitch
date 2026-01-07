# 🔄 envSwitch

> A blazing-fast environment switcher for legacy AngularJS apps. Replaces slow Gulp tasks with a native Go binary.

**50-100x faster than Gulp** — No Node.js required at runtime.

---

## 📋 Table of Contents

- [Why envSwitch?](#-why-envswitch)
- [Quick Start](#-quick-start)
- [Installation](#-installation)
- [Usage](#-usage)
  - [Interactive Mode (Recommended)](#interactive-mode-recommended)
  - [Command Line Mode](#command-line-mode)
- [macOS Gatekeeper Workaround](#-macos-gatekeeper-workaround)
- [Shell Aliases](#-shell-aliases)
- [Supported Formats](#-supported-formats)
- [Configuration Files](#-configuration-files)
- [Adding New Apps](#-adding-new-apps)

---

## 🚀 Why envSwitch?

| | Gulp (Node.js) | envSwitch (Go) |
|---|---|---|
| **Startup time** | 1-2 seconds | ~40ms |
| **Dependencies** | node_modules | None |
| **Memory** | ~100MB | ~10MB |
| **Distribution** | Requires Node.js | Single binary |

---

## ⚡ Quick Start

```bash
# Clone
git clone https://github.com/GonzaloSecades/envSwitch.git
cd envSwitch

# Build (auto-detects your OS/architecture)
./build-local.sh

# Run interactive mode
./envswitch -i
```

---

## 📦 Installation

### Prerequisites

- [Go 1.21+](https://go.dev/dl/) (only needed for building)

### Build for Your Platform

```bash
# macOS / Linux / Git Bash
chmod +x build-local.sh
./build-local.sh

# Windows PowerShell
.\build-local.ps1
```

### Build for All Platforms (maintainers)

```bash
./build.sh
# Creates: dist/envswitch.exe, dist/envswitch-mac-intel, dist/envswitch-mac-arm, dist/envswitch-linux
```

---

## 🎮 Usage

### Interactive Mode (Recommended)

```bash
./envswitch -i
```

Features:
- 🏗️ Select from saved apps
- ➕ Add new apps with guided setup
- 🚀 Quick switch with saved paths
- ✏️ Edit or delete app configurations
- 💾 Remembers your settings between sessions

### Command Line Mode

**Basic syntax:**
```bash
./envswitch --env <environment> [options]
```

**Options:**

| Flag | Description | Default |
|------|-------------|---------|
| `--env` | Environment name (required) | - |
| `--config-dir` | Path to config files folder | `./configs` |
| `--target` | Path to target file to modify | `./app/.../serverConfig.js` |
| `--format` | Output format: `serverConfig` or `envJs` | `serverConfig` |
| `--js` | Use `.js` config files (not `.json`) | `false` |
| `--dist` | Set `isDist` to `true` | `false` |
| `--dry-run` | Preview changes without modifying | `false` |
| `-i` | Interactive mode | `false` |

**Examples:**

```bash
# Switch to test environment (serverConfig format)
./envswitch --env test --js \
  --config-dir "/path/to/app/gulp/configs" \
  --target "/path/to/app/app/shared/services/web/serverConfig.js"

# Switch to stress environment (envJs format)
./envswitch --env stress --js --format envJs \
  --config-dir "/path/to/app/gulp/configs" \
  --target "/path/to/app/app/env.js"

# Preview changes without applying
./envswitch --env prod --js --dry-run \
  --config-dir "/path/to/configs" \
  --target "/path/to/serverConfig.js"
```

---

## 🍎 macOS Gatekeeper Workaround

If you see **"envswitch cannot be opened because it is from an unidentified developer"** or your organization restricts unsigned apps:

### Solution: Run via `go run`

Instead of compiling a binary, run directly from source:

```bash
cd /path/to/envSwitch

# Interactive mode
go run . -i

# Command line mode
go run . --env test --js --format envJs \
  --config-dir "/path/to/configs" \
  --target "/path/to/env.js"
```

This bypasses Gatekeeper because Go itself is a signed application.

---

## 🔧 Shell Aliases

Add these to your shell config (`~/.zshrc` or `~/.bashrc`):

### Option A: Using compiled binary (if Gatekeeper allows)

```bash
# Interactive mode
alias envswitch="/path/to/envSwitch/envswitch -i"

# Quick switch command
alias esw="/path/to/envSwitch/envswitch"
```

### Option B: Using `go run` (Gatekeeper workaround)

```bash
# Interactive mode
alias envswitch="go run /path/to/envSwitch/. -i"

# Quick switch command  
alias esw="go run /path/to/envSwitch/."
```

### Apply changes

```bash
source ~/.zshrc
```

### Now use from anywhere

```bash
# Interactive
envswitch

# Command line
esw --env test --js --format envJs \
  --config-dir "./gulp/configs" \
  --target "./app/env.js"
```

---

## 📁 Supported Formats

### `serverConfig` — Angular Factory

**Target file:** `serverConfig.js`

```javascript
angular.module('myApp').factory('serverConfig', function () {
    return {
        baseUrl: "https://api.example.com",
        questUrl: "https://quest.example.com",
        questFront: "https://front.example.com",
        recaptchaApiKey: "your-key",
        isDist: false
    }
})
```

**Replaced values:** `baseUrl`, `questUrl`, `questFront`, `recaptchaApiKey`, `isDist`

---

### `envJs` — Variable Declarations

**Target file:** `env.js`

```javascript
var urls = {"quest":"...","agents":"...","bo":"...","tpv":"...","vault":"...","front":"..."}; var recaptchaKey = "..."; var isDist = false; var walkMeUrl= "..."
```

**Replaced values:** `urls` object, `recaptchaKey`, `isDist`, `walkMeUrl`

---

## 📝 Configuration Files

Config files should be named `config.<env>.js` and placed in your config directory.

### For `serverConfig` format

```javascript
module.exports = function () {
    return {
        server: 'https://api.example.com',
        questServer: 'https://quest.example.com',
        questFront: 'https://front.example.com',
        google: {
            recaptcha: 'your-recaptcha-key'
        }
    }
}
```

### For `envJs` format

```javascript
module.exports = function () {
    return {
        server: {
            quest: 'https://quest.example.com',
            agents: 'https://agents.example.com',
            bo: 'https://bo.example.com',
            tpv: 'https://tpv.example.com',
            vault: 'https://vault.example.com',
            front: 'https://front.example.com'
        },
        google: {
            recaptcha: 'your-recaptcha-key'
        },
        walkmeUrl: 'https://walkme.example.com/script.js'
    }
}
```

---

## ➕ Adding New Apps

### Via Interactive Mode

1. Run `./envswitch -i` (or `go run . -i`)
2. Select **"➕ Add New App..."**
3. Follow the prompts:
   - App name
   - Config directory path
   - Target file path
   - Use JS configs? (yes/no)
   - Format (serverConfig / envJs)

### Via Command Line

Just run with your paths — no pre-registration needed:

```bash
./envswitch --env test --js --format envJs \
  --config-dir "/new/app/configs" \
  --target "/new/app/env.js"
```

---

## 🗂️ Project Structure

```
envSwitch/
├── main.go           # CLI entry point & flags
├── cli.go            # Interactive TUI (Bubble Tea)
├── jsconfig.go       # JS config file parser
├── go.mod
│
├── build-local.sh    # Build for current platform
├── build.sh          # Build all platforms
│
└── configs/          # Sample config files
    ├── config.test.json
    └── config.stress.json
```

---

## 📄 License

MIT

---

<p align="center">
  Made with 💙💛
</p>
