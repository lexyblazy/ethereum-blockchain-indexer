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
	rpcNodeUrl := flag.String("rpcnodeurl", "", "the rpc node url")
	port := flag.String("port", "", "the server port")
	dataDir := flag.String("datadir", "", "path to database directory")
	dbCache := flag.Int("dbcache", 1<<29, "size of the rocksdb cache")
	fromHeight := flag.Int("fromHeight", -1, "where to start indexing blocks from")
	toHeight := flag.Int("toHeight", -1, "where to stop")

	flag.Parse()

	if len(*port) == 0 {
		log.Fatal("port is missing")
	}

	index, err := db.NewConn(*dataDir, *dbCache)

	if err != nil {
		log.Fatal("Failed to open database connection", err.Error())
	}

	bc := bchain.NewBlockChain(*rpcNodeUrl)

	syncWorker := worker.NewWorker(bc, index, *rpcNodeUrl)
	syncWorker.Sync(*fromHeight, *toHeight)

	server.Start(index, *port)

}
