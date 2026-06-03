package repository

import (
	"database/sql"
	"fmt"

	"github.com/gowild/server/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.QueryRow(
		`INSERT INTO users (username, email, password, phone) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`,
		user.Username, user.Email, user.Password, user.Phone,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) FindByID(id int) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(
		`SELECT id, username, email, password, avatar, phone, total_dist, total_rides, created_at, updated_at FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.Phone, &user.TotalDist, &user.TotalRides, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(
		`SELECT id, username, email, password, avatar, phone, total_dist, total_rides, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.Phone, &user.TotalDist, &user.TotalRides, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return user, nil
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(
		`SELECT id, username, email, password, avatar, phone, total_dist, total_rides, created_at, updated_at FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.Phone, &user.TotalDist, &user.TotalRides, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	_, err := r.db.Exec(
		`UPDATE users SET username=$1, avatar=$2, phone=$3, updated_at=CURRENT_TIMESTAMP WHERE id=$4`,
		user.Username, user.Avatar, user.Phone, user.ID,
	)
	return err
}

func (r *UserRepository) UpdateStats(userID int, dist float64) error {
	_, err := r.db.Exec(
		`UPDATE users SET total_dist = total_dist + $1, total_rides = total_rides + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`,
		dist, userID,
	)
	return err
}
