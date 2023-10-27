package utils

import (
	"fmt"
	"github.com/markgenuine/ever-client-go/domain"
	"github.com/pkg/errors"
)

const scPath = "../smartcontracts/"

func ReadContract() (abi *domain.Abi, tvc []byte, err error) {
	abi, err = GetAbi("device")
	if err != nil {
		err = errors.Wrap(err, "GetAbi")
		return
	}

	tvc, err = ReadFile(scPath + "device.tvc")
	if err != nil {
		err = errors.Wrapf(err, "readFile(%s)", scPath+"device.tvc")
	}
	return
}

func GetAbi(cType string) (*domain.Abi, error) {
	path := scPath + fmt.Sprintf("_%s/%s.abi.json", cType, cType)
	ac := &domain.AbiContract{}
	err := ReadJSONFile(path, ac)
	if err != nil {
		return nil, errors.Wrapf(err, "ReadJSONFile(%s)", path)
	}
	return domain.NewAbiContract(ac), nil
}
