package server

import (
	"indexer/bchain"
	"indexer/db"
	// "context"
	// "encoding/json"
	// "errors"
	"log"
	"net/http"
	// "strconv"
)

type Server struct {
	http    *http.Server
	handler *ApiHandler
	db      *db.RocksDB
}

func NewServer(db *db.RocksDB, bc *bchain.BlockChain, port string) *Server {
	serveMux := http.NewServeMux()

	httpInterface := &http.Server{
		Addr:    port,
		Handler: serveMux,
	}

	ap := &ApiHandler{
		bc: bc,
	}

	return &Server{
		http:    httpInterface,
		handler: ap,
		db:      db,
	}

}

func (s *Server) Start() {
	s.loadJsonAPIRoutes()
	log.Println("Server is up and running on ", s.http.Addr, "🚀🚀🚀")
	log.Fatal(s.http.ListenAndServe())
}

func (s *Server) loadJsonAPIRoutes() {
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
