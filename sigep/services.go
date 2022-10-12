package sigep

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sehogas/qr-reader/models"
)

func SendToServerBulk(url string, apiKey string, accessBulk models.AccessBulk) error {

	client := &http.Client{Timeout: time.Second * 5}

	bodyReq, err := json.Marshal(accessBulk)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyReq))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", apiKey)

	res, err := client.Do(request)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("[%d] %s", res.StatusCode, http.StatusText(res.StatusCode))
	}
	defer res.Body.Close()

	return nil
}

func SendToServer(url string, apiKey string, access models.Access) (models.AccessDataResponse, error) {
	var accessData models.AccessDataResponse

	client := &http.Client{Timeout: time.Second * 3}

	bodyReq, err := json.Marshal(access)
	if err != nil {
		return accessData, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyReq))
	if err != nil {
		return accessData, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", apiKey)

	res, err := client.Do(request)
	if err != nil {
		return accessData, err
	}

	if res.StatusCode != http.StatusCreated {
		return accessData, fmt.Errorf("[%d] %s", res.StatusCode, http.StatusText(res.StatusCode))
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return accessData, err
	}
	err = json.Unmarshal(body, &accessData)
	if err != nil {
		return accessData, err
	}

	return accessData, nil
}

func GetCardsFromServer(url string, apiKey string, fecDesde time.Time, anulados bool) ([]models.Card, time.Time, error) {
	var cards []models.Card
	var syncTime time.Time = time.Now()
	var iAnulados int8

	if anulados {
		iAnulados = 1
	} else {
		iAnulados = 0
	}

	client := &http.Client{Timeout: time.Second * 5}

	req := fmt.Sprintf("%s/%s/%d", url, fecDesde.Format("20060102150405"), iAnulados)
	request, err := http.NewRequest(http.MethodGet, req, nil)
	if err != nil {
		return cards, syncTime, err
	}

	request.Header.Set("x-api-key", apiKey)
	res, err := client.Do(request)
	if err != nil {
		return cards, syncTime, err
	}

	if res.StatusCode != http.StatusOK {
		return cards, syncTime, fmt.Errorf("[%d] %s", res.StatusCode, http.StatusText(res.StatusCode))
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return cards, syncTime, err
	}

	var response models.QRUpdatesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return cards, syncTime, err
	}

	return response.Cards, response.ServerTime, nil
}

func PrintData(data models.AccessDataResponse) {
	var sEventual string = "NO"
	if data.Eventual {
		sEventual = "SI"
	}
	log.Printf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
		fmt.Sprintf("Tipo y Nro.Doc.  : %s %s", data.DocumentType, data.DocumentNumber),
		fmt.Sprintf("%20sApellido y Nombre: %s", "", data.PersonName),
		fmt.Sprintf("%20sVigencia Tarjeta : %s - %s", "", data.DateFrom.Format("02/01/2006"), data.DateTo.Format("02/01/2006")),
		fmt.Sprintf("%20sProntuario P.N.A.: %s", "", data.PNA),
		fmt.Sprintf("%20sEventual         : %s", "", sEventual),
		fmt.Sprintf("%20sColor de Tarjeta : %s", "", GetColorName(data.Color)),
		fmt.Sprintf("%20sDominio vehículo : %s", "", data.LicensePlate),
	)
}

func GetColorName(code string) string {
	var color string
	switch code {
	case "VE":
		color = "VERDE"
	case "RO":
		color = "ROJO"
	case "AM":
		color = "AMARILLO"
	case "NA":
		color = "NARANJA"
	default:
		color = "OTRO"
	}
	return color
}
