package link

import (
	"fmt"
	"sync"
)

type ChatConnection struct {
	ConnID string
	ConnIP string
	DocCount int
	InUse bool
	Linker *Linker
	Executor *Executor
	Lock sync.Mutex
	IsDemo bool
}

func NewChatDemo(uid string, srcIP string, connStore *ConnectionStore) *ChatConnection {
	existingChat := connStore.DemoExists(srcIP)
	if existingChat != nil {
		fmt.Println("Found existing chat")
		return existingChat
	}

	newChat := &ChatConnection{
		ConnID: uid,
		ConnIP: srcIP,
		Linker: CreateLinker(uid, connStore.UploadsDir),
		Executor: NewExecutor(uid),
		IsDemo: true,
	}
	fmt.Println("Returning new connection")
	connStore.RegisterConnection(newChat)
	return newChat
}

func NewChatConnection(uid string, srcIP string, connStore *ConnectionStore) *ChatConnection {
	return &ChatConnection{
		ConnID: uid,
		ConnIP: srcIP,
		Linker: CreateLinker(uid, connStore.UploadsDir),
		Executor: NewExecutor(uid),
	}
}

func (chat *ChatConnection) ToggleUse() {
	chat.InUse = !chat.InUse
}