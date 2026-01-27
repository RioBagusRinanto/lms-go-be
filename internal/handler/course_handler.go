package handler

import (
	"net/http"
	"strconv"

	"lms-go-be/internal/models"
	"lms-go-be/internal/repository"
	"lms-go-be/internal/service"
	"lms-go-be/internal/utils"

	"github.com/gin-gonic/gin"
)

// CourseHandler handles course endpoints
type CourseHandler struct {
	courseService *service.CourseService
	auditLogRepo  *repository.SystemAuditLogRepository
}

// NewCourseHandler creates a new course handler
func NewCourseHandler(courseService *service.CourseService, auditLogRepo *repository.SystemAuditLogRepository) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
		auditLogRepo:  auditLogRepo,
	}
}

// GetAllCourses gets all published courses
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	courses, total, err := h.courseService.GetAllCourses(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve courses", err.Error())
		return
	}

	dtos := make([]interface{}, len(courses))
	for i, course := range courses {
		dtos[i] = service.ConvertCourseToDTO(&course)
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Courses retrieved successfully", dtos, page, pageSize, total)
}

// GetCourse gets a single course
func (h *CourseHandler) GetCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid course ID", err.Error())
		return
	}

	course, err := h.courseService.GetCourse(uint(courseID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Course not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Course retrieved successfully", service.ConvertCourseToDTO(course))
}

// SearchCourses searches for courses
func (h *CourseHandler) SearchCourses(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Search query required", "q parameter is required")
		return
	}

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	courses, total, err := h.courseService.SearchCourses(query, page, 10)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Search failed", err.Error())
		return
	}

	dtos := make([]interface{}, len(courses))
	for i, course := range courses {
		dtos[i] = service.ConvertCourseToDTO(&course)
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Search completed", dtos, page, 10, total)
}

// GetByCategory gets courses by category
func (h *CourseHandler) GetByCategory(c *gin.Context) {
	category := c.Param("category")
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	courses, total, err := h.courseService.GetCoursesByCategory(category, page, 10)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve courses", err.Error())
		return
	}

	dtos := make([]interface{}, len(courses))
	for i, course := range courses {
		dtos[i] = service.ConvertCourseToDTO(&course)
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Courses retrieved successfully", dtos, page, 10, total)
}

// CreateCourse creates a new course (admin/instructor only)
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	instructorID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	var req service.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	course, err := h.courseService.CreateCourse(instructorID.(uint), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create course", err.Error())
		return
	}

	// Audit log
	_ = h.auditLogRepo.Create(&models.SystemAuditLog{
		UserID:     &instructorID.(uint),
		Action:     "course_created",
		EntityType: "course",
		EntityID:   &course.ID,
	})

	utils.SuccessResponse(c, http.StatusCreated, "Course created successfully", service.ConvertCourseToDTO(course))
}

// UpdateCourse updates a course
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid course ID", err.Error())
		return
	}

	var req service.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	course, err := h.courseService.UpdateCourse(uint(courseID), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update course", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Course updated successfully", service.ConvertCourseToDTO(course))
}

// DeleteCourse deletes a course
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid course ID", err.Error())
		return
	}

	if err := h.courseService.courseRepo.Delete(uint(courseID)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete course", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Course deleted successfully", nil)
}

// PublishCourse publishes a course
func (h *CourseHandler) PublishCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid course ID", err.Error())
		return
	}

	course, err := h.courseService.PublishCourse(uint(courseID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to publish course", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Course published successfully", service.ConvertCourseToDTO(course))
}

// AddReview adds a review to a course
type AddReviewRequest struct {
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	ReviewText string `json:"review_text"`
}

// AddReview adds a course review
func (h *CourseHandler) AddReview(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid course ID", err.Error())
		return
	}

	var req AddReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	review, err := h.courseService.AddReview(userID.(uint), uint(courseID), req.Rating, req.ReviewText)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to add review", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Review added successfully", review)
}

// GetReviews gets reviews for a course
func (h *CourseHandler) GetReviews(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid course ID", err.Error())
		return
	}

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	reviews, total, err := h.courseService.GetCourseReviews(uint(courseID), page, 10)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve reviews", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Reviews retrieved successfully", reviews, page, 10, total)
}
