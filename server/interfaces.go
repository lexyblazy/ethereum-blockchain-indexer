package server

import "math/big"

type ChainIDResponse struct {
	ChainId *big.Int `json:"chainId"`
}

type AccountBalanceResponse struct {
	Balance string `json:"balance"`
}

type AccountNonceResponse struct {
	Nonce uint64 `json:"nonce"`
}

type APIErrorResponse struct {
	Message string `json:"message"`
}

type TransactionCreateRequestBody struct {
	RawHex string `json:"rawHex"`
}

type TransactionCreateResponse struct {
	TxId string `json:"txId"`
}

type GasPriceResponse struct{
	Price *big.Int `json:"gasPrice"`
}