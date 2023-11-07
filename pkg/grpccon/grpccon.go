package grcpcon

import (
	"errors"
	"log"
	"net"

	"github.com/guatom999/Go-MicroService/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	GrcpClientFactoryHandler interface {
	}

	grcpClientFactory struct {
	}

	grpcAuth struct {
	}
)

func (g *grcpClientFactory) Auth() {

}

func NewGrpcClient(host string) (GrcpClientFactoryHandler, error) {

	options := make([]grpc.DialOption, 0)

	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	clientConn, err := grpc.Dial(host, options...)

	if err != nil {
		log.Printf("Error: Grpc client connection failed:%v", err)
		return nil, errors.New("error: grpc client connection failed:%")
	}

	return clientConn, nil
}

func NewGrpcServer(cfg *config.Config, host string) (*grpc.Server, net.Listener) {
	options := make([]grpc.ServerOption, 0)

	grpcServer := grpc.NewServer(options...)

	lis, err := net.Listen("tcp", host)

	if err != nil {
		log.Printf("Error: Failed to listen: %v", err)
	}

	return grpcServer, lis

}
