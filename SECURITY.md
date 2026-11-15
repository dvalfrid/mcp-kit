# Security Policy

## Supported Versions

Only the latest release receives security fixes. A fix is released as a patch version at once, and every server using `mcp-kit` is updated and released.

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Use [GitHub private vulnerability reporting](https://github.com/dvalfrid/mcp-kit/security/advisories/new).

Include the type of issue, steps to reproduce, the affected version, and the potential impact.

| | Target |
|---|---|
| Acknowledgement | Within 7 days |
| Patch release | Within 14 days of confirmed vulnerability |

## Threat model

`mcp-kit` is a library used by local MCP servers that hold the owner's passwords, tokens and personal data. It does not reach the network. What it protects:

| Asset | Threat | Protection |
|---|---|---|
| Passwords, tokens, key passphrases | Leaked through logs, errors, tool output or JSON | Held as `secret.Secret`, which prints `[redacted]` everywhere; tested for `%v`, `%s`, `%#v`, slog and JSON |
| Passwords in `.env` | Echoed in an error about a malformed line | `.env` errors name the line number only |
| One server's secrets | Read by another server through the credential store | Every server uses its own service name; an empty name is refused. **macOS caveat:** go-keyring creates Keychain items through `/usr/bin/security`, so other programs running as the same user can read them without a prompt, no worse than a file only the owner can read |
| User data | Leaked through a panic log | A panic value is logged by type only, unless it is a runtime error |

In scope: anything that gets around one of the protections above. Out of scope: the servers' own guardrails (report them in that server's repository).
