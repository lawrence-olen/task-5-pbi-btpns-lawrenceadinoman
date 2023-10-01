package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/crocox/final-project/app"
	"github.com/crocox/final-project/database"
	"github.com/crocox/final-project/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Register User
func Register(c *gin.Context) {
	// Get the data off request body
	var register app.UserRegisterInput

	// Request Body JSON
	if c.Bind(&register) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body JSON",
		})

		return
	}

	// Character Length
	if len(register.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Password must be at least 6 Characters",
		})

		return
	}

	// Hash the Password
	hash, err := bcrypt.GenerateFromPassword([]byte(register.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Create the user
	db := database.ConnectionDB()

	user := models.User{Username: register.Username, Email: register.Email, Password: string(hash)}
	result := db.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email have been Registered",
		})

		return
	}

	// Respond
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Register User Successful",
	})
}

// Login User
func Login(c *gin.Context) {
	// Get the data Email and Password off request body
	var login app.UserLoginInput

	if c.Bind(&login) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body JSON",
		})

		return
	}

	// Look up requested user
	var user models.User
	db := database.ConnectionDB()
	db.First(&user, "email = ?", login.Email) // SELECT * FROM users WHERE email = "email";

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})

		return
	}

	// Compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})

		return
	}

	// Generate a JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})

		return
	}

	// Send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*7*1, "", "", false, true)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Token Created",
		"token":   "Check in data cookies",
	})
}

// Validate Email
func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
