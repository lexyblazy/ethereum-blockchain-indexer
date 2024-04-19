package worker

import (
	"fmt"
	"indexer/bchain"
	"indexer/db"
	"log"
	// "sync"
	// "time"
	// "log"
)

type Worker struct {
	bc           *bchain.BlockChain
	db           *db.RocksDB
	workersCount int
}

func NewWorker(bc *bchain.BlockChain, index *db.RocksDB, workersCount int) *Worker {
	return &Worker{
		bc:           bc,
		db:           index,
		workersCount: workersCount,
	}
}


func (w *Worker) Sync(fromHeight, toHeight int) {

	h, _, _ := w.db.GetBestBlock()

	if fromHeight > toHeight {
		toHeight = fromHeight
	}

	if h == 0 {
		fmt.Println("No block in the database, this is a fresh sync")

		// sync block from startHeight to the remote best block

		rpcBestBlock, err := w.bc.GetBestBlockHeight()

		if err != nil {
			log.Fatal("Failed to get rpcBestBlock for Fresh sync")
		}

		log.Println("Remote RPC block best height is", rpcBestBlock)

	}

}
