package everscale

import (
	"encoding/base64"
	"encoding/json"
	"github.com/markgenuine/ever-client-go/domain"
	log "github.com/ndmsystems/golog"
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
	message, err := Ever.Abi.EncodeMessage(cd.deployOptions)
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
	resp, err := Ever.Processing.ProcessMessage(params, nil)
	if err != nil {
		return err
	}
	log.Debug(resp)
	return nil
}
