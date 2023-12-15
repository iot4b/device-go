package everscale

import (
	"fmt"
	"github.com/markgenuine/ever-client-go/domain"
	"github.com/pkg/errors"
)

// Execute a [method] on a contract [name] deployed to [address]
func Execute(name, address, method string, input interface{}, signer *domain.Signer) ([]byte, error) {
	fmt.Println("executing", method, "on", name, "contract at address", address)

	abi, err := getAbi(name)
	if err != nil {
		return nil, errors.Wrap(err, "getAbi")
	}

	result, err := processMessage(abi, address, method, input, signer)
	if err != nil {
		return nil, errors.Wrap(err, "processMessage")
	}

	fmt.Println(string(result.Decoded.Output))
	return result.Decoded.Output, nil
}
