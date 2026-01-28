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

// EnrollmentHandler handles enrollment endpoints
type EnrollmentHandler struct {
	enrollmentService *service.EnrollmentService
	auditLogRepo      *repository.SystemAuditLogRepository
}

// NewEnrollmentHandler creates a new enrollment handler
func NewEnrollmentHandler(enrollmentService *service.EnrollmentService, auditLogRepo *repository.SystemAuditLogRepository) *EnrollmentHandler {
	return &EnrollmentHandler{
		enrollmentService: enrollmentService,
		auditLogRepo:      auditLogRepo,
	}
}

// Enroll enrolls a user in a course
func (h *EnrollmentHandler) Enroll(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	var req service.EnrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	userIDValue := userID.(uint)
	enrollment, err := h.enrollmentService.EnrollUser(userIDValue, req.CourseID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Enrollment failed", err.Error())
		return
	}

	// Audit log
	_ = h.auditLogRepo.Create(&models.SystemAuditLog{
		UserID:     &userIDValue,
		Action:     "course_enroll",
		EntityType: "enrollment",
		EntityID:   &enrollment.ID,
	})

	utils.SuccessResponse(c, http.StatusCreated, "Enrolled successfully", service.ConvertEnrollmentToDTO(enrollment))
}

// GetMyEnrollments gets all enrollments for the current user
func (h *EnrollmentHandler) GetMyEnrollments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	enrollments, total, err := h.enrollmentService.GetUserEnrollments(userID.(uint), page, 10)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve enrollments", err.Error())
		return
	}

	dtos := make([]interface{}, len(enrollments))
	for i, enrollment := range enrollments {
		dtos[i] = service.ConvertEnrollmentToDTO(&enrollment)
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Enrollments retrieved successfully", dtos, page, 10, total)
}

// GetInProgressCourses gets in-progress courses
func (h *EnrollmentHandler) GetInProgressCourses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	enrollments, err := h.enrollmentService.GetInProgressCourses(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve courses", err.Error())
		return
	}

	dtos := make([]interface{}, len(enrollments))
	for i, enrollment := range enrollments {
		dtos[i] = service.ConvertEnrollmentToDTO(&enrollment)
	}

	utils.SuccessResponse(c, http.StatusOK, "In-progress courses retrieved successfully", dtos)
}

// GetCompletedCourses gets completed courses
func (h *EnrollmentHandler) GetCompletedCourses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	enrollments, total, err := h.enrollmentService.GetCompletedCourses(userID.(uint), page, 10)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve courses", err.Error())
		return
	}

	dtos := make([]interface{}, len(enrollments))
	for i, enrollment := range enrollments {
		dtos[i] = service.ConvertEnrollmentToDTO(&enrollment)
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Completed courses retrieved successfully", dtos, page, 10, total)
}

// GetMandatoryCourses gets mandatory courses
func (h *EnrollmentHandler) GetMandatoryCourses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	enrollments, err := h.enrollmentService.GetMandatoryCourses(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve courses", err.Error())
		return
	}

	dtos := make([]interface{}, len(enrollments))
	for i, enrollment := range enrollments {
		dtos[i] = service.ConvertEnrollmentToDTO(&enrollment)
	}

	utils.SuccessResponse(c, http.StatusOK, "Mandatory courses retrieved successfully", dtos)
}

// ProgressHandler handles progress endpoints
type ProgressHandler struct {
	progressService *service.ProgressService
	auditLogRepo    *repository.SystemAuditLogRepository
}

// NewProgressHandler creates a new progress handler
func NewProgressHandler(progressService *service.ProgressService, auditLogRepo *repository.SystemAuditLogRepository) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
		auditLogRepo:    auditLogRepo,
	}
}

// TrackProgress tracks video progress
func (h *ProgressHandler) TrackProgress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	var req struct {
		CourseID        uint `json:"course_id" binding:"required"`
		LessonID        uint `json:"lesson_id" binding:"required"`
		WatchedDuration int  `json:"watched_duration" binding:"min=0"`
		TotalDuration   int  `json:"total_duration" binding:"min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	progress, err := h.progressService.TrackProgress(userID.(uint), req.CourseID, req.LessonID, req.WatchedDuration, req.TotalDuration)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to track progress", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Progress tracked successfully", progress)
}

// GetCourseProgress gets course progress
func (h *ProgressHandler) GetCourseProgress(c *gin.Context) {
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

	progresses, err := h.progressService.GetCourseProgress(userID.(uint), uint(courseID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve progress", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Progress retrieved successfully", progresses)
}

// GetLessonProgress gets lesson progress
func (h *ProgressHandler) GetLessonProgress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	lessonID, err := strconv.ParseUint(c.Param("lessonId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid lesson ID", err.Error())
		return
	}

	// This would need course ID from context or request
	// Simplified version - in production, get courseID properly
	var courseID uint
	if cID := c.Query("course_id"); cID != "" {
		if parsed, err := strconv.ParseUint(cID, 10, 32); err == nil {
			courseID = uint(parsed)
		}
	}

	progress, err := h.progressService.GetLessonProgress(userID.(uint), courseID, uint(lessonID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Progress not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Progress retrieved successfully", progress)
}

// QuizHandler handles quiz endpoints
type QuizHandler struct {
	quizService  *service.QuizService
	auditLogRepo *repository.SystemAuditLogRepository
}

// NewQuizHandler creates a new quiz handler
func NewQuizHandler(quizService *service.QuizService, auditLogRepo *repository.SystemAuditLogRepository) *QuizHandler {
	return &QuizHandler{
		quizService:  quizService,
		auditLogRepo: auditLogRepo,
	}
}

// StartAttempt starts a quiz attempt
func (h *QuizHandler) StartAttempt(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	var req struct {
		QuizID uint `json:"quiz_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	attempt, err := h.quizService.StartAttempt(userID.(uint), req.QuizID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to start quiz", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Quiz attempt started", attempt)
}

// SubmitAttempt submits quiz answers
func (h *QuizHandler) SubmitAttempt(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	attemptID, err := strconv.ParseUint(c.Param("attemptId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid attempt ID", err.Error())
		return
	}

	var req struct {
		QuizID    uint              `json:"quiz_id" binding:"required"`
		Answers   map[string]string `json:"answers"`
		TimeSpent int               `json:"time_spent"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Convert string keys to uint
	answers := make(map[uint]string)
	for k, v := range req.Answers {
		if parsed, err := strconv.ParseUint(k, 10, 32); err == nil {
			answers[uint(parsed)] = v
		}
	}

	attempt, err := h.quizService.SubmitAttempt(userID.(uint), req.QuizID, uint(attemptID), answers, req.TimeSpent)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to submit quiz", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Quiz submitted successfully", attempt)
}

// GetAttempts gets quiz attempts
func (h *QuizHandler) GetAttempts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	quizID, err := strconv.ParseUint(c.Param("quizId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID", err.Error())
		return
	}

	attempts, err := h.quizService.GetUserAttempts(userID.(uint), uint(quizID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve attempts", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Attempts retrieved successfully", attempts)
}
