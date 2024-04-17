// this is responsible for connecting to and querying the RPC node directly
package bchain

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

func newEthClient(rpcNodeUrl string) *ethclient.Client {

	fmt.Println("RPC_NODE_URL", rpcNodeUrl)

	ec, err := ethclient.Dial(rpcNodeUrl)

	if err != nil {
		log.Fatal("Failed to connect to RPC_NODE_URL ", rpcNodeUrl)
	}

	return ec

}

func dialRpc(url string) (*rpc.Client, *ethclient.Client, error) {
	rpcClient, err := rpc.Dial(url)

	if err != nil {
		return nil, nil, err
	}

	ethereumClient := ethclient.NewClient(rpcClient)

	return rpcClient, ethereumClient, err
}

type BlockChain struct {
	client *ethclient.Client
}

func (b *BlockChain) GetBlockHeadersByHash(blockHash string) (*types.Header, error) {
	return b.client.HeaderByHash(context.Background(), common.HexToHash(blockHash))
}

func (b *BlockChain) GetBestBlockHeight() (uint64, error) {
	ctx := context.Background()

	return b.client.BlockNumber(ctx)
}

func (b *BlockChain) GetChainId() (*big.Int, error) {
	ctx := context.Background()

	return b.client.ChainID(ctx)
}

func (b *BlockChain) GetTransaction(txId string) (*types.Transaction, error) {

	tx, _, err := b.client.TransactionByHash(context.Background(), common.HexToHash(txId))

	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (b *BlockChain) CreateTransaction(rawHex string) (*types.Transaction, error) {
	rawTxBytes, err := hex.DecodeString(rawHex)

	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &tx)

	return tx, nil
}

func (b *BlockChain) SendTransaction(tx *types.Transaction) (string, error) {
	err := b.client.SendTransaction(context.Background(), tx)

	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func (b *BlockChain) GetTransactionReceipt(txId string) (*types.Receipt, error) {
	return b.client.TransactionReceipt(context.Background(), common.HexToHash(txId))
}

func (b *BlockChain) GetCurrentGasPrice() (*big.Int, error) {
	return b.client.SuggestGasPrice(context.Background())
}

func (b *BlockChain) GetAddressBalance(address string) (*big.Int, error) {
	return b.client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
}

func (b *BlockChain) GetNonce(address string) (uint64, error) {
	return b.client.PendingNonceAt(context.Background(), common.HexToAddress(address))
}

func NewBlockChain(rpcNodeUrl string) *BlockChain {

	client := newEthClient(rpcNodeUrl)

	return &BlockChain{client: client}
}

// func (b *BlockChain)
