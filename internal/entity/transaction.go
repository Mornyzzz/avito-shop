package entity

type CoinTransaction struct {
	ID       int    `json:"id"`
	FromUser string `json:"fromUser"`
	ToUser   string `json:"toUser"`
	Amount   int    `json:"amount"`
}

type ReceivedTransaction struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentTransaction struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Received []ReceivedTransaction `json:"received"`
	Sent     []SentTransaction     `json:"sent"`
}
