package redis

import "time"

type Config struct {
	Addr            string
	UserName        string
	Password        string
	Database        int
	SentinelAddrs   []string
	MasterName      string
	Sentinels       bool
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}
