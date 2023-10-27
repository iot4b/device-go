package everscale

import "device-go/utils"

type Giver struct {
	Address string
	Public  string
	Secret  string
}

// request for sendTransaction method of giver contract
type sendTransaction struct {
	Dest    string `json:"dest"`              // dest address
	Value   int    `json:"value"`             // amount in nano EVER
	Bounce  bool   `json:"bounce"`            // false for contract deploy
	Flags   int    `json:"flags,omitempty"`   // ???
	Payload string `json:"payload,omitempty"` // ???
}

func (g *Giver) SendTokens(address string, amount int) error {
	signer := NewSigner(g.Public, g.Secret)

	abi, err := utils.GetAbi("giver")
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
