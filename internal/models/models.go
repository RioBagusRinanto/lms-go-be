package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a system user (learner, instructor, admin, etc.)
type User struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Email              string         `gorm:"uniqueIndex;not null" json:"email"`
	Password           string         `gorm:"not null" json:"-"`
	FirstName          string         `gorm:"not null" json:"first_name"`
	LastName           string         `gorm:"not null" json:"last_name"`
	Department         string         `json:"department"`
	Role               string         `gorm:"not null;default:'learner'" json:"role"` // learner, instructor, admin, hr_personnel
	IsActive           bool           `gorm:"not null;default:true" json:"is_active"`
	LastLoginAt        *time.Time     `json:"last_login_at"`
	ProfileImageURL    string         `json:"profile_image_url"`
	GMFCCoins          int64          `gorm:"not null;default:0" json:"gmfc_coins"`
	CurrentBadgeLevel  string         `gorm:"default:'bronze'" json:"current_badge_level"` // bronze, silver, gold, platinum
	TotalLearningHours float64        `gorm:"default:0" json:"total_learning_hours"`
	CurrentStreak      int            `gorm:"default:0" json:"current_streak"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Enrollments      []Enrollment      `gorm:"foreignKey:UserID"`
	UserProgress     []UserProgress    `gorm:"foreignKey:UserID"`
	QuizAttempts     []QuizAttempt     `gorm:"foreignKey:UserID"`
	Certificates     []Certificate     `gorm:"foreignKey:UserID"`
	CoinTransactions []CoinTransaction `gorm:"foreignKey:UserID"`
	BadgeProgresses  []BadgeProgress   `gorm:"foreignKey:UserID"`
}

// Course represents a training course
type Course struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Title            string         `gorm:"not null;index" json:"title"`
	Description      string         `gorm:"type:text" json:"description"`
	Category         string         `gorm:"not null;index" json:"category"`
	InstructorID     uint           `gorm:"not null" json:"instructor_id"`
	ThumbnailURL     string         `json:"thumbnail_url"`
	DurationMinutes  int            `gorm:"not null" json:"duration_minutes"`
	DifficultyLevel  string         `gorm:"default:'beginner'" json:"difficulty_level"` // beginner, intermediate, advanced
	PassingScore     int            `gorm:"default:70" json:"passing_score"`
	IsMandatory      bool           `gorm:"default:false;index" json:"is_mandatory"`
	MandatoryDueDate *time.Time     `json:"mandatory_due_date"`
	MaxEnrollments   int            `json:"max_enrollments"`
	IsPublished      bool           `gorm:"default:false;index" json:"is_published"`
	EnrollmentCount  int            `gorm:"default:0" json:"enrollment_count"`
	CompletionCount  int            `gorm:"default:0" json:"completion_count"`
	AverageRating    float64        `gorm:"default:0" json:"average_rating"`
	CoinsReward      int            `gorm:"default:100" json:"coins_reward"` // Coins earned on completion
	BadgeReward      string         `json:"badge_reward"`                    // Badge earned on completion
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Instructor   User           `gorm:"foreignKey:InstructorID"`
	Lessons      []Lesson       `gorm:"foreignKey:CourseID"`
	Quizzes      []Quiz         `gorm:"foreignKey:CourseID"`
	Enrollments  []Enrollment   `gorm:"foreignKey:CourseID"`
	Certificates []Certificate  `gorm:"foreignKey:CourseID"`
	Reviews      []CourseReview `gorm:"foreignKey:CourseID"`
}

// Lesson represents a single lesson within a course
type Lesson struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CourseID      uint           `gorm:"not null;index" json:"course_id"`
	Title         string         `gorm:"not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	ContentType   string         `gorm:"not null" json:"content_type"` // video, document, interactive
	VideoURL      *string        `json:"video_url"`
	VideoDuration int            `json:"video_duration_minutes"`
	OrderNumber   int            `gorm:"not null" json:"order_number"`
	IsPublished   bool           `gorm:"default:true" json:"is_published"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Course          Course           `gorm:"foreignKey:CourseID"`
	Materials       []LessonMaterial `gorm:"foreignKey:LessonID"`
	UserProgressLog []UserProgress   `gorm:"foreignKey:LessonID"`
	Quiz            *Quiz            `gorm:"foreignKey:LessonID"`
}

// LessonMaterial represents downloadable materials for a lesson
type LessonMaterial struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	LessonID      uint           `gorm:"not null;index" json:"lesson_id"`
	MaterialName  string         `gorm:"not null" json:"material_name"`
	MaterialType  string         `gorm:"not null" json:"material_type"` // pdf, ppt, pptx, doc, zip
	FileURL       string         `gorm:"not null" json:"file_url"`
	FileSizeBytes int64          `json:"file_size_bytes"`
	Version       int            `gorm:"default:1" json:"version"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Lesson Lesson `gorm:"foreignKey:LessonID"`
}

// Quiz represents a quiz/assessment for a lesson
type Quiz struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CourseID      uint           `gorm:"not null;index" json:"course_id"`
	LessonID      *uint          `gorm:"index" json:"lesson_id"`
	Title         string         `gorm:"not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	PassingScore  int            `gorm:"default:70" json:"passing_score"`
	TimeLimit     int            `json:"time_limit_minutes"` // 0 means no limit
	Attempts      int            `gorm:"default:3" json:"allowed_attempts"`
	QuestionCount int            `gorm:"default:0" json:"question_count"`
	IsPublished   bool           `gorm:"default:false" json:"is_published"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Course    Course        `gorm:"foreignKey:CourseID"`
	Lesson    *Lesson       `gorm:"foreignKey:LessonID"`
	Questions []Question    `gorm:"foreignKey:QuizID"`
	Attempts  []QuizAttempt `gorm:"foreignKey:QuizID"`
}

// Question represents a single question in a quiz
type Question struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	QuizID       uint           `gorm:"not null;index" json:"quiz_id"`
	QuestionText string         `gorm:"type:text;not null" json:"question_text"`
	QuestionType string         `gorm:"not null" json:"question_type"` // mcq, true_false, short_answer, fill_blank
	OrderNumber  int            `gorm:"not null" json:"order_number"`
	Points       int            `gorm:"default:1" json:"points"`
	IsPublished  bool           `gorm:"default:true" json:"is_published"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Quiz    Quiz             `gorm:"foreignKey:QuizID"`
	Options []QuestionOption `gorm:"foreignKey:QuestionID"`
	Answers []QuestionAnswer `gorm:"foreignKey:QuestionID"`
}

// QuestionOption represents an option for MCQ or True/False questions
type QuestionOption struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	QuestionID  uint           `gorm:"not null;index" json:"question_id"`
	OptionText  string         `gorm:"type:text;not null" json:"option_text"`
	IsCorrect   bool           `gorm:"default:false" json:"is_correct"`
	OrderNumber int            `gorm:"not null" json:"order_number"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Question Question `gorm:"foreignKey:QuestionID"`
}

// QuestionAnswer represents instructor's answer to short answer/fill blank questions
type QuestionAnswer struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	QuestionID  uint           `gorm:"not null;index" json:"question_id"`
	CorrectText string         `gorm:"type:text;not null" json:"correct_text"`
	IsPartialOK bool           `gorm:"default:false" json:"is_partial_ok"` // For fill blanks
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Question Question `gorm:"foreignKey:QuestionID"`
}

// Enrollment represents a user's enrollment in a course
type Enrollment struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UserID           uint           `gorm:"not null;index" json:"user_id"`
	CourseID         uint           `gorm:"not null;index" json:"course_id"`
	EnrolledAt       time.Time      `gorm:"autoCreateTime" json:"enrolled_at"`
	CompletedAt      *time.Time     `json:"completed_at"`
	LastAccessedAt   *time.Time     `json:"last_accessed_at"`
	CompletionStatus string         `gorm:"default:'not_started'" json:"completion_status"` // not_started, in_progress, completed
	OverallProgress  int            `gorm:"default:0" json:"overall_progress"`              // 0-100
	FinalScore       int            `json:"final_score"`
	IsPassed         bool           `gorm:"default:false" json:"is_passed"`
	IsOverdue        bool           `gorm:"default:false;index" json:"is_overdue"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User   User   `gorm:"foreignKey:UserID"`
	Course Course `gorm:"foreignKey:CourseID"`
}

// UserProgress tracks user's progress on individual lessons
type UserProgress struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	UserID             uint           `gorm:"not null;index" json:"user_id"`
	CourseID           uint           `gorm:"not null;index" json:"course_id"`
	LessonID           uint           `gorm:"not null;index" json:"lesson_id"`
	WatchedDuration    int            `json:"watched_duration_seconds"`
	TotalDuration      int            `json:"total_duration_seconds"`
	IsCompleted        bool           `gorm:"default:false;index" json:"is_completed"`
	CompletedAt        *time.Time     `json:"completed_at"`
	LastAccessedAt     *time.Time     `json:"last_accessed_at"`
	ProgressPercentage int            `gorm:"default:0" json:"progress_percentage"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User   User   `gorm:"foreignKey:UserID"`
	Course Course `gorm:"foreignKey:CourseID"`
	Lesson Lesson `gorm:"foreignKey:LessonID"`
}

// QuizAttempt represents a user's attempt at a quiz
type QuizAttempt struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	QuizID           uint           `gorm:"not null;index" json:"quiz_id"`
	UserID           uint           `gorm:"not null;index" json:"user_id"`
	AttemptNumber    int            `gorm:"not null" json:"attempt_number"`
	StartedAt        time.Time      `gorm:"autoCreateTime" json:"started_at"`
	SubmittedAt      *time.Time     `json:"submitted_at"`
	Score            int            `json:"score"`
	MaxScore         int            `json:"max_score"`
	Percentage       int            `json:"percentage"`
	IsPassed         bool           `json:"is_passed"`
	TimeSpentSeconds int            `json:"time_spent_seconds"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Quiz    Quiz              `gorm:"foreignKey:QuizID"`
	User    User              `gorm:"foreignKey:UserID"`
	Answers []QuizAnswerEntry `gorm:"foreignKey:QuizAttemptID"`
}

// QuizAnswerEntry represents a user's answer to a specific question in a quiz attempt
type QuizAnswerEntry struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	QuizAttemptID uint           `gorm:"not null;index" json:"quiz_attempt_id"`
	QuestionID    uint           `gorm:"not null" json:"question_id"`
	UserAnswer    string         `gorm:"type:text" json:"user_answer"` // Can be option ID or text
	IsCorrect     *bool          `json:"is_correct"`                   // nil for pending grading
	PointsEarned  int            `json:"points_earned"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	QuizAttempt QuizAttempt `gorm:"foreignKey:QuizAttemptID"`
	Question    Question    `gorm:"foreignKey:QuestionID"`
}

// Certificate represents a certificate earned by a user
type Certificate struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	UserID            uint           `gorm:"not null;index" json:"user_id"`
	CourseID          uint           `gorm:"not null;index" json:"course_id"`
	CertificateNumber string         `gorm:"uniqueIndex;not null" json:"certificate_number"` // Unique certificate ID
	IssuedAt          time.Time      `gorm:"autoCreateTime" json:"issued_at"`
	ExpiresAt         *time.Time     `json:"expires_at"`
	Score             int            `json:"score"`
	CertificateURL    string         `json:"certificate_url"` // URL to download
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User   User   `gorm:"foreignKey:UserID"`
	Course Course `gorm:"foreignKey:CourseID"`
}

// CoinTransaction tracks GMFC coin transactions
type CoinTransaction struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"not null;index" json:"user_id"`
	Amount          int64          `gorm:"not null" json:"amount"`           // Can be negative for spending
	TransactionType string         `gorm:"not null" json:"transaction_type"` // earned, spent, redeemed, admin_adjustment
	Reason          string         `json:"reason"`                           // e.g., "Course Completion", "Quiz Score"
	ReferenceID     *uint          `json:"reference_id"`                     // e.g., CourseID or QuizID
	ReferenceType   string         `json:"reference_type"`                   // e.g., "course", "quiz"
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User User `gorm:"foreignKey:UserID"`
}

// Badge represents achievement badges
type Badge struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	IconURL     string         `json:"icon_url"`
	Level       string         `gorm:"not null" json:"level"`     // bronze, silver, gold, platinum
	Criteria    string         `gorm:"type:text" json:"criteria"` // JSON describing criteria
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	UserBadges []BadgeProgress `gorm:"foreignKey:BadgeID"`
}

// BadgeProgress tracks user's progress towards and achievement of badges
type BadgeProgress struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	BadgeID   uint           `gorm:"not null;index" json:"badge_id"`
	Progress  int            `gorm:"default:0" json:"progress"` // 0-100
	IsEarned  bool           `gorm:"default:false;index" json:"is_earned"`
	EarnedAt  *time.Time     `json:"earned_at"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User  User  `gorm:"foreignKey:UserID"`
	Badge Badge `gorm:"foreignKey:BadgeID"`
}

// CourseReview represents a user's review of a course
type CourseReview struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CourseID   uint           `gorm:"not null;index" json:"course_id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	Rating     int            `gorm:"not null" json:"rating"` // 1-5
	ReviewText string         `gorm:"type:text" json:"review_text"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Course Course `gorm:"foreignKey:CourseID"`
	User   User   `gorm:"foreignKey:UserID"`
}

// LearningReport represents aggregated learning data for reporting
type LearningReport struct {
	ID                      uint           `gorm:"primaryKey" json:"id"`
	UserID                  *uint          `json:"user_id"`                     // Null for organization-wide reports
	DepartmentID            *uint          `json:"department_id"`               // Null for user-specific reports
	ReportType              string         `gorm:"not null" json:"report_type"` // user, department, organization, course
	TotalEnrollments        int            `json:"total_enrollments"`
	TotalCompletions        int            `json:"total_completions"`
	CompletionRate          float64        `json:"completion_rate"` // percentage
	AvgCompletionTime       float64        `json:"avg_completion_time_hours"`
	AvgScores               float64        `json:"avg_scores"`
	CertificatesIssued      int            `json:"certificates_issued"`
	MandatoryCoursesPending int            `json:"mandatory_courses_pending"`
	TotalLearningHours      float64        `json:"total_learning_hours"`
	GeneratedAt             time.Time      `gorm:"autoCreateTime" json:"generated_at"`
	CreatedAt               time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt               gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User *User `gorm:"foreignKey:UserID"`
}

// DownloadLog tracks file downloads for compliance
type DownloadLog struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	MaterialID   uint           `gorm:"not null;index" json:"material_id"`
	DownloadedAt time.Time      `gorm:"autoCreateTime" json:"downloaded_at"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User     User           `gorm:"foreignKey:UserID"`
	Material LessonMaterial `gorm:"foreignKey:MaterialID"`
}

// SystemAuditLog logs important system events for compliance
type SystemAuditLog struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     *uint          `json:"user_id"`
	Action     string         `gorm:"not null;index" json:"action"` // login, course_enroll, course_complete, quiz_submit, etc
	EntityType string         `gorm:"not null" json:"entity_type"`  // User, Course, Quiz, etc
	EntityID   *uint          `json:"entity_id"`
	Details    string         `gorm:"type:text" json:"details"` // JSON additional details
	IPAddress  string         `json:"ip_address"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User *User `gorm:"foreignKey:UserID"`
}
