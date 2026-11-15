## Summary

<!-- What does this change, and why? -->

Closes #

## Testing

<!-- Which tests cover it? -->

## Checklist

- [ ] Commit subject follows [Conventional Commits](https://www.conventionalcommits.org/) (`type(scope): subject`)
- [ ] `./tools/task.ps1 verify` passes locally
- [ ] New behaviour has tests; a bug fix has a test that failed before the fix
- [ ] No guardrail weakened ([STANDARDS.md — Guardrails](../STANDARDS.md#guardrails)); no network code, no service code ([dependency rules](../docs/dependency-rules.md))
- [ ] A breaking change is marked `!` and says what servers must change
- [ ] No new dependency, or it is approved and listed in STANDARDS.md
