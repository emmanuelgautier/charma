# charma

[![CI](https://github.com/emmanuelgautier/charma/actions/workflows/ci.yml/badge.svg)](https://github.com/emmanuelgautier/charma/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/emmanuelgautier/charma.svg)](https://pkg.go.dev/github.com/emmanuelgautier/charma)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**charma** converts text to styled ASCII art. It works as both a standalone CLI tool and an importable Go library.

## Features

- 10 FIGlet fonts (standard, doom, slant, banner, and more)
- Unicode box borders (single, double, rounded, bold, ASCII)
- ANSI colors and two-color gradients
- Four output formats: **terminal**, **txt**, **png**, **svg**
- Auto-detects TTY; strips ANSI codes in pipes/CI

## Installation

### Homebrew (macOS/Linux)

```bash
brew install emmanuelgautier/tap/charma
```

### Snap (Linux)

```bash
sudo snap install charma
```

### Chocolatey (Windows)

```powershell
choco install charma
```

### Docker

```bash
# Run directly
docker run --rm ghcr.io/emmanuelgautier/charma "Hello World"

# With flags
docker run --rm ghcr.io/emmanuelgautier/charma "Hello World" --font doom --border double --color cyan
```

### Go install

```bash
go install github.com/emmanuelgautier/charma/cmd/charma@latest
```

### Download a release binary

Pre-built binaries for Linux, macOS, and Windows are available on the [releases page](https://github.com/emmanuelgautier/charma/releases). Release artifacts are signed with [cosign](https://github.com/sigstore/cosign) — verify with:

```bash
cosign verify-blob \
  --certificate charma_<version>_checksums.txt.pem \
  --signature charma_<version>_checksums.txt.sig \
  --certificate-identity-regexp "https://github.com/emmanuelgautier/charma" \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  charma_<version>_checksums.txt
```

## CLI Usage

```bash
charma "Hello World"
charma "Hello World" --font doom --border double --color cyan
charma "v2.0"       --font slant --border rounded --output png --out-file banner.png
charma "API"        --font big   --gradient "blue:cyan" --border bold --align center
charma --list-fonts
```

### All flags

| Flag | Default | Description |
|------|---------|-------------|
| `--font` | `standard` | FIGlet font: standard, big, doom, isometric1, slant, block, 3d, shadow, banner, bulbhead |
| `--border` | `none` | Border style: none, single, double, rounded, bold, ascii |
| `--color` | `default` | Text color: red, green, blue, cyan, magenta, yellow, white |
| `--border-color` | `default` | Border color (same choices as `--color`) |
| `--align` | `left` | Text alignment: left, center, right |
| `--padding` | `1` | Inner padding inside border box |
| `--output` | `terminal` | Output format: terminal, txt, png, svg |
| `--out-file` | `./output.<ext>` | Output file path |
| `--width` | terminal width | Max width in characters |
| `--gradient` | | Two-color gradient e.g. `"red:blue"` |
| `--bg-color` | `black` | Background color for PNG/SVG output |
| `--no-color` | `false` | Strip all ANSI codes |
| `--list-fonts` | | Print available fonts and exit |
| `--version` | | Print version, commit, and build date |

## Library Usage

```go
import "github.com/emmanuelgautier/charma"

opts := charma.DefaultOptions()
opts.Font   = "doom"
opts.Border = "double"
opts.Color  = "cyan"

result, err := charma.Generate("Hello", opts)
if err != nil {
    log.Fatal(err)
}

fmt.Println(result.Styled)          // ANSI-styled for terminal
fmt.Println(result.Lines[0])        // plain-text first line
```

## Development

```bash
git clone https://github.com/emmanuelgautier/charma
cd charma

make build        # compile binary → bin/charma
make test         # run all tests
make test-race    # run with race detector
make coverage     # generate coverage.html
make lint         # golangci-lint
make test-update  # regenerate golden test files
make snapshot     # local GoReleaser snapshot
```

## Project structure

```
charma/
├── charma.go              # Public library API
├── cmd/charma/main.go     # CLI entry point (Cobra)
├── internal/
│   ├── renderer/            # FIGlet rendering via go-figure
│   ├── border/              # Box drawing via box-cli-maker
│   ├── color/               # ANSI color + gradient via lipgloss
│   └── output/              # terminal / txt / png / svg writers
└── testdata/                # Golden files
```

## License

MIT — see [LICENSE](LICENSE).
