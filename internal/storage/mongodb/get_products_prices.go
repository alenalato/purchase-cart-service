package mongodb

import (
	"context"
	"github.com/alenalato/purchase-cart-service/internal/common"
	"github.com/alenalato/purchase-cart-service/internal/logger"
	"github.com/alenalato/purchase-cart-service/internal/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

func (m *MongoDB) GetProductPrices(ctx context.Context, productIds []int) (map[int]storage.ProductPrice, error) {
	if len(productIds) == 0 {
		return nil, nil
	}

	findCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.D{{"product_id", bson.D{{"$in", productIds}}}}
	cursor, findErr := m.database.Collection(ProductPriceCollection).Find(
		findCtx,
		filter,
	)
	if findErr != nil {
		findErr = common.NewError(findErr, common.ErrTypeInternal)
		logger.Log.Errorf("Error finding product prices for ids %v: %s", productIds, findErr.Error())

		return nil, findErr
	}

	var productPrices []storage.ProductPrice
	decodeErr := cursor.All(ctx, &productPrices)

	if decodeErr != nil {
		logger.Log.Errorf("Error decoding product prices: %s", decodeErr.Error())

		return nil, common.NewError(decodeErr, common.ErrTypeInternal)
	}

	pricesByProduct := make(map[int]storage.ProductPrice)
	for _, prodPrice := range productPrices {
		pricesByProduct[prodPrice.ProductId] = prodPrice
	}

	return pricesByProduct, nil
}
