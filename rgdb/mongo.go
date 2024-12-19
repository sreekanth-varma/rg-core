package rgdb

import (
	"context"
	"log/slog"
	"os"

	"github.com/sreekanth-varma/rg-core/rgutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	client *mongo.Client
)

func Connect(ctx *context.Context) rgutil.Err {
	// Load .env file
	mongoconn := options.Client().ApplyURI(os.Getenv("db_url"))
	var err error
	client, err = mongo.Connect(*ctx, mongoconn)
	if err != nil {
		slog.Error("db: connection failed", "error", err.Error())
		return rgutil.ErrUnavailable
	}
	err = client.Ping(*ctx, readpref.Primary())
	if err != nil {
		slog.Error("db: ping failed", "error", err.Error())
		return rgutil.ErrUnavailable
	}

	slog.Info("db: connected")
	return rgutil.ErrNil
}
func Disconnect(ctx *context.Context) {
	client.Disconnect(*ctx)
}
