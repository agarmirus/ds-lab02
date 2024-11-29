package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/agarmirus/ds-lab02/internal/controllers"
	"github.com/agarmirus/ds-lab02/internal/database"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
	"github.com/agarmirus/ds-lab02/internal/services"
)

type paymentConfigDataStruct struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	ConnStr string `json:"connDb"`
}

func readConfig(path string, configData *paymentConfigDataStruct) (err error) {
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

func buildService(configData *paymentConfigDataStruct) (controller controllers.IController, err error) {
	paymentDAO := database.NewPostgresPaymentDAO(configData.ConnStr)
	service := services.NewPaymentService(paymentDAO)
	controller = controllers.NewPaymentController(configData.Host, configData.Port, service)

	return controller, nil
}

func main() {
	var configData paymentConfigDataStruct
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
