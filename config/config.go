package config

import (
	"fmt"
	"github.com/imtanmoy/logx"
	"strings"

	"github.com/spf13/viper"
)

// Config contains env variables
type Config struct {
	ENVIRONMENT           string `mapstructure:"environment"`
	DEBUG                 bool   `mapstructure:"debug"`
	JwtSecretKey          string `mapstructure:"jwt_secret_key"`
	JwtAccessTokenExpires int    `mapstructure:"jwt_access_token_expires"`
	SERVER                Server
	DB                    DB
}

type Server struct {
	HOST string `mapstructure:"host"`
	PORT int    `mapstructure:"port"`
}

type DB struct {
	HOST     string `mapstructure:"host"`
	PORT     int    `mapstructure:"port"`
	USERNAME string `mapstructure:"username"`
	PASSWORD string `mapstructure:"password"`
	DBNAME   string `mapstructure:"db_name"`
}

// Conf is global configuration file
var Conf Config

// InitConfig initialze the Conf
func InitConfig() {
	config, err := initViper("./")
	if err != nil {
		logx.Fatal(err)
		return
	}
	Conf = *config
}

func initViper(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(configPath)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return &config, nil
}
