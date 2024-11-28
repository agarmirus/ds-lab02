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

type reservConfigDataStruct struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	ConnStr string `json:"connDb"`
}

func readConfig(path string, configData *reservConfigDataStruct) (err error) {
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

func buildService(configData *reservConfigDataStruct) (controller controllers.IController, err error) {
	hotelDAO := database.NewPostgresHotelDAO(configData.ConnStr)
	reservDAO := database.NewPostgresReservationDAO(configData.ConnStr)
	service := services.NewReservationService(reservDAO, hotelDAO)
	controller = controllers.NewReservationController(configData.Host, configData.Port, service)

	return controller, nil
}

func main() {
	var configData reservConfigDataStruct
	err := readConfig(`config/reservation.json`, &configData)

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
