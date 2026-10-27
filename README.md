🟢 PulseGuard
Know when your websites go down — before your users do.
Concurrent Website & API Health Monitoring | Go 1.22+ | Responsive Web Dashboard
A lightweight, local-first monitoring project built with the Go standard library. Track HTTP availability, inspect response times and manage your monitored endpoints from a browser.

Status: Functional educational MVP  ·  Language: Go  ·  License: MIT
✨ Highlights
Feature	What it does
🌐 Endpoint monitoring	Add HTTP/HTTPS websites or API URLs to monitor
⚡ Concurrent checks	Monitor up to 5 sites simultaneously in one cycle
📊 Dashboard	Responsive dark UI with latest status, uptime and latency
🕒 Monitoring history	Keep the most recent 48 health checks per endpoint
💾 Local persistence	Save configured targets and check history to JSON
🚨 Log-based alerts	Record failed checks to the server terminal
🔌 REST API	Add, list and remove monitoring targets
🧪 Automated tests	Go unit tests and GitHub Actions CI


🧰 Tech stack
Go · net/http · goroutines · JSON · HTML · CSS · JavaScript · GitHub Actions

No external Go dependencies are required.
🚀 Getting started
Prerequisite: Go 1.22 or later. Confirm your version:
go version
Clone this repository, enter the folder and run:
git clone https://github.com/YOUR_USERNAME/pulseguard-go.git
cd pulseguard-go
go run .
Visit http://127.0.0.1:8080 and add an endpoint such as https://example.com.
Customize the port and check interval
On Linux or macOS:
PORT=9090 CHECK_INTERVAL_SECONDS=15 go run .
On Windows PowerShell:
$env:PORT="9090"
$env:CHECK_INTERVAL_SECONDS="15"
go run .
🧠 How it works
Your Browser ───────────────► Go HTTP Server ─────► Local JSON Store
    │                               │
    └──── GET /api/targets ◄────────┤
                                    └────► Concurrent HTTP checks
                                                │
                                     HTTP status + latency
                                                │
                                         Dashboard history
Monitoring semantics: HTTP status codes from 200 through 399 are treated as online. Each check tracks its timestamp, status, response time and any connection error. Uptime is computed from the stored rolling sample, not a long-term SLA.
🔌 REST API
Method	Endpoint	Purpose
GET	/api/targets	Get configured targets and recent checks
POST	/api/targets	Register a website or API
DELETE	/api/targets/{id}	Remove an existing target


Create a target with curl:
curl -X POST http://127.0.0.1:8080/api/targets \
  -H 'Content-Type: application/json' \
  -d '{"name":"Example","url":"https://example.com"}'
📁 Project structure
pulseguard-go/
├── .github/workflows/ci.yml    # Automated Go tests
├── web/index.html              # Embedded dashboard
├── main.go                     # Server, scheduler, storage and APIs
├── main_test.go                # Unit tests
├── go.mod
├── .gitignore
└── LICENSE
The application creates data/sites.json at runtime. This is intentionally excluded from Git.
✅ Testing
go test ./...
go vet ./...
The repository also runs these checks in GitHub Actions on push and pull requests.
🗺️ Roadmap (not implemented yet)
- [ ] Email or Telegram notification integrations
- [ ] SQLite and extended historical analytics
- [ ] Authentication and role-based accounts
- [ ] Safer externally hosted deployments with network restrictions
- [ ] Downloadable SLA reports and uptime charts
🔒 Important limitations
This is a local development tool, not a hardened hosted monitoring platform. The server binds to 127.0.0.1 by default. Its URL monitoring capability should not be exposed publicly without authentication, network egress controls and SSRF protection. Alerting is currently in the terminal logs, not email. A local instance supports a maximum of 25 monitored URLs.
📸 Screenshots
Run the application locally, capture your own dashboard and add screenshots here. No fabricated screenshots are included.
🤝 Contributions & license
Ideas, bug reports and improvements are welcome. Licensed under the MIT License — see LICENSE.
GitHub description: A Go-powered website and API uptime monitor with concurrent health checks, latency tracking, and a responsive dashboard.
