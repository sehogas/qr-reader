package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sehogas/qr-reader/util"
)

var clientID string
var zoneID string
var eventID string
var mode int

func main() {

	flag.StringVar(&clientID, "client-id", "", "String value. Instance")
	flag.StringVar(&zoneID, "zone", "", "String value. Zone")
	flag.StringVar(&eventID, "event", "", "String value. Event. (I=In,O=Out)")
	flag.IntVar(&mode, "mode", 1, "Int value. Read mode. Default=1. (1=Only persons,2:Persons+Vehicles)")
	flag.Parse()

	if clientID == "" {
		log.Fatal("Required [client-id] parameter")
	}

	if zoneID == "" {
		log.Fatal("Required [zone] parameter")
	}

	if eventID == "" {
		log.Fatal("Required [event] parameter")
	} else {
		if eventID == "I" || eventID == "O" {
			if eventID == "I" {
				eventID = "E"
			}
			if eventID == "O" {
				eventID = "S"
			}
		} else {
			log.Fatal("Invalid [event] parameter")
		}
	}

	if mode < 1 || mode > 2 {
		log.Fatal("Invalid mode")
	}

	deviceID, err := strconv.Atoi(os.Getenv("QRREADER_DEVICEID"))
	if err != nil {
		deviceID = 0 //webcam default
	}
	fromFile := os.Getenv("QRREADER_RTSP")
	dbPath := os.Getenv("QRREADER_DB")
	if dbPath == "" {
		dbPath = "./data.db"
	}

	if fromFile == "" {
		fmt.Printf("DeviceID: %d \n", deviceID)
	} else {
		fmt.Printf("RTSP: %s \n", fromFile)
	}
	//fmt.Printf("DB: %s \n", dbPath)
	fmt.Printf("Mode: %d \n", mode)
	fmt.Printf("ZoneID: %s \n", zoneID)
	fmt.Printf("EventID: %s \n", eventID)

	repo := util.NewRepository("sqlite3", dbPath)
	defer repo.Close()

	repo.InsertOrReplaceCards(util.TestCards())

	lectorQR := util.NewLectorQR(deviceID, fromFile, repo, mode, clientID, zoneID, eventID)
	lectorQR.Start()
}
