# todoApp

A todo REST API in Go (gorilla/mux + GORM + PostgreSQL). See `go-postgres/`.

## Configuration

Database settings are read from environment variables. For local runs:

```bash
cp go-postgres/.env.example go-postgres/.env
# edit go-postgres/.env with your own values
```

`.env` is git-ignored and must never be committed.
