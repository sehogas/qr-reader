package util

import (
	"log"
	"time"

	"github.com/sehogas/qr-reader/models"
	"gocv.io/x/gocv"
)

type LectorQR struct {
	DeviceID int
	FromFile string
	Repo     *Repository
}

func NewLectorQR(deviceID int, fromFile string, repo *Repository) *LectorQR {
	return &LectorQR{
		DeviceID: deviceID,
		FromFile: fromFile,
		Repo:     repo,
	}
}

func (s *LectorQR) Start() {
	var camera *gocv.VideoCapture
	var err error

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

	pts := gocv.NewMat()
	defer pts.Close()

	straight_qrcode := gocv.NewMat()
	defer straight_qrcode.Close()

	//window := gocv.NewWindow("QR detector")
	//defer window.Close()

	qrCodeDetector := gocv.NewQRCodeDetector()
	defer qrCodeDetector.Close()

	var decoded string
	var previous string

	log.Println("Reading camera...")

	for {
		if ok := camera.Read(&frame); !ok {
			log.Println("Could not read the camera")
			return
		}

		if frame.Empty() {
			continue
		}

		camera.Read(&frame)

		decoded = qrCodeDetector.DetectAndDecode(frame, &pts, &straight_qrcode)
		if decoded != "" {
			if decoded != previous {
				previous = decoded
				log.Printf("Read QR [%s]\n", decoded)
				if s.Repo.ValidCard(decoded) {
					s.Repo.InsertAccess(models.Access{Code: decoded,
						AccessDate:   time.Now(),
						Zone:         "AP",
						Event:        "E",
						Synchronized: false})
					log.Println("Access granted")
				} else {
					log.Println("Access denied!")
				}
			}
		}

		//window.IMShow(frame)
		//window.WaitKey(10)
	}

}
