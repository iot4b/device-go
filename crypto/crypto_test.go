package crypto

import (
	"device-go/everscale"
	"testing"
)

const (
	everEndpoint = "https://devnet.evercloud.dev/5c41d775a6ab4bacb3cc25666b93de60/graphql"
	keysPath     = "../device.keys.json"
)

func Test(t *testing.T) {
	everscale.Init([]string{everEndpoint})
	defer everscale.Destroy()

	Init(keysPath)

	original := "unsigned"

	signed := Keys.Sign(original)

	unsigned, valid := Keys.Verify(signed)

	if unsigned != original {
		t.Errorf("have: %s, want: %s", unsigned, original)
	}

	if !valid {
		t.Errorf("signature is not valid")
	}
}
