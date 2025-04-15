package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserData struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Avatar   string `json:"avatar"`
}

func Register(ctx *gin.Context) {
	var data UserData
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := orm.User{
		Email:    data.Email,
		Username: data.Username,
		Avatar:   data.Avatar,
		Password: string(encryptedPassword),
	}
	orm.Db.Create(&user)
	if user.ID > 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "Registered"})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Username and Emaill has registered!"})
	}
	return

}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(ctx *gin.Context) {
	var data UserLogin
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	var RowDataDb orm.User
	orm.Db.Where("username = ? or email = ?", data.Username, data.Email).First(&RowDataDb)
	if RowDataDb.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "You're not registered"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(RowDataDb.Password), []byte(data.Password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Your password are wrong!"})

	} else {
		hmacSampleSecret := os.Getenv("JWT_SECRAT_KEY")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID": RowDataDb.ID,
			"exp":    time.Now().Add(time.Minute * 24).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, _ := token.SignedString([]byte(hmacSampleSecret))
		ctx.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
	}

	return
}
