//go:generate mockgen -source=servicer.go -destination=servicer_mock_test.go -package=server
package server

type Servicer interface {
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}
