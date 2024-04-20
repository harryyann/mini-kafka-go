package server

import (
	"context"
	"go.uber.org/zap"
	"strings"

	"mini-kafka-go/pkg/config"
	"mini-kafka-go/pkg/log"
)

type KafkaServer struct {
	config    config.KafkaConfig
	listeners map[*SocketServer]context.CancelFunc
}

func NewKafkaServer(ctx context.Context, c *config.KafkaConfig) KafkaServer {
	listeners := strings.Split(c.Listeners, ",")
	if len(listeners) == 0 {
		log.Logger().Fatal("listener must be configured at least once")
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

func (s *KafkaServer) GracefullyStop() {
	for _, cancel := range s.listeners {
		cancel()
	}

	err := s.cleanup()
	if err != nil {
		log.Logger().Error("Clean up error", zap.Error(err))
	}
}

func (s *KafkaServer) cleanup() error {

	return nil
}
