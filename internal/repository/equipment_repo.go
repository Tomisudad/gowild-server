package repository

import (
	"database/sql"
	"fmt"

	"github.com/gowild/server/internal/model"
)

type EquipmentRepository struct {
	db *sql.DB
}

func NewEquipmentRepository(db *sql.DB) *EquipmentRepository {
	return &EquipmentRepository{db: db}
}

func (r *EquipmentRepository) Create(equip *model.Equipment) error {
	return r.db.QueryRow(
		`INSERT INTO equipment (user_id, name, icon, detail, status, category) 
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`,
		equip.UserID, equip.Name, equip.Icon, equip.Detail, equip.Status, equip.Category,
	).Scan(&equip.ID, &equip.CreatedAt, &equip.UpdatedAt)
}

func (r *EquipmentRepository) FindByID(id int) (*model.Equipment, error) {
	equip := &model.Equipment{}
	err := r.db.QueryRow(
		`SELECT id, user_id, name, icon, detail, status, category, created_at, updated_at 
		 FROM equipment WHERE id=$1`, id,
	).Scan(&equip.ID, &equip.UserID, &equip.Name, &equip.Icon, &equip.Detail,
		&equip.Status, &equip.Category, &equip.CreatedAt, &equip.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("查询装备失败: %w", err)
	}
	return equip, nil
}

func (r *EquipmentRepository) ListByUser(userID int) ([]model.Equipment, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, name, icon, detail, status, category, created_at, updated_at 
		 FROM equipment WHERE user_id=$1 ORDER BY category, created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equips []model.Equipment
	for rows.Next() {
		var e model.Equipment
		if err := rows.Scan(&e.ID, &e.UserID, &e.Name, &e.Icon, &e.Detail,
			&e.Status, &e.Category, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		equips = append(equips, e)
	}
	return equips, nil
}

func (r *EquipmentRepository) Update(equip *model.Equipment) error {
	_, err := r.db.Exec(
		`UPDATE equipment SET name=$1, icon=$2, detail=$3, status=$4, updated_at=CURRENT_TIMESTAMP WHERE id=$5`,
		equip.Name, equip.Icon, equip.Detail, equip.Status, equip.ID,
	)
	return err
}

func (r *EquipmentRepository) UpdateStatus(id int, status string) error {
	_, err := r.db.Exec(
		`UPDATE equipment SET status=$1, updated_at=CURRENT_TIMESTAMP WHERE id=$2`,
		status, id,
	)
	return err
}

func (r *EquipmentRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM equipment WHERE id=$1`, id)
	return err
}

func (r *EquipmentRepository) GetCategories(userID int) ([]string, error) {
	rows, err := r.db.Query(
		`SELECT DISTINCT category FROM equipment WHERE user_id=$1 ORDER BY category`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// 消耗品相关

func (r *EquipmentRepository) CreateConsumable(c *model.Consumable) error {
	return r.db.QueryRow(
		`INSERT INTO consumables (user_id, name, percent) VALUES ($1,$2,$3) RETURNING id, created_at, updated_at`,
		c.UserID, c.Name, c.Percent,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *EquipmentRepository) ListConsumables(userID int) ([]model.Consumable, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, name, percent, created_at, updated_at 
		 FROM consumables WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consumables []model.Consumable
	for rows.Next() {
		var c model.Consumable
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Percent, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		consumables = append(consumables, c)
	}
	return consumables, nil
}

func (r *EquipmentRepository) UpdateConsumable(c *model.Consumable) error {
	_, err := r.db.Exec(
		`UPDATE consumables SET name=$1, percent=$2, updated_at=CURRENT_TIMESTAMP WHERE id=$3`,
		c.Name, c.Percent, c.ID,
	)
	return err
}

func (r *EquipmentRepository) DeleteConsumable(id int) error {
	_, err := r.db.Exec(`DELETE FROM consumables WHERE id=$1`, id)
	return err
}

func (r *EquipmentRepository) FindConsumableByID(id int) (*model.Consumable, error) {
	c := &model.Consumable{}
	err := r.db.QueryRow(
		`SELECT id, user_id, name, percent, created_at, updated_at FROM consumables WHERE id=$1`, id,
	).Scan(&c.ID, &c.UserID, &c.Name, &c.Percent, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("查询消耗品失败: %w", err)
	}
	return c, nil
}
