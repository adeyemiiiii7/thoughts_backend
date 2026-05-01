# Thoughts Backend

Simple Go starter project with:

- `chi` for routing
- PostgreSQL for the database
- `gorm` for working with Go structs as database tables
- Homebrew PostgreSQL for local database setup

## 1. Start PostgreSQL

```bash
brew services start postgresql@18
```

## 2. Set your environment variables

Copy `.env.example` into `.env`:

```bash
cp .env.example .env
```

## 3. Install dependencies

```bash
go mod tidy
```

## 4. Run the server

```bash
go run .
```

## Routes

- `GET /` returns a simple welcome message
- `GET /health` checks if the app can reach PostgreSQL

## GORM in this project

The model lives in `models/thought.go`.

When the app starts, GORM runs `AutoMigrate(&models.Thought{})`.
That means it checks whether the `thoughts` table exists and creates or updates it to match your model.

The database connection is kept simple for learning:

```env
DATABASE_URL=postgresql://your_username@localhost:5432/thoughts?sslmode=disable
```
