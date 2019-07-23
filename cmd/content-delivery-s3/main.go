package main

import (
	"context"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/janminder/content-delivery-s3-backend/api/rest"
	"github.com/janminder/content-delivery-s3-backend/api/services"
	"github.com/janminder/content-delivery-s3-backend/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var conf *viper.Viper

func main() {

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Setup Logging
	setupLogger()

	// Check Profile Selection
	profile := os.Getenv("PROFILE")

	if profile != "" {
		switch profile {
			case "cloud":
				profile = "cloud"
			default:
				profile = "dev"
		}
	} else {
		profile = "dev"
	}

	log.Info("Active Profile: ", profile)
	conf = config.LoadConfig(profile)

	port := os.Getenv("PORT")

	if port != "" {
		conf.Set("server.port", port)
	}

	// Create Service
	fileService := services.NewFileService(conf)

	// Create HTTP Server
	httpServer := rest.NewServer(fileService, conf)
	echoServer := httpServer.InitializeHandler()

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := echoServer.Start(":" + strconv.Itoa(conf.GetInt("server.port"))); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	echoServer.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}

func setupLogger() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
}
