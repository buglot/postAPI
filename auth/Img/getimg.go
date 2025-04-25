package img

import (
	"net/http"

	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
)

func GetImg(ctx *gin.Context) {
	filename := ctx.Param("filename")
	fullPath := "./uploads/" + filename
	userID := ctx.MustGet("userID").(uint)
	var img orm.Image
	if err := orm.Db.Where("url = ?", fullPath).First(&img).Error; err != nil {
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
		Where("follower = ? AND followee = ?", userID, post.UserID).
		First(&follow).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "You are not allowed to view this image"})
		return
	}

	ctx.File(fullPath)

}
