# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## What this is

`mcp-kit`: shared Go packages for the owner's local MCP servers (proton-mail-mcp, proton-calendar-mcp, later ones). A library, no binaries. **The rules between the repositories are in [docs/dependency-rules.md](docs/dependency-rules.md) and outrank convenience:** the kit never imports a server, never reaches the network, holds no service code, and code moves up only when a second server needs it.

## Skills (invoke with /name)

| Skill | When to use |
|---|---|
| `/commit-workflow` | Opening an issue, branching, committing, opening a PR |

## Commands

```powershell
./tools/task.ps1 verify     # full check: tidy, fmt-check, lint, test, vuln — run before every commit
./tools/task.ps1 fmt        # format (modifies files)
./tools/task.ps1 lint       # golangci-lint, zero issues required
./tools/task.ps1 test       # all tests
./tools/task.ps1 setup      # install git hooks (first time)
```

Go is at `C:\Program Files\Go\bin`; if a shell doesn't find `go`, prepend it to `PATH`. Tools are pinned in `tools/go.mod` and run via `go tool "-modfile=tools/go.mod" <tool>` (quote it in PowerShell).

## Public repository

This repo is **public**. Never commit secrets, real personal data (email, events, contacts, addresses, names), or local paths with user names. Test data is invented (`alice@example.com`). Check the diff before every commit.

## Hard rules

- No network: `net` and `net/http` are rejected by depguard. Network code stays in the server.
- No service code (Proton, IMAP, CalDAV, …) and no MCP tool texts for a service.
- stdout belongs to the servers' MCP protocol: log with `log/slog` to stderr, never `fmt.Print*`.
- Secrets are `secret.Secret`; never log or return their value.
- A breaking API change is `feat!:`/`fix!:` and the release notes say what servers must change.
- Never weaken a guardrail or its test. No new dependencies without the owner's approval (STANDARDS.md#dependencies).
- Never commit to `main` directly.

## Commit message format

Mandatory — `release-please` parses commit subjects.

```
<type>(<scope>): <subject>

<optional body>

Closes #N
```

- **type:** `feat`, `fix`, `perf`, `docs`, `refactor`, `test`, `build`, `ci`, `chore`, `style`.
- **scope:** the package: `secret`, `keystore`, `dotenv`, `logging`, or `ci`, `deps`.
- **subject:** imperative, lower-case start, no trailing period.

## Evidence, not guesses

Label every cause, diagnosis and "fixed" claim **verifierat** (what was observed) or **hypotes** (what would confirm it). "Fixed" means reproduced, changed, and gone in the same scenario.

## Layout

| Path | Role |
|---|---|
| `secret/` | `Secret` type: never prints its value |
| `keystore/` | OS credential store per service name; `Memory` for tests |
| `dotenv/` | `.env` reading; `Lookup`, `Chain`, `FromMap` |
| `logging/` | stderr logger; `LogPanic` |
| `docs/dependency-rules.md` | Rules between the repositories |
