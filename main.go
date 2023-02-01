package main

import (
	"flag"
	"log"
	// "fmt"
	"ganache/indexer/server"
)

func main() {

	rpcNodeUrlPtr := flag.String("rpcnodeurl", "", "the rpc node url")
	portPtr := flag.String("port", "", "the server port")

	flag.Parse()

	rpcNodeUrl := *rpcNodeUrlPtr
	port := *portPtr

	if len(port) == 0 {
		log.Fatal("port is missing")
	}

	server.Start(rpcNodeUrl, port)
}
