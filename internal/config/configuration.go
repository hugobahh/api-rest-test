package config

import (
	"context"
	"errors"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type Configuration struct {
	Port            int           `env:"PORT,default=5003"`
	Version         string        `env:"VERSION,default=1.0"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,default=5s"`
	MySQLDataSource MySQLDataSource
}

type MySQLDataSource struct {
	Host                   string        `env:"MYSQL_HOST,required"`
	User                   string        `env:"MYSQL_USER,required"`
	Password               string        `env:"MYSQL_PASSWORD,required"`
	Name                   string        `env:"MYSQL_NAME,required"`
	Port                   int           `env:"MYSQL_PORT,default=1433"`
	ReadOnly               bool          `env:"MYSQL_READONLY,default=true"`
	Timeout                time.Duration `env:"MYSQL_TIMEOUT,required"`
	MaxOpenConnections     int           `env:"MYSQL_MAX_OPEN_CONNECTIONS,default=25"`
	MaxIdleConnections     int           `env:"MYSQL_MAX_IDLE_CONNECTIONS,default=25"`
	MaxLifetimeConnections time.Duration `env:"MYSQL_MAX_LIFETIME_CONNECTIONS,default=5m"`
	TimeoutQuery           time.Duration `env:"MYSQL_TIMEOUT_QUERY,required"`
}

func NewConfiguration() (*Configuration, error) {
	env := new(Configuration)
	err := envconfig.Process(context.Background(), env, envconfig.MutatorFunc(func(ctx context.Context, originalKey,
		resolvedKey, originalValue, currentValue string) (newValue string, stop bool, err error) {
		if currentValue == "" {
			return "", true, errors.New("can not be an empty value " + originalKey)
		}
		return currentValue, false, nil
	}))
	return env, err
}
