// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// PostDspHandlerFunc turns a function with the right signature into a post dsp handler
type PostDspHandlerFunc func(PostDspParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostDspHandlerFunc) Handle(params PostDspParams) middleware.Responder {
	return fn(params)
}

// PostDspHandler interface for that can handle valid post dsp params
type PostDspHandler interface {
	Handle(PostDspParams) middleware.Responder
}

// NewPostDsp creates a new http.Handler for the post dsp operation
func NewPostDsp(ctx *middleware.Context, handler PostDspHandler) *PostDsp {
	return &PostDsp{Context: ctx, Handler: handler}
}

/* PostDsp swagger:route POST /dsp postDsp

Get the optimal recipe

Get the optimal recipe

*/
type PostDsp struct {
	Context *middleware.Context
	Handler PostDspHandler
}

func (o *PostDsp) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostDspParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
