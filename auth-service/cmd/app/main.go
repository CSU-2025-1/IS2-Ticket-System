package main

import (
	"auth-service/config"
	"auth-service/internal/app"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	defaultConfigPath = "./config/config.yaml"
)

func main() {
	var path string

	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	if path == "" {
		path = defaultConfigPath
	}

	cfg, err := config.Parse(path)
	if err != nil {
		log.Fatal("parse config error: ", err.Error())
		return
	}

	ctx := context.Background()

	application, err := app.Build(ctx, cfg)
	if err != nil {
		log.Fatal("build application error: ", err.Error())
		return
	}

	go func(application *app.App) {
		if err = application.Run(ctx); err != nil {
			log.Fatal("application run error: ", err.Error())
			return
		}
	}(application)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	application.Shutdown(ctx)

}
