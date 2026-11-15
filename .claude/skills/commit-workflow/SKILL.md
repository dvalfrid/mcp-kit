---
name: commit-workflow
description: Full sequence for shipping a change in mcp-kit — GitHub issue, branch, implement with tests, verify, Conventional Commit, pull request. Run when starting work on a feature/bug or when ready to commit.
---

This is the sequence for every change. Do not skip steps or reorder them.

## 1. Open a GitHub issue before starting work

```powershell
gh issue create --title "..." --body "..." --label enhancement --milestone 0.2   # or --label bug
```

- **Milestone:** the release it will ship in (`gh api repos/dvalfrid/mcp-kit/milestones`). Close a milestone when its release is out.
- **Moving code up from a server:** name both servers that need it and link their issues ([docs/dependency-rules.md §3](../../../docs/dependency-rules.md)). One server is not enough.
- Never real personal data in an issue.

## 2. Branch from `main`

```powershell
git switch main; git pull; git switch -c <type>/<short-name>
```

## 3. Implement with tests

- The package stays general: no server name, no service words, no network (depguard enforces the last).
- A guardrail gets a test that fails if it breaks. A bug fix gets a test that fails without the fix.
- A breaking API change: `feat!:` and say in the body what servers must change.
- Touching a guardrail: stop and ask the owner first.

## 4. Verify

```powershell
./tools/task.ps1 fmt
./tools/task.ps1 verify     # zero lint issues, all tests green
```

Report the result honestly: **verifierat** (what was observed) or **hypotes** (what would confirm it).

## 5. Commit

```
<type>(<scope>): <subject>

<optional body>

Closes #N
```

Scope is the package (`secret`, `keystore`, `dotenv`, `logging`) or `ci`, `deps`. The `commit-msg` hook rejects a subject that doesn't match.

## 6. Push and open a pull request

```powershell
git push -u origin HEAD
gh pr create --fill
gh pr merge --auto --squash  # the owner's standing approval: auto-merge every PR
```

Never merge around the ruleset (`--admin`). If a check fails, fix it on the branch.

## 7. After a release

A new `mcp-kit` version reaches the servers through Dependabot. For a breaking change or a security fix, open an issue in each server that uses the changed package.
