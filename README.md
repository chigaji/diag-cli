# DiagCLI


A cross‑platform (Linux/macOS/Windows) system diagnostics CLI.


## Features

- System: CPU, memory, disk, host, uptime, temperatures (if available)
- Network: interfaces, IPs, I/O counters, quick ping, public IP (optional)
- Processes: top by CPU/MEM, filter by name, list listening ports (root/admin)
- Output: table (pretty), JSON, or YAML
- Config: YAML + ENV overrides (prefix `DIAG_`)
- Logging: zerolog with levels, colored console
- Extensible: clean interfaces for alternate backends or remote targets


## Quick Start

```bash
make build
```

## Usage

```bash
# how to use instructions
./bin/diag --help 



# Show system metrics
./bin/diag sys --all --c config.yaml

# Show network diagnostics
./bin/diag net --iface eth0

# Show top processes
./bin/diag proc --top 10 --sort cpu
```
