package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "library-project/config"
	_ "library-project/docs"
    "library-project/internal/handler"
    "library-project/internal/middleware"
    "library-project/internal/models"
    "library-project/internal/repository"
    "library-project/internal/service"
	"library-project/internal/utils"

    "github.com/gin-gonic/gin"
	"github.com/google/uuid" 
    _ "github.com/lib/pq"
    
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Library Management API
// @version 1.0
// @description Library Management System API with JWT authentication
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    db, err := sql.Open("postgres", cfg.Database.DSN())
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

	createSuperAdmin(db)

    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    if err := os.MkdirAll(cfg.Upload.Path, os.ModePerm); err != nil {
        log.Fatalf("Failed to create upload directory: %v", err)
    }

    userRepo := repository.NewUserRepository(db)
    bookRepo := repository.NewBookRepository(db)
    categoryRepo := repository.NewCategoryRepository(db)
    commentRepo := repository.NewCommentRepository(db)
    likeRepo := repository.NewLikeRepository(db)
    savedRepo := repository.NewSavedBookRepository(db)

    authService := service.NewAuthService(userRepo, cfg)
    bookService := service.NewBookService(bookRepo, categoryRepo, likeRepo, savedRepo, commentRepo)

    authHandler := handler.NewAuthHandler(authService)
    bookHandler := handler.NewBookHandler(bookService, cfg)

    r := gin.Default()

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    api := r.Group("/api/v1")
    {
        auth := api.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
        }

        protected := api.Group("")
        protected.Use(middleware.AuthMiddleware(cfg))
        {
            categories := protected.Group("/categories")
            {
                categories.GET("", bookHandler.GetAllCategories)
                categories.POST("", 
                    middleware.RoleMiddleware(models.RoleOwner), 
                    bookHandler.CreateCategory)
            }

            books := protected.Group("/books")
            {
                books.GET("", bookHandler.GetAllBooks)
                books.GET("/category", bookHandler.GetBooksByCategory)
                books.GET("/saved", 
                    middleware.RoleMiddleware(models.RoleMember), 
                    bookHandler.GetSavedBooks)
                books.GET("/:id", bookHandler.GetBook)
                books.POST("", 
                    middleware.RoleMiddleware(models.RoleOwner), 
                    bookHandler.CreateBook)
                books.PUT("/:id", 
                    middleware.RoleMiddleware(models.RoleOwner), 
                    bookHandler.UpdateBook)
                books.DELETE("/:id", 
                    middleware.RoleMiddleware(models.RoleOwner), 
                    bookHandler.DeleteBook)
                books.POST("/:id/save", 
                    middleware.RoleMiddleware(models.RoleMember), 
                    bookHandler.SaveBook)
                books.DELETE("/:id/unsave", 
                    middleware.RoleMiddleware(models.RoleMember), 
                    bookHandler.UnsaveBook)
                books.POST("/:id/like", bookHandler.LikeBook)
                books.GET("/:id/comments", bookHandler.GetComments)
                books.POST("/:id/comments", bookHandler.AddComment)
            }
        }
    }

    addr := fmt.Sprintf(":%s", cfg.Server.Port)
    log.Printf("Server starting on %s", addr)
    log.Printf("Swagger documentation: http://localhost%s/swagger/index.html", addr)
    
    if err := r.Run(addr); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}


func createSuperAdmin(db *sql.DB) {
	email := os.Getenv("SUPER_ADMIN_EMAIL")
	password := os.Getenv("SUPER_ADMIN_PASSWORD")

	if email == "" || password == "" {
		log.Println("Super admin credentials not configured, skipping...")
		return // Skip if not configured
	}

	// Check if exists
	userRepo := repository.NewUserRepository(db)
	existing, err := userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("Error checking for existing admin: %v", err)
		return
	}
	if existing != nil {
		log.Println("Super admin already exists")
		return // Already exists
	}

	// Create super admin
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return
	}

	user := &models.User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  hashedPassword,
		FirstName: "Super",
		LastName:  "Admin",
		Role:      models.RoleOwner,
	}

	if err := userRepo.Create(user); err != nil {
		log.Printf("Error creating super admin: %v", err)
		return
	}

	log.Printf("✓ Super admin created: %s", email)
}