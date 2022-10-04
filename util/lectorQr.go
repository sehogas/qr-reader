package util

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/liyue201/goqr"
	"github.com/sehogas/qr-reader/models"
	"gocv.io/x/gocv"
)

const (
	_PERSON  = "Person"
	_VEHICLE = "Vehicle"
	_OTHER   = "Other"
)

type LectorQR struct {
	DeviceID       int
	FromFile       string
	Repo           *Repository
	Mode           int
	ClientID       string
	Zone           string
	Event          string
	PathWavGranted string
	PathWavDenied  string
}

func NewLectorQR(deviceID int, fromFile string, repo *Repository, mode int, clientID, zone, event, pathWavGranted, pathWavDenied string) *LectorQR {
	return &LectorQR{
		DeviceID:       deviceID,
		FromFile:       fromFile,
		Repo:           repo,
		Mode:           mode,
		ClientID:       clientID,
		Zone:           zone,
		Event:          event,
		PathWavGranted: pathWavGranted,
		PathWavDenied:  pathWavDenied,
	}
}

func (s *LectorQR) Start() {

	var camera *gocv.VideoCapture
	var err error
	var keyPrev int
	var key int

	if s.FromFile == "" {
		camera, err = gocv.VideoCaptureDevice(s.DeviceID)
	} else {
		camera, err = gocv.VideoCaptureFile(s.FromFile)
	}
	if err != nil {
		log.Println(err)
		return
	}
	defer camera.Close()

	frame := gocv.NewMat()
	defer frame.Close()

	//frameGray := gocv.NewMat()
	//defer frameGray.Close()

	//pts := gocv.NewMat()
	//defer pts.Close()

	//straight_qrcode := gocv.NewMat()
	//defer straight_qrcode.Close()

	green := gocv.IMRead("./assets/access-granted.png", gocv.IMReadColor)
	defer green.Close()

	red := gocv.IMRead("./assets/access-denied.png", gocv.IMReadColor)
	defer red.Close()

	window := gocv.NewWindow(s.ClientID)
	defer window.Close()

	//qrCodeDetector := gocv.NewQRCodeDetector()
	//defer qrCodeDetector.Close()

	var decoded string
	var prev1 string
	var prev2 string
	var accessGranted bool
	wait := 1

	log.Println("Reading camera...")

	for {
		if ok := camera.Read(&frame); !ok {
			log.Println("Could not read the camera")
			break
		}

		if frame.Empty() {
			continue
		}

		//decoded = qrCodeDetector.DetectAndDecode(frame, &pts, &straight_qrcode)

		//gocv.CvtColor(frame, &frameGray, gocv.ColorBGRToGray)

		img, err := frame.ToImage()

		if err == nil {
			qrCodes, err := goqr.Recognize(img)
			if err == nil {
				if len(qrCodes) == 1 {
					decoded = string(qrCodes[0].Payload)
					if decoded != "" {
						if s.Mode == 1 { //Only person
							if prev1 != decoded && IsPerson(decoded) {
								prev1 = decoded
								log.Printf("Card 1 [%s] is %s \n", prev1, PersonOrVehicle(prev1))

								if accessGranted, _ = s.Repo.ValidCard(prev1); accessGranted {

									access := models.Access{
										Code1:      prev1,
										Code2:      "",
										AccessDate: time.Now(),
										Zone:       s.Zone,
										Event:      s.Event,
									}
									message := models.AccessZone{
										ClientID: s.ClientID,
										Access:   []models.Access{access},
									}

									if !s.SendToServer(&message) { //Envia al servidor
										err := s.Repo.InsertAccess(&access) //Sino intento grabar en base de datos local
										if err != nil {
											log.Println("Local storage: ", err)
											break
										}
										log.Println("Recorded to local storage: OK")
									}
									log.Println("Access granted!")
									green.CopyTo(&frame)
									// err = wavGranted.Play()
									// if err != nil {
									// 	log.Println(err)
									// }
								} else {
									log.Println("Access denied!")
									red.CopyTo(&frame)
									// err = wavDenied.Play()
									// if err != nil {
									// 	log.Println(err)
									// }
								}
								accessGranted = false
								wait = 2000
							}
						}

						if s.Mode == 2 { //Persons + Vehicles
							if prev1 == "" && prev2 != decoded {
								prev1 = decoded
								prev2 = ""
								log.Printf("Card 1 [%s] is %s\nWaiting for next card...\n", prev1, PersonOrVehicle(prev1))
							} else {
								if prev2 == "" && prev1 != decoded && PersonOrVehicle(decoded) != _OTHER && PersonOrVehicle(prev1) != PersonOrVehicle(decoded) {
									prev2 = decoded
									log.Printf("Card 2 [%s] is %s\n", prev2, PersonOrVehicle(prev2))
								}
							}

							if prev1 != "" && prev2 != "" {
								accessGranted, _ = s.Repo.ValidCard(prev1)
								if accessGranted {
									accessGranted, _ = s.Repo.ValidCard(prev2)
									if accessGranted {
										//log.Println("Access granted!")
										//green.CopyTo(&frame)

										var code1, code2 string
										if IsPerson(prev1) {
											code1 = prev1
											code2 = prev2
										} else {
											code1 = prev2
											code2 = prev1
										}

										access := models.Access{
											Code1:      code1,
											Code2:      code2,
											AccessDate: time.Now(),
											Zone:       s.Zone,
											Event:      s.Event,
										}
										message := models.AccessZone{
											ClientID: s.ClientID,
											Access:   []models.Access{access},
										}

										if !s.SendToServer(&message) { //Envia al servidor
											err := s.Repo.InsertAccess(&access) //Sino intento grabar en base de datos local
											if err != nil {
												log.Println("Local storage: ", err)
												break
											}
											log.Println("Recorded to local storage: OK")
										}
										log.Println("Access granted!")
										green.CopyTo(&frame)
										// wavGranted.Play()
									} else {
										log.Println("Access denied!")
										red.CopyTo(&frame)
										// wavDenied.Play()
									}
								} else {
									log.Println("Access denied!")
									red.CopyTo(&frame)
									// wavDenied.Play()
								}
								prev1 = ""
								accessGranted = false
								wait = 2000
							}
						}
					}
				}
			}
		}

		window.IMShow(frame)
		key = window.WaitKey(wait)
		if key == 27 {
			if key == 27 && key == keyPrev {
				break
			}
		}
		keyPrev = key
		wait = 1
	}
	log.Println("Done")
}

func IsVehicle(qr string) bool {
	return (len(qr) == 40 && qr[0:3] == "002")
}

func IsPerson(qr string) bool {
	return (len(qr) == 40 && qr[0:3] == "001")
}

func PersonOrVehicle(qr string) string {
	if len(qr) == 40 {
		if qr[0:3] == "001" {
			return _PERSON
		}
		if qr[0:3] == "002" {
			return _VEHICLE
		}
	}
	return "Other"
}

func (s *LectorQR) SendToServer(message *models.AccessZone) bool {
	body, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return false
	}
	res, err := http.Post("https://backend.dpp.gob.ar/api/v1/prueba", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		return false
	}
	if !(res.StatusCode == 200 || res.StatusCode == 201) {
		log.Println("SendToServer() error:", http.StatusText(res.StatusCode))
		return false
	}
	return true
}
