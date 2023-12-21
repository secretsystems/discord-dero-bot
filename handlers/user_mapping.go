package handlers

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var userMappings map[string]string
var userMappingsMutex sync.Mutex

func loadUserMappings() {
	file, err := os.OpenFile("userMappings.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Error opening user mappings file: %v", err)
		return
	}
	defer file.Close()

	// Check if the file is empty before decoding
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info: %v", err)
		return
	}
	if fileInfo.Size() == 0 {
		userMappings = make(map[string]string)
		return
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&userMappings)
	if err != nil {
		log.Printf("Error decoding user mappings: %v", err)
	}
}

func saveUserMappings() {
	file, err := os.Create("userMappings.json")
	if err != nil {
		log.Printf("Error creating user mappings file: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(userMappings)
	if err != nil {
		log.Printf("Error encoding user mappings: %v", err)
	}
}
