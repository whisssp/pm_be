package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	_ "pm/docs"
	"pm/infrastructure/config"
	"pm/infrastructure/controllers/handlers"
	"pm/infrastructure/controllers/middleware"
	"pm/infrastructure/jobs"
	"pm/infrastructure/persistences/base"
	"pm/utils"
)

type Server struct {
	Port        string
	Persistence *base.Persistence
	appConfig   *config.AppConfig
}

func InitServer(port string, persistence *base.Persistence) *Server {
	return &Server{
		Port:        port,
		Persistence: persistence,
		appConfig:   config.Configs,
	}
}

func (s *Server) Run() {
	router := gin.New()
	router.SetTrustedProxies([]string{"0.0.0.0/0"})
	router.RemoteIPHeaders = []string{"X-Forwarded-For", "X-Real-IP"}

	s.InitHelpers()

	s.SetUpRoutes(router)

	c := cron.New()
	err := c.AddFunc("*/10 * * * *", func() {
		jobs.LoadProductToRedis(s.Persistence)
	})

	if err != nil {
		fmt.Println("Error adding cron job:", err)
	}

	go c.Run()

	err = router.Run(fmt.Sprintf(":%s", s.Port))
	if err != nil {
		log.Fatal("error running server", err)
		return
	}
}

func (s *Server) SetUpRoutes(router *gin.Engine) {
	productHandler := handlers.NewProductHandler(s.Persistence)
	categoryHandler := handlers.NewCategoryHandler(s.Persistence)
	fileHandler := handlers.NewFileHandler(s.Persistence)
	userHandler := handlers.NewUserHandler(s.Persistence)
	orderHandler := handlers.NewOrderHandler(s.Persistence)

	productRoute := NewProductRoutes(s.Persistence, productHandler)
	categoryRoute := NewCategoryRoutes(s.Persistence, categoryHandler)
	fileRoute := NewFileRoutes(s.Persistence, fileHandler)
	userRoute := NewUserRoutes(s.Persistence, userHandler)
	orderRoute := NewOrderRoutes(s.Persistence, orderHandler)

	router.Use(middleware.HoneycombHandler(), middleware.ErrorHandlingMiddleware())

	v1 := router.Group("/api/v1")

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	productRoute.RegisterRoutes(v1)
	categoryRoute.RegisterRoutes(v1)
	fileRoute.RegisterRoutes(v1)
	userRoute.RegisterRoutes(v1)
	orderRoute.RegisterRoutes(v1)
}

func (s *Server) InitHelpers() {
	utils.InitCacheHelper(s.Persistence)
	utils.InitSupabaseStorage(s.Persistence)
	utils.InitValidatorHelper()
	utils.InitJwtHelper(s.Persistence, s.appConfig.JwtConfig)
}