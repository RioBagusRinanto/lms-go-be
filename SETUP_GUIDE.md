# LMS Backend Setup & Installation Guide

## Prerequisites

- **Go**: Version 1.21 or higher
  - Download from: https://golang.org/dl/
  - Verify: `go version`

- **PostgreSQL**: Version 12 or higher
  - Download from: https://www.postgresql.org/download/
  - Verify: `psql --version`

- **Git**: For version control
  - Download from: https://git-scm.com/
  - Verify: `git --version`

## Step 1: Create PostgreSQL Database

### Using PostgreSQL Command Line:

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE lms_db;

# Create user (optional, for security)
CREATE USER lms_user WITH PASSWORD 'secure_password';

# Grant privileges
GRANT ALL PRIVILEGES ON DATABASE lms_db TO lms_user;

# Exit
\q
```

### Using pgAdmin (GUI):
1. Open pgAdmin
2. Right-click on "Databases"
3. Create > Database
4. Name: `lms_db`
5. Save

## Step 2: Project Setup

### Clone or Navigate to Project

```bash
cd d:\workspace\learning\golang\gin\lms-go-be
```

### Verify Project Structure

```bash
dir
# Should show: main.go, go.mod, config/, models/, etc.
```

## Step 3: Configure Environment

### Edit `.env` File

```bash
# Windows
notepad .env

# Or use VS Code
code .env
```

### Set Configuration

```env
# Server
ENV=development
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres           # or your PostgreSQL user
DB_PASSWORD=password       # Your PostgreSQL password
DB_NAME=lms_db

# JWT
JWT_SECRET=your-super-secret-key-change-in-production
```

## Step 4: Download Dependencies

```bash
cd d:\workspace\learning\golang\gin\lms-go-be

# Download all dependencies
go mod download

# Tidy up dependencies
go mod tidy

# Verify everything is installed
go list -m all
```

## Step 5: Run the Application

### Development Mode

```bash
go run main.go
```

Expected output:
```
Database connected successfully!
Running database migrations...
Database migrations completed successfully!
Seeding database with initial data...
Users seeded successfully
Categories seeded successfully
Courses seeded successfully
Lessons seeded successfully
Quizzes seeded successfully
Database seeding completed successfully!
Starting server on port 8080...
[GIN-debug] Loaded HTML Templates (2): ...
```

### Build for Production

```bash
# Build executable
go build -o lms-api.exe main.go

# Run built executable
.\lms-api.exe
```

## Step 6: Verify Installation

### Health Check Endpoint

```bash
# Open new terminal and run:
curl http://localhost:8080/api/health

# Expected response:
# {
#   "status": "ok",
#   "message": "LMS Backend API is running"
# }
```

### Test Authentication

```bash
# Register new user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Test",
    "last_name": "User",
    "email": "test@example.com",
    "password": "password123",
    "role": "learner"
  }'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## Troubleshooting

### Issue: "Database connection error"

**Solution:**
```bash
# Check PostgreSQL is running
# Windows Services: Check if "PostgreSQL" service is running
# Or restart with:
pg_ctl -D "C:\Program Files\PostgreSQL\15\data" start

# Verify database exists
psql -U postgres -l
```

### Issue: "Port 8080 already in use"

**Solution:**
```bash
# Windows: Find and kill process using port 8080
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Or change port in .env:
PORT=8081
```

### Issue: "Module not found" error

**Solution:**
```bash
# Clean cache and reinstall
go clean -modcache
go get -u ./...
go mod tidy
```

### Issue: "Password authentication failed"

**Solution:**
```bash
# Verify PostgreSQL password
# Reset password if needed:
psql -U postgres
ALTER USER postgres WITH PASSWORD 'newpassword';
\q

# Update .env with new password
```

## Using Test Data

The application automatically seeds initial test data:

### Test Users (from seeders):
```
Email: john@example.com | Password: password123 | Role: learner
Email: jane@example.com | Password: password123 | Role: instructor
Email: admin@example.com | Password: admin123 | Role: admin
Email: alice@example.com | Password: password123 | Role: learner
```

### Sample Data:
- 5 Users
- 5 Categories
- 4 Courses (2 mandatory)
- 12 Lessons
- 4 Quizzes

## Database Reset

### Clear and Reseed Database

```go
// In main.go, uncomment these lines:
// seeders.ClearDatabase(db)
// seeders.SeedDatabase(db)

// Then run:
go run main.go
```

### Or manually drop and recreate:

```bash
# Connect to PostgreSQL
psql -U postgres

# Drop and recreate database
DROP DATABASE lms_db;
CREATE DATABASE lms_db;

# Exit and run application
\q
go run main.go
```

## Development Workflow

### 1. Create New Feature

```go
// 1. Update models if needed (models/models.go)
// 2. Migrations run automatically on startup
// 3. Create repository (repositories/)
// 4. Create service (services/)
// 5. Create handler (handlers/)
// 6. Add routes (main.go - setupRouter)
```

### 2. Test APIs

Use any REST client:
- **Postman**: https://www.postman.com/downloads/
- **Thunder Client**: VS Code Extension
- **curl**: Command line
- **Insomnia**: https://insomnia.rest/

### 3. Code Quality

```bash
# Format code
go fmt ./...

# Lint code (install golangci-lint first)
golangci-lint run

# Run tests
go test ./...

# Test with coverage
go test -cover ./...
```

## Production Deployment

### Docker Deployment

1. **Create Dockerfile**:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o lms-api .

FROM alpine:latest
RUN apk add --no-cache postgresql-client
WORKDIR /root
COPY --from=builder /app/lms-api .
COPY .env .
EXPOSE 8080
CMD ["./lms-api"]
```

2. **Build and Run**:
```bash
docker build -t lms-api:latest .
docker run -p 8080:8080 --env-file .env lms-api:latest
```

### Environment Variables for Production

```env
ENV=production
PORT=8080
DB_HOST=production-db-host
DB_PORT=5432
DB_USER=lms_user
DB_PASSWORD=<strong-password>
DB_NAME=lms_db
JWT_SECRET=<generate-strong-secret>
```

### Generate Strong JWT Secret

```bash
# Bash/Linux/Mac:
openssl rand -base64 32

# Windows PowerShell:
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Maximum 256 }))
```

## Performance Optimization

### Database Connection Pooling
Already configured in GORM with:
- Max Open Connections: 25
- Max Idle Connections: 5
- Connection Max Lifetime: 5 minutes

### API Response Compression
Gin automatically compresses responses using gzip when appropriate.

### Database Indexing
Automatically created on:
- user.email
- course.status
- enrollment.user_id, course_id
- lesson_progress.user_id, lesson_id

## Monitoring & Logging

### Current Logging
- Database connections
- Server startup
- Data seeding

### For Production, Add:
```go
import "log"
import "github.com/sirupsen/logrus" // Recommended logging library

// Setup structured logging
// Log all API requests/responses
// Log database queries (slow query logs)
// Log errors with stack traces
```

## Next Steps

1. ✅ **Complete Backend API Implementation**
   - Quiz/Assessment endpoints
   - Certificate generation
   - Gamification endpoints
   - Admin analytics endpoints

2. **Frontend Integration**
   - Build React/Vue frontend
   - Connect to API endpoints
   - Implement authentication flow

3. **Testing**
   - Unit tests for services
   - Integration tests for API
   - Load testing

4. **Deployment**
   - Setup CI/CD pipeline
   - Deploy to cloud (AWS, GCP, Azure)
   - Setup monitoring and alerts

## Additional Resources

- **GORM Documentation**: https://gorm.io/
- **Gin Framework**: https://gin-gonic.com/
- **PostgreSQL**: https://www.postgresql.org/docs/
- **JWT**: https://jwt.io/
- **Go Best Practices**: https://golang.org/doc/effective_go

## Support

For issues or questions:
1. Check this guide's Troubleshooting section
2. Check application logs
3. Review code comments
4. Refer to package documentation

---

**Happy coding! 🚀**
