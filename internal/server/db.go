//go:generate mockgen -source=db.go -destination=db_mock_test.go -package=server
package server

type DB interface {
	// Computed
	// Supplier(context.Context, string) (*data.Supplier, error)
	// PatchSupplier(context.Context, data.Supplier) (*data.Supplier, error)
	// PutSupplier(context.Context, data.Supplier) (*data.Supplier, error)
	// DeleteSupplier(context.Context, string) error

	// // Insurances
	// Insurance(context.Context, string) (*data.Insurance, error)
	// InsuranceByLabelID(context.Context, string) (*data.Insurance, error)
	// PatchInsurance(context.Context, data.Insurance) (*data.Insurance, error)
	// PutInsurance(context.Context, data.Insurance) (*data.Insurance, error)
	// DeleteInsurance(context.Context, string) error
}
