# telesieve

**telesieve** is a local CLI that scans log dumps for noisy telemetry waste and sensitive data (PII). Everything runs on your machine—no uploads, no API calls.

## Install

Pre-built binaries for Linux and macOS (amd64/arm64) are on
[GitHub Releases](https://github.com/sysrqio/telesieve/releases).

```bash
go build -o telesieve ./cmd/telesieve
```

Requires Go 1.22+.

### Cut a release

Tag and push a version; CI runs tests and publishes release assets:

```bash
git tag v0.1.0 && git push origin v0.1.0
```

## Usage

### Scan a log file

```bash
telesieve scan ./logs/app.log
telesieve scan ./logs/app.log.gz --provider splunk --output json
telesieve scan ./logs/app.log --mask-pii=false
```

| Flag | Default | Description |
|------|---------|-------------|
| `--provider` | `datadog` | Cost model: `datadog`, `splunk`, `dynatrace`, `custom` |
| `--ingestion-rate` | `0.10` | USD per GB ingested (override vendor defaults) |
| `--indexing-rate` | `1.70` | USD per GB indexed |
| `--output` | `terminal` | `terminal`, `json`, or `markdown` |
| `--mask-pii` | `true` | Mask sensitive values in output; exit code **2** if critical PII is found |

**Exit codes:** `0` success, `1` error, `2` critical PII detected (when `--mask-pii` is enabled).

Supported inputs: JSON lines, logfmt, plain text, and gzip-compressed files (`.gz`).

Waste detection includes Kubernetes probe noise (`/healthz`, `/livez`, `/readyz`), Redis `PING`/`PONG`, and common debug/heartbeat chatter.

Monthly cost figures assume the dump represents one day of volume, projected over 30 days. Tune `--ingestion-rate` and `--indexing-rate` to match your vendor contract.

### Generate OpenTelemetry Collector snippets

```bash
telesieve generate ./logs/app.log --drop-health --mask-tokens
```

Writes filter and transform processor YAML to stdout for dropping health checks and masking tokens.

### Version

```bash
telesieve version
```

## Development

```bash
make test
make build
```

## License

MIT — see [LICENSE](LICENSE).
