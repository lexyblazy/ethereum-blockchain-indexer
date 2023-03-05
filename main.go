package main

import (
	"flag"
	// "fmt"
	"ganache/indexer/db"
	"ganache/indexer/server"
	"log"
)

func main() {
	rpcNodeUrl := flag.String("rpcnodeurl", "", "the rpc node url")
	port := flag.String("port", "", "the server port")
	dataDir := flag.String("datadir", "", "path to database directory")
	dbCache := flag.Int("dbcache", 1<<29, "size of the rocksdb cache")

	flag.Parse()

	if len(*port) == 0 {
		log.Fatal("port is missing")
	}

	db.NewConn(*dataDir, *dbCache)
	server.Start(*rpcNodeUrl, *port)
}
