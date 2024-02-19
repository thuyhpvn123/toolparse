package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"github.com/gorilla/websocket"
	c_config "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/config"
	cc_config "github.com/meta-node-blockchain/meta-node/cmd/chiabai/config"
	"github.com/meta-node-blockchain/meta-node/cmd/chiabai/core"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/syndtr/goleveldb/leveldb"

)
var levelDb *leveldb.DB

type Server struct {
	sync.Mutex
	contractABI map[string]*core.ContractABI
	config      *c_config.ClientConfig
	clients           ClientList

}
type Message1 struct {
	Command string      `json:"command"`
	Data    interface{} `json:"data"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (server *Server) Init(config *c_config.ClientConfig) *Server {
	// init subscriber
	server.config = config
	server.contractABI = make(map[string]*core.ContractABI)
	var wg sync.WaitGroup
	for _, contract := range core.Contracts {
		wg.Add(1)
		go server.getABI(&wg, contract)
	}
	wg.Wait()
	// connected clients
	server.clients.data = make(map[*websocket.Conn]Client)
	//Open levelDb
	fmt.Println("the end")

	return &Server{
		contractABI: server.contractABI,
		config:      config,
	}
}


func (server *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	config, err := c_config.LoadConfig(cc_config.CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.ClientConfig)
	client := Client{
		ws: conn,
		server: server,
		config: cConfig,
		cliMap:    make(map[string]*Cli),
		sendChan: make(chan Message1,1),
	}
	client.init()
	
	log.Println("Client Connected successfully") //write on server terminal
	server.clients.Lock()
	server.clients.data[conn] = client
	fmt.Println("server.clients.data:",server.clients.data)
	server.clients.Unlock()
	defer server.clients.Remove(conn)
	//listen websocket
	
	client.handleListen()
}

func (server *Server) getABI(wg *sync.WaitGroup, contract core.Contract) {
	var temp core.ContractABI
	temp.InitContract(contract)
	server.Lock()
	server.contractABI[contract.Name] = &temp
	server.Unlock()
	wg.Done()
}
func (server *Server)InitDB(){
	leveldb, err := leveldb.OpenFile("./db/keys", nil)
	if err != nil {
		panic(err)
	}
	levelDb=leveldb
}
