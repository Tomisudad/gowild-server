package model

import (
	"time"
)

// ============================================
// 实体模型
// ============================================

// User 用户
type User struct {
	ID         int       `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	Email      string    `json:"email" db:"email"`
	Password   string    `json:"-" db:"password"`
	Avatar     string    `json:"avatar" db:"avatar"`
	Phone      string    `json:"phone" db:"phone"`
	TotalDist  float64   `json:"total_dist" db:"total_dist"`
	TotalRides int       `json:"total_rides" db:"total_rides"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Route 路线
type Route struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Distance    float64   `json:"distance" db:"distance"`
	Duration    int       `json:"duration" db:"duration"`
	Elevation   int       `json:"elevation" db:"elevation"`
	Difficulty  string    `json:"difficulty" db:"difficulty"`
	Coordinates string    `json:"coordinates" db:"coordinates"`
	IsFavorite  bool      `json:"is_favorite" db:"is_favorite"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// RouteSegment 路线分段
type RouteSegment struct {
	ID        int     `json:"id" db:"id"`
	RouteID   int     `json:"route_id" db:"route_id"`
	Name      string  `json:"name" db:"name"`
	Distance  float64 `json:"distance" db:"distance"`
	Elevation int     `json:"elevation" db:"elevation"`
	Terrain   string  `json:"terrain" db:"terrain"`
	Tips      string  `json:"tips" db:"tips"`
	OrderNum  int     `json:"order_num" db:"order_num"`
}

// Equipment 装备
type Equipment struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Icon      string    `json:"icon" db:"icon"`
	Detail    string    `json:"detail" db:"detail"`
	Status    string    `json:"status" db:"status"`
	Category  string    `json:"category" db:"category"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Consumable 消耗品
type Consumable struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Percent   int       `json:"percent" db:"percent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Ride 骑行记录
type Ride struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	RouteID     *int      `json:"route_id" db:"route_id"`
	Name        string    `json:"name" db:"name"`
	Distance    float64   `json:"distance" db:"distance"`
	Duration    int       `json:"duration" db:"duration"`
	AvgSpeed    float64   `json:"avg_speed" db:"avg_speed"`
	MaxSpeed    float64   `json:"max_speed" db:"max_speed"`
	Elevation   int       `json:"elevation" db:"elevation"`
	Calories    int       `json:"calories" db:"calories"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     *time.Time `json:"end_time" db:"end_time"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	IsPaused    bool      `json:"is_paused" db:"is_paused"`
	Coordinates string    `json:"coordinates" db:"coordinates"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// RidePoint 骑行轨迹点
type RidePoint struct {
	ID        int       `json:"id" db:"id"`
	RideID    int       `json:"ride_id" db:"ride_id"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	Speed     float64   `json:"speed" db:"speed"`
	Elevation int       `json:"elevation" db:"elevation"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

// Team 车队
type Team struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatorID int       `json:"creator_id" db:"creator_id"`
	Code      string    `json:"code" db:"code"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TeamMember 车队成员
type TeamMember struct {
	ID       int       `json:"id" db:"id"`
	TeamID   int       `json:"team_id" db:"team_id"`
	UserID   int       `json:"user_id" db:"user_id"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
	Color    string    `json:"color" db:"color"`
	Username string    `json:"username" db:"username"`
	Avatar   string    `json:"avatar" db:"avatar"`
}

// TeamMessage 车队消息
type TeamMessage struct {
	ID        int       `json:"id" db:"id"`
	TeamID    int       `json:"team_id" db:"team_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Content   string    `json:"content" db:"content"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TeamLocation 车队成员位置
type TeamLocation struct {
	ID        int       `json:"id" db:"id"`
	TeamID    int       `json:"team_id" db:"team_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	Speed     float64   `json:"speed" db:"speed"`
	Distance  float64   `json:"distance" db:"distance"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================
// 数据传输对象（DTO）
// ============================================

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	Token        string   `json:"token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	User         UserInfo `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID         int     `json:"id"`
	Username   string  `json:"username"`
	Email      string  `json:"email"`
	Avatar     string  `json:"avatar"`
	Phone      string  `json:"phone"`
	TotalDist  float64 `json:"total_dist"`
	TotalRides int     `json:"total_rides"`
}

// UserStats 用户统计
type UserStats struct {
	TotalDist   float64            `json:"total_dist"`
	TotalRides  int                `json:"total_rides"`
	TotalTime   int                `json:"total_time"`
	AvgSpeed    float64            `json:"avg_speed"`
	MaxSpeed    float64            `json:"max_speed"`
	TotalCalories int              `json:"total_calories"`
	MonthlyDist []MonthlyDistItem  `json:"monthly_dist"`
}

// MonthlyDistItem 月度里程项
type MonthlyDistItem struct {
	Month string  `json:"month"`
	Dist  float64 `json:"dist"`
}

// RideStats 骑行统计
type RideStats struct {
	TotalRides    int     `json:"total_rides"`
	TotalDist     float64 `json:"total_dist"`
	TotalDuration int     `json:"total_duration"`
	AvgSpeed      float64 `json:"avg_speed"`
	MaxSpeed      float64 `json:"max_speed"`
	TotalCalories int     `json:"total_calories"`
	ThisMonthRides int    `json:"this_month_rides"`
	ThisMonthDist  float64 `json:"this_month_dist"`
	LongestRide    float64 `json:"longest_ride"`
}

// RideCompareResult 骑行对比结果
type RideCompareResult struct {
	RideID      int     `json:"ride_id"`
	Duration    int     `json:"duration"`
	AvgSpeed    float64 `json:"avg_speed"`
	MaxSpeed    float64 `json:"max_speed"`
	Calories    int     `json:"calories"`
	ComparedToAvg float64 `json:"compared_to_avg"`
}

// EquipmentStatusRequest 装备状态更新请求
type StatusUpdateRequest struct {
	Status string `json:"status" binding:"required"`
}

// CreateRouteRequest 创建路线请求
type CreateRouteRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Distance    float64 `json:"distance" binding:"required"`
	Duration    int     `json:"duration" binding:"required"`
	Elevation   int     `json:"elevation"`
	Difficulty  string  `json:"difficulty" binding:"required"`
	Coordinates string  `json:"coordinates" binding:"required"`
	IsPublic    bool    `json:"is_public"`
	Segments    []SegmentDTO `json:"segments"`
}

// SegmentDTO 路线分段DTO
type SegmentDTO struct {
	Name      string  `json:"name"`
	Distance  float64 `json:"distance"`
	Elevation int     `json:"elevation"`
	Terrain   string  `json:"terrain"`
	Tips      string  `json:"tips"`
	OrderNum  int     `json:"order_num"`
}

// StartRideRequest 开始骑行请求
type StartRideRequest struct {
	RouteID     *int    `json:"route_id"`
	Name        string  `json:"name" binding:"required"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

// EndRideRequest 结束骑行请求
type EndRideRequest struct {
	Distance    float64 `json:"distance" binding:"required"`
	Duration    int     `json:"duration" binding:"required"`
	AvgSpeed    float64 `json:"avg_speed"`
	MaxSpeed    float64 `json:"max_speed"`
	Calories    int     `json:"calories"`
	Coordinates string  `json:"coordinates"`
}

// CreateTeamRequest 创建车队请求
type CreateTeamRequest struct {
	Name string `json:"name" binding:"required"`
}

// JoinTeamRequest 加入车队请求
type JoinTeamRequest struct {
	Code string `json:"code" binding:"required"`
}

// TeamMessageRequest 车队消息请求
type TeamMessageRequest struct {
	Content string `json:"content" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

// TeamLocationRequest 车队位置更新请求
type TeamLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Speed     float64 `json:"speed"`
	Distance  float64 `json:"distance"`
}

// SOSRequest SOS请求
type SOSRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Message   string  `json:"message"`
}

// CreateEquipmentRequest 创建装备请求
type CreateEquipmentRequest struct {
	Name     string `json:"name" binding:"required"`
	Icon     string `json:"icon"`
	Detail   string `json:"detail"`
	Status   string `json:"status" binding:"required"`
	Category string `json:"category" binding:"required"`
}

// UpdateEquipmentRequest 更新装备请求
type UpdateEquipmentRequest struct {
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Detail string `json:"detail"`
	Status string `json:"status"`
}

// CreateConsumableRequest 创建消耗品请求
type CreateConsumableRequest struct {
	Name    string `json:"name" binding:"required"`
	Percent int    `json:"percent" binding:"required,min=0,max=100"`
}

// UpdateConsumableRequest 更新消耗品请求
type UpdateConsumableRequest struct {
	Name    string `json:"name"`
	Percent int    `json:"percent" binding:"min=0,max=100"`
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
}

// ListResponse 分页列表响应
type ListResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// WeatherInfo 天气信息
type WeatherInfo struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Condition string  `json:"condition"`
	Wind      string  `json:"wind"`
	WindSpeed float64 `json:"wind_speed"`
	Humidity  int     `json:"humidity"`
	UVIndex   int     `json:"uv_index"`
	AQI       int     `json:"aqi"`
	Icon      string  `json:"icon"`
}

// TodoItem 待办事项
type TodoItem struct {
	Type      string `json:"type"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Source    string `json:"source"`
	Color     string `json:"color"`
	ActionURL string `json:"action_url"`
}

// GPXData GPX导入数据
type GPXData struct {
	Name        string        `json:"name"`
	Distance    float64       `json:"distance"`
	Elevation   int           `json:"elevation"`
	Coordinates string        `json:"coordinates"`
	Segments    []GPXSegment  `json:"segments"`
}

// GPXSegment GPX分段
type GPXSegment struct {
	Name      string  `json:"name"`
	Distance  float64 `json:"distance"`
	Elevation int     `json:"elevation"`
	Points    []LatLng `json:"points"`
}

// LatLng 经纬度
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	Ele int     `json:"ele"`
}

// WebSocket消息
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// TeamLocationUpdate WebSocket Team位置更新
type TeamLocationUpdate struct {
	TeamID    int              `json:"team_id"`
	Locations []TeamLocation  `json:"locations"`
}

// RouteAnalysis 路线分析
type RouteAnalysis struct {
	Difficulty    string  `json:"difficulty"`
	DistanceKm    float64 `json:"distance_km"`
	ElevationGain int     `json:"elevation_gain"`
	AvgGrade      float64 `json:"avg_grade"`
	MaxGrade      float64 `json:"max_grade"`
	EstTime       int     `json:"est_time"`
	Segments      int     `json:"segments"`
	SurfaceDesc   string  `json:"surface_desc"`
}
