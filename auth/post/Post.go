package post

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreatePostInput struct {
	Message        string   `json:"message"`
	Accessname     string   `json:"accessname" binding:"required"`
	TypeofPostname string   `json:"typeofpostname" binding:"required"`
	Images         []string `json:"images"`
}

func CreatePost(ctx *gin.Context) {
	userid := ctx.MustGet("userID")
	var input CreatePostInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"messeage": err.Error()})
		return
	}
	var access orm.Access
	result := orm.Db.Model(&orm.Access{}).Where("name = ?", input.Accessname).First(&access)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
		return
	}
	var typepost orm.TypeofPost
	result = orm.Db.Model(&orm.TypeofPost{}).Where("name = ?", input.TypeofPostname).First(&typepost)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
		return
	}
	idStr := fmt.Sprintf("%v", userid)
	userIDuintm, _ := strconv.ParseUint(idStr, 10, 32)
	post := orm.Post{
		Url:          strings.Replace(uuid.NewString(), "-", "", 4),
		UserID:       uint(userIDuintm),
		Message:      input.Message,
		AccessID:     access.ID,
		TypeofPostID: typepost.ID,
	}
	for _, imgUrl := range input.Images {
		post.Image = append(post.Image, orm.Image{Url: imgUrl})
	}
	if err := orm.Db.Create(&post).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"Message": "Posted",
	})
}
