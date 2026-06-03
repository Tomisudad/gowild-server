package handler

import (
	"net/http"
	"strconv"

	"github.com/gowild/server/internal/model"
	"github.com/gowild/server/internal/service"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	svc *service.TeamService
}

func NewTeamHandler(svc *service.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	userID := c.GetInt("userID")
	var req model.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	team, err := h.svc.Create(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

func (h *TeamHandler) JoinTeam(c *gin.Context) {
	userID := c.GetInt("userID")
	var req model.JoinTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	team, err := h.svc.Join(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) LeaveTeam(c *gin.Context) {
	userID := c.GetInt("userID")
	teamID, _ := strconv.Atoi(c.Param("id"))

	if err := h.svc.Leave(teamID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "已退出车队"})
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("id"))
	team, err := h.svc.GetTeam(teamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "车队不存在"})
		return
	}

	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) GetMyTeam(c *gin.Context) {
	userID := c.GetInt("userID")
	team, err := h.svc.GetMyTeam(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "你尚未加入任何车队"})
		return
	}

	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) GetMembers(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("id"))
	members, err := h.svc.GetMembers(teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if members == nil {
		members = []model.TeamMember{}
	}
	c.JSON(http.StatusOK, members)
}

func (h *TeamHandler) SendMessage(c *gin.Context) {
	userID := c.GetInt("userID")
	teamID, _ := strconv.Atoi(c.Param("id"))

	var req model.TeamMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	msg, err := h.svc.SendMessage(teamID, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)
}

func (h *TeamHandler) GetMessages(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("id"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := h.svc.GetMessages(teamID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if messages == nil {
		messages = []model.TeamMessage{}
	}
	c.JSON(http.StatusOK, messages)
}

func (h *TeamHandler) UpdateLocation(c *gin.Context) {
	userID := c.GetInt("userID")
	teamID, _ := strconv.Atoi(c.Param("id"))

	var req model.TeamLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	if err := h.svc.UpdateLocation(teamID, userID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "位置已更新"})
}

func (h *TeamHandler) GetLocations(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("id"))
	locs, err := h.svc.GetLocations(teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if locs == nil {
		locs = []model.TeamLocation{}
	}
	c.JSON(http.StatusOK, locs)
}
