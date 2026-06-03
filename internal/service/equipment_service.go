package service

import (
	"errors"

	"github.com/gowild/server/internal/model"
	"github.com/gowild/server/internal/repository"
)

type EquipmentService struct {
	repo *repository.EquipmentRepository
}

func NewEquipmentService(repo *repository.EquipmentRepository) *EquipmentService {
	return &EquipmentService{repo: repo}
}

func (s *EquipmentService) Create(userID int, req model.CreateEquipmentRequest) (*model.Equipment, error) {
	if req.Icon == "" {
		req.Icon = "bike"
	}

	equip := &model.Equipment{
		UserID:   userID,
		Name:     req.Name,
		Icon:     req.Icon,
		Detail:   req.Detail,
		Status:   req.Status,
		Category: req.Category,
	}

	if err := s.repo.Create(equip); err != nil {
		return nil, err
	}

	return equip, nil
}

func (s *EquipmentService) GetByID(id int) (*model.Equipment, error) {
	return s.repo.FindByID(id)
}

func (s *EquipmentService) ListByUser(userID int) ([]model.Equipment, error) {
	return s.repo.ListByUser(userID)
}

func (s *EquipmentService) Update(id int, req model.UpdateEquipmentRequest) error {
	equip, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if req.Name != "" {
		equip.Name = req.Name
	}
	if req.Icon != "" {
		equip.Icon = req.Icon
	}
	if req.Detail != "" {
		equip.Detail = req.Detail
	}
	if req.Status != "" {
		equip.Status = req.Status
	}

	return s.repo.Update(equip)
}

func (s *EquipmentService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *EquipmentService) UpdateStatus(id int, status string) error {
	validStatuses := map[string]bool{"good": true, "maintenance": true, "broken": true, "retired": true}
	if !validStatuses[status] {
		return errors.New("无效的装备状态，允许: good, maintenance, broken, retired")
	}
	return s.repo.UpdateStatus(id, status)
}

func (s *EquipmentService) GetCategories(userID int) ([]string, error) {
	return s.repo.GetCategories(userID)
}

// 消耗品相关

func (s *EquipmentService) CreateConsumable(userID int, req model.CreateConsumableRequest) (*model.Consumable, error) {
	c := &model.Consumable{
		UserID:  userID,
		Name:    req.Name,
		Percent: req.Percent,
	}
	if err := s.repo.CreateConsumable(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *EquipmentService) ListConsumables(userID int) ([]model.Consumable, error) {
	return s.repo.ListConsumables(userID)
}

func (s *EquipmentService) UpdateConsumable(id int, req model.UpdateConsumableRequest) error {
	c, err := s.repo.FindConsumableByID(id)
	if err != nil {
		return err
	}

	if req.Name != "" {
		c.Name = req.Name
	}
	c.Percent = req.Percent

	return s.repo.UpdateConsumable(c)
}

func (s *EquipmentService) DeleteConsumable(id int) error {
	return s.repo.DeleteConsumable(id)
}
