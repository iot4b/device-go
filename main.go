package main

import (
	"device-go/packages/aliver"
	"device-go/packages/config"
	"device-go/packages/crypto"
	"device-go/packages/events"
	"device-go/packages/handlers"
	"device-go/packages/registration"
	"device-go/packages/storage"
	"os"
	"os/exec"
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

func main() {
	var rootCmd = &cobra.Command{
		Use:   "device",
		Short: "Device CLI",
		Run:   runDevice,
	}
	rootCmd.PersistentFlags().StringVar(&env, "env", "prod", "Set environment")
	rootCmd.PersistentFlags().StringVar(&port, "port", "5684", "Set coala port")

	config.Init(env)
	log.Init(config.Bool("debug"))

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the device",
		Run:   initDevice,
	}

	rootCmd.AddCommand(initCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runDevice(_ *cobra.Command, _ []string) {
	crypto.Init(config.Get("localFiles.keys"))

	log.Info("Waiting for contract data file...")
	storage.WaitForData()

	// сервер для запросов от клиентов и нод
	server := coalago.NewServer()
	server.GET("/info", handlers.GetInfo)
	server.POST("/cmd", handlers.ExecCmd)
	server.POST("/update", handlers.Update)
	server.POST("/sign", handlers.Sign)

	for {
		// регистрируем устройство на ноде
		// если ошибка, то повторяем цикл регистрации
		err := registration.Register()
		if err == nil {
			// стартуем сервер
			go listen(server)

			// delay for alive correct work
			time.Sleep(time.Second)

			// начинаем слать alive пакеты, чтобы сохранять соединение для udp punching
			go aliver.Run(server, storage.Device.Address.String(), config.Time("timeout.alive"))

			if storage.Device.Events {
				// delay after first alive to store ip:port and send event
				time.Sleep(time.Second)
				events.Send(new(events.Start))
			}

			time.Sleep(config.Time("timeout.registerRepeat"))
			go registration.Repeat()

			break
		} else {
			log.Error(err)
			time.Sleep(3 * time.Second)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func listen(server *coalago.Server) {
	if err := server.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

func initDevice(_ *cobra.Command, _ []string) {
	log.Info("Init Device")
	storage.Init(
		config.Get("localFiles.contract"),
		config.Get("localFiles.init"),
		config.Get("everscale.elector"),
		config.Get("everscale.vendor"),
		config.Get("everscale.deviceAPI"),
		config.Get("info.type"),
		config.Get("info.version"))
	if storage.Device.Address != "" {
		log.Info("Device contract is already deployed. Address:")
		log.Info(storage.Device.Address)
		return
	}
	if !isServiceRunning() {
		log.Info("iot4b-device service is not running.")
		log.Info("to start it run the following command in a separate terminal:")
		log.Info("brew services start iot4b-device")
		for {
			if isServiceRunning() {
				break
			}
			time.Sleep(time.Second)
		}
	}

	log.Info("Waiting for contract deployment...")
	for {
		storage.Update()
		if storage.Device.Address != "" {
			log.Info("Device contract is deployed.")
			log.Infof("Address: %s", storage.Device.Address)
			log.Infof("Group:   %s", storage.Device.Group)
			log.Infof("Elector: %s", storage.Device.Elector)
			return
		}
		time.Sleep(time.Second)
	}
}

func isServiceRunning() bool {
	out, err := exec.Command("pgrep", "-x", "iot4b-device").Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) != ""
}
