package controllers

import (
	"goproject/database"
	"goproject/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetArticles(c *gin.Context) {
    var articles []models.GoArticle
    if err := database.DB.Preload("User").Find(&articles).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve articles"})
        return
    }

    c.JSON(http.StatusOK, articles)
}

func CreateArticle(c *gin.Context) {
    var article models.GoArticle
    userID := c.MustGet("user_id").(uint)

    if err := c.ShouldBindJSON(&article); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    article.UserID = userID
    if err := database.DB.Create(&article).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create article"})
        return
    }

    c.JSON(http.StatusOK, article)
}

func UpdateArticle(c *gin.Context) {
    var article models.GoArticle
    userID := c.MustGet("user_id").(uint)
    articleID := c.Param("id")

    if err := database.DB.Where("id = ? AND user_id = ?", articleID, userID).First(&article).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
        return
    }

    if err := c.ShouldBindJSON(&article); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := database.DB.Save(&article).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update article"})
        return
    }

    c.JSON(http.StatusOK, article)
}

func DeleteArticle(c *gin.Context) {
    var article models.GoArticle
    userID := c.MustGet("user_id").(uint)
    articleID := c.Param("id")

    if err := database.DB.Where("id = ? AND user_id = ?", articleID, userID).First(&article).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
        return
    }

    if err := database.DB.Delete(&article).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete article"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Article deleted"})
}
