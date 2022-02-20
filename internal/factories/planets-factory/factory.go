package planetsfactory

import (
	"github.com/gin-gonic/gin"
	planetcontroller "github.com/gmaschi/b2w-sw-planets/internal/controllers/planet"
	planetsdb "github.com/gmaschi/b2w-sw-planets/internal/services/datastore/mongodb/planets-db"
)

type (
	Factory struct {
		store          planetsdb.Store
		planetsHandler planetsHandler
		Router         *gin.Engine
	}

	planetsHandler struct {
		planetsController *planetcontroller.Controller
	}
)

func New(store planetsdb.Store) (*Factory, error) {
	factory := &Factory{
		store: store,
		planetsHandler: planetsHandler{
			planetsController: planetcontroller.New(store),
		},
	}
	router := gin.Default()

	factory.setupRoutes(router)

	factory.Router = router
	return factory, nil
}

func (f *Factory) setupRoutes(router *gin.Engine) {
	planetsV1 := router.Group("/v1/planets")
	{
		planetsV1.POST("", f.planetsHandler.planetsController.Create)
		planetsV1.GET("/:id", f.planetsHandler.planetsController.Planet)
		planetsV1.GET("", f.planetsHandler.planetsController.List)
		planetsV1.DELETE("/:id", f.planetsHandler.planetsController.Delete)
	}
}

func (f *Factory) Start(address string) error {
	return f.Router.Run(address)
}
