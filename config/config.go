package config

import (
	"github.com/caarlos0/env/v6"
)

type ConfigurationEnvironment struct {
	DatabaseEnvironment
	RedisEnvironment
	ElasticEnvironment
	ApplicationEnvironment
	GcpEnvironment
}

type DatabaseEnvironment struct {
	MaxIdle     int    `env:"DBMaxIdleConn"`
	MaxOpenConn int    `env:"DBMaxOpenConn"`
	DBType      string `env:"DBType"`
	DBName      string `env:"DBName"`
	DBUser      string `env:"DBUser"`
	DBPass      string `env:"DBPass"`
	DBHost      string `env:"DBHost"`
	DBSocket    string `env:"DBSocket"`
}

type RedisEnvironment struct {
	RedisAddress  string `env:"REDIS_ADDRESS"`
	RedisPassword string `env:"REDIS_PASSWORD"`
}

type ElasticEnvironment struct {
	ElasticUrl string `env:"ElasticUrl" envDefault:"http://127.0.0.1:9200"`
}

type GcpEnvironment struct {
	GcpProjectId          string `env:"NotificationProjectId"`
	GcpCredentialFileName string `env:"GcpCredentialFileName"`
	GcpStorageBucketName  string `env:"GcpStorageBucketName"`
}

type ApplicationEnvironment struct {
	Debug bool `env:"Debug"`
}

var Env = ConfigurationEnvironment{}

func LoadEnv() {
	if err := env.Parse(&Env); err != nil {
		panic(err.Error())
	}
}
