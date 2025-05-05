package db

import (
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/stdlib"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	MaxIdleConnections = 10
	MaxOpenConnections = 100
)

func NewDBConnection(connection string, serviceName string) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             1.0 * time.Second, // Slow SQL threshold
				LogLevel:                  logger.Error,      // Log level
				IgnoreRecordNotFoundError: true,              // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,              // Disable color
			},
		),
		PrepareStmt: true,
	}
	sqltrace.Register("pgx", &stdlib.Driver{}, sqltrace.WithAnalytics(true))
	sqlDB, err := sqltrace.Open("pgx", connection, sqltrace.WithServiceName(serviceName), sqltrace.WithAnalytics(true))
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(MaxIdleConnections)
	sqlDB.SetMaxOpenConns(MaxOpenConnections)

	db, err := gormtrace.Open(postgres.New(postgres.Config{Conn: sqlDB}), gormConfig,
		gormtrace.WithServiceName(serviceName),
		gormtrace.WithAnalytics(true))
	if err != nil {
		return nil, err
	}

	return db, nil
}
