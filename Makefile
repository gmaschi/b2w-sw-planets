test:
	go test -v -cover ./...

mock:
	mockgen -package mockedstore -destination internal/services/datastore/mocks/mongodb/planets-db/mockedStore.go github.com/gmaschi/b2w-sw-planets/internal/services/datastore/mongodb/planets-db Store

.PHONY: test mock