# Rules between the repositories

These rules hold for every repository in the family: `mcp-kit` and the MCP servers that use it (`proton-mail-mcp`, `proton-calendar-mcp`, and later ones). Each server's `CLAUDE.md` links here. A server's own `STANDARDS.md` may add rules, never loosen these.

## The family

```text
            mcp-kit                      shared Go packages, no service code
          ▲     ▲      ▲
          │     │      │                 a server imports mcp-kit, never the other way
proton-mail-mcp  proton-calendar-mcp  …  one repository per server
```

Each server is its own repository, with its own version, releases, installer, secrets and security rules. A user who wants one server installs only that one.

## 1. Direction

1. A server may import `mcp-kit`. `mcp-kit` never imports a server.
2. A server never imports another server. Code two servers need goes into `mcp-kit` (rule 3), or is copied.
3. Inside a server, everything stays under `internal/`, so nothing else can import it. `mcp-kit` is the only shared code.

## 2. What belongs in mcp-kit

**In:** code that is the same for every server, whatever service it talks to.

- secrets that never print (`secret`), the OS credential store (`keystore`), reading `.env` and chaining settings (`dotenv`), logging to stderr and panic logging (`logging`);
- later, when a second server needs them: the per-run change limit, marking untrusted content in tool output, the `about` tool, the setup window, saving files into one folder, fetching from the public internet.

**Out:** anything tied to one service or one server's rules.

- No service protocol or API: IMAP, SMTP, Proton's API, CalDAV, BlueBubbles, signal-cli stay in their server.
- **No network access at all.** The linter rejects `net` and `net/http` in `mcp-kit`. A package that must reach the network stays in its server, where that server's rules and review cover it.
- No MCP tool definitions or tool texts for a service. Those are the server's interface and carry its guardrails.
- No server's guardrail decisions (what may be sent, deleted, written). `mcp-kit` gives the building blocks; each server decides how strict to be.

## 3. When code moves up into mcp-kit

Code moves up when **a second server needs it**, not before. One user is not enough to know what is general.

1. Open an issue in `mcp-kit` naming both servers and what each needs.
2. Move the code **with its tests**, and make it general: no server name, no service words, no defaults that only fit one server (for example the credential store takes the service name as a field).
3. Release `mcp-kit` (a new minor version before 1.0).
4. In each server, a separate pull request switches to the `mcp-kit` package and deletes its own copy. The servers' tests must pass unchanged, which shows the behaviour didn't change.

Moving a package never weakens it. If the server's copy had a guardrail test, the `mcp-kit` package gets the same test.

## 4. Dependencies

1. `mcp-kit` keeps its dependencies few: every server inherits them. Standard library first.
2. A new dependency in `mcp-kit` needs the owner's approval, with the reason in the pull request and in `STANDARDS.md#dependencies`. Approval in one repository doesn't carry over to another.
3. Approved in `mcp-kit` today: `github.com/zalando/go-keyring` (and what it pulls in).
4. Each server lists its own approved dependencies in its own `STANDARDS.md`.

## 5. Versions

1. `mcp-kit` follows [Semantic Versioning](https://semver.org/) with tags `vX.Y.Z`, made by Release Please. Before 1.0 a breaking change bumps the minor version, and the release notes say what servers must change.
2. A server requires a **tagged** `mcp-kit` version in `go.mod`. Never a branch or a pseudo-version in a release, and never a committed `replace` directive.
3. To work on both at once, use a `go.work` file in the server's folder (git-ignored), for example `use . ../mcp-kit`. CI never sees it. While the version the server requires isn't tagged yet, the workspace also needs `go work edit "-replace=github.com/dvalfrid/mcp-kit@vX.Y.Z=../mcp-kit"`: a replace in `go.work` is local, never in `go.mod`.
4. Dependabot in each server proposes `mcp-kit` updates. A server takes them through its full `verify`, like any other dependency.
5. A security fix in `mcp-kit` gets a patch release at once, and every server is updated and released.

## 6. Shared tooling

The servers and `mcp-kit` share the same setup: `tools/task.ps1`, pinned tools in `tools/go.mod`, `.golangci.yml`, Release Please, Dependabot, signed commits, the `commit-workflow` skill, Conventional Commits. When one of them improves (a newer linter, a better workflow step), the change is copied to the others in the same week, each in its own pull request. There is no shared CI config: each repository's workflows stand alone and can be read on their own.

## 7. What stays the same in every repository

- Public repository rules: no secrets, no real personal data, no local paths with user names.
- stdout is the MCP protocol: log to stderr only.
- Secrets are held as `secret.Secret`, stored only in the OS credential store, `.env` (git-ignored) or the environment, and typed by the owner in a terminal, never through the AI client.
- Content from a service (email, events, contacts, messages) is untrusted input and is marked as such in tool output.
- Every guardrail has a test that fails if it breaks.
- Claims about cause and fixes are labelled verified or hypothesis.
