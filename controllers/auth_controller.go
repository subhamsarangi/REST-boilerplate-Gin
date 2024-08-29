package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"

	"goproject/config"
	"goproject/database"
	"goproject/models"
	"goproject/utils"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Register(c *gin.Context) {
	var input models.GoUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(input); err != nil {
		validationError, errorStatus := utils.UserDataValidationError(err)
		c.JSON(errorStatus, gin.H{"error": validationError})
		return
	}

	if err := input.HashPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	if err := database.DB.Create(&input).Error; err != nil {
		insertError, errorStatus := utils.UserDataInsertError(err)
		c.JSON(errorStatus, gin.H{"error": insertError})
		return
	}

	c.JSON(http.StatusOK, input)
}

func Login(c *gin.Context) {
	var input models.GoUser
	var user models.GoUser

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryField := "username"
	if strings.Contains(input.Username, "@") {
		queryField = "email"
	}

	if err := database.DB.Where(fmt.Sprintf("%s = ?", queryField), input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/email or password"})
		return
	}

	if !user.CheckPassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"type":    "access",
		"iss":     "goblog",
		"aud":     "googlers",
		"iat":     time.Now().Unix(),
		"sub":     user.Username,
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Profile(c *gin.Context) {
	var user models.GoUser
	userID := c.MustGet("user_id").(uint)

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
