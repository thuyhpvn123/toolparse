package core

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	. "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type ContractABI struct {
	Name    string
	Address string
	Abi     ABI
}
type Event2 struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func (contract *ContractABI) InitContract(info Contract) {
	reader, err := os.Open("./abi/" + info.Name + ".json")
	if err != nil {
		log.Fatalf("Error occured while reading %s", "./abi/"+info.Name+".json")
	}
	contract.Abi, err = JSON(reader)
	if err != nil {
		log.Fatalf("Error occured while init abi %s", info.Name)
	}
	contract.Address = info.Address
	contract.Name = info.Name
	fmt.Println("Init contract ", info.Name)
}
func (contract *ContractABI) Decode(name, data string) interface{} {
	bytes, err := hex.DecodeString(data)
	if err != nil {
		log.Fatalf("Error occured while convert data to byte[] - Data: %s", data)
	}
	result := make(map[string]interface{})
	err = contract.Abi.UnpackIntoMap(result, name, bytes)
	if err != nil {
		log.Fatalf("Error occured while unpack %s - %s \n %s \n %s", name, err, data, bytes)
	}
	return result
}

func (contract *ContractABI) Encode(name string, args ...interface{}) []byte {
	formatedData := contract.formatPreEncode(contract.Abi.Methods[name].Inputs, args)
	data, err := contract.Abi.Pack(name, formatedData[:]...)
	if err != nil {
		log.Fatalf("Error occured while pack %s - %s", name, err)
	}
	return data
}

func (contract *ContractABI) formatPreEncode(args Arguments, data []interface{}) []interface{} {
	i := 0
	temp := make([]interface{}, len(args))
	for _, arg := range args {
		temp[i] = formatData(arg.Type.String(), data[i])
		i++
	}
	return temp
}

func (contract *ContractABI) DecodeEvent(name string, data string, topics map[int]interface{}) map[string]interface{} {
	result := contract.Decode(name, data).(map[string]interface{})
	contract.formatTopics(name, result, topics)
	return result
}
func (contract *ContractABI) formatTopics(name string, result map[string]interface{}, topics map[int]interface{}) {
	i := 1

	for _, arg := range contract.Abi.Events[name].Inputs {
		if arg.Indexed {
			bytes := common.FromHex(topics[i].(string))
			switch arg.Type.T {
			case IntTy, UintTy:
				result[arg.Name] = ReadInteger(arg.Type, bytes)
			case BoolTy:
				result[arg.Name] = readBool(bytes)
			case AddressTy:
				result[arg.Name] = common.BytesToAddress(bytes)
			}
			i++
		}
	}
}
func readBool(word []byte) bool {
	for _, b := range word[:31] {
		if b != 0 {
			return false
		}
	}
	switch word[31] {
	case 0:
		return false
	case 1:
		return true
	default:
		return false
	}
}


// format utils
func formatData(dataType string, data interface{}) interface{} {
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
	case "address[]":
		var addressList []common.Address
		for _, item := range data.([]string) {
			addressList = append(addressList, common.HexToAddress(item))
		}
		return addressList
	case "string[]":
		var list []string
		for _, item := range data.([]string) {
			list = append(list, formatData("string", item).(string))
			

		}
		return list
		// var out []string
		// for i := 0; i < len(data.([]interface{})); i++ {
		// 	out = append(out, data.([]interface{})[i].(string))
		// }

		// return out

	case "uint256[]":
		var list []interface{}
		for _, item := range data.([]string) {
			list = append(list, formatData("uint256", item))
		}
		return list
	case "uint", "uint256":
		nubmer := big.NewInt(0)
		nubmer, ok := nubmer.SetString(data.(string), 10)
		if !ok {
			log.Warn("Format big int: error")
			return nil
		}
		return nubmer
	default:
		return nil
	}
}

// func (contract *ContractABI) formatToNumber(data interface{}) string {
// 	number := data.(*big.Int).String()
// 	return number
// }

func StringToBytes32(input string) [32]byte {
	var result [32]byte

	// Convert the string to a byte slice
	byteSlice := []byte(input)

	// Truncate or pad the byte slice to 32 bytes
	copy(result[:], byteSlice)

	return result
}