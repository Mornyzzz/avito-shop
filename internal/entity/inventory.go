package entity

// Inventory представляет инвентарь пользователя
type Inventory struct {
	Username string `json:"username"`
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

// InventoryItem представляет элемент инвентаря
type InventoryItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
