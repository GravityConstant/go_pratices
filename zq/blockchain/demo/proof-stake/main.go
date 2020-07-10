package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/joho/godotenv"
)

var (
	Blockchain      = []Block{}
	tempBlocks      = []Block{}
	candidateBlocks = make(chan Block)
	validators      = make(map[string]int)
	announcements   = make(chan string)
	mutex           = &sync.Mutex{}
)

type Block struct {
	Index     int
	Timestamp string
	BMP       int
	Hash      string
	PrevHash  string
	Validator string
}

// calculate hash
func calculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// calculate block hash
func calculateBlockHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BMP) + block.PrevHash
	return calculateHash(record)
}

// generate block
func generateBlock(oldBlock Block, BMP int, address string) (Block, error) {
	var block Block

	block.Index = oldBlock.Index + 1
	block.Timestamp = time.Now().String()
	block.BMP = BMP
	block.PrevHash = oldBlock.Hash
	block.Hash = calculateBlockHash(block)
	block.Validator = address

	return block, nil
}

// is block valid
func isBlockValid(newBlock, oldBlock Block) bool {
	if newBlock.Index != oldBlock.Index+1 {
		return false
	}
	if newBlock.PrevHash != oldBlock.Hash {
		return false
	}
	if newBlock.Hash != calculateBlockHash(newBlock) {
		return false
	}
	return true
}

func pickWinner() {
	time.Sleep(30 * time.Second)
	mutex.Lock()
	temp := tempBlocks
	mutex.Unlock()

	lotteryPool := []string{}
	if len(temp) > 0 {
	OUTTER:
		for _, block := range temp {
			for _, node := range lotteryPool {
				if node == block.Validator {
					continue OUTTER
				}
			}
			// 查看是否投注了
			mutex.Lock()
			setValidators := validators
			mutex.Unlock()
			k, ok := setValidators[block.Validator]
			if ok {
				for i := 0; i < k; i++ {
					lotteryPool = append(lotteryPool, block.Validator)
				}
			}
		}

		// 摇奖
		if len(lotteryPool) > 0 {
			s := rand.NewSource(time.Now().Unix())
			r := rand.New(s)
			lotteryWinner := lotteryPool[r.Intn(len(lotteryPool))]
			// save to blockchain and pub
			for _, block := range temp {
				if block.Validator == lotteryWinner {
					mutex.Lock()
					Blockchain = append(Blockchain, block)
					mutex.Unlock()
					for range validators {
						announcements <- "\nwinning validator:" + lotteryWinner + "\n"
					}
					break
				}
			}
		}
	}
	// reset tempBlocks
	mutex.Lock()
	tempBlocks = []Block{}
	mutex.Unlock()
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// pub announcements
	go func() {
		msg := <-announcements
		io.WriteString(conn, msg)
	}()

	var address string
	// get balance
	io.WriteString(conn, "Enter a balance:")
	balanceScan := bufio.NewScanner(conn)
	for balanceScan.Scan() {
		balance, err := strconv.Atoi(balanceScan.Text())
		if err != nil {
			fmt.Printf("%v not a number: %v", balanceScan.Text(), err)
			return
		}
		// set validators
		address = calculateHash(time.Now().String())
		validators[address] = balance
		fmt.Println(validators)
		break
	}

	// get bpm
	io.WriteString(conn, "Enter a new BPM:")
	bpmScan := bufio.NewScanner(conn)
	go func() {
		for bpmScan.Scan() {
			bpm, err := strconv.Atoi(bpmScan.Text())
			if err != nil {
				fmt.Printf("%v not a number: %v", bpmScan.Text(), err)
				delete(validators, address)
				conn.Close()
			}
			// get old last block
			mutex.Lock()
			oldLastBlock := Blockchain[len(Blockchain)-1]
			mutex.Unlock()
			// generate block
			newBlock, err := generateBlock(oldLastBlock, bpm, address)
			if err != nil {
				log.Println(err)
				continue
			}
			// check valid
			if isBlockValid(newBlock, oldLastBlock) {
				candidateBlocks <- newBlock
			}
			io.WriteString(conn, "\nEnter a new BPM:")
		}
	}()

	// pub to clients
	for {
		time.Sleep(time.Minute)
		mutex.Lock()
		output, err := json.Marshal(Blockchain)
		mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(conn, string(output)+"\n")
	}
}

func main() {
	// set os env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// generate genesis block
	genesisBlock := Block{}
	genesisBlock = Block{0, time.Now().String(), 0, calculateBlockHash(genesisBlock), "", ""}
	Blockchain = append(Blockchain, genesisBlock)
	spew.Dump(Blockchain)

	go func() {
		// add to tempBlocks
		for candidate := range candidateBlocks {
			mutex.Lock()
			tempBlocks = append(tempBlocks, candidate)
			mutex.Unlock()
		}
	}()

	go func() {
		// produce a winner
		pickWinner()
	}()

	// listen and handle
	httpAddr := os.Getenv("PORT")
	server, err := net.Listen("tcp", ":"+httpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	log.Println("Listening on", httpAddr)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}

}
