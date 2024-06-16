package config

import (
	"github.com/spf13/viper"

	"github.com/rs/zerolog/log"
)

type AppConfig struct {
	DbConfig   DbConf   `mapstructure:"db"`
	LogConfig  LogConf  `mapstructure:"log"`
	HttpPort   int      `mapstructure:"http-port"`
	FileConfig FileConf `mapstructure:"file"`
}

type DbConf struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"database"`
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

type LogConf struct {
	Level      string `mapstructure:"level"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file-path"`
	MaxSize    int    `mapstructure:"max-size"`
	MaxBackups int    `mapstructure:"max-backups"`
	MaxAge     int    `mapstructure:"max-age"`
}

type FileConf struct {
	VideoPath  string `mapstructure:"video-path"`
	ConcatPath string `mapstructure:"concat-path"`
}

var appConfig AppConfig

func LoadConfig() (AppConfig, error) {

	viper.SetConfigName("app_config") // name of config file (without extension)
	viper.SetConfigType("yaml")       // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./configs/") // path to look for the config file in
	viper.AddConfigPath(".")          // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Error().Msg(err.Error()) // Handle errors reading the config file
		return appConfig, err
	}

	config := viper.Sub("config")
	err = config.Unmarshal(&appConfig)
	if err != nil {
		log.Error().Msg("" + err.Error())
		return AppConfig{}, err
	}

	return appConfig, nil
}

func GetConfig() *AppConfig {
	return &appConfig
}
