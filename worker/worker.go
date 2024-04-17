package worker

import (
	"fmt"
	"indexer/bchain"
	"indexer/db"
	"log"
)

type Worker struct {
	bc *bchain.BlockChain
	db *db.RocksDB
}

func NewWorker(bc *bchain.BlockChain, index *db.RocksDB, nodeUrl string) *Worker {
	return &Worker{
		bc: bc,
		db: index,
	}
}

func (w *Worker) Sync(fromheight, toHeight int) {

	latestBlock, err := w.bc.GetBestBlockHeight()
	if err != nil {
		log.Fatal("Failed to get latest block")
	}

	fmt.Println("latestBlock", latestBlock, "fromHeight", fromheight, "toHeight", toHeight)
}
