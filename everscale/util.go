package everscale

import (
	"github.com/markgenuine/ever-client-go/domain"
)

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
