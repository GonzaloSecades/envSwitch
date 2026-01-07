# envswitch

A fast Go-based environment switcher for AngularJS config files with an **8-bit styled interactive CLI** 💙💛

## Why Go instead of Gulp?

| Aspect | Gulp (Node.js) | Go |
|--------|---------------|-----|
| Startup time | ~500ms-2s (module loading) | ~5ms |
| Memory | ~50-100MB | ~5-10MB |
| Dependencies | Many npm packages | Zero |
| Distribution | Requires Node.js | Single binary |

## Installation

```bash
cd envswitch
go build -o envswitch.exe .
```

## Usage

### Interactive Mode (Recommended) 🎮

Run the interactive CLI with Boca Juniors themed colors:

```bash
go run . -i
```

This launches an 8-bit styled interface that guides you through:
1. **App Selection** - Choose which app to configure (e.g., "The Vault")
2. **Config Directory** - Paste the absolute path to your config files
3. **Target File** - Paste the absolute path to the target file to modify
4. **Environment Selection** - Choose between test, stress, etc.
5. **Confirmation** - Review and execute the switch

### Command Line Mode

```bash
# Basic usage (like gulp switch --env test)
envswitch --env test

# With custom paths
envswitch --env test --config-dir ./my-configs --target ./src/serverConfig.js

# For distribution builds
envswitch --env prod --dist

# Use JavaScript config files instead of JSON
envswitch --env test --js --config-dir ./path/to/configs --target ./path/to/target.js

# Dry run - see changes without applying them
envswitch --env test --dry-run
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-i` | `false` | Run in interactive mode with visual CLI |
| `--env` | (required in CLI mode) | Environment name: test, stress, cfg, prod, etc. |
| `--config-dir` | `./configs` | Directory containing `config.{env}.json` files |
| `--target` | `./app/shared/services/web/serverConfig.js` | Target file to modify |
| `--dist` | `false` | Set `isDist` to true |
| `--js` | `false` | Use `.js` config files instead of `.json` |
| `--dry-run` | `false` | Show what would be changed without modifying |

## Config Files

Create JSON config files in your config directory:

```
configs/
├── config.test.json
├── config.stress.json
├── config.cfg.json
└── config.prod.json
```

Each config file should have this structure:

```json
{
  "server": "https://api.example.com",
  "questServer": "https://quest.example.com",
  "questFront": "https://front.example.com",
  "firebase": {
    "apiKey": "...",
    "authDomain": "...",
    "databaseURL": "...",
    "storageBucket": "...",
    "messagingSenderId": "..."
  },
  "google": {
    "mapsKey": "...",
    "analytics": "...",
    "recaptcha": "..."
  }
}
```

## Migrating from Gulp

Your old Gulp command:
```bash
gulp switch --env test
```

Becomes:
```bash
envswitch --env test
```

## Performance

Benchmarks show this Go implementation is **50-100x faster** than the Gulp equivalent for single-file operations due to:

1. No module loading overhead
2. Compiled native code
3. Direct file I/O (no streams)
4. Single-pass regex replacement

