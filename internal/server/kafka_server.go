package server

import (
	"context"
	"fmt"
	"mini-kafka-go/pkg/config"
	"os"
	"strings"
)

type KafkaServer struct {
	config    config.KafkaConfig
	listeners map[*SocketServer]context.CancelFunc
}

func NewKafkaServer(ctx context.Context, c *config.KafkaConfig) KafkaServer {
	listeners := strings.Split(c.Listeners, ",")
	if len(listeners) == 0 {
		fmt.Println("listener must be configured at least once")
		os.Exit(1)
	}
	localListeners := map[*SocketServer]context.CancelFunc{}
	for _, l := range listeners {
		serverCtx, cancel := context.WithCancel(ctx)
		socketServer := DefaultSocketServer(serverCtx, c, l)
		localListeners[&socketServer] = cancel
	}
	return KafkaServer{
		listeners: localListeners,
	}
}

func (s *KafkaServer) Startup() {
	for ss := range s.listeners {
		go ss.Startup()
	}
}

func (s *KafkaServer) GracefulStop() {
	for _, cancel := range s.listeners {
		cancel()
	}
	err := s.cleanup()
	if err != nil {
		fmt.Println("Clean up error")
	}
}

func (s *KafkaServer) cleanup() error {

	return nil
}
