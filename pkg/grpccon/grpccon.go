package grpccon

import (
	"context"
	"errors"
	"log"
	"net"

	authPb "github.com/guatom999/Go-MicroService/modules/auth/authPb"
	inventoryPb "github.com/guatom999/Go-MicroService/modules/inventory/inventoryPb"
	playerPb "github.com/guatom999/Go-MicroService/modules/player/playerPb"
	"github.com/guatom999/Go-MicroService/pkg/jwtauth"

	itemPb "github.com/guatom999/Go-MicroService/modules/item/itemPb"

	"github.com/guatom999/Go-MicroService/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type (
	GrcpClientFactoryHandler interface {
		Auth() authPb.AuthGrpcServiceClient
		Inventory() inventoryPb.InventoryGrpcServiceClient
		Player() playerPb.PlayerGrpcServiceClient
		Item() itemPb.ItemGrpcServiceClient
	}

	grcpClientFactory struct {
		client *grpc.ClientConn
	}

	grpcAuth struct {
		secretKey string
	}
)

func (g *grpcAuth) unaryAuthorization(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("Error: Metadata not found")
		return nil, errors.New("error:metadata not found")
	}

	authHeader, ok := md["auth"]
	if !ok {
		log.Printf("Error: Metadata not found")
		return nil, errors.New("error:metadata not found")
	}

	if len(authHeader) == 0 {
		log.Printf("Error: Metadata not found")
		return nil, errors.New("error:metadata not found")
	}

	clamis, err := jwtauth.ParseToken(g.secretKey, string(authHeader[0]))
	if err != nil {
		log.Printf("Error: ParseToken failed")
		return nil, errors.New("error: token is invalid")
	}

	log.Printf("claims: %v", clamis)

	return handler(ctx, req)
}

func (g *grcpClientFactory) Auth() authPb.AuthGrpcServiceClient {
	return authPb.NewAuthGrpcServiceClient(g.client)
}

func (g *grcpClientFactory) Inventory() inventoryPb.InventoryGrpcServiceClient {
	return inventoryPb.NewInventoryGrpcServiceClient(g.client)
}

func (g *grcpClientFactory) Player() playerPb.PlayerGrpcServiceClient {
	return playerPb.NewPlayerGrpcServiceClient(g.client)
}

func (g *grcpClientFactory) Item() itemPb.ItemGrpcServiceClient {
	return itemPb.NewItemGrpcServiceClient(g.client)
}

func NewGrpcClient(host string) (GrcpClientFactoryHandler, error) {

	options := make([]grpc.DialOption, 0)

	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	clientConn, err := grpc.Dial(host, options...)

	if err != nil {
		log.Printf("Error: Grpc client connection failed:%v", err)
		return nil, errors.New("error: grpc client connection failed:%")
	}

	return &grcpClientFactory{client: clientConn}, nil
}

func NewGrpcServer(cfg *config.Jwt, host string) (*grpc.Server, net.Listener) {
	options := make([]grpc.ServerOption, 0)

	grpcAuth := &grpcAuth{
		secretKey: cfg.ApiSecretKey,
	}

	options = append(options, grpc.UnaryInterceptor(grpcAuth.unaryAuthorization))

	grpcServer := grpc.NewServer(options...)

	lis, err := net.Listen("tcp", host)

	if err != nil {
		log.Printf("Error: Failed to listen: %v", err)
	}

	return grpcServer, lis

}
