package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"indexer/bchain"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
)

const maxUint32 = ^uint32(0)

type ApiHandler struct {
	bc *bchain.BlockChain
}

func (a *APIErrorResponse) toJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (ap *ApiHandler) jsonWrapper(handler func(r *http.Request) (interface{}, int, error)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		if IsPostRequest(r) && !IsJsonContentType(r) {
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

	chainId, chainIdErr := ap.bc.GetChainId()

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

	balance, _ := ap.bc.GetAddressBalance(address)

	return &AccountBalanceResponse{
		Balance: balance.String(),
	}, http.StatusOK, nil

}

func (ap *ApiHandler) getNonce(r *http.Request) (interface{}, int, error) {
	address := GetParamFromRequestURL(r.URL.Path)

	if len(address) == 0 {
		return nil, http.StatusBadRequest, errors.New("address is required")
	}

	nonce, _ := ap.bc.GetNonce(address)

	return &AccountNonceResponse{
		Nonce: nonce,
	}, http.StatusOK, nil

}

func (ap *ApiHandler) getBlock(r *http.Request) (interface{}, int, error) {

	blockIdOrHash := GetParamFromRequestURL(r.URL.Path)

	var blockHash string
	var blockHeight uint32

	var block *types.Block
	var getBlockError error

	if len(blockIdOrHash) == 0 {
		return nil, http.StatusBadRequest, errors.New("blockHash is required")
	}

	value, err := strconv.Atoi(blockIdOrHash)

	if err != nil {
		blockHash = blockIdOrHash
		block, getBlockError = ap.bc.GetBlockByHash(blockHash)
	}

	if err == nil && value < int(maxUint32) {
		blockHeight = uint32(value)
		block, getBlockError = ap.bc.GetBlockByHeight(blockHeight)
	}

	if getBlockError != nil {
		return nil, http.StatusNotFound, err

	}

	return block.Header(), http.StatusOK, nil

}

func (ap *ApiHandler) createTransaction(r *http.Request) (interface{}, int, error) {

	var txCreateBody TransactionCreateRequestBody

	err := json.NewDecoder(r.Body).Decode(&txCreateBody)

	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	tx, err := ap.bc.CreateTransaction(txCreateBody.RawHex)

	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	txHash, err := ap.bc.SendTransaction(tx)

	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &TransactionCreateResponse{
		TxId: txHash,
	}, http.StatusOK, nil

}

func (ap *ApiHandler) getTransaction(r *http.Request) (interface{}, int, error) {
	txId := GetParamFromRequestURL(r.URL.Path)

	tx, err := ap.bc.GetTransaction((txId))

	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return tx, http.StatusOK, nil
}

func (ap *ApiHandler) getTransactionReceipt(r *http.Request) (interface{}, int, error) {
	txId := GetParamFromRequestURL(r.URL.Path)

	txReceipt, err := ap.bc.GetTransactionReceipt(txId)

	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return txReceipt, http.StatusOK, nil
}

func (ap *ApiHandler) getCurrentGasPrice(r *http.Request) (interface{}, int, error) {

	gasPrice, err := ap.bc.GetCurrentGasPrice()

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &GasPriceResponse{Price: gasPrice}, http.StatusOK, nil
}
