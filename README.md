# OpenAgent

Cursor extension to run multi-agent workflows (LangGraph first) and manage configs.

## Plan (next steps)
1) **Backend service (Go)**: local HTTP service to run workflows and stream logs.
2) **Extension integration**: command to call backend `/run` and show logs.
3) **Config management**: workspace config file + settings UI entry.
4) **Logs view**: output channel + optional tree/timeline view.

## MVP scope
- Run workflow + show logs
- Manage configs/env

## How to build (current)
- Extension: `npm i` then `npm run compile`
- Backend: `make backend`

## Endpoints (current stub)
- `GET /health`
- `POST /run`
