package model

// Part представляет деталь из InventoryService
type Part struct {
	UUID          string  `json:"uuid"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int64   `json:"stock_quantity"`
	Category      string  `json:"category"`
}