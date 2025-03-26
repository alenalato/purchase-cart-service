package main

import (
	"context"
	"fmt"
	"github.com/alenalato/purchase-cart-service/internal/logger"
	"github.com/alenalato/purchase-cart-service/internal/storage"
	"github.com/alenalato/purchase-cart-service/internal/storage/mongodb"
	"github.com/govalues/decimal"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"time"
)

func main() {
	defer func(Log *zap.SugaredLogger) {
		_ = Log.Sync()
	}(logger.Log)

	envErr := godotenv.Load(fmt.Sprintf("%s/.env", os.Getenv("CONFIG_DIR")))
	if envErr != nil {
		logger.Log.Fatalf("Error loading .env file: %v", envErr)
	}

	ctx := context.Background()

	mongoDbStorage, mongodbErr := mongodb.NewMongoDB(
		os.Getenv("MONGODB_URI"),
		os.Getenv("MONGODB_DATABASE"),
	)
	if mongodbErr != nil {
		logger.Log.Fatalf("could not initialize MongoDB storage: %v", mongodbErr)
	}
	defer func(mongoDbStorage *mongodb.MongoDB, ctx context.Context) {
		err := mongoDbStorage.Close(ctx)
		if err != nil {
			logger.Log.Errorf("could not close MongoDB storage: %v", err)
		}
	}(mongoDbStorage, ctx)

	insertCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	asDecimal := func(value float64) decimal.Decimal {
		dv, _ := decimal.NewFromFloat64(value)

		return dv
	}

	productPrices := []interface{}{
		storage.ProductPrice{
			ProductId: 1,
			Price:     asDecimal(2),
			VatClass:  1,
		},
		storage.ProductPrice{
			ProductId: 2,
			Price:     asDecimal(1.5),
			VatClass:  1,
		},
		storage.ProductPrice{
			ProductId: 3,
			Price:     asDecimal(3),
			VatClass:  1,
		},
		storage.ProductPrice{
			ProductId: 4,
			Price:     asDecimal(1.2),
			VatClass:  1,
		},
	}

	_, seedErr := mongoDbStorage.GetDataBase().Collection(mongodb.ProductPriceCollection).InsertMany(
		insertCtx,
		productPrices,
	)
	if seedErr != nil {
		logger.Log.Errorf("could not seed product prices: %v", seedErr)
	}
}
