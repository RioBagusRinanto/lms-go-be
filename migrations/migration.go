package migrations

import (
	"log"

	"lms-go-be/models"

	"gorm.io/gorm"
)

// AutoMigrate runs all migrations
// This function creates all database tables based on models
func AutoMigrate(db *gorm.DB) {
	log.Println("Running database migrations...")

	// Create tables in proper order (respecting foreign key relationships)
	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Course{},
		&models.Lesson{},
		&models.LessonProgress{},
		&models.Enrollment{},
		&models.Quiz{},
		&models.Question{},
		&models.QuestionOption{},
		&models.QuizAttempt{},
		&models.CourseMaterial{},
		&models.Certificate{},
		&models.CoinHistory{},
		&models.BadgeHistory{},
		&models.Notification{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Database migrations completed successfully!")
}
