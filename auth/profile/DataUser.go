package profile

import (
	"net/http"

	"github.com/buglot/postAPI/auth/post"
	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
)

type Username struct {
	Username    string `json:"Username"`
	Url         string `json:"Url"`
	Avatar      string `json:"Avatar"`
	Email       string `json:"Email"`
}

func Profile(ctx *gin.Context) {
	userid := ctx.MustGet("userID")
	var user orm.User
	err := orm.Db.Where("id = ?", userid).First(&user).Error
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"mesesage": err.Error(),
		})
		return
	}
	data := Username{
		Username: user.Username,
		Url:      user.Url,
		Email:    user.Email,
		Avatar:   user.Avatar,
	}
	ctx.IndentedJSON(http.StatusOK, data)
	return
}

func MyPost(ctx *gin.Context) {
	userid := ctx.MustGet("userID")
	var Post []orm.Post
	err := orm.Db.Preload("Post").Preload("User").Joins("post.user_id = user.id and user.id = ?", userid).Find(&Post).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	var send []post.SentData
	for _, post_i := range Post {
		var images []string
		for _, img := range post_i.Image {
			images = append(images, img.Url) // สมมุติว่า struct Image มีฟิลด์ชื่อ Url
		}
		data := post.SentData{
			Name:         post_i.User.Username,
			Avatar:       post_i.User.Avatar,
			TypeofAccess: post_i.Access.Name,
			TypeofPost:   post_i.TypeofPost.Name,
			Message:      post_i.Message,
			Url:          post_i.Url,
			Date:         post_i.CreatedAt.Format("2006-01-02 15:04:05"),
			Images:       images,
			ErrorMessage: "",
		}
		send = append(send, data)
	}
	ctx.IndentedJSON(http.StatusOK, send)
}
