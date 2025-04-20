package post

import (
	"net/http"

	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
)

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
	ctx.IndentedJSON(http.StatusOK, posts)
}
