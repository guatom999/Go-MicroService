package inventory

import (
	"github.com/guatom999/Go-MicroService/modules/item"
)

type (
	UpdateInventory struct {
		PlayerId string `json:"player" validate:"required,max=64"`
		ItemId   string `json:"item_id" validate:"required,max=64"`
	}

	ItemInInventory struct {
		InventoryId string `json:"inventory_id"`
		PlayerId    string `json:"player_id"`
		*item.ItemShowCase
	}
)
