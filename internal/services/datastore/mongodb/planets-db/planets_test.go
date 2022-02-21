package planetsdb

import (
	"context"
	"fmt"
	errorsmodel "github.com/gmaschi/b2w-sw-planets/internal/models/planet/errors-model"
	"github.com/gmaschi/b2w-sw-planets/pkg/tools/random"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func createRandomPlanet(t *testing.T) Planet {
	planets := []struct {
		name   string
		movies int
	}{
		{
			name:   "Tatooine",
			movies: 5,
		},
		{
			name:   "Kamino",
			movies: 1,
		},
		{
			name:   "Stewjon",
			movies: 0,
		},
		{
			name:   "Utapau",
			movies: 1,
		},
		{
			name:   "Alderaan",
			movies: 2,
		},
	}

	rand.Seed(time.Now().UnixNano())
	planetIndex := rand.Intn(len(planets))
	arg := CreatePlanetParams{
		Name:    planets[planetIndex].name,
		Terrain: random.String(6),
		Climate: random.String(5),
	}

	planet, err := testStore.CreatePlanet(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Name, planet.Name)
	require.Equal(t, arg.Terrain, planet.Terrain)
	require.Equal(t, arg.Climate, planet.Climate)
	require.Equal(t, planets[planetIndex].movies, planet.Movies)
	return planet
}

func TestCreatePlanet(t *testing.T) {
	createRandomPlanet(t)
}

func TestGetPlanet(t *testing.T) {
	planet := createRandomPlanet(t)
	gotPlanet, err := testStore.GetPlanet(context.Background(), planet.ID.Hex())
	require.NoError(t, err)

	require.Equal(t, planet.Name, gotPlanet.Name)
	require.Equal(t, planet.Terrain, gotPlanet.Terrain)
	require.Equal(t, planet.Climate, gotPlanet.Climate)
	require.Equal(t, planet.Movies, gotPlanet.Movies)
}

func TestDeletePlanet(t *testing.T) {
	planet := createRandomPlanet(t)
	err := testStore.DeletePlanet(context.Background(), planet.ID.Hex())
	require.NoError(t, err)

	deletedPlanet, err := testStore.GetPlanet(context.Background(), planet.ID.Hex())
	require.Error(t, err)
	require.EqualError(t, err, fmt.Errorf("get planet: %s", errorsmodel.PlanetDoesNotExist).Error())
	require.Empty(t, deletedPlanet)
}

func TestListPlanets(t *testing.T) {
	n := 5
	for i := 0; i < n; i++ {
		createRandomPlanet(t)
	}

	testCases := []struct {
		name     string
		listArgs ListPlanetParams
	}{
		{
			name: "unfilteredList",
			listArgs: ListPlanetParams{
				Name: "",
			},
		},
		{
			name: "filteredList",
			listArgs: ListPlanetParams{
				Name: "Tatooine",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			planetsList, err := testStore.ListPlanets(context.Background(), tc.listArgs)
			if len(planetsList) == 0 {
				require.Error(t, err, fmt.Errorf("list planets: %s", errorsmodel.PlanetDoesNotExist))
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, planetsList)
			}
			for _, planet := range planetsList {
				require.NotEmpty(t, planet)
			}
		})
	}
}
