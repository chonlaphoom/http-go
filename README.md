[![ci](https://github.com/chonlaphoom/http-go/actions/workflows/ci.yml/badge.svg)](https://github.com/chonlaphoom/http-go/actions/workflows/ci.yml)

**HTTP Go (playful side project)**

- A tiny, opinionated Go HTTP service for chirps/users/auth built as a learning side-project and a neat little demo.

**What it does**
- **Chirps**: create/read/update simple short messages.
- **Users**: registration + auth flow.
- **Auth**: JWT tokens + bcrypt password hashing.

**Quick Start**
- **Clone**: `git clone https://github.com/chonlaphoom/http-go && cd http-go`
- **Setup DB**: run Postgres (v15 recommended) and apply migrations (this repo uses `goose` + `sql` migrations in `sql/schema`).
- **Env**: copy `.env` (or set) variables used by `load_env.go` such as `DATABASE_URL` and `JWT_SECRET`.
- **Build & Run**: `go build -o http-go && ./http-go` or use `./build_and_run.sh`.

**Deps / Tools**
- **Go**: `go 1.24`
- **DB**: `postgresql@15`
- **Migrations**: `goose`
- **SQL gen**: `sqlc`

**Useful files**
- `main.go` - server bootstrap
- `chirps.go`, `users.go`, `json.go` - handlers and helpers
- `internal/database` - DB layer and generated SQL (via `sqlc`)
- `internal/auth` - JWT + password logic
- `sql/schema` - migrations

**Endpoints (high level)**
- `POST /users` - create user
- `POST /login` - get JWT
- `GET/POST /chirps` - list/create chirps

**Notes
- This is a learning playground. Expect small, opinionated choices and readable code.
