package profile

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/buglot/postAPI/lib"
	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileSendURL struct {
	Username    string `json:"Username"`
	Email       string `json:"Email"`
	Url         string `json:"Url"`
	Avatar      string `json:"Avatar"`
	IsMyProfile bool   `json:"IsMyProfile"`
	Follow      uint   `json:"Follow"`
	Following   uint   `json:"Following"`
	Followed    bool   `json:"Followed"`
}

func GetProfile(ctx *gin.Context) {
	url := ctx.Query("url")
	userid := ctx.MustGet("userID")
	var data orm.User
	err := orm.Db.Where("url = ?", url).First(&data).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Can't find url"})
		return
	}
	sendr := ProfileSendURL{
		Username:    data.Username,
		Url:         data.Url,
		Email:       data.Email,
		Avatar:      data.Avatar,
		IsMyProfile: false,
	}
	useridint, _ := strconv.Atoi(fmt.Sprintf("%v", userid))
	if data.ID == uint(useridint) {
		sendr.IsMyProfile = true
	}
	var foll_count int64
	err = orm.Db.Model(&orm.Follow{}).Where("follower_id = ?", data.ID).Count(&foll_count).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Can't find url"})
		return
	}
	var follwee_count int64
	err = orm.Db.Model(&orm.Follow{}).Where("followee_id = ?", data.ID).Count(&follwee_count).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Can't find url"})
		return
	}
	sendr.Follow = uint(foll_count)
	sendr.Following = uint(follwee_count)
	sendr.Followed = true
	var me orm.Follow
	err = orm.Db.Model(&orm.Follow{}).Where("followee_id = ? and  follower_id=?", data.ID, userid).First(&me).Error
	if err != nil {
		sendr.Followed = false
	}
	ctx.IndentedJSON(http.StatusOK, sendr)
	return
}

type Post_Follow struct {
	Url string `json:"Url" binding:"required"`
}

func Follow_friend(ctx *gin.Context) {
	userid := lib.AnyToUInt(ctx.MustGet("userID"))
	// Bind the request JSON
	var req Post_Follow
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Find the followee by URL
	var followee orm.User
	if err := orm.Db.Where("url = ?", req.Url).First(&followee).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Prevent self-follow
	if followee.ID == userid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "You cannot follow yourself"})
		return
	}

	// Check if already following
	var existing orm.Follow
	err := orm.Db.Where("follower_id = ? AND followee_id = ?", userid, followee.ID).First(&existing).Error

	if err == nil {
		// Already following -> Unfollow
		if err := orm.Db.Delete(&existing).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Unfollowed"})
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Not following -> Follow
		newFollow := orm.Follow{
			FollowerID: userid,
			FolloweeID: followee.ID,
		}
		if err := orm.Db.Create(&newFollow).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Followed"})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check follow status"})
	}

}
