package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sehogas/qr-reader/util"
)

var clientID string
var zoneID string
var eventID string
var deviceID int
var mode int
var db string
var rtsp string
var pathWavGranted string
var pathWavDenied string

func main() {

	log.Println("Iniciando...")

	flag.StringVar(&clientID, "client-id", "", "String value. Instance")
	flag.StringVar(&zoneID, "zone-id", "", "String value. Zone")
	flag.StringVar(&eventID, "event-id", "", "String value. Event. (I=In,O=Out)")
	flag.IntVar(&mode, "mode", 1, "Optional integer value. Read mode. Default=1. (1=Only persons,2:Persons+Vehicles)")
	flag.IntVar(&deviceID, "device-id", 0, "Optional integer value. Device identifier. Default=0. (0=Webcam default)")
	flag.StringVar(&db, "db-name", "data.db", "Optional string value. Database name")
	flag.StringVar(&rtsp, "rtsp", "", "Optional string value. Url video stream (rtsp)")
	flag.StringVar(&pathWavGranted, "path-wav-granted", "./assets/access-granted.wav", "Optional string value. File path .wav for access granted")
	flag.StringVar(&pathWavDenied, "path-wav-denied", "./assets/access-denied.wav", "Optional string value. File path .wav for access denied")

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

	if rtsp == "" {
		fmt.Printf("DeviceID: %d \n", deviceID)
	} else {
		fmt.Printf("RTSP: %s \n", rtsp)
	}
	fmt.Printf("Mode: %d \n", mode)
	fmt.Printf("ZoneID: %s \n", zoneID)
	fmt.Printf("EventID: %s \n", eventID)

	//Init audios

	_ = util.NewSound(pathWavDenied)
	//defer wavDenied.Close()

	//_ = util.NewSound(pathWavGranted)
	//defer wavGranted.Close()

	// defer repo.Close()

	// repo.InsertOrReplaceCards(util.TestCards())

	// lectorQR := util.NewLectorQR(deviceID, rtsp, repo, mode, clientID, zoneID, eventID, pathWavGranted, pathWavDenied)
	// lectorQR.Start()
	fmt.Scanln()
}
