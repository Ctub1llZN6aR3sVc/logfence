# logfence

Lightweight structured log filtering and routing daemon for containerized environments.

---

## Installation

```bash
go install github.com/yourorg/logfence@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/logfence.git && cd logfence && go build -o logfence .
```

---

## Usage

Run the daemon with a config file:

```bash
logfence --config logfence.yaml
```

Example `logfence.yaml`:

```yaml
inputs:
  - type: stdin

filters:
  - field: level
    match: error|warn

outputs:
  - type: stdout
  - type: file
    path: /var/log/app/errors.log
```

Pipe container logs directly into logfence:

```bash
docker logs -f my-container | logfence --config logfence.yaml
```

---

## Features

- Structured log parsing (JSON, logfmt)
- Field-based filtering with regex support
- Multiple output targets (stdout, file, HTTP endpoint)
- Minimal resource footprint — designed for sidecar deployments

---

## License

MIT © yourorg