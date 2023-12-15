//go:generate mockgen -source=db.go -destination=db_mock_test.go -package=server
package server

import (
	"alex-api/internal/data"
)

type DB interface {
	GetTwitterMedia(skip *int64, limit *int64) ([]data.TwitterMedia, error)
	GetBDORecipes(skip *int64, limit *int64) ([]data.Recipe, error)
	GetDSPRecipes(skip *int64, limit *int64) ([]data.Recipe, error)
	GetPageView(domain string, path string) (data.PageView, error)
	CreatePageView(pageView data.PageView) error
	UpdatePageView(pageView data.PageView) error
}
