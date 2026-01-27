package main

import (
	"log"
	"os"

	"lms-go-be/config"
	"lms-go-be/handlers"
	"lms-go-be/middleware"
	"lms-go-be/migrations"
	"lms-go-be/repositories"
	"lms-go-be/seeders"
	"lms-go-be/services"

	"github.com/gin-gonic/gin"
)

// main is the entry point of the application
// It initializes configuration, database, and starts the server
func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize database
	db := config.InitDatabase()
	log.Println("Database connected successfully!")

	// Run migrations to create tables
	migrations.AutoMigrate(db)

	// Seed database with initial data (comment out after first run)
	// seeders.ClearDatabase(db) // Uncomment to clear existing data
	seeders.SeedDatabase(db)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	enrollmentRepo := repositories.NewEnrollmentRepository(db)
	lessonRepo := repositories.NewLessonRepository(db)
	lessonProgressRepo := repositories.NewLessonProgressRepository(db)
	coinHistoryRepo := repositories.NewCoinHistoryRepository(db)
	badgeHistoryRepo := repositories.NewBadgeHistoryRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	courseService := services.NewCourseService(courseRepo, enrollmentRepo, lessonRepo)
	enrollmentService := services.NewEnrollmentService(enrollmentRepo, courseRepo, lessonRepo, lessonProgressRepo, userRepo)
	gamificationService := services.NewGamificationService(userRepo, coinHistoryRepo, badgeHistoryRepo, enrollmentRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Setup Gin router
	router := setupRouter(authHandler, courseService, enrollmentService, gamificationService)

	// Get port from config
	appConfig := config.LoadConfig()
	port := appConfig.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRouter configures all routes and middleware
func setupRouter(
	authHandler *handlers.AuthHandler,
	courseService *services.CourseService,
	enrollmentService *services.EnrollmentService,
	gamificationService *services.GamificationService,
) *gin.Engine {
	// Set Gin mode from environment
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlingMiddleware())
	router.Use(middleware.LoggingMiddleware())

	// ===== PUBLIC ROUTES (No Authentication Required) =====
	public := router.Group("/api")
	{
		auth := public.Group("/auth")
		{
			// Authentication endpoints
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Health check endpoint
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "LMS Backend API is running",
			})
		})
	}

	// ===== PROTECTED ROUTES (Authentication Required) =====
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Auth endpoints (protected)
		auth := protected.Group("/auth")
		{
			auth.GET("/profile", authHandler.GetProfile)
			auth.PUT("/profile", authHandler.UpdateProfile)
			auth.POST("/change-password", authHandler.ChangePassword)
			auth.POST("/refresh-token", authHandler.RefreshToken)
			auth.GET("/stats", authHandler.GetStats)
		}

		// Dashboard endpoints
		dashboard := protected.Group("/dashboard")
		{
			dashboard.GET("/mandatory-courses", getDashboardMandatoryCourses(courseService, enrollmentService))
			dashboard.GET("/in-progress-courses", getDashboardInProgressCourses(enrollmentService))
			dashboard.GET("/completed-courses", getDashboardCompletedCourses(enrollmentService))
			dashboard.GET("/coins", getUserCoins(gamificationService))
			dashboard.GET("/badges", getUserBadges(gamificationService))
			dashboard.GET("/stats", getDashboardStats(enrollmentService))
		}

		// Course endpoints
		courses := protected.Group("/courses")
		{
			courses.GET("", getCourses(courseService))
			courses.GET("/mandatory", getMandatoryCourses(courseService))
			courses.GET("/search", searchCourses(courseService))
			courses.GET("/:id", getCourseDetails(courseService))
			courses.POST("/:id/enroll", enrollCourse(enrollmentService))
		}

		// Enrollment endpoints
		enrollments := protected.Group("/enrollments")
		{
			enrollments.GET("", getUserEnrollments(enrollmentService))
			enrollments.GET("/:id", getEnrollmentDetails(enrollmentService))
		}

		// Admin-only routes
		admin := protected.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin", "hc_admin"))
		{
			// Course management
			admin.POST("/courses", createCourse(courseService))
			admin.PUT("/courses/:id", updateCourse(courseService))
			admin.DELETE("/courses/:id", deleteCourse(courseService))
			admin.POST("/courses/:id/publish", publishCourse(courseService))
			admin.GET("/courses/:id/stats", getCourseStats(courseService))

			// Analytics
			admin.GET("/analytics/overview", getAnalyticsOverview())
			admin.GET("/analytics/courses", getCourseAnalytics())
			admin.GET("/analytics/users", getUserAnalytics())
		}

		// Instructor-only routes
		instructor := protected.Group("/instructor")
		instructor.Use(middleware.RoleMiddleware("instructor"))
		{
			instructor.GET("/courses", getInstructorCourses(courseService))
			instructor.POST("/courses", createCourse(courseService))
			instructor.PUT("/courses/:id", updateCourse(courseService))
		}
	}

	return router
}

// Handler functions - Placeholder implementations

func getDashboardMandatoryCourses(courseService *services.CourseService, enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		enrollments, err := enrollmentService.GetMandatoryCoursesForUser(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, enrollments)
	}
}

func getDashboardInProgressCourses(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		enrollments, err := enrollmentService.GetInProgressCourses(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, enrollments)
	}
}

func getDashboardCompletedCourses(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		enrollments, err := enrollmentService.GetCompletedCourses(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, enrollments)
	}
}

func getUserCoins(gamificationService *services.GamificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		coins, err := gamificationService.GetUserCoins(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, coins)
	}
}

func getUserBadges(gamificationService *services.GamificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		badges, err := gamificationService.GetUserBadges(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, badges)
	}
}

func getDashboardStats(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		stats, err := enrollmentService.GetUserStats(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, stats)
	}
}

func getCourses(courseService *services.CourseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := 1
		pageSize := 10
		courses, total, err := courseService.GetAllCourses(page, pageSize, map[string]interface{}{})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": courses, "total": total})
	}
}

func getMandatoryCourses(courseService *services.CourseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		courses, err := courseService.GetMandatoryCourses()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, courses)
	}
}

func searchCourses(courseService *services.CourseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Query("q")
		courses, total, err := courseService.SearchCourses(keyword, 1, 10)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": courses, "total": total})
	}
}

func getCourseDetails(courseService *services.CourseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		course, err := courseService.GetCourseByID(courseID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, course)
	}
}

func enrollCourse(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		courseID := c.Param("id")
		enrollment, err := enrollmentService.EnrollUser(userID, courseID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, enrollment)
	}
}

func getUserEnrollments(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		enrollments, err := enrollmentService.GetUserCourses(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, enrollments)
	}
}

func getEnrollmentDetails(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Get enrollment details"})
	}
}

func createCourse(courseService *services.CourseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(201, gin.H{"message": "Course created"})
	}
}

func updateCourse(courseService *services.CourseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Course updated"})
	}
}

func deleteCourse(courseService *services.CourseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Course deleted"})
	}
}

func publishCourse(courseService *services.CourseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Course published"})
	}
}

func getCourseStats(courseService *services.CourseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		stats, err := courseService.GetCourseStats(courseID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, stats)
	}
}

func getAnalyticsOverview() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Analytics overview"})
	}
}

func getCourseAnalytics() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Course analytics"})
	}
}

func getUserAnalytics() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "User analytics"})
	}
}

func getInstructorCourses(courseService *services.CourseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Instructor courses"})
	}
}
