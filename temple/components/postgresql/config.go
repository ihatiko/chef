package postgresql

import "time"

type Config struct {
	Host               string
	Port               int
	User               string
	Password           string
	Database           string
	SSLMode            string
	PgDriver           string
	AutoMigrate        bool
	QueryExecMode      string
	MaxOpenConnections int
	ConnMaxLifetime    time.Duration
	MaxIdleConnections int
	ConnMaxIdleTime    time.Duration
}
