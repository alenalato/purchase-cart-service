package mongodb

import (
	"context"
	"fmt"
	"github.com/alenalato/purchase-cart-service/internal/common"
	"github.com/alenalato/purchase-cart-service/internal/logger"
	"github.com/alenalato/purchase-cart-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var testMongoStorage *MongoDB

const testDbName = "test"

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct test pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pull mongodb docker image
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "4.4",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	var dbClient *mongo.Client

	// exponential backoff-retry
	err = pool.Retry(func() error {
		dbClient, err = newMongoDBClient(fmt.Sprintf(
			"mongodb://%s:%s",
			resource.Container.NetworkSettings.Gateway,
			resource.GetPort("27017/tcp"),
		))
		if err != nil {
			return err
		}
		return dbClient.Ping(context.TODO(), nil)
	})

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	testMongoStorage, err = NewMongoDB(dbClient, testDbName)

	defer func() {
		// kill and remove the container
		if err = pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
		// disconnect mongodb client
		if err = dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// seed product prices
	seedTestProductPrices()

	// run tests
	m.Run()
}

func seedTestProductPrices() {
	productPrices := []interface{}{
		storage.ProductPrice{
			ProductId: 1,
			Price:     common.AsDecimal(2),
			VatClass:  1,
		},
		storage.ProductPrice{
			ProductId: 2,
			Price:     common.AsDecimal(1.5),
			VatClass:  1,
		},
		storage.ProductPrice{
			ProductId: 3,
			Price:     common.AsDecimal(3),
			VatClass:  1,
		},
		storage.ProductPrice{
			ProductId: 4,
			Price:     common.AsDecimal(1.2),
			VatClass:  1,
		},
	}

	insertCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, seedErr := testMongoStorage.GetDataBase().Collection(ProductPriceCollection).InsertMany(
		insertCtx,
		productPrices,
	)
	if seedErr != nil {
		logger.Log.Errorf("could not seed test product prices: %v", seedErr)
	}
}

func TestMongoDB_productIdUniqueIndex(t *testing.T) {
	insertCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := testMongoStorage.database.Collection(ProductPriceCollection).
		InsertOne(insertCtx, storage.ProductPrice{
			// a product with id 1 already exists in seeds
			ProductId: 1,
			Price:     common.AsDecimal(2),
			VatClass:  1,
		})

	assert.Error(t, err)
	assert.True(t, mongo.IsDuplicateKeyError(err))
}
