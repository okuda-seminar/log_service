// container.go
package container

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// startRabbitMQContainer starts a RabbitMQ container and returns the connection string.
func StartRabbitMQContainer(ctx context.Context) (testcontainers.Container, string, error) {
	req := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3-management",
		ExposedPorts: []string{"5672/tcp", "15672/tcp"},
		WaitingFor:   wait.ForLog("Server startup complete"),
	}
	rabbitContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	host, err := rabbitContainer.Host(ctx)
	if err != nil {
		return nil, "", err
	}

	port, err := rabbitContainer.MappedPort(ctx, "5672")
	if err != nil {
		return nil, "", err
	}

	connectionString := fmt.Sprintf("amqp://guest:guest@%s:%s/", host, port.Port())
	return rabbitContainer, connectionString, nil
}
