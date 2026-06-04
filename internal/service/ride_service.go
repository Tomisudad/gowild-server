package service

import (
	"errors"
	"time"

	"github.com/gowild/server/internal/model"
	"github.com/gowild/server/internal/repository"
)

type RideService struct {
	repo     *repository.RideRepository
	userSvc  *UserService
}

func NewRideService(repo *repository.RideRepository) *RideService {
	return &RideService{repo: repo}
}

func (s *RideService) SetUserService(userSvc *UserService) {
	s.userSvc = userSvc
}

func (s *RideService) Start(userID int, req model.StartRideRequest) (*model.Ride, error) {
	// 妫€鏌ユ槸鍚︽湁娲昏穬鐨勯獞琛?	active, err := s.repo.FindActiveByUser(userID)
	if err == nil && active != nil {
		return nil, errors.New("宸叉湁杩涜涓殑楠戣锛岃鍏堢粨鏉熷綋鍓嶉獞琛?)
	}

	ride := &model.Ride{
		UserID:    userID,
		RouteID:   req.RouteID,
		Name:      req.Name,
		StartTime: time.Now(),
		IsActive:  true,
		IsPaused:  false,
	}

	if err := s.repo.Create(ride); err != nil {
		return nil, err
	}

	// 娣诲姞鍒濆杞ㄨ抗鐐?	if req.Latitude != 0 || req.Longitude != 0 {
		point := &model.RidePoint{
			RideID:    ride.ID,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
			Timestamp: time.Now(),
		}
		s.repo.AddPoint(point)
	}

	return ride, nil
}

func (s *RideService) End(rideID, userID int, dist float64, duration int, avgSpeed, maxSpeed float64, calories int, coordinates string) error {
	ride, err := s.repo.FindByID(rideID)
	if err != nil {
		return err
	}
	if ride.UserID != userID {
		return errors.New("鏃犳潈鎿嶄綔姝ら獞琛岃褰?)
	}
	if !ride.IsActive {
		return errors.New("楠戣宸茬粨鏉?)
	}

	if err := s.repo.EndRide(rideID, dist, duration, avgSpeed, maxSpeed, calories, coordinates); err != nil {
		return err
	}

	// 鏇存柊鐢ㄦ埛缁熻
	if s.userSvc != nil {
		s.userSvc.UpdateStats(userID, dist)
	}

	return nil
}

func (s *RideService) Pause(rideID, userID int) error {
	ride, err := s.repo.FindByID(rideID)
	if err != nil {
		return err
	}
	if ride.UserID != userID {
		return errors.New("鏃犳潈鎿嶄綔姝ら獞琛岃褰?)
	}
	if ride.IsPaused {
		return errors.New("楠戣宸叉殏鍋?)
	}

	return s.repo.PauseRide(rideID)
}

func (s *RideService) Resume(rideID, userID int) error {
	ride, err := s.repo.FindByID(rideID)
	if err != nil {
		return err
	}
	if ride.UserID != userID {
		return errors.New("鏃犳潈鎿嶄綔姝ら獞琛岃褰?)
	}
	if !ride.IsPaused {
		return errors.New("楠戣鏈殏鍋?)
	}

	return s.repo.ResumeRide(rideID)
}

func (s *RideService) GetByID(rideID int) (*model.Ride, error) {
	return s.repo.FindByID(rideID)
}

func (s *RideService) GetActive(userID int) (*model.Ride, error) {
	return s.repo.FindActiveByUser(userID)
}

func (s *RideService) ListByUser(userID, page, pageSize int) ([]model.Ride, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.repo.ListByUser(userID, page, pageSize)
}

func (s *RideService) GetPoints(rideID int) ([]model.RidePoint, error) {
	return s.repo.GetPoints(rideID)
}

func (s *RideService) GetStats(userID int) (*model.RideStats, error) {
	return s.repo.GetStats(userID)
}

// CompareRides 瀵规瘮鍚屼竴璺嚎鐨勯獞琛岃褰?func (s *RideService) CompareRides(routeID, userID int) ([]model.RideCompareResult, error) {
	rides, err := s.repo.ListByRoute(routeID, userID)
	if err != nil {
		return nil, err
	}

	if len(rides) == 0 {
		return nil, errors.New("璇ヨ矾绾挎殏鏃犻獞琛岃褰?)
	}

	// 璁＄畻骞冲潎鍊?	var totalDur, totalCal int
	var avgSpdSum, maxSpdSum float64
	for _, r := range rides {
		totalDur += r.Duration
		avgSpdSum += r.AvgSpeed
		maxSpdSum += r.MaxSpeed
		totalCal += r.Calories
	}
	n := float64(len(rides))
	avgSpeed := avgSpdSum / n

	var results []model.RideCompareResult
	for _, r := range rides {
		comparedToAvg := 0.0
		if avgSpeed > 0 {
			comparedToAvg = (r.AvgSpeed - avgSpeed) / avgSpeed * 100
		}
		results = append(results, model.RideCompareResult{
			RideID:        r.ID,
			Duration:      r.Duration,
			AvgSpeed:      r.AvgSpeed,
			MaxSpeed:      r.MaxSpeed,
			Calories:      r.Calories,
			ComparedToAvg: comparedToAvg,
		})
	}

	return results, nil
}

// TriggerSOS 瑙﹀彂绱ф€ユ眰鍔?func (s *RideService) TriggerSOS(userID int, lat, lng float64, msg string) error {
	// TODO: 闆嗘垚SMS/鎺ㄩ€佹湇鍔?	// 鍙戦€丼OS淇℃伅鑷崇揣鎬ヨ仈绯讳汉
	return nil
}
