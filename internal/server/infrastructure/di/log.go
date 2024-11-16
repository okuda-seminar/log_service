package di

import (
	"go.uber.org/dig"

	"log_service/internal/server/domain"
	"log_service/internal/server/infrastructure/mysql/db"
	"log_service/internal/server/infrastructure/mysql/repository"
	"log_service/internal/server/infrastructure/rabbitmq"
	"log_service/internal/server/presentation"
	"log_service/internal/server/usecase"
)

func BuildLogContainer() (*dig.Container, error) {
	container := dig.New()
	if err := container.Provide(db.Connect); err != nil {
		return nil, err
	}

	if err := container.Provide(repository.NewLogRepository, dig.As(new(domain.ILogRepository))); err != nil {
		return nil, err
	}

	if err := container.Provide(usecase.NewInsertLogUseCase, dig.As(new(usecase.IInsertLogUseCase))); err != nil {
		return nil, err
	}

	if err := container.Provide(usecase.NewListLogsUseCase, dig.As(new(usecase.IListLogsUseCase))); err != nil {
		return nil, err
	}

	if err := container.Provide(rabbitmq.Connect); err != nil {
		return nil, err
	}

	if err := container.Provide(presentation.NewAMQPLogHandler); err != nil {
		return nil, err
	}

	if err := container.Provide(presentation.NewHttpLogHandler); err != nil {
		return nil, err
	}

	return container, nil
}
