package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sehogas/qr-reader/sigep"
	"github.com/sehogas/qr-reader/util"
)

func main() {

	qrReaderKey := os.Getenv("QR_READER_KEY")
	if qrReaderKey == "" {
		log.Fatal("No enviroment key present")
	}

	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	enviromentFileEcryptedParam := runCmd.String("file", "", "Environment file encrypted")

	encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
	fileEncryptParam := encryptCmd.String("file", "", "File to encrypt")

	decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
	fileDecryptParam := decryptCmd.String("file", "", "File to decrypt")

	if len(os.Args) < 2 {
		log.Fatal("Missing parameters")
	}

	// Check arguments
	switch os.Args[1] {
	case "encrypt":
		encryptCmd.Parse(os.Args[2:])
		if *fileEncryptParam == "" {
			log.Fatal("File to encrypt required")
		}
	case "decrypt":
		decryptCmd.Parse(os.Args[2:])
		if *fileDecryptParam == "" {
			log.Fatal("File to encrypt required")
		}
	case "run":
		runCmd.Parse(os.Args[2:])
		if *enviromentFileEcryptedParam == "" {
			log.Fatal("Encrypted enviroment file required")
		}
	default:
		fmt.Println("Expected 'encrypt' or 'decrypt' subcommands")
		os.Exit(1)
	}

	log.Println("Starting...")
	switch os.Args[1] {
	case "encrypt":
		util.EncryptFile(*fileEncryptParam, []byte(qrReaderKey), ".encrypted")
	case "decrypt":
		util.DecryptFile(*fileDecryptParam, []byte(qrReaderKey), ".decrypted")
	case "run":
		mConfig, err := util.GetConfigFromEncryptedFile(*enviromentFileEcryptedParam, []byte(qrReaderKey))
		if err != nil {
			log.Fatal(err)
		}
		util.CheckConfig(mConfig)

		repo := util.NewRepository("sqlite3", mConfig["DB"])
		defer repo.Close()

		ticker := time.NewTicker(1 * time.Minute)
		done := make(chan bool)
		go inBackground(done, ticker, mConfig, repo)

		lectorQR := util.NewLectorQR(mConfig, repo)
		lectorQR.Start()

		ticker.Stop()
		done <- true

	}

	log.Println("Finish")
}

func inBackground(c chan bool, ticker *time.Ticker, cfg map[string]string, repo *util.Repository) {
	for {
		select {
		case <-c:
			return
		case t := <-ticker.C:
			log.Println("Sync at ", t.UTC())

			cards, err := sigep.GetCardsFromServer(cfg["URL_GET_CARDS"], cfg["API_KEY"], repo.Config.LastUpdateCards, consultarAnulados)
			if err != nil {
				log.Println("*** Error getting cards from server ***")
			}
			log.Printf("Total cards read: %d\n", len(cards))
			repo.SyncCards(cards)
			totalCards, err := repo.TotalCards()
			if err != nil {
				log.Println("*** Error checking total cards ***")
			}
			log.Printf("Total local cards: %d \n", totalCards)
		}
	}
}
