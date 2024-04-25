package grpc

import "sync"

var ConfigSingleton Config
var TransportSingleton Transport
var Once sync.Once
