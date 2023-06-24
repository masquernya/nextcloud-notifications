package storage

import (
	"encoding/json"
	"os"
)

type Storage struct {
	LoginUsername string
	LoginPassword string

	AdminUsername string
	AdminPassword string
	NotifiedTasks map[string]bool
}

var storage *Storage

func Get() *Storage {
	if storage == nil {
		loadStorage()
	}
	return storage
}

func Save() {
	bits, err := json.Marshal(storage)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("storage.json", bits, 0644)
	if err != nil {
		panic(err)
	}
}

func loadStorage() {
	_, err := os.Stat("storage.json")
	if err != nil {
		if os.IsNotExist(err) || os.IsPermission(err) {
			storage = &Storage{}
			return
		}
		panic(err)
	}
	bits, err := os.ReadFile("storage.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bits, &storage)
	if err != nil {
		panic(err)
	}

}

func TrySetNotified(id string) bool {
	me := Get()
	if me.NotifiedTasks == nil {
		me.NotifiedTasks = make(map[string]bool)
	}
	if _, ok := me.NotifiedTasks[id]; ok {
		return false
	}
	me.NotifiedTasks[id] = true
	Save()
	return true
}
