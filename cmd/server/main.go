package main

import (
	"context"
	"errors"
	"fmt"
	httpservice "github.com/alenalato/purchase-cart-service/internal/http"
	"syscall"

	//bookingflowbizlogic "github.com/alenalato/nuitee-rest-demo/rest-server/internal/bizlogic/bookingflow"
	//staticdatabizlogic "github.com/alenalato/nuitee-rest-demo/rest-server/internal/bizlogic/staticdata"
	//bookingflowservice "github.com/alenalato/nuitee-rest-demo/rest-server/internal/http/bookingflow"
	//staticdataservice "github.com/alenalato/nuitee-rest-demo/rest-server/internal/http/staticdata"
	"github.com/alenalato/purchase-cart-service/internal/logger"
	//"github.com/alenalato/nuitee-rest-demo/rest-server/internal/storage/mockedmongo"
	//"github.com/alenalato/nuitee-rest-demo/rest-server/internal/storage/mockedpostgres"
	api "github.com/alenalato/purchase-cart-service/internal/api/go"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	defer logger.Log.Sync()

	httpAddress := fmt.Sprintf(
		"%s:%s",
		os.Getenv("HTTP_LISTEN_HOST"),
		os.Getenv("HTTP_LISTEN_PORT"),
	)

	httpService := httpservice.NewPurchaseCartAPIService()
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
