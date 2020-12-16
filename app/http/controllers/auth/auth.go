package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaobingge/dbger/app/utils"
	"net/http"
)

type UserInfo struct {
	Username string `form:"username"`
	Password  string `form:"password"`
}

func Login(c *gin.Context) {
	// 用户发送用户名和密码过来
	var user UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "无效的参数",
		})
		return
	}
	// 校验用户名和密码是否正确
	if user.Username == "bingge" && user.Password == "123456" {
		// 生成Token
		tokenString, _ := utils.GenerateToken(user.Username)
		c.JSON(http.StatusOK, gin.H{
			"code": 1000,
			"msg":  "success",
			"data": gin.H{"token": tokenString},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2002,
		"msg":  "鉴权失败",
	})
	return
}

func Register(c *gin.Context)  {
	var user UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "无效的参数",
		})
		return
	}
	passwordbyte, err := utils.GeneratePassword(user.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2002,
			"msg":  "加密出错了",
		})
	}else{
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "success",
			"data": gin.H{"password": passwordbyte},
		})
	}
	return
}




