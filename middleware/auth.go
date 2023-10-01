package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/crocox/final-project/database"
	"github.com/crocox/final-project/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Auth(c *gin.Context) {
	// fmt.Println("In Middleware")
	// Get the cookie off request
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Decode/Validate it
	// Parse takes the token string and a function for looking up the key. The latter is especially
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Find the user with token
		var user models.User
		users := database.ConnectionDB()
		users.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Attach to request
		c.Set("user", claims["sub"])

		// Continue
		c.Next()

		// fmt.Println(claims["foo"], claims["nbf"])
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
