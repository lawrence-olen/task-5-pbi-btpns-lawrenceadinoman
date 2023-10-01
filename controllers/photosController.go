package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/crocox/final-project/app"
	"github.com/crocox/final-project/database"
	"github.com/crocox/final-project/models"
	"github.com/gin-gonic/gin"
)

func UploadPhoto(c *gin.Context) {
	// Get the data off request body
	var upload app.UploadRequestPhotos
	user, _ := c.Get("user")
	file, err := c.FormFile("photo_url")

	// fmt.Println(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File Photo Not Found",
		})

		return
	}

	timeNow := time.Now().Unix()
	filename := fmt.Sprintf("%d-%s", timeNow, strings.ReplaceAll(file.Filename, " ", "-"))
	filePath := "http://localhost:8080/api/photos/" + filename

	// Create a new photo record
	upload.Title = c.PostForm("title")
	upload.Caption = c.PostForm("caption")
	upload.PhotoUrl = filePath
	// fmt.Println(filePath)

	photo := models.Photo{
		Title:    upload.Title,
		Caption:  upload.Caption,
		PhotoUrl: upload.PhotoUrl,
		UserID:   user.(float64),
	}

	db := database.ConnectionDB()
	err = c.SaveUploadedFile(file, "./uploads/"+filename)

	// Save the photo to the Database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to Save File to Database",
		})

		return
	}

	if err := db.Create(&photo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Successfully Upload Photo",
	})
}

func GetPhoto(c *gin.Context) {
	var photo []models.Photo
	user, _ := c.Get("user")

	db := database.ConnectionDB()
	result := db.Where("user_id = ?", user).Preload("User").Find(&photo)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch data Photos",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Successfully to fetch data Photos",
		"data":    photo,
	})
}

func UpdatePhoto(c *gin.Context) {
	user, _ := c.Get("user")
	id := c.Param("id")
	file, err := c.FormFile("photo_url")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File Photo Not Found",
		})

		return
	}

	timeNow := time.Now().Unix()
	filename := fmt.Sprintf("%d-%s", timeNow, strings.ReplaceAll(file.Filename, " ", "-"))
	filePath := "http://localhost:8080/api/photos/" + filename

	// Update a new data photo
	update := app.UploadRequestPhotos{
		Title:    c.PostForm("title"),
		Caption:  c.PostForm("caption"),
		PhotoUrl: filePath,
	}

	photo := models.Photo{
		Title:    update.Title,
		Caption:  update.Caption,
		PhotoUrl: update.PhotoUrl,
	}

	db := database.ConnectionDB()
	// result := db.Where("id = ?", id).Where("user_id", user)

	if err := db.Where("id = ?", id).Where("user_id", user).Preload("User").First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Photo Not Found",
		})

		return
	}

	// Menguppdate Photo ke Database
	if err := db.Save(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})

		return
	}

	if db.Model(&photo).Where("user_id = ?", user).Updates(&photo).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Data Photo Not Found",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Successfully Updated Photo",
		"data":    photo,
	})
}

func DeletePhoto(c *gin.Context) {
	var photo models.Photo
	var users models.User
	user, _ := c.Get("user")
	id := c.Param("id")

	db := database.ConnectionDB()

	if err := db.Where("user_id = ?", user).Where("id = ?", id).First(&photo).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Data Photo Not Found",
		})

		return
	}

	//Validasi user tidak dapat menghapus photo yang dibuat user lain
	if users.ID != uint(photo.UserID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "No access to delete photo",
			"data":    nil,
		})

		return
	}

	// Menghapus Photo dari Database
	if err := db.Debug().Delete(&photo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Successfully to Delete data Photos",
	})
}
