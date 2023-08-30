package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

const (
	EnvDebug   = gin.DebugMode
	EnvRelease = gin.ReleaseMode
)

type Config struct {
	Env      string `env:"ENV" env-default:"release"`
	DbConn   string `env:"DB_CONN" env-required:"true"`
	HTTPAddr string `env:"HTTP_ADDR" env-default:":8888"`
}

func MustLoad() Config {
	path := os.Getenv("CONFIG_PATH")
	var cfg Config
	var err error
	if path != "" {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			panic(fmt.Sprintf("config file does not exist: %s", path))
		}
		err = cleanenv.ReadConfig(path, &cfg)
	} else {
		err = cleanenv.ReadEnv(&cfg)
	}
	if err != nil {
		panic(fmt.Sprintf("cannot read config: %s", err))
	}

	return cfg
}
