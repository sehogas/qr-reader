package util

import (
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"time"

	"image/color"

	"github.com/sehogas/qr-reader/backend"
	"github.com/sehogas/qr-reader/models"
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
	DeviceID           int
	FromFile           string
	Repo               *Repository
	Mode               int
	ClientID           string
	Zone               string
	EventCode          string
	EventName          string
	FileWavGranted     string
	FileWavDenied      string
	APIKey             string
	DebugMode          *bool
	FileDetectProtoTxt string
	FileDetectCaffe    string
	FileSuperProtoTxt  string
	FileSuperCaffe     string
	TmpDir             string
	UrlBackend         string
	quad1              gocv.Mat
	quad2              gocv.Mat
	quad3              gocv.Mat
	quad4              gocv.Mat
	quadUp             gocv.Mat
	quadDown           gocv.Mat
	quadTmp            gocv.Mat
	frameResize        gocv.Mat
	Width              int
	Height             int
}

var QRPrev1 string
var QRPrev2 string
var QR1 string
var QR2 string

func NewLectorQR(cfg map[string]string, repo *Repository, tmpDir string, debugMode *bool) *LectorQR {
	deviceID, _ := strconv.Atoi(cfg["DEVICE_ID"])
	mode, _ := strconv.Atoi(cfg["MODE"])

	var eventName string
	if cfg["EVENT_ID"] == "E" {
		eventName = "ENTRADA"
	} else {
		eventName = "SALIDA"
	}

	return &LectorQR{
		DeviceID:           deviceID,
		FromFile:           cfg["RTSP"],
		Repo:               repo,
		Mode:               mode,
		ClientID:           cfg["CLIENT_ID"],
		Zone:               cfg["ZONE_ID"],
		EventCode:          cfg["EVENT_ID"],
		EventName:          eventName,
		FileWavGranted:     cfg["FILE_WAV_GRANTED"],
		FileWavDenied:      cfg["FILE_WAV_DENIED"],
		APIKey:             cfg["API_KEY"],
		DebugMode:          debugMode,
		FileDetectProtoTxt: cfg["FILE_DETECT_PROTO_TXT"],
		FileDetectCaffe:    cfg["FILE_DETECT_CAFFE"],
		FileSuperProtoTxt:  cfg["FILE_SUPER_PROTO_TXT"],
		FileSuperCaffe:     cfg["FILE_SUPER_CAFFE"],
		TmpDir:             tmpDir,
		UrlBackend:         cfg["URL_BACKEND"],
		Width:              640,
		Height:             480,
	}
}

func (s *LectorQR) Start() {
	var err error

	wavDenied := NewSound(s.FileWavDenied)
	wavDenied.Play()
	wavGranted := NewSound(s.FileWavGranted)
	wavGranted.Play()

	var camera *gocv.VideoCapture
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

	if !camera.Read(&frame) {
		log.Println("*** No se pudo leer la cámara ***")
		return
	}
	s.Width = frame.Cols()
	s.Height = frame.Rows()

	log.Printf("Resolución video: %d x %d\n", s.Width, s.Height)

	black, err := gocv.ImageToMatRGB(CreateImage(s.Width, s.Height, color.RGBA{0, 0, 0, 0}))
	if err != nil {
		log.Println("*** Error generando background black ***")
		return
	}
	defer black.Close()

	green, err := gocv.ImageToMatRGB(CreateImage(s.Width, s.Height, color.RGBA{0, 255, 0, 0}))
	if err != nil {
		log.Println("*** Error generando background green ***")
		return
	}
	defer green.Close()

	red, err := gocv.ImageToMatRGB(CreateImage(s.Width, s.Height, color.RGBA{255, 0, 0, 0}))
	if err != nil {
		log.Println("*** Error generando background green ***")
		return
	}
	defer red.Close()

	s.quad1 = gocv.NewMat()
	defer s.quad1.Close()
	s.quad2 = gocv.NewMat()
	defer s.quad2.Close()
	s.quad3 = gocv.NewMat()
	defer s.quad3.Close()
	s.quad4 = gocv.NewMat()
	defer s.quad4.Close()

	s.quadUp = gocv.NewMat()
	defer s.quadUp.Close()
	s.quadDown = gocv.NewMat()
	defer s.quadDown.Close()
	s.quadTmp = gocv.NewMat()
	defer s.quadTmp.Close()
	s.frameResize = gocv.NewMat()
	defer s.frameResize.Close()

	black.CopyTo(&s.quad1)
	black.CopyTo(&s.quad2)
	black.CopyTo(&s.quad3)
	black.CopyTo(&s.quad4)

	mats := make([]gocv.Mat, 0)

	frameClearBuffer := gocv.NewMat()
	defer frameClearBuffer.Close()

	window := gocv.NewWindow(s.ClientID)
	defer window.Close()
	window.ResizeWindow(s.Width, s.Height)

	var done chan bool

	var decoded string
	var decoded2 string
	// var img image.Image
	var status string
	var code1 string
	var code2 string
	var QRsCount int

	var existPhoto1 bool
	var existPhoto2 bool

	photo1 := gocv.NewMat()
	defer photo1.Close()
	photo2 := gocv.NewMat()
	defer photo2.Close()

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

		frame.CopyTo(&s.quad4)

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
				status, code1, code2 = s.AccessGranted2(&decoded, &decoded2)

				switch status {
				case _ACCESS_GRANTED:
					if *s.DebugMode {
						log.Printf("### ACCESO PERMITIDO ###\n")
					}

					existPhoto1 = s.GetPhotoPerson(code1, &photo1)
					if code2 != "" {
						existPhoto2 = s.GetPhotoVehicle(code2, &photo2)
					}

					green.CopyTo(&s.quad3)
					if existPhoto1 {
						photo1.CopyTo(&s.quad1)
					} else {
						black.CopyTo(&s.quad1)
					}

					if existPhoto2 {
						photo2.CopyTo(&s.quad2)
					} else {
						black.CopyTo(&s.quad2)
					}

					wavGranted.Play()
					//wait = 2000

				case _ACCESS_DENIED:
					if *s.DebugMode {
						log.Printf("### ACCESO DENEGADO ###\n")
					}

					existPhoto1 = s.GetPhotoPerson(code1, &photo1)
					if code2 != "" {
						existPhoto2 = s.GetPhotoVehicle(code2, &photo2)
					}

					if decoded2 != "" {
						fmt.Printf("%s | %-6s | %-25s | %-23s | * ACCESO DENEGADO *\n",
							time.Now().Format("02/01/2006 15:04:05"),
							s.EventName,
							decoded[len(decoded)-12:],
							decoded2[len(decoded2)-12:])
					} else {
						fmt.Printf("%s | %-6s | %-51s | * ACCESO DENEGADO *\n",
							time.Now().Format("02/01/2006 15:04:05"),
							s.EventName,
							decoded)
					}

					red.CopyTo(&s.quad3)
					s.PutTitle(&s.quad3, "ACCESO DENEGADO!", 1)
					s.PutText(&s.quad3, "", 2)
					s.PutText(&s.quad3, "QR invalido, vencido o anulado", 3)

					if existPhoto1 {
						photo1.CopyTo(&s.quad1)
					} else {
						black.CopyTo(&s.quad1)
					}
					if existPhoto2 {
						photo2.CopyTo(&s.quad2)
					} else {
						black.CopyTo(&s.quad2)
					}

					wavDenied.Play()
					//wait = 2000

				case _CONTINUE:

				case _ERROR:
					break
				}

				decoded = ""
				decoded2 = ""
			}
		}

		window.IMShow(*s.GetFrame())

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
			*s.DebugMode = !*s.DebugMode
			log.Printf("MODO DEBUG: %v\n", *s.DebugMode)
		}

		wait = 1
	}
	log.Println("Cámara cerrada")
}

func (s *LectorQR) GetFrame() *gocv.Mat {
	gocv.Hconcat(s.quad1, s.quad2, &s.quadUp)
	gocv.Hconcat(s.quad3, s.quad4, &s.quadDown)
	gocv.Vconcat(s.quadUp, s.quadDown, &s.quadTmp)
	gocv.Resize(s.quadTmp, &s.frameResize, image.Pt(s.Width, s.Height), 0, 0, gocv.InterpolationArea)
	return &s.frameResize
}

func (s *LectorQR) PutTitle(background *gocv.Mat, text string, linea int) *gocv.Mat {
	gocv.PutText(background, text, image.Point{10, 50 * linea}, gocv.FontHersheyTriplex, 1.1, color.RGBA{0, 0, 0, 0}, 2)
	return background
}

func (s *LectorQR) PutTitleWithColor(background *gocv.Mat, text string, linea int, color color.RGBA) *gocv.Mat {
	gocv.PutText(background, text, image.Point{10, 50 * linea}, gocv.FontHersheyTriplex, 1.1, color, 2)
	return background
}

func (s *LectorQR) PutText(background *gocv.Mat, text string, linea int) *gocv.Mat {
	gocv.PutText(background, text, image.Point{10, 35 * linea}, gocv.FontHersheyDuplex, 1, color.RGBA{0, 0, 0, 0}, 2)
	return background
}

func (s *LectorQR) PutTextWithColor(background *gocv.Mat, text string, linea int, color color.RGBA) *gocv.Mat {
	gocv.PutText(background, text, image.Point{10, 35 * linea}, gocv.FontHersheyDuplex, 1, color, 2)
	return background
}

func (s *LectorQR) GetPhotoPerson(code string, photo *gocv.Mat) bool {
	card, err := s.Repo.InfoCard(code)
	if err != nil {
		if *s.DebugMode {
			log.Println("Error obteniendo información de la tarjeta")
		}
		return false
	}

	if card.Photo == "" {
		if *s.DebugMode {
			log.Println("No existe información de la tarjeta de la persona")
		}
		return false
	}
	f := fmt.Sprintf("%s\\%s", s.TmpDir, card.Photo)
	if !CheckFileExists(f) {
		bytes, err := backend.GetPhotoPerson(s.UrlBackend, s.APIKey, card.Photo)
		if err != nil {
			if *s.DebugMode {
				log.Println("Error obteniendo foto del servidor. ", err)
			}
			return false
		}
		err = SaveJPG(bytes, f)
		if err != nil {
			if *s.DebugMode {
				log.Println("Error grabando foto local. ", err)
			}
			return false
		}
	}

	photoTmp := gocv.IMRead(f, gocv.IMReadColor)
	defer photoTmp.Close()
	if photoTmp.Cols() != s.Width || photoTmp.Rows() != s.Height {
		gocv.Resize(photoTmp, &photoTmp, image.Pt(s.Width, s.Height), 0, 0, gocv.InterpolationArea)
	}
	photoTmp.CopyTo(photo)
	return !photoTmp.Empty()
}

func (s *LectorQR) GetPhotoVehicle(code string, photo *gocv.Mat) bool {
	card, err := s.Repo.InfoCard(code)
	if err != nil {
		if *s.DebugMode {
			log.Println("Error obteniendo información de la tarjeta")
		}
		return false
	}
	if card.Photo == "" {
		if *s.DebugMode {
			log.Println("No existe información de la tarjeta del vehículo")
		}
		return false
	}
	f := fmt.Sprintf("%s\\%s", s.TmpDir, card.Photo)
	if !CheckFileExists(f) {
		if os.IsNotExist(err) {
			// File or directory does not exist
			bytes, err := backend.GetPhotoVehicle(s.UrlBackend, s.APIKey, card.Photo)
			if err != nil {
				if *s.DebugMode {
					log.Println("Error obteniendo foto del servidor. ", err)
				}
				return false
			}
			err = SaveJPG(bytes, f)
			if err != nil {
				if *s.DebugMode {
					log.Println("Error grabando foto local. ", err)
				}
				return false
			}
		} else {
			if *s.DebugMode {
				log.Println(err)
			}
			return false
		}
	}
	photoTmp := gocv.IMRead(f, gocv.IMReadColor)
	defer photoTmp.Close()
	if photoTmp.Cols() != s.Width || photoTmp.Rows() != s.Height {
		gocv.Resize(photoTmp, &photoTmp, image.Pt(s.Width, s.Height), 0, 0, gocv.InterpolationArea)
	}
	photoTmp.CopyTo(photo)
	return !photoTmp.Empty()
}

func (s *LectorQR) SaveAccessGranted(code1, code2 string) {

	access := models.Access{
		UUID:       s.Repo.NewUUID(),
		Code1:      code1,
		Code2:      code2,
		AccessDate: time.Now(),
		Zone:       s.Zone,
		Event:      s.EventCode,
	}

	if *s.DebugMode {
		log.Printf("access request: %v\n", access)
	}

	go func(url_backend, APIKey string, access *models.Access) {

		dataAccess, err := backend.SendToServer(url_backend, s.APIKey, *access)
		if err != nil {
			if *s.DebugMode {
				LogError("*** Error enviando movimiento al servidor ***", err, *s.DebugMode)
			}
			err := s.Repo.InsertAccess(access)
			if err != nil {
				LogError("Amacenamiento local: ERROR", err, *s.DebugMode)
			} else {
				if *s.DebugMode {
					log.Println("Almacenamiento en local: OK")
				}

				if code2 != "" {
					fmt.Printf("%s | %-6s | %-25s | %-23s | * SINCRONIZACION PENDIENTE *\n",
						time.Now().Format("02/01/2006 15:04:05"),
						s.EventName,
						code1[len(code1)-12:],
						code2[len(code2)-12:])
				} else {
					fmt.Printf("%s | %-6s | %-51s | * SINCRONIZACION PENDIENTE *\n",
						time.Now().Format("02/01/2006 15:04:05"),
						s.EventName,
						code1)
				}

				s.PutTitle(&s.quad3, "ACCESO PERMITIDO", 1)
				s.PutText(&s.quad3, "", 2)
				s.PutText(&s.quad3, fmt.Sprintf("%s: %s", s.EventName, access.AccessDate.Format("02/01/2006 15:04:05")), 3)
			}

		} else {
			if *s.DebugMode {
				backend.PrintData(dataAccess)
			}
			fmt.Printf("%s | %-6s | %-25s | %s %s | %-8s | %s\n",
				access.AccessDate.Format("02/01/2006 15:04:05"),
				s.EventName,
				dataAccess.PersonName,
				dataAccess.DocumentType,
				dataAccess.DocumentNumber,
				GetCardColor(dataAccess.Color),
				dataAccess.LicensePlate)

			s.PutTitle(&s.quad3, "ACCESO PERMITIDO", 1)
			s.PutText(&s.quad3, "", 2)
			s.PutText(&s.quad3, fmt.Sprintf("%s: %s", s.EventName, access.AccessDate.Format("02/01/2006 15:04:05")), 3)
			s.PutText(&s.quad3, fmt.Sprintf("%s %s", dataAccess.DocumentType, dataAccess.DocumentNumber), 4)
			s.PutText(&s.quad3, dataAccess.PersonName, 5)
			//s.PutText(&s.quad3, "DIRECCION PROVINCIAL DE PUERTOS", 6)
			if dataAccess.Eventual {
				s.PutText(&s.quad3, fmt.Sprintf("EVENTUAL - PNA: %s", dataAccess.PNA), 7)
			} else {
				s.PutText(&s.quad3, fmt.Sprintf("PERMANENTE - PNA: %s", dataAccess.PNA), 7)
			}
			s.PutText(&s.quad3, fmt.Sprintf("VIGENCIA: %s %s", dataAccess.DateFrom.Format("02/01/2006"), dataAccess.DateTo.Format("02/01/2006")), 8)
			s.PutText(&s.quad3, fmt.Sprintf("COLOR TARJ. PERSONA: %s", GetCardColor(dataAccess.Color)), 9)
			if dataAccess.LicensePlate != "" {
				s.PutText(&s.quad3, fmt.Sprintf("VEHICULO: %s", dataAccess.LicensePlate), 11)
			}
		}
	}(s.UrlBackend, s.APIKey, &access)
}

func (s *LectorQR) AccessGranted2(decoded *string, decoded2 *string) (status, code1, code2 string) {

	if s.Mode == 1 { //Only person
		if *decoded != QRPrev1 && IsPerson(*decoded) {
			QRPrev1 = *decoded
			if *s.DebugMode {
				log.Printf("QR [%s] es %s\n", *decoded, PersonOrVehicle(*decoded))
			}
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

		if *s.DebugMode {
			log.Printf("QR1 [%s] es %s, QR2 [%s] es %s\n", *decoded, PersonOrVehicle(*decoded), *decoded2, PersonOrVehicle(*decoded2))
		}

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
