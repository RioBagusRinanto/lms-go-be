package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User represents a system user (learner, instructor, or admin)
type User struct {
	ID            string         `gorm:"primaryKey" json:"id"`
	FirstName     string         `json:"first_name" gorm:"index"`
	LastName      string         `json:"last_name" gorm:"index"`
	Email         string         `json:"email" gorm:"uniqueIndex;index"`
	Password      string         `json:"-"` // Never expose password in JSON
	ProfileImage  string         `json:"profile_image"`
	Bio           string         `json:"bio"`
	Role          string         `json:"role" gorm:"index"` // learner, instructor, admin, hc_admin
	Department    string         `json:"department"`
	Status        string         `json:"status"` // active, inactive, suspended
	GMFCCoins     int64          `json:"gmfc_coins" gorm:"default:0"`
	BadgeLevel    string         `json:"badge_level" gorm:"default:'Bronze'"`
	CurrentStreak int            `json:"current_streak" gorm:"default:0"`
	IsVerified    bool           `json:"is_verified" gorm:"default:false"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Enrollments  []Enrollment   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	CoursesOwned []Course       `gorm:"foreignKey:InstructorID;constraint:OnDelete:SET NULL" json:"-"`
	QuizAttempts []QuizAttempt  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	CoinHistory  []CoinHistory  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	BadgeHistory []BadgeHistory `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate generates UUID for new user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// Category represents course categories
type Category struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" gorm:"uniqueIndex;index"`
	Slug      string    `json:"slug" gorm:"uniqueIndex;index"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Courses []Course `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL" json:"-"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// Course represents a training course
type Course struct {
	ID            string         `gorm:"primaryKey" json:"id"`
	Title         string         `json:"title" gorm:"index"`
	Description   string         `json:"description" gorm:"type:text"`
	Thumbnail     string         `json:"thumbnail"`
	TotalDuration int            `json:"total_duration"` // in minutes
	Level         string         `json:"level"`          // beginner, intermediate, advanced
	Status        string         `json:"status"`         // draft, published, archived
	IsMandatory   bool           `json:"is_mandatory" gorm:"index;default:false"`
	PassingScore  int            `json:"passing_score" gorm:"default:70"`
	InstructorID  string         `json:"instructor_id" gorm:"index"`
	CategoryID    string         `json:"category_id" gorm:"index"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Instructor  *User            `gorm:"foreignKey:InstructorID;constraint:OnDelete:SET NULL" json:"-"`
	Category    *Category        `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL" json:"-"`
	Lessons     []Lesson         `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"lessons,omitempty"`
	Quizzes     []Quiz           `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"-"`
	Enrollments []Enrollment     `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"-"`
	Materials   []CourseMaterial `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"-"`
}

func (c *Course) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// Lesson represents individual lessons within a course
type Lesson struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CourseID    string         `json:"course_id" gorm:"index"`
	Title       string         `json:"title" gorm:"index"`
	Description string         `json:"description" gorm:"type:text"`
	VideoURL    string         `json:"video_url"`
	Duration    int            `json:"duration"` // in minutes
	OrderIndex  int            `json:"order_index"`
	IsLocked    bool           `json:"is_locked" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Course         *Course          `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"-"`
	LessonProgress []LessonProgress `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE" json:"-"`
}

func (l *Lesson) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}

// LessonProgress tracks user progress in lessons
type LessonProgress struct {
	ID             string     `gorm:"primaryKey" json:"id"`
	UserID         string     `json:"user_id" gorm:"index;uniqueIndex:idx_user_lesson"`
	LessonID       string     `json:"lesson_id" gorm:"index;uniqueIndex:idx_user_lesson"`
	WatchDuration  int        `json:"watch_duration"` // in seconds
	IsCompleted    bool       `json:"is_completed" gorm:"index"`
	CompletionDate *time.Time `json:"completion_date"`
	LastAccessDate time.Time  `json:"last_access_date"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relations
	User   *User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Lesson *Lesson `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE" json:"-"`
}

func (lp *LessonProgress) BeforeCreate(tx *gorm.DB) error {
	if lp.ID == "" {
		lp.ID = uuid.New().String()
	}
	return nil
}

// Enrollment tracks user enrollment in courses
type Enrollment struct {
	ID             string         `gorm:"primaryKey" json:"id"`
	UserID         string         `json:"user_id" gorm:"index;uniqueIndex:idx_user_course"`
	CourseID       string         `json:"course_id" gorm:"index;uniqueIndex:idx_user_course"`
	Status         string         `json:"status"`   // enrolled, in_progress, completed, suspended
	Progress       int            `json:"progress"` // 0-100 percentage
	CompletionDate *time.Time     `json:"completion_date"`
	FinalScore     *int           `json:"final_score"`
	EnrolledAt     time.Time      `json:"enrolled_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User   *User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Course *Course `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"-"`
}

func (e *Enrollment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}

// Quiz represents a quiz within a course
type Quiz struct {
	ID             string         `gorm:"primaryKey" json:"id"`
	CourseID       string         `json:"course_id" gorm:"index"`
	Title          string         `json:"title" gorm:"index"`
	Description    string         `json:"description" gorm:"type:text"`
	TimeLimit      int            `json:"time_limit"` // in minutes, 0 means unlimited
	PassingScore   int            `json:"passing_score" gorm:"default:70"`
	MaxAttempts    int            `json:"max_attempts" gorm:"default:0"` // 0 means unlimited
	IsPublished    bool           `json:"is_published" gorm:"default:false"`
	RandomizeOrder bool           `json:"randomize_order" gorm:"default:false"`
	ShowAnswers    bool           `json:"show_answers" gorm:"default:true"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Course       *Course       `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"-"`
	Questions    []Question    `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"questions,omitempty"`
	QuizAttempts []QuizAttempt `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"-"`
}

func (q *Quiz) BeforeCreate(tx *gorm.DB) error {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return nil
}

// Question represents a quiz question
type Question struct {
	ID           string         `gorm:"primaryKey" json:"id"`
	QuizID       string         `json:"quiz_id" gorm:"index"`
	QuestionText string         `json:"question_text" gorm:"type:text"`
	Type         string         `json:"type"` // mcq, true_false, short_answer, fill_blank, matching, essay
	Points       int            `json:"points" gorm:"default:1"`
	OrderIndex   int            `json:"order_index"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Quiz    *Quiz            `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"-"`
	Options []QuestionOption `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"options,omitempty"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) error {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return nil
}

// QuestionOption represents answer options for questions
type QuestionOption struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	QuestionID string    `json:"question_id" gorm:"index"`
	Text       string    `json:"text" gorm:"type:text"`
	IsCorrect  bool      `json:"is_correct" gorm:"default:false"`
	OrderIndex int       `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	Question *Question `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"-"`
}

func (qo *QuestionOption) BeforeCreate(tx *gorm.DB) error {
	if qo.ID == "" {
		qo.ID = uuid.New().String()
	}
	return nil
}

// QuizAttempt represents a user's attempt at a quiz
type QuizAttempt struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	UserID    string         `json:"user_id" gorm:"index"`
	QuizID    string         `json:"quiz_id" gorm:"index"`
	StartTime time.Time      `json:"start_time"`
	EndTime   *time.Time     `json:"end_time"`
	Score     *int           `json:"score"`
	Status    string         `json:"status"`                      // in_progress, completed, submitted
	Responses datatypes.JSON `json:"responses" gorm:"type:jsonb"` // Store answer responses
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Quiz *Quiz `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"-"`
}

func (qa *QuizAttempt) BeforeCreate(tx *gorm.DB) error {
	if qa.ID == "" {
		qa.ID = uuid.New().String()
	}
	return nil
}

// CourseMaterial represents downloadable materials for a course
type CourseMaterial struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CourseID    string         `json:"course_id" gorm:"index"`
	Title       string         `json:"title" gorm:"index"`
	FileURL     string         `json:"file_url"`
	FileType    string         `json:"file_type"` // ppt, pdf, zip, doc
	FileSize    int64          `json:"file_size"`
	Version     int            `json:"version" gorm:"default:1"`
	Description string         `json:"description" gorm:"type:text"`
	OrderIndex  int            `json:"order_index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Course *Course `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"-"`
}

func (cm *CourseMaterial) BeforeCreate(tx *gorm.DB) error {
	if cm.ID == "" {
		cm.ID = uuid.New().String()
	}
	return nil
}

// Certificate represents a course completion certificate
type Certificate struct {
	ID                string     `gorm:"primaryKey" json:"id"`
	UserID            string     `json:"user_id" gorm:"index"`
	CourseID          string     `json:"course_id" gorm:"index"`
	CertificateNumber string     `json:"certificate_number" gorm:"uniqueIndex;index"`
	IssuedDate        time.Time  `json:"issued_date"`
	ExpiryDate        *time.Time `json:"expiry_date"`
	FileURL           string     `json:"file_url"`
	VerificationURL   string     `json:"verification_url"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Relations
	User   *User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Course *Course `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"-"`
}

func (c *Certificate) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// CoinHistory tracks GMFC coin transactions
type CoinHistory struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	UserID      string    `json:"user_id" gorm:"index"`
	Amount      int64     `json:"amount"`
	Reason      string    `json:"reason"`       // course_completion, quiz_passed, streak_bonus, etc
	ReferenceID string    `json:"reference_id"` // enrollment_id, quiz_attempt_id, etc
	CreatedAt   time.Time `json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (ch *CoinHistory) BeforeCreate(tx *gorm.DB) error {
	if ch.ID == "" {
		ch.ID = uuid.New().String()
	}
	return nil
}

// BadgeHistory tracks badge achievements
type BadgeHistory struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	UserID     string    `json:"user_id" gorm:"index"`
	BadgeLevel string    `json:"badge_level"` // Bronze, Silver, Gold, Platinum
	AchievedAt time.Time `json:"achieved_at"`
	CreatedAt  time.Time `json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (bh *BadgeHistory) BeforeCreate(tx *gorm.DB) error {
	if bh.ID == "" {
		bh.ID = uuid.New().String()
	}
	return nil
}

// Notification represents system notifications
type Notification struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	UserID    string         `json:"user_id" gorm:"index"`
	Title     string         `json:"title"`
	Message   string         `json:"message" gorm:"type:text"`
	Type      string         `json:"type"` // course, achievement, deadline, etc
	IsRead    bool           `json:"is_read" gorm:"index;default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}
