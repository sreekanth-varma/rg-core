package server

import (
	"bufio"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sreekanth-varma/rg-core/rgmiddleware"
	rgutil "github.com/sreekanth-varma/rg-core/rgutil"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client *mongo.Client
)

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

func initServer(options *Options) rgutil.Err {
	if !options.WebServerEnabled {
		slog.Info("server: server not enabled")
		return rgutil.ErrNil
	}

	module := os.Getenv("module")
	if module == "" {
		log.Fatal("module is not defined")
	}

	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.Use(rgmiddleware.CORSMiddleware())
	baseRouter := server.Group(module)

	if options.WebServerPreHandler != nil {
		options.WebServerPreHandler(server, baseRouter)
	}

	slog.Info("rgserver: starting on port " + options.WebServerPort)
	if err := server.Run(":" + options.WebServerPort); err != nil {
		slog.Error("rgserver: failed to start", "error", err.Error())
		return rgutil.ErrNil
	}

	return rgutil.ErrNil
}
