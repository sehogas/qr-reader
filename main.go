package main

import (
	"flag"
	"fmt"
	"log"

	"gocv.io/x/gocv"
)

func lectorQR(camara *gocv.VideoCapture) {
	var decoded string
	var dataOld string

	frame := gocv.NewMat()
	defer frame.Close()

	pts := gocv.NewMat()
	defer pts.Close()

	straight_qrcode := gocv.NewMat()
	defer straight_qrcode.Close()

	window := gocv.NewWindow("Detector QR")
	defer window.Close()

	qrCodeDetector := gocv.NewQRCodeDetector()
	defer qrCodeDetector.Close()

	log.Println("Iniciando lectura de la cámara...")

	for {
		if ok := camara.Read(&frame); !ok {
			log.Printf("No se pudo leer frame de cámara")
			return
		}
		if frame.Empty() {
			continue
		}
		camara.Read(&frame)

		decoded = qrCodeDetector.DetectAndDecode(frame, &pts, &straight_qrcode)
		if decoded != "" {
			if decoded != dataOld {
				dataOld = decoded
				log.Println(decoded)
			}
		}

		window.IMShow(frame)
		window.WaitKey(1)
	}

}

func main() {
	var deviceID int
	var fromFile string

	flag.IntVar(&deviceID, "device-id", 0, "integer value, webcam device ID")
	flag.StringVar(&fromFile, "from-file", "", "string value, url: rtsp://user:pass@host:port/stream1, rtsp://dahua:admin@192.168.88.108:554/cam/realmonitor?channel=1&subtype=1")
	flag.Parse()

	jobs := make(chan int, 4)
	done := make(chan bool)

	var camara *gocv.VideoCapture
	var err error

	if fromFile == "" {
		camara, err = gocv.VideoCaptureDevice(deviceID)
	} else {
		camara, err = gocv.VideoCaptureFile(fromFile)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	defer camara.Close()

	go lectorQR(camara)

	close(jobs)
	fmt.Println("Todos los trabajos enviados")

	<-done
	log.Println("Fin")
}
