package post

import (
	"net/http"

	"github.com/buglot/postAPI/lib"
	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
)

type likesend struct {
	Like      bool  `json:"Like"`
	LikeCount int64 `json:"LikeCount"`
}
type likeget struct {
	Url string `json:"url" binding:"required"`
}

func Like(ctx *gin.Context) {
	userID := lib.AnyToUInt(ctx.MustGet("userID"))
	var url likeget
	ctx.ShouldBindJSON(&url)
	var post orm.Post
	if err := orm.Db.Where("url = ?", url.Url).First(&post).Error; err != nil {
		ctx.JSON(404, gin.H{"message": "Post not found"})
		return
	}
	var existingLike orm.LikePost
	err := orm.Db.Where("user_id = ? AND post_id = ?", userID, post.ID).First(&existingLike).Error
	if err == nil {
		if err := orm.Db.Delete(&existingLike).Error; err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": "Cannot unlike post"})
			return
		}
		var likeCount int64
		orm.Db.Model(&orm.LikePost{}).Where("post_id = ?", post.ID).Count(&likeCount)
		var data = likesend{
			Like:      false,
			LikeCount: likeCount,
		}
		ctx.IndentedJSON(http.StatusOK, data)
		return
	}
	like := orm.LikePost{
		UserID: userID,
		PostID: post.ID,
	}
	if err := orm.Db.Create(&like).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot like post"})
		return
	}
	var likeCount int64
	orm.Db.Model(&orm.LikePost{}).Where("post_id = ?", post.ID).Count(&likeCount)

	var data = likesend{
		Like:      true,
		LikeCount: likeCount,
	}
	ctx.IndentedJSON(http.StatusOK, data)
}
