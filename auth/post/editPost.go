package post

import (
	"fmt"
	"net/http"

	"github.com/buglot/postAPI/lib"
	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
)

type EditPostInput struct {
	CreatePostInput
	Url string `json:"Url" binding:"required"`
}

func EditPost(ctx *gin.Context) {
	userid := lib.AnyToUInt(ctx.MustGet("userID"))
	var dataEdit EditPostInput
	if err := ctx.ShouldBindJSON(&dataEdit); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	fmt.Println(dataEdit.Images)
	var dataDB orm.Post
	if err := orm.Db.Preload("Image").Where("url = ?", dataEdit.Url).First(&dataDB).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if dataDB.UserID != userid {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "you're not premission",
		})
		return
	}
	var access orm.Access
	result := orm.Db.Model(&orm.Access{}).Where("name = ?", dataEdit.Accessname).First(&access)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
		return
	}
	var typepost orm.TypeofPost
	result = orm.Db.Model(&orm.TypeofPost{}).Where("name = ?", dataEdit.TypeofPostname).First(&typepost)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
		return
	}
	dataDB.Message = dataEdit.Message
	dataDB.TypeofPostID = typepost.ID
	dataDB.AccessID = access.ID
	newImageMap := make(map[string]bool)
	for _, img := range dataEdit.Images {

		newImageMap[img] = true
	}
	for _, oldImg := range dataDB.Image {

		fmt.Println(oldImg.Url)
		if !newImageMap[oldImg.Url] {
			fmt.Println(oldImg.Url + " delete")
			orm.Db.Delete(&oldImg)
		}
	}
	err := orm.Db.Save(&dataDB).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}
