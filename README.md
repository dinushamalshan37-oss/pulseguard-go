
# 🟢 PulseGuard
### Smart Website & API Monitoring System

> Monitor your websites. Track performance. Detect downtime.

PulseGuard is a lightweight website and API monitoring application developed using **Go (Golang)**. It helps developers monitor website availability, analyze response times, and inspect endpoint health through a modern web dashboard.

---

## 🚀 Features

- 🌐 Website and API health monitoring
- ⚡ Concurrent HTTP requests using Go
- 📊 Modern responsive monitoring dashboard
- 🟢 Online and offline status detection
- ⏱️ Response time and latency tracking
- 📈 Recent uptime statistics
- 💾 Local JSON data persistence
- 🔌 REST API for managing endpoints
- 🚨 Console-based downtime alerts
- 🧪 Automated unit testing

## 🛠️ Technologies Used

| Technology | Purpose |
|---|---|
| Go | Backend development |
| net/http | HTTP requests and web server |
| Goroutines | Concurrent monitoring |
| JSON | Local data storage |
| HTML5 | Dashboard structure |
| CSS3 | Dashboard styling |
| JavaScript | Frontend interactions |
| GitHub Actions | Continuous integration |

## 📁 Project Structure

```text
pulseguard-go/
├── main.go
├── main_test.go
├── go.mod
├── web/
│   └── index.html
├── .github/
│   └── workflows/
│       └── ci.yml
├── .gitignore
├── LICENSE
└── README.md
```

## ⚙️ Installation

### Requirements
- Go 1.22 or newer
- Git
- Modern web browser

### Clone Repository

```bash
git clone https://github.com/YOUR_USERNAME/pulseguard-go.git
cd pulseguard-go
```

### Run Application

```bash
go run .
```

Open your browser:

```text
http://127.0.0.1:8080
```

## 🔌 REST API

| Method | Endpoint | Description |
|---|---|---|
| GET | /api/targets | List monitored websites |
| POST | /api/targets | Add a monitoring target |
| DELETE | /api/targets/{id} | Delete a monitoring target |

### Example Request

```bash
curl -X POST http://127.0.0.1:8080/api/targets \
-H "Content-Type: application/json" \
-d '{"name":"Example","url":"https://example.com"}'
```

## 🧪 Testing

```bash
go test ./...
go vet ./...
```

## 🧠 How It Works

1. Users add website or API URLs.
2. Go performs scheduled HTTP health checks.
3. Response status and latency are recorded.
4. Monitoring results are stored locally.
5. The dashboard displays recent statistics.
6. Failed checks are reported in server logs.

## 🔮 Future Improvements

- [ ] Email and Telegram notifications
- [ ] SQLite database integration
- [ ] Authentication and user accounts
- [ ] Advanced analytics and charts
- [ ] Exportable uptime reports
- [ ] Cloud deployment with security controls

## ⚠️ Limitations

This project is an educational MVP intended for local development. Monitoring history is limited, alerts are console-based, and authentication is not implemented.

## 🤝 Contributing

Contributions, suggestions, and improvements are welcome.

## 📜 License

This project is licensed under the MIT License.

---

### ⭐ Support

If you find this project useful, consider giving it a star!

**Built with Go | Backend Development | DevOps | Monitoring**
