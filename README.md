<div align="center">

# 🌤️ Go API Aggregator

### AI-Powered Weather Assistant

[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-29.4.3-2496ED?style=flat-square&logo=docker&logoColor=white)](https://docker.com)
[![GitHub Actions](https://img.shields.io/badge/CI%2FCD-GitHub_Actions-2088FF?style=flat-square&logo=github-actions&logoColor=white)](https://github.com/features/actions)
[![HTTPS](https://img.shields.io/badge/HTTPS-Let's_Encrypt-003A70?style=flat-square&logo=letsencrypt&logoColor=white)](https://letsencrypt.org)
[![Live](https://img.shields.io/badge/Live-agungperdana.store-22C55E?style=flat-square&logo=vercel&logoColor=white)](https://agungperdana.store)

A lightweight REST API that aggregates real-time weather data and Google Gemini AI to deliver contextual outdoor activity recommendations — deployed on a 512MB VPS with full CI/CD automation.

**[🌐 Live Demo](https://agungperdana.store)**

</div>

---

## ✨ Features

- **Real-time Weather** — fetch current temperature, humidity, and conditions for any city worldwide
- **AI-Powered Advice** — Google Gemini AI generates contextual activity recommendations based on live weather
- **Aggregator Architecture** — combines multiple external APIs into a single clean endpoint
- **Ultra-Lightweight** — Go binary uses ~1 MB RAM at idle inside a minimal Alpine Docker container
- **Full HTTPS** — SSL certificate via Let's Encrypt with automatic renewal
- **Zero-Downtime CI/CD** — auto build & deploy on every `git push` via GitHub Actions

---

## 🏗️ Architecture

```
User Browser
     │
     │  HTTPS :443
     ▼
┌─────────────────────┐
│   Nginx (Reverse    │  ← SSL termination
│   Proxy + SSL)      │
└─────────┬───────────┘
          │ proxy_pass :3000
          ▼
┌─────────────────────┐
│  Docker Container   │  ← memory: 150MB limit
│    Go App Binary    │     cpu: 0.5 cores
└────────┬────────────┘
         │
    ┌────┴────┐
    ▼         ▼
Weather    Gemini
  API      AI API
```

**CI/CD Pipeline:**
```
git push → GitHub Actions Runner (7GB RAM)
               │
               ├── docker build (multi-stage)
               ├── push image → ghcr.io
               └── SSH → VPS → docker pull → restart
```

---

## 🚀 API Endpoints

| Method | Endpoint | Parameter | Response |
|--------|----------|-----------|----------|
| `GET` | `/health` | — | `OK` |
| `GET` | `/api/weather` | `?city=Jakarta` | JSON weather data |
| `GET` | `/api/summary` | `?city=Jakarta` | Weather + AI advice |
| `POST` | `/api/ai` | `{"prompt": "..."}` | AI response |
| `GET` | `/` | — | Web UI |

### Example Responses

**`GET /api/weather?city=Jakarta`**
```json
{
  "city": "Jakarta",
  "temperature": 29.85,
  "description": "broken clouds"
}
```

**`GET /api/summary?city=Bali`**
```json
{
  "city": "Bali",
  "temperature": 28.1,
  "description": "few clouds",
  "advice": "Sangat cocok untuk aktivitas outdoor. Suhu 28°C tergolong nyaman..."
}
```

---

## 🛠️ Tech Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| Language | Go 1.22 | High-performance backend |
| Container | Docker + Alpine | Lightweight, portable runtime |
| Registry | GitHub Container Registry | Docker image storage |
| CI/CD | GitHub Actions | Automated build & deploy |
| Web Server | Nginx 1.24 | Reverse proxy & SSL termination |
| SSL | Let's Encrypt + Certbot | Free HTTPS, auto-renew |
| VPS | DigitalOcean (Ubuntu 24.04) | 512 MB RAM production server |
| AI | Google Gemini Flash Latest| Contextual weather advice |
| Weather | OpenWeatherMap API | Real-time weather data |

---

## 📁 Project Structure

```
go-api-aggregator/
├── cmd/
│   └── server/
│       └── main.go          # Entry point, routing
├── internal/
│   ├── client/
│   │   ├── ai_client.go     # Gemini AI API client (10s timeout)
│   │   └── weather_client.go
│   ├── config/
│   │   └── config.go        # Env config loader
│   ├── handler/
│   │   ├── ai.go            # POST /api/ai
│   │   ├── health.go        # GET /health
│   │   └── weather.go       # GET /api/weather, /api/summary
│   └── service/
│       ├── ai_service.go
│       └── weather_service.go
├── web/
│   └── index.html           # Frontend UI
├── .github/
│   └── workflows/
│       └── deploy.yml       # CI/CD pipeline
├── Dockerfile               # Multi-stage build
├── go.mod
├── go.sum
└── .env.example             # Environment template
```

---

## ⚡ Performance on 512MB VPS

One of the key challenges of this project was running a full production stack on a constrained 512MB VPS. Here's how it was achieved:

| Component | RAM Usage |
|-----------|-----------|
| OS + systemd | ~150 MB |
| Docker daemon | ~30 MB |
| Nginx | ~20 MB |
| **Go App (hard limit)** | **≤ 150 MB** |
| DigitalOcean agent | ~10 MB |
| **Buffer remaining** | **~98 MB** |

**Key strategies:**
- **Build on GitHub runner** (7GB RAM), not on VPS — avoids build-time RAM spikes
- **Multi-stage Dockerfile** with `-ldflags="-w -s"` reduces binary size ~30%
- **Hard container limits** via `--memory=150m --cpus=0.5` prevent resource exhaustion
- **Go's efficiency** — only ~1 MB RAM at idle vs ~50 MB for Node.js equivalent
- **Auto image pruning** after every deploy to prevent disk bloat

---

## 🔧 Local Development

### Prerequisites
- Go 1.22+
- Docker (optional)

### Setup

```bash
# Clone the repository
git clone https://github.com/AgungPerdana-IT/go-api-aggregator.git
cd go-api-aggregator

# Copy env template
cp .env.example .env
# Fill in your API keys in .env

# Run locally
go run cmd/server/main.go

# Or with Docker
docker build -t aggregator .
docker run -p 3000:3000 --env-file .env aggregator
```

### Environment Variables

```env
PORT=3000
WEATHER_API_KEY=your_openweathermap_key
AI_API_KEY=your_gemini_api_key
```

Get your free API keys:
- **Weather**: [openweathermap.org](https://openweathermap.org/api)
- **AI**: [aistudio.google.com](https://aistudio.google.com)

---

## 🚢 Deployment

This project uses **GitHub Actions** for fully automated CI/CD.

### How it works

Every push to `main` triggers:

1. **Build** — `docker buildx build` on GitHub runner (not on VPS)
2. **Push** — image pushed to `ghcr.io`
3. **Deploy** — SSH into VPS, pull new image, restart container

### Required GitHub Secrets

| Secret | Description |
|--------|-------------|
| `VPS_HOST` | VPS IP address |
| `VPS_USERNAME` | SSH username (e.g. `root`) |
| `VPS_SSH_PRIVATE_KEY` | SSH private key for VPS access |

### VPS Setup (one-time)

```bash
# Install Docker
apt update && apt install -y docker-ce docker-ce-cli containerd.io

# Clone repo (for .env file location)
git clone https://github.com/AgungPerdana-IT/go-api-aggregator.git

# Create .env
nano /root/go-api-aggregator/.env

# Install Nginx + SSL
apt install nginx certbot python3-certbot-nginx -y
certbot --nginx -d yourdomain.com
```

---

## 💾 Backup & Restore

Even though the app is stateless, automated `.env` backup is configured as a production simulation:

```bash
# Backup runs every 12 hours via crontab
0 */12 * * * /root/backup.sh

# Backups stored at /root/backup/.env.YYYY-MM-DD-HH-MM-SS
# Auto-deleted after 3 days

# Restore from latest backup
bash /root/restore.sh
```

---

## 📄 License

MIT License — feel free to use this as a reference for your own projects.

---

<div align="center">

**Built with Go · Deployed on DigitalOcean · Automated with GitHub Actions**

[🌐 Live Demo](https://agungperdana.store) · [📁 Source Code](https://github.com/AgungPerdana-IT/go-api-aggregator)

</div>

<div align="center">
## 👤 Author

**[Agung Perdana]**
QA Engineer | Manual & Automation Testing

[![LinkedIn](https://img.shields.io/badge/LinkedIn-Connect-0077B5?style=flat-square&logo=linkedin)](https://linkedin.com/in/agung-perdana-it)
[![GitHub](https://img.shields.io/badge/GitHub-Follow-181717?style=flat-square&logo=github)](https://github.com/AgungPerdana-IT)

</div>