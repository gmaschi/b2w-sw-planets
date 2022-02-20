docker-run-db:
	docker run --name sw-planets -p 27017:27017 -d mongo

docker-rm-db:
	docker rm -f sw-planets

test:
	go test -v -cover ./...

mock:
	mockgen -package mockedstore -destination internal/services/datastore/mocks/mongodb/planets-db/mockedStore.go github.com/gmaschi/b2w-sw-planets/internal/services/datastore/mongodb/planets-db Store

.PHONY: docker-run-db docker-rm-db test mock