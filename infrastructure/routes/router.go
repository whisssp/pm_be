package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	_ "pm/docs"
	"pm/infrastructure/controllers/handlers"
	"pm/infrastructure/jobs"
	"pm/infrastructure/persistences/base"
	"pm/utils"
)

type Server struct {
	Port        string
	Persistence *base.Persistence
}

func InitServer(port string, persistence *base.Persistence) *Server {
	return &Server{
		Port:        port,
		Persistence: persistence,
	}
}

func (s *Server) Run() {
	router := gin.Default()
	InitHelpers(s.Persistence)

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

	productRoute := NewProductRoutes(productHandler)
	categoryRoute := NewCategoryRoutes(categoryHandler)
	fileRoute := NewFileRoutes(fileHandler)

	v1 := router.Group("/api/v1")

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	productRoute.RegisterRoutes(v1)
	categoryRoute.RegisterRoutes(v1)
	fileRoute.RegisterRoutes(v1)
}

func InitHelpers(p *base.Persistence) {
	utils.InitCacheHelper(p)
	utils.InitSupabaseStorage(p)
}