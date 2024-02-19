package core

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	// . "github.com/ethereum/go-ethereum/accounts/abi"
)

type Account struct {
	Address string
	Private string
}

var PORT int
var Contracts = [...]Contract{
	{Name: "kventure", Address: "0xf28B2FDA437D4717915FfDA6b858B6A29fBAE74f"},
	// {Name: "codepool", Address: "0x9eeB553fC6d5D6f9b35a2728121317DCc2eA6C38"},
	// {Name: "product", Address: "0x79965F966Cdb14130127EEB9E8411869893ba26A"},


}

// var accounts = [...]Account{
// 	{
// 		Address: "45c75cfb8e20a8631c134555fa5d61fcf3e602f2",
// 		Private: "36e1aa979f98c7154fb2491491ec044ccac099651209ccfbe2561746dbe29ebb",
// 	},

// }
//f8eaba3eb679f6defbe78ce8dd5229ec3622f2a7
func GetPORT() int {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	// log.Info("PORT: ", os.Getenv("PORT"))
	PORT, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
	return PORT
}

type Contract struct {
	Name    string
	Address string
}
