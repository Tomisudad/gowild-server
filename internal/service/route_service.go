package service

import (
	"errors"

	"github.com/gowild/server/internal/model"
	"github.com/gowild/server/internal/repository"
)

type RouteService struct {
	repo *repository.RouteRepository
}

func NewRouteService(repo *repository.RouteRepository) *RouteService {
	return &RouteService{repo: repo}
}

func (s *RouteService) Create(userID int, req model.CreateRouteRequest) (*model.Route, error) {
	route := &model.Route{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Distance:    req.Distance,
		Duration:    req.Duration,
		Elevation:   req.Elevation,
		Difficulty:  req.Difficulty,
		Coordinates: req.Coordinates,
		IsFavorite:  false,
		IsPublic:    req.IsPublic,
	}

	if err := s.repo.Create(route); err != nil {
		return nil, err
	}

	// 创建分段
	for _, seg := range req.Segments {
		segment := &model.RouteSegment{
			RouteID:  route.ID,
			Name:     seg.Name,
			Distance: seg.Distance,
			Elevation: seg.Elevation,
			Terrain:  seg.Terrain,
			Tips:     seg.Tips,
			OrderNum: seg.OrderNum,
		}
		if err := s.repo.CreateSegment(segment); err != nil {
			return nil, err
		}
	}

	return route, nil
}

func (s *RouteService) GetByID(id int) (*model.Route, error) {
	return s.repo.FindByID(id)
}

func (s *RouteService) ListByUser(userID, page, pageSize int) ([]model.Route, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.repo.ListByUser(userID, page, pageSize)
}

func (s *RouteService) ListPublic(page, pageSize int) ([]model.Route, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.repo.ListPublic(page, pageSize)
}

func (s *RouteService) Update(id int, route *model.Route) error {
	route.ID = id
	return s.repo.Update(route)
}

func (s *RouteService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *RouteService) ToggleFavorite(id int) (bool, error) {
	return s.repo.ToggleFavorite(id)
}

func (s *RouteService) GetSegments(routeID int) ([]model.RouteSegment, error) {
	return s.repo.GetSegments(routeID)
}

// AnalyzeRoute 分析路线数据
func (s *RouteService) AnalyzeRoute(id int) (*model.RouteAnalysis, error) {
	route, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	segments, err := s.repo.GetSegments(id)
	if err != nil {
		return nil, err
	}

	analysis := &model.RouteAnalysis{
		Difficulty:    route.Difficulty,
		DistanceKm:    route.Distance,
		ElevationGain: route.Elevation,
		Segments:      len(segments),
	}

	// 计算平均坡度
	if route.Distance > 0 {
		analysis.AvgGrade = (float64(route.Elevation) / (route.Distance * 1000)) * 100
	} else {
		analysis.AvgGrade = 0
	}

	// 估计最大坡度（从分段数据）
	for _, seg := range segments {
		if seg.Distance > 0 {
			grade := (float64(seg.Elevation) / (seg.Distance * 1000)) * 100
			if grade > analysis.MaxGrade {
				analysis.MaxGrade = grade
			}
		}
	}

	// 估计时间
	avgSpeedKph := 18.0
	switch route.Difficulty {
	case "easy":
		avgSpeedKph = 22.0
	case "medium":
		avgSpeedKph = 18.0
	case "hard":
		avgSpeedKph = 14.0
	case "extreme":
		avgSpeedKph = 10.0
	}
	analysis.EstTime = int((route.Distance / avgSpeedKph) * 60)

	// 路面描述
	surfaceTypes := make(map[string]int)
	for _, seg := range segments {
		surfaceTypes[seg.Terrain]++
	}
	if len(surfaceTypes) > 0 {
		maxType := ""
		maxCount := 0
		for t, c := range surfaceTypes {
			if c > maxCount {
				maxType = t
				maxCount = c
			}
		}
		analysis.SurfaceDesc = maxType
	} else {
		analysis.SurfaceDesc = "paved"
	}

	return analysis, nil
}

// ImportGPX 导入GPX文件（placeholder, 实际解析由handler完成）
func (s *RouteService) ImportGPX(userID int, gpx model.GPXData) (*model.Route, error) {
	if gpx.Name == "" {
		return nil, errors.New("路线名称为空")
	}

	route := &model.Route{
		UserID:      userID,
		Name:        gpx.Name,
		Distance:    gpx.Distance,
		Elevation:   gpx.Elevation,
		Difficulty:  "medium",
		Coordinates: gpx.Coordinates,
		IsPublic:    false,
	}

	if err := s.repo.Create(route); err != nil {
		return nil, err
	}

	for i, seg := range gpx.Segments {
		segment := &model.RouteSegment{
			RouteID:  route.ID,
			Name:     seg.Name,
			Distance: seg.Distance,
			Elevation: seg.Elevation,
			Terrain:  "paved",
			OrderNum: i + 1,
		}
		if err := s.repo.CreateSegment(segment); err != nil {
			return nil, err
		}
	}

	return route, nil
}
