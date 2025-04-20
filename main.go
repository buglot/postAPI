package main

import (
	"log"
	"time"

	"github.com/buglot/postAPI/auth"
	"github.com/buglot/postAPI/auth/post"
	"github.com/buglot/postAPI/middleware"
	"github.com/buglot/postAPI/orm"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := gin.Default()
	orm.InitDB()
	orm.RoleDefault()
	orm.AccessAndTypePostDefault()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)
	Authen := router.Group("/auth", middleware.Auth())
	Authen.Static("/img", "./uploads")
	Authen.GET("/getPost", post.GetPost)
	Authen.POST("/imgupload", post.Uploads)
	Authen.POST("/Post")
	router.Run("localhost:8080")
}
