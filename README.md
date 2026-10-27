PulseGuard 🟢
Smart API & website health monitor • Go • local-first

A working, zero-dependency HTTP monitoring dashboard written in Go. Add URLs, monitor uptime and response time, inspect recent checks and receive alert events in the terminal. Built as a student portfolio MVP.

Features
Concurrent checks (up to 5 simultaneously) every 30 seconds by default
HTTP status, latency, last 48 checks per URL and rolling uptime percent
Dark responsive web dashboard
JSON persistence in data/sites.json (no DB setup needed)
Local logging for failures, REST API, graceful shutdown, GitHub Actions tests
Quick start
Requires: Go 1.22+

go run .
Open http://127.0.0.1:8080, add a URL and allow the first monitoring cycle. For a custom interval / port:

PORT=9090 CHECK_INTERVAL_SECONDS=15 go run .
For Windows PowerShell:

$env:PORT="9090"; $env:CHECK_INTERVAL_SECONDS="15"; go run .
API
Method	Path	Purpose
GET	/api/targets	List monitors and recent checks
POST	/api/targets	Add: {"name":"Example","url":"https://example.com"}
DELETE	/api/targets/{id}	Remove monitor
curl -X POST http://127.0.0.1:8080/api/targets -H 'Content-Type: application/json' -d '{"name":"Example","url":"https://example.com"}'
Tests
go test ./...
go vet ./...
Important limitations / next milestones
Local development only. The app deliberately binds to 127.0.0.1; monitoring arbitrary URLs on a public server without authentication and network restrictions exposes SSRF risks.
Alerts appear in server logs, not email. Add a notification provider and tests before claiming email alerts.
Results are limited to 48 checks per target; rolling uptime is not a historical SLA calculation.
Future: SQLite, accounts, authenticated multi-user deployments, advanced charts, email alerts.
Suggested GitHub metadata
Description: Concurrent Go website and API uptime monitor with a responsive dashboard and local persistence. Topics: golang, monitoring, api, devops, dashboard, portfolio-project

License
MIT — see LICENSE.
