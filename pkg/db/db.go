package db

import (
	"backend_test/pkg/config"
	"errors"
	"fmt"
	golog "log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/labstack/gommon/log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Handler struct {
	DB *gorm.DB
	Tx *gorm.DB
}

func Init() *Handler {
	return InitWithPort(int(config.Data.Db.Port))
}

func InitWithPort(port int) *Handler {
	gormlog := logger.New(
		golog.New(os.Stdout, "\r\n", golog.LstdFlags),
		logger.Config{
			LogLevel:                  logger.Info,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
		},
	)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		config.Data.Db.Host,
		config.Data.Db.Username,
		config.Data.Db.Password,
		config.Data.Db.Name,
		port)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormlog})
	if err != nil {
		log.Fatal("Failed to open database connection: ", err)
	}
	return &Handler{DB: conn, Tx: conn}
}

func Migrate(handler *Handler) {
	if handler == nil {
		log.Fatal("Database handler is nil, cannot run migrations")
	}
	instance, err := handler.DB.DB()
	if err != nil {
		log.Fatal("Get db instance for migrations error: ", err)
	}
	driver, err := migratepostgres.WithInstance(instance, &migratepostgres.Config{})
	if err != nil {
		log.Fatal("Get postgres driver for migrations error: ", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatal("Create new migration instance error: ", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal("Migration up error: ", err)
		}
	}
}
