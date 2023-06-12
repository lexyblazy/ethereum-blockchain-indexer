package worker

import (
	"context"
	"fmt"
	"indexer/db"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Worker struct {
	ec *ethclient.Client
	db *db.RocksDB
}

func NewWorker(ec *ethclient.Client, index *db.RocksDB, nodeUrl string) *Worker {
	return &Worker{
		ec: ec,
		db: index,
	}
}

func (w *Worker) Sync() {

	ctx := context.Background()
	latestBlock, err := w.ec.BlockNumber(ctx)
	if err != nil {
		log.Fatal("Failed to get latest block")
	}

	fmt.Println("latestBlock", latestBlock)
}
