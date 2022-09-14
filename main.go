package main

import (
	"flag"
	"fmt"

	"gocv.io/x/gocv"
)

func main() {
	var deviceID int
	var fromFile string
	var dataOld string

	flag.IntVar(&deviceID, "device-id", 0, "integer value, webcam device ID")
	flag.StringVar(&fromFile, "from-file", "", "string value, url: rtsp://user:pass@host:port/stream1, rtsp://dahua:admin@192.168.88.108:554/cam/realmonitor?channel=1&subtype=1")
	flag.Parse()

	var webcam *gocv.VideoCapture
	var err error

	if fromFile == "" {
		webcam, err = gocv.VideoCaptureDevice(deviceID)
	} else {
		webcam, err = gocv.VideoCaptureFile(fromFile)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	frame := gocv.NewMat()
	defer frame.Close()

	pts := gocv.NewMat()
	defer pts.Close()

	straight_qrcode := gocv.NewMat()
	defer straight_qrcode.Close()

	window := gocv.NewWindow("Detector QR")
	defer window.Close()

	//	window2 := gocv.NewWindow("QR Detectado")
	//	defer window2.Close()

	qrCodeDetector := gocv.NewQRCodeDetector()
	defer qrCodeDetector.Close()

	//green := color.RGBA{0, 255, 0, 0} // green

	var decoded string

	if fromFile != "" {
		fmt.Printf("start reading camera from file: %v\n", fromFile)
	} else {
		fmt.Printf("start reading camera device: %v\n", deviceID)
	}
	for {
		if ok := webcam.Read(&frame); !ok {
			fmt.Printf("No se pudo leer el dispositivo %d\n", 0)
			return
		}
		if frame.Empty() {
			continue
		}
		webcam.Read(&frame)

		decoded = qrCodeDetector.DetectAndDecode(frame, &pts, &straight_qrcode)
		if decoded != "" {
			if decoded != dataOld {
				dataOld = decoded
				fmt.Println(decoded)
			}
		}

		window.IMShow(frame)
		window.WaitKey(1)
	}

}
