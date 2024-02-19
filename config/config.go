package config

import (
	// "encoding/hex"
	"encoding/json"
	"os"
	"github.com/meta-node-blockchain/meta-node/pkg/bls"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	p_common "github.com/meta-node-blockchain/meta-node/pkg/common"
	// "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/config"
	"github.com/meta-node-blockchain/meta-node/types"
)
// const CONFIG_FILE_PATH = "config.json"

type Connection struct {
	Address string `json:address`
	Ip      string `json:ip`
	Port    int    `json:port`
	Type    string `json:type`
}

type Config struct {
	// Address                    string       `json:"address"`
	// ByteAddress                []byte       `json:"-"`
	// BytePrivateKey             []byte       `json:"-"`
	// Ip                         string       `json:"ip"`
	// Port                       int          `json:"port"`
	// NodeType_                 string       `json:"node_type"`
	// HashPerSecond              int          `json:"hash_per_second"`
	// TickPerSecond              int          `json:"tick_per_second"`
	// TickPerSlot                int          `json:"tick_per_slot"`
	// BlockStackSize             int          `json:"block_stack_size"`
	// TimeOutTicks               int          `json:"time_out_ticks"` // how many tick validator should wait before create virture block
	// TransactionPerHash         int          `json:"transaction_per_hash"`
	// NumberOfValidatePohRoutine int          `json:"number_of_validate_poh_routine"`
	// AccountDBPath              string       `json:"account_db_path"`
	// SecretKey                  string       `json:"secret_key"`
	// TransferFee                int          `json:"transfer_fee"`
	// GuaranteeAmount            int          `json:"guarantee_amount"`
	// TransactionFee             *uint256.Int `json:"-"`
	// TransactionFeeHex          string       `json:"transaction_fee"`

	// Version_          string     `json:"version"`
	// BytePublicKey    []byte     `json:"-"`
	// ParentConnection Connection `json:"parent_connection"`
	ServerAddress_    string     `json:"server_address"`
	// PrivateKey_       string     `json:"private_key"`

	// TcpIp   string `json:"tcp_ip"`
	// TcpPort int    `json:"tcp_port"`

	// ParentAddress           string `json:"parent_address"`
	// ParentConnectionAddress string `json:"parent_connection_address"`
	// ParentConnectionType    string `json:"parent_connection_type"`

	// ConnectionAddress_       string `json:"connection_address"`
	// PublicConnectionAddress_ string `json:"public_connection_address"`
	PrivateKey_ string `json:"private_key"`

	ConnectionAddress_       string `json:"connection_address"`
	PublicConnectionAddress_ string `json:"public_connection_address"`
	DnsLink_                 string `json:"dns_link"`

	Version_          string       `json:"version"`
	TransactionFeeHex string       `json:"transaction_fee"`
	TransactionFee    *uint256.Int `json:"-"`

	ParentAddress           string `json:"parent_address"`
	ParentConnectionAddress string `json:"parent_connection_address"`
	ParentConnectionType    string `json:"parent_connection_type"`
}

const (
	CONFIG_FILE_PATH = "config/conf.json"
)

// func LoadConfig(configPath string) (types.Config, error) {
// 	var config Config
// 	// raw, err := ioutil.ReadFile("config/conf.json")
// 	raw, err := os.ReadFile(configPath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = json.Unmarshal(raw, &config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	byteAddress, err := hex.DecodeString(config.Address)
// 	if err != nil {
// 		return nil, err
// 	}
// 	config.ByteAddress = byteAddress
// 	// log.Printf("Config loaded: %v\n", config)
// 	// config.BytePrivateKey, config.BytePublicKey, config.ByteAddress = ccrypto.GenerateKeyPairFromSecretKey(config.SecretKey)
// 	config.TransactionFee = uint256.NewInt(0).SetBytes(common.FromHex(config.TransactionFeeHex))

// 	return config, nil
// }
func LoadConfig(configPath string) (types.Config, error) {
	// general config
	var config Config
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(raw, config)
	if err != nil {
		return nil, err
	}
	config.TransactionFee = uint256.NewInt(0).SetBytes(common.FromHex(config.TransactionFeeHex))
	return config, nil
}


// var AppConfig,err = loadConfig(CONFIG_FILE_PATH)

// func (config Config) Version() string {
// 	return config.Version_
// }

// func (config Config) Pubkey() []byte {
// 	return config.BytePublicKey
// }

// func (config Config) GetPrivateKey() []byte {
// 	return config.BytePrivateKey
// }

// func (config Config) PrivateKey() []byte {
// 	return common.FromHex(config.PrivateKey_)
// }
// func (config Config) NodeType() string {
// 	return p_common.CLIENT_CONNECTION_TYPE
// }
// func (config Config) ConnectionAddress() string {
// 	return config.ConnectionAddress_
// }

// func (config Config) PublicConnectionAddress() string {
// 	return config.PublicConnectionAddress_
// }

func (c Config) ConnectionAddress() string {
	return c.ConnectionAddress_
}

func (c Config) PublicConnectionAddress() string {
	return c.PublicConnectionAddress_
}

func (c Config) Version() string {
	return c.Version_
}

func (c Config) PrivateKey() []byte {
	return common.FromHex(c.PrivateKey_)
}

func (c Config) Address() common.Address {
	_, _, address := bls.GenerateKeyPairFromSecretKey(c.PrivateKey_)
	return address
}

func (c Config) NodeType() string {
	return p_common.CLIENT_CONNECTION_TYPE
}

func (c Config) DnsLink() string {
	return c.DnsLink_
}
func (c Config) ServerAddress() string {
	return c.ServerAddress_
}

