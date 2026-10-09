 
<div align="center">

<img src="assets/banner.png" alt="Ymirr Banner" width="100%">


<img src="assets/logoG.png" alt="Ymirr Mascot" width="160"> 

# YMIRR

**A lightweight, modular Go toolkit for Turso and HTTP middleware.**

Built with Go · Developed by **CRAJ Labs**

</div>

---

## About Ymirr

**Ymirr** is a lightweight Go library designed to simplify repetitive backend development tasks.

Inspired by Ymir, the primordial giant of Norse mythology, Ymirr aims to provide foundational building blocks for modern Go applications.

The project focuses on three principles:

- **Simplicity:** Minimal configuration and straightforward APIs.
- **Modularity:** Import only the packages you need.
- **Flexibility:** Choose between convenient defaults and advanced configuration.

## Features

### Turso Integration

- Turso/libSQL database connectivity.
- Environment-based configuration.
- Custom connection settings.
- SQLX integration.
- Connection pooling configuration.

### HTTP Middleware

- Gin-compatible authentication middleware.
- Session-based authentication.
- Configurable authentication options.
- Extensible session validation interfaces.
- CSRF protection integration.

## Installation

Install Ymirr in your Go project:

```bash
go get github.com/IntentoDeIngenieroAgustin/Ymirr@latest
```

## Quick Start

### Turso — Simple Configuration

Configure your environment variables:

```env
TURSO_DATABASE_URL=libsql://your-database.turso.io
TURSO_AUTH_TOKEN=your-auth-token
```

Connect to Turso:

```go
package main

import (
    "context"
    "log"

    "github.com/IntentoDeIngenieroAgustin/ymirr/turso"
)

func main() {
    db, err := turso.ConnectFromEnv(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    log.Println("Connected to Turso")
}
```

### Turso — Advanced Configuration

```go
db, err := turso.Connect(ctx, turso.Config{
    URL:             os.Getenv("TURSO_DATABASE_URL"),
    Token:           os.Getenv("TURSO_AUTH_TOKEN"),
    MaxOpenConns:    10,
    MaxIdleConns:    5,
    ConnMaxLifetime: 30 * time.Minute,
})
```

### Authentication Middleware

Ymirr provides reusable authentication middleware for Gin.

```go
router := gin.Default()

// authService must implement middleware.SessionValidator.
router.Use(middleware.RequireAuth(authService))
```

Custom configuration:

```go
router.Use(middleware.RequireAuthWithConfig(
    authService,
    middleware.AuthConfig{
        CookieName: "session_id",
        ContextKey: "user_id",
    },
))
```

Session validation is implemented by the consuming application.

For cookie-based authentication, configure CSRF protection separately using a validated token strategy.

## Project Structure

```text
ymirr/
├── assets/
│   ├── banner.png
│   └── logo.png
├── turso/
│   ├── config.go
│   └── connection.go
├── middleware/
│   ├── auth.go
│   └── csrf.go
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

## Development

Clone the repository:

```bash
git clone https://github.com/IntentoDeIngenieroAgustin/ymirr.git
cd ymirr
```

Install dependencies:

```bash
go mod tidy
```

Run tests:

```bash
go test ./...
```

## Roadmap

- [x] Initial Turso integration
- [x] Configurable database connections
- [x] Authentication middleware prototype
- [ ] Complete CSRF integration and security tests
- [ ] Database integration tests
- [ ] Additional HTTP middleware
- [ ] Database migrations
- [ ] Extended documentation
- [ ] Rust-powered project generator CLI

## License

Ymirr is intended to be distributed under the MIT License. See the [LICENSE](LICENSE) file once added to the repository.

---

<div align="center">

**YMIRR — From nothing, build everything.**

Made by [CRAJ Labs](https://www.crajlabs.tech/)

</div>
