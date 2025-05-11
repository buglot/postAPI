package post

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/buglot/postAPI/lib"
	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
)

type SentData struct {
	Name         string   `json:"Name" binding:"required"`
	Avatar       string   `json:"Avatar" binding:"required"`
	UserUrl      string   `json:"UserUrl"`
	TypeofAccess string   `json:"TypeofAccess" binding:"required"`
	TypeofPost   string   `json:"TypeofPost" binding:"required"`
	Message      string   `json:"Message"`
	Date         string   `json:"Date" binding:"required"`
	Images       []string `json:"Images" binding:"required"`
	ErrorMessage string   `json:"errormessage"`
	Url          string   `json:"Url"`
	IsMyPost     bool     `json:"IsMyPost"`
	IntLike      int64    `json:"IntLike"`
	Liked        bool     `json:"Liked"`
}

func GetPost(ctx *gin.Context) {
	var userid = ctx.MustGet("userID")
	useridint, _ := strconv.Atoi(fmt.Sprintf("%v", userid))
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
			images = append(images, img.Url)
		}
		var existingLike orm.LikePost
		err := orm.Db.Where("user_id = ? AND post_id = ?", useridint, post.ID).First(&existingLike).Error
		var likeCount int64
		orm.Db.Model(&orm.LikePost{}).Where("post_id = ?", post.ID).Count(&likeCount)
		converted := SentData{
			Name:         post.User.Username,
			UserUrl:      post.User.Url,
			Avatar:       post.User.Avatar,
			TypeofAccess: post.Access.Name,
			TypeofPost:   post.TypeofPost.Name,
			Message:      post.Message,
			Date:         post.CreatedAt.Format("2006-01-02 15:04:05"),
			Images:       images,
			Url:          post.Url,
			ErrorMessage: "",
			IsMyPost:     post.User.ID == uint(useridint),
			IntLike:      likeCount,
		}
		if err == nil {
			converted.Liked = true
		} else {
			converted.Liked = false
		}
		data = append(data, converted)
	}
	ctx.IndentedJSON(http.StatusOK, data)
}

func GetPostInProfile(ctx *gin.Context) {
	userid := ctx.MustGet("userID")
	url := ctx.Query("url")
	var user orm.User
	err := orm.Db.Where("url = ?", url).First(&user).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	var posts []orm.Post
	useridint, _ := strconv.Atoi(fmt.Sprintf("%v", userid))
	if user.ID == uint(useridint) { //if same userid
		result := orm.Db.
			Joins("JOIN accesses ON accesses.id = posts.access_id").
			Joins("JOIN typeof_posts ON typeof_posts.id = posts.typeof_post_id").
			Joins("JOIN users ON users.id = posts.user_id").
			Preload("User").
			Preload("Access").
			Preload("Image").
			Where("users.url = ?", url).
			Find(&posts)
		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
			return
		}

	} else {
		result := orm.Db.
			Preload("User").
			Preload("Access").
			Preload("Image").
			Preload("TypeofPost").
			Joins("JOIN accesses ON accesses.id = posts.access_id").
			Joins("LEFT JOIN follows ON follows.followee_id = posts.user_id AND follows.follower_id = ?", userid).
			Joins("JOIN users ON users.id = posts.user_id").
			Where("users.url = ? AND accesses.name = ? OR (accesses.name = ? AND follows.id IS NOT NULL)", url, "public", "follow").
			Find(&posts)
		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
			return
		}
	}
	var data []SentData
	for _, post := range posts {
		var images []string
		for _, img := range post.Image {
			images = append(images, img.Url) // สมมุติว่า struct Image มีฟิลด์ชื่อ Url
		}
		var likeCount int64

		orm.Db.Model(&orm.LikePost{}).Where("post_id = ?", post.ID).Count(&likeCount)
		var existingLike orm.LikePost
		err := orm.Db.Where("user_id = ? AND post_id = ?", useridint, post.ID).First(&existingLike).Error
		converted := SentData{
			Name:         post.User.Username,
			UserUrl:      post.User.Url,
			Avatar:       post.User.Avatar,
			TypeofAccess: post.Access.Name,
			TypeofPost:   post.TypeofPost.Name,
			Message:      post.Message,
			Date:         post.CreatedAt.Format("2006-01-02 15:04:05"),
			Images:       images,
			Url:          post.Url,
			ErrorMessage: "",
			IsMyPost:     post.User.ID == uint(useridint),
			IntLike:      likeCount,
		}
		if err == nil {
			converted.Liked = true
		} else {
			converted.Liked = false
		}
		data = append(data, converted)
	}
	ctx.IndentedJSON(http.StatusOK, data)

}

func GetPostURL(ctx *gin.Context) {
	url := ctx.Query("url")
	userid := lib.AnyToUInt(ctx.MustGet("userID"))

	var post orm.Post
	orm.Db.Preload("User").
		Preload("Access").
		Preload("TypeofPost").
		Preload("Image").
		Joins("JOIN users on users.id = posts.user_id").
		Joins("JOIN accesses on accesses.id = posts.access_id").
		Where("posts.url = ?", url).First(&post)
	var images []string
	for _, img := range post.Image {
		images = append(images, img.Url) // สมมุติว่า struct Image มีฟิลด์ชื่อ Url
	}
	var likeCount int64
	orm.Db.Model(&orm.LikePost{}).Where("post_id = ?", post.ID).Count(&likeCount)
	data := SentData{
		Name:         post.User.Username,
		UserUrl:      post.User.Url,
		Avatar:       post.User.Avatar,
		TypeofAccess: post.Access.Name,
		TypeofPost:   post.TypeofPost.Name,
		Message:      post.Message,
		Date:         post.CreatedAt.Format("2006-01-02 15:04:05"),
		Images:       images,
		Url:          post.Url,
		ErrorMessage: "",
		IsMyPost:     post.User.ID == userid,
		IntLike:      likeCount,
	}
	var existingLike orm.LikePost
	err := orm.Db.Where("user_id = ? AND post_id = ?", userid, post.ID).First(&existingLike).Error
	if err == nil {
		data.Liked = true
	} else {
		data.Liked = false
	}
	if userid == post.UserID {
		ctx.IndentedJSON(http.StatusOK, data)
		return
	}
	if post.Access.Name == "public" {
		ctx.IndentedJSON(http.StatusOK, data)
		return
	}
	var follow orm.Follow
	if err := orm.Db.
		Where("follower_id = ? AND followee_id = ?", userid, post.UserID).
		First(&follow).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "You are not allowed to view this image"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, data)
}
