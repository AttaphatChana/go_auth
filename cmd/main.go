package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"test_auth/adapters/http_handler"
	"test_auth/adapters/persistence"
	"test_auth/adapters/security"
	"test_auth/application/user"
	"test_auth/config"
	"test_auth/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config error:", err)
	}

	// Database connection
	db, err := pgxpool.New(context.Background(), cfg.DBHost)
	if err != nil {
		db, err = config.RetryDB(cfg)
		if err != nil {
			log.Fatal("RetryDB timeout")
		}
	}
	defer db.Close()
	log.Println("Connected to database")

	// Ensure users table exists
	_, err = db.Exec(context.Background(), `
    CREATE TABLE IF NOT EXISTS users (
      id TEXT PRIMARY KEY,
      username TEXT UNIQUE NOT NULL,
      password TEXT NOT NULL
    );
  `)
	if err != nil {
		log.Fatal("Failed to ensure users table exists:", err)
	}
	log.Println("Ensured users table exists")

	// Redis connection
	//rdb, err := config.RetryRedis(cfg)
	//if err != nil {
	//	log.Fatal("RetryRedis timeout")
	//}
	//defer rdb.Close()
	//log.Println("Connected to Redis")

	// Initialize repositories and services
	userRepo := &persistence.PostgresUserRepo{Pool: db}
	jwtSvc := &security.JwtService{Secret: []byte(cfg.JWTSecret)}

	// Initialize use cases
	regUC := &user.RegisterUser{Repo: userRepo}
	loginUC := &user.LoginUser{Repo: userRepo, TokenService: jwtSvc}

	// Setup HTTP handlers
	authH := &http_handler.AuthHandler{RegisterUC: regUC, LoginUC: loginUC}

	// Setup Gin router
	r := gin.Default()
	r.POST("/signup", authH.Register)
	r.POST("/login", authH.Login)

	// Protected routes
	api := r.Group("/api", middleware.JWTMiddleware(jwtSvc))
	api.GET("/profile", func(c *gin.Context) {
		userID := c.GetString("userID")
		c.JSON(200, gin.H{"userID": userID})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
