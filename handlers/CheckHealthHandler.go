package handlers

import (
	"addi/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

//Health route returns OK
func CheckHealthHandler(operations.CheckHealthParams) middleware.Responder {
	return operations.NewCheckHealthOK().WithPayload("OK")
}
