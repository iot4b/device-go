package everscale

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/markgenuine/ever-client-go/domain"
	"github.com/pkg/errors"
)

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "os.Open(%s)", path)
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrapf(err, "io.ReadAll(%s)", path)
	}
	return bytes, nil
}

func getAbiFromFile(path string) (*domain.Abi, error) {
	byteAbi, err := readFile(path)
	ac := &domain.AbiContract{}
	err = json.Unmarshal(byteAbi, &ac)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal(byteAbi, &ac)")
	}
	return domain.NewAbiContract(ac), nil
}

func processMessage(abi *domain.Abi, address, method string, input interface{}, signer *domain.Signer) (*domain.ResultOfProcessMessage, error) {
	return ever.Processing.ProcessMessage(&domain.ParamsOfProcessMessage{
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

func ReadContract(path, name string) (abi *domain.Abi, tvc []byte, err error) {
	abi, err = getAbiFromFile(fmt.Sprintf("%s/%s.abi.json", path, name))
	if err != nil {
		err = errors.Wrapf(err, "getAbiFromFile(%s)", name+".abi.json")
		return
	}

	tvc, err = readFile(fmt.Sprintf("%s/%s.tvc", path, name))
	if err != nil {
		err = errors.Wrapf(err, "readFile(%s)", name+".tvc")
	}
	return
}
