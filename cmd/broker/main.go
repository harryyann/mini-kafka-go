package main

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v3"

	"mini-kafka-go/internal/server"
	"mini-kafka-go/pkg/config"
	"mini-kafka-go/pkg/log"
)

const DefaultConfigFile = "config/server.yaml"

var configPath string

func init() {
	var mode string
	flag.StringVar(&mode, "mode", log.DEVELOPMENT, "Running mode, develop or produce")
	flag.StringVar(&configPath, "config", DefaultConfigFile, "Configuration file path")
	flag.Parse()
	err := log.InitLog(mode)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	if configPath == "" {
		configPath = DefaultConfigFile
	}
	fileBytes, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Load config file server.yaml failed:")
		os.Exit(1)
	}
	conf := config.NewKafkaConfig()
	err = yaml.Unmarshal(fileBytes, conf)
	if err != nil {
		fmt.Println("Parse config failed")
		os.Exit(1)
	}

	ctx := context.Background()
	s := server.NewKafkaServer(ctx, conf)
	s.Startup()
	signalCtx, signalCancel := context.WithCancel(ctx)
	signalChan := make(chan struct{})
	go signalListen(signalCtx, signalChan)
	select {
	case <-ctx.Done():
		signalCancel()
		return
	case <-signalChan:
		log.Logger().Info("Interrupt signal! Server will stop gracefully")
		s.GracefullyStop()
		defer func(l *zap.SugaredLogger) {
			err := l.Sync()
			if err != nil {
				log.Logger().Error("Sync logger failed", zap.Error(err))
			}
		}(log.Logger())
		os.Exit(0)
	}
}

func signalListen(signalCtx context.Context, signalChan chan struct{}) {
	log.Logger().Info("Interrupt signal listening...")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	select {
	case s := <-c:
		log.Logger().Info("Received stop signal", s.String())
		signalChan <- struct{}{}
		return
	case <-signalCtx.Done():
		return
	}
}
