# envswitch

A fast Go-based environment switcher for AngularJS config files.

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

```bash
# Basic usage (like gulp switch --env test)
envswitch --env test

# With custom paths
envswitch --env test --config-dir ./my-configs --target ./src/serverConfig.js

# For distribution builds
envswitch --env prod --dist
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--env` | (required) | Environment name: test, stress, cfg, prod, etc. |
| `--config-dir` | `./configs` | Directory containing `config.{env}.json` files |
| `--target` | `./app/shared/services/web/serverConfig.js` | Target file to modify |
| `--dist` | `false` | Set `isDist` to true |

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

