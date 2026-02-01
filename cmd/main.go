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
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
    // Initialize logger first
    logger := utils.InitLogger()
    logger.Info("Starting Library Management API...")

    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    logger.Info("Configuration loaded successfully")

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
    refreshTokenRepo := repository.NewRefreshTokenRepository(db)

    authService := service.NewAuthService(userRepo, refreshTokenRepo, cfg)
    bookService := service.NewBookService(bookRepo, categoryRepo, likeRepo, savedRepo, commentRepo)

    authHandler := handler.NewAuthHandler(authService)
    bookHandler := handler.NewBookHandler(bookService, cfg)

    // Create Gin router without default middleware
    r := gin.New()

    // Add custom middleware
    r.Use(gin.Recovery()) // Recovery middleware

    r.Use(middleware.CORSMiddleware())
    r.Use(middleware.LoggerMiddleware()) // Custom logger middleware

    logger.Info("Middleware configured successfully")

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    api := r.Group("/api/v1")
    {
        auth := api.Group("/auth")
        auth.Use(middleware.AuthRateLimitMiddleware()) // Rate limit: 5 req/min
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
            auth.POST("/refresh", authHandler.RefreshToken)
            auth.POST("/logout", middleware.AuthMiddleware(cfg, userRepo), authHandler.Logout)
        }

        protected := api.Group("")
        protected.Use(middleware.AuthMiddleware(cfg, userRepo))
        protected.Use(middleware.APIRateLimitMiddleware()) // Rate limit: 100 req/min
        {
            categories := protected.Group("/categories")
            {
                categories.GET("", bookHandler.GetAllCategories)
                categories.POST("",
                    middleware.RoleMiddleware(models.RoleOwner),
                    bookHandler.CreateCategory)
                categories.DELETE("/:id",
                    middleware.RoleMiddleware(models.RoleOwner),
                    bookHandler.DeleteCategory)
            }

            books := protected.Group("/books")
            {
                books.GET("", bookHandler.GetAllBooks)
                books.GET("/category", bookHandler.GetBooksByCategory)
                books.GET("/my-books",
                    middleware.RoleMiddleware(models.RoleMember),
                    bookHandler.GetSavedBooks)
                books.GET("/:id", bookHandler.GetBook)
                books.GET("/:id/download", bookHandler.DownloadBook)
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
    logger.WithField("address", addr).Info("Server starting")
    logger.WithField("url", fmt.Sprintf("http://localhost%s/swagger/index.html", addr)).Info("Swagger documentation available")

    if err := r.Run(addr); err != nil {
        logger.WithError(err).Fatal("Failed to start server")
    }
}


func createSuperAdmin(db *sql.DB) {
	logger := utils.Logger
	email := os.Getenv("SUPER_ADMIN_EMAIL")
	password := os.Getenv("SUPER_ADMIN_PASSWORD")

	if email == "" || password == "" {
		logger.Warn("Super admin credentials not configured, skipping...")
		return // Skip if not configured
	}

	// Check if exists
	userRepo := repository.NewUserRepository(db)
	existing, err := userRepo.FindByEmail(email)
	if err != nil {
		logger.WithError(err).Error("Error checking for existing admin")
		return
	}
	if existing != nil {
		logger.Info("Super admin already exists")
		return // Already exists
	}

	// Create super admin
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		logger.WithError(err).Error("Error hashing password")
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
		logger.WithError(err).Error("Error creating super admin")
		return
	}

	logger.WithField("email", email).Info("✓ Super admin created successfully")
}