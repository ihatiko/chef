package grpc

import "time"

const (
	reflectDefault           = true
	healtzDefault            = true
	timeoutDefault           = time.Second * 5
	timeDefault              = time.Second * 10
	maxConnectionIdleDefault = time.Hour
	naxConnectionAgeDefault  = time.Hour
)
