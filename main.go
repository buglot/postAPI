package main

import (
	"fmt"
	"time"

	"github.com/buglot/postAPI/auth"
	img "github.com/buglot/postAPI/auth/Img"
	"github.com/buglot/postAPI/auth/post"
	"github.com/buglot/postAPI/auth/profile"
	"github.com/buglot/postAPI/middleware"
	"github.com/buglot/postAPI/orm"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
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
	router.Static("/img/public/", "./img/profile")
	Authen.GET("/img/:filename", img.GetImg)
	Authen.GET("/getPost", post.GetPost)
	Authen.POST("/like", post.Like)
	Authen.POST("/comment", post.Comment)
	Authen.POST("/imgupload", post.Uploads)
	Authen.POST("/Post", post.CreatePost)
	Authen.GET("/PostUrl", post.GetPostURL)
	Authen.GET("/Profile", profile.Profile)
	Authen.GET("/ProfileUrl", profile.GetProfile)
	Authen.GET("/GetPostInProfile", post.GetPostInProfile)
	Authen.GET("/GetComments", post.GetComments)
	Authen.POST("/follow", profile.Follow_friend)
	router.Run("0.0.0.0:8080")
}
