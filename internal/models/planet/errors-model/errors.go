package errorsmodel

const (
	FailedToFetchRecord  = "failed to fetch record"
	FailedToInsertRecord = "failed to insert record"

	InvalidPlanetName = "invalid planet name"
	InvalidID         = "invalid ID"

	PlanetAlreadyExists = "planet already exists"
	PlanetDoesNotExist  = "planet does not exist"

	CouldNotDeleteItem = "could not delete item"

	FailedToUnmarshalRecord = "failed to unmarshal record"
	FailedToMarshalItem     = "failed to marshal item"
)
