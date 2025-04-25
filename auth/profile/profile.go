package profile

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
)

type ProfileSendURL struct {
	Username    string `json:"Username"`
	Email       string `json:"Email"`
	Url         string `json:"Url"`
	Avatar      string `json:"Avatar"`
	IsMyProfile bool   `json:"IsMyProfile"`
}

func GetProfile(ctx *gin.Context) {
	url := ctx.Query("url")
	userid := ctx.MustGet("userID")
	var data orm.User
	err := orm.Db.Where("url = ?", url).First(&data).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Can't find url"})
		return
	}
	sendr := ProfileSendURL{
		Username:    data.Username,
		Url:         data.Url,
		Email:       data.Email,
		Avatar:      data.Avatar,
		IsMyProfile: false,
	}
	useridint, _ := strconv.Atoi(fmt.Sprintf("%v", userid))
	if data.ID == uint(useridint) {
		sendr.IsMyProfile = true
	}
	ctx.IndentedJSON(http.StatusOK, sendr)
	return
}
