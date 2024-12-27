package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// SaveToFile saves data to a file in JSON format
func SaveToFile(filePath string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	fmt.Printf("Data saved to %s\n", filePath)
	return nil
}

// LoadFromFile loads data from a JSON file
func LoadFromFile(filePath string, result interface{}) error {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	err = json.Unmarshal(fileData, result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	fmt.Printf("Data loaded from %s\n", filePath)
	return nil
}
