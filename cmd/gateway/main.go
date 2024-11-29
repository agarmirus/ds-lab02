package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/agarmirus/ds-lab02/internal/controllers"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
	"github.com/agarmirus/ds-lab02/internal/services"
)

type gatewayConfigDataStruct struct {
	Host        string `json:"host"`
	LoayltyHost string `json:"loayltyHost"`
	PaymentHost string `json:"paymentHost"`
	ReservHost  string `json:"reservHost"`
	Port        int    `json:"port"`
	LoyaltyPort int    `json:"loyaltyPort"`
	PaymentPort int    `json:"paymentPort"`
	ReservPort  int    `json:"reservPort"`
}

func readConfig(path string, configData *gatewayConfigDataStruct) (err error) {
	configFile, err := os.Open(path)

	if err != nil {
		return err
	}

	defer configFile.Close()

	configJSON, err := io.ReadAll(configFile)

	if err != nil {
		return err
	}

	return json.Unmarshal(configJSON, configData)
}

func buildService(configData *gatewayConfigDataStruct) (controller controllers.IController, err error) {
	service := services.NewGatewayService(
		configData.ReservHost,
		configData.ReservPort,
		configData.PaymentHost,
		configData.PaymentPort,
		configData.LoayltyHost,
		configData.LoyaltyPort,
	)

	controller = controllers.NewGatewayController(
		configData.Host,
		configData.Port,
		service,
	)

	return controller, nil
}

func main() {
	var configData gatewayConfigDataStruct
	err := readConfig(`/configs/config.json`, &configData)

	if err != nil {
		panic(errors.New(serverrors.ErrConfigRead))
	}

	controller, err := buildService(&configData)

	if err != nil {
		panic(errors.New(serverrors.ErrServiceBuild))
	}

	err = controller.Prepare()

	if err != nil {
		panic(errors.New(serverrors.ErrControllerPrepare))
	}

	controller.Run()
}
