package controllers

import (
	"sync"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)
type ClientList struct {
	sync.RWMutex
	data map[*websocket.Conn]Client
}
func (clients *ClientList) Remove(ws *websocket.Conn) {
	clients.Lock()
	delete(clients.data, ws)
	clients.Unlock()
	log.Warn("Client disconnection")
}