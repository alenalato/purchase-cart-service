package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/alenalato/purchase-cart-service/internal/bizlogic/order"
	httpservice "github.com/alenalato/purchase-cart-service/internal/http"
	"github.com/alenalato/purchase-cart-service/internal/storage/mongodb"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"syscall"

	api "github.com/alenalato/purchase-cart-service/internal/api/go"
	"github.com/alenalato/purchase-cart-service/internal/bizlogic/pricing"
	"github.com/alenalato/purchase-cart-service/internal/logger"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	defer func(Log *zap.SugaredLogger) {
		_ = Log.Sync()
	}(logger.Log)

	envErr := godotenv.Load(fmt.Sprintf("%s/.env", os.Getenv("CONFIG_DIR")))
	if envErr != nil {
		logger.Log.Fatalf("Error loading .env file: %v", envErr)
	}

	httpAddress := fmt.Sprintf(
		"%s:%s",
		os.Getenv("HTTP_LISTEN_HOST"),
		os.Getenv("HTTP_LISTEN_PORT"),
	)

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

	vatManager, vatErr := pricing.NewFixedVat(10.0)
	if vatErr != nil {
		logger.Log.Fatalf("could not initialize VAT manager: %v", vatErr)
	}

	pricingManager := pricing.NewLogic(vatManager, mongoDbStorage)
	orderManager := order.NewLogic(pricingManager, mongoDbStorage)

	httpService := httpservice.NewPurchaseCartAPIService(orderManager)

	apiController := api.NewPurchaseCartAPIController(httpService)
	router := api.NewRouter(apiController)

	srv := &http.Server{
		Addr:    httpAddress,
		Handler: router,
	}

	closeServer := make(chan struct{})

	go func() {
		if srvErr := srv.ListenAndServe(); !errors.Is(srvErr, http.ErrServerClosed) {
			logger.Log.Errorf("could not listen HTTP(%s): %v", httpAddress, srvErr)
			close(closeServer)
		}
	}()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGTERM, os.Interrupt)
		<-sigint
		close(closeServer)
	}()

	logger.Log.Infof("http server initialized on %s", httpAddress)

	<-closeServer

	logger.Log.Infof("waiting for HTTP server to close")

	if shutDownErr := srv.Shutdown(context.Background()); shutDownErr != nil {
		logger.Log.Errorf("could not shutdown HTTP server: %v", shutDownErr)
	}
}
