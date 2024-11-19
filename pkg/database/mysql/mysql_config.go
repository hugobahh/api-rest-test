package mysql

import (
	"time"
)

type DataSource struct {
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
