package recipecalc

import (
	"alex-api/internal/data"
	"context"
	"os"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoadBDORecipes(t *testing.T) {
	os.Setenv("MONGO_CONNECTION_URI", "mongodb://localhost:27017")
	log := logrus.New().WithContext(context.TODO())

	ctrl := gomock.NewController(t)
	dbMock := NewMockDB(ctrl)
	dbMock.EXPECT().GetBDORecipes(nil, nil).Return([]data.Recipe{
		{
			Name: "recipe1",
			MarketData: &data.MarketData{
				Price: 100,
			},
		},
	}, nil)

	recipes := LoadBDORecipes(log, dbMock)
	assert.Len(t, recipes, 1)
	assert.NotEmpty(t, recipes["recipe1"])
	assert.Equal(t, 100.0, recipes["recipe1"][0].MarketData.Price)
}
