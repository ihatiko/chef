package postgresql

import "time"

type Config struct {
	Port               int
	Host               string
	User               string
	Password           string
	Database           string
	SSLMode            string
	PgDriver           string
	AutoMigrate        bool
	QueryExecMode      string
	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifetime    time.Duration
	ConnMaxIdleTime    time.Duration
}
