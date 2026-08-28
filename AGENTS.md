# AGENTS.md

## Project Overview

`gbrowse` is a CLI tool that opens a Git repository's file or directory in the browser.
Given a path (or `FILE:LINUM`), it resolves the remote origin URL, the current commit hash,
and the relative path inside the repository, then constructs a GitHub blob URL and opens it.

```
gbrowse [flags] [target]

  target = PATH | FILE:LINUM
```

There are two binaries in this repository:

| Binary | Path | Purpose |
|--------|------|---------|
| `gbrowse` | `cmd/gbrowse` | Main CLI — opens/prints the GitHub URL |
| `gbrowse-git` | `cmd/gbrowse-git` | Mock `git` command used in integration tests |

---

## Repository Layout

```
gbrowse/
├── cmd/
│   ├── gbrowse/        # Main entry point (main.go, env.go, main_test.go)
│   └── gbrowse-git/    # Mock git binary for testing (main.go, main_test.go)
├── browse/             # OS-specific browser opener (darwin/linux/windows)
├── ctxlog/             # Context-aware structured logger (slog wrapper)
├── env/                # Helper: env.GetOr()
├── execx/              # Helper: runs an external command and returns stdout
├── git/                # git.Git interface + implementation using execx
├── parse/              # Parses target strings and remote URLs
├── urlx/               # Assembles the final GitHub blob URL
├── dist/               # Build output (gitignored)
├── Makefile
├── go.mod / go.sum
└── README.md
```

### Package Responsibilities

| Package | Role |
|---------|------|
| `browse` | Opens a URL in the system browser; platform-specific via build tags |
| `ctxlog` | Wraps `log/slog` with context propagation; use `ctxlog.From(ctx)` to retrieve logger |
| `env` | `env.GetOr(key, default)` — thin wrapper around `os.Getenv` |
| `execx` | `execx.Run(ctx, cmd, args...)` — runs a process and returns trimmed stdout |
| `git` | `git.Git` interface with methods like `RemoteOriginURL`, `CommitHash`, `ShowPrefix`, `RelativePath` |
| `parse` | `parse.ReadTarget(s)` parses `PATH` or `FILE:LINUM`; `parse.ReadRepoURL` normalises remote URLs |
| `urlx` | `urlx.Build(ctx, git, target)` assembles the final blob URL |

---

## Build & Run

```sh
# Build the main binary to dist/gbrowse
make

# Or build manually
go build -trimpath -race -v -o dist/gbrowse ./cmd/gbrowse
```

> The output binary is `dist/gbrowse`. The `dist/` directory is gitignored.

---

## Testing

```sh
# Run all tests (with race detector and coverage)
make test
# Equivalent: go test -cover -race ./...
```

### Integration Test Architecture

`cmd/gbrowse/main_test.go` contains an end-to-end test (`TestEndToEnd`) that:

1. Builds `gbrowse` and `gbrowse-git` binaries into a temp directory.
2. Runs `gbrowse -print [target]` with `GIT=<path-to-gbrowse-git>` to inject a mock git.
3. Configures the mock via `GBROWSE_GIT_CONFIG=<json>` (see `cmd/gbrowse-git/main.go`).

**Mock git JSON fields** (all strings):

| Field | Corresponding `git` subcommand |
|-------|-------------------------------|
| `default_branch` | `git remote show origin` → `HEAD branch:` |
| `remote_origin_url` | `git config --get remote.origin.url` |
| `head_object_name` | `git rev-parse --abbrev-ref @` |
| `show_prefix` | `git rev-parse --show-prefix` |
| `relative_path` | `git ls-files --full-name <path>` |
| `describe_tag` | `git describe --tags --abbrev=0` |
| `show_current` | `git branch --show-current` |
| `commit_hash` | `git rev-parse @` |

---

## Code Generation

Two tools are used for code generation, managed via `go tool` directives in `go.mod`:

| Tool | Directive | Output |
|------|-----------|--------|
| `goconfig` | `//go:generate go tool goconfig ...` | `git/config_generated.go` — option builder for `git.Config` |
| `dataclass` | `//go:generate go tool dataclass ...` | `parse/parameter_dataclass_generated.go` — `InternalTarget` value type |

```sh
# Regenerate all generated files
make generate

# Remove all generated files
make clean-generated
# Equivalent: find . -name "*_generated.go" -type f -delete
```

> **Never edit `*_generated.go` files by hand.** Run `make generate` after changing the
> `//go:generate` directive or its arguments.

---

## Environment Variables

### `gbrowse`

| Variable | Default | Description |
|----------|---------|-------------|
| `GIT` | `git` | Path to the git executable |
| `DEBUG` | _(unset)_ | Set to any non-empty value to enable debug (JSON) logging |

### `gbrowse-git` (mock)

| Variable | Default | Description |
|----------|---------|-------------|
| `GBROWSE_GIT_CONFIG` | `{}` | JSON object mapping git subcommands to mock output values |

---

## Vulnerability Check

```sh
make vuln
# Equivalent: go tool govulncheck ./...
```

---

## Development Guidelines

1. **Read before writing.** Explore relevant packages (`git`, `urlx`, `parse`) before modifying URL-building logic.
2. **Keep generated files clean.** Do not edit `*_generated.go` files; regenerate them with `make generate` instead.
3. **Tests use the mock git binary.** When adding new `git.Git` interface methods, also add the corresponding field to `cmd/gbrowse-git/main.go`'s `config` struct and `intoMappingTuples()`.
4. **Platform-specific code lives in `browse/`.** Use `_darwin.go` / `_linux.go` / `_windows.go` build-tag files; do not use `//go:build` guards in a single file.
5. **Logger via context.** Always propagate the logger through `ctx` using `ctxlog.With` / `ctxlog.From`; do not create package-level loggers.
6. **Run tests before committing.**

   ```sh
   make vet && make test
   ```

7. **URL format.** The generated blob URL follows:
   `<remote_origin_url>/blob/<commit_hash>/<show_prefix><relative_path>[#L<linum>]`
