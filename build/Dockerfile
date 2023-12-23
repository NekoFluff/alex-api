ARG BASE_IMAGE=gcr.io/distroless/static
# Step 1 build from golang alpine
FROM golang:1.21-alpine as builder
WORKDIR /app
COPY cmd ./cmd
COPY internal ./internal
COPY go.mod go.sum ./
COPY vendor ./vendor

# Build the Go binary
ARG GO_BINARY
ENV GO_BINARY ${GO_BINARY:-./cmd/alex-api}
RUN CGO_ENABLED=0 go build -mod vendor -buildvcs=false -o ./bin/app $GO_BINARY

# Step 2 from distroless
FROM $BASE_IMAGE as app
COPY --from=builder /app/bin/app /app
ENTRYPOINT ["/app"]
