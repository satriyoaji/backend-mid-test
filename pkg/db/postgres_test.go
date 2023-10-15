package db

import (
	"backend_test/pkg/config"
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {
	err := config.LoadWithPath("./../../configs/config-test.yml")
	if err != nil {
		log.Fatal("Load config error: ", err)
	}
	m.Run()
}

func TestPostgres(t *testing.T) {
	ctx := context.Background()

	var dbname = config.Data.Db.Name
	var user = config.Data.Db.Username
	var password = config.Data.Db.Password

	port, err := nat.NewPort("tcp", fmt.Sprintf("%d", config.Data.Db.Port))
	require.NoError(t, err)

	container, err := SetupPostgres(ctx,
		WithPort(port.Port()),
		WithInitialDatabase(user, password, dbname),
		WithWaitStrategy(wait.ForLog(PostgresReadyMsg).WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	containerPort, err := container.MappedPort(ctx, port)
	assert.NoError(t, err)

	host, err := container.Host(ctx)
	assert.NoError(t, err)

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, containerPort.Port(), user, password, dbname)

	// perform assertions
	db, err := sql.Open("postgres", connStr)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	defer db.Close()

	result, err := db.Exec("CREATE TABLE IF NOT EXISTS test (id int, name varchar(255));")
	assert.NoError(t, err)
	assert.NotNil(t, result)

	result, err = db.Exec("INSERT INTO test (id, name) VALUES (1, 'test');")
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestContainerWithWaitForSQL(t *testing.T) {
	var dbname = config.Data.Db.Name
	var user = config.Data.Db.Username
	var password = config.Data.Db.Password

	ctx := context.Background()

	var port = fmt.Sprintf("%d/tcp", config.Data.Db.Port)
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port.Port(), dbname)
	}

	t.Run("default query", func(t *testing.T) {
		container, err := SetupPostgres(
			ctx,
			WithPort(port),
			WithInitialDatabase(user, password, dbname),
			WithWaitStrategy(wait.ForSQL(nat.Port(port), "postgres", dbURL)))
		require.NoError(t, err)
		require.NotNil(t, container)
	})
	t.Run("custom query", func(t *testing.T) {
		container, err := SetupPostgres(
			ctx,
			WithPort(port),
			WithInitialDatabase(user, password, dbname),
			WithWaitStrategy(wait.ForSQL(nat.Port(port), "postgres", dbURL).WithStartupTimeout(time.Second*5).WithQuery("SELECT 10")),
		)
		require.NoError(t, err)
		require.NotNil(t, container)
	})
	t.Run("custom bad query", func(t *testing.T) {
		container, err := SetupPostgres(
			ctx,
			WithPort(port),
			WithInitialDatabase(user, password, dbname),
			WithWaitStrategy(wait.ForSQL(nat.Port(port), "postgres", dbURL).WithStartupTimeout(time.Second*5).WithQuery("SELECT 'a' from b")),
		)
		require.Error(t, err)
		require.Nil(t, container)
	})
}
