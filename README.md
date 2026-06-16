# GOTTH Auction

A real-time auction platform built with Go, where multiple users can place bids simultaneously and see live price updates — no page refreshes needed.

Built as a portfolio project to demonstrate backend engineering with Go: WebSocket concurrency, JWT authentication with token rotation, PostgreSQL via GORM, and server-side rendering with Templ.

**Stack:** Go · Gin · PostgreSQL · GORM · WebSockets · JWT (RS256) · Templ · HTMX · Tailwind CSS

---

## Features

- **Real-time bidding** — WebSocket hub per auction broadcasts live price and bid history to all connected clients
- **JWT authentication** — RS256-signed access/refresh token pair stored in HTTP-only cookies with automatic silent refresh
- **User accounts** — register, login, profile photo upload, bid history, won auctions
- **Auction listings** — active and expired auctions, category filtering
- **Server-side rendering** — type-safe HTML via Templ templates, dynamic partials via HTMX

---

## Architecture

```
gotth-auction/
├── main.go                   # Server bootstrap, route wiring
├── initializers/             # DB connection, env loading (Viper)
├── models/                   # GORM models: User, Auction, Bid, Category
├── controllers/              # Business logic: auth, auctions, WebSocket
├── handlers/                 # HTTP handlers: parse request → call controller → render template
├── routes/                   # Route groups and middleware assignment
├── middleware/               # JWT validation + silent token refresh
├── templates/                # Templ components (compiled to *_templ.go)
├── utils/                    # JWT helpers (RS256), bcrypt, file upload
├── migrate/                  # GORM AutoMigrate
└── seeder/                   # Sample data for local development
```

### WebSocket hub pattern

Each auction gets a dedicated `AuctionHub` goroutine that owns a register/unregister/broadcast channel set. When a user connects to `/ws/:id`, the route layer looks up an existing hub for that auction ID or creates a new one. The hub's `Run()` loop handles all client lifecycle and message fanout, keeping concurrent bid writes safe without explicit locking on the broadcast path.

### JWT token rotation

Access tokens expire in 15 minutes. The `deserialize-user` middleware validates the access token on every request; if expired, it automatically issues a new one from the refresh token (7-day TTL) before the handler runs. Both tokens are RS256-signed using separate key pairs.

---

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL

### 1. Clone and install dependencies

```bash
git clone https://github.com/sinanaas/gotth-auction.git
cd gotth-auction
go mod download
```

### 2. Generate RSA key pairs

The app uses separate RS256 key pairs for access and refresh tokens.

```bash
# Access token keys
openssl genrsa -out access_private.pem 2048
openssl rsa -in access_private.pem -pubout -out access_public.pem

# Refresh token keys
openssl genrsa -out refresh_private.pem 2048
openssl rsa -in refresh_private.pem -pubout -out refresh_public.pem

# Base64-encode them (paste output into app.env)
base64 -i access_private.pem
base64 -i access_public.pem
base64 -i refresh_private.pem
base64 -i refresh_public.pem
```

### 3. Configure environment

Create `app.env` in the project root:

```env
# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
POSTGRES_DB=gotth_auction

# Server
PORT=8000
CLIENT_ORIGIN=http://localhost:8000

# JWT — paste base64-encoded PEM content from step 2
ACCESS_TOKEN_PRIVATE_KEY=
ACCESS_TOKEN_PUBLIC_KEY=
REFRESH_TOKEN_PRIVATE_KEY=
REFRESH_TOKEN_PUBLIC_KEY=

# Token expiry
ACCESS_TOKEN_EXPIRED_IN=15m
REFRESH_TOKEN_EXPIRED_IN=168h
ACCESS_TOKEN_MAXAGE=900
REFRESH_TOKEN_MAXAGE=604800

# Sessions
SESSION_SECRET_KEY=change-me
```

### 4. Enable the UUID extension

GORM uses `uuid_generate_v4()` for primary keys, which requires the `uuid-ossp` PostgreSQL extension:

```bash
psql -U postgres -d gotth-auction -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp";'
```

### 5. Run migrations and seed data

```bash
go run migrate/migrate.go
go run seeder/seeder.go   # optional — loads sample auctions and a test user
```

### 5. Start the server

```bash
go run main.go
```

Open [http://localhost:8000](http://localhost:8000).

For hot-reloading during development, install [Air](https://github.com/air-verse/air) and run `air`.

---

## API / Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Home — active and expired auctions |
| GET | `/auction/:id` | Auction detail page |
| GET/POST | `/login` | Login |
| GET/POST | `/register` | Register |
| GET/POST | `/profile` | View and update profile |
| GET | `/history` | User's won auctions |
| GET | `/about` | About page |
| GET | `/ws/:id` | WebSocket endpoint for real-time bidding |

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.22 |
| HTTP framework | Gin |
| Database | PostgreSQL |
| ORM | GORM |
| Templating | Templ |
| Real-time | gorilla/websocket |
| Auth | golang-jwt (RS256) |
| Frontend interactivity | HTMX |
| Styling | Tailwind CSS |
| Config | Viper |
