package database

import (
	"fmt"
	"log"
	"time"

	"lms-go-be/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed seeds the database with initial data
func Seed(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Seed users
	if err := seedUsers(db); err != nil {
		return fmt.Errorf("error seeding users: %v", err)
	}

	// Seed courses
	if err := seedCourses(db); err != nil {
		return fmt.Errorf("error seeding courses: %v", err)
	}

	// Seed lessons
	if err := seedLessons(db); err != nil {
		return fmt.Errorf("error seeding lessons: %v", err)
	}

	// Seed quizzes
	if err := seedQuizzes(db); err != nil {
		return fmt.Errorf("error seeding quizzes: %v", err)
	}

	// Seed enrollments
	if err := seedEnrollments(db); err != nil {
		return fmt.Errorf("error seeding enrollments: %v", err)
	}

	// Seed badges
	if err := seedBadges(db); err != nil {
		return fmt.Errorf("error seeding badges: %v", err)
	}

	log.Println("Database seeding completed successfully")
	return nil
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}
	return string(hashedPassword)
}

// seedUsers seeds initial users
func seedUsers(db *gorm.DB) error {
	users := []models.User{
		{
			Email:      "admin@lms.com",
			Password:   hashPassword("admin@123"),
			FirstName:  "Admin",
			LastName:   "User",
			Department: "Management",
			Role:       "admin",
			IsActive:   true,
			GMFCCoins:  5000,
		},
		{
			Email:      "instructor@lms.com",
			Password:   hashPassword("instructor@123"),
			FirstName:  "John",
			LastName:   "Instructor",
			Department: "Training",
			Role:       "instructor",
			IsActive:   true,
			GMFCCoins:  2000,
		},
		{
			Email:      "hr@lms.com",
			Password:   hashPassword("hr@123"),
			FirstName:  "Sarah",
			LastName:   "HR",
			Department: "Human Resources",
			Role:       "hr_personnel",
			IsActive:   true,
			GMFCCoins:  1000,
		},
		{
			Email:      "learner1@lms.com",
			Password:   hashPassword("learner@123"),
			FirstName:  "Alice",
			LastName:   "Johnson",
			Department: "Engineering",
			Role:       "learner",
			IsActive:   true,
			GMFCCoins:  500,
		},
		{
			Email:      "learner2@lms.com",
			Password:   hashPassword("learner@123"),
			FirstName:  "Bob",
			LastName:   "Smith",
			Department: "Sales",
			Role:       "learner",
			IsActive:   true,
			GMFCCoins:  300,
		},
		{
			Email:      "learner3@lms.com",
			Password:   hashPassword("learner@123"),
			FirstName:  "Carol",
			LastName:   "Williams",
			Department: "Engineering",
			Role:       "learner",
			IsActive:   true,
			GMFCCoins:  750,
		},
	}

	// Check if users already exist
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Users already seeded, skipping...")
		return nil
	}

	if err := db.Create(&users).Error; err != nil {
		return err
	}

	log.Printf("Seeded %d users\n", len(users))
	return nil
}

// seedCourses seeds initial courses
func seedCourses(db *gorm.DB) error {
	var instructorID uint
	if err := db.Model(&models.User{}).Where("role = ?", "instructor").Select("id").Row().Scan(&instructorID); err != nil {
		return fmt.Errorf("error fetching instructor: %v", err)
	}

	dueDate := time.Now().AddDate(0, 0, 30)

	courses := []models.Course{
		{
			Title:            "Go Programming Fundamentals",
			Description:      "Learn the basics of Go programming language including variables, functions, and packages",
			Category:         "Programming",
			InstructorID:     instructorID,
			DurationMinutes:  120,
			DifficultyLevel:  "beginner",
			PassingScore:     70,
			IsMandatory:      true,
			MandatoryDueDate: &dueDate,
			MaxEnrollments:   100,
			IsPublished:      true,
			CoinsReward:      150,
			AverageRating:    4.5,
		},
		{
			Title:           "Advanced Gin Framework",
			Description:     "Master the Gin web framework for building high-performance REST APIs",
			Category:        "Web Development",
			InstructorID:    instructorID,
			DurationMinutes: 150,
			DifficultyLevel: "intermediate",
			PassingScore:    75,
			IsMandatory:     false,
			MaxEnrollments:  50,
			IsPublished:     true,
			CoinsReward:     200,
			AverageRating:   4.7,
		},
		{
			Title:            "Database Design with PostgreSQL",
			Description:      "Learn to design and optimize databases using PostgreSQL and GORM",
			Category:         "Databases",
			InstructorID:     instructorID,
			DurationMinutes:  180,
			DifficultyLevel:  "intermediate",
			PassingScore:     70,
			IsMandatory:      true,
			MandatoryDueDate: &dueDate,
			MaxEnrollments:   75,
			IsPublished:      true,
			CoinsReward:      175,
			AverageRating:    4.6,
		},
		{
			Title:           "RESTful API Best Practices",
			Description:     "Master best practices for designing and building RESTful APIs",
			Category:        "API Design",
			InstructorID:    instructorID,
			DurationMinutes: 100,
			DifficultyLevel: "intermediate",
			PassingScore:    70,
			IsMandatory:     false,
			MaxEnrollments:  80,
			IsPublished:     true,
			CoinsReward:     125,
			AverageRating:   4.4,
		},
		{
			Title:           "Testing in Go",
			Description:     "Write effective unit tests and integration tests in Go",
			Category:        "Testing",
			InstructorID:    instructorID,
			DurationMinutes: 90,
			DifficultyLevel: "intermediate",
			PassingScore:    65,
			IsMandatory:     false,
			MaxEnrollments:  60,
			IsPublished:     true,
			CoinsReward:     100,
			AverageRating:   4.3,
		},
	}

	var count int64
	db.Model(&models.Course{}).Count(&count)
	if count > 0 {
		log.Println("Courses already seeded, skipping...")
		return nil
	}

	if err := db.Create(&courses).Error; err != nil {
		return err
	}

	log.Printf("Seeded %d courses\n", len(courses))
	return nil
}

// seedLessons seeds initial lessons
func seedLessons(db *gorm.DB) error {
	var courses []models.Course
	if err := db.Find(&courses).Error; err != nil {
		return err
	}

	if len(courses) == 0 {
		return fmt.Errorf("no courses found")
	}

	var lessons []models.Lesson
	for _, course := range courses {
		lessons = append(lessons,
			models.Lesson{
				CourseID:      course.ID,
				Title:         fmt.Sprintf("%s - Lesson 1", course.Title),
				Description:   "Introduction to the course",
				ContentType:   "video",
				VideoDuration: 30,
				OrderNumber:   1,
				IsPublished:   true,
			},
			models.Lesson{
				CourseID:      course.ID,
				Title:         fmt.Sprintf("%s - Lesson 2", course.Title),
				Description:   "Core concepts",
				ContentType:   "video",
				VideoDuration: 45,
				OrderNumber:   2,
				IsPublished:   true,
			},
			models.Lesson{
				CourseID:      course.ID,
				Title:         fmt.Sprintf("%s - Lesson 3", course.Title),
				Description:   "Practical examples",
				ContentType:   "video",
				VideoDuration: 40,
				OrderNumber:   3,
				IsPublished:   true,
			},
		)
	}

	var count int64
	db.Model(&models.Lesson{}).Count(&count)
	if count > 0 {
		log.Println("Lessons already seeded, skipping...")
		return nil
	}

	if err := db.Create(&lessons).Error; err != nil {
		return err
	}

	log.Printf("Seeded %d lessons\n", len(lessons))
	return nil
}

// seedQuizzes seeds initial quizzes
func seedQuizzes(db *gorm.DB) error {
	var courses []models.Course
	if err := db.Find(&courses).Error; err != nil {
		return err
	}

	if len(courses) == 0 {
		return fmt.Errorf("no courses found")
	}

	quizzes := []models.Quiz{
		{
			CourseID:      courses[0].ID,
			Title:         fmt.Sprintf("%s - Final Quiz", courses[0].Title),
			Description:   "Test your knowledge",
			PassingScore:  70,
			TimeLimit:     60,
			Attempts:      3,
			QuestionCount: 10,
			IsPublished:   true,
		},
	}

	if len(courses) > 1 {
		quizzes = append(quizzes, models.Quiz{
			CourseID:      courses[1].ID,
			Title:         fmt.Sprintf("%s - Final Quiz", courses[1].Title),
			Description:   "Test your knowledge",
			PassingScore:  75,
			TimeLimit:     75,
			Attempts:      3,
			QuestionCount: 15,
			IsPublished:   true,
		})
	}

	var count int64
	db.Model(&models.Quiz{}).Count(&count)
	if count > 0 {
		log.Println("Quizzes already seeded, skipping...")
		return nil
	}

	if err := db.Create(&quizzes).Error; err != nil {
		return err
	}

	log.Printf("Seeded %d quizzes\n", len(quizzes))

	// Seed questions for the first quiz
	if err := seedQuestions(db, quizzes[0].ID); err != nil {
		return err
	}

	return nil
}

// seedQuestions seeds questions for a quiz
func seedQuestions(db *gorm.DB, quizID uint) error {
	questions := []models.Question{
		{
			QuizID:       quizID,
			QuestionText: "What is Go's main purpose?",
			QuestionType: "mcq",
			OrderNumber:  1,
			Points:       1,
			IsPublished:  true,
		},
		{
			QuizID:       quizID,
			QuestionText: "Is Go a compiled language?",
			QuestionType: "true_false",
			OrderNumber:  2,
			Points:       1,
			IsPublished:  true,
		},
	}

	if err := db.Create(&questions).Error; err != nil {
		return err
	}

	// Seed options for the first question
	if len(questions) > 0 {
		options := []models.QuestionOption{
			{
				QuestionID:  questions[0].ID,
				OptionText:  "Building web applications",
				IsCorrect:   false,
				OrderNumber: 1,
			},
			{
				QuestionID:  questions[0].ID,
				OptionText:  "Building systems software and concurrent applications",
				IsCorrect:   true,
				OrderNumber: 2,
			},
			{
				QuestionID:  questions[0].ID,
				OptionText:  "Mobile app development",
				IsCorrect:   false,
				OrderNumber: 3,
			},
		}

		if err := db.Create(&options).Error; err != nil {
			return err
		}
	}

	// Seed options for the second question (true/false)
	if len(questions) > 1 {
		options := []models.QuestionOption{
			{
				QuestionID:  questions[1].ID,
				OptionText:  "True",
				IsCorrect:   true,
				OrderNumber: 1,
			},
			{
				QuestionID:  questions[1].ID,
				OptionText:  "False",
				IsCorrect:   false,
				OrderNumber: 2,
			},
		}

		if err := db.Create(&options).Error; err != nil {
			return err
		}
	}

	return nil
}

// seedEnrollments seeds initial enrollments
func seedEnrollments(db *gorm.DB) error {
	var users []models.User
	var courses []models.Course

	db.Where("role = ?", "learner").Find(&users)
	db.Find(&courses)

	if len(users) == 0 || len(courses) == 0 {
		return fmt.Errorf("not enough users or courses for enrollment")
	}

	var enrollments []models.Enrollment
	for i, user := range users {
		for j := 0; j < len(courses) && j < 3; j++ {
			status := "not_started"
			progress := 0
			var completedAt *time.Time

			if i%2 == 0 {
				status = "in_progress"
				progress = 50
			}

			if i%3 == 0 && j == 0 {
				status = "completed"
				progress = 100
				now := time.Now()
				completedAt = &now
			}

			enrollments = append(enrollments, models.Enrollment{
				UserID:           user.ID,
				CourseID:         courses[j].ID,
				CompletionStatus: status,
				OverallProgress:  progress,
				CompletedAt:      completedAt,
				IsPassed:         status == "completed",
			})
		}
	}

	var count int64
	db.Model(&models.Enrollment{}).Count(&count)
	if count > 0 {
		log.Println("Enrollments already seeded, skipping...")
		return nil
	}

	if err := db.Create(&enrollments).Error; err != nil {
		return err
	}

	log.Printf("Seeded %d enrollments\n", len(enrollments))
	return nil
}

// seedBadges seeds initial badges
func seedBadges(db *gorm.DB) error {
	badges := []models.Badge{
		{
			Name:        "Bronze Learner",
			Description: "Complete your first course",
			Level:       "bronze",
			Criteria:    `{"type": "courses_completed", "value": 1}`,
		},
		{
			Name:        "Silver Learner",
			Description: "Complete 5 courses",
			Level:       "silver",
			Criteria:    `{"type": "courses_completed", "value": 5}`,
		},
		{
			Name:        "Gold Learner",
			Description: "Complete 10 courses",
			Level:       "gold",
			Criteria:    `{"type": "courses_completed", "value": 10}`,
		},
		{
			Name:        "Platinum Expert",
			Description: "Complete 20 courses and maintain 90% average",
			Level:       "platinum",
			Criteria:    `{"type": "courses_completed", "value": 20, "avg_score": 90}`,
		},
		{
			Name:        "Quiz Master",
			Description: "Pass 5 quizzes with perfect score",
			Level:       "gold",
			Criteria:    `{"type": "perfect_quizzes", "value": 5}`,
		},
		{
			Name:        "Consistent Learner",
			Description: "Maintain 30 day learning streak",
			Level:       "silver",
			Criteria:    `{"type": "learning_streak", "days": 30}`,
		},
	}

	var count int64
	db.Model(&models.Badge{}).Count(&count)
	if count > 0 {
		log.Println("Badges already seeded, skipping...")
		return nil
	}

	if err := db.Create(&badges).Error; err != nil {
		return err
	}

	log.Printf("Seeded %d badges\n", len(badges))
	return nil
}

// CleanDatabase clears all data from the database (use with caution!)
func CleanDatabase(db *gorm.DB) error {
	log.Println("WARNING: Cleaning database...")
	tables := []string{
		"system_audit_logs",
		"download_logs",
		"learning_reports",
		"course_reviews",
		"badge_progresses",
		"badges",
		"coin_transactions",
		"certificates",
		"quiz_answer_entries",
		"quiz_attempts",
		"question_answers",
		"question_options",
		"questions",
		"quizzes",
		"user_progresses",
		"enrollments",
		"lesson_materials",
		"lessons",
		"courses",
		"users",
	}

	for _, table := range tables {
		if err := db.Exec("DELETE FROM " + table).Error; err != nil {
			log.Printf("Warning: failed to clean table %s: %v", table, err)
		}
	}

	log.Println("Database cleaned")
	return nil
}
