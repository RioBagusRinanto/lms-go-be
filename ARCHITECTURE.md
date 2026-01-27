# LMS Backend Architecture & Code Organization

## Overview

This document describes the professional clean architecture used in the LMS Backend, following best practices for maintainability, testability, and scalability.

## Architecture Layers

```
┌─────────────────────────────────────────┐
│           HTTP Clients / UI             │  ← External Interface
├─────────────────────────────────────────┤
│            Middleware Layer             │  ← Cross-cutting concerns
│  (Authentication, CORS, Logging, etc)   │
├─────────────────────────────────────────┤
│          Handlers/Controllers           │  ← Request/Response Processing
│  (Parse requests, call services)        │
├─────────────────────────────────────────┤
│         Services/Business Logic         │  ← Core Business Rules
│  (Validation, Calculations, Rules)      │
├─────────────────────────────────────────┤
│      Repositories/Data Access           │  ← Data Abstraction
│  (Database queries, transactions)       │
├─────────────────────────────────────────┤
│         Models/Entities                 │  ← Data Structures
│  (Database schemas, data types)         │
├─────────────────────────────────────────┤
│         Database (PostgreSQL)           │  ← Persistent Storage
└─────────────────────────────────────────┘
```

## Layer Responsibilities

### 1. **Models Layer** (`models/models.go`)

Defines database entities and their relationships.

**Responsibility**: Data structure definition
**Key Concepts**:
- GORM model tags for database mapping
- Relationships (HasMany, BelongsTo, etc.)
- Field validation tags
- Lifecycle hooks (BeforeCreate, etc.)

**Example**:
```go
type User struct {
    ID        string    `gorm:"primaryKey"`
    Email     string    `gorm:"uniqueIndex"`
    Password  string    `json:"-"`  // Never expose password
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 2. **Repository Layer** (`repositories/`)

Provides data access abstraction, isolating business logic from database details.

**Responsibility**: 
- Execute database queries
- Map database results to models
- Handle database transactions
- Implement pagination and filtering

**Pattern**: Repository Pattern
- One repository per entity
- Dependency injection via constructor
- Interface-based design (can be mocked for testing)

**Example**:
```go
type UserRepository struct {
    db *gorm.DB
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
    var user models.User
    if err := r.db.First(&user, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

**Benefits**:
- ✅ Database changes don't affect business logic
- ✅ Easy to mock for unit testing
- ✅ Consistent data access patterns
- ✅ Single source of truth for queries

### 3. **Service Layer** (`services/`)

Contains business logic and application rules.

**Responsibility**:
- Validate input data
- Implement business rules
- Orchestrate repository calls
- Handle errors and edge cases
- Perform calculations and transformations

**Pattern**: Service Pattern
- One service per business domain
- May depend on multiple repositories
- No HTTP context knowledge
- Pure business logic

**Example**:
```go
type CourseService struct {
    courseRepo     *repositories.CourseRepository
    enrollmentRepo *repositories.EnrollmentRepository
}

func (s *CourseService) EnrollUser(userID, courseID string) error {
    // Validation
    if userID == "" || courseID == "" {
        return errors.New("invalid input")
    }
    
    // Business logic
    course, err := s.courseRepo.GetByID(courseID)
    if err != nil {
        return err
    }
    
    // More logic...
}
```

**Benefits**:
- ✅ Business logic centralized and reusable
- ✅ Easy to test (mock repositories)
- ✅ Framework-agnostic
- ✅ Clear separation of concerns

### 4. **Handler Layer** (`handlers/`)

Processes HTTP requests and responses.

**Responsibility**:
- Parse HTTP requests
- Validate request format
- Call appropriate services
- Format and return HTTP responses
- Handle HTTP-specific errors

**Pattern**: Handler/Controller Pattern
- One handler per entity or domain
- Thin logic (delegates to services)
- HTTP-aware only
- Request validation using Gin bindings

**Example**:
```go
type AuthHandler struct {
    authService *services.AuthService
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequestError(c, "invalid request", err.Error())
        return
    }
    
    token, err := h.authService.Login(req.Email, req.Password)
    if err != nil {
        utils.UnauthorizedError(c, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "login successful", gin.H{
        "token": token,
    })
}
```

**Benefits**:
- ✅ Clean HTTP abstraction
- ✅ Easy to test (mock services)
- ✅ Consistent response formatting
- ✅ Proper error handling

### 5. **Middleware Layer** (`middleware/`)

Cross-cutting concerns that apply to multiple routes.

**Responsibility**:
- Authentication and authorization
- Request/response logging
- CORS handling
- Error handling
- Request/response transformation

**Example**:
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        claims, err := utils.ValidateToken(token)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        
        c.Set("user_id", claims.ID)
        c.Next()
    }
}
```

**Benefits**:
- ✅ DRY (Don't Repeat Yourself)
- ✅ Consistent behavior across routes
- ✅ Easy to enable/disable features
- ✅ Chainable middleware

### 6. **Utils Layer** (`utils/`)

Shared utility functions.

**Responsibility**:
- Cryptographic operations (password hashing, JWT)
- Common calculations
- Data validation
- Response formatting

**Organization**:
- `helpers.go`: General utilities
- `jwt.go`: JWT token management
- `response.go`: API response formatting

## Data Flow Example: User Registration

```
1. HTTP Request (POST /api/auth/register)
   ↓
2. Middleware Layer
   - CORSMiddleware() - Handle CORS
   - LoggingMiddleware() - Log request
   ↓
3. Handler Layer (AuthHandler.Register)
   - Parse JSON request
   - Validate request format
   ↓
4. Service Layer (AuthService.Register)
   - Validate email format
   - Check if user exists
   - Hash password
   ↓
5. Repository Layer (UserRepository.Create)
   - Insert into database
   ↓
6. Models Layer
   - GORM executes SQL INSERT
   ↓
7. Response Back Up The Stack
   - Formatted as SuccessResponse
   ↓
8. HTTP Response (JSON)
```

## Design Patterns Used

### 1. **Dependency Injection**
```go
// Services receive dependencies via constructor
func NewCourseService(
    courseRepo *repositories.CourseRepository,
    enrollmentRepo *repositories.EnrollmentRepository,
) *CourseService {
    return &CourseService{
        courseRepo: courseRepo,
        enrollmentRepo: enrollmentRepo,
    }
}
```

**Benefits**: Loose coupling, easy testing, flexible configuration

### 2. **Repository Pattern**
```go
// Abstracts database access
type CourseRepository interface {
    Create(course *models.Course) error
    GetByID(id string) (*models.Course, error)
    Update(course *models.Course) error
}
```

**Benefits**: Database independence, easier testing

### 3. **Middleware Pattern**
```go
// Chain of responsibility
router.Use(middleware.CORSMiddleware())
router.Use(middleware.AuthMiddleware())
```

**Benefits**: Composable, reusable, chainable

### 4. **Error Handling**
```go
// Consistent error response
utils.BadRequestError(c, "message", "details")
utils.NotFoundError(c, "message")
utils.InternalServerError(c, "message", "details")
```

**Benefits**: Consistent API error responses, easy client handling

## File Organization Strategy

```
lms-go-be/
├── main.go                    # Entry point, route setup
├── config/                    # Configuration & database
│   ├── app.go                # App configuration
│   └── database.go           # Database initialization
├── models/                    # Data models
│   └── models.go             # All entities
├── repositories/              # Data access layer
│   ├── user_repository.go
│   ├── course_repository.go
│   ├── enrollment_repository.go
│   └── other_repositories.go
├── services/                  # Business logic
│   ├── auth_service.go       # Authentication/authorization
│   ├── course_service.go     # Course & enrollment logic
│   └── gamification_service.go # Gamification logic
├── handlers/                  # HTTP handlers
│   ├── auth_handler.go
│   └── (more handlers added as needed)
├── middleware/                # Middleware functions
│   └── auth.go
├── migrations/                # Database migrations
│   └── migration.go
├── seeders/                   # Database seeders
│   └── seeder.go
└── utils/                     # Utilities
    ├── helpers.go            # General utilities
    ├── jwt.go                # JWT operations
    └── response.go           # Response formatting
```

## Code Quality Principles

### 1. **Single Responsibility Principle**
Each file/function has one reason to change.

```go
// ❌ Bad: UserRepository doing too much
func (r *UserRepository) CreateAndEnroll(user, courseID) {}

// ✅ Good: Separate responsibilities
func (r *UserRepository) Create(user) {}
func (r *EnrollmentRepository) Create(enrollment) {}
```

### 2. **Dependency Inversion**
Depend on abstractions, not implementations.

```go
// ❌ Bad: Direct dependency
type Service struct {
    db *gorm.DB
}

// ✅ Good: Repository abstraction
type Service struct {
    repo *repositories.UserRepository
}
```

### 3. **DRY (Don't Repeat Yourself)**
Extract common logic into utilities or helpers.

```go
// Instead of repeating validation in multiple handlers
// Use centralized validation function
if !utils.ValidateEmail(email) {
    return errors.New("invalid email")
}
```

### 4. **SOLID Principles**
- **S**ingle Responsibility: One reason to change
- **O**pen/Closed: Open for extension, closed for modification
- **L**iskov Substitution: Subtypes should be substitutable
- **I**nterface Segregation: Small, focused interfaces
- **D**ependency Inversion: Depend on abstractions

## Testing Strategy

### Unit Tests Example
```go
func TestUserService_Register(t *testing.T) {
    // Arrange: Setup mocks
    mockRepo := &MockUserRepository{}
    service := services.NewAuthService(mockRepo)
    
    // Act: Call function
    user, err := service.Register("John", "Doe", "john@example.com", "pass", "learner")
    
    // Assert: Check results
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if user.Email != "john@example.com" {
        t.Errorf("Expected email john@example.com, got %s", user.Email)
    }
}
```

## Documentation Standards

### Function Comments
```go
// CreateCourse creates a new course
// Parameters:
//   - course: course model to create
// Returns: created course and error if any
func (s *CourseService) CreateCourse(course *models.Course) (*models.Course, error) {
    // Implementation
}
```

### Inline Comments for Complex Logic
```go
// Calculate badge level based on completed courses
// Each badge level requires more courses than the previous
newBadgeLevel := utils.CalculateBadgeLevel(completedCount)
```

## Performance Considerations

### Database
- **Indexing**: Applied to frequently queried fields
- **Eager Loading**: Use Preload() to avoid N+1 queries
- **Pagination**: Implemented for large datasets
- **Connection Pooling**: Configured in GORM

### API
- **Caching**: Can be added at service layer
- **Rate Limiting**: Can be added via middleware
- **Compression**: Enabled by default in Gin
- **Lazy Loading**: Fields loaded on-demand

### Example: Eager Loading
```go
// ❌ Bad: N+1 query problem
courses := courseRepo.GetAll()
for _, course := range courses {
    instructor := getUserByID(course.InstructorID) // Extra query each iteration
}

// ✅ Good: Eager loading
courses := db.Preload("Instructor").Find(&courses)
```

## Security Best Practices

1. **Password Security**
   - Hashed with bcrypt
   - Never logged or exposed in API

2. **JWT Tokens**
   - Signed with secret key
   - 7-day expiration
   - Validated on each request

3. **SQL Injection Prevention**
   - GORM parameterized queries
   - No string concatenation for SQL

4. **Authorization**
   - Role-based access control
   - Checked in middleware

5. **Input Validation**
   - Required fields validated
   - Email format validated
   - Type-checked via Gin bindings

## Extending the Application

### Adding New Entity (e.g., Assignment)

1. **Add Model** (`models/models.go`)
```go
type Assignment struct {
    ID        string
    CourseID  string
    Title     string
    // ...
}
```

2. **Create Repository** (`repositories/assignment_repository.go`)
```go
type AssignmentRepository struct {
    db *gorm.DB
}
```

3. **Create Service** (`services/assignment_service.go`)
```go
type AssignmentService struct {
    repo *repositories.AssignmentRepository
}
```

4. **Create Handler** (`handlers/assignment_handler.go`)
```go
type AssignmentHandler struct {
    service *services.AssignmentService
}
```

5. **Add Routes** (`main.go`)
```go
assignments := router.Group("/api/assignments")
assignments.POST("", assignmentHandler.Create)
assignments.GET("/:id", assignmentHandler.GetByID)
```

6. **Migration Runs Automatically** (on startup)

---

## Summary

This architecture provides:
- ✅ **Clean Code**: Clear separation of concerns
- ✅ **Testability**: Easy to mock and test each layer
- ✅ **Maintainability**: Changes isolated to relevant layers
- ✅ **Scalability**: Easy to add new features
- ✅ **Professional**: Follows Go best practices
- ✅ **Documented**: Comments explain complex logic

By following these patterns, the codebase remains professional, maintainable, and scalable as the project grows.
