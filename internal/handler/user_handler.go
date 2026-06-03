package handler

import (
	"net/http"

	"github.com/gowild/server/internal/middleware"
	"github.com/gowild/server/internal/model"
	"github.com/gowild/server/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc       *service.UserService
	jwtSecret string
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc, jwtSecret: svc.GetJWTSecret()}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	user, err := h.svc.Register(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "注册失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "注册成功", "user": model.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	user, err := h.svc.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "登录失败", "detail": err.Error()})
		return
	}

	// 生成访问令牌
	token, err := middleware.GenerateToken(user.ID, user.Username, h.jwtSecret, 24*3600*1e9) // 简化: 24小时
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	// 生成刷新令牌
	refreshToken, _ := middleware.GenerateToken(user.ID, user.Username, h.jwtSecret, 30*24*3600*1e9)

	c.JSON(http.StatusOK, model.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    86400,
		User: model.UserInfo{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			Avatar:     user.Avatar,
			Phone:      user.Phone,
			TotalDist:  user.TotalDist,
			TotalRides: user.TotalRides,
		},
	})
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	// TODO: 实现刷新令牌逻辑
	c.JSON(http.StatusOK, gin.H{"message": "刷新成功"})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt("userID")
	user, err := h.svc.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, model.UserInfo{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Avatar:     user.Avatar,
		Phone:      user.Phone,
		TotalDist:  user.TotalDist,
		TotalRides: user.TotalRides,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt("userID")
	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	if err := h.svc.UpdateProfile(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *UserHandler) GetStats(c *gin.Context) {
	userID := c.GetInt("userID")
	user, err := h.svc.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, model.UserStats{
		TotalDist:   user.TotalDist,
		TotalRides:  user.TotalRides,
	})
}
