package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vivekanand/labour-thekedar-backend/internal/config"
	"github.com/vivekanand/labour-thekedar-backend/internal/database"
	"github.com/vivekanand/labour-thekedar-backend/internal/handler"
	"github.com/vivekanand/labour-thekedar-backend/internal/middleware"
	"github.com/vivekanand/labour-thekedar-backend/internal/repository"
	"github.com/vivekanand/labour-thekedar-backend/internal/service"
	"github.com/vivekanand/labour-thekedar-backend/pkg/otp"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Run migrations
	log.Println("Running database migrations...")
	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")

	// Connect to database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize OTP provider (mock for development)
	otpProvider := otp.NewMockProvider(true)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.Pool)
	projectRepo := repository.NewProjectRepository(db.Pool)
	labourRepo := repository.NewLabourRepository(db.Pool)
	workDayRepo := repository.NewWorkDayRepository(db.Pool)
	paymentRepo := repository.NewPaymentRepository(db.Pool)

	// Initialize services
	authService := service.NewAuthService(userRepo, otpProvider, cfg.JWTSecret)
	projectService := service.NewProjectService(projectRepo, labourRepo)
	labourService := service.NewLabourService(labourRepo)
	workDayService := service.NewWorkDayService(workDayRepo, labourRepo)
	paymentService := service.NewPaymentService(paymentRepo, labourRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService)
	labourHandler := handler.NewLabourHandler(labourService, projectService)
	workDayHandler := handler.NewWorkDayHandler(workDayService, projectService)
	paymentHandler := handler.NewPaymentHandler(paymentService, projectService)

	// Setup router
	r := gin.Default()

	// Add CORS middleware
	r.Use(corsMiddleware())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/send-otp", authHandler.SendOTP)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// Projects
			projects := protected.Group("/projects")
			{
				projects.GET("", projectHandler.List)
				projects.POST("", projectHandler.Create)
				projects.GET("/:id", projectHandler.Get)
				projects.PUT("/:id", projectHandler.Update)
				projects.DELETE("/:id", projectHandler.Delete)

				// Project labours
				projects.GET("/:id/labours", labourHandler.ListByProject)
				projects.POST("/:id/labours", labourHandler.AssignToProject)
				projects.DELETE("/:id/labours/:labour_id", labourHandler.RemoveFromProject)

				// Project attendance
				projects.GET("/:id/attendance", workDayHandler.List)
				projects.POST("/:id/attendance", workDayHandler.Create)

				// Project payments
				projects.GET("/:id/payments", paymentHandler.ListByProject)
				projects.POST("/:id/payments", paymentHandler.Create)

				// Labour balance in project
				projects.GET("/:id/labours/:labour_id/balance", paymentHandler.GetBalance)
			}

			// Labours
			labours := protected.Group("/labours")
			{
				labours.GET("", labourHandler.List)
				labours.POST("", labourHandler.Create)
				labours.GET("/:id", labourHandler.Get)
				labours.PUT("/:id", labourHandler.Update)
				labours.DELETE("/:id", labourHandler.Delete)
				labours.GET("/:id/payments", paymentHandler.ListByLabour)
			}

			// Attendance (for update/delete by ID)
			attendance := protected.Group("/attendance")
			{
				attendance.PUT("/:id", workDayHandler.Update)
				attendance.DELETE("/:id", workDayHandler.Delete)
			}

			// Payments (for delete by ID)
			payments := protected.Group("/payments")
			{
				payments.DELETE("/:id", paymentHandler.Delete)
			}
		}
	}

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
