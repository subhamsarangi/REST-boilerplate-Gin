package routes

import (
	"goproject/controllers"
	"goproject/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    auth := r.Group("/auth")
    {
        auth.POST("/register", controllers.Register)
        auth.POST("/login", controllers.Login)
        auth.GET("/profile", middleware.JWTMiddleware(), controllers.Profile)
    }

    articles := r.Group("/")
    {
        articles.GET("/articles", controllers.GetArticles)
        articles.POST("/article", middleware.JWTMiddleware(), controllers.CreateArticle)
        articles.PUT("/article/:id/update", middleware.JWTMiddleware(), controllers.UpdateArticle)
        articles.DELETE("/article/:id/delete", middleware.JWTMiddleware(), controllers.DeleteArticle)
    }

    return r
}
