package app

import "context"

type App struct {
	diContainer *diContainer
}

// New creates the assembly application and initializes its dependencies.
func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

// Run starts the background assembly consumer.
func (a *App) Run(ctx context.Context) error {
	return a.diContainer.OrderConsumer(ctx).RunOrderConsumer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initMetrics,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	return a.diContainer.InitMetrics(ctx)
}
