package inventoryHandlers

import (
	"context"

	inventoryPb "github.com/guatom999/Go-MicroService/modules/inventory/inventoryPb"
	"github.com/guatom999/Go-MicroService/modules/inventory/inventoryUseCases"
)

type (
	inventoryGrpcHandler struct {
		inventoryPb.UnimplementedInventoryGrpcServiceServer
		inventoryUseCase inventoryUseCases.IInventoryUseCaseService
	}
)

func NewInventoryGrpcHandler(inventoryUseCase inventoryUseCases.IInventoryUseCaseService) *inventoryGrpcHandler {
	return &inventoryGrpcHandler{inventoryUseCase: inventoryUseCase}
}

func (g *inventoryGrpcHandler) IsAvaliableToSell(ctx context.Context, req *inventoryPb.IsAvaliableToSellReq) (*inventoryPb.IsAvaliableToSellRes, error) {
	return nil, nil
}
