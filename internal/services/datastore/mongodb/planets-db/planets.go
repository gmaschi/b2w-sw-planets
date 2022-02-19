package planetsdb

import (
	"context"
	"encoding/json"
	"fmt"
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

// CreatePlanet creates a new planet resource with the specified arguments
func (ms *MongoDBStore) CreatePlanet(ctx context.Context, arg CreatePlanetParams) (Planet, error) {
	var retPlanet Planet
	// get planet movie appearances
	movies, err := getMovieAppearances(arg.Name)
	if err != nil {
		return retPlanet, err
	}

	// get collection
	collection := ms.mongodbClient.Database(databaseName).Collection(planetsCollectionName)
	// planet to insert
	planetToAdd := Planet{
		Name:    arg.Name,
		Terrain: arg.Terrain,
		Climate: arg.Climate,
		Movies:  movies,
	}

	res, err := collection.InsertOne(ctx, planetToAdd)
	if err != nil {
		return retPlanet, err
	}
	objectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return retPlanet, err
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
	//// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", objectId}}
	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// GetPlanet finds a planet based on it's ID
func (ms *MongoDBStore) GetPlanet(ctx context.Context, id string) (Planet, error) {
	collection := ms.mongodbClient.Database(databaseName).Collection(planetsCollectionName)
	var planet Planet
	// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return planet, err
	}
	filter := bson.D{{"_id", objectId}}
	err = collection.FindOne(ctx, filter).Decode(&planet)
	if err != nil {
		return planet, err
	}
	return planet, nil
}

type ListPlanetParams struct {
	Name string `json:"name"`
	//Limit  int32  `json:"limit"`
	//Offset int32  `json:"offset"`
}

// ListPlanets list planets by pagination or based on the name search
func (ms *MongoDBStore) ListPlanets(ctx context.Context, arg ListPlanetParams) ([]Planet, error) {
	// TODO: implement pagination
	// TODO: deal with multiple entries with the same planet name
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

	var planets []Planet

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var planet Planet
		err := cur.Decode(&planet)
		if err != nil {
			return nil, err
		}
		planets = append(planets, planet)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return planets, nil
}

// getMovieAppearances gets the total number of movies that a planet has appeared in
func getMovieAppearances(name string) (int, error) {
	baseUrlSearch := "https://swapi.dev/api/planets/?search="
	searchQuery := baseUrlSearch + name

	req, err := http.NewRequest(http.MethodGet, searchQuery, nil)
	if err != nil {
		return -1, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()

	var planetInfo swApiPlanetInfo

	err = json.NewDecoder(res.Body).Decode(&planetInfo)
	if err != nil {
		return -1, err
	}

	if planetInfo.Count != 1 || planetInfo.Results[0].Name != name {
		return -1, fmt.Errorf("getMovieAppearances: invalid planet name: %s", name)
	}

	return len(planetInfo.Results[0].Films), nil
}
