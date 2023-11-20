package server

import (
	"log"

	"github.com/guatom999/Go-MicroService/modules/item/itemHandlers"
	itemPb "github.com/guatom999/Go-MicroService/modules/item/itemPb"
	"github.com/guatom999/Go-MicroService/modules/item/itemRepositories"
	"github.com/guatom999/Go-MicroService/modules/item/itemUseCases"
	"github.com/guatom999/Go-MicroService/pkg/grpccon"
)

func (s *server) itemService() {
	itemRepository := itemRepositories.NewItemRepository(s.db)
	itemUseCase := itemUseCases.NewItemUseCase(itemRepository)
	itemHtppHandler := itemHandlers.NewItemHttpHandler(s.cfg, itemUseCase)
	itemGrpcHandler := itemHandlers.NewItemGrpcHandler(itemUseCase)

	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.ItemUrl)

		itemPb.RegisterItemGrpcServiceServer(grpcServer, itemGrpcHandler)

		log.Printf("Item Grpc server listening on: %s", s.cfg.Grpc.ItemUrl)
		grpcServer.Serve(lis)

	}()

	_ = itemGrpcHandler

	item := s.app.Group("/item_v1")

	// Health Check
	item.GET("", s.healthCheckService)

	item.POST("/item", s.middleware.JwtAuthorization(s.middleware.RbacAuthorization(itemHtppHandler.CreateItem, []int{1, 0})))

	item.GET("/item/:item_id", itemHtppHandler.FindOneItem)

	item.GET("/item", itemHtppHandler.FindManyItem)
}
