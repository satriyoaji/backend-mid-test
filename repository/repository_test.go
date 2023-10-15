package repository

import (
	"backend_test/entity"
	"backend_test/pkg/config"
	"backend_test/pkg/db"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type InitialData struct {
	persons []entity.Employee
}

var (
	conn        *gorm.DB
	repo        Repository
	initialData = getInitialData()
)

func getInitialData() InitialData {
	now := time.Now()
	data := InitialData{}
	for i := 0; i < 2; i++ {
		data.persons = append(data.persons, entity.Employee{
			ID:        uint(i) + 1,
			FirstName: fmt.Sprintf("First Employee %d", i+1),
			LastName:  fmt.Sprintf("Last %d", i+1),
			CreatedAt: now.Add(time.Duration(i) * time.Hour),
			UpdatedAt: now.Add(time.Duration(i) * time.Hour),
		})
	}

	return data
}

func resetData() {
	conn.Where("1=1").Delete(&entity.Employee{})
	insertData()
}

func insertData() {
	conn.Create(initialData.persons)
}

func migrateDatabase() {
	err := conn.AutoMigrate(
		&entity.Employee{},
	)
	if err != nil {
		log.Fatal("Auto migrate error: ", err)
	}
	insertData()
}

func startDatabase() (context.Context, *db.PostgresContainer, nat.Port) {
	ctx := context.Background()
	port, err := nat.NewPort("tcp", fmt.Sprintf("%d", config.Data.Db.Port))
	if err != nil {
		log.Fatal("Create port error: ", err)
	}
	container, err := db.SetupPostgres(ctx,
		db.WithPort(port.Port()),
		db.WithInitialDatabase(config.Data.Db.Username, config.Data.Db.Password, config.Data.Db.Name),
		db.WithWaitStrategy(wait.ForLog(db.PostgresReadyMsg).WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatal("Start test database error: ", err)
	}
	extPort, err := container.MappedPort(ctx, port)
	if err != nil {
		log.Fatal("Get mapped port error: ", err)
	}
	return ctx, container, extPort
}

func connectDatabase(port int) {
	h := db.InitWithPort(port)
	conn = h.DB
	repo = Default(h)
}

func TestMain(m *testing.M) {
	err := config.LoadWithPath("./../configs/config-test.yml")
	if err != nil {
		log.Fatal("Load config error: ", err)
	}
	ctx, container, extPort := startDatabase()
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Fatal("Stop test database error: ", err)
		}
	}()
	connectDatabase(extPort.Int())
	migrateDatabase()
	m.Run()
}

func TestWithAlias(t *testing.T) {
	col := "name"
	assert.Equal(t, "m.name", withAlias(col, "m"))
}
func TestWithPercentAround(t *testing.T) {
	val := "1"
	assert.Equal(t, "%1%", withPercentAround(val))
}
func TestWithPercentAfter(t *testing.T) {
	val := "1"
	assert.Equal(t, "1%", withPercentAfter(val))
}
func TestWithPercentBefore(t *testing.T) {
	val := "1"
	assert.Equal(t, "%1", withPercentBefore(val))
}
