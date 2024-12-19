package server

import (
	"context"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/sreekanth-varma/rg-core/rgcache"
	"github.com/sreekanth-varma/rg-core/rgdb"
	"github.com/sreekanth-varma/rg-core/rgutil"
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
	WebServerPreHandler func(*gin.Engine, *gin.RouterGroup)
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

		MongoEnabled: true,
		PGEnabled:    true,

		WebServerEnabled:    true,
		WebServerPort:       "8080",
		WebServerPreHandler: nil,
	}
}

func Start(options Options) {
	LoadConfig()
	ctx = context.TODO()

	if options.MongoEnabled {
		if err := rgdb.Connect(&ctx); err != rgutil.ErrNil {
			panic("rgserver: mongo connection failed")
		}
		defer rgdb.Disconnect(&ctx)
	}

	if err := InitCache(&options); err != rgutil.ErrNil {
		panic("rgserver: failed to start. cache connect failed")
	}
	if err := initServer(&options); err != rgutil.ErrNil {
		panic("rgserver: failed to start. webserver(CORS) init failed")
	}

}

func InitCache(options *Options) rgutil.Err {
	if !options.CacheEnabled {
		log.Println("cache not enabled")
		return rgutil.ErrNil
	}

	if options.CachePreHandler != nil {
		options.CachePreHandler()
	}

	if err := rgcache.Init(); err != rgutil.ErrNil {
		slog.Error("server: cache failed to start")
		return rgutil.ErrNil
	}

	if options.CachePostHandler != nil {
		options.CachePostHandler()
	}

	return rgutil.ErrNil
}
