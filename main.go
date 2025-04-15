package main

import (
	"log"

	"github.com/buglot/postAPI/auth"
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
	router.Use(cors.Default())
	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)
	router.Run("localhost:8080")
}
