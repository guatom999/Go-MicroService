package server

import (
	"github.com/guatom999/Go-MicroService/modules/item/itemHandlers"
	"github.com/guatom999/Go-MicroService/modules/item/itemRepositories"
	"github.com/guatom999/Go-MicroService/modules/item/itemUseCases"
)

func (s *server) itemService() {
	itemRepository := itemRepositories.NewItemRepository(s.db)
	itemUseCase := itemUseCases.NewItemUseCase(itemRepository)
	itemHtppHandler := itemHandlers.NewItemHttpHandler(s.cfg, itemUseCase)
	itemGrpcHandler := itemHandlers.NewItemGrpcHandler(itemUseCase)

	_ = itemHtppHandler
	_ = itemGrpcHandler

	item := s.app.Group("/item_v1")

	// Health Check
	_ = item
}
