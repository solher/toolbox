package api

import (
	"context"
	"net"
	"time"

	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// NewDefaultGRPCServer returns a new default gRPC server.
func NewDefaultGRPCServer(addr string) *GRPCServer {
	return &GRPCServer{
		Addr:       addr,
		grpcServer: grpc.NewServer(),
	}
}

// GRPCServer is a simplae wrapper around a gRPC server.
type GRPCServer struct {
	Addr       string
	grpcServer *grpc.Server
}

// ListenAndServe allows a similar API as with HTTP.
func (s *GRPCServer) ListenAndServe() error {
	grpcConn, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	return s.grpcServer.Serve(grpcConn)
}

// Shutdown allows a homogeous API with the rest of the toolbox.
func (s *GRPCServer) Shutdown(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	return nil
}

// GetServer returns the concrete implementation of the gRPC server.
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.grpcServer
}

// NewDefaultGRPCClientConn is a simplae wrapper around a gRPC client connection.
func NewDefaultGRPCClientConn(ctx context.Context, addr string, connTimeout time.Duration, l log.Logger) (*grpc.ClientConn, error) {
	l = log.With(l, "addr", addr)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),                   // Allows no SSL.
		grpc.WithBackoffMaxDelay(connTimeout), // Enables backoff.
	}

	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, err
	}

	go func() {
		timeout := time.After(connTimeout)
		for {
			if s := conn.GetState(); s == connectivity.Ready {
				l.Log("msg", "connection ready (gRPC)")
				return
			}
			select {
			case <-timeout:
				l.Log("msg", "connection expired (gRPC)", "timeout", connTimeout.String())
				return
			default:
			}
			time.Sleep(time.Second)
		}
	}()

	return conn, nil
}
