package planetsdb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Planet struct {
		ID      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		Name    string             `bson:"name" json:"name"`
		Terrain string             `bson:"terrain" json:"terrain"`
		Climate string             `bson:"climate" json:"climate"`
		Movies  int                `bson:"movies" json:"movies"`
	}

	Querier interface {
		CreatePlanet(ctx context.Context, arg CreatePlanetParams) (Planet, error)
		DeletePlanet(ctx context.Context, id string) error
		GetPlanet(ctx context.Context, id string) (Planet, error)
		ListPlanets(ctx context.Context, arg ListPlanetParams) ([]Planet, error)
	}

	swApiPlanetInfo struct {
		Count   int `json:"count"`
		Results []struct {
			Name  string   `json:"name"`
			Films []string `json:"films"`
		} `json:"results"`
	}
)
