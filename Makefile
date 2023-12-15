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

#################################################################################
# TEST
#################################################################################
test:
	go test -cover ./...
	golangci-lint run ./...

test-coverage:
	go test -coverpkg ./internal/... -coverprofile coverage.out ./... && go tool cover -html=coverage.out
