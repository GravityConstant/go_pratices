package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"github.com/davecgh/go-spew/spew"

	"github.com/gorilla/mux"
)

var (
	Blockchain = []Block{}
	mutex      = &sync.Mutex{}
	difficulty = 1
)

type Block struct {
	Index      int
	Timestamp  string
	BMP        int
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      string
}

type Message struct {
	BMP int
}

// calculate hash
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BMP) + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// generate block
func generateBlock(oldBlock Block, BMP int) Block {
	var newBlock Block

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.BMP = BMP
	newBlock.Difficulty = difficulty
	newBlock.PrevHash = oldBlock.Hash

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(calculateHash(newBlock), " do more work!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(calculateHash(newBlock), " work done!")
			newBlock.Hash = calculateHash(newBlock)
			break
		}
	}

	return newBlock
}

// hash valid
func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

// check
func isBlockValid(newBlock, oldBlock Block) bool {
	if newBlock.Index != oldBlock.Index+1 {
		return false
	}
	if newBlock.PrevHash != oldBlock.Hash {
		return false
	}
	if newBlock.Hash != calculateHash(newBlock) {
		return false
	}
	return true
}

func run() error {
	router := makeRouter()
	httpAddr := os.Getenv("ADDR")
	server := http.Server{
		Addr:           ":" + httpAddr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Listening on", httpAddr)
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func makeRouter() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", handleGetBlockchain).Methods("GET")
	router.HandleFunc("/", handleWriteBlock).Methods("POST")

	return router
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m Message

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		log.Printf("json error: %v\n", err)
		responseWithJson(w, r, http.StatusInternalServerError, r.Body)
		return
	}
	fmt.Println("message:", m)
	mutex.Lock()
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], m.BMP)
	mutex.Unlock()
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		Blockchain = append(Blockchain, newBlock)
		spew.Dump(newBlock)
	}

	responseWithJson(w, r, http.StatusCreated, newBlock)
}

func responseWithJson(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		genesis := Block{0, time.Now().String(), 0, "", "", difficulty, ""}
		spew.Dump(genesis)

		mutex.Lock()
		Blockchain = append(Blockchain, genesis)
		mutex.Unlock()
	}()
	log.Fatal(run())
}
