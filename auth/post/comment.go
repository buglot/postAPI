package post

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/buglot/postAPI/lib"
	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
)

type CommentGet struct {
	Url     string `json:"Url" binding:"required"`
	Message string `json:"Message" binding:"required"`
}
type CommentSend struct {
	UserUrl     string `json:"Url"`
	Message     string `json:"Message"`
	IsMyComment bool   `json:"IsMyComment"`
	Avatar      string `json:"Avatar"`
	UserName    string `json:"UserName"`
}

func Comment(ctx *gin.Context) {
	userID := lib.AnyToUInt(ctx.MustGet("userID"))
	var datacomment CommentGet
	if err := ctx.ShouldBindJSON(&datacomment); err != nil {
		ctx.JSON(404, gin.H{"message": "Post not found"})
		return
	}
	fmt.Println(datacomment.Message, datacomment.Url)
	var post orm.Post
	if err := orm.Db.Where("url = ?", datacomment.Url).First(&post).Error; err != nil {
		ctx.JSON(404, gin.H{"message": "Post not found"})
		return
	}

	comment := orm.Comment{
		UserID:  userID,
		PostID:  post.ID,
		Comment: datacomment.Message,
	}

	if err := orm.Db.Create(&comment).Error; err != nil {
		ctx.JSON(500, gin.H{"message": "Cannot comment"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Comment added successfully"})
}
func GetComments(ctx *gin.Context) {
	url := ctx.Query("url")
	userid := lib.AnyToUInt(ctx.MustGet("userID"))
	lParam := ctx.DefaultQuery("l", "5")
	limit, err := strconv.Atoi(lParam)
	if err != nil || limit <= 0 {
		limit = 5
	}
	offset := limit - 5
	if offset < 0 {
		offset = 0
	}
	var post orm.Post
	if err := orm.Db.Where("url = ?", url).First(&post).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "Post not found"})
		return
	}
	var comments1 []orm.Comment
	if err := orm.Db.Preload("User").Where("post_id = ?", post.ID).
		Order("created_at ASC").
		Offset(offset).
		Limit(5).
		Find(&comments1).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot fetch comments"})
		return
	}
	var data []CommentSend
	for _, com := range comments1 {
		var data1 = CommentSend{
			UserUrl:     com.User.Url,
			IsMyComment: com.User.ID == userid,
			Avatar:      com.User.Avatar,
			UserName:    com.User.Username,
			Message:     com.Comment,
		}
		data = append(data, data1)
	}

	ctx.JSON(http.StatusOK, data)
}
