package server

import (
	"context"
	"fmt"
	"mini-kafka-go/pkg/config"
	"mini-kafka-go/pkg/log"
	"mini-kafka-go/pkg/protocol"
	"net"
	"sync"
	"time"
)

type LocalConnect struct {
	id          string
	conn        net.Conn
	lastTime    int64
	kafkaConfig config.KafkaConfig
}

func (c *LocalConnect) handleRequest(ctx context.Context, s *SocketServer) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	// TODO 请求大小暂时设为1024b

	for {
		var buffer [1024]byte
		n, err := c.conn.Read(buffer[:])
		if err != nil {
			s.removeClient(c)
			s.onError("Read bytes from clients failed", err)
			break
		}
		if n <= 0 {
			s.onError("Remote clients closed", err)
			break
		} else if n > 1024 {
			// TODO 返回请求过大的错误
		}
		c.lastTime = time.Now().Unix()
		fmt.Println(string(buffer[:n]))
		r, err := protocol.DeserializeKafkaRequest(buffer[:n])
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch r.GetApiKey() {
		case protocol.PRODUCE:
			fmt.Println("produce request")
			p, ok := r.(protocol.ProduceRequest)
			if !ok {
				fmt.Println("parse produce request failed")
			}
			fmt.Printf("%+v\n", p)
		case protocol.FETCH:

		default:
			fmt.Println("Unrecognized request")
		}
	}
}

func DefaultSocketServer(ctx context.Context, kafkaConfig *config.KafkaConfig, listen string) SocketServer {
	return SocketServer{
		ctx:                  ctx,
		config:               kafkaConfig,
		address:              listen,
		clients:              make(map[string]*LocalConnect),
		clientTimeoutSeconds: 10,
		maxClientNum:         1024,
		currentClientNum:     0,
		onError: func(msg string, err error) {
			fmt.Println(msg)
		},
		onStart: func(server *SocketServer) {
			log.Logger().Info("Socket server listening on " + server.address)
		},
		onConnect: func(client *LocalConnect) {
			fmt.Println("Client connected!", client.id)
		},
		onClientClose: func(client *LocalConnect) {
			fmt.Println("Client closed!", client.id)
		},
		onMessage: func(client *LocalConnect) {
			fmt.Println("Client message received!", client.id)
		},
		mutex: sync.Mutex{},
	}
}

type SocketServer struct {
	ctx                  context.Context
	config               *config.KafkaConfig
	address              string
	clients              map[string]*LocalConnect
	clientTimeoutSeconds int64
	maxClientNum         int
	currentClientNum     int
	onError              func(msg string, err error)
	onStart              func(server *SocketServer)
	onConnect            func(client *LocalConnect)
	onClientClose        func(client *LocalConnect)
	onMessage            func(client *LocalConnect)
	mutex                sync.Mutex
}

func (s *SocketServer) Startup() {
	listener, err := net.Listen("tcp4", s.address)
	if err != nil {
		s.onError("Startup socket server failed", err)
	}
	defer listener.Close()
	s.onStart(s)
	go s.checkClient(s.ctx)
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.onError("Socket server accept from clients failed", err)
			continue
		}
		if s.currentClientNum > s.maxClientNum {
			conn.Close()
			continue
		}
		client := s.makeClient(conn)
		s.onConnect(client)
		s.addClient(client)
		go client.handleRequest(s.ctx, s)
	}
}

func (s *SocketServer) makeClient(conn net.Conn) *LocalConnect {
	client := LocalConnect{
		id:       conn.RemoteAddr().String(),
		conn:     conn,
		lastTime: time.Now().Unix(),
	}
	return &client
}

func (s *SocketServer) checkClient(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if s.clientTimeoutSeconds < 1 {
			time.Sleep(5 * time.Second)
			continue
		}
		s.mutex.Lock()
		now := time.Now().Unix()
		for _, v := range s.clients {
			if now-v.lastTime > s.clientTimeoutSeconds {
				v.conn.Close()
			}
		}
		s.mutex.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func (s *SocketServer) removeClient(client *LocalConnect) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.clients[client.id]; ok {
		client.conn.Close()
		delete(s.clients, client.id)
		s.onClientClose(client)
		s.currentClientNum--
	}
}

func (s *SocketServer) addClient(client *LocalConnect) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.clients[client.id] = client
	s.currentClientNum++
}
