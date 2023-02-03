package server

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"

	"net/http"
)

type ApiHandler struct {
	ethClient *ethclient.Client
}

func (a *APIErrorResponse) toJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (ap *ApiHandler) jsonWrapper(handler func(r *http.Request) (interface{}, int, error)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		if !IsJsonContentType(r) {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("content Type is not application/json"))

			return

		}

		data, httpStatus, err := handler(r)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader((httpStatus))

		if err != nil {
			fmt.Println(err)
			apiError, _ := (&APIErrorResponse{
				Message: err.Error(),
			}).toJSON()

			w.Write(apiError)

		} else {
			json.NewEncoder(w).Encode(data)

		}

	}
}

func (ap *ApiHandler) getChainId(r *http.Request) (interface{}, int, error) {

	ctx := context.Background()

	chainId, chainIdErr := ap.ethClient.ChainID(ctx)

	if chainIdErr != nil {

		return nil, http.StatusInternalServerError, errors.New("failed to get chainId")
	}

	result := &ChainIDResponse{
		ChainId: chainId,
	}

	return result, http.StatusOK, nil
}

func (ap *ApiHandler) getAccountBalance(r *http.Request) (interface{}, int, error) {

	address := GetParamFromRequestURL(r.URL.Path)

	if len(address) == 0 {
		return nil, http.StatusBadRequest, errors.New("address is required")
	}

	ctx := context.Background()

	balance, _ := ap.ethClient.BalanceAt(ctx, common.HexToAddress(address), nil)

	return &AccountBalanceResponse{
		Balance: balance.String(),
	}, http.StatusOK, nil

}

func (ap *ApiHandler) getNonce(r *http.Request) (interface{}, int, error) {
	address := GetParamFromRequestURL(r.URL.Path)

	if len(address) == 0 {
		return nil, http.StatusBadRequest, errors.New("address is required")
	}

	nonce, _ := ap.ethClient.PendingNonceAt(context.Background(), common.HexToAddress(address))

	return &AccountNonceResponse{
		Nonce: nonce,
	}, http.StatusOK, nil

}

func (ap *ApiHandler) getBlock(r *http.Request) (interface{}, int, error) {

	blockHash := GetParamFromRequestURL(r.URL.Path)

	if len(blockHash) == 0 {
		return nil, http.StatusBadRequest, errors.New("blockHash is required")
	}

	block, err := ap.ethClient.HeaderByHash(context.Background(), common.HexToHash(blockHash))

	if err != nil {
		return nil, http.StatusNotFound, err

	}

	return block, http.StatusOK, err
}

func (ap *ApiHandler) createTransaction(r *http.Request) (interface{}, int, error) {

	var txCreateBody TransactionCreateRequestBody

	err := json.NewDecoder(r.Body).Decode(&txCreateBody)

	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	fmt.Println("txCreateBodyRawHex", txCreateBody.RawHex)

	rawTxBytes, decodeStringError := hex.DecodeString(txCreateBody.RawHex)

	if decodeStringError != nil {
		return nil, http.StatusInternalServerError, decodeStringError
	}

	tx := new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &tx)

	sendTxErr := ap.ethClient.SendTransaction(context.Background(), tx)

	if sendTxErr != nil {
		return nil, http.StatusInternalServerError, sendTxErr
	}

	return &TransactionCreateResponse{
		TxId: tx.Hash().Hex(),
	}, http.StatusOK, nil

}

func (ap *ApiHandler) getTransaction(r *http.Request) (interface{}, int, error) {
	txId := GetParamFromRequestURL(r.URL.Path)

	fmt.Println("TransactionId", txId)

	tx, _, err := ap.ethClient.TransactionByHash(context.Background(), common.HexToHash(txId))

	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return tx, http.StatusOK, err
}

func (ap *ApiHandler) getTransactionReceipt(r *http.Request) (interface{}, int, error) {
	txId := GetParamFromRequestURL(r.URL.Path)

	fmt.Println("TransactionId", txId)

	txReceipt, err := ap.ethClient.TransactionReceipt(context.Background(), common.HexToHash(txId))

	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return txReceipt, http.StatusOK, err
}
