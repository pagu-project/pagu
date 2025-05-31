package grpc

import (
	"context"
	"net"

	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/pkg/log"
	pagu "github.com/pagu-project/pagu/pkg/proto/gen/go"
	"google.golang.org/grpc"
)

type Server struct {
	ctx      context.Context
	cancel   context.CancelFunc
	listener net.Listener
	address  string
	engine   *engine.BotEngine
	grpc     *grpc.Server
	cfg      *Config
}

func NewServer(ctx context.Context, eng *engine.BotEngine, cfg *Config) *Server {
	return &Server{
		ctx:    ctx,
		engine: eng,
		cfg:    cfg,
	}
}

func (s *Server) Start() error {
	log.Info("Starting gRPC Server")
	listener, err := net.Listen("tcp", s.cfg.Listen)
	if err != nil {
		return err
	}

	s.startListening(listener)

	return nil
}

func (s *Server) startListening(listener net.Listener) {
	opts := make([]grpc.UnaryServerInterceptor, 0)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(opts...))

	server := newPaguServer(s)

	pagu.RegisterPaguServer(grpcServer, server)

	s.listener = listener
	s.address = listener.Addr().String()
	s.grpc = grpcServer

	log.Info("gRPC Server Started Listening", "address", listener.Addr().String())
	go func() {
		if err := s.grpc.Serve(listener); err != nil {
			log.Error("error on grpc serve", "error", err)
		}
	}()
}

func (s *Server) Stop() error {
	log.Info("Stopping gRPC Server", "addr", s.address)

	s.cancel()

	if s.grpc != nil {
		s.grpc.Stop()
	}

	return s.listener.Close()
}
