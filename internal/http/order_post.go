package http

import (
	"context"
	"errors"
	api "github.com/alenalato/purchase-cart-service/internal/api/go"
	"github.com/alenalato/purchase-cart-service/internal/common"
	"github.com/alenalato/purchase-cart-service/internal/logger"
)

// V1OrderPost creates an order
func (s *PurchaseCartAPIService) V1OrderPost(
	ctx context.Context,
	createOrderRequest api.CreateOrderRequest,
) (api.ImplResponse, error) {
	// create order
	order, createErr := s.orderManager.CreateOrder(
		ctx,
		s.converter.fromApiCreateOrderRequestToModel(ctx, createOrderRequest),
	)
	if createErr != nil {
		var createErrCommon common.Error
		if errors.As(createErr, &createErrCommon) {
			// validation errors are 400, not found are 422
			if createErrCommon.GetType() == common.ErrTypeInvalidArgument {
				return api.Response(400, api.Error{
					Code:    errorCodeInvalidArgument.String(),
					Message: createErrCommon.Error(),
				}), nil
			} else if createErrCommon.GetType() == common.ErrTypeNotFound {
				return api.Response(422, api.Error{
					Code:    errorCodeUnprocessableEntity.String(),
					Message: createErrCommon.Error(),
				}), nil
			}
		}

		// other errors are internal
		return api.Response(500, api.Error{
			Code: errorCodeInternal.String(),
		}), nil
	}
	if order == nil {
		logger.Log.Error("Unexpected nil order w/o error")

		return api.Response(500, api.Error{
			Code: errorCodeInternal.String(),
		}), nil
	}

	// convert order to response
	resBody, convertErr := s.converter.fromModelOrderToApi(ctx, order)
	if convertErr != nil {
		return api.Response(500, api.Error{
			Code: errorCodeInternal.String(),
		}), nil
	}

	return api.Response(201, resBody), nil
}
