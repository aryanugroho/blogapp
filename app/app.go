package app

import (
	"context"

	"github.com/aryanugroho/blogapp/config"
	"github.com/aryanugroho/blogapp/infrastructure"
)

type App struct {
	Infrastructure infrastructure.Infrastructure
}

type UseCase interface {
	CreatePost(ctx context.Context, payload *PostPayload) (*PostPayload, error)
}

func NewApp(ctx context.Context) (*App, error) {
	infrastructure, err := infrastructure.NewInfra(ctx, *config.All())
	if err != nil {
		return nil, err
	}

	app := &App{
		Infrastructure: infrastructure,
	}

	return app, nil
}
