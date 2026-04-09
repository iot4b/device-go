package registration

import (
	"device-go/packages/storage"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	log "github.com/ndmsystems/golog"
)

const (
	serviceName         = "iot4b"
	serviceHeartbeatTTL = 10 * time.Second
)

func Setup() {
	if storage.Device.Address != "" {
		fmt.Println("Device contract address is already configured:")
		fmt.Println(storage.Device.Address)
		return
	}

	if !isServiceRunning() {
		fmt.Printf("%s service is not running.\n", serviceName)
		fmt.Printf("Start the service first so it can create a pairing session.\n")
		if runtime.GOOS == "darwin" {
			fmt.Printf("brew services start %s\n", serviceName)
		} else if runtime.GOOS == "linux" {
			fmt.Printf("systemctl start %s\n", serviceName)
		} else {
			fmt.Printf("You can also run it locally, for example: go run .\n")
		}
		return
	}

	fmt.Println("Waiting for pairing session...")
	lastCode := ""
	renderedWaiting := false

	for {
		if err := storage.Update(); err != nil {
			log.Error(err)
		}
		if storage.Device.Address != "" && storage.Device.NodeIpPort != "" {
			fmt.Println("Device pairing completed.")
			fmt.Printf("Address: %s\n", storage.Device.Address)
			fmt.Printf("Node: %s\n", storage.Device.Node)
			return
		}
		if storage.Pairing.Code != "" && storage.Pairing.Code != lastCode {
			lastCode = storage.Pairing.Code
			renderPairingSetup(storage.Pairing.Code)
			renderedWaiting = false
		} else if storage.Pairing.Code == "" && lastCode != "" {
			lastCode = ""
			renderPairingSetup("")
			renderedWaiting = true
		} else if storage.Pairing.Code == "" && !renderedWaiting {
			renderPairingSetup("")
			renderedWaiting = true
		}
		if !isServiceRunning() {
			fmt.Printf("%s service stopped before pairing completed.\n", serviceName)
			return
		}
		time.Sleep(time.Second)
	}
}

func renderPairingSetup(code string) {
	clearConsole()
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("Setup Device")
	fmt.Println()
	var displayCode string
	if code == "" {
		displayCode = "....."
	} else {
		displayCode = code
	}
	fmt.Println("+---------------+")
	fmt.Printf("|     %s     |\n", displayCode)
	fmt.Println("+---------------+")
	fmt.Println()
	fmt.Println("Enter this code in the app to bind the device.")
	fmt.Println("The code refreshes about once per minute.")
	fmt.Println("========================================")
}

func clearConsole() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err == nil {
		return
	}
	fmt.Print("\033[H\033[2J")
}

func isServiceRunning() bool {
	if storage.Pairing.HeartbeatAt == 0 {
		return false
	}
	lastSeen := time.Unix(storage.Pairing.HeartbeatAt, 0)
	return time.Since(lastSeen) <= serviceHeartbeatTTL
}
