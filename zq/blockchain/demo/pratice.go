package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/davecgh/go-spew/spew"

	"github.com/gorilla/mux"
)

var Blockchain = []Block{}

type Block struct {
	Index     int
	Timestamp string
	BMP       int
	Hash      string
	PrevHash  string
}

// 计算hash
func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BMP) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// 生成块
func generateBlock(oldBlock Block, BMP int) (Block, error) {
	var block Block

	block.Index = oldBlock.Index + 1
	block.Timestamp = time.Now().String()
	block.BMP = BMP
	block.PrevHash = oldBlock.Hash
	block.Hash = calculateHash(block)

	return block, nil
}

// 校验块
func isBlockValid(newBlock, oldBlock Block) bool {
	if newBlock.Index != oldBlock.Index+1 {
		return false
	}
	if newBlock.Hash != calculateHash(newBlock) {
		return false
	}
	if newBlock.PrevHash != oldBlock.Hash {
		return false
	}
	return true
}

// 块取代
func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

// web服务
func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", httpAddr)
	s := http.Server{
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
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	spew.Dump(Blockchain)
	io.WriteString(w, string(bytes))
}

type Message struct {
	BMP  int
	what string
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], m.BMP)
	if err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, m)
		return
	}
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		spew.Dump(newBlock)
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	bytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(code)
		w.Write([]byte("http 500: internal error"))
		return
	}
	w.WriteHeader(code)
	w.Write(bytes)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		genesisBlock := Block{0, time.Now().String(), 0, "", ""}
		Blockchain = append(Blockchain, genesisBlock)
		spew.Dump(genesisBlock)
	}()
	log.Fatal(run())
}
