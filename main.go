package main

import (
	"flag"
	"indexer/bchain"
	"indexer/db"
	"indexer/server"
	"indexer/worker"
	"log"
)

func main() {
	rpcNodeUrl := flag.String("rpcnodeurl", "http://127.0.0.1:7545", "the rpc node url")
	port := flag.String("port", "", "the server port")
	dataDir := flag.String("datadir", "", "path to database directory")
	dbCache := flag.Int("dbcache", 1<<29, "size of the rocksdb cache")
	fromHeight := flag.Int("fromHeight", -1, "where to start indexing blocks from")
	toHeight := flag.Int("toHeight", -1, "where to stop")
	workers := flag.Int("workers", 8, "number of workers to process blocks in bulk mode")

	flag.Parse()

	if len(*port) == 0 {
		log.Fatal("port is missing")
	}

	index, err := db.NewConn(*dataDir, *dbCache)

	if err != nil {
		log.Fatal("Failed to open database connection", err.Error())
	}

	bc := bchain.NewBlockChain(*rpcNodeUrl)

	syncWorker := worker.NewWorker(bc, index, *workers)
	syncWorker.Sync(*fromHeight, *toHeight)

	s := server.NewServer(index, bc, *port)
	s.Start()

}
