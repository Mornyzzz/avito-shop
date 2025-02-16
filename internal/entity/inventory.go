package entity

type Inventory struct {
	Username string `json:"username"`
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

type InventoryItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
