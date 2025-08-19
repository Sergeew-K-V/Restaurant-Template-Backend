# Restaurant Backend

A Go-based backend service for restaurant management.

## Features

- User authentication
- Database migrations
- RESTful API endpoints

## CLI Commands

The application now has a clear separation of concerns:

### Database Migrations (migrate.go)

Run all pending database migrations:

```bash
go run migrate.go
```

Check migration status:

```bash
go run migrate.go -migrate-status
```

Show help:

```bash
go run migrate.go -help
```

Show version:

```bash
go run migrate.go -version
```

### HTTP Server (main.go)

Start the HTTP server:

```bash
go run src/main.go
```

## Available Commands

### migrate.go - Database Migration Tool

- `go run migrate.go` - Run all pending migrations (default)
- `go run migrate.go -migrate` - Run migrations explicitly
- `go run migrate.go -migrate-status` - Show migration status
- `go run migrate.go -version` - Show version information
- `go run migrate.go -help` - Show help message

### main.go - HTTP Server

- `go run src/main.go` - Start the HTTP server

## Development

### Prerequisites

- Go 1.19+
- PostgreSQL database

### Setup

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Configure environment variables
4. Run migrations: `go run migrate.go`
5. Start server: `go run src/main.go`

### Database Migrations

Migrations are stored in `src/database/migrations/` and follow the naming convention:
`{version}-{description}.sql` (e.g., `1.0-create_users_table.sql`)

The migration system automatically:

- Creates a `migrations` table to track executed migrations
- Runs migrations in version order
- Skips already executed migrations
- Provides transaction safety

### Project Structure

- `migrate.go` - Database migration tool with CLI interface
- `src/main.go` - HTTP server application
- `src/database/migrator.go` - Migration logic and database operations
- `src/database/migrations/` - SQL migration files
