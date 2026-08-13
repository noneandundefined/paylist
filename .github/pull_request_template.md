## Summary

<!-- Describe what this PR changes and why. -->

## PR title (important)

Use a **conventional commit** title — it becomes the release version on merge:

| Title prefix | Release bump | Example |
|--------------|--------------|---------|
| `fix:` | patch (0.1.0 → 0.1.1) | `fix: correct subscription summary currency` |
| `feat:` | minor (0.1.0 → 0.2.0) | `feat: add category donut chart to analytics` |
| `feat!:` or `BREAKING CHANGE:` in body | major (0.1.0 → 1.0.0) | `feat!: change tracked subscription list API` |
| `chore:`, `ci:`, `docs:` | no release tag | `chore: update workflows` |

> Use **Squash merge** so the PR title becomes the commit on `main`.

## Checklist

- [ ] PR title follows the table above
- [ ] Backend builds (`go build ./...` in `http-server/`)
- [ ] Frontend lint passes (`pnpm lint` in `www-client/`)
- [ ] Frontend build passes (`pnpm build` in `www-client/`)
- [ ] Locales updated for all supported languages (en, ru, de, es) if UI text changed
- [ ] README or deployment docs updated if behavior or setup changed

## Related issues

<!-- Closes #123 -->

<!-- paylist-preview -->
<!-- /paylist-preview -->
