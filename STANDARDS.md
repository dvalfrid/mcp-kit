# Code Standards

This document defines the coding, security and testing standards for `mcp-kit`. All contributors and AI assistants follow it. The rules between this library and the servers that use it are in [docs/dependency-rules.md](docs/dependency-rules.md).

`mcp-kit` is used by servers that hold the owner's passwords, tokens and personal data. A mistake here reaches every server, so the [guardrails](#guardrails) outrank every other rule.

## Contents

- [Tools and commands](#tools-and-commands)
- [Go](#go)
- [Guardrails](#guardrails)
- [Testing](#testing)
- [Dependencies](#dependencies)
- [AI-generated code](#ai-generated-code)

---

## Tools and commands

All tasks go through one script, so local runs and CI do the same thing. It needs PowerShell 7 (`pwsh`).

| Purpose | Command |
|---|---|
| Format (modifies files) | `./tools/task.ps1 fmt` |
| Check formatting (CI) | `./tools/task.ps1 fmt-check` |
| Lint (`golangci-lint`, includes `go vet`) | `./tools/task.ps1 lint` |
| Run all tests | `./tools/task.ps1 test` (add `-Race` where cgo is available, as in CI) |
| Known vulnerabilities (`govulncheck`) | `./tools/task.ps1 vuln` |
| Full verification (tidy, fmt-check, lint, test, vuln) | `./tools/task.ps1 verify` |
| First-time setup (git hooks) | `./tools/task.ps1 setup` |

Tool versions are pinned in [`tools/go.mod`](tools/go.mod) and run with `go tool -modfile=tools/go.mod …`. Run `./tools/task.ps1 verify` before every commit.

---

## Go

The same Go rules as the servers:

- **Formatting:** `gofmt` and `goimports` through `golangci-lint fmt`. Imports in three groups: standard library, third party, this module.
- **Naming:** [Effective Go](https://go.dev/doc/effective_go#names) and the [code review comments](https://go.dev/wiki/CodeReviewComments). Initialisms in caps (`URL`, `ID`, `MCP`). No stutter.
- **Doc comments:** every package and every exported identifier. Explain why, not what.
- **Errors:** returned, wrapped with `%w`, compared with `errors.Is`/`As`. An error never contains a secret's value or a line of user data.
- **Context:** functions that do I/O take a `context.Context` first.
- **Logging:** `log/slog` to stderr only; `forbidigo` rejects `fmt.Print*`.
- **State:** no package-level mutable state.
- **Layout:** one package per responsibility, at the module root (`secret/`, `keystore/`). No `internal/` here: the packages are the API. No `pkg/`.
- **API:** small and concrete. A function earns its place when a server calls it. An interface earns its place when there are two implementations (like `keystore.OS` and `keystore.Memory`).
- **Linter suppressions:** every `//nolint` names the linter and says why.
- **Go version:** the `go` line is the minimum; the `toolchain` line is the newest patch release. Keep both in step with the servers.

---

## Guardrails

Hard rules. A change that breaks one is not merged. Each has a test.

| Rule | Enforced by |
|---|---|
| No network access in the kit | depguard rejects `net` and `net/http` |
| A `Secret` never shows its value (`%v`, `%s`, `%#v`, slog, JSON) | `TestSecretNeverShows` |
| A `.env` error never echoes the line (it could hold a password) | `TestReadNeverEchoesALine` |
| One server never reads another's secrets: the credential store is keyed by service name, and an empty name is refused | `TestOSKeepsServicesApart`, `TestOSNeedsAService` |
| A panic is logged without its value, unless it is a runtime error | `TestLogPanicHidesTheValue` |
| stdout is never written | forbidigo |

### Public repository

The repository is public. No secrets, no real personal data, no local paths with user names. Test data is invented (`example.com`, `example.org`, `example.net`).

---

## Testing

- Tests sit next to the code, table-driven where it fits.
- No test needs a real credential store or the network: `keystore` tests use go-keyring's in-memory mock.
- Every guardrail has a test that fails if the guardrail breaks.
- A bug fix comes with a test that fails without the fix.
- The suite runs with `-race` in CI, and on Windows and macOS too.

---

## Dependencies

Standard library first. Every server inherits what is listed here, so the list stays short. Approved:

| Module | For |
|---|---|
| `github.com/zalando/go-keyring` (with `danieljoos/wincred`, `godbus/dbus/v5`, `golang.org/x/sys`) | The OS credential store on Windows, macOS and Linux, without cgo. Already approved for proton-mail-mcp (2026-11-11); its macOS weakness is in SECURITY.md |

Development tools are pinned in `tools/go.mod`: golangci-lint, govulncheck, lefthook.

Anything else needs explicit approval from the owner, with a reason in the PR. `go mod tidy` must leave `go.mod` and `go.sum` unchanged (checked in CI).

---

## AI-generated code

AI-written code meets the same standards and review as any other. Keep it simple and concrete; don't add dependencies, restructure packages or change the API without being asked. Never weaken a guardrail or its test to make something pass: stop and ask.
