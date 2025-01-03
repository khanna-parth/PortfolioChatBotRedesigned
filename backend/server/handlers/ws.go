package handlers

import (
	"fmt"
	"log"
	// "net"
	"net/http"
	"server/helper"
	"server/link"
	"time"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request, connStore *link.ConnectionStore) {
	if !connStore.AvailableInPool(0) {
		http.Error(w, "All server resources are in use", http.StatusServiceUnavailable)
		fmt.Println("Denying websocket. All resources are full")
		return
	}

	srcIP, _ := helper.GetClientData(r)
	if srcIP == "" {
		fmt.Println("Denying websocket. No srcIP")
		http.Error(w, "Could not extract source", http.StatusInternalServerError)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket: ", err)
		log.Println("Error upgrading to websocket: ", err)
		return
	}

	uid := uuid.NewString()

	if srcIP == "::1" {
		uid = "pkhanna"
	}

		
	presetErr, _ := connStore.IsPreset(uid)
	for presetErr {
		uid = uuid.NewString()
		presetErr, _ = connStore.IsPreset(uid)
	}

	currConn := link.NewChatConnection(uid, srcIP, connStore)

	connStore.RegisterConnection(currConn)

	defer func() {
		conn.Close()
		currConn.Linker.Clean()
		connStore.RemoveConnection(currConn.ConnID)
	}()

	helper.SendWebsocket(conn, map[string]string{"connID": currConn.ConnID})

	log.Printf("New websocket connection established for %s\n", srcIP)
	
	timeout := 2 * time.Hour
	timer := time.NewTimer(timeout)
	_ = timer

	// go func() {
	// 	<-timer.C
	// 	log.Println("WebSocket connection expired. Closing.")
	// 	conn.Close()
	// }()


	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		fmt.Printf("Received: %s\n", p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}