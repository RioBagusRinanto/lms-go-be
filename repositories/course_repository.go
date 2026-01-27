package repositories

import (
	"lms-go-be/models"

	"gorm.io/gorm"
)

// CourseRepository handles course data access
type CourseRepository struct {
	db *gorm.DB
}

// NewCourseRepository creates a new course repository instance
func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// Create creates a new course
// Parameters:
//   - course: course model to create
//
// Returns: created course and error if any
func (r *CourseRepository) Create(course *models.Course) (*models.Course, error) {
	if err := r.db.Create(course).Error; err != nil {
		return nil, err
	}
	return course, nil
}

// GetByID retrieves a course by ID with relations
// Parameters:
//   - id: course ID
//
// Returns: course model and error if not found
func (r *CourseRepository) GetByID(id string) (*models.Course, error) {
	var course models.Course
	if err := r.db.Preload("Lessons").Preload("Instructor").Preload("Category").
		First(&course, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

// Update updates an existing course
// Parameters:
//   - course: course model with updated data
//
// Returns: updated course and error if any
func (r *CourseRepository) Update(course *models.Course) (*models.Course, error) {
	if err := r.db.Save(course).Error; err != nil {
		return nil, err
	}
	return course, nil
}

// GetAll retrieves all courses with pagination
// Parameters:
//   - page: page number
//   - pageSize: number of records per page
//   - filters: map of filter key-value pairs (status, category_id, instructor_id, is_mandatory)
//
// Returns: slice of courses, total count, and error if any
func (r *CourseRepository) GetAll(page, pageSize int, filters map[string]interface{}) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	query := r.db.Where("status = ?", "published")

	// Apply filters
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Model(&models.Course{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Preload("Lessons").Preload("Instructor").Preload("Category").
		Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// GetMandatoryCourses retrieves all mandatory courses
// Returns: slice of mandatory courses and error if any
func (r *CourseRepository) GetMandatoryCourses() ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.Where("is_mandatory = ? AND status = ?", true, "published").
		Preload("Lessons").
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

// GetCoursesByInstructor retrieves courses by instructor ID
// Parameters:
//   - instructorID: instructor user ID
//
// Returns: slice of courses and error if any
func (r *CourseRepository) GetCoursesByInstructor(instructorID string) ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.Where("instructor_id = ?", instructorID).
		Preload("Lessons").
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

// GetCoursesByCategory retrieves courses by category ID
// Parameters:
//   - categoryID: category ID
//
// Returns: slice of courses and error if any
func (r *CourseRepository) GetCoursesByCategory(categoryID string) ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.Where("category_id = ? AND status = ?", categoryID, "published").
		Preload("Lessons").
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

// Delete deletes a course by ID
// Parameters:
//   - id: course ID
//
// Returns: error if any
func (r *CourseRepository) Delete(id string) error {
	return r.db.Delete(&models.Course{}, "id = ?", id).Error
}

// GetTotalCourses returns total number of published courses
// Returns: total count and error if any
func (r *CourseRepository) GetTotalCourses() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Course{}).Where("status = ?", "published").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetTotalDuration returns total duration of a course
// Parameters:
//   - courseID: course ID
//
// Returns: total duration in minutes and error if any
func (r *CourseRepository) GetTotalDuration(courseID string) (int, error) {
	var course models.Course
	if err := r.db.First(&course, "id = ?", courseID).Error; err != nil {
		return 0, err
	}
	return course.TotalDuration, nil
}

// SearchCourses searches courses by keyword
// Parameters:
//   - keyword: search keyword
//   - page: page number
//   - pageSize: number of records per page
//
// Returns: slice of courses, total count, and error if any
func (r *CourseRepository) SearchCourses(keyword string, page, pageSize int) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	query := r.db.Where("status = ? AND (title ILIKE ? OR description ILIKE ?)", "published", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Model(&models.Course{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Preload("Lessons").Preload("Instructor").Preload("Category").
		Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}
