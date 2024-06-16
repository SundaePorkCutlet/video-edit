package db

import (
	"context"
	"database/sql"
	"fmt"
	"stockfoilo_test/internal/config"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

var dbPool *sql.DB

const (
	MaxOpenConnections    = 20
	MaxIdleConnections    = 20
	ConnectionMaxIdleTime = 5 * time.Minute // minutes
)

func InitDbConnection(appConfig *config.AppConfig) {

	dbAddress := fmt.Sprintf("%s:%d", appConfig.DbConfig.Host, appConfig.DbConfig.Port)

	cfg := mysql.Config{
		User:                 appConfig.DbConfig.User,
		Passwd:               appConfig.DbConfig.Password,
		Net:                  "tcp",
		Addr:                 dbAddress,
		DBName:               appConfig.DbConfig.DbName,
		AllowNativePasswords: true,
	}

	var err error
	dbPool, err = sql.Open(appConfig.DbConfig.Driver, cfg.FormatDSN())
	if err != nil {
		log.Error().Msgf("InitDbConnection:: database connection error: %s", err.Error())
		log.Panic().Msgf("InitDbConnection:: database connection error: %s", err.Error())
		return
	}

	dbPool.SetMaxOpenConns(MaxOpenConnections)
	dbPool.SetMaxIdleConns(MaxIdleConnections)
	dbPool.SetConnMaxIdleTime(ConnectionMaxIdleTime)

	dbPool.Ping()

	log.Info().Msg("InitDbConnection:: database connected successfully!!!")
}

func GetDbConnection(ctx context.Context) *DbCtx {
	return &DbCtx{
		DB:  dbPool,
		Ctx: ctx,
	}
}
