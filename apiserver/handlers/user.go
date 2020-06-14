package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gogo-cmdb/apiserver/forms"
	"gogo-cmdb/apiserver/models"
	"gogo-cmdb/apiserver/utils"
	"path/filepath"
	"strconv"
	"strings"
)

func GetUser(c *gin.Context) {

	userIdStr := c.DefaultQuery("userId", utils.NoneUserId)
	userName := c.DefaultQuery("userName", utils.UnknownUserName)
	if userIdStr == utils.NoneUserId && userName == utils.UnknownUserName {
		utils.Logger.Info("userId or userName must be set")
		utils.BadResponse(c, "用户查询条件不能为空")
		return
	}

	var (
		userId int
		err    error
	)
	if userIdStr != utils.NoneUserId {
		userId, err = strconv.Atoi(userIdStr)
		if err != nil {
			utils.Logger.Error(err)
			utils.BadResponse(c, "非法用户ID")
			return
		}

		user, err := models.DefalutUserManager.GetUserById(userId)
		if err != nil {
			utils.Logger.Error(err)
			utils.BadResponse(c, "获取用户失败")
			return
		}
		c.JSON(200, gin.H{
			"status": "ok",
			"msg":    "获取用户成功",
			"data":   user,
		})
	}

	if userName != utils.UnknownUserName {
		user, err := models.DefalutUserManager.GetUserByName(userName)
		if err != nil {
			utils.Logger.Error(err)
			utils.BadResponse(c, "非法用户名")
			return
		}
		c.JSON(200, gin.H{
			"status": "ok",
			"msg":    "获取用户成功",
			"data":   user,
		})
	}
}

func GetUserDetailById(c *gin.Context) {
	userId := utils.GetUserIdFromUrl(c)
	if userId == utils.InvalidUid {
		return
	}
	userDetails, err := models.DefalutUserManager.GetUserDetailById(userId)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "获取用户详情失败")
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "获取用户详情成功",
		"data":   userDetails,
	})

}

func GetUserList(c *gin.Context) {
	var p, s int
	var err error
	deafultPageNum := utils.GlobalConfig.GetInt("server.default_page_num")
	deafultPageSize := utils.GlobalConfig.GetInt("server.default_page_size")
	page := c.DefaultQuery("page", "none")
	size := c.DefaultQuery("size", "none")

	p, err = strconv.Atoi(page)
	if err != nil {
		p = deafultPageNum
	}
	s, err = strconv.Atoi(size)
	if err != nil {
		s = deafultPageSize
	}

	total, users, pagination, err := models.DefalutUserManager.GetUserList(p, s)
	if err != nil {
		utils.Logger.Error(err)
		c.JSON(200, gin.H{
			"status": "ok",
			"msg":    "获取用户列表失败",
			"data": gin.H{
				"total":          0,
				"users":          nil,
				"currentPageNum": -1,
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "获取用户列表成功",
		"data": gin.H{
			"total":          total,
			"users":          users,
			"currentPageNum": pagination.CurrentPageNum,
		},
	})

}

func CreateUser(c *gin.Context) {
	urf := &forms.UserRegisterForm{}
	err := c.ShouldBindJSON(urf)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "非法表单数据")
		return
	}

	_, err = models.DefalutUserManager.GetUserByName(urf.UserNmae)
	if err == nil {
		utils.Logger.Info("用户已经存在")
		utils.BadResponse(c, "用户已经存在")
		return
	}
	uid, err := models.DefalutUserManager.CreateUser(urf)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "用户创建失败")
		return
	}

	utils.Logger.Info(fmt.Sprintf("Create new user success [userId:%d]", uid))

	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "用户创建成功",
		"data":   gin.H{"userId": uid,},
	})
}

func UpdateDetailById(c *gin.Context) {
	userId := utils.GetUserIdFromUrl(c)
	if userId == utils.InvalidUid {
		return
	}
	uuf := &forms.DetailUpdateForm{}
	err := c.ShouldBindJSON(uuf)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "非法表单数据")
		return
	}

	_, err = models.DefalutUserManager.GetUserById(userId)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "用户不存在")
		return
	}

	err = models.DefalutUserManager.UpdateDetailById(userId, uuf)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "用户详情更新失败")
		return
	}

	utils.Logger.Info(fmt.Sprintf("Update user details success [userId:%d]", userId))
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "用户详情更新成功",
		"data":   gin.H{},
	})
}

func UpdateUserTypeById(c *gin.Context) {
	userId := utils.GetUserIdFromUrl(c)
	if userId == utils.InvalidUid {
		return
	}
	utuf := &forms.UserTypeUpdateForm{}
	err := c.ShouldBindJSON(utuf)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "非法表单数据")
		return
	}

	_, err = models.DefalutUserManager.GetUserById(userId)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "用户不存在")
		return
	}

	err = models.DefalutUserManager.UpdateUserTypeById(userId, utuf)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "用户类型更新失败")
		return
	}

	utils.Logger.Info(fmt.Sprintf("Update user Type success [userId:%d userType:%s]", userId, utils.UserTypeToStr(utils.IntToUserType(utuf.UserType))))
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "用户类型更新成功",
		"data":   gin.H{},
	})
}

func UpdateUserStatusById(c *gin.Context) {
	userId := utils.GetUserIdFromUrl(c)
	if userId == utils.InvalidUid {
		return
	}
	usuf := &forms.UserStatusUpdateForm{}
	err := c.ShouldBindJSON(usuf)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "非法表单数据")
		return
	}

	_, err = models.DefalutUserManager.GetUserById(userId)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "用户不存在")
		return
	}

	err = models.DefalutUserManager.UpdateUserStatusById(userId, usuf)
	if err != nil {
		utils.BadResponse(c, "用户状态更新失败")
		return
	}

	utils.Logger.Info(fmt.Sprintf("Update user status success [userId:%d userType:%s]", userId, utils.UserStatusToStr(utils.IntToUserStatus(usuf.UserStatus))))

	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "用户状态更新成功",
		"data":   gin.H{},
	})
}

func UpdatePasswordById(c *gin.Context) {
	userId := utils.GetUserIdFromUrl(c)
	if userId == utils.InvalidUid {
		return
	}
	upuf := &forms.PasswordUpdateForm{}
	err := c.ShouldBindJSON(upuf)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "非法表单数据")
		return
	}

	_, err = models.DefalutUserManager.GetUserById(userId)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "用户不存在")
		return
	}

	err = models.DefalutUserManager.UpdatePasswordById(userId, upuf)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "改密失败")
		return
	}

	utils.Logger.Info(fmt.Sprintf("Change password success [userId:%d]", userId))

	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "改密成功",
		"data":   gin.H{},
	})

}

var host = utils.GlobalConfig.GetString("server.host")
var port = utils.GlobalConfig.GetString("server.port")
var userStaticUrl = utils.GlobalConfig.GetString("filesystem.user_static_url")
var userStaticPath = utils.GlobalConfig.GetString("filesystem.user_static_dir_name")

func UploadAvatar(c *gin.Context) {
	userId := utils.GetUserIdFromUrl(c)
	if userId == utils.InvalidUid {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("parse form file failed, %s", err))
		utils.BadResponse(c, "解析上传文件失败")
		return
	}

	contentType := file.Header.Get("Content-Type")
	cs := strings.Split(contentType, "/")
	var suffix = ""
	if len(cs) > 1 {
		suffix = cs[1]
	} else {
		utils.BadResponse(c, "获取文件MIMETYPE失败")
		return
	}

	switch suffix {
	case "jpeg":
		suffix = ".jpg"
	case "png":
		suffix = ".png"
	default:
		utils.BadResponse(c, "仅支持 jpeg/png")
		return
	}

	upAbsDir := filepath.Join(utils.BaseDir, "upload", userStaticPath, fmt.Sprintf("%d", userId))

	// map[Content-Disposition:[form-data; name="file"; filename="head-img-01.jpg"] Content-Type:[image/jpeg]]
	// fmt.Println(file.Header)

	//newFileName := fmt.Sprintf("%s%s", utils.GenUuid(), suffix)
	newFileName, err := utils.SaveSingleFile(c, file, upAbsDir, suffix)
	if err != nil {
		utils.Logger.Error(err)
		return
	}

	imgUrl := fmt.Sprintf("http://%s:%s/img%s", host, port, fmt.Sprintf("%s/%s/%s", userStaticUrl, fmt.Sprintf("%d", userId), newFileName))

	err = models.DefalutUserManager.UpdateAvatarById(imgUrl, userId)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "文件路径写入数据库失败")
		return
	}

	utils.Logger.Info(fmt.Sprintf("upload file success [userId:%d]", userId))
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "upload file success",
		"data": gin.H{
			"fileUrl": imgUrl,
		},
	})

}
