package util

import (
	"log"
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CheckConfig(m map[string]string) {

	if m == nil {
		log.Fatal("Archivo de configuración de entorno inválido")
	}

	v, e := m["CLIENT_ID"]
	if !e || v == "" {
		log.Fatal("Se requiere el parámetro CLIENT_ID")
	}

	v, e = m["ZONE_ID"]
	if !e || v == "" {
		log.Fatal("Se requiere el parámetro ZONE_ID")
	}

	v, e = m["EVENT_ID"]
	if !e || v == "" {
		log.Fatal("Se requiere el parámetro EVENT_ID")
	}
	if v == "I" || v == "O" {
		if v == "I" {
			m["EVENT_ID"] = "E"
		}
		if v == "O" {
			m["EVENT_ID"] = "S"
		}
	} else {
		log.Fatal("El parámetro EVENT_ID es inválido")
	}

	v, e = m["MODE"]
	if !e || v == "" {
		log.Fatal("Se requiere el parámetro MODE")
	}
	if !(v == "1" || v == "2") {
		log.Fatal("El parámetro MODE es inválido")
	}

	v, e = m["URL_BACKEND"]
	if !e || v == "" {
		log.Fatal("Se requiere el parámetro URL_BACKEND")
	}

	v, e = m["API_KEY"]
	if !e || v == "" {
		log.Fatal("Se requiere el parámetro API_KEY")
	}

	v, e = m["DEVICE_ID"]
	if !e || v == "" {
		log.Fatal("Se requiere el parámetro DEVICE_ID")
	}

	_, e = m["RTSP"]
	if !e {
		log.Fatal("Se requiere el parámetro RTSP aunque esté vacío")
	}

	v, e = m["DB"]
	if !e || v == "" {
		log.Fatal("Se requiere el parámetro DB")
	}

	v, e = m["FILE_WAV_GRANTED"]
	if !e || v == "" {
		log.Fatal("Se requiere el parámetro FILE_WAV_GRANTED")
	}
	if !FileExists(v) {
		log.Fatal("El archivo configurado en el parámetro FILE_WAV_GRANTED no existe")
	}

	v, e = m["FILE_WAV_DENIED"]
	if !e || v == "" {
		log.Fatal("Se requiere el parámetro FILE_WAV_DENIED")
	}
	if !FileExists(v) {
		log.Fatal("El archivo configurado en el parámetro FILE_WAV_DENIED no existe")
	}

	v, e = m["FILE_DETECT_PROTO_TXT"]
	if !e || v == "" {
		log.Fatal("Se requiere el parámetro FILE_DETECT_PROTO_TXT")
	}
	if !FileExists(v) {
		log.Fatal("El archivo configurado en el parámetro FILE_DETECT_PROTO_TXT no existe")
	}

	v, e = m["FILE_DETECT_CAFFE"]
	if !e || v == "" {
		log.Fatal("Se requiere el parámetro FILE_DETECT_CAFFE")
	}
	if !FileExists(v) {
		log.Fatal("El archivo configurado en el parámetro FILE_DETECT_CAFFE no existe")
	}

	v, e = m["FILE_SUPER_PROTO_TXT"]
	if !e || v == "" {
		log.Fatal("Se requiere el parámetro FILE_SUPER_PROTO_TXT")
	}
	if !FileExists(v) {
		log.Fatal("El archivo configurado en el parámetro FILE_SUPER_PROTO_TXT no existe")
	}

	v, e = m["FILE_SUPER_CAFFE"]
	if !e || v == "" {
		log.Fatal("Se requiere el parámetro FILE_SUPER_CAFFE")
	}
	if !FileExists(v) {
		log.Fatal("El archivo configurado en el parámetro FILE_SUPER_CAFFE no existe")
	}

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
	return _UNKNOWN
}

func LogError(text string, err error, modeDebug bool) {
	log.Println(text)
	if modeDebug {
		log.Printf("DEBUG: %s\n", err.Error())
	}
}
