package planetcontroller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	planetsfactory "github.com/gmaschi/b2w-sw-planets/internal/factories/planets-factory"
	planetmodel "github.com/gmaschi/b2w-sw-planets/internal/models/planet"
	mockedstore "github.com/gmaschi/b2w-sw-planets/internal/services/datastore/mocks/mongodb/planets-db"
	planetsdb "github.com/gmaschi/b2w-sw-planets/internal/services/datastore/mongodb/planets-db"
	"github.com/gmaschi/b2w-sw-planets/pkg/tools/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestCreate tests the Create planet controller
func TestCreate(t *testing.T) {
	planet := randomPlanet()

	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mockedstore.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"name":    planet.Name,
				"terrain": planet.Terrain,
				"climate": planet.Climate,
			},
			buildStubs: func(store *mockedstore.MockStore) {
				arg := planetsdb.CreatePlanetParams{
					Name:    planet.Name,
					Terrain: planet.Terrain,
					Climate: planet.Climate,
				}
				store.EXPECT().
					CreatePlanet(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(planet, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchCreate(t, recorder.Body, planet)
			},
		},
		{
			name: "BadRequest",
			body: map[string]interface{}{
				"name":    "",
				"terrain": planet.Terrain,
				"climate": planet.Climate,
			},
			buildStubs: func(store *mockedstore.MockStore) {
				store.EXPECT().
					CreatePlanet(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: map[string]interface{}{
				"name":    planet.Name,
				"terrain": planet.Terrain,
				"climate": planet.Climate,
			},
			buildStubs: func(store *mockedstore.MockStore) {
				store.EXPECT().
					CreatePlanet(gomock.Any(), gomock.Any()).
					Times(1).
					Return(planetsdb.Planet{}, mongo.ErrClientDisconnected)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockedstore.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server, err := planetsfactory.New(store)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/planets"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			// check response
			server.Router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

// TestPlanet tests the Planet controller
func TestPlanet(t *testing.T) {
	planet := randomPlanet()
	testCases := []struct {
		name          string
		planetID      string
		buildStubs    func(store *mockedstore.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			planetID: planet.ID.Hex(),
			buildStubs: func(store *mockedstore.MockStore) {
				store.EXPECT().
					GetPlanet(gomock.Any(), gomock.Eq(planet.ID.Hex())).
					Times(1).
					Return(planet, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPlanet(t, recorder.Body, planet)
			},
		},
		{
			name:     "BadRequest",
			planetID: "inval!d-$ID#",
			buildStubs: func(store *mockedstore.MockStore) {
				store.EXPECT().
					GetPlanet(gomock.Any(), gomock.Eq(planet.ID.Hex())).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "NotFound",
			planetID: planet.ID.Hex(),
			buildStubs: func(store *mockedstore.MockStore) {
				store.EXPECT().
					GetPlanet(gomock.Any(), gomock.Eq(planet.ID.Hex())).
					Times(1).
					Return(planetsdb.Planet{}, mongo.ErrNoDocuments)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			planetID: planet.ID.Hex(),
			buildStubs: func(store *mockedstore.MockStore) {
				store.EXPECT().
					GetPlanet(gomock.Any(), gomock.Eq(planet.ID.Hex())).
					Times(1).
					Return(planetsdb.Planet{}, mongo.ErrClientDisconnected)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockedstore.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server, err := planetsfactory.New(store)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/planets/%s", tc.planetID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// check response
			server.Router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

// TestDelete tests the Planet controller
func TestDelete(t *testing.T) {
	planet := randomPlanet()
	testCases := []struct {
		name          string
		planetID      string
		buildStubs    func(store *mockedstore.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			planetID: planet.ID.Hex(),
			buildStubs: func(store *mockedstore.MockStore) {
				store.EXPECT().
					DeletePlanet(gomock.Any(), gomock.Eq(planet.ID.Hex())).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:     "BadRequest",
			planetID: "inval!d-$ID#",
			buildStubs: func(store *mockedstore.MockStore) {
				store.EXPECT().
					DeletePlanet(gomock.Any(), gomock.Eq(planet.ID.Hex())).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			planetID: planet.ID.Hex(),
			buildStubs: func(store *mockedstore.MockStore) {
				store.EXPECT().
					DeletePlanet(gomock.Any(), gomock.Eq(planet.ID.Hex())).
					Times(1).
					Return(mongo.ErrClientDisconnected)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockedstore.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server, err := planetsfactory.New(store)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/planets/%s", tc.planetID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			// check response
			server.Router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

// TestList tests the Planet controller
func TestList(t *testing.T) {
	n := 5
	planetsSlice := make([]planetsdb.Planet, 0, 10)
	var planet planetsdb.Planet
	for i := 0; i < n; i++ {
		planet = randomPlanet()
		planetsSlice = append(planetsSlice, planet)
	}
	testCases := []struct {
		name     string
		listData struct {
			name string
		}
		buildStubs    func(store *mockedstore.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OKUnfiltered",
			listData: struct {
				name string
			}{
				name: "",
			},
			buildStubs: func(store *mockedstore.MockStore) {
				listArgs := planetsdb.ListPlanetParams{
					Name: "",
				}
				store.EXPECT().
					ListPlanets(gomock.Any(), gomock.Eq(listArgs)).
					Times(1).
					Return(planetsSlice, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchList(t, recorder.Body, planetsSlice)
			},
		},
		{
			name: "OKFiltered",
			listData: struct {
				name string
			}{
				name: "randomName",
			},
			buildStubs: func(store *mockedstore.MockStore) {
				listArgs := planetsdb.ListPlanetParams{
					Name: "randomName",
				}
				store.EXPECT().
					ListPlanets(gomock.Any(), gomock.Eq(listArgs)).
					Times(1).
					Return(planetsSlice, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchList(t, recorder.Body, planetsSlice)
			},
		},
		{
			name: "BadRequest",
			listData: struct {
				name string
			}{
				name: "invalid-n@me#",
			},
			buildStubs: func(store *mockedstore.MockStore) {
				listArgs := planetsdb.ListPlanetParams{
					Name: "invalid-n@me#",
				}
				store.EXPECT().
					ListPlanets(gomock.Any(), gomock.Eq(listArgs)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			listData: struct {
				name string
			}{
				name: "",
			},
			buildStubs: func(store *mockedstore.MockStore) {
				listArgs := planetsdb.ListPlanetParams{
					Name: "",
				}
				store.EXPECT().
					ListPlanets(gomock.Any(), gomock.Eq(listArgs)).
					Times(1).
					Return(nil, mongo.ErrNoDocuments)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			listData: struct {
				name string
			}{
				name: "",
			},
			buildStubs: func(store *mockedstore.MockStore) {
				listArgs := planetsdb.ListPlanetParams{
					Name: "",
				}
				store.EXPECT().
					ListPlanets(gomock.Any(), gomock.Eq(listArgs)).
					Times(1).
					Return(nil, mongo.ErrClientDisconnected)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockedstore.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server, err := planetsfactory.New(store)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			var url string
			if tc.listData.name == "" {
				url = fmt.Sprintf("/planets")
			} else {
				url = fmt.Sprintf("/planets?name=%s", tc.listData.name)
			}
			fmt.Println("url: ", url, "-")
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// check response
			server.Router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomPlanet() planetsdb.Planet {
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

	return planetsdb.Planet{
		ID:      primitive.NewObjectID(),
		Name:    planets[planetIndex].name,
		Terrain: random.String(6),
		Climate: random.String(5),
		Movies:  planets[planetIndex].movies,
	}
}

func requireBodyMatchCreate(t *testing.T, body *bytes.Buffer, planet planetsdb.Planet) {
	var gotPlanet planetmodel.CreateResponse

	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	err = json.Unmarshal(data, &gotPlanet)
	require.NoError(t, err)

	require.Equal(t, planet.ID, gotPlanet.ID)
	require.Equal(t, planet.Name, gotPlanet.Name)
	require.Equal(t, planet.Terrain, gotPlanet.Terrain)
	require.Equal(t, planet.Climate, gotPlanet.Climate)
	require.Equal(t, planet.Movies, gotPlanet.Movies)
}

func requireBodyMatchPlanet(t *testing.T, body *bytes.Buffer, planet planetsdb.Planet) {
	var gotPlanet planetmodel.CreateResponse

	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	err = json.Unmarshal(data, &gotPlanet)
	require.NoError(t, err)

	require.Equal(t, planet.ID, gotPlanet.ID)
	require.Equal(t, planet.Name, gotPlanet.Name)
	require.Equal(t, planet.Terrain, gotPlanet.Terrain)
	require.Equal(t, planet.Climate, gotPlanet.Climate)
	require.Equal(t, planet.Movies, gotPlanet.Movies)
}

func requireBodyMatchList(t *testing.T, body *bytes.Buffer, planets []planetsdb.Planet) {
	var gotPlanets []planetmodel.CreateResponse

	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	err = json.Unmarshal(data, &gotPlanets)
	require.NoError(t, err)

	require.Len(t, gotPlanets, len(planets))

	for i, planet := range gotPlanets {
		require.Equal(t, planets[i].ID, planet.ID)
		require.Equal(t, planets[i].Name, planet.Name)
		require.Equal(t, planets[i].Terrain, planet.Terrain)
		require.Equal(t, planets[i].Climate, planet.Climate)
		require.Equal(t, planets[i].Movies, planet.Movies)
	}
}
