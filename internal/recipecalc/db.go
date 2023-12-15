//go:generate mockgen -source=db.go -destination=db_mock_test.go -package=recipecalc
package recipecalc

import (
	"alex-api/internal/data"
)

type DB interface {
	GetDSPRecipes(skip *int64, limit *int64) ([]data.Recipe, error)
	GetBDORecipes(skip *int64, limit *int64) ([]data.Recipe, error)
}
