package seeders

import (
	"log"
	"time"

	"lms-go-be/models"
	"lms-go-be/utils"

	"gorm.io/gorm"
)

// SeedDatabase populates database with initial data
// This function creates test users, categories, courses, and lessons
func SeedDatabase(db *gorm.DB) {
	log.Println("Seeding database with initial data...")

	// Seed users
	seedUsers(db)

	// Seed categories
	seedCategories(db)

	// Seed courses
	seedCourses(db)

	// Seed lessons
	seedLessons(db)

	// Seed quizzes
	seedQuizzes(db)

	log.Println("Database seeding completed successfully!")
}

// seedUsers creates sample users
func seedUsers(db *gorm.DB) {
	users := []models.User{
		{
			FirstName:     "John",
			LastName:      "Doe",
			Email:         "john@example.com",
			Password:      hashPassword("password123"),
			Role:          "learner",
			Status:        "active",
			Department:    "Engineering",
			GMFCCoins:     0,
			BadgeLevel:    "Bronze",
			CurrentStreak: 0,
			IsVerified:    true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			FirstName:     "Jane",
			LastName:      "Smith",
			Email:         "jane@example.com",
			Password:      hashPassword("password123"),
			Role:          "instructor",
			Status:        "active",
			Department:    "Training",
			GMFCCoins:     0,
			BadgeLevel:    "Gold",
			CurrentStreak: 0,
			IsVerified:    true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			FirstName:     "Admin",
			LastName:      "User",
			Email:         "admin@example.com",
			Password:      hashPassword("admin123"),
			Role:          "admin",
			Status:        "active",
			GMFCCoins:     0,
			BadgeLevel:    "Diamond",
			CurrentStreak: 0,
			IsVerified:    true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			FirstName:     "Alice",
			LastName:      "Johnson",
			Email:         "alice@example.com",
			Password:      hashPassword("password123"),
			Role:          "learner",
			Status:        "active",
			Department:    "Sales",
			GMFCCoins:     150,
			BadgeLevel:    "Silver",
			CurrentStreak: 5,
			IsVerified:    true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", user.Email, err)
		}
	}

	log.Println("Users seeded successfully")
}

// seedCategories creates course categories
func seedCategories(db *gorm.DB) {
	categories := []models.Category{
		{
			Name: "Development",
			Slug: "development",
			Icon: "code",
		},
		{
			Name: "Business",
			Slug: "business",
			Icon: "briefcase",
		},
		{
			Name: "Design",
			Slug: "design",
			Icon: "palette",
		},
		{
			Name: "Leadership",
			Slug: "leadership",
			Icon: "users",
		},
		{
			Name: "Compliance",
			Slug: "compliance",
			Icon: "shield",
		},
	}

	for _, category := range categories {
		if err := db.Create(&category).Error; err != nil {
			log.Printf("Failed to seed category %s: %v", category.Name, err)
		}
	}

	log.Println("Categories seeded successfully")
}

// seedCourses creates sample courses
func seedCourses(db *gorm.DB) {
	// Get sample instructor
	var instructor models.User
	db.Where("email = ?", "jane@example.com").First(&instructor)

	// Get categories
	var devCategory, businessCategory models.Category
	db.Where("slug = ?", "development").First(&devCategory)
	db.Where("slug = ?", "business").First(&businessCategory)

	courses := []models.Course{
		{
			Title:         "Golang for Beginners",
			Description:   "Learn Go programming language from scratch. Master the fundamentals and build real-world applications.",
			Level:         "beginner",
			Status:        "published",
			IsMandatory:   true,
			PassingScore:  70,
			InstructorID:  instructor.ID,
			CategoryID:    devCategory.ID,
			TotalDuration: 360,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Title:         "Advanced Go Development",
			Description:   "Master advanced Go concepts including concurrency, channels, and system programming.",
			Level:         "advanced",
			Status:        "published",
			IsMandatory:   false,
			PassingScore:  75,
			InstructorID:  instructor.ID,
			CategoryID:    devCategory.ID,
			TotalDuration: 480,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Title:         "Business Ethics and Compliance",
			Description:   "Essential training on corporate ethics, compliance requirements, and professional conduct.",
			Level:         "beginner",
			Status:        "published",
			IsMandatory:   true,
			PassingScore:  80,
			InstructorID:  instructor.ID,
			CategoryID:    businessCategory.ID,
			TotalDuration: 120,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Title:         "Project Management Fundamentals",
			Description:   "Learn essential project management skills and methodologies.",
			Level:         "intermediate",
			Status:        "published",
			IsMandatory:   false,
			PassingScore:  70,
			InstructorID:  instructor.ID,
			CategoryID:    businessCategory.ID,
			TotalDuration: 240,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	for _, course := range courses {
		if err := db.Create(&course).Error; err != nil {
			log.Printf("Failed to seed course %s: %v", course.Title, err)
		}
	}

	log.Println("Courses seeded successfully")
}

// seedLessons creates lessons for courses
func seedLessons(db *gorm.DB) {
	var courses []models.Course
	db.Find(&courses)

	for i, course := range courses {
		lessons := []models.Lesson{
			{
				CourseID:    course.ID,
				Title:       "Introduction to " + course.Title,
				Description: "Get started with this course and understand the basics.",
				Duration:    45,
				OrderIndex:  1,
				IsLocked:    false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				CourseID:    course.ID,
				Title:       "Core Concepts",
				Description: "Learn the fundamental concepts and principles.",
				Duration:    60,
				OrderIndex:  2,
				IsLocked:    false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				CourseID:    course.ID,
				Title:       "Advanced Topics",
				Description: "Dive deeper into advanced topics and real-world applications.",
				Duration:    60,
				OrderIndex:  3,
				IsLocked:    false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		for _, lesson := range lessons {
			if err := db.Create(&lesson).Error; err != nil {
				log.Printf("Failed to seed lesson %s: %v", lesson.Title, err)
			}
		}
	}

	log.Println("Lessons seeded successfully")
}

// seedQuizzes creates quizzes for courses
func seedQuizzes(db *gorm.DB) {
	var courses []models.Course
	db.Find(&courses)

	for _, course := range courses {
		quiz := models.Quiz{
			CourseID:       course.ID,
			Title:          course.Title + " Final Assessment",
			Description:    "Test your knowledge of " + course.Title,
			TimeLimit:      60,
			PassingScore:   70,
			MaxAttempts:    3,
			IsPublished:    true,
			RandomizeOrder: true,
			ShowAnswers:    true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := db.Create(&quiz).Error; err != nil {
			log.Printf("Failed to seed quiz for course %s: %v", course.Title, err)
		}
	}

	log.Println("Quizzes seeded successfully")
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) string {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return ""
	}
	return hashed
}

// ClearDatabase deletes all data from tables
// Use with caution - this is for development/testing only
func ClearDatabase(db *gorm.DB) {
	log.Println("Clearing database...")

	db.Exec("DELETE FROM notifications")
	db.Exec("DELETE FROM badge_history")
	db.Exec("DELETE FROM coin_history")
	db.Exec("DELETE FROM certificates")
	db.Exec("DELETE FROM quiz_attempts")
	db.Exec("DELETE FROM question_options")
	db.Exec("DELETE FROM questions")
	db.Exec("DELETE FROM quizzes")
	db.Exec("DELETE FROM course_materials")
	db.Exec("DELETE FROM lesson_progress")
	db.Exec("DELETE FROM lessons")
	db.Exec("DELETE FROM enrollments")
	db.Exec("DELETE FROM courses")
	db.Exec("DELETE FROM categories")
	db.Exec("DELETE FROM users")

	log.Println("Database cleared successfully!")
}
