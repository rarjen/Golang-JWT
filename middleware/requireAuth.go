package middleware

import (
	"fmt"
	"go-jwt/initializer"
	"go-jwt/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// func RequireAuth(c *gin.Context) {
// 	fmt.Println("In middleware")

// 	tokenString, err := c.Cookie("Authorization")
// 	if err != nil {
// 		c.AbortWithStatus(http.StatusUnauthorized)
// 	}

// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(os.Getenv("SECRET")), nil
// 	})

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		if float64(time.Now().Unix()) > claims["exp"].(float64) {
// 			c.AbortWithStatus(http.StatusUnauthorized)
// 		}

// 		var user models.User

// 		initializer.DB.First(&user, claims["id"])

// 		if user.ID == 0 {
// 			c.AbortWithStatus(http.StatusUnauthorized)
// 		}

// 		c.Set("user", user)

// 		c.Next()
// 	} else {
// 		c.AbortWithStatus(http.StatusUnauthorized)
// 	}
// }

func RequireAuth(c *gin.Context) {
	fmt.Println("In middleware")

	// Get the Authorization header
	authHeader := c.GetHeader("Authorization")

	// Check if the Authorization header is properly formatted (starts with "Bearer")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Extract the token from the header (remove "Bearer " prefix)
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse the token
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	// Check token validity and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check token expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Find the user in the database
		var user models.User
		initializer.DB.First(&user, claims["id"])

		// If user does not exist, abort
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Attach user to the context
		c.Set("user", user)

		// Continue to the next handler
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
