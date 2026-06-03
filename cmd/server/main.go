package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gowild/server/internal/config"
	"github.com/gowild/server/internal/handler"
	"github.com/gowild/server/internal/middleware"
	"github.com/gowild/server/internal/repository"
	"github.com/gowild/server/internal/service"
	"github.com/gowild/server/internal/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	dbCfg := &repository.Config{
		DatabaseURL:         cfg.DatabaseURL,
		DatabaseMaxOpenConns: cfg.DatabaseMaxOpenConns,
		DatabaseMaxIdleConns: cfg.DatabaseMaxIdleConns,
		DatabaseMaxLifetime:  cfg.DatabaseMaxLifetime,
	}
	db, err := repository.InitDB(dbCfg)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer db.Close()

	// 初始化仓储层
	userRepo := repository.NewUserRepository(db)
	routeRepo := repository.NewRouteRepository(db)
	equipRepo := repository.NewEquipmentRepository(db)
	rideRepo := repository.NewRideRepository(db)
	teamRepo := repository.NewTeamRepository(db)

	// 初始化服务层
	userSvc := service.NewUserService(userRepo, cfg.JWTSecret)
	routeSvc := service.NewRouteService(routeRepo)
	equipSvc := service.NewEquipmentService(equipRepo)
	rideSvc := service.NewRideService(rideRepo)
	rideSvc.SetUserService(userSvc)
	teamSvc := service.NewTeamService(teamRepo)

	// 初始化处理器
	userH := handler.NewUserHandler(userSvc)
	routeH := handler.NewRouteHandler(routeSvc)
	equipH := handler.NewEquipmentHandler(equipSvc)
	rideH := handler.NewRideHandler(rideSvc)
	teamH := handler.NewTeamHandler(teamSvc)

	// 初始化Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"version":   "1.0.0",
			"timestamp": time.Now().Unix(),
		})
	})

	// 公开路由
	api := r.Group("/api/v1")
	{
		// 认证
		auth := api.Group("/auth")
		{
			auth.POST("/register", userH.Register)
			auth.POST("/login", userH.Login)
			auth.POST("/refresh", userH.RefreshToken)
		}

		// 公开路线
		api.GET("/routes/public", routeH.ListPublicRoutes)
	}

	// 需要认证的路由
	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// 用户
		auth.GET("/user/profile", userH.GetProfile)
		auth.PUT("/user/profile", userH.UpdateProfile)
		auth.GET("/user/stats", userH.GetStats)

		// 路线
		auth.GET("/routes", routeH.ListRoutes)
		auth.POST("/routes", routeH.CreateRoute)
		auth.GET("/routes/:id", routeH.GetRoute)
		auth.PUT("/routes/:id", routeH.UpdateRoute)
		auth.DELETE("/routes/:id", routeH.DeleteRoute)
		auth.POST("/routes/:id/favorite", routeH.ToggleFavorite)
		auth.POST("/routes/import/gpx", routeH.ImportGPX)
		auth.GET("/routes/:id/segments", routeH.GetSegments)
		auth.GET("/routes/:id/analysis", routeH.GetRouteAnalysis)
		auth.GET("/weather", routeH.GetWeather)

		// 骑行
		auth.POST("/rides/start", rideH.StartRide)
		auth.POST("/rides/:id/end", rideH.EndRide)
		auth.POST("/rides/:id/pause", rideH.PauseRide)
		auth.POST("/rides/:id/resume", rideH.ResumeRide)
		auth.GET("/rides", rideH.ListRides)
		auth.GET("/rides/active", rideH.GetActiveRide)
		auth.GET("/rides/:id", rideH.GetRide)
		auth.GET("/rides/:id/points", rideH.GetRidePoints)
		auth.GET("/rides/stats", rideH.GetStats)
		auth.GET("/rides/compare/:routeId", rideH.CompareRides)
		auth.POST("/rides/sos", rideH.TriggerSOS)

		// 装备
		auth.GET("/equipment", equipH.ListEquipment)
		auth.POST("/equipment", equipH.CreateEquipment)
		auth.GET("/equipment/:id", equipH.GetEquipment)
		auth.PUT("/equipment/:id", equipH.UpdateEquipment)
		auth.DELETE("/equipment/:id", equipH.DeleteEquipment)
		auth.PATCH("/equipment/:id/status", equipH.UpdateEquipmentStatus)
		auth.GET("/equipment/categories", equipH.GetCategories)
		auth.GET("/equipment/todos", equipH.ListTodos)

		// 消耗品
		auth.GET("/consumables", equipH.ListConsumables)
		auth.POST("/consumables", equipH.CreateConsumable)
		auth.PUT("/consumables/:id", equipH.UpdateConsumable)
		auth.DELETE("/consumables/:id", equipH.DeleteConsumable)

		// 车队
		auth.POST("/teams", teamH.CreateTeam)
		auth.POST("/teams/join", teamH.JoinTeam)
		auth.POST("/teams/:id/leave", teamH.LeaveTeam)
		auth.GET("/teams/mine", teamH.GetMyTeam)
		auth.GET("/teams/:id", teamH.GetTeam)
		auth.GET("/teams/:id/members", teamH.GetMembers)
		auth.POST("/teams/:id/messages", teamH.SendMessage)
		auth.GET("/teams/:id/messages", teamH.GetMessages)
		auth.POST("/teams/:id/location", teamH.UpdateLocation)
		auth.GET("/teams/:id/locations", teamH.GetLocations)
	}

	// WebSocket
	r.GET("/ws", func(c *gin.Context) {
		// 从query参数获取token
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证token"})
			return
		}

		claims, err := middleware.ParseToken(token, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			return
		}

		ws.HandleWebSocket(c, claims.UserID)
	})

	// 启动服务器
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("去野服务器启动在端口 %s (环境: %s)", cfg.Port, cfg.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭失败: %v", err)
	}

	log.Println("服务器已安全关闭")
}
