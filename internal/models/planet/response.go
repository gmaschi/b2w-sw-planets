package planetmodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CreateResponse struct {
		ID      primitive.ObjectID `json:"_id"`
		Name    string             `json:"name"`
		Terrain string             `json:"terrain"`
		Climate string             `json:"climate"`
		Movies  int                `json:"movies"`
	}

	GetResponse struct {
		ID      primitive.ObjectID `json:"_id"`
		Name    string             `json:"name"`
		Terrain string             `json:"terrain"`
		Climate string             `json:"climate"`
		Movies  int                `json:"movies"`
	}

	ListResponse struct {
		ID      primitive.ObjectID `json:"_id"`
		Name    string             `json:"name"`
		Terrain string             `json:"terrain"`
		Climate string             `json:"climate"`
		Movies  int                `json:"movies"`
	}

	DeleteResponse struct {
		ID primitive.ObjectID `json:"_id"`
	}
)
