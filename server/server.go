package server

import (
	"fmt"
	"indexer/db"
	// "context"
	// "encoding/json"
	// "errors"
	"log"
	"net/http"

	// "strconv"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Server struct {
	http    *http.Server
	handler *ApiHandler
	db      *db.RocksDB
}

func Start(ec *ethclient.Client, db *db.RocksDB, port string) {
	serveMux := http.NewServeMux()

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
		db:      db,
	}

	serverInstance.loadRoutes()

	log.Println("Server is up and running on ", port, "🚀🚀🚀")
	log.Fatal(serverInstance.http.ListenAndServe())
}

func (s *Server) loadRoutes() {
	serveMux := s.http.Handler.(*http.ServeMux)

	serveMux.HandleFunc("/", s.handler.jsonWrapper(func(r *http.Request) (interface{}, int, error) {
		return "Hello there", http.StatusOK, nil
	}))

	serveMux.HandleFunc("/api/chainId", s.handler.jsonWrapper(s.handler.getChainId))

	serveMux.HandleFunc("/api/balance/", s.handler.jsonWrapper(s.handler.getAccountBalance))
	serveMux.HandleFunc("/api/nonce/", s.handler.jsonWrapper(s.handler.getNonce))

	serveMux.HandleFunc("/api/block/", s.handler.jsonWrapper(s.handler.getBlock))

	serveMux.HandleFunc("/api/tx", s.handler.jsonWrapper(s.handler.createTransaction))

	serveMux.HandleFunc("/api/tx/", s.handler.jsonWrapper(s.handler.getTransaction))

	serveMux.HandleFunc("/api/tx-receipt/", s.handler.jsonWrapper(s.handler.getTransactionReceipt))
	serveMux.HandleFunc("/api/gas-price", s.handler.jsonWrapper(s.handler.getCurrentGasPrice))

}
