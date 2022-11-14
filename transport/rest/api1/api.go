package api1

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aryanugroho/blogapp/app"
	"github.com/aryanugroho/blogapp/config"
	"github.com/aryanugroho/blogapp/internal/httpapi/middleware"
	"github.com/aryanugroho/blogapp/internal/logger"
	"github.com/gorilla/mux"
)

type Routes struct {
	Root    *mux.Router // ''
	ApiRoot *mux.Router // 'api/v1'

	Post *mux.Router // 'api/v1/post'
}

type API struct {
	app        *app.App
	BaseRoutes *Routes
}

func Init(ctx context.Context, srv *app.App) (*API, error) {
	conf := config.All()

	apiv1 := API{
		app: srv,
		BaseRoutes: &Routes{
			Root: mux.NewRouter(),
		},
	}
	apiv1.BaseRoutes.Root.Use(middleware.CorrelationLoggingMiddleware)
	apiv1.BaseRoutes.Root.Use(middleware.AuthAPI)
	apiv1.BaseRoutes.Root.Use(middleware.APM)

	apiv1.BaseRoutes.ApiRoot = apiv1.BaseRoutes.Root.PathPrefix(fmt.Sprintf("/%s/api/v1", conf.App.ApiPrefix)).Subrouter()

	// init routes
	apiv1.InitAPI()

	server := http.Server{
		Handler: apiv1.BaseRoutes.Root,
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
	}

	server.ListenAndServe()

	logger.Info(ctx, "server started on port "+fmt.Sprintf("%d", conf.Server.Port))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	logger.Info(ctx, "shutting down")
	server.Shutdown(ctx)
	logger.Info(ctx, "all server stopped!")

	return &apiv1, nil
}

func callOnInterrupt(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
	cancel()
}
