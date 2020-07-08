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

func GetUserByName(c *gin.Context) {

	userName := c.Query("userName")
	if userName == "" {
		utils.Logger.Info("userName must be set")
		EmptyUserResponse(c, "用户查询条件不能为空")
		return
	}

	user, err := models.DefaultUserManager.GetUserByName(userName)
	if err != nil {
		utils.Logger.Error(err)
		EmptyUserResponse(c, "非法用户名")
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "获取用户成功",
		"data":   []*models.User{user},
	})
}

func GetUserById(c *gin.Context) {

	userIdStr := c.Query("userId")
	if userIdStr == "" {
		utils.Logger.Info("userId must be set")
		EmptyUsersResponse(c, "用户查询条件不能为空")
		return
	}

	var (
		userId int
		err    error
		user   *models.User
	)
	userId, err = strconv.Atoi(userIdStr)
	if err != nil {
		utils.Logger.Error(err)
		EmptyUserResponse(c, "非法用户ID")
		return
	}
	user, err = models.DefaultUserManager.GetUserById(userId)
	if err != nil {
		utils.Logger.Error(err)
		EmptyUserResponse(c, "获取用户失败")
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "获取用户成功",
		"data":   []*models.User{user},
	})

}

func GetUserList(c *gin.Context) {
	var (
		err  error
		p, s int
	)
	defaultPageNum := utils.GlobalConfig.GetInt("server.default_page_num")
	defaultPageSize := utils.GlobalConfig.GetInt("server.default_page_size")
	page := c.Query("page")
	size := c.Query("size")

	p, err = strconv.Atoi(page)
	if err != nil {
		p = defaultPageNum
	}
	s, err = strconv.Atoi(size)
	if err != nil {
		s = defaultPageSize
	}

	total, users, pagination, err := models.DefaultUserManager.GetUserList(p, s)
	if err != nil {
		utils.Logger.Error(err)
		EmptyUsersResponse(c, "获取用户列表失败")
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

func GetUserDetailById(c *gin.Context) {
	userId := utils.GetUserIdFromUrl(c)
	if userId == utils.InvalidUid {
		EmptyUserDetailsResponse(c, "获取用户详情失败")
		return
	}
	userDetails, err := models.DefaultUserManager.GetUserDetailById(userId)
	if err != nil {
		utils.Logger.Error(err)
		EmptyUserDetailsResponse(c, "获取用户详情失败")
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "获取用户详情成功",
		"data":   userDetails,
	})
}

func CreateUser(c *gin.Context) {
	urf := &forms.UserRegisterForm{}
	err := c.ShouldBindJSON(urf)
	if err != nil {
		utils.Logger.Error(err)
		BadResponse(c, "非法表单数据")
		return
	}

	_, err = models.DefaultUserManager.GetUserByName(urf.UserNmae)
	if err == nil {
		utils.Logger.Info("用户已经存在")
		BadResponse(c, "用户已经存在")
		return
	}
	uid, err := models.DefaultUserManager.CreateUser(urf)
	if err != nil {
		utils.Logger.Error(err)
		BadResponse(c, "用户创建失败")
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
		BadResponse(c, "非法表单数据")
		return
	}

	_, err = models.DefaultUserManager.GetUserById(userId)
	if err != nil {
		utils.Logger.Info(err)
		BadResponse(c, "用户不存在")
		return
	}

	err = models.DefaultUserManager.UpdateDetailById(userId, uuf)
	if err != nil {
		utils.Logger.Error(err)
		BadResponse(c, "用户详情更新失败")
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
		BadResponse(c, "非法表单数据")
		return
	}

	_, err = models.DefaultUserManager.GetUserById(userId)
	if err != nil {
		utils.Logger.Error(err)
		BadResponse(c, "用户不存在")
		return
	}

	err = models.DefaultUserManager.UpdateUserTypeById(userId, utuf)
	if err != nil {
		utils.Logger.Error(err)
		BadResponse(c, "用户类型更新失败")
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
		BadResponse(c, "非法表单数据")
		return
	}

	_, err = models.DefaultUserManager.GetUserById(userId)
	if err != nil {
		utils.Logger.Error(err)
		BadResponse(c, "用户不存在")
		return
	}

	err = models.DefaultUserManager.UpdateUserStatusById(userId, usuf)
	if err != nil {
		BadResponse(c, "用户状态更新失败")
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
		BadResponse(c, "非法用户id")
		return
	}
	upuf := &forms.PasswordUpdateForm{}
	err := c.ShouldBindJSON(upuf)
	if err != nil {
		utils.Logger.Error(err)
		BadResponse(c, "非法表单数据")
		return
	}

	_, err = models.DefaultUserManager.GetUserById(userId)
	if err != nil {
		utils.Logger.Error(err)
		BadResponse(c, "用户不存在")
		return
	}

	err = models.DefaultUserManager.UpdatePasswordById(userId, upuf)
	if err != nil {
		utils.Logger.Error(err)
		BadResponse(c, "改密失败")
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
		BadResponse(c, "解析上传文件失败")
		return
	}

	contentType := file.Header.Get("Content-Type")
	cs := strings.Split(contentType, "/")
	var suffix = ""
	if len(cs) > 1 {
		suffix = cs[1]
	} else {
		BadResponse(c, "获取文件MIMETYPE失败")
		return
	}

	switch suffix {
	case "jpeg":
		suffix = ".jpg"
	case "png":
		suffix = ".png"
	default:
		BadResponse(c, "仅支持 jpeg/png")
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

	err = models.DefaultUserManager.UpdateAvatarById(imgUrl, userId)
	if err != nil {
		utils.Logger.Error(err)
		BadResponse(c, "文件路径写入数据库失败")
		return
	}

	utils.Logger.Info(fmt.Sprintf("upload file success [userId:%d]", userId))
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "upload file success",
		"data":   gin.H{"fileUrl": imgUrl,},
	})

}
