# Postman Collection - Quick Reference

Professional testing collection for LMS Backend API with automated test cases and environment variables.

## 📦 Files Included

| File | Purpose |
|------|---------|
| `LMS-Backend-Collection.postman_collection.json` | Complete API collection with 50+ endpoints |
| `LMS-Backend-Environment.postman_environment.json` | Environment variables and configuration |
| `POSTMAN_TESTING_GUIDE.md` | Detailed testing guide with scenarios |

## ⚡ Quick Import (2 Minutes)

### Step 1: Open Postman
- Launch Postman application
- Click "Import" button (top-left corner)

### Step 2: Import Collection
- Select: `LMS-Backend-Collection.postman_collection.json`
- Click "Import"

### Step 3: Import Environment
- Click "Import" again
- Select: `LMS-Backend-Environment.postman_environment.json`
- Click "Import"

### Step 4: Select Environment
- Top-right dropdown (next to Send button)
- Select: "LMS Backend - Development Environment"

### Step 5: Start Testing
- Open any folder in the collection
- Select a request
- Click "Send"

---

## 🎯 What's Included

### 10 Main Categories

1. **Authentication** (5 requests)
   - Register, Login, Profile, Update, Change Password

2. **Courses** (8 requests)
   - List, Search, Filter, Details, Create, Update, Publish, Review

3. **Enrollments** (5 requests)
   - Enroll, View All, In-Progress, Completed, Mandatory

4. **Progress Tracking** (3 requests)
   - Track Progress, Course Progress, Lesson Progress

5. **Quizzes** (3 requests)
   - Start Quiz, Submit Answers, View Attempts

6. **Gamification** (4 requests)
   - Coins, Transactions, Badges, Leaderboard

7. **Dashboard** (1 request)
   - Complete Dashboard with all metrics

8. **User Management** (2 requests)
   - User Profile, User Statistics

9. **Admin Operations** (3 requests)
   - List Users, Search Users, Audit Logs

10. **Test Workflow** (10 requests)
    - Complete end-to-end workflow

---

## 🔄 Common Workflows

### 1️⃣ Complete User Journey (10 minutes)

```
1. Run: "1. Register Test User"          → Creates new account
2. Run: "2. Login and Get Token"         → Gets JWT token (auto-saved)
3. Run: "3. Browse Available Courses"    → Lists all courses
4. Run: "4. Enroll in Course"            → Enrolls in first course
5. Run: "5. Track Video Progress"        → Watches 90% of video
6. Run: "6. Start Quiz"                  → Begins quiz attempt
7. Run: "7. Submit Quiz"                 → Submits answers & gets score
8. Run: "8. View Dashboard"              → Sees all updated stats
9. Run: "9. Check Coins and Badges"      → Views gamification rewards
10. Run: "10. View Leaderboard"          → Sees competitive ranking
```

### 2️⃣ Authentication Testing (2 minutes)

```
Go to: 1. Authentication
├─ Register User              ✓ New account
├─ Login User                 ✓ Token obtained
├─ Get Current User Profile   ✓ User info
├─ Update User Profile        ✓ Profile updated
└─ Change Password            ✓ Password changed
```

### 3️⃣ Course Management (5 minutes)

```
Go to: 2. Courses
├─ Get All Courses           ✓ Browse all (public)
├─ Search Courses            ✓ Search functionality
├─ Get by Category           ✓ Filter by category
├─ Get Course Details        ✓ Full course info
├─ Create Course (Admin)     ✓ Admin-only
├─ Update Course (Admin)     ✓ Admin-only
├─ Publish Course (Admin)    ✓ Admin-only
└─ Add Course Review         ✓ User review
```

### 4️⃣ Learning Progress (5 minutes)

```
Go to: 4. Progress Tracking
├─ Track Video Progress      ✓ Update watched %
├─ Get Course Progress       ✓ Overall completion
└─ Get Lesson Progress       ✓ Per-lesson status
```

### 5️⃣ Quiz Assessment (5 minutes)

```
Go to: 5. Quizzes & Assessments
├─ Start Quiz Attempt        ✓ Begin quiz
├─ Submit Quiz Attempt       ✓ Get score
└─ Get Quiz Attempts         ✓ View history
```

### 6️⃣ Admin Operations (5 minutes)

```
Go to: 9. Admin Operations
├─ Get All Users             ✓ List with pagination
├─ Search Users              ✓ Find specific user
└─ Get Audit Logs            ✓ Compliance logs
```

---

## 📊 Pre-configured Requests

### Authentication Flow
```
Register → Login → Get Token (auto-saved) → Use in all subsequent requests
```

### Environment Variables Auto-Set
| Variable | Set By | Where |
|----------|--------|-------|
| `auth_token` | Login endpoint | "Tests" tab |
| `user_id` | Login endpoint | "Tests" tab |
| `course_id` | Course retrieval | "Tests" tab |
| `enrollment_id` | Enrollment endpoint | "Tests" tab |
| `attempt_id` | Start Quiz | "Tests" tab |

### Bearer Token Auto-Applied
All protected requests automatically use `{{auth_token}}` from environment.

---

## 🚀 First-Time Setup

### Prerequisites
- ✅ Postman installed (free version OK)
- ✅ LMS Backend running (`go run cmd/main.go`)
- ✅ PostgreSQL running with database created
- ✅ `.env` file configured

### Setup Steps

1. **Import Collection & Environment** (see Quick Import above)

2. **Verify Base URL**
   - Click "Environment" dropdown
   - Select "LMS Backend - Development Environment"
   - Check `base_url` = `http://localhost:8080`
   - (Update if running on different host/port)

3. **Optional: Add Admin Token**
   - If testing admin endpoints, login as admin first:
     ```json
     Email: admin@lms.com
     Password: admin@123
     ```
   - Copy the token from response
   - Set `admin_token` variable in environment

4. **Run Test Workflow**
   - Go to "10. Test Workflow" folder
   - Run the requests in order
   - All should pass with green checkmarks

---

## ✅ Automated Test Cases

### Every Request Includes Tests For:
- ✓ HTTP Status Code (200, 201, 400, 401, etc.)
- ✓ Response Structure (required fields)
- ✓ Data Types (strings, numbers, arrays)
- ✓ Token Extraction (for auth endpoints)

### View Test Results
1. Send a request
2. Click "Tests" tab below response
3. Green checkmarks = passing
4. Red X = failing (check error message)

### Example Test Output
```
✓ Status code is 200
✓ Response has pagination
✓ Login successful and token obtained
```

---

## 🔑 Key Features

### 1. Automatic Variable Management
- Tokens automatically saved after login
- IDs automatically saved from responses
- Reused throughout collection

### 2. Comprehensive Endpoint Coverage
- 50+ endpoints tested
- Public, protected, and admin routes
- All HTTP methods (GET, POST, PUT, DELETE)

### 3. Built-in Documentation
- Each request has description
- Parameters documented in URL
- Body examples included

### 4. Pre-request Scripts
- Dynamic test email generation
- Timestamp-based unique values
- No manual data entry needed

### 5. Test Scripts
- Automatic validation
- Error detection
- Data extraction

---

## 🧪 Request Anatomy

### Example: Login Request

```json
{
  "name": "Login User",
  "request": {
    "method": "POST",
    "url": "{{base_url}}/api/v1/public/auth/login",
    "body": {
      "email": "{{test_email}}",
      "password": "TestPass123!"
    }
  },
  "tests": {
    "verify_token": "jsonData.data.token exists",
    "save_token": "pm.environment.set('auth_token', token)"
  }
}
```

### How Variables Work
- `{{base_url}}` → `http://localhost:8080`
- `{{test_email}}` → `testuser1234567890@example.com`
- `{{auth_token}}` → JWT token (auto-saved)

---

## 📝 Environment Variables Reference

### Default Values

| Variable | Default | Purpose |
|----------|---------|---------|
| `base_url` | localhost:8080 | API endpoint |
| `course_id` | 1 | Default course ID |
| `lesson_id` | 1 | Default lesson ID |
| `quiz_id` | 1 | Default quiz ID |
| `admin_email` | admin@lms.com | Admin credentials |
| `admin_password` | admin@123 | Admin credentials |
| `learner_email` | learner1@lms.com | Learner credentials |
| `learner_password` | learner@123 | Learner credentials |

### Runtime Values (Auto-Set)

| Variable | Set By | Used In |
|----------|--------|---------|
| `auth_token` | Login endpoint | All protected requests |
| `user_id` | Login response | User-specific endpoints |
| `enrollment_id` | Enroll endpoint | Progress endpoints |
| `attempt_id` | Start Quiz | Submit Quiz endpoint |
| `test_email` | Dynamic script | Registration |

---

## 🎓 Learning Resources

### In This Collection
1. **Complete API examples** - All 50+ endpoints
2. **Test cases** - Automated validation
3. **Error handling** - Common issues covered
4. **Best practices** - Proper request structure

### Additional Resources
- [POSTMAN_TESTING_GUIDE.md](POSTMAN_TESTING_GUIDE.md) - Detailed guide
- [API_REFERENCE.md](API_REFERENCE.md) - Complete API docs
- [README.md](README.md) - Backend overview

---

## 🔧 Customization

### Change Base URL
1. Click "Environment" dropdown
2. Click "LMS Backend - Development Environment"
3. Edit `base_url` value
4. Save

### Change Default Course
1. Modify `course_id` in environment
2. Update in any request that needs it

### Add New Test Cases
1. Right-click request
2. Select "Add Test"
3. Write JavaScript test
4. Save

### Create Custom Environment
1. Click "Environment" → Create New
2. Add variables
3. Use in requests

---

## 🐛 Troubleshooting

### Collection Won't Import
**Solution:** 
- Ensure JSON is valid (use jsonlint.com)
- Try re-downloading collection
- Clear Postman cache: Settings → Data

### Variables Not Working
**Solution:**
- Verify environment is selected (dropdown, top-right)
- Check variable name spelling
- Run login to set auth_token
- Reload request (Ctrl+R)

### Getting 401 Unauthorized
**Solution:**
- Login again (auth token may be expired)
- Check auth_token is in environment
- Verify Bearer token is in Authorization tab

### Base URL Not Working
**Solution:**
- Verify backend is running
- Check base_url doesn't have trailing slash
- Verify firewall isn't blocking port 8080
- Check if running on different port

---

## 📊 Test Coverage

### Endpoints by Category

| Category | Count | Status |
|----------|-------|--------|
| Authentication | 5 | ✅ |
| Courses | 8 | ✅ |
| Enrollments | 5 | ✅ |
| Progress | 3 | ✅ |
| Quizzes | 3 | ✅ |
| Gamification | 4 | ✅ |
| Dashboard | 1 | ✅ |
| User Management | 2 | ✅ |
| Admin | 3 | ✅ |
| Workflows | 10 | ✅ |
| **TOTAL** | **44** | ✅ |

### HTTP Methods Covered

| Method | Count | Examples |
|--------|-------|----------|
| GET | 25 | Retrieve data |
| POST | 15 | Create/submit |
| PUT | 3 | Update |
| DELETE | 1 | Delete |

### Response Codes Validated

| Code | Purpose | Validated |
|------|---------|-----------|
| 200 | Success | ✅ |
| 201 | Created | ✅ |
| 400 | Bad Request | ✅ |
| 401 | Unauthorized | ✅ |
| 403 | Forbidden | ✅ |
| 404 | Not Found | ✅ |

---

## 💡 Pro Tips

### 1. Bulk Test All Endpoints
- Click collection name
- Click "Run" button
- All requests execute in sequence
- View summary report

### 2. Export Test Results
- After running collection
- Click "..." menu
- Select "Export Results"
- Save as JSON report

### 3. Monitor Network Activity
- View → Show Network Activity
- See all requests/responses
- Check response times
- Profile performance

### 4. Generate Code Snippets
- Right-click request
- Select language (cURL, Python, JavaScript, etc.)
- Copy generated code

### 5. Create Custom Scripts
- Click "Tests" tab
- Write JavaScript with pm library
- Example:
  ```javascript
  pm.test("Status is 200", function() {
      pm.response.to.have.status(200);
  });
  ```

---

## 📋 Checklist Before Testing

- [ ] Postman installed
- [ ] Collection imported
- [ ] Environment imported
- [ ] Environment selected (dropdown)
- [ ] Backend running (`go run cmd/main.go`)
- [ ] PostgreSQL running
- [ ] `.env` file configured
- [ ] Base URL correct
- [ ] Ready to test!

---

## 🎉 Success Indicators

When everything is working:
- ✅ Register request succeeds (200)
- ✅ Login returns token
- ✅ Dashboard loads all data
- ✅ Tests show green checkmarks
- ✅ Variables auto-save
- ✅ Protected endpoints return 200 (not 401)

---

## 📞 Need Help?

1. **Check error message** in response body
2. **Review POSTMAN_TESTING_GUIDE.md** for scenarios
3. **Check API_REFERENCE.md** for endpoint details
4. **Verify backend is running** and database configured
5. **View Postman Console** (bottom-left) for detailed logs

---

**Happy Testing! 🚀**

For detailed testing guide: See [POSTMAN_TESTING_GUIDE.md](POSTMAN_TESTING_GUIDE.md)  
For API documentation: See [API_REFERENCE.md](API_REFERENCE.md)  
For setup instructions: See [README.md](README.md)
