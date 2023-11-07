package inventory

import (
	"github.com/guatom999/Go-MicroService/modules/item"
	"github.com/guatom999/Go-MicroService/modules/models"
)

type (
	UpdateInventory struct {
		PlayerId string `json:"player" validate:"required,max=64"`
		ItemId   string `json:"item_id" validate:"required,max=64"`
	}

	ItemInInventory struct {
		InventoryId string `json:"inventory_id"`
		*item.ItemShowCase
	}

	PlayerInventory struct {
		PlayerId            string `json:"player_id"`
		*models.PaginateRes `json:"data"`
	}
)
