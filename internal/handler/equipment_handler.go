package handler

import (
	"net/http"
	"strconv"

	"github.com/gowild/server/internal/model"
	"github.com/gowild/server/internal/service"

	"github.com/gin-gonic/gin"
)

type EquipmentHandler struct {
	svc *service.EquipmentService
}

func NewEquipmentHandler(svc *service.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{svc: svc}
}

func (h *EquipmentHandler) ListEquipment(c *gin.Context) {
	userID := c.GetInt("userID")
	equips, err := h.svc.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if equips == nil {
		equips = []model.Equipment{}
	}
	c.JSON(http.StatusOK, equips)
}

func (h *EquipmentHandler) CreateEquipment(c *gin.Context) {
	userID := c.GetInt("userID")
	var req model.CreateEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	equip, err := h.svc.Create(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, equip)
}

func (h *EquipmentHandler) GetEquipment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	equip, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "装备不存在"})
		return
	}

	c.JSON(http.StatusOK, equip)
}

func (h *EquipmentHandler) UpdateEquipment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req model.UpdateEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	if err := h.svc.Update(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *EquipmentHandler) DeleteEquipment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *EquipmentHandler) UpdateEquipmentStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req model.StatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	if err := h.svc.UpdateStatus(id, req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}

func (h *EquipmentHandler) GetCategories(c *gin.Context) {
	userID := c.GetInt("userID")
	categories, err := h.svc.GetCategories(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if categories == nil {
		categories = []string{}
	}
	c.JSON(http.StatusOK, categories)
}

// 消耗品相关

func (h *EquipmentHandler) ListConsumables(c *gin.Context) {
	userID := c.GetInt("userID")
	consumables, err := h.svc.ListConsumables(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if consumables == nil {
		consumables = []model.Consumable{}
	}
	c.JSON(http.StatusOK, consumables)
}

func (h *EquipmentHandler) CreateConsumable(c *gin.Context) {
	userID := c.GetInt("userID")
	var req model.CreateConsumableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	consumable, err := h.svc.CreateConsumable(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, consumable)
}

func (h *EquipmentHandler) UpdateConsumable(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req model.UpdateConsumableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	if err := h.svc.UpdateConsumable(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *EquipmentHandler) DeleteConsumable(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.DeleteConsumable(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ListTodos 获取待办事项
func (h *EquipmentHandler) ListTodos(c *gin.Context) {
	userID := c.GetInt("userID")

	// 从装备状态和消耗品中生成待办
	var todos []model.TodoItem

	// 装备维护提醒
	equips, _ := h.svc.ListByUser(userID)
	for _, e := range equips {
		switch e.Status {
		case "maintenance":
			todos = append(todos, model.TodoItem{
				Type:      "equipment",
				ID:        e.ID,
				Name:      e.Name + " 需要维修",
				Source:    "equipment",
				Color:     "#E74C3C",
				ActionURL: "/equipment/" + strconv.Itoa(e.ID),
			})
		case "broken":
			todos = append(todos, model.TodoItem{
				Type:      "equipment",
				ID:        e.ID,
				Name:      e.Name + " 已损坏",
				Source:    "equipment",
				Color:     "#C0392B",
				ActionURL: "/equipment/" + strconv.Itoa(e.ID),
			})
		}
	}

	// 消耗品不足提醒
	consumables, _ := h.svc.ListConsumables(userID)
	for _, c := range consumables {
		if c.Percent <= 20 {
			todos = append(todos, model.TodoItem{
				Type:      "consumable",
				ID:        c.ID,
				Name:      c.Name + " 仅剩 " + strconv.Itoa(c.Percent) + "%",
				Source:    "consumable",
				Color:     "#F39C12",
				ActionURL: "/consumables/" + strconv.Itoa(c.ID),
			})
		}
	}

	if todos == nil {
		todos = []model.TodoItem{}
	}
	c.JSON(http.StatusOK, todos)
}
