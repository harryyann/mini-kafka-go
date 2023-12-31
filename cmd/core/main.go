package main

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"mini-kafka-go/internal/server"
	"mini-kafka-go/pkg/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fileBytes, err := os.ReadFile("config/server.yaml")
	if err != nil {
		fmt.Println("read file failed")
		os.Exit(1)
	}
	c := config.NewKafkaConfig()
	err = yaml.Unmarshal(fileBytes, c)
	if err != nil {
		fmt.Println("parse config failed")
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
		fmt.Println("Interrupt by signal")
		os.Exit(1)
	}
}

func signalListen(signalCtx context.Context, signalChan chan struct{}) {
	fmt.Println("Signal listen started")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	for {
		select {
		case <-c:
			fmt.Println("Received stop signal")
			signalChan <- struct{}{}
			return
		case <-signalCtx.Done():
			return
		default:
		}
	}
}
