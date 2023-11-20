package itemHandlers

import (
	"context"

	itemPb "github.com/guatom999/Go-MicroService/modules/item/itemPb"
	"github.com/guatom999/Go-MicroService/modules/item/itemUseCases"
)

type (
	itemGrpcHandler struct {
		itemPb.UnimplementedItemGrpcServiceServer
		itemUseCase itemUseCases.IItemUseCaseService
	}
)

func NewItemGrpcHandler(itemUseCase itemUseCases.IItemUseCaseService) *itemGrpcHandler {
	return &itemGrpcHandler{itemUseCase: itemUseCase}
}

func (g *itemGrpcHandler) FindItemsInIds(ctx context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	return g.itemUseCase.FindItemsInIds(ctx, req)
}
