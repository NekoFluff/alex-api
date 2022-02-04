// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DSPMaterial d s p material
//
// swagger:model DSPMaterial
type DSPMaterial struct {

	// count
	Count float64 `json:"count,omitempty"`

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this d s p material
func (m *DSPMaterial) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this d s p material based on context it is used
func (m *DSPMaterial) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DSPMaterial) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DSPMaterial) UnmarshalBinary(b []byte) error {
	var res DSPMaterial
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}