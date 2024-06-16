package main

import (
	"fmt"
	"os"
	"stockfoilo_test/internal/config"

	"stockfoilo_test/internal/db"
	"stockfoilo_test/internal/router"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("config load error: %s", err.Error())
	}
	log.Debug().Msgf("config: %v", config)

	engine := gin.New()
	setupLogConfig(&config)
	db.InitDbConnection(&config)

	router.InitRoute(engine)
	port := fmt.Sprintf(":%d", config.HttpPort)

	err = engine.Run(port)
	if err != nil {
		log.Fatal().Msgf("config load error: %s", err.Error())
	}

}

func setupLogConfig(config *config.AppConfig) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = "timestamp"
	log.Logger = log.With().Caller().Logger()

	isConsoleOutput := config.LogConfig.Output == "console"
	if isConsoleOutput {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		zerolog.CallerFieldName = "caller"
	} else {
		logFilePath := config.LogConfig.FilePath
		if logFilePath == "" {
			log.Fatal().Msg("log file path empty.")
			os.Exit(-1)
		}

		log.Logger = log.Output(&lumberjack.Logger{
			Filename:   fmt.Sprintf(logFilePath, time.Now().Format("2006-01-02")),
			MaxSize:    config.LogConfig.MaxSize,    // megabytes
			MaxAge:     config.LogConfig.MaxAge,     //max no. of days to retain old log files
			MaxBackups: config.LogConfig.MaxBackups, //max no. of old log files
		})
	}
}
