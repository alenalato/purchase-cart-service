package mongodb

import (
	"context"
	"github.com/alenalato/purchase-cart-service/internal/common"
	"github.com/alenalato/purchase-cart-service/internal/logger"
	"github.com/alenalato/purchase-cart-service/internal/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

func (m *MongoDB) CreateOrder(ctx context.Context, orderDetails storage.OrderDetails) (*storage.Order, error) {
	insertCtx, cancelInsert := context.WithTimeout(ctx, 10*time.Second)
	defer cancelInsert()

	collection := m.database.Collection(OrderCollection)

	insertRes, insertErr := collection.InsertOne(insertCtx, orderDetails)
	if insertErr != nil {
		if mongo.IsDuplicateKeyError(insertErr) {
			insertErr = common.NewError(insertErr, common.ErrTypeAlreadyExists)
		} else {
			insertErr = common.NewError(insertErr, common.ErrTypeInternal)
		}
		logger.Log.Errorf("Error creating order: %s", insertErr.Error())

		return nil, insertErr
	}

	findCtx, cancelFind := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFind()

	filter := bson.D{{Key: "_id", Value: insertRes.InsertedID}}

	var order storage.Order
	findErr := collection.FindOne(findCtx, filter).Decode(&order)
	if findErr != nil {
		logger.Log.Errorf("Error creating order: %s", findErr.Error())

		return nil, common.NewError(findErr, common.ErrTypeInternal)
	}
	order.Id = insertRes.InsertedID.(bson.ObjectID).Hex()

	return &order, nil
}
