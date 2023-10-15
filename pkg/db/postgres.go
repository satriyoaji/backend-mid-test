package db

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const PostgresReadyMsg = "database system is ready to accept connections"

// PostgresContainer represents the postgres container type used in the module
type PostgresContainer struct {
	testcontainers.Container
}

type PostgresContainerOption func(req *testcontainers.ContainerRequest)

func WithWaitStrategy(strategies ...wait.Strategy) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.WaitingFor = wait.ForAll(strategies...).WithDeadline(1 * time.Minute)
	}
}

func WithPort(port string) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.ExposedPorts = append(req.ExposedPorts, port)
	}
}

func WithInitialDatabase(user string, password string, dbName string) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.Env["POSTGRES_USER"] = user
		req.Env["POSTGRES_PASSWORD"] = password
		req.Env["POSTGRES_DB"] = dbName
	}
}

// SetupPostgres creates an instance of the postgres container type
func SetupPostgres(ctx context.Context, opts ...PostgresContainerOption) (*PostgresContainer, error) {
	skipReaper, _ := strconv.ParseBool(os.Getenv("TESTCONTAINERS_RYUK_DISABLED"))
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13-alpine3.13",
		Env:          map[string]string{},
		ExposedPorts: []string{},
		Cmd:          []string{"postgres", "-c", "fsync=off"},
		SkipReaper:   skipReaper,
	}

	for _, opt := range opts {
		opt(&req)
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{Container: container}, nil
}
