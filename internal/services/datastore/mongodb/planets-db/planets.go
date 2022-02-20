package planetsdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	errorsmodel "github.com/gmaschi/b2w-sw-planets/internal/models/planet/errors-model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

const (
	databaseName          = "star-wars"
	planetsCollectionName = "planets"
)

type CreatePlanetParams struct {
	Name    string `json:"name"`
	Terrain string `json:"terrain"`
	Climate string `json:"climate"`
}

// TODO: handle multiple requests to insert the same planet name

// CreatePlanet creates a new planet resource with the specified arguments
func (ms *MongoDBStore) CreatePlanet(ctx context.Context, arg CreatePlanetParams) (Planet, error) {
	var retPlanet Planet
	movies, err := getMovieAppearances(arg.Name)
	if err != nil {
		return retPlanet, fmt.Errorf("create planet: %s", err.Error())
	}

	collection := ms.mongodbClient.Database(databaseName).Collection(planetsCollectionName)

	planetToAdd := Planet{
		Name:    arg.Name,
		Terrain: arg.Terrain,
		Climate: arg.Climate,
		Movies:  movies,
	}

	res, err := collection.InsertOne(ctx, planetToAdd)
	if err != nil {
		return retPlanet, fmt.Errorf("create planet: %s", errorsmodel.FailedToInsertRecord)
	}
	objectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return retPlanet, fmt.Errorf("create planet: %s", errorsmodel.FailedToInsertRecord)
	}

	retPlanet = Planet{
		ID:      objectID,
		Name:    arg.Name,
		Terrain: arg.Terrain,
		Climate: arg.Climate,
		Movies:  movies,
	}
	return retPlanet, nil
}

// DeletePlanet deletes an existing planet from the collection based on the id
func (ms *MongoDBStore) DeletePlanet(ctx context.Context, id string) error {
	collection := ms.mongodbClient.Database(databaseName).Collection(planetsCollectionName)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("delete planet: %s", errorsmodel.InvalidID)
	}

	filter := bson.D{{"_id", objectId}}
	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("delete planet: %s", errorsmodel.CouldNotDeleteItem)
	}
	return nil
}

// GetPlanet finds a planet based on the ID
func (ms *MongoDBStore) GetPlanet(ctx context.Context, id string) (Planet, error) {
	collection := ms.mongodbClient.Database(databaseName).Collection(planetsCollectionName)
	var planet Planet

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return planet, fmt.Errorf("get planet: %s", errorsmodel.InvalidID)
	}
	filter := bson.D{{"_id", objectId}}
	err = collection.FindOne(ctx, filter).Decode(&planet)
	if err != nil {
		return planet, fmt.Errorf("get planet: %s", errorsmodel.FailedToFetchRecord)
	}
	return planet, nil
}

type ListPlanetParams struct {
	Name string `json:"name"`
}

// ListPlanets list filtered planets based on "name" query or list all of them
func (ms *MongoDBStore) ListPlanets(ctx context.Context, arg ListPlanetParams) ([]Planet, error) {
	// TODO: implement pagination
	collection := ms.mongodbClient.Database(databaseName).Collection(planetsCollectionName)

	var filter bson.D
	trimmedName := strings.TrimSpace(arg.Name)
	if trimmedName == "" {
		filter = bson.D{}
	} else {
		filter = bson.D{{"name", trimmedName}}
	}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var planets []Planet

	for cur.Next(ctx) {
		var planet Planet
		err := cur.Decode(&planet)
		if err != nil {
			return nil, fmt.Errorf("list planets: %s", errorsmodel.FailedToUnmarshalRecord)
		}
		planets = append(planets, planet)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	if len(planets) == 0 {
		return nil, errors.New(errorsmodel.PlanetDoesNotExist)
	}

	return planets, nil
}

// getMovieAppearances gets the total number of movies that a planet has appeared in
func getMovieAppearances(name string) (int, error) {
	baseUrlSearch := "https://swapi.dev/api/planets/?search="
	searchQuery := baseUrlSearch + name

	req, err := http.NewRequest(http.MethodGet, searchQuery, nil)
	if err != nil {
		return -1, errors.New(errorsmodel.FailedToFetchRecord)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, errors.New(errorsmodel.FailedToFetchRecord)
	}
	defer res.Body.Close()

	var planetInfo swApiPlanetInfo

	err = json.NewDecoder(res.Body).Decode(&planetInfo)
	if err != nil {
		return -1, errors.New(errorsmodel.FailedToUnmarshalRecord)
	}

	if planetInfo.Count != 1 || planetInfo.Results[0].Name != name {
		return -1, fmt.Errorf("%s: %s", errorsmodel.InvalidPlanetName, name)
	}

	return len(planetInfo.Results[0].Films), nil
}
