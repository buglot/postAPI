package post

import (
	"net/http"

	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
)

type SentData struct {
	Name         string   `json:"Name" binding:"required"`
	Avatar       string   `json:"Avatar" binding:"required"`
	TypeofAccess string   `json:"TypeofAccess" binding:"required"`
	TypeofPost   string   `json:"TypeofPost" binding:"required"`
	Message      string   `json:"Message"`
	Date         string   `json:"Date" binding:"required"`
	Images       []string `json:"Images"`
	ErrorMessage string   `json:"errormessage"`
	Url          string   `json:"url"`
}

func GetPost(ctx *gin.Context) {
	userid := ctx.MustGet("userID")
	var posts []orm.Post
	result := orm.Db.
		Preload("User").
		Preload("Access").
		Preload("Image").
		Preload("TypeofPost").
		Joins("JOIN accesses ON accesses.id = posts.access_id").
		Joins("LEFT JOIN follows ON follows.followee_id = posts.user_id AND follows.follower_id = ?", userid).
		Where("accesses.name = ? OR (accesses.name = ? AND follows.id IS NOT NULL)", "public", "follow").
		Find(&posts)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}
	var data []SentData
	for _, post := range posts {
		var images []string
		for _, img := range post.Image {
			images = append(images, img.Url) // สมมุติว่า struct Image มีฟิลด์ชื่อ Url
		}
		converted := SentData{
			Name:         post.User.Username,
			Avatar:       post.User.Avatar,
			TypeofAccess: post.Access.Name,
			TypeofPost:   post.TypeofPost.Name,
			Message:      post.Message,
			Date:         post.CreatedAt.Format("2006-01-02 15:04:05"),
			Images:       images,
			Url: post.Url,
			ErrorMessage: "",
		}
		data = append(data, converted)
	}
	ctx.IndentedJSON(http.StatusOK, data)
}
