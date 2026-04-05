package main

import (
	"device-go/packages/aliver"
	"device-go/packages/api"
	"device-go/packages/buildinfo"
	"device-go/packages/config"
	"device-go/packages/crypto"
	"device-go/packages/events"
	"device-go/packages/handlers"
	"device-go/packages/registration"
	"device-go/packages/storage"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
	"github.com/spf13/cobra"
)

const serviceName = "iot4b"

var env string  // environment (config name)
var port string // coala port
var showVersion bool

func main() {
	var rootCmd = &cobra.Command{
		Use:   "iot4b",
		Short: "Device CLI",
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Println(buildinfo.Summary())
				return
			}
			runDevice(cmd, args)
		},
	}
	rootCmd.PersistentFlags().StringVar(&env, "env", "iot4b", "Set environment")
	rootCmd.PersistentFlags().StringVar(&port, "port", "5684", "Set coala port")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Print version information")
	rootCmd.ParseFlags(os.Args[1:])
	if showVersion || isVersionCommand(os.Args[1:]) {
		fmt.Println(buildinfo.Summary())
		return
	}

	config.Init(env)
	log.Init(config.Bool("debug"))

	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup the device",
		Run:   deviceSetup,
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show device status",
		Run:   deviceStatus,
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show build version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(buildinfo.Summary())
		},
	}

	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(versionCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runDevice(_ *cobra.Command, _ []string) {
	crypto.Init(config.Get("localFiles.keys"))

	log.Info("Init Local Storage")
	storage.Init(
		config.Get("localFiles.contract"),
		config.Get("localFiles.setup"),
		config.Get("everscale.elector"),
		config.Get("everscale.vendor"),
		config.Get("everscale.deviceAPI"),
		config.Get("info.type"),
		buildinfo.Version)

	waitForContractAddress()

	// сервер для запросов от клиентов и нод
	server := coalago.NewServer()
	server.GET("/info", handlers.GetInfo)
	server.POST("/cmd", handlers.ExecCmd)
	server.POST("/update", handlers.Update)
	server.POST("/sign", handlers.Sign)

	go listen(server)
	go sigterm()

	for {
		//start := storage.Device.NodeIpPort == ""
		// регистрируем устройство на ноде
		// если ошибка, то повторяем цикл регистрации
		err := registration.Register()
		if err == nil {
			//go sendEvents(start)

			// начинаем слать alive пакеты, чтобы сохранять соединение для udp punching
			aliver.Run(server, config.Time("timeout.alive"))
		} else {
			log.Error(err)
			time.Sleep(5 * time.Second)
		}
	}
}

func listen(server *coalago.Server) {
	if err := server.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

func sigterm() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
}

func sendEvents(start bool) {
	if storage.Device.Events {
		if start {
			events.Send(new(events.Start))
		}
		events.Send(new(events.Register))
	}
}

func deviceSetup(_ *cobra.Command, _ []string) {
	fmt.Println("Setup Device")
	log.Init(false)
	crypto.Init(config.Get("localFiles.keys"))
	storage.Init(
		config.Get("localFiles.contract"),
		config.Get("localFiles.setup"),
		config.Get("everscale.elector"),
		config.Get("everscale.vendor"),
		config.Get("everscale.deviceAPI"),
		config.Get("info.type"),
		buildinfo.Version)
	printDevicePublicKey()
	if storage.Device.Address != "" {
		fmt.Println("Device contract address is already configured:")
		fmt.Println(storage.Device.Address)
		return
	}

	if err := storage.PromptForContractAddress(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Device contract address saved.")
	fmt.Printf("Address: %s\n", storage.Device.Address)
	if !isServiceRunning() {
		fmt.Printf("%s service is not running.\n", serviceName)
		fmt.Printf("to start it run the following command:\n")
		if runtime.GOOS == "darwin" {
			fmt.Printf("brew services start %s\n", serviceName)
		} else if runtime.GOOS == "linux" {
			fmt.Printf("systemctl start %s\n", serviceName)
		}
	}
}

func isServiceRunning() bool {
	out, err := exec.Command("pgrep", "-x", serviceName).Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) != ""
}

func deviceStatus(_ *cobra.Command, _ []string) {
	crypto.Init(config.Get("localFiles.keys"))
	storage.Init(
		config.Get("localFiles.contract"),
		config.Get("localFiles.setup"),
		config.Get("everscale.elector"),
		config.Get("everscale.vendor"),
		config.Get("everscale.deviceAPI"),
		config.Get("info.type"),
		buildinfo.Version)
	if storage.Device.Address == "" {
		fmt.Println("Status:  Not configured")
		fmt.Printf("Name:    %s\n", storage.Device.Name)
		fmt.Printf("PubKey:  %s\n", crypto.Keys.PublicSign)
		return
	}
	fmt.Print("Status:  ")
	input := map[string]string{"address": string(storage.Device.Address)}
	_, err := api.GET("device/info", input)
	if err != nil {
		fmt.Println("Offline")
	} else {
		fmt.Println("Online")
	}
	fmt.Printf("Name:    %s\n", storage.Device.Name)
	fmt.Printf("Address: %s\n", storage.Device.Address)
	fmt.Printf("Node:    %s\n", storage.Device.Node)
	fmt.Printf("Group:   %s\n", storage.Device.Group)
	fmt.Printf("Elector: %s\n", storage.Device.Elector)
	balance := api.GetBalance()
	fmt.Printf("Balance: %.9f TON\n", balance)
}

func waitForContractAddress() {
	if storage.HasContractAddress() {
		return
	}

	printDevicePublicKey()
	log.Info("Device contract address is not configured yet.")
	log.Info("Run `iot4b setup` to save the deployed contract address. The setup file will be created automatically.")
	for !storage.HasContractAddress() {
		time.Sleep(time.Second)
		if err := storage.Update(); err != nil {
			log.Error(err)
		}
	}
}

func printDevicePublicKey() {
	fmt.Println("Device public key:")
	fmt.Println(crypto.Keys.PublicSign)
	fmt.Println("Copy this key into the app, deploy the device, then paste the deployed contract address here.")
	fmt.Println("The setup file is managed automatically. You do not need to create or copy it by hand.")
}

func isVersionCommand(args []string) bool {
	if len(args) == 0 {
		return false
	}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		return arg == "version"
	}
	return false
}
