package database

import (
	"fmt"
	"log"

	"lms-go-be/internal/config"
	"lms-go-be/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the database connection and runs auto migrations
func InitDB(cfg *config.Config) *gorm.DB {
	dsn := cfg.Database.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	// Run auto migrations
	AutoMigrate(db)

	return db
}

// AutoMigrate runs all database migrations
func AutoMigrate(db *gorm.DB) {
	models := []interface{}{
		&models.User{},
		&models.Course{},
		&models.Lesson{},
		&models.LessonMaterial{},
		&models.Quiz{},
		&models.Question{},
		&models.QuestionOption{},
		&models.QuestionAnswer{},
		&models.Enrollment{},
		&models.UserProgress{},
		&models.QuizAttempt{},
		&models.QuizAnswerEntry{},
		&models.Certificate{},
		&models.CoinTransaction{},
		&models.Badge{},
		&models.BadgeProgress{},
		&models.CourseReview{},
		&models.LearningReport{},
		&models.DownloadLog{},
		&models.SystemAuditLog{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			log.Fatalf("Failed to migrate model: %v", err)
		}
	}

	// Create indexes
	CreateIndexes(db)

	log.Println("Database migrations completed successfully")
}

// CreateIndexes creates additional database indexes for performance
func CreateIndexes(db *gorm.DB) {
	indexes := map[string]string{
		"idx_user_email":                "CREATE INDEX IF NOT EXISTS idx_user_email ON users(email);",
		"idx_course_instructor":         "CREATE INDEX IF NOT EXISTS idx_course_instructor ON courses(instructor_id);",
		"idx_enrollment_user_course":    "CREATE INDEX IF NOT EXISTS idx_enrollment_user_course ON enrollments(user_id, course_id);",
		"idx_lesson_course":             "CREATE INDEX IF NOT EXISTS idx_lesson_course ON lessons(course_id);",
		"idx_user_progress_user_course": "CREATE INDEX IF NOT EXISTS idx_user_progress_user_course ON user_progresses(user_id, course_id);",
		"idx_quiz_attempt_user":         "CREATE INDEX IF NOT EXISTS idx_quiz_attempt_user ON quiz_attempts(user_id);",
		"idx_certificate_user_course":   "CREATE INDEX IF NOT EXISTS idx_certificate_user_course ON certificates(user_id, course_id);",
		"idx_coin_transaction_user":     "CREATE INDEX IF NOT EXISTS idx_coin_transaction_user ON coin_transactions(user_id);",
		"idx_badge_progress_user":       "CREATE INDEX IF NOT EXISTS idx_badge_progress_user ON badge_progresses(user_id);",
		"idx_course_review_user":        "CREATE INDEX IF NOT EXISTS idx_course_review_user ON course_reviews(user_id);",
		"idx_system_audit_log_user":     "CREATE INDEX IF NOT EXISTS idx_system_audit_log_user ON system_audit_logs(user_id);",
		"idx_system_audit_log_action":   "CREATE INDEX IF NOT EXISTS idx_system_audit_log_action ON system_audit_logs(action);",
	}

	for name, query := range indexes {
		if err := db.Exec(query).Error; err != nil {
			log.Printf("Warning: Failed to create index %s: %v", name, err)
		}
	}
}

// CheckConnection checks if the database is reachable
func CheckConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	return nil
}

// Close closes the database connection
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
