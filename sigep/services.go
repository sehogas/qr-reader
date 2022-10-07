package sigep

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sehogas/qr-reader/models"
)

func GetCardsFromServer(url string, apiKey string, fecDesde time.Time, anulados bool) ([]models.Card, error) {
	var cards []models.Card

	var iAnulados int8
	if anulados {
		iAnulados = 1
	} else {
		iAnulados = 0
	}

	client := &http.Client{Timeout: time.Second * 5}
	req := fmt.Sprintf("%s/%s/%d", url, fecDesde.Format("20060102150405"), iAnulados)
	request, err := http.NewRequest("GET", req, nil)
	if err != nil {
		return cards, err
	}

	request.Header.Set("x-api-key", apiKey)
	res, err := client.Do(request)
	if err != nil {
		return cards, err
	}

	if res.StatusCode != 200 {
		return cards, fmt.Errorf("[%d] %s", res.StatusCode, http.StatusText(res.StatusCode))
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return cards, err
	}
	err = json.Unmarshal(body, &cards)
	if err != nil {
		return cards, err
	}

	return cards, nil
}
