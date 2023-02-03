package server

import (
	"fmt"
	// "context"
	// "encoding/json"
	// "errors"
	"log"
	"net/http"
	// "strconv"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Server struct {
	http    *http.Server
	handler *ApiHandler
}

const DEFAULT_RPC_NODE_URL = "http://127.0.0.1:7545"

func Start(rpcNodeUrl string, port string) {
	serveMux := http.NewServeMux()

	var nodeUrl string

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

	httpInterface := &http.Server{
		Addr:    port,
		Handler: serveMux,
	}

	ap := &ApiHandler{
		ethClient: ec,
	}

	serverInstance := &Server{
		http:    httpInterface,
		handler: ap,
	}

	serverInstance.loadRoutes()
	log.Println("Server is up and running on ", port, "🚀🚀🚀")
	log.Fatal(serverInstance.http.ListenAndServe())
}

func (s *Server) loadRoutes() {
	serveMux := s.http.Handler.(*http.ServeMux)

	serveMux.HandleFunc("/", s.handler.jsonWrapper(func(r *http.Request) (interface{}, error) {
		return "Hello there", nil
	}))

	serveMux.HandleFunc("/api/chainId", s.handler.jsonWrapper(s.handler.getChainId))

	serveMux.HandleFunc("/api/balance/", s.handler.jsonWrapper(s.handler.getAccountBalance))
	serveMux.HandleFunc("/api/nonce/", s.handler.jsonWrapper(s.handler.getNonce))

	serveMux.HandleFunc("/api/block/", s.handler.jsonWrapper(s.handler.getBlock))

	serveMux.HandleFunc("/api/tx", s.handler.jsonWrapper(s.handler.createTransaction))

	serveMux.HandleFunc("/api/tx/", s.handler.jsonWrapper(s.handler.getTransaction))

}

func dialRpc(url string) (*rpc.Client, *ethclient.Client, error) {
	rpcClient, err := rpc.Dial(url)

	if err != nil {
		return nil, nil, err
	}

	ethereumClient := ethclient.NewClient(rpcClient)

	return rpcClient, ethereumClient, err
}
