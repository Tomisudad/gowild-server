package repository

import (
	"database/sql"
	"fmt"

	"github.com/gowild/server/internal/model"
)

type RouteRepository struct {
	db *sql.DB
}

func NewRouteRepository(db *sql.DB) *RouteRepository {
	return &RouteRepository{db: db}
}

func (r *RouteRepository) Create(route *model.Route) error {
	return r.db.QueryRow(
		`INSERT INTO routes (user_id, name, description, distance, duration, elevation, difficulty, coordinates, is_favorite, is_public) 
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id, created_at, updated_at`,
		route.UserID, route.Name, route.Description, route.Distance, route.Duration,
		route.Elevation, route.Difficulty, route.Coordinates, route.IsFavorite, route.IsPublic,
	).Scan(&route.ID, &route.CreatedAt, &route.UpdatedAt)
}

func (r *RouteRepository) FindByID(id int) (*model.Route, error) {
	route := &model.Route{}
	err := r.db.QueryRow(
		`SELECT id, user_id, name, description, distance, duration, elevation, difficulty, 
		        coordinates, is_favorite, is_public, created_at, updated_at FROM routes WHERE id=$1`, id,
	).Scan(&route.ID, &route.UserID, &route.Name, &route.Description, &route.Distance,
		&route.Duration, &route.Elevation, &route.Difficulty, &route.Coordinates,
		&route.IsFavorite, &route.IsPublic, &route.CreatedAt, &route.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("查询路线失败: %w", err)
	}
	return route, nil
}

func (r *RouteRepository) ListByUser(userID, page, pageSize int) ([]model.Route, int, error) {
	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM routes WHERE user_id=$1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := r.db.Query(
		`SELECT id, user_id, name, description, distance, duration, elevation, difficulty, 
		        coordinates, is_favorite, is_public, created_at, updated_at 
		 FROM routes WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var routes []model.Route
	for rows.Next() {
		var route model.Route
		if err := rows.Scan(&route.ID, &route.UserID, &route.Name, &route.Description,
			&route.Distance, &route.Duration, &route.Elevation, &route.Difficulty,
			&route.Coordinates, &route.IsFavorite, &route.IsPublic, &route.CreatedAt, &route.UpdatedAt); err != nil {
			return nil, 0, err
		}
		routes = append(routes, route)
	}

	return routes, total, nil
}

func (r *RouteRepository) ListPublic(page, pageSize int) ([]model.Route, int, error) {
	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM routes WHERE is_public=true`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := r.db.Query(
		`SELECT id, user_id, name, description, distance, duration, elevation, difficulty,
		        coordinates, is_favorite, is_public, created_at, updated_at
		 FROM routes WHERE is_public=true ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var routes []model.Route
	for rows.Next() {
		var route model.Route
		if err := rows.Scan(&route.ID, &route.UserID, &route.Name, &route.Description,
			&route.Distance, &route.Duration, &route.Elevation, &route.Difficulty,
			&route.Coordinates, &route.IsFavorite, &route.IsPublic, &route.CreatedAt, &route.UpdatedAt); err != nil {
			return nil, 0, err
		}
		routes = append(routes, route)
	}

	return routes, total, nil
}

func (r *RouteRepository) Update(route *model.Route) error {
	_, err := r.db.Exec(
		`UPDATE routes SET name=$1, description=$2, distance=$3, duration=$4, elevation=$5,
		 difficulty=$6, coordinates=$7, is_favorite=$8, is_public=$9, updated_at=CURRENT_TIMESTAMP WHERE id=$10`,
		route.Name, route.Description, route.Distance, route.Duration, route.Elevation,
		route.Difficulty, route.Coordinates, route.IsFavorite, route.IsPublic, route.ID,
	)
	return err
}

func (r *RouteRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM routes WHERE id=$1`, id)
	return err
}

func (r *RouteRepository) ToggleFavorite(id int) (bool, error) {
	var fav bool
	err := r.db.QueryRow(
		`UPDATE routes SET is_favorite = NOT is_favorite, updated_at = CURRENT_TIMESTAMP WHERE id=$1 RETURNING is_favorite`,
		id,
	).Scan(&fav)
	return fav, err
}

func (r *RouteRepository) CreateSegment(seg *model.RouteSegment) error {
	return r.db.QueryRow(
		`INSERT INTO route_segments (route_id, name, distance, elevation, terrain, tips, order_num) 
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		seg.RouteID, seg.Name, seg.Distance, seg.Elevation, seg.Terrain, seg.Tips, seg.OrderNum,
	).Scan(&seg.ID)
}

func (r *RouteRepository) GetSegments(routeID int) ([]model.RouteSegment, error) {
	rows, err := r.db.Query(
		`SELECT id, route_id, name, distance, elevation, terrain, tips, order_num 
		 FROM route_segments WHERE route_id=$1 ORDER BY order_num`, routeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segments []model.RouteSegment
	for rows.Next() {
		var seg model.RouteSegment
		if err := rows.Scan(&seg.ID, &seg.RouteID, &seg.Name, &seg.Distance,
			&seg.Elevation, &seg.Terrain, &seg.Tips, &seg.OrderNum); err != nil {
			return nil, err
		}
		segments = append(segments, seg)
	}
	return segments, nil
}

func (r *RouteRepository) GetTotalDistance(userID int) (float64, error) {
	var dist sql.NullFloat64
	err := r.db.QueryRow(`SELECT SUM(distance) FROM routes WHERE user_id=$1`, userID).Scan(&dist)
	if err != nil {
		return 0, err
	}
	if dist.Valid {
		return dist.Float64, nil
	}
	return 0, nil
}
