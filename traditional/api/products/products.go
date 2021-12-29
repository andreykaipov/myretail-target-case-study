package products

// ProductResponse represents our _myRetail_ response bodies.
type ProductResponse struct {
	ID           int           `json:"id" bson:"_id"`
	Name         string        `json:"name"`
	CurrentPrice *currentPrice `json:"current_price"`
}

// ProductRequest represents Put request bodies made to our API.
type ProductRequest struct {
	ID           int           `json:"id" bson:"_id"`
	CurrentPrice *currentPrice `json:"current_price" validate:"required"`
}

type currentPrice struct {
	Value        float64 `json:"value" validate:"required"`
	CurrencyCode string  `json:"currency_code" validate:"required"`
}
