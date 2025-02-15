package entity

// InfoResponse представляет информацию о монетах, инвентаре и истории транзакций
type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}
