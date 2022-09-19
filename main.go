package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sehogas/qr-reader/util"
)

func main() {
	deviceID, err := strconv.Atoi(os.Getenv("DEVICEID"))
	if err != nil {
		deviceID = 0 //webcam default
	}
	fromFile := os.Getenv("RTSP")
	dbPath := os.Getenv("DB")
	if dbPath == "" {
		dbPath = "./data.db"
	}
	fmt.Printf("DEVICEID: %d \n", deviceID)
	fmt.Printf("RTSP: %s \n", fromFile)
	fmt.Printf("DB: %s \n", dbPath)

	repo := util.NewRepository("sqlite3", dbPath)
	defer repo.Close()

	repo.InsertOrReplaceCards(util.TestCards())

	lectorQR := util.NewLectorQR(deviceID, fromFile, repo)
	lectorQR.Start()
}
