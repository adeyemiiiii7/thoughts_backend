package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"thoughts_backend_api/db"
	"thoughts_backend_api/services/auth"
	"thoughts_backend_api/services/comments"
	"thoughts_backend_api/services/reactions"
	"thoughts_backend_api/services/thoughts"
	"thoughts_backend_api/shared"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=adeyemialadesuyi dbname=thoughts port=5432 sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-me"
		log.Println("JWT_SECRET not set, using insecure development fallback")
	}

	gormDB, err := db.InitDB(dbURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("failed to get raw database connection: %v", err)
	}
	defer sqlDB.Close()

	r := chi.NewRouter()
	authHandler := auth.NewHandler(gormDB, jwtSecret)
	commentHandler := comments.NewHandler(gormDB)
	reactionHandler := reactions.NewHandler(gormDB)
	thoughtHandler := thoughts.NewHandler(gormDB)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Thoughts backend with GORM is running"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			http.Error(w, "database is not reachable", http.StatusServiceUnavailable)
			return
		}

		w.Write([]byte("ok"))
	})

	r.Post("/auth/signup", authHandler.Signup)
	r.Post("/auth/login", authHandler.Login)
	r.Post("/auth/forgot-password", authHandler.ForgotPassword)
	r.Post("/auth/reset-password", authHandler.ResetPassword)
	r.Get("/auth/verify-email", authHandler.VerifyEmail)
	r.Get("/thoughts", thoughtHandler.List)
	r.Get("/thoughts/{id}/comments", commentHandler.ListByThought)

	r.Group(func(r chi.Router) {
		r.Use(shared.AuthMiddleware(gormDB, []byte(jwtSecret)))
		r.Get("/users/profile", authHandler.GetProfile)
		r.Post("/auth/change-password", authHandler.ChangePassword)
		r.Put("/users/interests", authHandler.UpdateInterests)
		r.Post("/thoughts", thoughtHandler.Create)
		r.Post("/thoughts/{id}/comments", commentHandler.Create)
		r.Post("/thoughts/{id}/reactions", reactionHandler.CreateOrUpdate)
		r.Post("/comments/{id}/replies", commentHandler.ReplyComment)
		r.Delete("/thoughts/{id}", thoughtHandler.Delete)
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("server running on http://localhost:%s", port)
	log.Fatal(server.ListenAndServe())
}
