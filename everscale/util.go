package everscale

import (
	"device-go/dsm"
	"device-go/utils"
	"fmt"
	"github.com/markgenuine/ever-client-go/domain"
	"github.com/pkg/errors"
)

const scPath = "../smartcontracts/build/"

func ReadContract() (abi *domain.Abi, tvc []byte, err error) {
	abi, err = getAbi("Device")
	if err != nil {
		err = errors.Wrap(err, "getAbi")
		return
	}

	tvc, err = utils.ReadFile(scPath + "Device.tvc")
	if err != nil {
		err = errors.Wrapf(err, "readFile(%s)", scPath+"Device.tvc")
	}
	return
}

func getAbi(name string) (*domain.Abi, error) {
	path := scPath + name + ".abi.json"
	abi := &domain.AbiContract{}
	err := utils.ReadJSONFile(path, abi)
	if err != nil {
		return nil, errors.Wrapf(err, "ReadJSONFile(%s)", path)
	}
	return domain.NewAbiContract(abi), nil
}

func processMessage(abi *domain.Abi, address, method string, input interface{}, signer *domain.Signer) (*domain.ResultOfProcessMessage, error) {
	return Ever.Processing.ProcessMessage(&domain.ParamsOfProcessMessage{
		MessageEncodeParams: &domain.ParamsOfEncodeMessage{
			Address: address,
			Abi:     abi,
			CallSet: &domain.CallSet{
				FunctionName: method,
				Input:        input,
			},
			Signer: signer,
		},
		SendEvents: false,
	}, nil)
}

func NewSigner(public, secret string) *domain.Signer {
	return domain.NewSigner(domain.SignerKeys{Keys: &domain.KeyPair{
		Public: public,
		Secret: secret,
	}})
}

// execute a [method] on a contract [name] deployed to [address]
func execute(name string, address dsm.EverAddress, method string, input interface{}, signer *domain.Signer) ([]byte, error) {
	fmt.Println("executing", method, "on", name, "contract at address", address)

	abi, err := getAbi(name)
	if err != nil {
		return nil, errors.Wrap(err, "getAbi")
	}

	result, err := processMessage(abi, string(address), method, input, signer)
	if err != nil {
		return nil, errors.Wrap(err, "processMessage")
	}

	fmt.Println(string(result.Decoded.Output))
	return result.Decoded.Output, nil
}
