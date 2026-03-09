package handlers

import (
	"homework04/models"
	"homework04/services"
	"homework04/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
	jwtSecret   []byte
}

func NewUserHandler(userService *services.UserService, jwtSecret []byte) *UserHandler {
	return &UserHandler{userService: userService, jwtSecret: jwtSecret}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "请求参数错误: "+err.Error())
		return
	}

	if err := h.userService.CreateUser(req.Username, req.Password, req.Email); err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, gin.H{
		"message": "注册成功",
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "请求参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Authenticate(req.Username, req.Password)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, h.jwtSecret)
	if err != nil {
		utils.Error(c, utils.ErrInternalServer.Code, "生成令牌失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
