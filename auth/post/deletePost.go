package post

import (
	"net/http"

	"github.com/buglot/postAPI/lib"
	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
)

type DeletePostInput struct {
	Url string `json:"Url" binding:"required"`
}

func DeletePost(ctx *gin.Context) {
	userid := lib.AnyToUInt(ctx.MustGet("userID"))
	var data DeletePostInput
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}
	var post orm.Post
	if err := orm.Db.Where("url = ? and user_id = ?", data.Url, userid).Delete(&post).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "deleted",
	})
}
