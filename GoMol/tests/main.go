package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	// Replace these values with your actual query and database information
	query := ">query\nPROTEINSEQUENCE"
	database := "nr"

	// Define the BLAST service endpoint
	blastURL := "https://blast.ncbi.nlm.nih.gov/Blast.cgi"

	// Build the POST request body
	requestBody := fmt.Sprintf("CMD=Put&PROGRAM=blastp&DATABASE=%s&QUERY=%s", database, query)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", blastURL, bytes.NewBufferString(requestBody))
	if err != nil {
		log.Fatal("Failed to create HTTP request:", err)
	}

	// Set the content type for the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Failed to perform HTTP request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response body:", err)
	}

	// Print the BLAST result
	fmt.Println("BLAST result:")
	fmt.Println(string(body))
}
