package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"samplechain/blockchain"
	"time"
)

type TransactionMessage struct {
	From  string
	To    string
	Value int
}

func Run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", os.Getenv("ADDR"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/transaction", handleTransactionBlock).Methods("POST")
	muxRouter.HandleFunc("/balance/{addr}", handleBalance).Methods("GET")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {

	bytes, err := json.MarshalIndent(blockchain.GetBlockChains(), "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleTransactionBlock(w http.ResponseWriter, r *http.Request) {
	var message TransactionMessage

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))

	err = json.Unmarshal(body, &message)
	if err != nil {
		panic(err)
	}
	log.Println(message)
	defer r.Body.Close()

	blockchain.Send(message.From, message.To, message.Value)

	respondWithJSON(w, r, http.StatusCreated, blockchain.GetBlockChains())

}

func handleBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["addr"]

	balance := blockchain.GetBalance(address)
	response := fmt.Sprintf(`{"balance":"%d"}`, balance)

	respondWithJSON(w, r, http.StatusOK, response)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
