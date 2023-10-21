package server

import (
	"context"
	"fmt"
	"mini-kafka-go/pkg/config"
	"os"
	"strings"
)

type KafkaServer struct {
	ctx       context.Context
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
		ctx:       ctx,
		listeners: localListeners,
	}
}

func (s *KafkaServer) Startup() {
	for ss := range s.listeners {
		go ss.Startup()
	}
}

func (s *KafkaServer) Stop() {
	for _, cancel := range s.listeners {
		cancel()
	}
}
