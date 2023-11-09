package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// Transaction defines a blockchain transaction
type Transaction struct {
	Sender   string
	Receiver string
	Amount   float64
}

// Block defines a block in the blockchain
type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	Nonce        int
	PrevHash     string
	Hash         string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block
var mutex = &sync.Mutex{}

// calculateHash returns the hash of a block
func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%d%s", block.Index, block.Timestamp, block.Nonce, block.PrevHash)
	for _, tx := range block.Transactions {
		record += fmt.Sprintf("%s%s%f", tx.Sender, tx.Receiver, tx.Amount)
	}
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// generateBlock creates a new block using previous block's hash
func generateBlock(oldBlock Block, transactions []Transaction) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Transactions = transactions
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Nonce = 0 // Here you would add a real nonce finding process with Proof of Work

	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

// isBlockValid makes sure block is valid by checking index, and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// addTransactionToBlock adds a transaction to the block
func addTransactionToBlock(block *Block, transaction Transaction) {
	block.Transactions = append(block.Transactions, transaction)
}

// P2P Network Implementation

func handleConn(conn net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()

		var newBlock Block
		if err := json.Unmarshal([]byte(msg), &newBlock); err != nil {
			log.Printf("Failed to unmarshal received block: %v", err)
			continue
		}

		if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
			Blockchain = append(Blockchain, newBlock)
			fmt.Printf("Added block %d to the blockchain\n", newBlock.Index)
			bytes, err := json.MarshalIndent(Blockchain, "", "  ")
			if err != nil {
				log.Printf("Failed to marshal blockchain: %v", err)
				continue
			}
			fmt.Printf("\nBlockchain: %s\n", string(bytes))
		} else {
			fmt.Printf("Received invalid block\n")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error during connection: %v", err)
	}

	conn.Close()
}

func startServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Printf("Listening on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func main() {
	// create the genesis block
	genesisBlock := Block{0, time.Now().String(), []Transaction{}, 0, "", ""}
	genesisBlock.Hash = calculateHash(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	go startServer("8080")

	fmt.Println("Enter a new BPM:")
	var bpm int
	for {
		fmt.Scanf("%d\n", &bpm)

		// let's add a block
		transactions := []Transaction{
			{Sender: "A", Receiver: "B", Amount: float64(bpm)},
		}
		newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], transactions)
		if err != nil {
			fmt.Println("Error Generating Block")
			continue
		}
		if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
			Blockchain = append(Blockchain, newBlock)
		}

		bytes, err := json.Marshal(newBlock)
		if err != nil {
			log.Printf("Failed to marshal new block: %v", err)
			continue
		}

		// This should broadcast to all nodes but here we'll simulate by printing to stdout
		fmt.Printf("\nSending new block: %s\n", string(bytes))
	}
}
