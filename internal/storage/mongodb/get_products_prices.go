package mongodb

import (
	"context"
	"github.com/alenalato/purchase-cart-service/internal/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

func (m *MongoDB) GetProductPrices(ctx context.Context, productIds []int) (map[int]storage.ProductPrice, error) {
	// TODO handle empty productIds

	findCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.D{{"product_id", bson.D{{"$in", productIds}}}}
	cursor, findErr := m.database.Collection(ProductPriceCollection).Find(
		findCtx,
		filter,
	)
	if findErr != nil {
		return nil, findErr
	}

	var productPrices []storage.ProductPrice
	decodeErr := cursor.All(ctx, &productPrices)

	if decodeErr != nil {
		return nil, decodeErr
	}

	pricesByProduct := make(map[int]storage.ProductPrice)
	for _, prodPrice := range productPrices {
		pricesByProduct[prodPrice.ProductId] = prodPrice
	}

	return pricesByProduct, nil
}
