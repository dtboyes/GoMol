package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadPDB(url, destination string) error {
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

func main() {
	fmt.Print()
	pdbURL := "https://files.rcsb.org/download/" + os.Args[1] + ".pdb" // Replace 'xxxx' with the actual PDB code
	localPath := "pdbfiles/" + os.Args[1] + ".pdb"                     // Replace 'xxxx' with the actual PDB code

	err := downloadPDB(pdbURL, localPath)
	if err != nil {
		fmt.Println("Error downloading PDB file:", err)
	}
}
