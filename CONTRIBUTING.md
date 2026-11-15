# Contributing

This is a small personal project, run like a real one: issues, branches, pull requests, CI and releases. The rules are in [STANDARDS.md](STANDARDS.md); the rules between `mcp-kit` and the servers are in [docs/dependency-rules.md](docs/dependency-rules.md).

## Setting up a development machine

| Tool | Why | Windows |
|---|---|---|
| **Git** | | `winget install Git.Git` |
| **Go**, the `toolchain` version in [`go.mod`](go.mod) | Build and test | `winget install GoLang.Go` |
| **PowerShell 7** (`pwsh`) | The task runner `tools/task.ps1` | `winget install Microsoft.PowerShell` |
| **GitHub CLI** (`gh`) | Issues and pull requests | `winget install GitHub.cli`, then `gh auth login` |

`golangci-lint`, `govulncheck` and `lefthook` are pinned in [`tools/go.mod`](tools/go.mod) and run through `go tool`.

```powershell
git clone git@github.com:dvalfrid/mcp-kit.git
cd mcp-kit
./tools/task.ps1 setup     # git hooks: signing, fmt + lint before commit, commit-message check
./tools/task.ps1 verify    # everything should pass
```

Commits must be signed (SSH signing), as in [proton-mail-mcp's CONTRIBUTING](https://github.com/dvalfrid/proton-mail-mcp/blob/main/CONTRIBUTING.md#3-sign-your-commits). The pre-commit hook refuses until signing is configured.

### Working on mcp-kit and a server at the same time

Put a `go.work` in the server's folder (it is git-ignored there):

```powershell
cd ../proton-calendar-mcp
go work init . ../mcp-kit
```

The server then builds against your local `mcp-kit`. Release `mcp-kit` and update the server's `go.mod` before the server's pull request is merged; never commit a `replace` directive.

## Workflow

1. **Open an issue**, on the milestone of the release it will ship in. Moving code up from a server: name both servers ([dependency rules §3](docs/dependency-rules.md#3-when-code-moves-up-into-mcp-kit)).
2. **Branch from `main`:** `git switch -c <type>/<short-name>`.
3. **Implement**, with tests. A guardrail gets a test that fails if it breaks.
4. **Run `./tools/task.ps1 verify`.** Zero lint issues is the bar.
5. **Commit** with a Conventional Commit message (`Closes #N`).
6. **Open a pull request.** The **Verify** check must be green.

Never commit to `main` directly. PRs are squash-merged.

## Commit convention

```
<type>(<scope>): <subject>

<optional body>

Closes #N
```

- **type:** `feat`, `fix`, `perf`, `docs`, `refactor`, `test`, `build`, `ci`, `chore`, `style`. Only `feat`, `fix` and `perf` appear in the changelog.
- **scope:** the package (`secret`, `keystore`, `dotenv`, `logging`), or `ci`, `deps`.
- **breaking change:** `feat!:` or a `BREAKING CHANGE:` footer, with what servers must change.

## Versions and releases

- [Semantic Versioning](https://semver.org/). Before 1.0, a breaking change or a feature bumps the minor version.
- Release Please keeps a release PR open. Merging it tags `vX.Y.Z` and creates the GitHub release. There are no binaries: the servers pick up the version through Dependabot.
- The release PR is opened with Actions' own token, which starts no workflows. Push a signed empty commit to its branch to run the required checks as a pull request.

## Continuous integration

| Workflow | When | What |
|---|---|---|
| `verify.yml` | Every push to `main` and every PR | Full `verify` on Linux with `-race` (**required**), tests on Windows and macOS |
| `codeql.yml` | Push, PR, weekly | CodeQL for Go and Actions |
| `audit.yml` | Weekly | `govulncheck` against the current `go.sum` |
| `release-please.yml` | Push to `main` | Release PR and tag |
| Dependabot | Weekly | Grouped updates for actions, Go modules and pinned tools; no majors |

All actions are pinned to a commit SHA, with the version in a comment.

## AI tools

This project is largely written with Claude. Say so with a `Co-Authored-By:` line. `CLAUDE.md` is configuration for AI coding tools.
