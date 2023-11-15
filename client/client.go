package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// Transaction is the data structure to hold transaction data
type Transaction struct {
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
}

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Create a scanner to read input from console
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter Sender, Receiver, and Amount separated by space (e.g., Alice Bob 10):")

	// Get input from the user
	for scanner.Scan() {
		input := scanner.Text()

		// Split the input to get sender, receiver, and amount
		var sender, receiver string
		var amount float64
		_, err := fmt.Sscanf(input, "%s %s %f", &sender, &receiver, &amount)
		if err != nil {
			fmt.Println("Invalid input. Please follow the format 'Sender Receiver Amount'. Error:", err)
			continue
		}

		// Create a transaction and marshal it into JSON
		transaction := Transaction{Sender: sender, Receiver: receiver, Amount: amount}
		jsonData, err := json.Marshal(transaction)
		if err != nil {
			fmt.Println("Error marshaling transaction:", err)
			continue
		}

		// Send JSON data to the server
		_, err = conn.Write(jsonData)
		if err != nil {
			fmt.Println("Error sending data to server:", err)
			continue
		}

		// Wait for the server to respond back
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading data from server:", err)
			continue
		}

		fmt.Println("Server response:", string(buffer[:n]))

		// Prompt the user for another input
		fmt.Println("Enter another transaction or CTRL+C to exit:")
	}
}
