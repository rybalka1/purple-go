package main

import (
	"encoding/json"
	"os"
	"sync"
)

var mu sync.Mutex

const storageFile = "storage.json"

type Verification struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

func saveVerification(v Verification) error {
	mu.Lock()
	defer mu.Unlock()

	var list []Verification
	_ = loadAll(&list)
	list = append(list, v)

	data, _ := json.MarshalIndent(list, "", "  ")
	return os.WriteFile(storageFile, data, 0644)
}

func loadAll(out *[]Verification) error {
	data, err := os.ReadFile(storageFile)
	if err != nil {
		return nil
	}
	return json.Unmarshal(data, out)
}

func findAndRemove(hash string) (bool, error) {
	mu.Lock()
	defer mu.Unlock()

	var list []Verification
	_ = loadAll(&list)

	found := false
	newList := []Verification{}
	for _, v := range list {
		if v.Hash == hash {
			found = true
			continue
		}
		newList = append(newList, v)
	}

	data, _ := json.MarshalIndent(newList, "", "  ")
	err := os.WriteFile(storageFile, data, 0644)
	return found, err
}
