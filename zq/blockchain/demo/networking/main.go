package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/joho/godotenv"
)

var Blockchain = []Block{}

type Block struct {
	Index     int
	Timestamp string
	BMP       int
	Hash      string
	PrevHash  string
}

var bcServer = make(chan []Block)
var mux = &sync.Mutex{}

// 计算hash
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BMP) + block.PrevHash
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

// 校验
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

// 替换新链
func replaceChain(newBlockchain []Block) {
	mux.Lock()
	if len(newBlockchain) > len(Blockchain) {
		Blockchain = newBlockchain
	}
	mux.Unlock()
}

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}

	genesisBlock := Block{0, time.Now().String(), 0, "", ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	httpAddr := os.Getenv("ADDR")

	l, err := net.Listen("tcp", ":"+httpAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on", httpAddr)
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, "Enter a BMP: ")
	scanner := bufio.NewScanner(conn)
	go func() {
		for {
			for scanner.Scan() {
				bmp, err := strconv.Atoi(scanner.Text())
				if err != nil {
					fmt.Printf("%v is not a int: %v\n", scanner.Text(), err)
					continue
				}
				newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], bmp)
				if err != nil {
					fmt.Printf("generete block error: %v\n", err)
					continue
				}
				if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
					newBlockchain := append(Blockchain, newBlock)
					replaceChain(newBlockchain)
				} else {
					fmt.Println("block is not valid")
				}

				bcServer <- Blockchain
				io.WriteString(conn, "\nEnter a BMP: ")
			}
		}
	}()

	go func() {
		for {
			time.Sleep(30 * time.Second)
			mux.Lock()
			bytes, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			mux.Unlock()
			io.WriteString(conn, string(bytes))
		}

	}()

	for _ = range bcServer {
		spew.Dump(Blockchain)
	}
}
