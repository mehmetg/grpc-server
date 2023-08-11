// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	api "grpc-server/api"
	"grpc-server/utilities"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement
type server struct {
	api.UnimplementedStringFunctionServer
	mu sync.Mutex // protects routeNotes
}

// Reverse implements grpc_server.Reverse
func (s *server) Reverse(_ context.Context, in *api.RequestMessage) (*api.ResponseMessage, error) {
	log.Printf("Received: %v", in.GetMessage())
	return &api.ResponseMessage{Message: utilities.Reverse(in.GetMessage()), CharCount: int64(len(in.GetMessage()))}, nil
}
func (s *server) BidiEcho(stream api.StringFunction_BidiEchoServer) error {
	var count int64 = 0
	var sb strings.Builder
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		s.mu.Lock()
		sb.WriteString(in.GetMessage())
		count++
		s.mu.Unlock()
		if err := stream.Send(&api.ResponseMessage{Message: sb.String(), CharCount: count}); err != nil {
			return err
		}

	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	api.RegisterStringFunctionServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
