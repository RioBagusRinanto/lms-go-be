package main

import (
	"fmt"
	"log"

	"lms-go-be/internal/config"
	"lms-go-be/internal/database"
	"lms-go-be/internal/handler"
	"lms-go-be/internal/middleware"
	"lms-go-be/internal/repository"
	"lms-go-be/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db := database.InitDB(cfg)

	// Check database connection
	if err := database.CheckConnection(db); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Seed database with initial data (comment out after first run)
	if err := database.Seed(db); err != nil {
		log.Printf("Warning: Database seeding failed: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	enrollmentRepo := repository.NewEnrollmentRepository(db)
	userProgressRepo := repository.NewUserProgressRepository(db)
	quizRepo := repository.NewQuizRepository(db)
	quizAttemptRepo := repository.NewQuizAttemptRepository(db)
	certificateRepo := repository.NewCertificateRepository(db)
	coinTransactionRepo := repository.NewCoinTransactionRepository(db)
	badgeRepo := repository.NewBadgeRepository(db)
	badgeProgressRepo := repository.NewBadgeProgressRepository(db)
	reviewRepo := repository.NewCourseReviewRepository(db)
	auditLogRepo := repository.NewSystemAuditLogRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	courseService := service.NewCourseService(courseRepo, enrollmentRepo, reviewRepo)
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, courseRepo, userProgressRepo, userRepo, coinTransactionRepo, certificateRepo)
	progressService := service.NewProgressService(userProgressRepo, enrollmentRepo, userRepo)
	gamificationService := service.NewGamificationService(coinTransactionRepo, badgeRepo, badgeProgressRepo, userRepo, certificateRepo)
	quizService := service.NewQuizService(quizRepo, quizAttemptRepo, enrollmentRepo, gamificationService)
	dashboardService := service.NewDashboardService(enrollmentRepo, userProgressRepo, certificateRepo, coinTransactionRepo, badgeProgressRepo, userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, cfg)
	courseHandler := handler.NewCourseHandler(courseService, auditLogRepo)
	enrollmentHandler := handler.NewEnrollmentHandler(enrollmentService, auditLogRepo)
	progressHandler := handler.NewProgressHandler(progressService, auditLogRepo)
	quizHandler := handler.NewQuizHandler(quizService, auditLogRepo)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	userHandler := handler.NewUserHandler(userRepo, gamificationService, badgeProgressRepo)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	if cfg.Server.Env == "development" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.RequestIDMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		handler.HealthCheck(c)
	})

	// Public routes (no auth required)
	public := router.Group("/api/v1/public")
	{
		// Auth endpoints
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)

		// Public course endpoints
		public.GET("/courses", courseHandler.GetAllCourses)
		public.GET("/courses/:id", courseHandler.GetCourse)
		public.GET("/courses/search", courseHandler.SearchCourses)
		public.GET("/courses/category/:category", courseHandler.GetByCategory)
	}

	// Protected routes (auth required)
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		// Auth endpoints
		auth := api.Group("/auth")
		{
			auth.GET("/me", authHandler.GetProfile)
			auth.PUT("/profile", authHandler.UpdateProfile)
			auth.POST("/change-password", authHandler.ChangePassword)
			auth.POST("/logout", authHandler.Logout)
		}

		// Dashboard endpoints
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("", dashboardHandler.GetDashboard)
		}

		// Course endpoints
		courses := api.Group("/courses")
		{
			courses.POST("/enroll", enrollmentHandler.Enroll)
			courses.GET("/my-enrollments", enrollmentHandler.GetMyEnrollments)
			courses.GET("/in-progress", enrollmentHandler.GetInProgressCourses)
			courses.GET("/completed", enrollmentHandler.GetCompletedCourses)
			courses.GET("/mandatory", enrollmentHandler.GetMandatoryCourses)
			courses.POST("/:courseId/reviews", courseHandler.AddReview)
			courses.GET("/:courseId/reviews", courseHandler.GetReviews)
		}

		// Progress endpoints
		progress := api.Group("/progress")
		{
			progress.POST("/track", progressHandler.TrackProgress)
			progress.GET("/course/:courseId", progressHandler.GetCourseProgress)
			progress.GET("/lesson/:lessonId", progressHandler.GetLessonProgress)
		}

		// Quiz endpoints
		quiz := api.Group("/quiz")
		{
			quiz.POST("/start", quizHandler.StartAttempt)
			quiz.POST("/submit/:attemptId", quizHandler.SubmitAttempt)
			quiz.GET("/:quizId/attempts", quizHandler.GetAttempts)
		}

		// User endpoints
		user := api.Group("/user")
		{
			user.GET("/profile/:userId", userHandler.GetUserProfile)
			user.GET("/leaderboard", userHandler.GetLeaderboard)
			user.GET("/coins", userHandler.GetCoins)
			user.GET("/coins/transactions", userHandler.GetCoinTransactions)
			user.GET("/badges", userHandler.GetBadges)
			user.GET("/badges/earned", userHandler.GetEarnedBadges)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin", "instructor"))
		{
			// Course management
			admin.POST("/courses", courseHandler.CreateCourse)
			admin.PUT("/courses/:id", courseHandler.UpdateCourse)
			admin.DELETE("/courses/:id", courseHandler.DeleteCourse)
			admin.POST("/courses/:id/publish", courseHandler.PublishCourse)

			// User management
			admin.GET("/users", userHandler.ListUsers)
			admin.GET("/users/:userId", userHandler.GetUserProfile)
			admin.POST("/users/:userId/adjust-coins", userHandler.AdjustCoins)
		}
	}

	// Start server
	address := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting LMS server on %s", address)

	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
