package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gowild/server/internal/model"
	"github.com/gowild/server/internal/service"

	"github.com/gin-gonic/gin"
)

type RideHandler struct {
	svc *service.RideService
}

func NewRideHandler(svc *service.RideService) *RideHandler {
	return &RideHandler{svc: svc}
}

func (h *RideHandler) StartRide(c *gin.Context) {
	userID := c.GetInt("userID")
	var req model.StartRideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	ride, err := h.svc.Start(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ride)
}

func (h *RideHandler) EndRide(c *gin.Context) {
	userID := c.GetInt("userID")
	id, _ := strconv.Atoi(c.Param("id"))

	var req model.EndRideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	if err := h.svc.End(id, userID, req.Distance, req.Duration, req.AvgSpeed, req.MaxSpeed, req.Calories, req.Coordinates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "骑行已结束"})
}

func (h *RideHandler) PauseRide(c *gin.Context) {
	userID := c.GetInt("userID")
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.svc.Pause(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "骑行已暂停"})
}

func (h *RideHandler) ResumeRide(c *gin.Context) {
	userID := c.GetInt("userID")
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.svc.Resume(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "骑行已恢复"})
}

func (h *RideHandler) GetRide(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ride, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "骑行记录不存在"})
		return
	}

	c.JSON(http.StatusOK, ride)
}

func (h *RideHandler) GetActiveRide(c *gin.Context) {
	userID := c.GetInt("userID")
	ride, err := h.svc.GetActive(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "无进行中的骑行"})
		return
	}

	c.JSON(http.StatusOK, ride)
}

func (h *RideHandler) ListRides(c *gin.Context) {
	userID := c.GetInt("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	rides, total, err := h.svc.ListByUser(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	c.JSON(http.StatusOK, model.ListResponse{
		Data:       rides,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *RideHandler) GetRidePoints(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	points, err := h.svc.GetPoints(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if points == nil {
		points = []model.RidePoint{}
	}
	c.JSON(http.StatusOK, points)
}

func (h *RideHandler) GetStats(c *gin.Context) {
	userID := c.GetInt("userID")
	stats, err := h.svc.GetStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *RideHandler) CompareRides(c *gin.Context) {
	userID := c.GetInt("userID")
	routeID, _ := strconv.Atoi(c.Param("routeId"))

	results, err := h.svc.CompareRides(routeID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *RideHandler) TriggerSOS(c *gin.Context) {
	userID := c.GetInt("userID")
	var req model.SOSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.svc.TriggerSOS(userID, req.Latitude, req.Longitude, req.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SOS发送失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SOS已发出"})
}
