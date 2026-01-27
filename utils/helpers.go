package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash of the password
// Parameters:
//   - password: plain text password to hash
//
// Returns: hashed password and error if any
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword compares a plain text password with a hash
// Parameters:
//   - hashedPassword: the hash to compare against
//   - password: the plain text password to check
//
// Returns: true if password matches, false otherwise
func ComparePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// ValidateEmail validates email format
// Parameters:
//   - email: email address to validate
//
// Returns: true if valid, false otherwise
func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// GenerateCertificateNumber generates a unique certificate number
// Parameters:
//   - userID: unique user identifier
//   - courseID: unique course identifier
//
// Returns: certificate number as string
func GenerateCertificateNumber(userID, courseID string) string {
	data := fmt.Sprintf("%s-%s-%d", userID, courseID, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:16]
}

// CalculateBadgeLevel calculates the badge level based on completed courses
// Parameters:
//   - completedCourses: number of completed courses
//
// Returns: badge level name
func CalculateBadgeLevel(completedCourses int) string {
	switch {
	case completedCourses < 5:
		return "Bronze"
	case completedCourses < 15:
		return "Silver"
	case completedCourses < 30:
		return "Gold"
	case completedCourses < 50:
		return "Platinum"
	default:
		return "Diamond"
	}
}

// CalculateCoinsEarned calculates coins earned for course completion
// Parameters:
//   - score: quiz/assessment score
//   - duration: course duration in minutes
//
// Returns: coins earned as int64
func CalculateCoinsEarned(score int, duration int) int64 {
	baseCoins := int64(duration / 10) // 1 coin per 10 minutes
	if score >= 90 {
		baseCoins += 50 // Bonus for high score
	} else if score >= 80 {
		baseCoins += 25
	} else if score >= 70 {
		baseCoins += 10
	}
	return baseCoins
}

// GetProgressColor returns color indicator for progress percentage
// Parameters:
//   - progress: progress percentage (0-100)
//
// Returns: color code as string
func GetProgressColor(progress int) string {
	switch {
	case progress <= 30:
		return "red"
	case progress <= 70:
		return "orange"
	case progress < 100:
		return "yellow"
	default:
		return "green"
	}
}

// IsValidURL checks if a string is a valid URL
// Parameters:
//   - str: string to validate
//
// Returns: true if valid URL, false otherwise
func IsValidURL(str string) bool {
	re := regexp.MustCompile(`^https?://[^\s]+$`)
	return re.MatchString(str)
}

// GetInitials returns initials from first and last name
// Parameters:
//   - firstName: user first name
//   - lastName: user last name
//
// Returns: initials as string
func GetInitials(firstName, lastName string) string {
	if len(firstName) > 0 && len(lastName) > 0 {
		return string(firstName[0]) + string(lastName[0])
	} else if len(firstName) > 0 {
		return string(firstName[0])
	}
	return ""
}

// RoundFloat rounds a float to specified decimal places
// Parameters:
//   - value: float value to round
//   - decimals: number of decimal places
//
// Returns: rounded float64 value
func RoundFloat(value float64, decimals int) float64 {
	ratio := float64(1)
	for i := 0; i < decimals; i++ {
		ratio *= 10
	}
	return float64(int64(value*ratio)) / ratio
}

// CalculateCompletionPercentage calculates course completion percentage
// Parameters:
//   - completedLessons: number of completed lessons
//   - totalLessons: total number of lessons
//
// Returns: completion percentage as integer
func CalculateCompletionPercentage(completedLessons, totalLessons int) int {
	if totalLessons == 0 {
		return 0
	}
	return (completedLessons * 100) / totalLessons
}

// GetBadgeRequirements returns requirements for each badge level
// Parameters:
//   - badgeLevel: badge level name
//
// Returns: number of courses required for that level
func GetBadgeRequirements(badgeLevel string) int {
	switch badgeLevel {
	case "Bronze":
		return 0
	case "Silver":
		return 5
	case "Gold":
		return 15
	case "Platinum":
		return 30
	case "Diamond":
		return 50
	default:
		return 0
	}
}
