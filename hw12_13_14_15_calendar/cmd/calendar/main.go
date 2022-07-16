package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/internal/app"
	"github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/internal/config"
	"github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/internal/server/http"
	"github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := config.New(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	logg, err := logger.New(cfg.GetString("logger.path"))
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := storage.Connect(cfg)
	if err != nil {
		logg.Error(err.Error())
	}

	calendar := app.New(logg, conn)

	server := internalhttp.NewServer(cfg, logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
