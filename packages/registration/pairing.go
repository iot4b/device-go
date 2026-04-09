package registration

import (
	"device-go/packages/crypto"
	"device-go/packages/dsm"
	"device-go/packages/storage"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

const (
	pairingRequestTries      = 5
	pairingCodeLength        = 5
	pairingStatusPending     = "pending"
	pairingStatusBound       = "bound"
	pairingHeartbeatInterval = 2 * time.Second
	pairingRequestTimeout    = 5 * time.Second
)

type pairingRegisterRequest struct {
	Name        string `json:"nm,omitempty"`
	Type        string `json:"t,omitempty"`
	Version     string `json:"ver,omitempty"`
	PublicSign  string `json:"k"`
	PublicNacl  string `json:"n"`
	PairingCode string `json:"code"`
}

type pairingResultResponse struct {
	Status          string `json:"s"`
	ContractAddress string `json:"a,omitempty"`
}

func Pair() error {
	if storage.HasContractAddress() {
		return nil
	}

	if !storage.HasPairingCode() {
		return startPairing()
	}

	result, err := pairingResult(storage.Pairing.NodeIpPort, crypto.Keys.PublicSign)
	if err != nil {
		if isPairingRecordMissing(err) {
			nodeHost := storage.Pairing.NodeIpPort
			nodeAddress := storage.Pairing.NodeAddress
			log.Infof("pairing record missing on node %s, rotating code", nodeHost)
			if clearErr := storage.ClearPairing(); clearErr != nil {
				return clearErr
			}
			if nodeHost != "" {
				if retryErr := startPairingOnNode(nodeHost, nodeAddress); retryErr == nil {
					log.Infof("pairing code rotated on current node %s", nodeHost)
					return nil
				} else {
					log.Errorf("failed to rotate pairing code on current node %s: %v", nodeHost, retryErr)
				}
			}
			log.Info("selecting a new node for pairing")
			return startPairing()
		}
		log.Debugf("pairing result request failed: %v", err)
		return nil
	}

	if result.Status != pairingStatusBound || result.ContractAddress == "" {
		return nil
	}

	return applyPairingBound(result.ContractAddress)
}

func PairingHeartbeat() {
	ticker := time.NewTicker(pairingHeartbeatInterval)
	defer ticker.Stop()

	for range ticker.C {
		if storage.Device.Address != "" {
			return
		}
		if err := storage.PairingTouch(); err != nil {
			log.Error(err)
		}
	}
}

func startPairing() error {
	nodeHost, nodeAddress, err := selectFastestNode()
	if err != nil {
		return err
	}

	return startPairingOnNode(nodeHost, nodeAddress)
}

func startPairingOnNode(nodeHost, nodeAddress string) error {
	for attempt := 0; attempt < pairingRequestTries; attempt++ {
		code := generatePairingCode()
		result, err := pairingRegister(nodeHost, pairingRegisterRequest{
			Name:        storage.Device.Name,
			Type:        storage.Device.Type,
			Version:     storage.Device.Version,
			PublicSign:  crypto.Keys.PublicSign,
			PublicNacl:  crypto.Keys.PublicNacl,
			PairingCode: code,
		})
		if err == nil {
			if result.Status == pairingStatusBound && result.ContractAddress != "" {
				return applyPairingBound(result.ContractAddress)
			}
			return storage.SetPairing(code, nodeHost, nodeAddress, pairingStatusPending)
		}

		if !isPairingConflict(err) {
			return err
		}
	}

	return fmt.Errorf("failed to allocate a unique pairing code after %d attempts", pairingRequestTries)
}

func generatePairingCode() string {
	upperBound := 1
	for i := 0; i < pairingCodeLength; i++ {
		upperBound *= 10
	}
	return fmt.Sprintf("%0*d", pairingCodeLength, rand.Intn(upperBound))
}

func pairingRegister(nodeHost string, payload pairingRegisterRequest) (pairingResultResponse, error) {
	body, err := doNodePairingRequest(nodeHost, coalago.POST, "/pairing/register", nil, payload)
	if err != nil {
		return pairingResultResponse{}, err
	}

	var result pairingResultResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return pairingResultResponse{}, err
	}

	return result, nil
}

func pairingResult(nodeHost, publicSign string) (pairingResultResponse, error) {
	body, err := doNodePairingRequest(nodeHost, coalago.GET, "/pairing/result", map[string]string{"k": publicSign}, nil)
	if err != nil {
		return pairingResultResponse{}, err
	}

	var result pairingResultResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return pairingResultResponse{}, err
	}

	return result, nil
}

func applyPairingBound(contractAddress string) error {
	storage.Device.Address = dsm.EverAddress(contractAddress)
	storage.Device.NodeIpPort = ""
	if err := storage.Save(); err != nil {
		return err
	}
	return storage.ClearPairing()
}

func doNodePairingRequest(nodeHost string, method coalago.CoapCode, path string, query map[string]string, payload any) ([]byte, error) {
	client := coalago.NewClient()
	msg := coalago.NewCoAPMessage(coalago.CON, method)
	msg.SetURIPath(path)
	msg.Timeout = pairingRequestTimeout

	for key, value := range query {
		msg.SetURIQuery(key, value)
	}

	if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		msg.SetStringPayload(string(body))
	}

	resp, err := client.Send(msg, nodeHost)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("nil response")
	}
	if resp.Code != coalago.CoapCodeContent {
		message := strings.TrimSpace(string(resp.Body))
		if message == "" {
			message = resp.Code.String()
		}
		return nil, fmt.Errorf("%s (%s)", message, resp.Code.String())
	}

	return resp.Body, nil
}

func isPairingConflict(err error) bool {
	if err == nil {
		return false
	}

	value := strings.ToLower(err.Error())
	return strings.Contains(value, "pairing code already in use") ||
		strings.Contains(value, "409")
}

func isPairingRecordMissing(err error) bool {
	if err == nil {
		return false
	}

	value := strings.ToLower(err.Error())
	return strings.Contains(value, "pairing record not found") ||
		strings.Contains(value, "not found")
}
