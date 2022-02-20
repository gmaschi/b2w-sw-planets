package planetcontroller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	planetmodel "github.com/gmaschi/b2w-sw-planets/internal/models/planet"
	errorsmodel "github.com/gmaschi/b2w-sw-planets/internal/models/planet/errors-model"
	planetsdb "github.com/gmaschi/b2w-sw-planets/internal/services/datastore/mongodb/planets-db"
	parseerrors "github.com/gmaschi/b2w-sw-planets/pkg/tools/parse-errors"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type Controller struct {
	store planetsdb.Store
}

// New creates a pointer to a Controller
func New(store planetsdb.Store) *Controller {
	return &Controller{
		store: store,
	}
}

// Create handles the request to create a new planet
func (c *Controller) Create(ctx *gin.Context) {
	var req planetmodel.CreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, parseerrors.ErrorResponse(err))
		return
	}
	createArgs := planetsdb.CreatePlanetParams{
		Name:    req.Name,
		Terrain: req.Terrain,
		Climate: req.Climate,
	}
	planet, err := c.store.CreatePlanet(ctx, createArgs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, parseerrors.ErrorResponse(err))
		return
	}

	res := planetmodel.CreateResponse(planet)
	ctx.JSON(http.StatusCreated, res)
}

// Planet handles the request to get a planet based on the ID
func (c *Controller) Planet(ctx *gin.Context) {
	var req planetmodel.GetRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, parseerrors.ErrorResponse(err))
		return
	}

	planet, err := c.store.GetPlanet(ctx, req.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, parseerrors.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, parseerrors.ErrorResponse(err))
		return
	}

	res := planetmodel.GetResponse(planet)
	ctx.JSON(http.StatusOK, res)
}

// Delete handles the request to delete a planet based on the ID
func (c *Controller) Delete(ctx *gin.Context) {
	var req planetmodel.DeleteRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, parseerrors.ErrorResponse(err))
		return
	}

	err := c.store.DeletePlanet(ctx, req.ID)
	// TODO: handle invalid hex
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, parseerrors.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, fmt.Sprintf("planet with id %s deleted", req.ID))
}

// List handles the request to list the planets
func (c *Controller) List(ctx *gin.Context) {
	var req planetmodel.ListRequest

	if err := ctx.ShouldBindQuery(&req); err != nil && req.Name != "" {
		ctx.JSON(http.StatusBadRequest, parseerrors.ErrorResponse(err))
		return
	}

	listArgs := planetsdb.ListPlanetParams{
		Name: req.Name,
	}

	planets, err := c.store.ListPlanets(ctx, listArgs)
	if err != nil {
		fmt.Println(err)
		if err.Error() == errorsmodel.PlanetDoesNotExist {
			ctx.JSON(http.StatusNotFound, parseerrors.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, parseerrors.ErrorResponse(err))
		return
	}

	k := len(planets)
	res := make([]planetmodel.ListResponse, 0, k)
	for _, planet := range planets {
		res = append(res, planetmodel.ListResponse(planet))
	}

	ctx.JSON(http.StatusOK, res)
}
