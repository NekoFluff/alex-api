//go:generate mockgen -source=db.go -destination=db_mock_test.go -package=server
package server

import "alex-api/internal/data"

type DB interface {
	GetTwitterMedia(skip *int64, limit *int64) ([]data.TwitterMedia, error)
	GetRecipes(skip *int64, limit *int64) ([]data.Recipe, error)
}
