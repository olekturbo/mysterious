package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run spam_client.go 'Your message'")
		return
	}
	message := os.Args[1]
	payload := map[string]string{"message": message}
	body, _ := json.Marshal(payload)

	resp, err := http.Post("http://localhost:5000/predict", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("Error closing response body: %v\n", closeErr)
		}
	}()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println(string(respBody))
}
