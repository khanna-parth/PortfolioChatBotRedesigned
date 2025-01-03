package link

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"sync"
)

type ConnectionStore struct {
	Pool []ChatConnection
	Lock sync.RWMutex
	UploadsDir string
	Script string
	Config map[string]string
	IPRequestCount map[string]int
	Presets map[string]string
}

func NewConnectionStore(config map[string]string) *ConnectionStore {
	return &ConnectionStore{
		Pool: []ChatConnection{},
		Lock: sync.RWMutex{},
		IPRequestCount: map[string]int{},
		// Presets: map[string]string{"PRESET-MAIN": config["PRESET_MAIN"]},
	}
}

func (connStore *ConnectionStore) DemoExists(srcIP string) (existChat *ChatConnection) {
	connStore.Lock.Lock()
	defer connStore.Lock.Unlock()

	for i := range connStore.Pool {
		if connStore.Pool[i].ConnIP == srcIP {
            return &connStore.Pool[i]
        }
	}
	return nil
}

func (connStore *ConnectionStore) PurgeDemoes() {
    connStore.Lock.Lock()
    defer connStore.Lock.Unlock()
    
    for idx := len(connStore.Pool) - 1; idx >= 0; idx-- {
        existingChat := &connStore.Pool[idx]
        if existingChat.IsDemo {
            connStore.Pool = append(connStore.Pool[:idx], connStore.Pool[idx+1:]...)
        }
    }
}

func (connStore *ConnectionStore) RegisterConnection(chat *ChatConnection) {
	connStore.Lock.Lock()
	defer connStore.Lock.Unlock()
	connStore.Pool = append(connStore.Pool, *chat)
	fmt.Printf("Registered: %+v\n", chat)
}

func (connStore *ConnectionStore) RemoveConnection(connID string) {
	connStore.Lock.Lock()
	defer connStore.Lock.Unlock()
	for i := range connStore.Pool {
		conn := &connStore.Pool[i]
		if conn.ConnID == connID {
			connStore.Pool = append(connStore.Pool[:i], connStore.Pool[i+1:]...)
			fmt.Printf("Removed closed connID: %s\n", connID)
			return
		}
	}
	fmt.Printf("Did not find connID to terminate: %s\n", connID)
}

func (connStore *ConnectionStore) ConnectionExists(connID string) bool {
	connStore.Lock.RLock()
	defer connStore.Lock.RUnlock()

	for i := range connStore.Pool {
		if connStore.Pool[i].ConnID == connID {
			return true
		}
	}

	return false
}

func (connStore *ConnectionStore) MatchConnection(connID string, srcIP string) bool {
	if len(connStore.Pool) == 0 {
		fmt.Println("Matching predetermined to fail. Pool are empty")
		return false
	}
	connStore.Lock.RLock()
	defer connStore.Lock.RUnlock()

	for i := range connStore.Pool {
		// fmt.Printf("Checking ids '%s' and '%s'\n", conn.ConnID, connID)
		// fmt.Printf("Checking ips '%s' and '%s'\n", conn.ConnIP, srcIP)
		if connStore.Pool[i].ConnID == connID && connStore.Pool[i].ConnIP == srcIP {
			return true
		}
	}

	return false	
}

func (connStore *ConnectionStore) AvailableInPool(variance int) bool {
	limit := connStore.Config["MAX_CONNECTIONS"]
	intLimit, err := strconv.Atoi(limit)
	if err != nil || intLimit == 0 {
		return true
	}
	connStore.Lock.Lock()
	defer connStore.Lock.Unlock()
	return len(connStore.Pool) < (intLimit + variance)
}

func (connStore *ConnectionStore) IsRequestMaxReached(srcIP string) bool {
	connStore.Lock.Lock()
	defer connStore.Lock.Unlock()

	intLimit, err := strconv.Atoi(connStore.Config["MAX_REQUEST_PER_IP"])
	if err != nil {
		log.Printf("MAX_REQUEST_PER_IP could not be parsed")
		return true
	}

	for key, val := range connStore.IPRequestCount {
		if key == srcIP {
			// fmt.Printf("%s's request count is %d\n", key, val)
			return val >= intLimit
		}
	}
	// fmt.Println("%s was not found in ip request count list")
	return false
}

func (connStore *ConnectionStore) AddRequestCount(forIP string) {
	connStore.Lock.Lock()
	defer connStore.Lock.Unlock()

	curr := connStore.IPRequestCount[forIP]
	connStore.IPRequestCount[forIP] = curr + 1

	log.Printf("Request count for %s was incremented to %d", forIP, curr + 1)
}

func (connStore *ConnectionStore) IsPromptValid(prompt *string) (valid bool, lengthLimit int) {
	limit := connStore.Config["MAX_PROMPT_LENGTH"]
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		return false, 0
	}

	return len(*prompt) < intLimit, intLimit
}

func (connStore *ConnectionStore) GetConnection(connID string) *ChatConnection {
	for i := range connStore.Pool {
		if connID == connStore.Pool[i].ConnID {
			return &connStore.Pool[i]
		}
	}

	return nil
}

func (connStore *ConnectionStore) IsPreset(id string) (status bool, presetPath string) {
	if id == "STATIC" { return true, ""}
	fmt.Printf("presets: %+v\n", connStore.Presets)
	for key, val := range connStore.Presets {
		if key == id {
			return true, val
		}
	}

	fmt.Printf("%s is not a demo preset\n", id)
	return false, ""
}


func (connStore *ConnectionStore) GetUploadsPath() string {
	absUploads, err := filepath.Abs(connStore.UploadsDir)
	if err != nil {
		return connStore.UploadsDir
	}

	return absUploads
}