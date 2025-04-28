package img

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
)

func GetImg(ctx *gin.Context) {
	filename := ctx.Param("filename")
	fullPath := "./uploads/" + filename
	userID := ctx.MustGet("userID")
	var img orm.Image
	if err := orm.Db.Where("url = ?", filename).First(&img).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Image not found"})
		return
	}
	var post orm.Post
	if err := orm.Db.
		Where("id = ?", img.PostID).
		Preload("Access").
		First(&post).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}
	uesridint, _ := strconv.Atoi(fmt.Sprintf("%v", userID))
	if post.UserID == uint(uesridint) {
		ctx.File(fullPath)
		return
	}
	if post.Access.Name == "public" {
		ctx.File(fullPath)
		return
	}
	if userID == post.UserID {
		ctx.File(fullPath)
		return
	}
	var follow orm.Follow
	if err := orm.Db.
		Where("follower_id = ? AND followee_id = ?", userID, post.UserID).
		First(&follow).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "You are not allowed to view this image"})
		return
	}

	ctx.File(fullPath)

}

