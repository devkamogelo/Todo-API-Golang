package main

import (
	"log"
	"todo_api/internal/config"
	"todo_api/internal/database"
	"todo_api/internal/handlers"
	"todo_api/internal/middleware"
	"todo_api/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.Load()

	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	rds := store.NewRedis(cfg)

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	defer pool.Close()

	var router = gin.Default()

	err = router.SetTrustedProxies(nil)
	if err != nil {
		return
	}
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "To Do Api is running",
			"status":   "success",
			"database": "connected",
		})
	})

	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg, rds))
	router.POST("/auth/refresh", handlers.RefreshHandler(cfg, rds))
	router.POST("/auth/logout", handlers.LogoutHandler(cfg, rds))

	protected := router.Group("/todos")
	protected.Use(middleware.AuthMiddleware(cfg, rds))

	{
		protected.POST("", handlers.CreateTodoHandler(pool))
		protected.GET("", handlers.GetAllTodosHandler(pool))
		protected.GET("/:id", handlers.GetTodoByID(pool))
		protected.PUT("/:id", handlers.UpdateTodo(pool))
		protected.DELETE("/:id", handlers.DeleteTodo(pool))
	}

	router.GET("/protected", middleware.AuthMiddleware(cfg, rds), handlers.TestProtectedHandler())

	err = router.Run(":" + cfg.Port)
	if err != nil {
		return
	}
}
