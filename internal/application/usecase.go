package application

import "context"

type UseCase[IN any, OUT any] interface {
	Execute(ctx context.Context, input IN) (OUT, error)
}

type UnitUseCase[IN any] interface {
	Execute(ctx context.Context, input IN) error
}

type NullaryUseCase[OUT any] interface {
	Execute(ctx context.Context) (OUT, error)
}
