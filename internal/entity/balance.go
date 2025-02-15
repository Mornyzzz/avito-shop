package entity

type Balance struct {
	Username string `json:"username"`
	Coins    int    `json:"coins"`
}

// BuyRequest представляет запрос на покупку предмета
type BuyRequest struct {
	Item string `json:"item"`
}
