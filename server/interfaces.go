package server

type ChainIDResponse struct {
	ChainId string `json:"chainId"`
}

type AccountBalanceResponse struct {
	Balance string `json:"balance"`
}


type APIErrorResponse struct {
	Message string `json:"message"`
	Status  uint16 `json:"status"`
}