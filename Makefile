#################################################################################
# BUILD
#################################################################################
gen:
	go generate ./...

#################################################################################
# RUN
#################################################################################
run:
	go run cmd/alex-api/main.go

scrape:
	go run cmd/scrapedsp/main.go

up:
	go mod vendor
	docker compose --file ./build/docker-compose.yml --project-directory . up --build -d
	rm -rf vendor

down:
	docker compose --file ./build/docker-compose.yml --project-directory . down

down-volumes:
	docker compose --file ./build/docker-compose.yml --project-directory . down --volumes

#################################################################################
# TEST
#################################################################################
test:
	go test -cover ./...
	golangci-lint run ./...

test-coverage:
	go test -coverpkg ./internal/... -coverprofile coverage.out ./... && go tool cover -html=coverage.out
