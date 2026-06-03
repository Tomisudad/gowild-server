package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gowild/server/internal/model"
)

type RideRepository struct {
	db *sql.DB
}

func NewRideRepository(db *sql.DB) *RideRepository {
	return &RideRepository{db: db}
}

func (r *RideRepository) Create(ride *model.Ride) error {
	return r.db.QueryRow(
		`INSERT INTO rides (user_id, route_id, name, distance, elevation, start_time, is_active, is_paused) 
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, created_at, updated_at`,
		ride.UserID, ride.RouteID, ride.Name, ride.Distance, ride.Elevation,
		ride.StartTime, true, false,
	).Scan(&ride.ID, &ride.CreatedAt, &ride.UpdatedAt)
}

func (r *RideRepository) FindByID(id int) (*model.Ride, error) {
	ride := &model.Ride{}
	err := r.db.QueryRow(
		`SELECT id, user_id, route_id, name, distance, duration, avg_speed, max_speed,
		        elevation, calories, start_time, end_time, is_active, is_paused, coordinates, created_at, updated_at 
		 FROM rides WHERE id=$1`, id,
	).Scan(&ride.ID, &ride.UserID, &ride.RouteID, &ride.Name, &ride.Distance,
		&ride.Duration, &ride.AvgSpeed, &ride.MaxSpeed, &ride.Elevation, &ride.Calories,
		&ride.StartTime, &ride.EndTime, &ride.IsActive, &ride.IsPaused,
		&ride.Coordinates, &ride.CreatedAt, &ride.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("查询骑行记录失败: %w", err)
	}
	return ride, nil
}

func (r *RideRepository) FindActiveByUser(userID int) (*model.Ride, error) {
	ride := &model.Ride{}
	err := r.db.QueryRow(
		`SELECT id, user_id, route_id, name, distance, duration, avg_speed, max_speed,
		        elevation, calories, start_time, end_time, is_active, is_paused, coordinates, created_at, updated_at
		 FROM rides WHERE user_id=$1 AND is_active=true ORDER BY start_time DESC LIMIT 1`, userID,
	).Scan(&ride.ID, &ride.UserID, &ride.RouteID, &ride.Name, &ride.Distance,
		&ride.Duration, &ride.AvgSpeed, &ride.MaxSpeed, &ride.Elevation, &ride.Calories,
		&ride.StartTime, &ride.EndTime, &ride.IsActive, &ride.IsPaused,
		&ride.Coordinates, &ride.CreatedAt, &ride.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("查询活跃骑行失败: %w", err)
	}
	return ride, nil
}

func (r *RideRepository) ListByUser(userID, page, pageSize int) ([]model.Ride, int, error) {
	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM rides WHERE user_id=$1 AND is_active=false`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := r.db.Query(
		`SELECT id, user_id, route_id, name, distance, duration, avg_speed, max_speed,
		        elevation, calories, start_time, end_time, is_active, is_paused, coordinates, created_at, updated_at
		 FROM rides WHERE user_id=$1 AND is_active=false ORDER BY start_time DESC LIMIT $2 OFFSET $3`,
		userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rides []model.Ride
	for rows.Next() {
		var ride model.Ride
		if err := rows.Scan(&ride.ID, &ride.UserID, &ride.RouteID, &ride.Name, &ride.Distance,
			&ride.Duration, &ride.AvgSpeed, &ride.MaxSpeed, &ride.Elevation, &ride.Calories,
			&ride.StartTime, &ride.EndTime, &ride.IsActive, &ride.IsPaused,
			&ride.Coordinates, &ride.CreatedAt, &ride.UpdatedAt); err != nil {
			return nil, 0, err
		}
		rides = append(rides, ride)
	}
	return rides, total, nil
}

func (r *RideRepository) EndRide(id int, dist float64, duration int, avgSpeed, maxSpeed float64, calories int, coordinates string) error {
	_, err := r.db.Exec(
		`UPDATE rides SET is_active=false, is_paused=false, end_time=$1, 
		 distance=$2, duration=$3, avg_speed=$4, max_speed=$5, calories=$6, coordinates=$7, updated_at=CURRENT_TIMESTAMP 
		 WHERE id=$8`,
		time.Now(), dist, duration, avgSpeed, maxSpeed, calories, coordinates, id,
	)
	return err
}

func (r *RideRepository) PauseRide(id int) error {
	_, err := r.db.Exec(
		`UPDATE rides SET is_paused=true, updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id,
	)
	return err
}

func (r *RideRepository) ResumeRide(id int) error {
	_, err := r.db.Exec(
		`UPDATE rides SET is_paused=false, updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id,
	)
	return err
}

func (r *RideRepository) AddPoint(point *model.RidePoint) error {
	return r.db.QueryRow(
		`INSERT INTO ride_points (ride_id, latitude, longitude, speed, elevation) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		point.RideID, point.Latitude, point.Longitude, point.Speed, point.Elevation,
	).Scan(&point.ID)
}

func (r *RideRepository) GetPoints(rideID int) ([]model.RidePoint, error) {
	rows, err := r.db.Query(
		`SELECT id, ride_id, latitude, longitude, speed, elevation, timestamp 
		 FROM ride_points WHERE ride_id=$1 ORDER BY timestamp`, rideID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []model.RidePoint
	for rows.Next() {
		var p model.RidePoint
		if err := rows.Scan(&p.ID, &p.RideID, &p.Latitude, &p.Longitude,
			&p.Speed, &p.Elevation, &p.Timestamp); err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	return points, nil
}

func (r *RideRepository) ListByRoute(routeID, userID int) ([]model.Ride, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, route_id, name, distance, duration, avg_speed, max_speed,
		        elevation, calories, start_time, end_time, is_active, is_paused, coordinates, created_at, updated_at
		 FROM rides WHERE route_id=$1 AND user_id=$2 AND is_active=false ORDER BY start_time DESC`, routeID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []model.Ride
	for rows.Next() {
		var ride model.Ride
		if err := rows.Scan(&ride.ID, &ride.UserID, &ride.RouteID, &ride.Name, &ride.Distance,
			&ride.Duration, &ride.AvgSpeed, &ride.MaxSpeed, &ride.Elevation, &ride.Calories,
			&ride.StartTime, &ride.EndTime, &ride.IsActive, &ride.IsPaused,
			&ride.Coordinates, &ride.CreatedAt, &ride.UpdatedAt); err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}
	return rides, nil
}

func (r *RideRepository) GetStats(userID int) (*model.RideStats, error) {
	stats := &model.RideStats{}
	err := r.db.QueryRow(
		`SELECT 
			COUNT(*) as total_rides,
			COALESCE(SUM(distance), 0) as total_dist,
			COALESCE(SUM(duration), 0) as total_duration,
			COALESCE(AVG(NULLIF(avg_speed, 0)), 0) as avg_speed,
			COALESCE(MAX(max_speed), 0) as max_speed,
			COALESCE(SUM(calories), 0) as total_calories,
			COALESCE(MAX(distance), 0) as longest_ride
		 FROM rides WHERE user_id=$1 AND is_active=false`, userID,
	).Scan(&stats.TotalRides, &stats.TotalDist, &stats.TotalDuration,
		&stats.AvgSpeed, &stats.MaxSpeed, &stats.TotalCalories, &stats.LongestRide)
	if err != nil {
		return nil, err
	}

	// 本月统计
	err = r.db.QueryRow(
		`SELECT COALESCE(COUNT(*), 0), COALESCE(SUM(distance), 0) 
		 FROM rides WHERE user_id=$1 AND is_active=false 
		 AND start_time >= date_trunc('month', CURRENT_DATE)`, userID,
	).Scan(&stats.ThisMonthRides, &stats.ThisMonthDist)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
