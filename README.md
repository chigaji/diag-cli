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
./bin/diag sys --all
./bin/diag net --iface eth0
./bin/diag proc --top 10 --sort cpu