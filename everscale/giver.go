package everscale

import (
	"github.com/pkg/errors"
)

// GetTokensFromGiver transfer a [value] of test nanotokens from giverAddress to [account]
func GetTokensFromGiver(g Giver, giverAbiFile, account string, value int) (err error) {
	signer := NewSigner(g.Public, g.Secret)

	abi, err := getAbiFromFile(giverAbiFile)
	if err != nil {
		return errors.Wrapf(err, "getAbiFromFile(%s)", giverAbiFile)
	}

	input := sendTransaction{
		Dest:   account,
		Value:  value,
		Bounce: false,
	}
	_, err = processMessage(
		abi,
		g.Address,
		"sendTransaction",
		input, signer)
	return
}
