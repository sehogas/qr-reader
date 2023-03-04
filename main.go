package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sehogas/qr-reader/backend"
	"github.com/sehogas/qr-reader/models"
	"github.com/sehogas/qr-reader/util"
)

var ticker *time.Ticker
var done chan bool

func main() {

	qrReaderKey := os.Getenv("QR_READER_KEY")
	if qrReaderKey == "" {
		log.Fatal("No se encontró la clave del producto")
	}

	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	enviromentFileEcryptedParam := runCmd.String("file", "", "Archivo de variables de entorno encriptado")

	encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
	fileEncryptParam := encryptCmd.String("file", "", "Archivo a encriptar")

	decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
	fileDecryptParam := decryptCmd.String("file", "", "Archivo a desencriptar")

	if len(os.Args) < 2 {
		log.Fatal("Faltan parámetros")
	}

	// Check arguments
	switch os.Args[1] {
	case "encrypt":
		encryptCmd.Parse(os.Args[2:])
		if *fileEncryptParam == "" {
			log.Fatal("Archivo a encriptar requerido")
		}
	case "decrypt":
		decryptCmd.Parse(os.Args[2:])
		if *fileDecryptParam == "" {
			log.Fatal("Archivo a desencriptar requerido")
		}
	case "run":
		runCmd.Parse(os.Args[2:])
		if *enviromentFileEcryptedParam == "" {
			log.Fatal("Archivo de variables de entorno encriptado requerido")
		}
	default:
		log.Fatal("Se esperaba el comando 'run' ó 'encrypt' ó 'decrypt'")
	}

	log.Println("Inicio del programa...")
	switch os.Args[1] {
	case "encrypt":
		util.EncryptFile(*fileEncryptParam, []byte(qrReaderKey), ".encrypted")
	case "decrypt":
		util.DecryptFile(*fileDecryptParam, []byte(qrReaderKey), ".decrypted")
	case "run":
		cfg, err := util.GetConfigFromEncryptedFile(*enviromentFileEcryptedParam, []byte(qrReaderKey))
		if err != nil {
			log.Fatal(err)
		}
		util.CheckConfig(cfg)

		tmpDir := fmt.Sprintf("%s\\qr-reader", os.TempDir())
		if _, err := os.Stat(tmpDir); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(tmpDir, os.ModePerm)
			if err != nil {
				log.Fatalln(err)
			}
		}

		repo := util.NewRepository("sqlite3", cfg["DB"])
		defer repo.Db.Close()

		if totalCards, _ := repo.TotalCards(); totalCards >= 0 {
			fnSync(time.Now().Local(), cfg, repo)
		} else {
			log.Printf("Total de tarjetas locales: %d \n", totalCards)
		}

		syncTime, err := strconv.ParseInt(cfg["SYNC_TIME"], 10, 64)
		if err != nil {
			log.Fatal("Error convirtiendo parámetro SYNC_TIME")
		}

		ticker = time.NewTicker(time.Duration(syncTime) * time.Minute)
		done = make(chan bool)
		go inBackground(cfg, repo, fnSync)

		lectorQR := util.NewLectorQR(cfg, repo, tmpDir)
		lectorQR.Start()

		ticker.Stop()
		done <- true
	}
	log.Println("Fin del programa")
}

// Esta función llama a sincronización cada cierto tiempo
func inBackground(cfg map[string]string, repo *util.Repository, f func(t time.Time, cfg map[string]string, repo *util.Repository)) {
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			f(t, cfg, repo)
		}
	}
}

// Esta función sincroniza las tarjetas y los movimientos pendientes con el servidor
func fnSync(t time.Time, cfg map[string]string, repo *util.Repository) {

	var totalAccessSync int
	var status string = "OK"

	consultarAnulados := !repo.Config.LastUpdateCards.Equal(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))

	cards, syncTime, err := backend.GetCardsFromServer(cfg["URL_BACKEND"], cfg["API_KEY"], repo.Config.LastUpdateCards, consultarAnulados)
	if err != nil {
		log.Println("*** Error consultando servidor para sincronizar tarjetas ***")
		status = "ERROR"
	} else {
		err = repo.SyncCards(cards, syncTime)
		if err != nil {
			log.Println("*** Error actualización tarjetas locales ***")
			status = "ERROR"
		}
	}

	tSync := time.Now()
	access, err := repo.GetAccessToSync(tSync)
	if err != nil {
		log.Println("*** Error consultando movimientos para enviar al servidor ***")
		status = "ERROR"
	} else {
		totalAccessSync = len(access)
		if totalAccessSync > 0 {
			var bOk bool = true

			err := backend.SendToServerBulk(cfg["URL_BACKEND"], cfg["API_KEY"], models.AccessBulk{SyncDate: tSync, ClientID: cfg["CLIENT_ID"], Access: access})
			if err != nil {
				bOk = false
				totalAccessSync = 0
				log.Println("*** Error enviando pendientes al servidor ***")
				status = "ERROR"
			}
			err = repo.SyncAccessUpdateDelete(tSync, bOk)
			if err != nil {
				log.Println("*** Error actualizando pendientes locales ***")
				status = "ERROR"
			}
		}
	}

	totalCards, err := repo.TotalCards()
	if err != nil {
		log.Println("*** Error consultando total de tarjetas locales ***")
		status = "ERROR"
	}

	TotalEarrings, err := repo.TotalEarrings()
	if err != nil {
		log.Println("*** Error consultando total de pendientes locales ***")
		status = "ERROR"
	}

	log.Printf("Sincronización %s:  [Total tarjetas: %d], [Total recibido: %d], [Total a enviado: %d], [Total pendientes: %d], [Fecha del servidor: %s]\n",
		status, totalCards, len(cards), totalAccessSync, TotalEarrings, syncTime.Format("2006-01-02 15:04:05"))
}
