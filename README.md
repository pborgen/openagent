# OpenAgent (Monorepo)

## Structure
```
apps/
  extension/   # Cursor/VS Code extension
  backend/     # Go backend service
packages/      # future shared packages
```

## Build
- Extension: `pnpm --dir apps/extension compile`
- Backend: `pnpm --dir apps/backend build`

## Config
`apps/extension/openagent.json`
