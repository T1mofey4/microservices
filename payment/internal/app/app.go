package app

import (
	"context"
	"net"

	"github.com/T1mofey4/microservices/payment/internal/config"
	"github.com/T1mofey4/microservices/platform/pkg/closer"
	"github.com/T1mofey4/microservices/platform/pkg/logger"
	"google.golang.org/grpc"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// func (a *App) Run(ctx context.Context) error {
// 	return a.r
// }

// func (a *App) initDeps(ctx context.Context) error {
// 	inits := []func(context.Context) error {

// 	}

// 	return nil
// }

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJSON(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}
