package controllers

import (
	// "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/meta-node-blockchain/meta-node/cmd/client/command"
	"github.com/meta-node-blockchain/meta-node/cmd/client/pkg/client_context"
	cc_config "github.com/meta-node-blockchain/meta-node/cmd/chiabai/config"
	c_config "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/config"

	"github.com/meta-node-blockchain/meta-node/cmd/client/pkg/controllers"
	c_network "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/network"
	client_types "github.com/meta-node-blockchain/meta-node/cmd/client/types"
	log "github.com/sirupsen/logrus"

	// "github.com/ethereum/go-ethereum/crypto"
	// "github.com/meta-node-blockchain/meta-node/cmd/chiabai/core"
	"github.com/meta-node-blockchain/meta-node/pkg/bls"
	p_common "github.com/meta-node-blockchain/meta-node/pkg/common"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/pkg/network"
	p_network "github.com/meta-node-blockchain/meta-node/pkg/network"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	p_transaction "github.com/meta-node-blockchain/meta-node/pkg/transaction"
	"github.com/meta-node-blockchain/meta-node/types"
)

// var cli *Cli
var defaultRelatedAddress [][]byte

var (
	connectionTypesForClient = []string{p_common.NODE_CONNECTION_TYPE}
)

type Cli struct {
	clientContext         *client_context.ClientContext
	mu                    sync.Mutex
	accountStateChan      chan types.AccountState
	receiptChan           chan types.Receipt
	chData                chan interface{}
	transactionController client_types.TransactionController

	server *Server
}
type Account struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}
type Info struct {
	Address string `json:"address"`
	FirstTimePay uint `json:firstTimePay`
	NextTimePay uint `json:nextTimePay`
	Childrens []string `json:Childrens`
	ChildrensMatrix []string `json:ChildrensMatrix`
	Line string `json:Line`
	LineMatrix string `json:LineMatrix`
	MtotalMember uint `json:MtotalMember`
	Rank uint `json:Rank`
	totalSubcriptionBonus uint `json:totalSubcriptionBonus`
	LineMtotalMatrixBonusatrix uint `json:totalMatrixBonus`
	totalMatchingBonus uint `json:totalMatchingBonus`
	totalSaleBonus uint `json:totalSaleBonus`
	totalGoodSaleBonus uint `json:totalGoodSaleBonus`
	totalExtraDiamondBonus uint `json:LineMattotalExtraDiamondBonusrix`
	totalExtraCrownDiamondBonus uint `json:totalExtraCrownDiamondBonus`
	totalSale uint `json:totalSale`
}

func (client *Client) ConnectWallet(
	callMap map[string]interface{},
	// config *c_config.ClientConfig,
	// server *Server,
) (*Cli, error) {
	Cconfig, err := c_config.LoadConfig(cc_config.CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	config := Cconfig.(*c_config.ClientConfig)
	clientContext := &client_context.ClientContext{
		Config: config,
	}
	// defer clientContext.Close()
	cli := Cli{
		clientContext:    clientContext,
		accountStateChan: make(chan types.AccountState, 1),
		receiptChan:      make(chan types.Receipt, 1),
		chData:           make(chan interface{}, 1),
		// server: server,
	}
	addressString := callMap["address"].(string)
	client.cliMap[addressString] = &cli
	privateKey := callMap["privateKey"].(string)
	clientContext.KeyPair = bls.NewKeyPair(common.FromHex(privateKey))
	clientContext.MessageSender = p_network.NewMessageSender(clientContext.KeyPair, config.Version())
	clientContext.ConnectionsManager = network.NewConnectionsManager()
	parentConn := network.NewConnection(
		common.HexToAddress(config.ParentAddress),
		config.ParentConnectionType,
		config.DnsLink(),
	)
	fmt.Println("parentConn:",parentConn)
	clientContext.Handler = c_network.NewHandler(
		cli.accountStateChan,
		cli.receiptChan,
		// cli.chData,
		// cli.eventChan,
	)
	clientContext.SocketServer = network.NewSockerServer(
		config,
		clientContext.KeyPair,
		clientContext.ConnectionsManager,
		clientContext.Handler,
	)
	err = parentConn.Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("error when connect to parent %v", err))
		return nil, err
	} else {
		// init connection
		clientContext.ConnectionsManager.AddParentConnection(parentConn)
		clientContext.SocketServer.OnConnect(parentConn)
		go clientContext.SocketServer.HandleConnection(parentConn)
	}
	cli.transactionController = controllers.NewTransactionController(
		clientContext,
	)
	return &cli, nil
}

func writeJSONToFile(filename string, data []Account) error {
	filePath, _ := filepath.Abs("frontend/public/chiabai/frontend/js/test500k/branch3/" + filename + ".json")
	// Create db if it doesn't exist
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}
func getDatas(filename string) []Account {
	dat, _ := os.ReadFile("frontend/public/chiabai/frontend/js/test100k/branch3/" + filename + ".json")
	scDatas := []Account{}
	err := json.Unmarshal(dat, &scDatas)
	if err != nil {
		panic(err)
	}
	return scDatas
}
func (client *Client) TryCall(callMap map[string]interface{}) interface{} {
	var result interface{}
	result = "TimeOut"
	result = client.call(callMap)

	// if result != "TimeOut" {
	// 	log.Info(" - Result: ", result)
	// 	functionName, ok := callMap["function-name"].(string)
	// 	if !ok {
	// 		functionName = ""
	// 	}
	// 	contractName, ok := callMap["contract"].(string)
	// 	if !ok {
	// 		contractName = ""
	// 	}
	// 	contract := client.server.contractABI[contractName]
	// 	// kq := contract.Decode(functionName, result.(string)).(map[string]interface{})[""]
	// 	kq := contract.Decode(functionName, result.(string))

	// 	client.sentToClient(functionName,kq)
	// return kq
	// }

	return result
}
func (client *Client)GetUserInfo(callMap map[string]interface{}){
	var result interface{}
	result = client.TryCall(callMap)
	if result != "TimeOut" {
		log.Info(" - Result: ", result)
		functionName, ok := callMap["function-name"].(string)
		if !ok {
			functionName = ""
		}
		contractName, ok := callMap["contract"].(string)
		if !ok {
			contractName = ""
		}
		contract := client.server.contractABI[contractName]
		kq := contract.Decode(functionName, result.(string)).(map[string]interface{})["userinfo"]
		// kq := contract.Decode(functionName, result.(string))

		client.sentToClient(functionName,kq)
		writeToJsonFile("GetUserInfo",kq)

	}

}
func (client *Client)ViewtreeMatrix(callMap map[string]interface{}){
	var result interface{}
	result = client.TryCall(callMap)
	if result != "TimeOut" {
		log.Info(" - Result: ", result)
		functionName, ok := callMap["function-name"].(string)
		if !ok {
			functionName = ""
		}
		contractName, ok := callMap["contract"].(string)
		if !ok {
			contractName = ""
		}
		contract := client.server.contractABI[contractName]
		kq := contract.Decode(functionName, result.(string)).(map[string]interface{})[""]
		// kq := contract.Decode(functionName, result.(string))

		client.sentToClient(functionName,kq)
		writeToJsonFile("viewtreeMatrix",kq)

	}
	// //write to json file
	// // writeToJsonFile("viewtreeMatrix",kq)
	// fromAddress := callMap["from"].(string)
	// filename := callMap["function-name"].(string)
	// // var codes []CodeRef

	// var filePath = "frontend/public/chiabai/frontend/file/" + filename + "Call.json"

	// // Check if the file exists
	// coderef := client.cliMap[fromAddress].TryCall(callMap)
	// // Create a new CodeRef element to add.
	// newInfor := InfoUser{
	// 	Address: fmt.Sprintf("%v", fromAddress),
	// 	Code:    fmt.Sprintf("%v", coderef),
	// }
	// writeCodeRef(codes, filePath, filename, newCodeRef)
}
func writeCodeRef(infos []Info, filePath string, filename string, newCodeRef Info) {

	if _, err := os.Stat(filePath); os.IsNotExist(err) {

		// Add the new element to the slice.
		infos = append(infos, newCodeRef)
		// Call the function to write the JSON data to the file
		if err := writeJSONToFileCodeRef(filePath, filename, infos); err != nil {
			fmt.Println("Error writing JSON to file:", err)

		} else {
			fmt.Println("JSON data written to", filename)
		}

	} else if err != nil {
		// Other error occurred while checking
		fmt.Println("Error checking file:", err)
		return
	} else {
		fmt.Println("File already exists:", filePath)
		// Read the existing JSON data from the file
		existingData, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading JSON file:", err)
			return
		}
		// Parse the JSON data into a Go data structure
		err = json.Unmarshal(existingData, &infos)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}
		// Add the new element to the slice.
		infos = append(infos, newCodeRef)
		// Call the function to write the JSON data to the file
		if err := writeJSONToFileCodeRef(filePath, filename, infos); err != nil {
			fmt.Println("Error writing JSON to file:", err)
		} else {
			fmt.Println("JSON data written to", filename)
		}

	}

}
func writeJSONToFileCodeRef(filePath string, filename string, data []Info) error {
	// filePath, _ := filepath.Abs("frontend/public/chiabai/frontend/coderef/test100k/branch3/" + filename + "CodeRef.json")
	// Create db if it doesn't exist
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}
func writeToJsonFile(functionName string,kq interface{}){
	
		fPath1:="frontend/public/chiabai/frontend/file/" + functionName + "Call.json"
		mess:=Message1{
			Command: functionName,
			Data: kq,
		}
		jsonData1, err := json.MarshalIndent(mess, "", "    ")
		if err != nil {
			panic("Error marshaling JSON")
		}

		err = os.WriteFile(fPath1, jsonData1, 0644)
		if err != nil {
			panic("Error writing to file:"+functionName)
		}
		logger.Debug("Data has been successfully written ")
}
func (client *Client) call(callMap map[string]interface{}) interface{} {
	fromAddress := callMap["from"].(string)
	toAddressStr, _ := callMap["to"].(string)
	toAddress := common.HexToAddress(toAddressStr)
	isCall, _ := callMap["is-call"].(bool)
	var action pb.ACTION
	if isCall {

		action = pb.ACTION_CALL_SMART_CONTRACT
	}
	// action  :=pb.ACTION_EMPTY
	relatedAddress := client.EnterRelatedAddress(callMap)
	inputStr, ok := callMap["input"].(string)
	if !ok {
		inputStr = ""
	}

	hexAmount, _ := callMap["amount"].(string)
	if hexAmount == "" {
		hexAmount = "0"
	}
	amount := uint256.NewInt(0).SetBytes(common.FromHex(hexAmount))
	var maxGas uint64
	maxGaskq, ok := callMap["gas"].(float64)
	if !ok {
		maxGas = 2000000
	} else {
		maxGas = uint64(maxGaskq)

	}

	var maxGasPriceGwei uint64
	maxGasPriceGweikq, ok := callMap["gasPrice"].(float64)
	if !ok {
		maxGasPriceGwei = 10
	} else {
		maxGasPriceGwei = uint64(maxGasPriceGweikq)

	}
	maxGasPrice := 100000000 * maxGasPriceGwei

	var maxTimeUse uint64
	maxTimeUsekq, ok := callMap["timeUse"].(float64)
	if !ok {
		maxTimeUse = 1000
	} else {
		maxTimeUse = uint64(maxTimeUsekq)
	}

	var data []byte
	if len(inputStr) > 0 {
		callData := p_transaction.NewCallData(common.FromHex(inputStr))
		data, _ = callData.Marshal()
	} else {
		data, _ = client.GetDataForCallSmartContract(callMap)
		// if err != nil {
		// 	panic(err)
		// }
	}
	_, err := client.SendTransaction(
		fromAddress,
		toAddress,
		amount,
		action,
		data,
		relatedAddress,
		maxGas,
		maxGasPrice,
		maxTimeUse,
	)

	if err != nil {
		log.Warn(err)
	} else {
		logger.Info("Done send transaction from " + fromAddress)
	}
	// if hashed == nil {
	// 	fmt.Println("hashed==nil")
	// 	return "TimeOut1"
	// }
	for {

		select {

		case receiver := <-client.cliMap[fromAddress].receiptChan:
			// a:=cli.clientContext.Handler.(*c_network.Handler).GetReceiptChan()
			// log.Info("Hash on server", common.BytesToHash(hash.([]byte)))
			// log.Info("Hash from chain", (receiver).(network.Receipt).Hash)
			// if (receiver).(network.Receipt).Hash != common.BytesToHash(hash.([]byte)) {
			// 	continue
			// }
			if (receiver).Status() == 2 {
				fmt.Println("Status threw")
				// writeReceiptThrewtoFile(receiver)
			}
			// fmt.Println("receiver:", receiver)
			return common.Bytes2Hex((receiver).Return())
			// return (receiver).(c_network.Receipt1).Value
		case <-time.After(5 * time.Second):
			return "TimeOut"
		}
	}

}
func (client *Client) GetDataForCallSmartContract(call map[string]interface{}) ([]byte, error) {
	kq := client.EncodeAbi(call)
	//    callData := transaction.NewCallData(kq)
	callData := p_transaction.NewCallData(kq)

	return callData.Marshal()
}

// func (cli *Cli) GetAccountState(address string ,sign cm.Sign) (state.IAccountState, error) {
// 	parentConn := caller.client.connectionsManager.GetParentConnection()
// 	caller.client.messageSenderMap[address].SendBytes(parentConn, command.GetAccountState, common.FromHex(address), sign)

// 	select {
// 	case accountState := <-caller.client.accountStateChan:
// 		return accountState, nil
// 	case <-time.After(5 * time.Second):
// 		return nil, ErrorGetAccountStateTimedOut
// 	}

// }


func (client *Client) SendTransaction(
	fromAddress string,
	toAddress common.Address,
	amount *uint256.Int,
	action pb.ACTION,
	data []byte,
	relatedAddress [][]byte,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,
) (chan types.Receipt, error) {

	client.cliMap[fromAddress].mu.Lock()
	defer client.cliMap[fromAddress].mu.Unlock()
	// get account state
	parentConn := client.cliMap[fromAddress].clientContext.ConnectionsManager.ParentConnection()
	client.cliMap[fromAddress].clientContext.MessageSender.SendBytes(
		parentConn,
		command.GetAccountState,
		client.cliMap[fromAddress].clientContext.KeyPair.Address().Bytes(),
		p_common.Sign{},
	)

	as := <-client.cliMap[fromAddress].accountStateChan
	lastHash := as.LastHash()
	pendingBalance := as.PendingBalance()

	// bRelatedAddresses := make([][]byte, len(relatedAddress))
	// for i, v := range relatedAddress {
	// 	bRelatedAddresses[i] = v.Bytes()
	// }
	transaction, err := client.cliMap[fromAddress].transactionController.SendTransaction(
		lastHash,
		toAddress,
		pendingBalance,
		amount,
		maxGas,
		maxGasPrice,
		maxTimeUse,
		action,
		data,
		relatedAddress,
	)

	logger.Info("Sending transaction", transaction)
	if err != nil {
		return nil, err
	}

	// receipt := <-cli.receiptChan
	// fmt.Println("receipt:",receipt)
	return client.cliMap[fromAddress].receiptChan, nil
}

func (cli *Cli) AccountState(address common.Address) (types.AccountState, error) {
	cli.mu.Lock()
	defer cli.mu.Unlock()
	// get account state
	parentConn := cli.clientContext.ConnectionsManager.ParentConnection()
	cli.clientContext.MessageSender.SendBytes(
		parentConn,
		command.GetAccountState,
		address.Bytes(),
		p_common.Sign{},
	)
	as := <-cli.accountStateChan
	return as, nil
}

// func (client *Client) Subcribe(storageHost string, smartContractAddress common.Address) (chan interface{}, error) {
// 	fmt.Println("1111111111111")
// 	storageConnection := network.NewConnection(common.Address{}, p_common.STORAGE_CONNECTION_TYPE, storageHost)
// 	err := storageConnection.Connect()
// 	if err != nil {
// 		logger.Error("Unable to connect to storage", err)
// 		return nil, fmt.Errorf("unable to connect to storage")
// 	}

// 	go client.clientContext.SocketServer.HandleConnection(storageConnection)

// 	err = client.clientContext.MessageSender.SendBytes(storageConnection, command.SubscribeToAddress, smartContractAddress.Bytes(), p_common.Sign{})
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to send subscribe")
// 	}

// 	eventChan := make(chan interface{}, 1)
// 	// cli.clientContext.Handler.(*c_network.Handler).SetEventChan(eventChan)
// 	return eventChan, nil
// }

func (client *Client) EnterRelatedAddress(call map[string]interface{}) [][]byte {
	var arrmap []map[string]interface{}
	arr, _ := call["relatedAddresses"].([]interface{})
	if call["relatedAddresses"] == nil || len(arr) == 0 {
		fmt.Println("111111111")
		return defaultRelatedAddress
	} else {
		fmt.Println("2222222222")

		for _, v := range arr {
			arrmap = append(arrmap, v.(map[string]interface{}))
		}

		var relatedAddStr []string

		for _, v := range arrmap {
			relatedAddStr = append(relatedAddStr, v["address"].(string))
		}
		var relatedAddress [][]byte

		// temp := strings.Split(relatedAddStr, ",")
		logger.Info("Temp Related Address")
		for _, addr := range relatedAddStr {
			addressHex := common.HexToAddress(addr)
			logger.Info(addressHex)
			relatedAddress = append(relatedAddress, addressHex.Bytes())
		}
		defaultRelatedAddress = append(defaultRelatedAddress, relatedAddress...)
		return relatedAddress

	}
}

func (client *Client) EncodeAbi(call map[string]interface{}) []byte {
	var inputArray []interface{}
	if call["inputArray"] == nil {
		inputArray = []interface{}{}
	} else {
		inputArray, _ = call["inputArray"].([]interface{})
	}
	functionName, _ := call["function-name"].(string)

	// abiData, ok := call["abiData"].(string)
	// if !ok {
	// 	logger.Error(fmt.Sprintf("error when get abiData %"))
	// 	panic(fmt.Sprintf("error when get abiData "))
	// }
	// abiJson, err := JSON(strings.NewReader(abiData))
	// if err != nil {
	// 	panic(err)
	// }
	contract, _ := call["contract"].(string)
	if contract == "" {
		contract = ""
	}
	var path string

	path = "./abi/" + contract + ".json"
	var out []byte

	if contract != "" {
		reader, err := os.Open(path)
		if err != nil {
			log.Fatalf("Error occured while reading contract %s", contract)
		}
		abiJson, err := JSON(reader)
		if err != nil {
			panic(err)
		}
		var abiTypes []interface{}
		for _, item := range inputArray {
			itemArr := encodeAbiItem(item)
			fmt.Println("itemArr:", itemArr)
			for _, v := range itemArr {
				abiTypes = append(abiTypes, v)
				fmt.Println("abiTypes:", abiTypes)
			}
		}
		out, err = abiJson.Pack(functionName, abiTypes[:]...)

		if err != nil {
			panic(err)
		}
		fmt.Println("out:", hex.EncodeToString(out))

	}
	return out
}

func encodeAbiItem(item interface{}) []interface{} {
	var result []interface{}
	var itemMap map[string]interface{}
	if err := json.Unmarshal([]byte(item.(string)), &itemMap); err != nil {
		log.Fatal(err)
	}
	itemType, _ := itemMap["type"].(string)
	switch itemType {
	case "tuple":
		var value []interface{}
		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", itemMap["value"])), &value); err != nil {
			log.Fatal(err)
		}

		var components []interface{}
		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", itemMap["components"])), &components); err != nil {
			log.Fatal(err)
		}

		var abiTypes []interface{}
		for i, component := range components {
			componentBytes, _ := json.Marshal(component)
			componentType, _ := component.(map[string]interface{})["type"].(string)
			if componentType == "tuple" || componentType == "tuple[]" {
				components[i].(map[string]interface{})["value"] = value[i]
				abiTypes = append(abiTypes, encodeAbiItem(componentBytes))
			} else {
				abiTypes = append(abiTypes, getAbiType(componentType, value[i]))
			}
		}
		result = abiTypes
	case "tuple[]":
		var value []interface{}
		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", itemMap["value"])), &value); err != nil {
			log.Fatal(err)
		}

		var components []interface{}
		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", itemMap["components"])), &components); err != nil {
			log.Fatal(err)
		}

		var tuples []interface{}
		for _, v := range value {
			vArray := v.([]interface{})
			var abiTypes []interface{}
			for j, component := range components {
				componentBytes, _ := json.Marshal(component)
				componentType, _ := component.(map[string]interface{})["type"].(string)
				components[j].(map[string]interface{})["value"] = vArray[j]
				if componentType == "tuple" || componentType == "tuple[]" {
					abiTypes = append(abiTypes, encodeAbiItem(componentBytes))
				} else {
					abiTypes = append(abiTypes, getAbiType(componentType, vArray[j]))
				}
			}
			tuples = append(tuples, abiTypes...)
		}
		result = tuples
	default:
		value := itemMap["value"]

		var arr []interface{}

		result1 := getAbiType(itemType, value)
		result = append(arr, result1)

	}
	return result
}
func getAbiType(dataType string, data interface{}) interface{} {
	if strings.Contains(dataType, "int") {
		params := big.NewInt(0)
		params, ok := params.SetString(fmt.Sprintf("%v", int64(data.(float64))), 10)

		if !ok {
			log.Warn("Format big int: error")
			return nil
		}
		return params

	} else {
		switch dataType {
		case "string":
			return data.(string)
		case "bool":
			return data.(bool)
		case "address":
			return common.HexToAddress(data.(string))
		case "uint8":
			intVar, err := strconv.Atoi(data.(string))
			if err != nil {
				log.Warn("Conver Uint8 fail", err)
				return nil
			}
			return uint8(intVar)
		// case "uint", "uint256":
		// 	nubmer := big.NewInt(0)
		// 	nubmer, ok := nubmer.SetString(data.(string), 10)
		// 	if !ok {
		// 		log.Warn("Format big int: error")
		// 		return nil
		// 	}
		// 	return nubmer
		case "array", "slice", "bytes32":
			fmt.Println("array nè")
			rv := reflect.ValueOf(data)
			var out []interface{}
			for i := 0; i < rv.Len(); i++ {
				out = append(out, rv.Index(i).Interface())
			}

			return out
		default:
			return data
		}
	}
}
