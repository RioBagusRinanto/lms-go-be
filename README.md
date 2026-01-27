# LMS (Learning Management System) Backend API

A professional, production-ready Learning Management System backend built with **Golang** and the **Gin web framework**. This comprehensive LMS platform supports course management, user progress tracking, gamification features, and comprehensive reporting capabilities.

## 📋 Features

### Core Features
- **User Management**: Role-based authentication (learners, instructors, admins)
- **Course Management**: Create, publish, and manage training courses
- **Video Lessons**: Structure courses with lessons and video content
- **Progress Tracking**: Automatic progress saving and completion tracking
- **Quizzes & Assessments**: Multiple question types with automatic and manual grading
- **Certificate Generation**: Automatic certificate creation upon course completion

### Gamification Features
- **GMFC Coins System**: Reward users with coins for achievements
- **Badge Levels**: Bronze, Silver, Gold, Platinum, Diamond badges
- **Learning Streaks**: Track consecutive days of learning
- **Leaderboards**: Competitive learning environment

### Dashboard & Reporting
- **Learner Dashboard**: View mandatory courses, in-progress courses, completed courses
- **Admin Analytics**: Comprehensive reporting and insights
- **Course Statistics**: Enrollment metrics, completion rates, performance analytics
- **User Statistics**: Learning hours, badges, coins, and achievements

## 🏗️ Project Structure

```
lms-go-be/
├── config/              # Configuration and database setup
│   ├── app.go          # Application configuration
│   └── database.go     # Database connection
├── models/             # Data models (User, Course, etc.)
│   └── models.go       # All database models
├── repositories/       # Data access layer
│   ├── user_repository.go
│   ├── course_repository.go
│   ├── enrollment_repository.go
│   └── other_repositories.go
├── services/           # Business logic layer
│   ├── auth_service.go
│   ├── course_service.go
│   └── gamification_service.go
├── handlers/           # HTTP request handlers
│   └── auth_handler.go
├── middleware/         # Custom middleware
│   └── auth.go        # Authentication & authorization
├── migrations/         # Database migrations
│   └── migration.go
├── seeders/           # Database seeders
│   └── seeder.go
├── utils/             # Utility functions
│   ├── helpers.go     # General helpers
│   ├── jwt.go         # JWT token management
│   └── response.go    # API response formatting
├── main.go            # Application entry point
├── go.mod             # Go module definition
├── .env               # Environment configuration
└── .env.example       # Example environment file
```

## 🛠️ Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin Gonic v1.9.1
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Database Migrations**: GORM Auto Migration

## 📦 Dependencies

Key packages used:
- `github.com/gin-gonic/gin` - Web framework
- `gorm.io/gorm` - ORM for database operations
- `gorm.io/driver/postgres` - PostgreSQL driver
- `github.com/golang-jwt/jwt/v5` - JWT token handling
- `golang.org/x/crypto` - Cryptographic functions
- `github.com/google/uuid` - UUID generation
- `github.com/joho/godotenv` - Environment variable loading

## 🚀 Getting Started

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 12+
- Git

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/your-username/lms-go-be.git
cd lms-go-be
```

2. **Install dependencies**
```bash
go mod download
go mod tidy
```

3. **Create PostgreSQL database**
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE lms_db;
```

4. **Configure environment variables**
```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your database credentials
nano .env  # or edit with your preferred editor
```

5. **Run the application**
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Authentication Endpoints

#### Register User
```bash
POST /api/auth/register
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "securepassword123",
  "role": "learner"
}
```

#### Login
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

Response:
```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "type": "Bearer"
  }
}
```

#### Get Profile
```bash
GET /api/auth/profile
Authorization: Bearer <token>
```

#### Update Profile
```bash
PUT /api/auth/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "first_name": "Jane",
  "bio": "Updated bio",
  "profile_image": "https://..."
}
```

### Dashboard Endpoints

#### Get Mandatory Courses
```bash
GET /api/dashboard/mandatory-courses
Authorization: Bearer <token>
```

#### Get In-Progress Courses
```bash
GET /api/dashboard/in-progress-courses
Authorization: Bearer <token>
```

#### Get Completed Courses
```bash
GET /api/dashboard/completed-courses
Authorization: Bearer <token>
```

#### Get User Coins
```bash
GET /api/dashboard/coins
Authorization: Bearer <token>
```

#### Get User Badges
```bash
GET /api/dashboard/badges
Authorization: Bearer <token>
```

#### Get Dashboard Stats
```bash
GET /api/dashboard/stats
Authorization: Bearer <token>
```

### Course Endpoints

#### Get All Courses
```bash
GET /api/courses
Authorization: Bearer <token>
```

#### Search Courses
```bash
GET /api/courses/search?q=golang
Authorization: Bearer <token>
```

#### Get Course Details
```bash
GET /api/courses/:id
Authorization: Bearer <token>
```

#### Enroll in Course
```bash
POST /api/courses/:id/enroll
Authorization: Bearer <token>
```

### Admin Endpoints

#### Create Course
```bash
POST /api/admin/courses
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Golang Fundamentals",
  "description": "Learn Go from scratch",
  "level": "beginner",
  "is_mandatory": true,
  "category_id": "...",
  "total_duration": 360
}
```

#### Get Course Statistics
```bash
GET /api/admin/courses/:id/stats
Authorization: Bearer <token>
```

## 🗄️ Database Schema

### Key Tables

**users**
- id (UUID primary key)
- first_name, last_name
- email (unique)
- password (hashed)
- role (learner, instructor, admin, hc_admin)
- gmfc_coins
- badge_level
- current_streak
- created_at, updated_at

**courses**
- id (UUID)
- title
- description
- instructor_id (FK to users)
- category_id (FK to categories)
- is_mandatory
- passing_score
- total_duration
- status (draft, published, archived)

**enrollments**
- id (UUID)
- user_id (FK to users)
- course_id (FK to courses)
- status (enrolled, in_progress, completed)
- progress (0-100)
- completion_date
- final_score

**lessons**
- id (UUID)
- course_id (FK to courses)
- title
- video_url
- duration
- order_index

**quizzes**
- id (UUID)
- course_id (FK to courses)
- title
- time_limit
- passing_score
- max_attempts

**certificates**
- id (UUID)
- user_id (FK to users)
- course_id (FK to courses)
- certificate_number
- issued_date
- file_url

## 🔐 Security Features

- **Password Hashing**: bcrypt with cost factor 10
- **JWT Authentication**: Token-based authentication with 7-day expiration
- **Role-Based Access Control**: Different endpoints for different user roles
- **Input Validation**: Request validation using Gin binding
- **SQL Injection Prevention**: GORM parameterized queries
- **CORS Support**: Configurable cross-origin requests

## 📝 Code Quality

The codebase follows clean architecture principles:

- **Models**: Define database schema
- **Repositories**: Data access abstraction layer
- **Services**: Business logic and rules
- **Handlers**: HTTP request/response processing
- **Middleware**: Cross-cutting concerns (auth, logging)
- **Utils**: Reusable utility functions

## 🧪 Testing

To run tests:
```bash
go test ./...
```

For coverage report:
```bash
go test -cover ./...
```

## 📈 Performance Optimizations

- Database indexing on frequently queried fields
- Connection pooling for database connections
- Request/response gzip compression
- Efficient database queries with proper joins
- Pagination for large datasets

## 🐛 Troubleshooting

### Database Connection Error
```
Error: Failed to connect to database
```
**Solution**: Check your `.env` file and ensure PostgreSQL is running and credentials are correct.

### JWT Token Expired
```
Error: token is expired
```
**Solution**: Request a new token using the `/api/auth/login` endpoint or refresh token endpoint.

### Port Already in Use
```
Error: listen tcp :8080: bind: address already in use
```
**Solution**: Change the PORT in `.env` file or kill the process using port 8080.

## 🚢 Deployment

### Docker Deployment

Create a `Dockerfile`:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o lms-api main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/lms-api .
COPY .env .
EXPOSE 8080
CMD ["./lms-api"]
```

Build and run:
```bash
docker build -t lms-api:latest .
docker run -p 8080:8080 lms-api:latest
```

### Environment Variables for Production
```env
ENV=production
JWT_SECRET=<strong-secret-key>
DB_HOST=<production-db-host>
DB_PASSWORD=<secure-password>
```

## 📄 API Response Format

All API responses follow this standard format:

**Success Response**:
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    ...
  }
}
```

**Error Response**:
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error information"
}
```

## 🔄 Workflow

1. **User Registration**: User registers with email and password
2. **User Login**: User logs in to get JWT token
3. **Course Enrollment**: User enrolls in available courses
4. **Progress Tracking**: System automatically tracks lesson completion
5. **Gamification**: Users earn coins and badges
6. **Certificate**: Upon course completion, certificate is generated
7. **Dashboard**: User views progress on personalized dashboard

## 📞 Support

For issues, questions, or contributions:
- Create an issue on GitHub
- Contact: support@lms-api.com
- Documentation: https://lms-api-docs.example.com

## 📄 License

This project is licensed under the MIT License - see LICENSE file for details.

## 🎯 Future Enhancements

- [ ] Video streaming with adaptive bitrate
- [ ] Advanced analytics and reporting
- [ ] Mobile application support
- [ ] Real-time notifications
- [ ] Discussion forums
- [ ] Peer review system
- [ ] Microlearning content
- [ ] Social learning features
- [ ] Integration with third-party services
- [ ] Advanced search with Elasticsearch

## 👥 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## ✨ Built with ❤️ by the LMS Team
