# mcp-kit

Shared Go packages for a family of local MCP servers: [proton-mail-mcp](https://github.com/dvalfrid/proton-mail-mcp), [proton-calendar-mcp](https://github.com/dvalfrid/proton-calendar-mcp) and later ones.

It holds only what every server does the same way. Service code (Proton, IMAP, calendars) stays in each server. `mcp-kit` never reaches the network.

| Package | What it does |
|---|---|
| `secret` | `Secret`: a string that never shows its value in logs, errors, `%v` or JSON |
| `keystore` | Secrets in the OS credential store (Windows Credential Manager, macOS Keychain, Linux Secret Service), one service name per server |
| `dotenv` | Reads `.env` files and chains settings: environment, then `.env`, then the server's own sources |
| `logging` | Logging to stderr (stdout is the MCP protocol) and panic logging that never logs user data |

## Use

```sh
go get github.com/dvalfrid/mcp-kit@latest
```

```go
store := keystore.OS{Service: "proton-calendar-mcp"}
token, err := store.Get("refresh-token") // a secret.Secret: printing it shows [redacted]
```

## Rules

What goes into `mcp-kit`, what stays in a server, and how versions work: [docs/dependency-rules.md](docs/dependency-rules.md).

## Development

```powershell
./tools/task.ps1 setup    # git hooks, once per clone
./tools/task.ps1 verify   # tidy, format, lint, test, vulnerabilities
```

See [CONTRIBUTING.md](CONTRIBUTING.md) and [STANDARDS.md](STANDARDS.md).

## License

[MIT](LICENSE) © Code By AB
