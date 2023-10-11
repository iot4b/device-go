package everscale

import (
	"device-go/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/markgenuine/ever-client-go/domain"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
)

type ContractBuilder struct {
	Public      string
	Secret      string
	Abi         *domain.Abi
	Tvc         []byte
	InitialData interface{}

	address       string
	signer        *domain.Signer
	deployOptions *domain.ParamsOfEncodeMessage
}

func (cd *ContractBuilder) InitDeployOptions() *ContractBuilder {
	initialData := json.RawMessage(`{}`)
	if cd.InitialData != nil {
		data, err := json.Marshal(cd.InitialData)
		if err == nil {
			initialData = data
		}
	}
	cd.signer = NewSigner(cd.Public, cd.Secret)
	cd.deployOptions = &domain.ParamsOfEncodeMessage{
		Abi:    cd.Abi,
		Signer: cd.signer,
		DeploySet: &domain.DeploySet{
			Tvc:         base64.StdEncoding.EncodeToString(cd.Tvc),
			InitialData: initialData,
		},
	}
	cd.address = cd.CalcWalletAddress()
	return cd
}

func (cd *ContractBuilder) CalcWalletAddress() string {
	message, err := ever.Abi.EncodeMessage(cd.deployOptions)
	if err != nil {
		log.Error(err)
		return ""
	}
	return message.Address
}

func (cd *ContractBuilder) Deploy(input interface{}) error {
	log.Debug(input)
	deployOptions := *cd.deployOptions
	deployOptions.CallSet = &domain.CallSet{
		FunctionName: "constructor",
		Input:        input,
	}
	params := &domain.ParamsOfProcessMessage{
		MessageEncodeParams: &deployOptions,
		SendEvents:          false,
	}
	resp, err := ever.Processing.ProcessMessage(params, nil)
	if err != nil {
		return err
	}
	log.Debug(resp)
	return nil
}

type Giver struct {
	Address string
	Public  string
	Secret  string
}

func (g *Giver) SendTokens(giverAniFile, address string, amount int) error {
	signer := NewSigner(g.Public, g.Secret)

	abi, err := getAbiFromFile(giverAniFile)
	if err != nil {
		return err
	}

	input := sendTransaction{
		Dest:   address,
		Value:  amount,
		Bounce: false,
	}
	_, err = processMessage(abi,
		g.Address,
		"sendTransaction",
		input,
		signer)
	return err
}

// Destroy client when finished
func Destroy() {
	ever.Client.Destroy()
}

func KeysFromFile() (public, secret string) {
	keysPath := "keys.device.json"
	keyPair, err := utils.ReadKeysFile(keysPath)
	if err != nil {
		keyPair, _ = GenerateKeyPair()
		data, _ := json.Marshal(keyPair)
		utils.SaveFile(keysPath, data)
	}
	public, secret = keyPair.Public, keyPair.Secret
	log.Debugf("keys from file. public: %s secret: %s", public, secret)
	return
}

func GenKeys() (public, secret string) {
	keyPair, _ := GenerateKeyPair()
	public, secret = keyPair.Public, keyPair.Secret

	log.Debugf("generate new keys. public: %s secret: %s", public, secret)
	return
}

// Execute a [method] on a contract [name] deployed to [address]
func Execute(name, address, method string, input interface{}) ([]byte, error) {
	fmt.Println("executing", method, "on", name, "contract at address", address)

	abiFile := fmt.Sprintf("contracts/%s.abi.json", name)
	abi, err := getAbiFromFile(abiFile)
	if err != nil {
		return nil, errors.Wrapf(err, "getAbiFromFile(%s)", abiFile)
	}

	result, err := processMessage(abi, address, method, input, domain.NewSigner(domain.SignerNone{}))
	if err != nil {
		return nil, errors.Wrap(err, "processMessage")
	}

	fmt.Println(string(result.Decoded.Output))
	return result.Decoded.Output, nil
}

// GetAccountInfo get balance and boc of the [address] account
func GetAccountInfo(address string) (AccountInfo, error) {
	res, err := ever.Net.Query(&domain.ParamsOfQuery{
		Query: fmt.Sprintf(`
		query {
		  blockchain {
			account(
			  address: "%s"
			) {
			   info {
				balance(format: DEC)
				boc
			  }
			}
		  }
		}`, address),
	})
	if err != nil {
		return AccountInfo{}, errors.Wrap(err, "ever.Net.Query")
	}

	result := &queryResult{}
	err = json.Unmarshal(res.Result, result)
	if err != nil {
		return AccountInfo{}, errors.Wrap(err, "json.Unmarshal")
	}

	info := result.Data.Blockchain.Account.Info
	info.Address = address

	return info, nil
}

func GenerateKeyPair() (domain.KeyPair, error) {
	keys, err := ever.Crypto.GenerateRandomSignKeys()
	return *keys, err
}
