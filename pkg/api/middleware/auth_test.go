package middleware

import (
	"juno/pkg/api/auth"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// Helper function to generate a JWT token for testing
func generateValidJWT(userID uuid.UUID, email string) string {
	claims := jwt.MapClaims{
		"id":    userID.String(),
		"email": email,
		"name":  "Test User",
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
	return tokenString
}

// TestAuthMiddleware tests the AuthMiddleware for various cases
func TestAuthMiddleware(t *testing.T) {
	// Set up the SECRET environment variable
	os.Setenv("SECRET", "mysecretkey")

	// Set up Gin
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())

	// Define a protected route
	r.GET("/protected", func(c *gin.Context) {
		u, ok := auth.UserFromContext(c.Request.Context())

		if !ok {
			c.JSON(http.StatusOK, gin.H{"message": "No user in context"})
			return
		}

		if u != nil {
			c.JSON(http.StatusOK, gin.H{"message": "Access granted", "user_id": u.ID.String()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "No user in context"})
		}
	})

	// Sub-test: Missing Authorization header
	t.Run("Missing Authorization Header", func(t *testing.T) {
		// Create an HTTP request without Authorization header
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Check if the response is 401 Unauthorized
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %d, but got %d", http.StatusUnauthorized, w.Code)
		}

		// Check if the response body contains the expected error message
		expectedBody := `{"error":"Authorization header required"}`
		if w.Body.String() != expectedBody {
			t.Errorf("Expected body %s, but got %s", expectedBody, w.Body.String())
		}
	})

	// Sub-test: Invalid token
	t.Run("Invalid Token", func(t *testing.T) {
		// Create an HTTP request with an invalid token
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "InvalidToken")
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Check if the response is 401 Unauthorized
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %d, but got %d", http.StatusUnauthorized, w.Code)
		}

		// Check if the response body contains the expected error message
		expectedBody := `{"error":"Invalid token"}`
		if w.Body.String() != expectedBody {
			t.Errorf("Expected body %s, but got %s", expectedBody, w.Body.String())
		}
	})

	// Sub-test: Valid token
	t.Run("Valid Token", func(t *testing.T) {
		// Generate a valid user and JWT token
		userID := uuid.New()
		token := generateValidJWT(userID, "test@example.com")

		// Create an HTTP request with a valid token
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", token)
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Check if the response is 200 OK
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
		}

		// Check if the response body contains the user ID
		expectedBody := `{"message":"Access granted","user_id":"` + userID.String() + `"}`
		if w.Body.String() != expectedBody {
			t.Errorf("Expected body %s, but got %s", expectedBody, w.Body.String())
		}
	})
}
