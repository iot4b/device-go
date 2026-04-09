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
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
	"github.com/spf13/cobra"
)

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

	// сервер для запросов от клиентов и нод
	server := coalago.NewServer()
	server.GET("/info", handlers.GetInfo)
	server.POST("/cmd", handlers.ExecCmd)
	server.POST("/update", handlers.Update)
	server.POST("/sign", handlers.Sign)

	go listen(server)
	go sigterm()
	go registration.PairingHeartbeat()

	for {
		var err error
		if storage.Device.Address == "" {
			err = registration.Pair()
			if err == nil {
				time.Sleep(2 * time.Second)
				continue
			}
		} else {
			err = registration.Register()
			if err == nil {
				//go sendEvents(start)

				// начинаем слать alive пакеты, чтобы сохранять соединение для udp punching
				aliver.Run(server, config.Time("timeout.alive"))
				continue
			}
		}

		if err != nil {
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
	registration.Setup()
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
		fmt.Println("Status:  Pairing")
		fmt.Printf("Name:         %s\n", storage.Device.Name)
		fmt.Printf("PubKey:       %s\n", crypto.Keys.PublicSign)
		if storage.Pairing.Code != "" {
			fmt.Printf("Pairing Code: %s\n", storage.Pairing.Code)
			fmt.Printf("Node:         %s\n", storage.Pairing.NodeIpPort)
		}
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
