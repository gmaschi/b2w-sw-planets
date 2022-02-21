# b2w-sw-planets-api

### Configurar o container com o MongoDB

- Para subir o container: make docker-run-db
- Para remover o container: make docker-rm-db

### Uso da API

- Rota: /v1/planets

#### Adicionar um planeta

- POST /v1/planets

#### Listar planetas

- GET /v1/planets (query "name" opcional para filtrar por nome)

#### Encontrar planeta por ID

- GET /v1/planets/:id

#### Remover planeta por ID

- DELETE /v1/planets/:id

#### Testes

- Para rodar os testes: make test