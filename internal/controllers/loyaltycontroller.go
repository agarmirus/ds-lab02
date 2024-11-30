package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
	"github.com/agarmirus/ds-lab02/internal/services"
)

type LoyaltyController struct {
	host string
	port int

	service services.ILoyaltyService
}

func NewLoyaltyController(
	host string,
	port int,
	service services.ILoyaltyService,
) IController {
	return &LoyaltyController{host, port, service}
}

func (controller *LoyaltyController) handleLoyaltyByUsernameGet(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] LoyaltyController.handleLoyaltyByUsernameGet. Handling loyalty by username GET request")

	username := req.Header.Get(`X-User-Name`)

	if strings.Trim(username, ` `) == `` {
		log.Println("[ERROR] LoyaltyController.handleLoyaltyByUsernameGet. Invalid username: " + username)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	loyalty, err := controller.service.ReadLoyaltyByUsername(username)

	if err != nil {
		log.Println("[ERROR] LoyaltyController.handleLoyaltyByUsernameGet. service.ReadLoyaltyByUsername returned error: ", err)
		if errors.Is(err, errors.New(serverrors.ErrEntityNotFound)) {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	loyaltyJSON, err := json.Marshal(loyalty)

	if err != nil {
		log.Println("[ERROR] LoyaltyController.handleLoyaltyByUsernameGet. Cannot convert result into JSON format: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Add(`Content-Type`, `application/json`)
	res.Write(loyaltyJSON)
}

func (controller *LoyaltyController) handleLoyaltyByIdPatch(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] LoyaltyController.handleLoyaltyByIdPatch. Handling loyalty by id PATCH request")

	loyaltyId, err := strconv.Atoi(req.PathValue(`loyaltyId`))

	if err != nil {
		log.Println("[ERROR] LoyaltyController.handleLoyaltyByIdPatch. Invalid loyalty id: ", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var reqBody []byte
	n, err := req.Body.Read(reqBody)

	if err != nil || n <= 0 {
		log.Println("[ERROR] LoyaltyController.handleLoyaltyByIdPatch. Error while reading request body: ", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var loyalty models.Loyalty
	err = json.Unmarshal(reqBody, &loyalty)

	if err != nil {
		log.Println("[ERROR] LoyaltyController.handleLoyaltyByIdPatch. Error while parsing JSON request body: ", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	loyalty.Id = loyaltyId
	_, err = controller.service.UpdateLoyaltyById(&loyalty)

	if err != nil {
		log.Println("[ERROR] LoyaltyController.handleLoyaltyByIdPatch. service.UpdateLoyaltyById returned error: ", err)
		if errors.Is(err, errors.New(serverrors.ErrEntityNotFound)) {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	res.WriteHeader(http.StatusOK)
}

func (controller *LoyaltyController) handleLoyaltyRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		log.Println("[INFO] LoyaltyController.handleLoyaltyByIdRequest. Got loyalty by username GET request")
		controller.handleLoyaltyByUsernameGet(res, req)
	} else {
		log.Println("[ERROR] LoyaltyController.handleLoyaltyRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *LoyaltyController) handleLoyaltyByIdRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `PATCH` {
		log.Println("[INFO] LoyaltyController.handleLoyaltyByIdRequest. Got loyalty by id PATCH request")
		controller.handleLoyaltyByIdPatch(res, req)
	} else {
		log.Println("[ERROR] LoyaltyController.handleLoyaltyByIdRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *LoyaltyController) handleHealthRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		log.Println("[INFO] LoyaltyController.handleHealthRequest. Got health GET request")
		res.WriteHeader(http.StatusOK)
	} else {
		log.Println("[ERROR] LoyaltyController.handleHealthRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *LoyaltyController) Prepare() error {
	http.HandleFunc(`/api/v1/loyalty`, controller.handleLoyaltyRequest)
	http.HandleFunc(`/api/v1/loyalty/{loyaltyId}`, controller.handleLoyaltyByIdRequest)

	http.HandleFunc(`/manage/health`, controller.handleHealthRequest)

	return nil
}

func (controller *LoyaltyController) Run() error {
	return http.ListenAndServe(
		fmt.Sprintf(
			`%s:%d`,
			controller.host,
			controller.port,
		),
		nil,
	)
}
