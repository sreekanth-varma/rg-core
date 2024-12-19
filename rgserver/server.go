package server

import (
	"context"

	"github.com/gin-gonic/gin"
)

var (
	ctx context.Context
)

type Options struct {
	ConfigEnabled       bool
	ConfigPreHandler    func()
	ConfigPostHandler   func()
	MQEnabled           bool
	MQPreHandler        func()
	MQPostHandler       func()
	CacheEnabled        bool
	CachePreHandler     func()
	CachePostHandler    func()
	ESEnabled           bool
	ESPreHandler        func()
	ESPostHandler       func()
	MongoEnabled        bool
	MongoPreHandler     func()
	MongoPostHandler    func()
	PGEnabled           bool
	PGPreHandler        func()
	PGPostHandler       func()
	OnReadyHandler      func()
	WebServerEnabled    bool
	WebServerPort       string
	WebServerPreHandler func(*gin.Engine)
}

func GetDefaultOptions() Options {
	return Options{
		ConfigEnabled:     true,
		ConfigPreHandler:  nil,
		ConfigPostHandler: nil,

		MQEnabled:     true,
		MQPreHandler:  nil,
		MQPostHandler: nil,

		CacheEnabled:     true,
		CachePreHandler:  nil,
		CachePostHandler: nil,

		ESEnabled:     true,
		ESPreHandler:  nil,
		ESPostHandler: nil,

		MongoEnabled: true,
		PGEnabled:    true,

		OnReadyHandler: nil,

		WebServerEnabled: true,
		WebServerPort:    "8080",
		// WebServerPreHandler: nil,
	}
}

func Start(options Options) {
	LoadConfig()
	ctx = context.TODO()

	// initConfig(&options)

	if options.MongoEnabled {
		Connect(&ctx)
		defer Disconnect(&ctx)
	}

	initCache(&options)

}
