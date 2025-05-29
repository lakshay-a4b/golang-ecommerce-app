package models

type CartProduct struct {
	ProductID int  	  `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
