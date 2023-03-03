package util

import (
	"log"
	"strconv"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/sehogas/qr-reader/models"
	"github.com/sehogas/qr-reader/sigep"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

const (
	_PERSON  = "PERSONA"
	_VEHICLE = "VEHICULO"
	_UNKNOWN = "DESCONOCIDO"
)

const (
	_ACCESS_GRANTED = "ACCESS_GRANTED"
	_ACCESS_DENIED  = "ACCESS_DENIED"
	_CONTINUE       = "CONTINUE"
	_ERROR          = "ERROR"
)

type LectorQR struct {
	DeviceID            int
	FromFile            string
	Repo                *Repository
	Mode                int
	ClientID            string
	Zone                string
	Event               string
	FileWavGranted      string
	FileWavDenied       string
	UrlPostAccess       string
	UrlPostAccessAPIKey string
	FileImageGranted    string
	FileImageDenied     string
	DebugMode           bool
	FileDetectProtoTxt  string
	FileDetectCaffe     string
	FileSuperProtoTxt   string
	FileSuperCaffe      string
}

var QRPrev1 string
var QRPrev2 string
var QR1 string
var QR2 string

// func NewLectorQR(deviceID int, fromFile string, repo *Repository, mode int, clientID, zone, event, pathWavGranted, pathWavDenied string) *LectorQR {
func NewLectorQR(cfg map[string]string, repo *Repository) *LectorQR {
	deviceID, _ := strconv.Atoi(cfg["DEVICE_ID"])
	mode, _ := strconv.Atoi(cfg["MODE"])
	return &LectorQR{
		DeviceID:            deviceID,
		FromFile:            cfg["RTSP"],
		Repo:                repo,
		Mode:                mode,
		ClientID:            cfg["CLIENT_ID"],
		Zone:                cfg["ZONE_ID"],
		Event:               cfg["EVENT_ID"],
		FileWavGranted:      cfg["FILE_WAV_GRANTED"],
		FileWavDenied:       cfg["FILE_WAV_DENIED"],
		UrlPostAccess:       cfg["URL_POST_ACCESS"],
		UrlPostAccessAPIKey: cfg["API_KEY"],
		FileImageGranted:    cfg["FILE_BACKGROUND_GRANTED"],
		FileImageDenied:     cfg["FILE_BACKGROUND_DENIED"],
		DebugMode:           false,
		FileDetectProtoTxt:  cfg["FILE_DETECT_PROTO_TXT"],
		FileDetectCaffe:     cfg["FILE_DETECT_CAFFE"],
		FileSuperProtoTxt:   cfg["FILE_SUPER_PROTO_TXT"],
		FileSuperCaffe:      cfg["FILE_SUPER_CAFFE"],
	}
}

func (s *LectorQR) Start() {

	wavDenied := NewSound(s.FileWavDenied)
	wavDenied.Play()
	wavGranted := NewSound(s.FileWavGranted)
	wavGranted.Play()

	green := gocv.IMRead(s.FileImageGranted, gocv.IMReadColor)
	defer green.Close()

	red := gocv.IMRead(s.FileImageDenied, gocv.IMReadColor)
	defer red.Close()

	var camera *gocv.VideoCapture
	var err error
	var key, keyPrev int

	if s.FromFile == "" {
		camera, err = gocv.VideoCaptureDevice(s.DeviceID)
		if err != nil {
			log.Printf("*** Error abriendo cámara (device ID.:%d) ***\n", s.DeviceID)
			return
		}
	} else {
		camera, err = gocv.VideoCaptureFile(s.FromFile)
		if err != nil {
			log.Println("*** Error abriendo cámara por protocolo RTSP ***")
			return
		}
	}
	defer camera.Close()

	frame := gocv.NewMat()
	defer frame.Close()

	mats := make([]gocv.Mat, 0)

	frameClearBuffer := gocv.NewMat()
	defer frameClearBuffer.Close()

	window := gocv.NewWindow(s.ClientID)
	defer window.Close()

	var done chan bool

	var decoded string
	var decoded2 string
	// var img image.Image
	var status string
	var code1 string
	var code2 string
	var QRsCount int

	wait := 1

	detector := contrib.NewWeChatQRCode(s.FileDetectProtoTxt, s.FileDetectCaffe, s.FileSuperProtoTxt, s.FileSuperCaffe)

	log.Println("Leyendo cámara...")

	for {
		if !camera.Read(&frame) {
			log.Println("*** No se pudo leer la cámara ***")
			break
		}
		if frame.Empty() {
			continue
		}

		//Esta rutina limpia el buffer de la cámara cuando es por rtsp
		//////////////////////////////////////////////////////////////
		done = make(chan bool)
		go func(f gocv.Mat) {
			for {
				select {
				case <-done:
					return
				default:
					camera.Read(&f)
				}
			}
		}(frameClearBuffer)
		///////////////////////////////////////////////////////////////

		qrCodes := detector.DetectAndDecode(frame, &mats)
		QRsCount = len(qrCodes)
		if QRsCount > 0 {
			decoded = string(qrCodes[0])
			if QRsCount == 2 {
				decoded2 = string(qrCodes[1])
			}
			if decoded != "" {
				//log.Printf("decoded: %s, decoded2: %s\n", decoded, decoded2)
				status, code1, code2 = s.AccessGranted2(&decoded, &decoded2)
				//log.Printf("status: %s, code1: %s, code2: %s\n", status, code1, code2)
				switch status {
				case _ACCESS_GRANTED:
					log.Printf("### ACCESO PERMITIDO ###\n")
					green.CopyTo(&frame)
					wavGranted.Play()
					wait = 2000
				case _ACCESS_DENIED:
					log.Printf("### ACCESO DENEGADO ###\n")
					red.CopyTo(&frame)
					wavDenied.Play()
					wait = 2000
				case _CONTINUE:
				case _ERROR:
					break
				}
				decoded = ""
				decoded2 = ""
			}
		}

		window.IMShow(frame)
		if status == _ACCESS_GRANTED {
			s.SaveAccessGranted(code1, code2)
			status = _CONTINUE
		}
		code1 = ""
		code2 = ""
		key = window.WaitKey(wait)
		done <- true
		if key == 27 {
			//Esc
			if key == keyPrev {
				break
			}
		}
		keyPrev = key
		switch key {
		case 100:
			//DebugMode
			s.DebugMode = !s.DebugMode
			log.Printf("DEBUG MODE: %v\n", s.DebugMode)
		}
		// if key != -1 {
		// 	fmt.Println(key)
		// }

		wait = 1
	}
	log.Println("Cámara cerrada")
}

func (s *LectorQR) SaveAccessGranted(code1, code2 string) {

	access := models.Access{
		UUID:       s.Repo.NewUUID(),
		Code1:      code1,
		Code2:      code2,
		AccessDate: time.Now(),
		Zone:       s.Zone,
		Event:      s.Event,
	}

	go func(url, APIKey string, access *models.Access) {
		dataAccess, err := sigep.SendToServer(url, s.UrlPostAccessAPIKey, *access)
		if err != nil {
			LogError("*** Error enviando movimiento al servidor ***", err, s.DebugMode)

			err := s.Repo.InsertAccess(access)
			if err != nil {
				LogError("Amacenamiento local: ERROR", err, s.DebugMode)
			}
			log.Println("Almacenamiento en local: OK")
		} else {
			sigep.PrintData(dataAccess)
		}
	}(s.UrlPostAccess, s.UrlPostAccessAPIKey, &access)
}

func (s *LectorQR) AccessGranted2(decoded *string, decoded2 *string) (status, code1, code2 string) {

	if s.Mode == 1 { //Only person
		if *decoded != QRPrev1 && IsPerson(*decoded) {
			QRPrev1 = *decoded
			log.Printf("QR [%s] es %s\n", *decoded, PersonOrVehicle(*decoded))
			if Ok, _ := s.Repo.ValidCard(QRPrev1); Ok {
				return _ACCESS_GRANTED, QRPrev1, ""
			} else {
				return _ACCESS_DENIED, QRPrev1, ""
			}
		} else {
			return _CONTINUE, "", ""
		}
	}

	if s.Mode == 2 { //Persons + Vehicles
		if (*decoded == QRPrev1 && *decoded2 == QRPrev2) || (*decoded == QRPrev2 && *decoded2 == QRPrev1) || (*decoded == *decoded2) || (*decoded2 == "") {
			return _CONTINUE, *decoded, *decoded2
		}
		QRPrev1 = *decoded
		QRPrev2 = *decoded2

		log.Printf("QR1 [%s] es %s, QR2 [%s] es %s\n", *decoded, PersonOrVehicle(*decoded), *decoded2, PersonOrVehicle(*decoded2))
		if (PersonOrVehicle(*decoded) == _PERSON && PersonOrVehicle(*decoded2) == _VEHICLE) || (PersonOrVehicle(*decoded) == _VEHICLE && PersonOrVehicle(*decoded2) == _PERSON) {
			accessGranted, _ := s.Repo.ValidCard(*decoded)
			if accessGranted {
				accessGranted, _ = s.Repo.ValidCard(*decoded2)
				if accessGranted {
					if IsPerson(QRPrev1) {
						code1 = QRPrev1
						code2 = QRPrev2
					} else {
						code1 = QRPrev2
						code2 = QRPrev1
					}
					return _ACCESS_GRANTED, code1, code2
				} else {
					QRPrev1 = ""
					QRPrev2 = ""
					return _ACCESS_DENIED, "", ""
				}
			} else {
				QRPrev1 = ""
				QRPrev2 = ""
				return _ACCESS_DENIED, "", ""
			}
		} else {
			return _ACCESS_DENIED, "", ""
		}
	}
	log.Println("Modo inválido!")
	return _ERROR, "", ""
}
