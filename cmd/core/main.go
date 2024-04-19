package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v3"

	"mini-kafka-go/internal/server"
	"mini-kafka-go/pkg/config"
	"mini-kafka-go/pkg/log"
)

const DefaultConfigFile = "config/server.yaml"

func main() {
	// TODO 加上cobra命令行参数解析

	fileBytes, err := os.ReadFile(DefaultConfigFile)
	if err != nil {
		fmt.Println("Load config file server.yaml failed:")
		os.Exit(1)
	}
	c := config.NewKafkaConfig()
	err = yaml.Unmarshal(fileBytes, c)
	if err != nil {
		fmt.Println("Parse config failed")
		os.Exit(1)
	}
	err = log.InitLog(c.Mode)
	if err != nil {
		fmt.Println("Initialize logger failed")
		os.Exit(1)
	}
	ctx := context.Background()
	s := server.NewKafkaServer(ctx, c)
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
		// TODO serer的优雅退出
		s.GracefulStop()
		os.Exit(0)
	}
}

func signalListen(signalCtx context.Context, signalChan chan struct{}) {
	log.Logger().Info("Interrupt signal listening")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	select {
	case <-c:
		log.Logger().Info("Received stop signal")
		signalChan <- struct{}{}
		return
	case <-signalCtx.Done():
		return
	}
}
