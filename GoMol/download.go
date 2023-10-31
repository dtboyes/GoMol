package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// downloadPDB downloads a PDB file from the RCSB PDB database using an HTTP request
// eliminates the need for the user to manually download PDB files
func downloadPDB(url, destination string) error {

	files, _ := os.ReadDir(destination)
	for _, file := range files {
		if file.Name() == destination {
			fmt.Println("File already exists")
			return nil
		}
	}

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Create the file
	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	fmt.Printf("PDB file downloaded successfully to: %s\n", destination)
	return nil
}
