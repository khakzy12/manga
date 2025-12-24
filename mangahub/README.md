# MangaHub (Distributed System Demo)

[![CI](https://github.com/<OWNER>/<REPO>/actions/workflows/ci.yml/badge.svg)](https://github.com/<OWNER>/<REPO>/actions/workflows/ci.yml)

> ⚠️ Replace `<OWNER>/<REPO>` in the badge URL above with your GitHub repository path (e.g., `your-username/mangahub`) so the badge shows your repository status.

🔧 Quick project overview and developer notes

## ✅ Seeded demo accounts
The seed script (`cmd/seed/main.go`) now adds two users for convenience:

- **admin** / **adminpass**  (role: `admin`) ✅
- **demo** / **demopass**    (role: `user`) ✅

Each seeded user includes a generated `id` (UUID) and a bcrypt-hashed password.

## ▶️ How to seed the database
From the project root run:

```bash
go run cmd/seed/main.go
```

This will create `data/mangahub.db` (if missing), create tables, and insert seed data (manga + users).

You can also run the bundled shortcut on Windows:

```powershell
./run_all.bat
```

## 🔐 Auth endpoints
- Register: `POST /auth/register`
  - Body: `{ "username": "<username>", "password": "<password>" }`
  - Success: `201 Created` with JSON `{ "message": "User created", "id": "<uuid>" }`

- Login: `POST /auth/login`
  - Returns a JWT token on success.

Example (PowerShell):

```powershell
# Register
Invoke-RestMethod -Uri 'http://localhost:8080/auth/register' -Method Post -ContentType 'application/json' -Body (ConvertTo-Json @{username='newuser'; password='s3cret'})

# Login
Invoke-RestMethod -Uri 'http://localhost:8080/auth/login' -Method Post -ContentType 'application/json' -Body (ConvertTo-Json @{username='demo'; password='demopass'})
```

Or using `curl` (on environments where `curl` behaves normally):

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"s3cret"}'
```

## Troubleshooting
- If you get `400` on register, check server logs for `❌ Register JSON Bind Error:` — often malformed JSON or missing fields.
- If `username already exists` is returned, try a different username or remove the DB (`data/mangahub.db`) and re-run the seed script to start fresh.

---

## Continuous Integration ✅
A GitHub Actions workflow is included at `.github/workflows/ci.yml`. It runs on push and pull request and performs:

- `gofmt` check (fails if formatting issues are found)
- `go vet` for static checks
- `golangci-lint` for linting (configured in `.golangci.yml`)
- `go test ./...` to run the test suite

If you'd like, I can expand the workflow to run linters (e.g., `golangci-lint`) or run tests on multiple Go versions.

If you'd like, I can also add these details to a `docs/` folder or create a `DEVELOPING.md` with contributing steps and tests. Would you like that? ✨