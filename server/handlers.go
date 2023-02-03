package server

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	// "fmt"

	// "fmt"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"

	// "log"
	"net/http"
	// "github.com/ethereum/go-ethereum/log"
)

type ApiHandler struct {
	ethClient *ethclient.Client
}

func (a *APIErrorResponse) toJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (ap *ApiHandler) jsonWrapper(handler func(r *http.Request) (interface{}, error)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := handler(r)

		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			apiError, _ := (&APIErrorResponse{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}).toJSON()

			w.Write(apiError)

		} else {
			json.NewEncoder(w).Encode(data)

		}

	}
}

func (ap *ApiHandler) getChainId(r *http.Request) (interface{}, error) {

	ctx := context.Background()

	chainId, chainIdErr := ap.ethClient.ChainID(ctx)

	if chainIdErr != nil {

		return nil, errors.New("failed to get chainId")
	}

	result := &ChainIDResponse{
		ChainId: chainId.String(),
	}

	return result, nil
}

func (ap *ApiHandler) getAccountBalance(r *http.Request) (interface{}, error) {

	address := GetAccountFromRequestURL(r.URL.Path)

	if len(address) == 0 {
		return nil, errors.New("address is required")
	}

	ctx := context.Background()

	balance, _ := ap.ethClient.BalanceAt(ctx, common.HexToAddress(address), nil)

	return &AccountBalanceResponse{
		Balance: balance.String(),
	}, nil

}

func (ap *ApiHandler) getNonce(r *http.Request) (interface{}, error) {
	address := GetAccountFromRequestURL(r.URL.Path)

	if len(address) == 0 {
		return nil, errors.New("address is required")
	}

	nonce, _ := ap.ethClient.PendingNonceAt(context.Background(), common.HexToAddress(address))

	return &AccountNonceResponse{
		Nonce: nonce,
	}, nil

}

func (ap *ApiHandler) getBlock(r *http.Request) (interface{}, error) {

	i := strings.LastIndexByte(r.URL.Path, '/')
	blockHash := (r.URL.Path[i+1:])

	if len(blockHash) == 0 {
		return nil, errors.New("blockHash is required")
	}

	block, err := ap.ethClient.HeaderByHash(context.Background(), common.HexToHash(blockHash))

	// block, err := ap.ethClient.HeaderByNumber(context.Background(), new(big.Int).SetUint64(0))

	return block, err
}

func (ap *ApiHandler) createTransaction(r *http.Request) (interface{}, error) {

	headerContentType := r.Header.Get("Content-Type")

	if headerContentType != "application/json" {

		return nil, errors.New("content Type is not application/json")
	}

	var txCreateBody TransactionCreateRequestBody

	err := json.NewDecoder(r.Body).Decode(&txCreateBody)

	if err != nil {
		return nil, err
	}

	fmt.Println("txCreateBodyRawHex", txCreateBody.RawHex)

	rawTxBytes, decodeStringError := hex.DecodeString(txCreateBody.RawHex)

	if decodeStringError != nil {
		return nil, decodeStringError
	}

	tx := new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &tx)

	sendTxErr := ap.ethClient.SendTransaction(context.Background(), tx)

	return &TransactionCreateResponse{
		TxId: tx.Hash().Hex(),
	}, sendTxErr

}
