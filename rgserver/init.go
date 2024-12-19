package server

import (
	"bufio"
	"context"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
	rgutil "github.com/sreekanth-varma/rg-core/rgutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	client *mongo.Client
	rdb    *redis.Client
)

var ()

func Connect(ctx1 *context.Context) {
	// Load .env file
	file, err := os.Open("app.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	defer file.Close()
	ctx = *ctx1
	log.Println("url : ", os.Getenv("db_url"))
	mongoconn := options.Client().ApplyURI(os.Getenv("db_url"))
	//	var err error
	client, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal("mongo: connection failed", err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("mongo: ping failed", err)
	}

	log.Println("mongo: connected")
}
func Disconnect(ctx *context.Context) {
	client.Disconnect(*ctx)
}

func DB() *mongo.Database {
	return client.Database(os.Getenv("DB_NAME"))
}

func LoadConfig() string {
	file, err := os.Open("app.env")
	if err != nil {
		log.Println("app.env not found. Config not loaded from file")
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// read each line
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), "")

		// ignore comments and empty lines
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}

		// separate key and value
		key, value, valid := strings.Cut(line, "=")
		// ignore invalid lines
		if !valid {
			log.Fatalf("config: invalid line:%v\n", line)
			continue
		}

		// prepare key, value fields
		key1 := strings.ToLower(strings.Trim(key, " "))
		value1 := strings.Trim(value, " ")

		// ignore if already exists in env
		_, found := os.LookupEnv(key)
		if found {
			continue
		}

		os.Setenv(key1, value1)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("config: failed to read, %v\n", err)
	}

	return ""
}

func GetEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func initCache(options *Options) rgutil.Err {
	if !options.CacheEnabled {
		log.Println("cache not enabled")
		return rgutil.ErrNil
	}

	if options.CachePreHandler != nil {
		options.CachePreHandler()
	}

	if err := Init(); err != rgutil.ErrNil {
		slog.Error("server: cache failed to start")
		return rgutil.ErrNil
	}

	if options.CachePostHandler != nil {
		options.CachePostHandler()
	}

	return rgutil.ErrNil
}

func Init() rgutil.Err {
	ctx = context.Background()

	// Get Redis host and password from environment variables
	host := os.Getenv("redis_url")
	if host == "" {
		return rgutil.NewErr("redis_url cannot be empty")
	}

	password := os.Getenv("redis_password")

	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Println("cache: connected")
			return nil
		},
	})

	// Test the connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		return rgutil.NewErr("failed to connect to Redis: " + err.Error())
	}

	return rgutil.ErrNil
}
