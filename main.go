package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"indexer/db"
	"indexer/server"
	"indexer/worker"
	"log"
)

func newEthClient(rpcNodeUrl string) *ethclient.Client {
	var nodeUrl string
	const DEFAULT_RPC_NODE_URL = "http://127.0.0.1:7545"

	if len(rpcNodeUrl) == 0 {
		nodeUrl = DEFAULT_RPC_NODE_URL
	} else {
		nodeUrl = rpcNodeUrl
	}

	fmt.Println("NodeUrl", nodeUrl)
	_, ec, err := dialRpc(nodeUrl)

	if err != nil {
		log.Fatal("Failed to connect to RPC_NODE ", nodeUrl)
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

func main() {
	rpcNodeUrl := flag.String("rpcnodeurl", "", "the rpc node url")
	port := flag.String("port", "", "the server port")
	dataDir := flag.String("datadir", "", "path to database directory")
	dbCache := flag.Int("dbcache", 1<<29, "size of the rocksdb cache")

	flag.Parse()

	if len(*port) == 0 {
		log.Fatal("port is missing")
	}

	ec := newEthClient(*rpcNodeUrl)
	index, err := db.NewConn(*dataDir, *dbCache)

	if err != nil {
		log.Fatal("Failed to open database connection", err.Error())
	}

	syncWorker := worker.NewWorker(ec, index, *rpcNodeUrl)
	syncWorker.Sync()

	server.Start(ec, index, *port)

}
