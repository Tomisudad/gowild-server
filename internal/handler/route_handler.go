package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gowild/server/internal/model"
	"github.com/gowild/server/internal/service"

	"github.com/gin-gonic/gin"
)

type RouteHandler struct {
	svc *service.RouteService
}

func NewRouteHandler(svc *service.RouteService) *RouteHandler {
	return &RouteHandler{svc: svc}
}

func (h *RouteHandler) ListRoutes(c *gin.Context) {
	userID := c.GetInt("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	routes, total, err := h.svc.ListByUser(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	c.JSON(http.StatusOK, model.ListResponse{
		Data:       routes,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *RouteHandler) ListPublicRoutes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	routes, total, err := h.svc.ListPublic(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	c.JSON(http.StatusOK, model.ListResponse{
		Data:       routes,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *RouteHandler) CreateRoute(c *gin.Context) {
	userID := c.GetInt("userID")
	var req model.CreateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	route, err := h.svc.Create(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, route)
}

func (h *RouteHandler) GetRoute(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	route, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "路线不存在"})
		return
	}

	c.JSON(http.StatusOK, route)
}

func (h *RouteHandler) UpdateRoute(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var route model.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	if err := h.svc.Update(id, &route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *RouteHandler) DeleteRoute(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *RouteHandler) ToggleFavorite(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	isFav, err := h.svc.ToggleFavorite(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_favorite": isFav})
}

func (h *RouteHandler) ImportGPX(c *gin.Context) {
	userID := c.GetInt("userID")
	var gpx model.GPXData
	if err := c.ShouldBindJSON(&gpx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	route, err := h.svc.ImportGPX(userID, gpx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "导入失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, route)
}

func (h *RouteHandler) GetSegments(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	segments, err := h.svc.GetSegments(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, segments)
}

func (h *RouteHandler) GetRouteAnalysis(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	analysis, err := h.svc.AnalyzeRoute(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetWeather 获取天气信息（模拟）
func (h *RouteHandler) GetWeather(c *gin.Context) {
	c.JSON(http.StatusOK, model.WeatherInfo{
		Temp:      24.5,
		FeelsLike: 26.0,
		Condition: "晴",
		Wind:      "西南风",
		WindSpeed: 12.5,
		Humidity:  55,
		UVIndex:   6,
		AQI:       42,
		Icon:      "sunny",
	})
}
