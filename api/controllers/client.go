package controllers

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	// "github.com/meta-node-blockchain/meta-node/cmd/chiabai/config"
	log "github.com/sirupsen/logrus"
	"github.com/meta-node-blockchain/meta-node/cmd/client/pkg/config"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"

	// controller_client "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/controllers"
	// "github.com/meta-node-blockchain/meta-node/pkg/bls"
	// "github.com/meta-node-blockchain/meta-node/pkg/network"
	// "github.com/meta-node-blockchain/meta-node/pkg/state"
)

type Client struct {
	ws     *websocket.Conn
	server *Server
	cliMap map[string]*Cli
	sync.Mutex
	config *config.ClientConfig
	sendChan                 chan Message1

}

func (client *Client) init() {
	go client.handleMessage()
	log.Info("End init client")
}
func (client *Client) handleListen() {
	for {
		// Read in a new message as JSON and map it to a Message object
		var msg map[string]interface{}
		err := client.ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			client.ws.Close()
			client.server.clients.Remove(client.ws)
			break
		}
		if msg["command"] == "ping" {
			err := client.ws.WriteJSON(Message1{Command: "test", Data: "Pong"})
			if err != nil {
				logger.Error("error in sendMessage with ws ", err)
				return
			}
		}
		// log.Info("Message from client: ", msg)
		client.handleCallChain(msg)
	}
}

// handle message struct tu chain tra ve va chuyen qua dang JSON gui toi cac client
func (client *Client) handleMessage() {
	client.ws.WriteJSON(
		Message1{Command: "message", Data: "Welcome "})
	for {
		msg := <-client.sendChan
		// msg1 := <-sendDataC
		log.Info(msg)
		err := client.ws.WriteJSON(msg)

		if err != nil {
			log.Printf("error: %v", err)
			client.ws.Close()
			client.server.clients.Remove(client.ws)

		}
	}
}
func (client *Client) handleCallChain(msg map[string]interface{}) {
	fmt.Println("msg là:", msg)
	// handle call
	map1,ok:=msg["value"]
	if !ok{
		logger.Error(fmt.Sprintf("error when handleCallChain map value"))
		return
	}	
	switch msg["command"] {
	case "test":
		go client.sentToClient("test", "ok")
	case "connect-wallet":

		call,ok := map1.(map[string]interface{})
		if !ok{
			logger.Error(fmt.Sprintf("error when handleCallChain map roomid"))
			return
		}			
		client.ConnectWallet(call)
	// case "connect-wallet-loop":
	// 	call,ok := map1.(map[string]interface{})
	// 	if !ok{
	// 		logger.Error(fmt.Sprintf("error when handleCallChain map roomid"))
	// 		return
	// 	}		
	// 	filename,_:=call["filename"].(string)	
	// 	client.ConnectWalletLoop(filename)
	// case "get-code-ref":
	// 	call,ok := map1.(map[string]interface{})
	// 	if !ok{
	// 		logger.Error(fmt.Sprintf("error when handleCallChain map roomid"))
	// 		return
	// 	}		
	// 	client.GetCodeRef(call)

	case "call":
	
		call,ok := map1.(map[string]interface{})
		if !ok{
			logger.Error(fmt.Sprintf("error when call transaction"))
			return
		}
		client.TryCall(call)
	case "getUserInfo":
		call,ok := map1.(map[string]interface{})
		if !ok{
			logger.Error(fmt.Sprintf("error when getCodeList"))
			return
		}
		client.GetUserInfo(call)
	case "viewtreeMatrix":
		call,ok := map1.(map[string]interface{})
		if !ok{
			logger.Error(fmt.Sprintf("error when viewtreeMatrix"))
			return
		}
		client.ViewtreeMatrix(call)

	default:
		log.Warn("Require call not match: ", msg)
	}
}
func (client *Client) sentToClient(command string, data interface{}) {
	client.sendChan <- Message1{command, data}
}