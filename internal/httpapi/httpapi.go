package httpapi

import (
	"context"
	"net/http"

	"github.com/aryanugroho/blogapp/internal/httpapi/middleware"
	"github.com/aryanugroho/blogapp/internal/httpapi/server"
)

func Serve(ctx context.Context, addr string, deps Dependencies) error {
	srv := server.New()
	srv.Register(Routes(deps))
	return srv.Serve(ctx, addr)
}

// Routes returns an array of server.Route with configured middleware for each.
func Routes(deps Dependencies) []server.Route {
	instrumentAPI := middleware.InstrumentStatsD()
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/ping",
			Handler: instrumentAPI(Ping()),
		},
	}
}

// Dependencies contains dependencies required for server and handlers.
type Dependencies struct{}
