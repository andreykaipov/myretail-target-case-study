package redsky

// Product represents a response from the RedSky product APIs.
type Product struct {
	Data struct {
		Product struct {
			Tcin string `json:"tcin"`
			Item struct {
				ProductDescription struct {
					Title                 string `json:"title"`
					DownstreamDescription string `json:"downstream_description"`
				} `json:"product_description"`
			} `json:"item"`
		} `json:"product"`
		Price struct {
			CurrentRetail float64 `json:"current_retail"`
		} `json:"price"`
	} `json:"data"`
}
