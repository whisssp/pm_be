package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"pm/application"
	_ "pm/docs"
	"pm/infrastructure/controllers/handlers"
	"pm/infrastructure/persistences/base"
)

type Server struct {
	Port        string
	Persistence *base.Persistence
}

type RouterService struct {
	Product ItemService
}

type ItemService struct {
	Usecase interface{}
	Handler interface{}
}

func InitServer(port string, persistence *base.Persistence) *Server {
	return &Server{
		Port:        port,
		Persistence: persistence,
	}
}

func (s *Server) Run() {
	router := gin.Default()
	routerService := SetupRouterService(s.Persistence)
	s.SetUpRoutes(router, routerService)
	err := router.Run(fmt.Sprintf(":%s", s.Port))
	if err != nil {
		log.Fatal("error running server", err)
		return
	}
}

func (s *Server) SetUpRoutes(router *gin.Engine, service *RouterService) {

	productRoute := NewProductRoutes(service.Product.Handler.(*handlers.ProductHandler))

	v1 := router.Group("/api/v1")

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	productRoute.RegisterRoutes(v1)
}

func SetupRouterService(p *base.Persistence) *RouterService {
	productUsecase := application.NewProductUsecase(p)
	productHandler := handlers.NewProductHandler(productUsecase)

	return &RouterService{
		Product: ItemService{
			Usecase: productUsecase,
			Handler: productHandler,
		},
	}

}